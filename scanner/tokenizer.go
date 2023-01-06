/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package scanner

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Element = *element

type element struct {
	flagged bool
	token   string
	tags    []string
	text    string
	location
}

func NewToken(token string, tags []string, location Location, flagged bool) Element {
	return &element{
		flagged:  flagged,
		token:    token,
		tags:     tags,
		location: location,
	}
}

func NewText(text string, location Location) Element {
	return &element{
		text:     text,
		location: location,
	}
}

func (e *element) IsText() bool {
	return e.token == ""
}

func (e *element) IsFlagged() bool {
	return e.flagged
}

func (e *element) IsToken() bool {
	return e.token != ""
}

func (e *element) Token() string {
	return e.token
}

func (e *element) Tags() []string {
	return e.tags
}

func (e *element) HasTags() bool {
	return len(e.tags) > 0
}

func (e *element) Tag(desc string) (string, error) {
	tgt := " for " + desc
	if desc == "" || desc == "tag" {
		tgt = ""
		desc = "tag"
	}
	if len(e.tags) == 0 {
		return "", e.Errorf("%s required", desc)
	}
	if len(e.tags) != 1 {
		return "", e.Errorf("found multiple tags%s", tgt)
	}
	if e.tags[0] == "" {
		return "", e.Errorf("non-empty tag required%s", desc)
	}
	return e.tags[0], nil
}

func (e *element) OptionalTag(desc string) (string, error) {
	if desc == "" || desc == "tag" {
		desc = ""
	} else {
		desc = " for " + desc
	}
	if len(e.tags) == 0 {
		return "", nil
	}
	if len(e.tags) != 1 {
		return "", e.Errorf("found multiple tags%s", desc)
	}
	if e.tags[0] == "" {
		return "", e.Errorf("non-empty tag required%s", desc)
	}
	return e.tags[0], nil
}

func (e *element) Text() string {
	return e.text
}

func (e *element) Append(s string) {
	if !e.IsToken() {
		e.text += s
	}
}

func (e *element) Location() Location {
	return e.location
}

func (e *element) String() string {
	if e.token != "" {
		tags := strings.Join(e.tags, "\", \"")
		if tags != "" {
			tags = "\"" + tags + "\""
		}
		return fmt.Sprintf("%s[%s]", e.token, tags)
	}
	return fmt.Sprintf(": %s", strings.ReplaceAll(e.text, "\n", "\n> "))
}

type Tokenizer = *tokenizer

type tokenizer struct {
	scanner Scanner
	stack   []Element
}

func NewTokenizer(source string, r io.Reader) Tokenizer {
	return &tokenizer{
		scanner: NewScanner(source, r),
	}
}

func (p *tokenizer) Push(e Element) {
	p.stack = append(p.stack, e)
}

func (p *tokenizer) requireNext() (rune, error) {
	r, err := p.scanner.Next()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return r, p.scanner.Errorf("unexpected EOF")
		}
		return r, p.scanner.Errorf("%s", err.Error())
	}
	return r, nil
}

func (p *tokenizer) NextElement() (Element, error) {

	if len(p.stack) > 0 {
		e := p.stack[len(p.stack)-1]
		p.stack = p.stack[:len(p.stack)-1]
		return e, nil
	}
	var next string

	text, loc, ok, err := p.parseTokenStart()

	if err != nil {
		if errors.Is(err, io.EOF) {
			if text != "" {
				return &element{text: text, location: loc}, nil
			}
			return nil, nil
		} else {
			return nil, err
		}
	}

	for !ok {
		next, loc, ok, err = p.parseTokenStart()
		text += next
		if errors.Is(err, io.EOF) {
			if text != "" {
				return &element{text: text, location: loc}, nil
			}
			return nil, nil
		}
	}
	if text != "" {
		return &element{
			text:     text,
			location: loc,
		}, nil
	}

	flagged := false
	p.scanner.Consume("{{")

	r := ' '
	for r == ' ' {
		r, err = p.requireNext()
		if err != nil {
			return nil, err
		}
	}
	if r == '*' {
		flagged = true
		r, err = p.requireNext()
		if err != nil {
			return nil, err
		}
	}
	if !unicode.IsLetter(r) {
		return nil, p.scanner.ErrorForPreviousf("invalid character %q in statement name", string(r))
	}
	token := string(r)
	for !p.scanner.Match("}}") {
		r, err = p.requireNext()
		if err != nil {
			return nil, err
		}
		if unicode.IsDigit(r) || unicode.IsLetter(r) {
			token += string(r)
		} else {
			if string(r) == " " {
				break
			} else {
				return nil, p.scanner.ErrorForPreviousf("invalid character %q in statement name", string(r))
			}
		}
	}

	var tags []string
	active := false
	quoted := false
	tag := ""
	masked := false

	for masked || quoted || !p.scanner.Consume("}}") {
		c, err := p.scanner.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				if quoted {
					return nil, p.scanner.Errorf("unexpected EOF in quoted tag")
				}
				if masked {
					return nil, p.scanner.Errorf("unexpected EOF after escape character")
				}
				return nil, p.scanner.Errorf("unexpected EOF in token")
			}
			return nil, p.scanner.Error(err)
		}
		if masked {
			masked = false
		} else {
			switch c {
			case '\\':
				masked = true
				continue
			case ' ':
				if !quoted {
					if active {
						tags = append(tags, tag)
						tag = ""
						active = false
					}
					continue
				}
			case '"':
				quoted = !quoted
				continue
			}
		}
		active = true
		tag += string(c)
	}
	if active {
		tags = append(tags, tag)
	}

	if p.scanner.Match("\\\n") {
		p.scanner.Consume("\\")
	} else {
		if p.scanner.Match("\n") && loc.column <= 1 {
			p.scanner.Consume("\n")
		}
	}
	return NewToken(token, tags, loc, flagged), nil
}

func (p *tokenizer) parseTokenStart() (string, location, bool, error) {
	mask := 0
	text := ""
	loc := p.scanner.Location()
	p.scanner.SkipComment()
	for p.scanner.Consume("\\") {
		mask++
	}
	p.scanner.SkipComment()
	if p.scanner.Match("{{") {
		for mask > 1 {
			text += "\\"
			mask -= 2
		}
		if mask != 0 {
			text += "{{"
			p.scanner.Consume("{{")
		} else {
			return text, p.scanner.Location(), true, nil
		}
	} else {
		for mask > 0 {
			text += "\\"
			mask--
		}
	}
	if text == "" {
		p.scanner.SkipComment()
		r, err := p.scanner.Next()
		if err != nil {
			return text, loc, false, err
		}
		text = string(r)
	}
	return text, loc, false, nil
}
