/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package term

import (
	"fmt"
	"strings"

	"github.com/mandelsoft/mdgen/scanner"
	"github.com/mandelsoft/mdgen/statements/termdef"
)

func init() {
	scanner.Tokens.RegisterStatement(NewStatement())

	scanner.Keywords.Register("description")

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
	td := scanner.Lookup[*termdef.TermdefNode](p)
	glossary := td != nil
	n := NewTermNode(p, e.Location(), tag, label, glossary, link)
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
	term     termdef.TermRef
	label    bool
	link     bool
	glossary bool
}

func NewTermNode(p scanner.Parser, location scanner.Location, tag string, label, glossary bool, link bool) TermNode {
	term := termdef.MapTermTag(tag)
	return &termnode{
		NodeBase: scanner.NewNodeBase(p.Document(), location),
		term:     term,
		label:    label,
		glossary: glossary,
		link:     link,
	}
}

func (n *termnode) Print(gap string) {
	mode := n.term.Mode()
	if n.glossary {
		mode = "glossary"
	} else if n.label {
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
	w := ctx.Writer()
	ref := nctx.term.GetLink()
	link, err := ctx.DetermineLink(ref)
	if err != nil {
		return err
	}
	if n.glossary {
		link = "#glossary/" + nctx.term.Tag()
	}
	if n.link {
		fmt.Fprintf(w, "<a href=\"%s\">", link)
	}
	if n.label {
		label := ctx.GetLinkInfo(ref).Label().Name()
		if label == "" {
			return n.Errorf("no label found for term %q", nctx.term.Tag())
		}
		fmt.Fprintf(w, "%s", label)
	} else {
		fmt.Fprintf(w, "%s", nctx.term.Format())
	}
	if n.link {
		fmt.Fprintf(w, "</a>")
	}
	return nil
}
