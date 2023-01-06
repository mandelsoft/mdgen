/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package section

import (
	"fmt"
	"html"
	"strings"

	"github.com/mandelsoft/mdgen/scanner"
)

func init() {
	scanner.Tokens.RegisterStatement(NewStatement())
}

type Statement struct {
	scanner.BracketedStatement[EscapeNode]
}

func NewStatement() *Statement {
	return &Statement{scanner.NewBracketedStatement[EscapeNode]("escape", true)}
}

func (s *Statement) Start(p scanner.Parser, e scanner.Element) (scanner.Element, error) {
	if e.HasTags() {
		return nil, e.Errorf("no arguments expected")
	}

	n := NewEscapeNode(p, e.Location())
	p.State = p.State.Sub(n)
	return p.NextElement()
}

////////////////////////////////////////////////////////////////////////////////

type nodeContextBase = scanner.NodeContextBase[*escapenode]

type EscapeNodeContext struct {
	nodeContextBase
	text *string
}

func NewEscapeContext(n *escapenode, ctx scanner.ResolutionContext) *EscapeNodeContext {
	return &EscapeNodeContext{
		nodeContextBase: scanner.NewNodeContextBase(n, ctx),
	}
}

type EscapeNode = *escapenode

type escapenode struct {
	scanner.NodeContainerBase
}

func NewEscapeNode(p scanner.Parser, location scanner.Location) EscapeNode {
	return &escapenode{
		NodeContainerBase: scanner.NewContainerBase("escape", p.Document(), location, p.State.Container),
	}
}

func (n *escapenode) Print(gap string) {
	fmt.Printf("%sESCAPE\n", gap)
	fmt.Printf("%s  nodes:\n", gap)
	n.NodeSequence.Print(gap + "  ")
}

func (n *escapenode) Register(ctx scanner.ResolutionContext) error {
	nctx := NewEscapeContext(n, ctx)
	ctx.SetNodeContext(n, nctx)
	return n.NodeSequence.Register(ctx)
}

func (n *escapenode) ResolveLabels(ctx scanner.ResolutionContext) error {
	return n.NodeSequence.ResolveLabels(ctx)
}

func (n *escapenode) ResolveValues(ctx scanner.ResolutionContext) error {
	nctx := scanner.GetNodeContext[*EscapeNodeContext](ctx, n)

	err := n.NodeSequence.ResolveValues(ctx)
	if err != nil {
		return err
	}

	buf := scanner.NewBufferContext(ctx)
	err = n.NodeSequence.Emit(buf)
	if err != nil {
		ctx.RegisterUnresolved(nctx, err)
		return nil
	}
	text := buf.String()
	nctx.text = &text
	return nil
}

func (n *escapenode) Emit(ctx scanner.ResolutionContext) error {
	w := ctx.Writer()
	nctx := scanner.GetNodeContext[*EscapeNodeContext](ctx, n)

	txt := strings.Split(*nctx.text, "</br>")
	for i, t := range txt {
		txt[i] = html.EscapeString(t)
	}
	fmt.Fprintf(w, "%s", strings.Join(txt, "</br>"))
	return nil
}
