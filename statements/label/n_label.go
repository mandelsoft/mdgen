/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package label

import (
	"fmt"

	"github.com/mandelsoft/mdgen/scanner"
	utils "github.com/mandelsoft/mdgen/utils"
)

func init() {
	scanner.Tokens.RegisterStatement(NewStatement())
}

type Statement struct {
	scanner.StatementBase
}

func NewStatement() scanner.Statement {
	return &Statement{scanner.NewStatementBase("label")}
}

func (s *Statement) Start(p scanner.Parser, e scanner.Element) (scanner.Element, error) {
	tag, err := e.OptionalTag("tag")
	if err != nil {
		return nil, err
	}
	if tag == "" {
		tag = p.State.LastTag()
	}
	if tag == "" {
		return nil, e.Errorf("tag required for label")
	}

	link, err := utils.ParseAbsoluteLink(tag, "", false)
	if err != nil {
		return nil, e.Errorf("%s", err.Error())
	}
	p.State.SetLastTag(tag)
	p.State.Container.AddNode(NewLabelNode(p.Document(), e.Location(), link))
	return p.NextElement()
}

////////////////////////////////////////////////////////////////////////////////

type LabelNodeContext struct {
	scanner.LinkContextInfoNode[*labelnode]
}

func NewLabelNodeContext(n *labelnode, ctx scanner.ResolutionContext) (*LabelNodeContext, error) {
	var err error
	c := &LabelNodeContext{}
	c.LinkContextInfoNode, err = scanner.NewLinkContextInfoNode(n, ctx)
	return c, err
}

type LabelNode = *labelnode

type labelnode struct {
	scanner.NodeBase
	link utils.Link
}

func NewLabelNode(d scanner.Document, location scanner.Location, link utils.Link) LabelNode {
	return &labelnode{
		NodeBase: scanner.NewNodeBase(d, location),
		link:     link,
	}
}

func (n *labelnode) GetLink() utils.Link {
	return n.link
}

func (n *labelnode) Print(gap string) {
	fmt.Printf("%sLABEL[%s]\n", gap, n.link)
}

func (n *labelnode) Register(ctx scanner.ResolutionContext) error {
	nctx, err := NewLabelNodeContext(n, ctx)
	if err != nil {
		return err
	}
	ctx.SetNodeContext(n, nctx)
	return nil
}

func (n *labelnode) ResolveLabels(ctx scanner.ResolutionContext) error {
	nctx := scanner.GetNodeContext[*LabelNodeContext](ctx, n)
	return nctx.Resolve(ctx)
}

func (n *labelnode) Emit(ctx scanner.ResolutionContext) error {
	nctx := scanner.GetNodeContext[*LabelNodeContext](ctx, n)

	w := ctx.Writer()
	l := nctx.Label()
	fmt.Fprintf(w, "%s", l.Name())
	return nil
}
