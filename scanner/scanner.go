/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package scanner

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"unicode/utf8"
)

type Scanner = *scanner

type scanner struct {
	reader    *bufio.Reader
	lookAhead string
	location
}

type Located interface {
	Source() string
	Location() Location
	Errorf(msg string, args ...interface{}) error
	Error(err error) error
}

type Location = location

type location struct {
	source string
	line   int
	column int
}

var _ Located = Location{}

func NewLocation(source string, line int) Location {
	return Location{source: source, line: line, column: 1}
}

func (l Location) Location() Location {
	return l
}

func (l location) Source() string {
	return l.source
}

func (l location) Line() int {
	return l.line
}

func (l location) Column() int {
	return l.column
}

func (l location) Errorf(msg string, args ...interface{}) error {
	return fmt.Errorf("%s: %s", l.String(), fmt.Sprintf(msg, args...))
}

func (l location) Error(err error) error {
	return fmt.Errorf("%s: %s", l.String(), err.Error())
}

func (l location) ErrorForPreviousf(msg string, args ...interface{}) error {
	l.column--
	return l.Errorf(msg, args...)
}

func (l location) String() string {
	if l.line == 0 {
		return l.source
	}
	if l.source == "" {
		return fmt.Sprintf("line %d, column %d", l.line, l.column)
	}
	return fmt.Sprintf("%s: line %d, column %d", l.source, l.line, l.column)
}

func (l location) SkipLine() Location {
	l.line++
	l.column = 1
	return l
}

func NewScanner(source string, r io.Reader) Scanner {
	return &scanner{
		location: NewLocation(source, 1),
		reader:   bufio.NewReader(r),
	}
}

func (s *scanner) inc(r rune) rune {
	if r == '\n' {
		s.line++
		s.column = 1
	} else {
		s.column++
	}
	return r
}

func (s *scanner) SkipComment() {
	firstcol := s.column
	cnt := 0
	lastcol := s.column
	for s.Match("/#") {
		for {
			lastcol = s.column
			r, err := s.next()
			if err != nil {
				return
			}
			if r == '\n' {
				cnt++
				break
			}
		}
	}
	if cnt > 0 && firstcol > 1 {
		s.lookAhead = "\n" + s.lookAhead
		s.column = lastcol
		s.line--
	}
}

func (s *scanner) Next() (rune, error) {
	return s.next()
}

func (s *scanner) next() (rune, error) {
	if s.lookAhead != "" {
		r, size := utf8.DecodeRuneInString(s.lookAhead)
		s.lookAhead = s.lookAhead[size:]
		return s.inc(r), nil
	}
	r, _, err := s.reader.ReadRune()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return 0, err
		}
		return 0, s.Errorf("%s", err.Error())
	}
	return s.inc(r), nil
}

func (s *scanner) Match(n string) bool {
	for len(s.lookAhead) < len(n) {
		r, _, err := s.reader.ReadRune()
		if err != nil {
			return false
		}
		s.lookAhead += string(r)
	}
	return s.lookAhead[:len(n)] == n
}

func (s *scanner) Consume(n string) bool {
	if s.Match(n) {
		s.lookAhead = s.lookAhead[len(n):]
		if n == "\n" {
			s.line++
			s.column = 1
		} else {
			s.column += len(n)
		}
		return true
	}
	return false
}

func (s *scanner) Location() Location {
	l := s.location
	return l
}
