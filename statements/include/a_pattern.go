/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package include

import (
	"fmt"
	"regexp"

	"github.com/mandelsoft/mdgen/scanner"
)

func init() {
	scanner.Keywords.Register("pattern", true)

}

var extractExpPat = regexp.MustCompile("^([a-zA-Z -]+)$")

func ParsePattern(p scanner.Parser, n *ContentHandler, e scanner.Element) (scanner.Element, error) {

	key, err := e.Tag("extraction key")
	if err != nil {
		return nil, err
	}
	if n.extract != nil {
		return nil, e.Errorf("range specification already set")
	}

	if !extractExpPat.MatchString(key) {
		return nil, e.Errorf("invalid range key %q", key)
	}
	n.extract = &PatternExtractor{key}
	return p.NextElement()
}

////////////////////////////////////////////////////////////////////////////////

type PatternExtractor struct {
	pattern string
}

func (i *PatternExtractor) Extract(data []byte) ([]byte, error) {
	_, start, err := i.match(data, "begin")
	if err != nil {
		return nil, err
	}
	end, _, err := i.match(data, "end")
	if err != nil {
		return nil, err
	}

	return data[start:end], nil
}

func (i *PatternExtractor) match(data []byte, key string) (int, int, error) {
	exp := regexp.MustCompile(fmt.Sprintf("(?m)^.*--- %s %s ---.*$", key, regexp.QuoteMeta(i.pattern)))

	matches := exp.FindAllIndex(data, -1)
	if len(matches) == 0 {
		return -1, -1, fmt.Errorf("%s pattern (%s) not found", key, i.pattern)
	}
	if len(matches) != 1 {
		return -1, -1, fmt.Errorf("%s pattern (%s) is not unique", key, i.pattern)
	}

	start := matches[0][0]
	if start > 0 && data[start-1] == '\n' {
		start--
	}
	if start > 0 && data[start-1] == '\r' {
		start--
	}

	end := matches[0][1]
	if len(data) > end && data[end] == '\r' {
		end++
	}
	if len(data) > end && data[end] == '\n' {
		end++
	}
	return start, end, nil
}
