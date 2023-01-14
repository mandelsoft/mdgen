/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package link

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/mandelsoft/mdgen/render"
	"github.com/mandelsoft/mdgen/scanner"
	"github.com/mandelsoft/mdgen/utils"
)

func init() {
	scanner.Tokens.RegisterStatement(NewStatement())
}

type Statement struct {
	scanner.StatementBase
}

func NewStatement() *Statement {
	return &Statement{scanner.NewStatementBase("ref")}
}

func (s *Statement) Start(p scanner.Parser, e scanner.Element) (scanner.Element, error) {
	tag, err := e.Tag("link")
	if err != nil {
		return nil, err
	}
	abbrev := false
	if strings.HasPrefix(tag, "*") {
		abbrev = true
		tag = tag[1:]
	}
	upper := false
	if strings.HasPrefix(tag, "^") {
		upper = true
		tag = tag[1:]
	}
	link, err := utils.ParseAbsoluteLink(tag, "", false)
	if err != nil {
		return nil, e.Errorf("%s", err.Error())
	}
	n := NewRefNode(p.Document(), e.Location(), link)
	n.abbrev = abbrev
	n.upper = upper
	p.State.Container.AddNode(n)
	p.State.SetLastTag(tag)
	return p.NextElement()
}

////////////////////////////////////////////////////////////////////////////////

type RefNodeContext struct {
	scanner.LinkContextInfoNode[*Refnode]
}

func NewRefNodeContext(n *Refnode, ctx scanner.ResolutionContext) (*RefNodeContext, error) {
	var err error
	c := &RefNodeContext{}
	c.LinkContextInfoNode, err = scanner.NewLinkContextInfoNode(n, ctx)
	return c, err
}

type RefNode = *Refnode

type Refnode struct {
	scanner.NodeBase
	abbrev bool
	upper  bool
	link   utils.Link
}

func NewRefNode(d scanner.Document, location scanner.Location, link utils.Link) RefNode {
	return &Refnode{
		NodeBase: scanner.NewNodeBase(d, location),
		link:     link,
	}
}

func (n *Refnode) GetLink() utils.Link {
	return n.link
}

func (t *Refnode) Print(gap string) {
	fmt.Printf("%sREF[%s]:\n", gap, t.link)
}

func (n *Refnode) Register(ctx scanner.ResolutionContext) error {
	nctx, err := NewRefNodeContext(n, ctx)
	if err != nil {
		return err
	}
	ctx.SetNodeContext(n, nctx)
	return nil
}

func (n *Refnode) ResolveLabels(ctx scanner.ResolutionContext) error {
	nctx := scanner.GetNodeContext[*RefNodeContext](ctx, n)

	return nctx.Resolve(ctx)
}

func (n *Refnode) Emit(ctx scanner.ResolutionContext) error {
	nctx := scanner.GetNodeContext[*RefNodeContext](ctx, n)
	link, err := nctx.Link(ctx)
	if err != nil {
		return err
	}
	label := func(ctx scanner.ResolutionContext) error {
		abbrev := ""
		if n.abbrev {
			abbrev = nctx.RefInfo.Abbrev()
			if abbrev != "" {
				if n.upper {
					r, i := utf8.DecodeRuneInString(abbrev)
					abbrev = string(unicode.ToTitle(r)) + abbrev[i:]
				}
				abbrev += " "
			}
		}

		label := nctx.RefInfo.Label().Name()
		if label != "" {
			label = abbrev + label
		}
		var c rune = 0x2192

		fmt.Fprintf(ctx.Writer(), "%s%s", string(c), label)
		return nil
	}

	return render.Current.Link(ctx, link, label)
}
