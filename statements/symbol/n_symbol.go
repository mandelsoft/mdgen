/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package include

import (
	"fmt"
	"strings"

	"github.com/mandelsoft/mdgen/scanner"
)

func init() {
	scanner.Tokens.RegisterStatement(NewStatement("br", "</br>"))
	scanner.Tokens.RegisterStatement(NewStatement("cs", "/#"))
}

type Statement struct {
	scanner.StatementBase
	symbol string
}

func NewStatement(name string, symbol string) scanner.Statement {
	return &Statement{scanner.NewStatementBase(name), symbol}
}

func (s *Statement) Start(p scanner.Parser, e scanner.Element) (scanner.Element, error) {
	if e.HasTags() {
		return nil, e.Errorf("no arguments possible")
	}
	p.State.Container.AddNode(NewNode(p.Document(), e.Location(), s.Name(), s.symbol))
	return p.NextElement()
}

////////////////////////////////////////////////////////////////////////////////

type Node = *node

type node struct {
	scanner.NodeBase
	name   string
	symbol string
}

func NewNode(d scanner.Document, location scanner.Location, name, symbol string) Node {
	return &node{
		NodeBase: scanner.NewNodeBase(d, location),
		name:     name,
		symbol:   symbol,
	}
}

func (n *node) Print(gap string) {
	fmt.Printf("%s%s %q\n", gap, strings.ToUpper(n.name), n.symbol)
}

func (n *node) Emit(ctx scanner.ResolutionContext) error {
	fmt.Fprintf(ctx.Writer(), n.symbol)
	return nil
}
