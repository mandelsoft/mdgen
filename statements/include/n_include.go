/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package include

import (
	"fmt"
	"os"

	"github.com/mandelsoft/filepath/pkg/filepath"

	"github.com/mandelsoft/mdgen/scanner"
)

func init() {
	scanner.Tokens.RegisterStatement(NewStatement())
}

type Statement struct {
	scanner.StatementBase
}

func NewStatement() scanner.Statement {
	return &Statement{scanner.NewStatementBase("include")}
}

func (s *Statement) Start(p scanner.Parser, e scanner.Element) (scanner.Element, error) {
	tag, err := e.Tag("include path")
	if err != nil {
		return nil, err
	}

	n := NewIncludeNode(p.State.Container, p.Document(), e.Location(), tag)
	p.State.Container.AddNode(n)
	return scanner.ParseElementsUntil(p, func(p scanner.Parser, e scanner.Element) (scanner.Element, error) {
		switch e.Token() {
		case "range":
			return ParseRange(p, &n.ContentHandler, e)
		case "pattern":
			return ParsePattern(p, &n.ContentHandler, e)
		case "filter":
			return ParseFilter(p, &n.ContentHandler, e)
		}
		return e, nil
	})
}

////////////////////////////////////////////////////////////////////////////////

type IncludeNodeContext struct {
	scanner.NodeContextBase[*includenode]
	file string
}

func NewIncludeNodeContext(n *includenode, ctx scanner.ResolutionContext, file string) *IncludeNodeContext {
	return &IncludeNodeContext{
		NodeContextBase: scanner.NewNodeContextBase(n, ctx),
		file:            file,
	}
}

type IncludeNode = *includenode

type includenode struct {
	scanner.NodeBase
	tag string

	ContentHandler
}

type ContentHandler struct {
	extract Extractor
	filter  Filter
}

func (n *ContentHandler) Process(data []byte) ([]byte, error) {
	var err error

	if n.extract != nil {
		data, err = n.extract.Extract(data)
		if err != nil {
			return nil, fmt.Errorf("cannot extract data: %w", err)
		}
	}
	if n.filter != nil {
		data, err = n.filter.Filter(data)
		if err != nil {
			return nil, fmt.Errorf("cannot filter data: %w", err)
		}
	}
	return data, nil
}

type Extractor interface {
	Extract(data []byte) ([]byte, error)
}

type Filter interface {
	Filter(data []byte) ([]byte, error)
}

func NewIncludeNode(p scanner.NodeContainer, d scanner.Document, location scanner.Location, tag string) IncludeNode {
	return &includenode{
		NodeBase: scanner.NewNodeBase(d, location),
		tag:      tag,
	}
}

func (n *includenode) Print(gap string) {
	fmt.Printf("%sINCLUDE %s\n", gap, n.tag)
}

func (n *includenode) Register(ctx scanner.ResolutionContext) error {
	file := n.tag
	if !filepath.IsAbs(n.tag) {
		file = filepath.Join(filepath.Dir(n.Source()), n.tag)
	}
	_, err := os.ReadFile(file)
	if err != nil {
		return n.Errorf("cannot read include file %q: %s", n.tag, err)
	}
	nctx := NewIncludeNodeContext(n, ctx, file)
	ctx.SetNodeContext(n, nctx)
	return nil
}

func (n *includenode) Emit(ctx scanner.ResolutionContext) error {
	nctx := scanner.GetNodeContext[*IncludeNodeContext](ctx, n)
	data, err := os.ReadFile(nctx.file)
	if err != nil {
		return n.Errorf("cannot read include file %q: %s", n.tag, err)
	}

	data, err = n.Process(data)
	if err != nil {
		return n.Errorf("%q: %s", n.tag, err)
	}

	fmt.Fprintf(ctx.Writer(), "%s\n", string(data))
	return nil
}
