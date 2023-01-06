/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package termdef

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/mandelsoft/mdgen/scanner"
	utils2 "github.com/mandelsoft/mdgen/utils"
)

type TermRef struct {
	tag    string
	plural bool
	upper  bool
}

func (t TermRef) Evaluate(ctx scanner.ResolutionContext) (TermRef, bool, error) {
	tag, explicit, err := scanner.EvaluateTag(ctx, t.tag)
	t.tag = tag
	return t, explicit, err
}

func (t *TermRef) Tag() string {
	return t.tag
}

func (t *TermRef) Mode() string {
	mode := "singular"
	if t.plural {
		mode = "plural"
	}
	if t.upper {
		return mode + ",upper"
	}
	return mode
}

func NewTerm(ref TermRef) *Term {
	return &Term{
		TermRef: ref,
	}
}

type Term struct {
	resolved *TermDefNodeContext
	TermRef
}

func (t *Term) IsFormatted() bool {
	return t.resolved.format != ""
}

func (t *Term) Tag() string {
	return t.tag
}

func (t *Term) Singular() string {
	return t.resolved.singular
}

func (t *Term) Plural() string {
	return t.resolved.plural
}

func (t *Term) Format() string {
	f := t.resolved.format
	return f + t.Get() + reverse(f)
}

func (t *Term) FormatSingular() string {
	f := t.resolved.format
	return f + t.Singular() + reverse(f)
}

func (t *Term) Get() string {
	term := t.resolved.singular
	if t.plural {
		term = t.resolved.plural
	}
	if t.upper {
		r, i := utf8.DecodeRuneInString(term)
		term = string(unicode.ToTitle(r)) + term[i:]
	}
	return term
}

func (t *Term) Resolve(ctx scanner.ResolutionContext) error {
	rctx := ctx.LookupTag(GT_TERM, t.tag)
	if rctx == nil {
		return fmt.Errorf("unknown term %q", t.tag)
	}
	t.resolved = rctx.(*TermDefNodeContext)
	return nil
}

func (t *Term) GetLink() utils2.Link {
	return t.resolved.GetLink()
}

////////////////////////////////////////////////////////////////////////////////

func MapTermTag(t string) TermRef {
	term := TermRef{}
	term.upper = false

	if term.plural = strings.HasPrefix(t, "*"); term.plural {
		t = t[1:]
	}
	if term.upper = strings.HasPrefix(t, "^"); term.upper {
		t = t[1:]
	} else {
		r, _ := utf8.DecodeRuneInString(t)
		term.upper = unicode.IsUpper(r)
	}
	term.tag = strings.ToLower(t)
	return term
}

func reverse(str string) string {
	result := ""
	for _, v := range str {
		result = string(v) + result
	}
	return result
}
