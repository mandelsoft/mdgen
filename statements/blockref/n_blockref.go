/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package blockref

import (
	"fmt"
	"strings"

	"github.com/mandelsoft/mdgen/scanner"
	utils2 "github.com/mandelsoft/mdgen/utils"
)

func init() {
	scanner.Tokens.RegisterStatement(NewStatement(), true, true)

}

type Statement struct {
	scanner.StatementBase
}

func NewStatement() scanner.Statement {
	return &Statement{scanner.NewStatementBase("blockref")}
}

func (s *Statement) Start(p scanner.Parser, e scanner.Element) (scanner.Element, error) {
	ref, err := e.Tag("block reference")
	if err != nil {
		return nil, err
	}
	tag := ref
	comps := strings.Split(tag, ":")
	switch len(comps) {
	case 1:
	case 2:
		ref = comps[1]
		tag = comps[0]
	default:
		return nil, e.Errorf("tags must be of [<tag>:]<ref>")
	}

	r, err := utils2.ParseLink(ref, true)
	if err != nil {
		return nil, e.Errorf("%s", err.Error())
	}
	n := NewBlockRefNode(p.Document(), e.Location(), r, tag)
	p.State.Container.AddNode(n)
	return scanner.ParseElementsUntil(p, func(p scanner.Parser, e scanner.Element) (scanner.Element, error) {
		switch e.Token() {
		case "arg":
			return ParseArg(p, n, e)
		}
		return e, nil
	})
}

////////////////////////////////////////////////////////////////////////////////

type blockRefContext struct {
	scanner.ResolutionContext

	nctx      *BlockRefNodeContext
	blockctx  scanner.BlockNodeContext
	callstack scanner.CallStack
	writer    scanner.Writer
}

func (c *blockRefContext) GetContextNodeContext() scanner.NodeContext {
	return c.nctx
}

func (c *blockRefContext) CallStack() scanner.CallStack {
	return c.callstack
}

var _ scanner.ResolutionContext = (*blockRefContext)(nil)

func BlockRefContext(ctx scanner.ResolutionContext, static scanner.Scope, bctx scanner.BlockNodeContext, nctx *BlockRefNodeContext) (*blockRefContext, error) {
	n := scanner.NewStaticContext(nil, ctx)
	c := &blockRefContext{
		ResolutionContext: n,
		nctx:              nctx,
		blockctx:          bctx,
	}
	scope := scanner.NewScope(ctx.GetScope(), static, c, bctx.Inventory(), nctx.tag)
	n.SetScope(scope)
	err := ctx.GetScope().AddSubScope(nctx.tag, scope)
	if err != nil {
		return nil, nctx.Errorf("%s", err.Error())
	}

	callstack, cycle := ctx.CallStack().Add(bctx.Name(), nctx.Location())
	if len(cycle) > 0 {
		return nil, nctx.Errorf("%s: recursive use of block %q (%s)", ctx.CallStack(), nctx.EffNode().ref, strings.Join(cycle, "->"))
	}
	c.callstack = callstack
	return c, nil
}

func (c *blockRefContext) Writer() scanner.Writer {
	return c.writer
}

func (c *blockRefContext) SetWriter(w scanner.Writer) {
	c.writer = scanner.NewIndentWriter(w)
}

////////////////////////////////////////////////////////////////////////////////

type BlockRefNodeContext struct {
	scanner.NodeContextBase[*blockrefnode]
	ctx *blockRefContext
	tag string
}

func NewBlockRefNodeContext(n *blockrefnode, ctx scanner.ResolutionContext, tag string) *BlockRefNodeContext {
	return &BlockRefNodeContext{
		NodeContextBase: scanner.NewNodeContextBase(n, ctx),
		ctx:             nil,
		tag:             tag,
	}
}

type BlockRefNode = *blockrefnode

type blockrefnode struct {
	scanner.NodeBase
	ref  utils2.Link
	tag  string
	args map[string]scanner.NodeSequence
}

func NewBlockRefNode(d scanner.Document, location scanner.Location, ref utils2.Link, tag string) BlockRefNode {
	return &blockrefnode{
		NodeBase: scanner.NewNodeBase(d, location),
		ref:      ref,
		tag:      tag,
		args:     map[string]scanner.NodeSequence{},
	}
}

func (n *blockrefnode) Print(gap string) {
	fmt.Printf("%sBLOCKREF %s[%s]\n", gap, n.tag, n.ref)
	for a, s := range n.args {
		fmt.Printf("%s  arg: %s\n", gap, a)
		s.Print(gap + "    ")
	}
}

func (n *blockrefnode) Register(ctx scanner.ResolutionContext) error {
	ref, err := n.ref.Abs(n.GetDocument().GetRefPath(), false)
	if err != nil {
		return err
	}
	b, s := ctx.LookupBlock(ref)
	if b == nil {
		return n.Errorf("block %q not found", n.ref)
	}

	tag := n.tag
	extend := false
	if tag == "" {
		tag = b.Tag()
		extend = true
	}
	name, err := ctx.GetScope().NextSubScopeName(tag, extend)
	if err != nil {
		return n.Errorf("%s", err)
	}
	nctx := NewBlockRefNodeContext(n, ctx, name)
	bctx, err := BlockRefContext(ctx, s, b, nctx)
	nctx.ctx = bctx
	if err != nil {
		return err
	}

	for pn := range n.args {
		if !b.HasParam(pn) {
			return n.Errorf("unknown parameter %q", pn)
		}
	}
	for _, pn := range b.GetParameterNames() {
		def := b.GetDefaultValue(pn)
		var value *scanner.Value
		if v, ok := n.args[pn]; ok {
			value = scanner.NewValue(ctx, v)
		} else {
			if def == nil {
				return n.Errorf("missing argument for parameter %q", pn)
			}
			value = def
		}
		bctx.SetValue(pn, value)
	}

	ctx.SetNodeContext(n, nctx)
	return b.GetNodeSequence().Register(bctx)
}

func (n *blockrefnode) ResolveLabels(ctx scanner.ResolutionContext) error {
	nctx := scanner.GetNodeContext[*BlockRefNodeContext](ctx, n)
	return nctx.ctx.blockctx.GetNodeSequence().ResolveLabels(nctx.ctx)
}

func (n *blockrefnode) ResolveValues(ctx scanner.ResolutionContext) error {
	nctx := scanner.GetNodeContext[*BlockRefNodeContext](ctx, n)
	return nctx.ctx.blockctx.GetNodeSequence().ResolveValues(nctx.ctx)
}

func (n *blockrefnode) Emit(ctx scanner.ResolutionContext) error {
	nctx := scanner.GetNodeContext[*BlockRefNodeContext](ctx, n)
	nctx.ctx.SetWriter(ctx.Writer())
	return nctx.ctx.blockctx.GetNodeSequence().Emit(nctx.ctx)
}

func (n *blockrefnode) EvaluateStatic(ctx scanner.ResolutionContext) error {
	nctx := scanner.GetNodeContext[*BlockRefNodeContext](ctx, n)
	nctx.ctx.SetWriter(ctx.Writer())
	return nctx.ctx.blockctx.GetNodeSequence().EvaluateStatic(nctx.ctx)
}
