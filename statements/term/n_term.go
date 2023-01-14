/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package term

import (
	"fmt"
	"strings"

	"github.com/mandelsoft/mdgen/render"
	"github.com/mandelsoft/mdgen/scanner"
	"github.com/mandelsoft/mdgen/statements/glossary"
	"github.com/mandelsoft/mdgen/statements/termdef"
)

func init() {
	scanner.Tokens.RegisterStatement(NewStatement())

	scanner.Keywords.Register("description", true)

}

type Statement struct {
	scanner.StatementBase
}

func NewStatement() scanner.Statement {
	return &Statement{scanner.NewStatementBase("term")}
}

func (s *Statement) Start(p scanner.Parser, e scanner.Element) (scanner.Element, error) {
	tag, err := e.Tag("term name")
	if err != nil {
		return nil, err
	}

	link := true
	if strings.HasPrefix(tag, "!") {
		link = false
		tag = tag[1:]
	}

	label := false
	if strings.HasPrefix(tag, "#") {
		label = true
		tag = tag[1:]
	}
	n := NewTermNode(p, e.Location(), tag, label, link)
	p.State.Container.AddNode(n)
	return p.NextElement()
}

////////////////////////////////////////////////////////////////////////////////

type TermNodeContext struct {
	scanner.NodeContextBase[*termnode]
	term *termdef.Term
}

func NewTermNodeContext(n *termnode, ctx scanner.ResolutionContext) (*TermNodeContext, error) {
	ref, _, err := n.term.Evaluate(ctx)
	if err != nil {
		return nil, n.Error(err)
	}
	return &TermNodeContext{
		NodeContextBase: scanner.NewNodeContextBase(n, ctx),
		term:            termdef.NewTerm(ref),
	}, nil
}

type TermNode = *termnode

type termnode struct {
	scanner.NodeBase
	term  termdef.TermRef
	label bool
	link  bool
}

func NewTermNode(p scanner.Parser, location scanner.Location, tag string, label, link bool) TermNode {
	term := termdef.MapTermTag(tag)
	return &termnode{
		NodeBase: scanner.NewNodeBase(p.Document(), location),
		term:     term,
		label:    label,
		link:     link,
	}
}

func (n *termnode) Print(gap string) {
	mode := n.term.Mode()
	if n.label {
		mode = "label"
	}
	fmt.Printf("%sTERM %s[%s]\n", gap, n.term.Tag(), mode)
}

func (n *termnode) Register(ctx scanner.ResolutionContext) error {
	nctx, err := NewTermNodeContext(n, ctx)
	if err != nil {
		return err
	}
	ctx.SetNodeContext(n, nctx)
	return nil
}

func (n *termnode) ResolveLabels(ctx scanner.ResolutionContext) error {
	nctx := scanner.GetNodeContext[*TermNodeContext](ctx, n)

	err := nctx.term.Resolve(ctx)
	if err != nil {
		return n.Errorf("%s", err)
	}
	return nil
}

func (n *termnode) Emit(ctx scanner.ResolutionContext) error {
	nctx := scanner.GetNodeContext[*TermNodeContext](ctx, n)
	ref := nctx.term.GetLink()
	link, err := ctx.DetermineLink(ref)
	if err != nil {
		return err
	}

	if ctx.Info(glossary.InfoKey) == true {
		link = "#glossary/" + nctx.term.Tag()
	}

	content := func(ctx scanner.ResolutionContext) error {
		out := ""
		if n.label {
			label := ctx.GetLinkInfo(ref).Label().Name()
			if label == "" {
				return n.Errorf("no label found for term %q", nctx.term.Tag())
			}
			out = label
		} else {
			out = nctx.term.Format()
		}
		fmt.Fprint(ctx.Writer(), out)
		return nil
	}

	if n.link {
		render.Current.Link(ctx, link, content)
	} else {
		content(ctx)
	}
	return nil
}
