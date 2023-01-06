/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package link

import (
	"fmt"

	"github.com/mandelsoft/mdgen/scanner"
	"github.com/mandelsoft/mdgen/utils"
)

func init() {
	scanner.Tokens.RegisterStatement(NewStatement())
}

type Statement struct {
	scanner.BracketedStatement[LinkNode]
}

func NewStatement() *Statement {
	return &Statement{scanner.NewBracketedStatement[LinkNode]("link", true)}
}

func (s *Statement) Start(p scanner.Parser, e scanner.Element) (scanner.Element, error) {
	tag, err := e.Tag("link")
	if err != nil {
		return nil, err
	}
	link, err := utils.ParseAbsoluteLink(tag, "", false)
	if err != nil {
		return nil, e.Errorf("%s", err.Error())
	}
	p.State = p.State.Sub(NewLinkNode(p.State.Container, p.Document(), e.Location(), link))
	p.State.SetLastTag(tag)
	return p.NextElement()
}

////////////////////////////////////////////////////////////////////////////////

type LinkNodeContext struct {
	scanner.LinkContextInfoNode[*linknode]
}

func NewLinkNodeContext(n *linknode, ctx scanner.ResolutionContext) (*LinkNodeContext, error) {
	var err error
	c := &LinkNodeContext{}
	c.LinkContextInfoNode, err = scanner.NewLinkContextInfoNode(n, ctx)
	return c, err
}

type LinkNode = *linknode

type linknode struct {
	scanner.NodeContainerBase
	link utils.Link
}

func NewLinkNode(p scanner.NodeContainer, d scanner.Document, location scanner.Location, link utils.Link) LinkNode {
	return &linknode{
		NodeContainerBase: scanner.NewContainerBase("link", d, location, p),
		link:              link,
	}
}

func (n *linknode) GetLink() utils.Link {
	return n.link
}

func (t *linknode) Print(gap string) {
	fmt.Printf("%sLINK[%s]:\n", gap, t.link)
	t.NodeContainerBase.Print(gap + "  ")
}

func (n *linknode) Register(ctx scanner.ResolutionContext) error {
	nctx, err := NewLinkNodeContext(n, ctx)
	if err != nil {
		return err
	}
	ctx.SetNodeContext(n, nctx)
	return n.NodeSequence.Register(ctx)
}

func (n *linknode) ResolveLabels(ctx scanner.ResolutionContext) error {
	nctx := scanner.GetNodeContext[*LinkNodeContext](ctx, n)

	err := nctx.Resolve(ctx)
	if err != nil {
		return err
	}
	return n.NodeSequence.ResolveLabels(ctx)
}

func (n *linknode) Emit(ctx scanner.ResolutionContext) error {
	nctx := scanner.GetNodeContext[*LinkNodeContext](ctx, n)
	link, err := nctx.Link(ctx)
	if err != nil {
		return err
	}
	w := ctx.Writer()
	fmt.Fprintf(w, "<a href=\"%s\">", link)
	err = n.NodeSequence.Emit(ctx)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "</a>")
	return nil
}
