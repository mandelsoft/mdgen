/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package toc

import (
	"fmt"
	"strings"

	"github.com/mandelsoft/mdgen/labels"
	"github.com/mandelsoft/mdgen/scanner"
	"github.com/mandelsoft/mdgen/utils"
)

func init() {
	scanner.Tokens.RegisterStatement(NewStatement())
}

type Statement struct {
	scanner.StatementBase
}

func NewStatement() scanner.Statement {
	return &Statement{scanner.NewStatementBase("toc")}
}

func (s *Statement) Start(p scanner.Parser, e scanner.Element) (scanner.Element, error) {
	typ, err := e.OptionalTag("tag")
	if err != nil {
		return nil, err
	}
	skip := false
	if strings.HasPrefix(typ, "*") {
		typ = typ[1:]
		skip = true
	}
	comps := strings.Split(typ, ":")
	if len(comps) > 2 {
		return nil, e.Errorf("invalid tag: only two tag components possible")
	}
	root := ""
	typ = comps[0]
	if len(comps) > 1 {
		root = comps[1]
	}
	if typ == "" {
		typ = scanner.SECTION_TYPE
	}

	var link *utils.Link
	if root != "" {
		tag := root
		l, err := utils.ParseAbsoluteLink(tag, "", false)
		if err != nil {
			return nil, e.Errorf("%s", err.Error())
		}
		link = &l
	}
	n := NewTOCNode(p.Document(), e.Location(), typ, link, skip)
	p.State.Container.AddNode(n)
	return p.NextElement()
}

////////////////////////////////////////////////////////////////////////////////

type TOCNodeContext struct {
	scanner.NodeContextBase[*tocnode]
	link *utils.Link
}

func NewTOCNodeContext(n *tocnode, ctx scanner.ResolutionContext) (*TOCNodeContext, error) {
	c := &TOCNodeContext{
		NodeContextBase: scanner.NewNodeContextBase(n, ctx),
	}
	if n.link != nil {
		link, err := n.link.Abs(ctx.GetDocument().GetRefPath(), false)
		if err != nil {
			return nil, n.Errorf("%s", err)
		}
		c.link = &link
		return c, err
	}
	return c, nil
}

type TOCNode = *tocnode

type tocnode struct {
	scanner.NodeBase
	typ  string
	link *utils.Link
	skip bool
}

func NewTOCNode(d scanner.Document, location scanner.Location, typ string, link *utils.Link, skip bool) *tocnode {
	return &tocnode{
		NodeBase: scanner.NewNodeBase(d, location),
		typ:      typ,
		link:     link,
		skip:     skip,
	}
}

func (n *tocnode) Print(gap string) {
	fmt.Printf("%sTOC %s\n", gap, n.typ)
}

func (n *tocnode) Register(ctx scanner.ResolutionContext) error {
	nctx, err := NewTOCNodeContext(n, ctx)
	if err != nil {
		return err
	}
	ctx.SetNodeContext(n, nctx)
	return nil
}

func (n *tocnode) Emit(ctx scanner.ResolutionContext) error {
	nctx := scanner.GetNodeContext[*TOCNodeContext](ctx, n)

	var list []TocEntry
	if nctx.link != nil {
		info := ctx.GetLinkInfo(*nctx.link)
		if info.Label().Id().Type() != n.typ {
			return n.Errorf("label %s is not of type %s, but %s", nctx.link, n.typ, info.Label().Id().Type())
		}
		list = TreeTOCIds(info.Context(), n.typ)
		prefix := info.Label().Id()
		for i := 0; i < len(list); i++ {
			if !prefix.IsPrefix(list[i].Id()) {
				list = append(list[:i], list[i+1:]...)
				i--
			}
		}
	} else {
		list = DocTOCIds(ctx, n.typ)
	}

	if len(list) == 0 {
		return nil
	}
	minlvl := list[0].Level()
	if n.skip {
		cnt := 0
		for _, e := range list {
			if e.Level() == minlvl {
				cnt++
			}
		}
		fmt.Printf("found %d top level %s entries\n", cnt, n.typ)
		if cnt == 1 {
			list = list[1:]
			if len(list) == 0 {
				return n.Errorf("no toc entries found")
			}
			minlvl = list[0].Level()
		}
	}

	w := ctx.Writer()
	for _, e := range list {
		info := e.info
		rt := info.Title()
		if rt == nil {
			return n.Errorf("unresolved title for %s:%s", info.GetRefPath(), info.Anchors()[0])
		}
		title := *rt
		link, err := ctx.DetermineLink(e.info.Link())
		if err != nil {
			return n.Errorf("cannot resolve link for %s %s: %s", n.typ, info.Label().Id(), err.Error())
		}
		if info.Label().Name() != "" {
			title = info.Label().Name() + " " + title
		}
		gap := fmt.Sprintf("%*s", (e.Level()-minlvl)*2+2, "")
		fmt.Fprintf(w, "%s [%s](%s)<br>\n", strings.ReplaceAll(gap, " ", "&nbsp;&nbsp;"), title, link)
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////

type TocEntry struct {
	info  scanner.TreeLabelInfo
	id    labels.LabelId
	level int
}

func (e TocEntry) Level() int {
	return e.level
}
func (e TocEntry) Id() scanner.TaggedId {
	return e.id
}

func DocTOCIds(ctx scanner.ResolutionContext, typ string) []TocEntry {
	return prepareTOCIds(ctx.GetLabelInfosForType(typ))
}

func TreeTOCIds(ctx scanner.ResolutionContext, typ string) []TocEntry {
	return prepareTOCIds(ctx.GetIdsForTypeInTree(typ))
}

func prepareTOCIds(m map[labels.LabelId]scanner.TreeLabelInfo) []TocEntry {
	var list labels.StructureList

	ctxs := map[labels.LabelId]scanner.TreeLabelInfo{}
	for id, info := range m {
		s := labels.Structured(id)
		ctxs[s.Id()] = info
		list = append(list, s)
	}
	list.Sort()

	var result []TocEntry
	for _, id := range list {
		result = append(result, TocEntry{id: id.Id(), level: id.Level(), info: ctxs[id.Id()]})
	}
	return result
}
