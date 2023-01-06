/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package block

import (
	"strings"

	"github.com/mandelsoft/mdgen/scanner"
)

func ParseParam(p scanner.Parser, n *blocknode, e scanner.Element) (scanner.Element, error) {
	name, err := e.Tag("parameter name")
	if err != nil {
		return nil, err
	}
	defaulted := e.IsFlagged()
	if name == "" {
		return nil, e.Errorf("non-empty tag required for parameter")
	}

	var names []string
	for _, tag := range strings.Split(name, ",") {
		tag = strings.TrimSpace(tag)
		if strings.HasPrefix(tag, "#") {
			return nil, e.Errorf("invalid parameter name %s", tag)
		}
		if _, ok := n.params[tag]; ok {
			return nil, e.Errorf("parameter %s already declared", tag)
		}
		names = append(names, tag)
		n.params[tag] = nil
	}

	if !defaulted {
		return p.NextElement()
	}

	e, seq, err := scanner.ParseSequence(p, e)
	if err != nil {
		return nil, err
	}
	for _, name := range names {
		n.params[name] = seq
	}
	return e, nil
}
