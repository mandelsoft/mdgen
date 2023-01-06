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
	"github.com/mandelsoft/mdgen/statements/anchor"
)

const FIGURE_TYPE = "figure"

func init() {
	scanner.Tokens.RegisterStatement(NewStatement())
}

type Statement struct {
	scanner.BracketedStatement[FigureNode]
}

func NewStatement() scanner.Statement {
	return &Statement{scanner.NewBracketedStatement[FigureNode]("figure", true)}
}

func (s *Statement) Start(p scanner.Parser, e scanner.Element) (scanner.Element, error) {

	tag := ""
	path := ""
	var attrs []string

	tags := e.Tags()
	switch len(e.Tags()) {
	case 0:
		return nil, e.Errorf("image path required")
	case 1:
		path = tags[0]
	case 2:
		path = tags[1]
		tag = tags[0]
	default:
		start := 0
		if strings.Contains(tags[1], "=") {
			path = tags[0]
			start = 1
		} else {
			path = tags[1]
			tag = tags[0]
			start = 2
		}
		for i, t := range tags[start:] {
			if !strings.Contains(t, "=") {
				return nil, e.Errorf("tag %d [%s] is no image attribute", i+start+1, t)
			}
			attrs = append(attrs, t)
		}
	}

	sid := p.State.NextId(FIGURE_TYPE).Id()
	n := NewFigureNode(p.State.Container, p.Document(), e.Location(), sid, tag, path, attrs)
	p.State = p.State.Sub(n)
	p.State.SetLastTag(tag)
	return p.NextElement()
}

////////////////////////////////////////////////////////////////////////////////

type FigureNode = *figurenode

type figurenode struct {
	anchor.AnchorNode
	path  string
	attrs []string
}

func NewFigureNode(p scanner.NodeContainer, d scanner.Document, location scanner.Location, sid scanner.TaggedId, tag string, path string, attrs []string) FigureNode {
	return &figurenode{
		AnchorNode: anchor.NewAnchorNode(p, d, location, sid, tag, false),
		path:       path,
		attrs:      attrs,
	}
}

func (n *figurenode) Print(gap string) {
	fmt.Printf("%sFIGURE %s[%s]: %s (%s)\n", gap, n.Id(), n.Tag(), n.path, strings.Join(n.attrs, ", "))
	n.NodeContainerBase.Print(gap + "  ")
}

func (n *figurenode) Emit(ctx scanner.ResolutionContext) error {
	w := ctx.Writer()

	nctx := scanner.GetNodeContext[*anchor.AnchorNodeContext](ctx, n.AnchorNode)
	info := ctx.GetReferencable(nctx.Id())

	path, err := ctx.DetermineLinkPath(n.Source(), n.path)
	if err != nil {
		return n.Errorf("cannot determine target path: %s", err)
	}

	fmt.Fprintf(w, "<div align=\"center\">\n")
	fmt.Fprintf(w, "<img src=\"%s\" alt=\"%s\" %s/>\n", path, *info.Title(), strings.Join(n.attrs, " "))
	err = n.AnchorNode.Emit(ctx)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "</div>\n")
	return nil
}
