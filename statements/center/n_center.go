/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sectionref

import (
	"fmt"

	"github.com/mandelsoft/mdgen/scanner"
)

func init() {
	scanner.Tokens.RegisterStatement(NewStatement(), true)
}

type Statement struct {
	scanner.BracketedStatement[Node]
}

func NewStatement() *Statement {
	return &Statement{scanner.NewBracketedStatement[Node]("center", true)}
}

func (s *Statement) Start(p scanner.Parser, e scanner.Element) (scanner.Element, error) {
	if e.HasTags() {
		return nil, e.Errorf("no arguments possible")
	}

	n := NewNode(p, e.Location())
	p.State = p.State.Sub(n)
	return p.NextElement()
}

////////////////////////////////////////////////////////////////////////////////

type NodeContext struct {
	scanner.NodeContextBase[Node]
}

func NewNodeContext(n Node, ctx scanner.ResolutionContext) (*NodeContext, error) {
	var err error
	c := &NodeContext{
		NodeContextBase: scanner.NewNodeContextBase(n, ctx),
	}
	return c, err
}

type Node = *node

type node struct {
	scanner.NodeContainerBase
}

func NewNode(p scanner.Parser, location scanner.Location) Node {
	return &node{
		NodeContainerBase: scanner.NewContainerBase("center", p.Document(), location, p.State.Container),
	}
}

func (n *node) Print(gap string) {
	fmt.Printf("%sCENTER:\n", gap)
	fmt.Printf("%s  nodes:\n", gap)
	n.NodeContainerBase.Print(gap + "    ")
}

func (n *node) Register(ctx scanner.ResolutionContext) error {

	nctx, err := NewNodeContext(n, ctx)
	if err != nil {
		return err
	}
	ctx.SetNodeContext(n, nctx)
	return n.NodeSequence.Register(ctx)
}

func (n *node) Emit(ctx scanner.ResolutionContext) error {
	if len(n.NodeSequence.GetNodes()) > 0 {
		w := ctx.Writer()
		fmt.Fprintf(w, "<div align=\"center\">\n\n")
		n.NodeSequence.Emit(ctx)
		fmt.Fprintf(w, "</div>\n")
	}

	return nil
}
