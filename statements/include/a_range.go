/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package include

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/mandelsoft/mdgen/scanner"
)

func init() {
	scanner.Keywords.Register("range", true)

}

var extractExpNum = regexp.MustCompile("^([0-9]+)?(?:(:)([0-9]+)?)?$")

func ParseRange(p scanner.Parser, n *ContentHandler, e scanner.Element) (scanner.Element, error) {

	key, err := e.Tag("extraction range")
	if err != nil {
		return nil, err
	}
	m := extractExpNum.FindSubmatch([]byte(key))
	if m == nil {
		return nil, e.Errorf("invalid range specification (%s): expected [<num>][:[<num>]]", key)
	}
	start := int64(0)
	end := int64(0)
	if m[1] != nil {
		start, err = strconv.ParseInt(string(m[1]), 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid start line: %w", err)
		}
		end = start
	}
	if m[2] != nil {
		if m[3] != nil {
			end, err = strconv.ParseInt(string(m[3]), 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid start line: %w", err)
			}
		} else {
			end = 0
		}
	}
	n.extract = &NumExtractor{
		start: int(start),
		end:   int(end),
	}
	return p.NextElement()
}

////////////////////////////////////////////////////////////////////////////////

type NumExtractor struct {
	start, end int
}

func (i *NumExtractor) Extract(data []byte) ([]byte, error) {
	lines := strings.Split(string(data), "\n")
	start := 0
	if i.start > 0 {
		start = i.start - 1
	}
	if start >= len(lines) {
		return nil, fmt.Errorf("start line %d after end of data (%d lines)", start, len(lines))
	}
	end := len(lines)
	if i.end > 0 {
		end = i.end
	}
	if end > len(lines) {
		return nil, fmt.Errorf("end line %d after end of file (%d lines", end, len(lines))
	}
	return []byte(strings.Join(lines[start:end], "\n")), nil
}
