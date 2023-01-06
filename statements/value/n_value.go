/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package value

import (
	"fmt"
	"strings"

	"github.com/mandelsoft/mdgen/scanner"
)

func init() {
	scanner.Tokens.RegisterStatement(NewStatement())
}

type Statement struct {
	scanner.StatementBase
}

func NewStatement() scanner.Statement {
	return &Statement{scanner.NewStatementBase("value")}
}

func (s *Statement) Start(p scanner.Parser, e scanner.Element) (scanner.Element, error) {
	tag, err := e.Tag("parameter name")
	if err != nil {
		return nil, err
	}
	attr := false
	if strings.HasPrefix(tag, "*") {
		attr = true
		tag = strings.TrimSpace(tag[1:])
	}
	if len(tag) == 0 {
		return nil, e.Errorf("empty tag")
	}

	if attr {
		if !scanner.ContextAttrs[tag] {
			return nil, e.Errorf("unknown scope attribute %q", tag)
		}
	} else {
		checkParam := func(b scanner.BlockNode) (bool, error) {
			if b == nil {
				return true, e.Errorf("parameter %q not defined in static scopes", tag)
			}
			return b.HasParam(tag), nil
		}

		if err = scanner.RequireNesting[scanner.BlockNode]("block", p, e, checkParam); err != nil {
			return nil, err
		}
	}
	n := NewValueNode(p.Document(), e.Location(), tag, attr)
	p.State.Container.AddNode(n)
	return p.NextElement()
}

////////////////////////////////////////////////////////////////////////////////

type ValueNodeContext struct {
	scanner.NodeContextBase[*valuenode]
	value *scanner.Value
	ctx   scanner.ResolutionContext
}

func NewValueNodeContext(n *valuenode, ctx scanner.ResolutionContext, v *scanner.Value) *ValueNodeContext {
	var vctx scanner.ResolutionContext
	if v != nil {
		vctx = scanner.NewStaticContext(v.GetContext().GetScope(), ctx)
	}
	return &ValueNodeContext{
		NodeContextBase: scanner.NewNodeContextBase(n, ctx),
		value:           v,
		ctx:             vctx,
	}
}

type ValueNode = *valuenode

type valuenode struct {
	scanner.NodeBase
	tag  string
	attr bool
}

func NewValueNode(d scanner.Document, location scanner.Location, tag string, attr bool) ValueNode {
	return &valuenode{
		NodeBase: scanner.NewNodeBase(d, location),
		tag:      tag,
		attr:     attr,
	}
}

func (n *valuenode) Print(gap string) {
	attr := ""
	if n.attr {
		attr = " (attr)"
	}
	fmt.Printf("%sVALUE %s%s\n", gap, n.tag, attr)
}

func (n *valuenode) Register(ctx scanner.ResolutionContext) error {
	var v *scanner.Value
	var err error
	if !n.attr {
		v = ctx.LookupValue(n.tag)
		if v == nil {
			n.Errorf("parameter %q not defined", n.tag)
		}
	}
	nctx := NewValueNodeContext(n, ctx, v)
	ctx.SetNodeContext(n, nctx)
	if !n.attr {
		err = nctx.value.Register(nctx.ctx)
	}
	return err
}

func (n *valuenode) ResolveLabels(ctx scanner.ResolutionContext) error {

	if !n.attr {
		nctx := scanner.GetNodeContext[*ValueNodeContext](ctx, n)
		return nctx.value.ResolveLabels(nctx.ctx)
	}
	return nil
}

func (n *valuenode) ResolveValues(ctx scanner.ResolutionContext) error {
	if !n.attr {
		nctx := scanner.GetNodeContext[*ValueNodeContext](ctx, n)
		return nctx.value.ResolveValues(nctx.ctx)
	}
	return nil
}

func (n *valuenode) Emit(ctx scanner.ResolutionContext) error {
	if !n.attr {
		nctx := scanner.GetNodeContext[*ValueNodeContext](ctx, n)
		return nctx.value.Emit(scanner.NewDelegationContext(ctx, nctx.ctx))
	}
	w := ctx.Writer()
	fmt.Fprintf(w, "%s", scanner.GetContextAttr(n.tag, ctx))
	return nil
}
