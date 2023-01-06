/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package title

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
	return &Statement{scanner.NewStatementBase("title")}
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
		return nil, e.Errorf("tag required for title")
	}
	link, err := utils.ParseAbsoluteLink(tag, "", false)
	if err != nil {
		return nil, e.Errorf("%s", err.Error())
	}
	p.State.SetLastTag(tag)
	p.State.Container.AddNode(NewTitleNode(p.Document(), e.Location(), link))
	return p.NextElement()
}

////////////////////////////////////////////////////////////////////////////////

type TitleNodeContext struct {
	scanner.LinkContextInfoNode[*titlenode]
}

func NewTitleNodeContext(n *titlenode, ctx scanner.ResolutionContext) (*TitleNodeContext, error) {
	var err error
	c := &TitleNodeContext{}
	c.LinkContextInfoNode, err = scanner.NewLinkContextInfoNode(n, ctx)
	return c, err
}

type TitleNode = *titlenode

type titlenode struct {
	scanner.NodeBase
	link utils.Link
}

func NewTitleNode(d scanner.Document, location scanner.Location, link utils.Link) TitleNode {
	return &titlenode{
		NodeBase: scanner.NewNodeBase(d, location),
		link:     link,
	}
}

func (n *titlenode) GetLink() utils.Link {
	return n.link
}

func (n *titlenode) Print(gap string) {
	fmt.Printf("%sTITLE[%s]\n", gap, n.link)
}

func (n *titlenode) Register(ctx scanner.ResolutionContext) error {
	nctx, err := NewTitleNodeContext(n, ctx)
	if err != nil {
		return err
	}
	ctx.SetNodeContext(n, nctx)
	return nil
}

func (n *titlenode) ResolveLabels(ctx scanner.ResolutionContext) error {
	nctx := scanner.GetNodeContext[*TitleNodeContext](ctx, n)
	return nctx.Resolve(ctx)
}

func (n *titlenode) Emit(ctx scanner.ResolutionContext) error {
	nctx := scanner.GetNodeContext[*TitleNodeContext](ctx, n)
	link, err := nctx.Link(ctx)
	if err != nil {
		return err
	}
	rt := nctx.Title()
	if rt == nil {
		return n.Errorf("unresolved title for %s", link)
	}
	w := ctx.Writer()
	fmt.Fprintf(w, "%s", *rt)
	return nil
}
