/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package section

import (
	"fmt"
	"strings"

	"github.com/mandelsoft/mdgen/scanner"
	"github.com/mandelsoft/mdgen/statements/sectionref"
	utils2 "github.com/mandelsoft/mdgen/utils"
)

func init() {
	scanner.Tokens.RegisterStatement(NewStatement())
}

type Statement struct {
	scanner.BracketedStatement[SectionNode]
}

func NewStatement() *Statement {
	return &Statement{scanner.NewBracketedStatement[SectionNode]("section", true)}
}

func (s *Statement) Start(p scanner.Parser, e scanner.Element) (scanner.Element, error) {
	tag, err := e.OptionalTag("tag")
	if err != nil {
		return nil, e.Errorf("%s", err.Error())
	}

	if err = scanner.ForbidNesting[sectionref.SectionRefNode]("sectionref", p, e); err != nil {
		return nil, err
	}

	sid := p.State.NextId(scanner.SECTION_TYPE).Id()
	n := NewSectionNode(sid, p, e.Location(), tag)
	err = p.State.Container.RegisterReferencable(n)
	if err != nil {
		return nil, err
	}
	p.State = p.State.Sub(n)
	p.State.SubId(scanner.SECTION_TYPE)

	stop := func(p scanner.Parser, e scanner.Element) bool {
		if !e.IsText() {
			// TODO: check for valid tokens
			return false
		}
		i := strings.Index(e.Text(), "\n")
		return i >= 0
	}
	e, ns, err := scanner.ParseSequenceUntil(p, e, stop)
	if err != nil {
		return nil, err
	}
	if !e.IsText() {
		return nil, e.Errorf("section title expected")
	}

	i := strings.Index(e.Text(), "\n")
	if i == 0 && len(ns.GetNodes()) == 0 {
		return nil, e.Errorf("section title must follow the section token")
	}
	n.title = ns
	switch {
	case i < 0:
		return nil, e.Errorf("section title must follow the section token followed by a newline")
	case i == 0:
		return scanner.NewText(e.Text()[i+1:], e.Location().SkipLine()), nil
	case i > 0:
		ns.AddNode(scanner.NewTextNode(p.Document(), e.Location().SkipLine(), e.Text()[:i]))
		if e.Text()[i+1:] != "" {
			return scanner.NewText(e.Text()[i+1:], e.Location().SkipLine()), nil
		}
	}
	return p.NextElement()
}

////////////////////////////////////////////////////////////////////////////////

type SectionNodeContext struct {
	*scanner.LabeledNodeContextBase[*sectionnode]
	ctx scanner.SubNumberRangeContext
}

func NewSectionNodeContext(n *sectionnode, ctx scanner.ResolutionContext, title scanner.NodeSequence) (*SectionNodeContext, error) {
	base, err := scanner.NewLabeledNodeContextBase(n, ctx, title)
	if err != nil {
		return nil, err
	}
	return &SectionNodeContext{
		LabeledNodeContextBase: base,
	}, nil
}

func (c *SectionNodeContext) GetLink() utils2.Link {
	return scanner.NewLink(c.ctx, c.ctx.GetReferencable(c.Id()).Anchors()...)
}

type SectionNode = *sectionnode

type sectionnode struct {
	scanner.TaggedNodeBase
	scanner.LabelRules
	scanner.NodeContainerBase
	title scanner.NodeSequence
}

var _ (scanner.LabelResolver) = (*sectionnode)(nil)

func NewSectionNode(sid scanner.TaggedId, p scanner.Parser, location scanner.Location, tag string) SectionNode {
	return &sectionnode{
		TaggedNodeBase:    scanner.NewTaggedNodeBase(sid, tag),
		NodeContainerBase: scanner.NewContainerBase("section", p.Document(), location, p.State.Container),
		LabelRules:        scanner.LabelRules{},
	}
}

func (n *sectionnode) Print(gap string) {
	fmt.Printf("%sSECTION %s[%s]\n", gap, n.Id(), n.Tag())
	fmt.Printf("%s  title:\n", gap)
	n.title.Print(gap + "  ")
	fmt.Printf("%s  nodes:\n", gap)
	n.NodeContainerBase.Print(gap + "  ")
}

func (n *sectionnode) Register(ctx scanner.ResolutionContext) error {
	var err error
	nctx, err := NewSectionNodeContext(n, ctx, n.title)
	if err != nil {
		return err
	}
	nctx.ctx = scanner.NewSubNumberRangeContext(scanner.SECTION_TYPE, ctx, nctx.IdRule(), nctx)
	ctx.SetNodeContext(n, nctx)
	return n.NodeSequence.Register(nctx.ctx)
}

func (n *sectionnode) ResolveLabels(ctx scanner.ResolutionContext) error {
	nctx := scanner.GetNodeContext[*SectionNodeContext](ctx, n)
	err := nctx.ResolveLabels(ctx)
	if err != nil {
		return err
	}
	subctx := nctx.ctx
	subctx.SetNumberRange(nctx.NumberRange().Sub())

	return n.NodeSequence.ResolveLabels(subctx)
}

func (n *sectionnode) ResolveValues(ctx scanner.ResolutionContext) error {
	nctx := scanner.GetNodeContext[*SectionNodeContext](ctx, n)

	err := nctx.ResolveValues(ctx)
	if err != nil {
		return err
	}
	return n.NodeSequence.ResolveValues(nctx.ctx)
}

func (n *sectionnode) Emit(ctx scanner.ResolutionContext) error {
	w := ctx.Writer()
	nctx := scanner.GetNodeContext[*SectionNodeContext](ctx, n)

	nctx.EmitAnchors(ctx)
	info := ctx.GetReferencable(nctx.Id())

	lvl := info.Label().Level()
	if lvl < 0 {
		fmt.Printf("WARN: invalid level %d for %s\n", lvl, info.Label().Id())
		lvl = 1
	}
	for lvl >= 0 {
		fmt.Fprintf(w, "#")
		lvl--
	}
	label := info.Label().Name()
	if label != "" {
		label += " "
	}
	fmt.Fprintf(w, " %s%s\n", label, *nctx.Title())
	return n.NodeSequence.Emit(nctx.ctx)
}
