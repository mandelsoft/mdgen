/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package include

import (
	"fmt"

	"github.com/mandelsoft/mdgen/scanner"
	"github.com/mandelsoft/mdgen/utils"
)

func init() {
	scanner.Tokens.RegisterStatement(NewStatement())
}

type Statement struct {
	scanner.StatementBase
}

func NewStatement() scanner.Statement {
	return &Statement{scanner.NewStatementBase("pagehistory")}
}

func (s *Statement) Start(p scanner.Parser, e scanner.Element) (scanner.Element, error) {
	if e.HasTags() {
		return nil, e.Errorf("no arguments possible")
	}
	p.State.Container.AddNode(NewNode(p.Document(), e.Location()))
	return p.NextElement()
}

////////////////////////////////////////////////////////////////////////////////

type Node = *node

type node struct {
	scanner.NodeBase
}

func NewNode(d scanner.Document, location scanner.Location) Node {
	return &node{
		NodeBase: scanner.NewNodeBase(d, location),
	}
}

func (n *node) Print(gap string) {
	fmt.Printf("%sPAGEHISTORY\n", gap)
}

func (n *node) Emit(ctx scanner.ResolutionContext) error {

	p := ctx.GetParentDocument()
	if p == nil {
		return nil
	}
	w := ctx.Writer()
	first := true
	for p != nil {
		if !first {
			fmt.Fprintf(w, "&nbsp;&#10230;&nbsp;")
		}
		first = false
		l := utils.NewLink(p.GetRefPath(), "")
		link, err := ctx.DetermineLink(l)
		if err != nil {
			return err
		}

		info := ctx.GetLinkInfo(l)
		fmt.Fprintf(w, "[%s](%s)", *info.Title(), link)
		p = p.GetParentDocument()
	}
	fmt.Fprintf(w, "\n\n---\n\n")
	return nil
}
