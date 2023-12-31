/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package execute

import (
	"fmt"
	"os/exec"

	"github.com/mandelsoft/filepath/pkg/filepath"

	"github.com/mandelsoft/mdgen/scanner"
	"github.com/mandelsoft/mdgen/statements/include"
)

func init() {
	scanner.Tokens.RegisterStatement(NewStatement())
}

type Statement struct {
	scanner.StatementBase
}

func NewStatement() scanner.Statement {
	return &Statement{scanner.NewStatementBase("execute")}
}

func (s *Statement) Start(p scanner.Parser, e scanner.Element) (scanner.Element, error) {
	if !e.HasTags() {
		return nil, e.Errorf("command missing")
	}

	n := NewExecuteNode(p.State.Container, p.Document(), e.Location(), e.Tags())
	p.State.Container.AddNode(n)
	return scanner.ParseElementsUntil(p, func(p scanner.Parser, e scanner.Element) (scanner.Element, error) {
		switch e.Token() {
		case "range":
			return include.ParseRange(p, &n.ContentHandler, e)
		case "pattern":
			return include.ParsePattern(p, &n.ContentHandler, e)
		case "filter":
			return include.ParseFilter(p, &n.ContentHandler, e)
		}
		return e, nil
	})
}

////////////////////////////////////////////////////////////////////////////////

type ExecuteNodeContext struct {
	scanner.NodeContextBase[*executenode]
	command []string
	dir     string
}

func NewExecuteNodeContext(n *executenode, ctx scanner.ResolutionContext, command []string, dir string) *ExecuteNodeContext {
	return &ExecuteNodeContext{
		NodeContextBase: scanner.NewNodeContextBase(n, ctx),
		command:         command,
		dir:             dir,
	}
}

type ExecuteNode = *executenode

type executenode struct {
	scanner.NodeBase
	tags []string

	include.ContentHandler
}

func NewExecuteNode(p scanner.NodeContainer, d scanner.Document, location scanner.Location, tags []string) ExecuteNode {
	return &executenode{
		NodeBase: scanner.NewNodeBase(d, location),
		tags:     tags,
	}
}

func (n *executenode) Print(gap string) {
	fmt.Printf("%sEXECUTE %v\n", gap, n.tags)
}

func (n *executenode) Register(ctx scanner.ResolutionContext) error {
	dir := filepath.Dir(n.Source())
	cmd := n.tags

	nctx := NewExecuteNodeContext(n, ctx, cmd, dir)
	ctx.SetNodeContext(n, nctx)
	return nil
}

func (n *executenode) Emit(ctx scanner.ResolutionContext) error {
	nctx := scanner.GetNodeContext[*ExecuteNodeContext](ctx, n)

	cmd := exec.Command(nctx.command[0], nctx.command[1:]...)
	cmd.Dir = nctx.dir
	data, err := cmd.Output()
	if err != nil {
		return n.Errorf("cannot execute %v: %s", nctx.command, err)
	}

	data, err = n.Process(data)
	if err != nil {
		return n.Errorf("%v: %s", n.tags, err)
	}

	fmt.Fprintf(ctx.Writer(), "%s\n", string(data))
	return nil
}
