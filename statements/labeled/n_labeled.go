/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package figure

import (
	"fmt"
	"strings"

	"github.com/mandelsoft/mdgen/scanner"
)

func init() {
	scanner.Tokens.RegisterStatement(NewStatement())

	scanner.Keywords.Register("content")
}

type Statement struct {
	scanner.BracketedStatement[LabeledNode]
}

func NewStatement() scanner.Statement {
	return &Statement{scanner.NewBracketedStatement[LabeledNode]("labeled", true)}
}

func (s *Statement) Start(p scanner.Parser, e scanner.Element) (scanner.Element, error) {

	if !e.HasTags() {
		return nil, e.Errorf("number range type required")
	}
	tags := e.Tags()
	tag := tags[0]
	if len(tags) > 2 {
		return nil, e.Errorf("too many arguments")
	}
	mode := "box"
	if len(tags) > 1 {
		mode = tags[1]
		switch mode {
		case "box":
		case "float":
		default:
			return nil, e.Errorf("mode must be box or float, but found %q", tags[1])
		}
	}

	comps := strings.Split(tag, ":")
	if len(comps) > 2 {
		return nil, e.Errorf("use tag with <number range>[:tag]")
	}
	typ := comps[0]
	if typ == "" {
		return nil, e.Errorf("nuber range name must be non-empty")
	}
	if len(comps) > 1 {
		tag = comps[1]
	} else {
		tag = ""
	}
	sid := p.State.NextId(typ).Id()

	n := NewLabeledNode(p.State.Container, p.Document(), e.Location(), sid, tag, mode == "box")

	stop := func(p scanner.Parser, e scanner.Element) bool {
		if e.Token() != "content" {
			return false
		}
		return true
	}
	e, ns, err := scanner.ParseSequenceUntil(p, e, stop)
	if err != nil {
		return nil, err
	}
	n.title = ns

	p.State = p.State.Sub(n)
	p.State.SetLastTag(tag)

	if e.Token() == "content" {
		return p.NextElement()
	}
	return e, nil
}

////////////////////////////////////////////////////////////////////////////////

type LabeledNodeContext = scanner.LabeledNodeContextBase[*labelednode]

func NewLabeledNodeContext(n *labelednode, ctx scanner.ResolutionContext, title scanner.NodeSequence) (*LabeledNodeContext, error) {
	return scanner.NewLabeledNodeContextBase(n, ctx, title)
}

////////////////////////////////////////////////////////////////////////////////

type LabeledNode = *labelednode

type labelednode struct {
	scanner.TaggedNodeBase
	scanner.NodeContainerBase
	title scanner.NodeSequence
	box   bool
}

func NewLabeledNode(p scanner.NodeContainer, d scanner.Document, location scanner.Location, sid scanner.TaggedId, tag string, box bool) LabeledNode {
	return &labelednode{
		TaggedNodeBase:    scanner.NewTaggedNodeBase(sid, tag),
		NodeContainerBase: scanner.NewContainerBase("labeled", d, location, p),
		box:               box,
	}
}

func (n *labelednode) Print(gap string) {
	fmt.Printf("%sLABELED %s[%s] %t\n", gap, n.Id(), n.Tag(), n.box)
	fmt.Printf("%s  title:\n", gap)
	n.title.Print(gap + "  ")
	fmt.Printf("%s  nodes:\n", gap)
	n.NodeContainerBase.Print(gap + "  ")
}

func (n *labelednode) Register(ctx scanner.ResolutionContext) error {
	nctx, err := NewLabeledNodeContext(n, ctx, n.title)
	if err != nil {
		return err
	}
	ctx.SetNodeContext(n, nctx)
	return n.NodeSequence.Register(ctx)
}

func (n *labelednode) ResolveLabels(ctx scanner.ResolutionContext) error {
	nctx := scanner.GetNodeContext[*LabeledNodeContext](ctx, n)
	err := nctx.ResolveLabels(ctx)
	if err != nil {
		return err
	}
	return n.NodeSequence.ResolveLabels(ctx)
}

func (n *labelednode) ResolveValues(ctx scanner.ResolutionContext) error {
	nctx := scanner.GetNodeContext[*LabeledNodeContext](ctx, n)
	err := nctx.ResolveValues(ctx)
	if err != nil {
		return err
	}
	return n.NodeSequence.ResolveValues(ctx)
}

func (n *labelednode) Emit(ctx scanner.ResolutionContext) error {
	w := ctx.Writer()

	nctx := scanner.GetNodeContext[*LabeledNodeContext](ctx, n)

	nctx.EmitAnchors(ctx)
	if n.box {
		fmt.Fprintf(w, "<div align=\"center\"><table><tr><td>\n\n")
	}
	err := n.NodeSequence.Emit(ctx)
	if err != nil {
		return err
	}
	if n.box {
		fmt.Fprintf(w, "</td></tr></table>\n")
	} else {
		fmt.Fprintf(w, "<div align=\"center\">\n")
	}
	nctx.EmitTitle(ctx)
	fmt.Fprintf(w, "</div>\n")
	return nil
}
