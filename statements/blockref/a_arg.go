/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package blockref

import (
	"github.com/mandelsoft/mdgen/scanner"
)

func init() {
	scanner.Keywords.Register("arg", true)
	scanner.Keywords.Register("endarg", true)

}

func ParseArg(p scanner.Parser, b BlockRefNode, e scanner.Element) (scanner.Element, error) {

	name, err := e.Tag("argument name")
	if err != nil {
		return nil, err
	}
	if b.args[name] != nil {
		return nil, e.Errorf("argument %s already defined", name)
	}

	e, seq, err := scanner.ParseSequence(p, e)
	if err != nil {
		return nil, err
	}
	b.args[name] = seq
	return e, nil
}
