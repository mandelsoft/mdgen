/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package include

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mandelsoft/mdgen/scanner"
)

func init() {
	scanner.Keywords.Register("filter", true)

}

func ParseFilter(p scanner.Parser, n *ContentHandler, e scanner.Element) (scanner.Element, error) {
	key, err := e.Tag("filter expression")
	if err != nil {
		return nil, err
	}
	if n.filter != nil {
		return nil, e.Errorf("filter already set")
	}

	m, err := regexp.Compile(key)
	if err != nil {
		return nil, e.Errorf("invalid filter key (%s): %s", key, err)
	}

	n.filter = &RegExpFilter{
		pattern: m,
	}
	return p.NextElement()
}

////////////////////////////////////////////////////////////////////////////////

type RegExpFilter struct {
	pattern *regexp.Regexp
}

func (i *RegExpFilter) Filter(data []byte) ([]byte, error) {
	if i == nil || i.pattern == nil {
		return data, nil
	}
	sep := ""
	if strings.HasPrefix(i.pattern.String(), "(?m)") {
		sep = "\n"
	}
	matches := i.pattern.FindAllSubmatch(data, -1)
	var result []byte
	for _, m := range matches {
		if len(m) != 2 {
			return nil, fmt.Errorf("regular expression must contain one matching group")
		}
		result = append(result, m[1]...)
		result = append(result, []byte(sep)...)
	}
	return result, nil
}
