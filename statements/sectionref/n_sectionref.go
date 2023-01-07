/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sectionref

import (
	"fmt"

	"github.com/mandelsoft/mdgen/scanner"
	"github.com/mandelsoft/mdgen/utils"
)

func init() {
	scanner.Tokens.RegisterStatement(NewStatement())
}

type Statement struct {
	scanner.BracketedStatement[SectionRefNode]
}

func NewStatement() *Statement {
	return &Statement{scanner.NewBracketedStatement[SectionRefNode]("sectionref", true)}
}

func (s *Statement) Start(p scanner.Parser, e scanner.Element) (scanner.Element, error) {
	tag, err := e.Tag("section ref")
	if err != nil {
		return nil, err
	}

	err = scanner.ForbidNesting[SectionRefNode]("sectionref", p, e)
	if err != nil {
		return nil, err
	}

	sid := p.State.NextId(scanner.SECTION_TYPE).Id()
	link, err := utils.ParseAbsoluteLink(tag, "", false)
	if err != nil {
		return nil, e.Errorf("%s", err.Error())
	}
	n := NewSectionRefNode(sid, p, e.Location(), link, tag)
	if e.IsFlagged() {
		p.State = p.State.Sub(n)
		p.State.SetLastTag(tag)
	} else {
		{
			p.State.Container.AddNode(n)
		}
	}
	return p.NextElement()
}

////////////////////////////////////////////////////////////////////////////////

type SectionRefNodeContext struct {
	scanner.LinkContextInfoNode[*sectionrefnode]
	id     scanner.TaggedId
	hlabel scanner.HierarchyLabel
}

func NewSectionRefNodeContext(n *sectionrefnode, ctx scanner.ResolutionContext, id scanner.TaggedId) (*SectionRefNodeContext, error) {
	var err error
	c := &SectionRefNodeContext{
		id:     id,
		hlabel: nil,
	}
	c.LinkContextInfoNode, err = scanner.NewLinkContextInfoNode(n, ctx)
	return c, err
}

type SectionRefNode = *sectionrefnode

type sectionrefnode struct {
	scanner.NodeContainerBase
	sid    scanner.TaggedId
	hlabel scanner.HierarchyLabel
	link   utils.Link
}

func NewSectionRefNode(sid scanner.TaggedId, p scanner.Parser, location scanner.Location, link utils.Link, tag string) SectionRefNode {
	return &sectionrefnode{
		NodeContainerBase: scanner.NewContainerBase("sectionref", p.Document(), location, p.State.Container),
		sid:               sid,
		link:              link,
	}
}

func (n *sectionrefnode) Id() scanner.TaggedId {
	return n.sid
}

func (n *sectionrefnode) GetLink() utils.Link {
	return n.link
}

func (n *sectionrefnode) Print(gap string) {
	fmt.Printf("%sSECTIONREF %s[%s]:\n", gap, n.sid, n.link)
	n.NodeContainerBase.Print(gap + "  ")
}

func (n *sectionrefnode) Register(ctx scanner.ResolutionContext) error {
	id := ctx.NextId(scanner.SECTION_TYPE).Id()
	ctx.RequestNumberRange(scanner.SECTION_TYPE)
	nctx, err := NewSectionRefNodeContext(n, ctx, id)
	if err != nil {
		return err
	}
	ctx.SetNodeContext(n, nctx)
	err = ctx.RequestDocument(nctx.GetLink(), ctx.GetDocument())
	if err != nil {
		return err
	}

	return n.NodeSequence.Register(ctx)
}

func (n *sectionrefnode) ResolveLabels(ctx scanner.ResolutionContext) error {
	nctx := scanner.GetNodeContext[*SectionRefNodeContext](ctx, n)
	err := nctx.Resolve(ctx)
	if err != nil {
		return err
	}
	d := ctx.GetDocumentForLink(nctx.GetLink())
	if d == nil {
		return n.Errorf("cannot resolve sectionref %s", nctx.GetLink())
	}
	ctx.SetNumberRangeFor(d, nctx.id, scanner.SECTION_TYPE, ctx.GetNumberRange(scanner.SECTION_TYPE))
	return n.NodeSequence.ResolveLabels(ctx)
}

func (n *sectionrefnode) Emit(ctx scanner.ResolutionContext) error {
	nctx := scanner.GetNodeContext[*SectionRefNodeContext](ctx, n)
	w := ctx.Writer()
	link, err := nctx.Link(ctx)
	if err != nil {
		return n.Location().Errorf("section ref: %s", err)
	}
	if len(n.NodeSequence.GetNodes()) > 0 {
		fmt.Fprintf(w, "<a href=\"%s\">", link)
		err = n.NodeSequence.Emit(ctx)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "</a>")
	}
	return nil
}
