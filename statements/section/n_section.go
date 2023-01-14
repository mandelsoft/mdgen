/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package section

import (
	"fmt"

	"github.com/mandelsoft/mdgen/scanner"
	"github.com/mandelsoft/mdgen/statements/sectionref"
	"github.com/mandelsoft/mdgen/statements/subrange"
)

func init() {
	scanner.Tokens.RegisterStatement(NewStatement())
}

type Statement struct {
	*subrange.Statement[*node]
}

func NewStatement() *Statement {
	s := &Statement{}
	s.Statement = subrange.NewStatement[*node]("section", s.newNode, true)
	return s
}

func (s *Statement) NewNode(p scanner.Parser, e scanner.Element, tag string) Node {
	n := &node{}
	n.Node = s.Statement.NewNode(n, p, e.Location(), scanner.SECTION_TYPE, tag)
	return n
}

func (s *Statement) newNode(p scanner.Parser, e scanner.Element) (*node, error) {
	tag, err := e.OptionalTag("tag")
	if err != nil {
		return nil, e.Errorf("%s", err.Error())
	}
	if err = scanner.ForbidNesting[sectionref.SectionRefNode]("sectionref", p, e); err != nil {
		return nil, err
	}
	return s.NewNode(p, e, tag), nil
}

////////////////////////////////////////////////////////////////////////////////

type Node = *node

type node struct {
	subrange.Node
}

func (n *node) Emit(ctx scanner.ResolutionContext) error {
	w := ctx.Writer()
	nctx := scanner.GetNodeContext[*subrange.NodeContext](ctx, n)

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
	return n.GetContent().Emit(nctx.GetContext())
}
