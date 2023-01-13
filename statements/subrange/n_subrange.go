/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package subrange

import (
	"fmt"
	"strings"

	"github.com/mandelsoft/mdgen/scanner"
	utils2 "github.com/mandelsoft/mdgen/utils"
)

func init() {
	scanner.Tokens.RegisterStatement(newSubRangeStatement())
}

func newSubRangeStatement() scanner.Statement {
	s := NewStatement[Node]("subrange", nil, false)
	s.creator = s.newNode
	return s
}

type Node interface {
	scanner.TaggedNode
	scanner.LabelResolver
	scanner.NodeContainer

	GetRangeType() string

	SetTitleSequence(ns scanner.NodeSequence)
}

////////////////////////////////////////////////////////////////////////////////

type NodeCreator[N Node] func(p scanner.Parser, e scanner.Element) (N, error)

type Statement[N Node] struct {
	scanner.BracketedStatement[N]
	titleRequired bool
	creator       NodeCreator[N]
}

func NewStatement[N Node](name string, c NodeCreator[N], titleRequired bool) *Statement[N] {
	scanner.Keywords.Register("subrange", titleRequired)
	return &Statement[N]{scanner.NewBracketedStatement[N](name, true), titleRequired, c}
}

func (s *Statement[N]) newNode(p scanner.Parser, e scanner.Element) (Node, error) {
	tag, err := e.OptionalTag("tag")
	if err != nil {
		return nil, err
	}

	comps := strings.Split(tag, ":")
	typ := comps[0]
	switch len(comps) {
	case 1:
		tag = ""
	case 2:
		tag = comps[1]
	default:
		return nil, e.Errorf("argument must be of <numberrange>[:<tag>]")
	}
	if typ == "" {
		return nil, e.Errorf("number range type may not be empty")
	}
	return s.NewNode(s.Name(), p, e.Location(), typ, tag), nil
}

func (s *Statement[N]) NewNode(name string, p scanner.Parser, location scanner.Location, typ, tag string) Node {
	sid := p.State.NextId(typ).Id()
	return &node{
		name:              name,
		typ:               typ,
		TaggedNodeBase:    scanner.NewTaggedNodeBase(sid, tag),
		NodeContainerBase: scanner.NewContainerBase(s.Name(), p.Document(), location, p.State.Container),
		LabelRules:        scanner.LabelRules{},
	}
}

func (s *Statement[N]) Start(p scanner.Parser, e scanner.Element) (scanner.Element, error) {
	n, err := s.creator(p, e)
	if err != nil {
		return nil, err
	}
	err = p.State.Container.RegisterReferencable(n)
	if err != nil {
		return nil, err
	}
	p.State = p.State.Sub(n)
	p.State.SubId(n.GetRangeType())

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
		return nil, e.Errorf("%s title expected", s.Name())
	}

	i := strings.Index(e.Text(), "\n")
	if i == 0 && len(ns.GetNodes()) == 0 && s.titleRequired {
		return nil, e.Errorf("%s title must follow the %s token", s.Name(), s.Name())
	}
	switch {
	case i < 0:
		return nil, e.Errorf("%s title must follow the %s token followed by a newline", s.Name(), s.Name())
	case i == 0:
		return scanner.NewText(e.Text()[i+1:], e.Location().SkipLine()), nil
	case i > 0:
		ns.AddNode(scanner.NewTextNode(p.Document(), e.Location().SkipLine(), e.Text()[:i]))
		if e.Text()[i+1:] != "" {
			return scanner.NewText(e.Text()[i+1:], e.Location().SkipLine()), nil
		}
	}
	n.SetTitleSequence(ns)
	return p.NextElement()
}

////////////////////////////////////////////////////////////////////////////////

type NodeContext struct {
	*scanner.LabeledNodeContextBase[Node]
	ctx scanner.SubNumberRangeContext
}

func NewSectionNodeContext(n Node, ctx scanner.ResolutionContext, title scanner.NodeSequence) (*NodeContext, error) {
	base, err := scanner.NewLabeledNodeContextBase(n, ctx, title)
	if err != nil {
		return nil, err
	}
	return &NodeContext{
		LabeledNodeContextBase: base,
	}, nil
}

func (c *NodeContext) GetLink() utils2.Link {
	return scanner.NewLink(c.ctx, c.ctx.GetReferencable(c.Id()).Anchors()...)
}

type node struct {
	scanner.TaggedNodeBase
	scanner.LabelRules
	scanner.NodeContainerBase
	title scanner.NodeSequence
	typ   string
	name  string
}

var _ (Node) = (*node)(nil)

func (n *node) GetRangeType() string {
	return n.typ
}

func (n *node) SetTitleSequence(ns scanner.NodeSequence) {
	n.title = ns
}

func (n *node) Print(gap string) {
	fmt.Printf("%s%s %s[%s]\n", gap, strings.ToUpper(n.name), n.Id(), n.Tag())
	fmt.Printf("%s  title:\n", gap)
	n.title.Print(gap + "  ")
	fmt.Printf("%s  nodes:\n", gap)
	n.NodeContainerBase.Print(gap + "  ")
}

func (n *node) Register(ctx scanner.ResolutionContext) error {
	var err error
	nctx, err := NewSectionNodeContext(n, ctx, n.title)
	if err != nil {
		return err
	}
	nctx.ctx = scanner.NewSubNumberRangeContext(n.typ, ctx, nctx.IdRule(), nctx)
	ctx.SetNodeContext(n, nctx)
	return n.NodeSequence.Register(nctx.ctx)
}

func (n *node) ResolveLabels(ctx scanner.ResolutionContext) error {
	nctx := scanner.GetNodeContext[*NodeContext](ctx, n)
	err := nctx.ResolveLabels(ctx)
	if err != nil {
		return err
	}
	subctx := nctx.ctx
	subctx.SetNumberRange(nctx.NumberRange().Sub())

	return n.NodeSequence.ResolveLabels(subctx)
}

func (n *node) ResolveValues(ctx scanner.ResolutionContext) error {
	nctx := scanner.GetNodeContext[*NodeContext](ctx, n)

	err := nctx.ResolveValues(ctx)
	if err != nil {
		return err
	}
	return n.NodeSequence.ResolveValues(nctx.ctx)
}

func (n *node) Emit(ctx scanner.ResolutionContext) error {
	nctx := scanner.GetNodeContext[*NodeContext](ctx, n)
	nctx.EmitAnchors(ctx)
	return n.NodeSequence.Emit(nctx.ctx)
}
