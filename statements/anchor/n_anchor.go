/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package anchor

import (
	"fmt"
	"strings"

	"github.com/mandelsoft/mdgen/scanner"
)

func init() {
	scanner.Tokens.RegisterStatement(NewStatement())
}

type Statement struct {
	scanner.BracketedStatement[AnchorNode]
}

func NewStatement() scanner.Statement {
	return &Statement{scanner.NewBracketedStatement[AnchorNode]("anchor", true)}
}

func (s *Statement) Start(p scanner.Parser, e scanner.Element) (scanner.Element, error) {
	tag, err := e.Tag("link")
	if err != nil {
		return nil, err
	}
	plain := !e.IsFlagged()
	typ := scanner.ANCHOR_TYPE
	comps := strings.Split(tag, ":")
	switch len(comps) {
	case 1:
	case 2:
		typ = comps[0]
	default:
		return nil, e.Errorf("invalid tag syntax ([<type>:]<anchor)")
	}

	omit := strings.HasPrefix(tag, "!")
	if omit {
		tag = tag[1:]
	}
	sid := p.State.NextId(typ).Id()
	n := NewAnchorNode(s.Name(), p.State.Container, p.Document(), e.Location(), sid, tag, omit)
	if !plain {
		p.State = p.State.Sub(n)
	} else {
		p.State.Container.AddNode(n)
	}
	p.State.SetLastTag(tag)
	return p.NextElement()
}

////////////////////////////////////////////////////////////////////////////////

type AnchorNodeContext = scanner.LabeledNodeContextBase[*anchornode]

func NewAnchorNodeContext(n *anchornode, ctx scanner.ResolutionContext, title scanner.NodeSequence) (*AnchorNodeContext, error) {
	return scanner.NewLabeledNodeContextBase(n, ctx, title)
}

type AnchorNode = *anchornode

type anchornode struct {
	scanner.TaggedNodeBase
	scanner.NodeContainerBase
	omit bool
}

func NewAnchorNode(name string, p scanner.NodeContainer, d scanner.Document, location scanner.Location, sid scanner.TaggedId, tag string, omit bool) AnchorNode {
	return &anchornode{
		TaggedNodeBase:    scanner.NewTaggedNodeBase(sid, tag),
		NodeContainerBase: scanner.NewContainerBase(name, d, location, p),
		omit:              omit,
	}
}

func (n *anchornode) Print(gap string) {
	omit := ""
	if n.omit {
		omit = " (omit)"
	}
	fmt.Printf("%sANCHOR %s[%s]%s\n", gap, n.Id(), n.Tag(), omit)
	n.NodeContainerBase.Print(gap + "  ")
}

func (n *anchornode) Register(ctx scanner.ResolutionContext) error {
	nctx, err := NewAnchorNodeContext(n, ctx, n.NodeSequence)
	if err != nil {
		return err
	}
	ctx.SetNodeContext(n, nctx)
	return nil
}

func (n *anchornode) ResolveLabels(ctx scanner.ResolutionContext) error {
	nctx := scanner.GetNodeContext[*AnchorNodeContext](ctx, n)
	return nctx.ResolveLabels(ctx)
}

func (n *anchornode) ResolveValues(ctx scanner.ResolutionContext) error {
	nctx := scanner.GetNodeContext[*AnchorNodeContext](ctx, n)
	return nctx.ResolveValues(ctx)
}

func (n *anchornode) Emit(ctx scanner.ResolutionContext) error {
	nctx := scanner.GetNodeContext[*AnchorNodeContext](ctx, n)
	nctx.EmitAnchors(ctx)
	if !n.omit {
		nctx.EmitTitle(ctx)
	}
	return nil
}
