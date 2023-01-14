/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package block

import (
	"fmt"
	"strings"

	"github.com/mandelsoft/mdgen/scanner"
	"github.com/mandelsoft/mdgen/statements/section"
	"github.com/mandelsoft/mdgen/statements/sectionref"
	"github.com/mandelsoft/mdgen/utils"
)

const BLOCK_TYPE = "block"

func init() {
	scanner.Tokens.RegisterStatement(NewStatement(), true, true)
}

type Statement struct {
	scanner.BracketedStatement[BlockNode]
}

func NewStatement() scanner.Statement {
	return &Statement{scanner.NewBracketedStatement[BlockNode]("block", false)}
}

func (s *Statement) Start(p scanner.Parser, e scanner.Element) (scanner.Element, error) {
	var err error
	var tag string

	if e.HasTags() {
		tag, err = e.Tag("block tag")
		if err != nil {
			return nil, err
		}
	}

	if err = scanner.ForbidNesting[sectionref.SectionRefNode]("sectionref", p, e); err != nil {
		return nil, err
	}
	if err = scanner.ForbidNesting[section.Node]("section", p, e); err != nil {
		return nil, err
	}

	if err = scanner.ForbidNestingInTypes(p, e, "arg", "param"); err != nil {
		return nil, err
	}

	sid := p.State.NextId(BLOCK_TYPE).Id()
	n := NewBlockNode(sid, p, e.Location(), tag)
	err = p.State.Container.RegisterBlock(n)
	if err != nil {
		return nil, err
	}
	p.State = p.State.Sub(n, n.name)
	p.State.SubId(BLOCK_TYPE)

	return scanner.ParseElementsUntil(p, func(p scanner.Parser, e scanner.Element) (scanner.Element, error) {
		switch e.Token() {
		case "param":
			return ParseParam(p, n, e)
		}
		return e, nil
	})
}

////////////////////////////////////////////////////////////////////////////////

type BlockNodeContext struct {
	scanner.NodeContextBase[*blocknode]
	params map[string]*scanner.Value
}

var _ scanner.BlockNodeContext = (*BlockNodeContext)(nil)

func NewBlockNodeContext(n *blocknode, ctx scanner.ResolutionContext) *BlockNodeContext {
	return &BlockNodeContext{
		NodeContextBase: scanner.NewNodeContextBase(n, ctx),
		params:          map[string]*scanner.Value{},
	}
}

func (c *BlockNodeContext) Tag() string {
	return c.EffNode().Tag()
}

func (c *BlockNodeContext) Name() string {
	return c.EffNode().name
}

func (c *BlockNodeContext) Inventory() scanner.Inventory {
	return c.EffNode().InventoryContainer
}

func (c *BlockNodeContext) GetDefaultValue(n string) *scanner.Value {
	return c.params[n]
}

func (c *BlockNodeContext) HasParam(n string) bool {
	_, ok := c.params[n]
	return ok
}

func (c *BlockNodeContext) GetParameterNames() []string {
	return utils.StringMapKeys(c.params)
}

func (c *BlockNodeContext) GetNodeSequence() scanner.NodeSequence {
	return c.EffNode().NodeSequence
}

type BlockNode = *blocknode

type blocknode struct {
	scanner.TaggedNodeBase
	scanner.NodeContainerBase
	static scanner.NodeContainer
	name   string
	params map[string]scanner.NodeSequence
}

var _ (scanner.LabelResolver) = (*blocknode)(nil)
var _ (scanner.BlockNode) = (*blocknode)(nil)

func NewBlockNode(sid scanner.TaggedId, p scanner.Parser, location scanner.Location, tag string) BlockNode {
	sep := "/"
	if strings.HasSuffix(p.State.ScopeName(), "#") {
		sep = ""
	}

	b := &blocknode{
		TaggedNodeBase: scanner.NewTaggedNodeBase(sid, tag),
		static:         p.State.Container,
		name:           p.State.ScopeName() + sep + tag,
		params:         map[string]scanner.NodeSequence{},
	}
	b.NodeContainerBase = scanner.NewContainerBase("block", p.Document(), location, scanner.NewInventoryScope(scanner.NewInventory()))
	return b
}

func (n *blocknode) Name() string {
	return n.name
}

func (n *blocknode) Inventory() scanner.Inventory {
	return n.InventoryContainer
}

func (c *blocknode) HasParam(n string) bool {
	_, ok := c.params[n]
	return ok
}

func (n *blocknode) Print(gap string) {
	fmt.Printf("%sBLOCK %s[%s]\n", gap, n.Id(), n.Tag())
	gap += "  "
	fmt.Printf("%sparameters:\n", gap)
	for _, k := range utils.StringMapKeys(n.params) {
		p := n.params[k]
		if p == nil {
			fmt.Printf("%s  %s (no default)\n", gap, k)
		} else {
			fmt.Printf("%s  %s:\n", gap, k)
			p.Print(gap + "    ")
		}
	}
	n.NodeContainerBase.Print(gap + "  ")
}

func (n *blocknode) RegisterAt(s scanner.Scope) error {
	if s.GetNodeContext(n) != nil {
		panic("blocknode already registered")
	}
	nctx := NewBlockNodeContext(n, s.GetContext())
	for k, v := range n.params {
		if v != nil {
			nctx.params[k] = scanner.NewValue(s.GetContext(), v)
		} else {
			nctx.params[k] = nil
		}
	}
	s.SetNodeContext(n, nctx)
	return nil
}

func (n *blocknode) ResolveLabels(ctx scanner.ResolutionContext) error {
	return nil
}

func (n *blocknode) Emit(ctx scanner.ResolutionContext) error {
	return nil
}
