/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package termdef

import (
	"fmt"
	"strings"

	"github.com/mandelsoft/mdgen/scanner"
	"github.com/mandelsoft/mdgen/statements/section"
	"github.com/mandelsoft/mdgen/statements/subrange"
	utils2 "github.com/mandelsoft/mdgen/utils"
)

const GT_TERM = "term"

func init() {
	scanner.Tokens.RegisterStatement(NewStatement(), true)
}

type Statement struct {
	scanner.BracketedStatement[TermDefNode]
}

func NewStatement() *Statement {
	return &Statement{scanner.NewBracketedStatement[TermDefNode]("termdef", true)}
}

func stripFormat(s, f string, format string) (string, string) {
	if len(s) > 1 && strings.HasPrefix(s, f) && strings.HasSuffix(s, f) {
		return s[1 : len(s)-1], format + f
	}
	return s, format
}

func (s *Statement) Start(p scanner.Parser, e scanner.Element) (scanner.Element, error) {
	tag, err := e.Tag("term name")
	if err != nil {
		return nil, err
	}
	skip := strings.HasPrefix(tag, "-")
	if skip {
		tag = tag[1:]
	}
	n := NewTermDefNode(p, e.Location(), tag, skip)

	stop := func(p scanner.Parser, e scanner.Element) bool {
		if e.Token() != "description" {
			return false
		}
		return true
	}
	e, ns, err := scanner.ParseSequenceUntil(p, e, stop)
	if err != nil {
		return nil, err
	}

	if e.Token() != "description" {
		return nil, e.Errorf("description token required")
	}
	n.termnodes = ns
	p.State = p.State.Sub(n)
	return p.NextElement()
}

////////////////////////////////////////////////////////////////////////////////

type TermDefNodeContext struct {
	scanner.NodeContextBase[*TermdefNode]
	ctx          scanner.ResolutionContext
	term         *Term
	format       string
	singular     string
	plural       string
	referencable *subrange.NodeContext
}

func NewTermDefNodeContext(n *TermdefNode, ctx scanner.ResolutionContext, rctx *subrange.NodeContext) (*TermDefNodeContext, error) {
	nctx := &TermDefNodeContext{
		NodeContextBase: scanner.NewNodeContextBase(n, ctx),
		ctx:             ctx,
		referencable:    rctx,
	}
	ref, explicit, err := n.term.Evaluate(ctx)
	if err != nil {
		return nil, n.Error(err)
	}
	nctx.term = &Term{
		resolved: nctx,
		TermRef:  ref,
	}
	//fmt.Printf("#### scope %s: register term %q[%s] (node %p) %t\n", ctx.GetScope().GetNamespace(), nctx.term.tag, n.term.tag, nctx, explicit)
	err = ctx.RegisterTag(GT_TERM, nctx.term.tag, nctx, explicit)
	return nctx, err
}

func (c *TermDefNodeContext) GetContext() scanner.ResolutionContext {
	return c.ctx
}

func (c *TermDefNodeContext) GetLink() utils2.Link {
	return c.referencable.GetLink()
}

func (c *TermDefNodeContext) GetNodeSequence() scanner.NodeSequence {
	return c.EffNode().NodeSequence
}

func (c *TermDefNodeContext) Term() *Term {
	return c.term
}

type TermDefNode = *TermdefNode

type TermdefNode struct {
	scanner.NodeContainerBase
	termnodes scanner.NodeSequence
	term      TermRef
	omit      bool
}

func NewTermDefNode(p scanner.Parser, location scanner.Location, tag string, omit bool) TermDefNode {
	return &TermdefNode{
		NodeContainerBase: scanner.NewContainerBase("termdef", p.Document(), location, p.State.Container),
		term:              MapTermTag(tag),
		omit:              omit,
	}
}

func (n *TermdefNode) Print(gap string) {
	omit := ""
	if n.omit {
		omit = " (omitted)"
	}
	fmt.Printf("%sTERMDEF %s%s\n", gap, n.term.tag, omit)
	if n.termnodes.GetNodes() != nil {
		fmt.Printf("%s  term nodes:\n", gap)
		n.termnodes.Print(gap + "    ")
	}
	if n.GetNodes() != nil {
		fmt.Printf("%s  description nodes:\n", gap)
		n.NodeSequence.Print(gap + "    ")
	}
}

func (n *TermdefNode) Register(ctx scanner.ResolutionContext) error {
	rctx := scanner.LookupNodeContext[*subrange.NodeContext, section.Node](ctx)
	if rctx == nil {
		return n.Errorf("no anchor section found for term")
	}

	nctx, err := NewTermDefNodeContext(n, ctx, rctx)
	if err != nil {
		return err
	}
	ctx.SetNodeContext(n, nctx)
	return n.termnodes.Register(ctx)
}

func (n *TermdefNode) ResolveLabels(ctx scanner.ResolutionContext) error {
	return n.termnodes.ResolveLabels(ctx)
}

func (n *TermdefNode) ResolveValues(ctx scanner.ResolutionContext) error {
	err := n.termnodes.ResolveValues(ctx)
	if err != nil {
		return err
	}

	nctx := scanner.GetNodeContext[*TermDefNodeContext](ctx, n)
	buf := scanner.NewBufferContext(ctx)
	err = n.termnodes.Emit(buf)
	if err != nil {
		ctx.RegisterUnresolved(nctx, err)
		return nil
	}
	text := strings.TrimSpace(buf.String())
	i := strings.Index(text, "\n")
	if i >= 0 {
		return n.termnodes.GetNodes()[0].Errorf("resolved term contains newline")
	}

	text = strings.TrimSpace(text)

	format := ""
	for {
		var n string
		n, format = stripFormat(text, "*", format)
		n, format = stripFormat(n, "_", format)
		n, format = stripFormat(n, "`", format)
		if text == n {
			break
		}
		text = n
	}
	singular := text
	plural := ""
	if i := strings.Index(singular, "/"); i > 0 {
		plural = strings.TrimSpace(singular[i+1:])
		singular = strings.TrimSpace(singular[:i])
	}
	if plural == "" {
		plural = utils2.Plural(singular)
	}
	nctx.singular = singular
	nctx.plural = plural
	nctx.format = format

	return nil
}

func (n *TermdefNode) Emit(ctx scanner.ResolutionContext) error {
	if !n.omit {
		nctx := scanner.GetNodeContext[*TermDefNodeContext](ctx, n)
		w := ctx.Writer()
		if nctx.term.IsFormatted() {
			fmt.Fprintf(w, "%s", nctx.term.Format())
		} else {
			fmt.Fprintf(w, "*%s*", nctx.term.Get())
		}
	}
	return nil
}
