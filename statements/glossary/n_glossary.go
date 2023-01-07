/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package glossary

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/mandelsoft/mdgen/scanner"
	"github.com/mandelsoft/mdgen/statements/termdef"
	"github.com/mandelsoft/mdgen/utils"
)

func init() {
	scanner.Tokens.RegisterStatement(NewStatement())
}

type Statement struct {
	scanner.StatementBase
}

func NewStatement() scanner.Statement {
	return &Statement{scanner.NewStatementBase("glossary")}
}

func (s *Statement) Start(p scanner.Parser, e scanner.Element) (scanner.Element, error) {
	var err error
	tag := ""
	if e.HasTags() {
		tag, err = e.Tag("tag")
		if err != nil {
			return nil, err
		}
	}
	n := NewGlossaryNode(p.Document(), e.Location(), tag)
	p.State.Container.AddNode(n)
	return p.NextElement()
}

////////////////////////////////////////////////////////////////////////////////

const InfoKey = "glossary"

type glossaryContext struct {
	*scanner.StaticContext
}

func newContext(orig scanner.ResolutionContext, ctx scanner.ResolutionContext) scanner.ResolutionContext {
	return &glossaryContext{
		StaticContext: scanner.NewStaticContext(ctx, orig),
	}
}

func (c *glossaryContext) Info(key string) interface{} {
	if key == InfoKey {
		return true
	}
	return c.StaticContext.Info(key)
}

////////////////////////////////////////////////////////////////////////////////

type GlossaryNode = *glossarynode

type glossarynode struct {
	scanner.NodeBase
	tag string
}

func NewGlossaryNode(d scanner.Document, location scanner.Location, tag string) *glossarynode {
	return &glossarynode{
		NodeBase: scanner.NewNodeBase(d, location),
		tag:      tag,
	}
}

func (n *glossarynode) Print(gap string) {
	fmt.Printf("%sGLOSSARY %s\n", gap, n.tag)
}

type Glossary map[string]map[string]*termdef.TermDefNodeContext

func writerBar(w scanner.Writer, header []string, g Glossary) {
	for _, h := range header {
		if g[h] == nil {
			fmt.Fprintf(w, "%s &nbsp;", h)
		} else {
			r, _ := utf8.DecodeRuneInString(h)
			l := unicode.ToLower(r)
			fmt.Fprintf(w, "[%s](%s) &nbsp;", h, string(l))
		}
	}
	fmt.Fprintf(w, "\n\n")
}

func (n *glossarynode) Register(ctx scanner.ResolutionContext) error {
	return nil
}

func (n *glossarynode) ResolveLabels(ctx scanner.ResolutionContext) error {
	tags := ctx.GetGlobalTags(termdef.GT_TERM)
	for _, c := range tags {
		nctx := c.(*termdef.TermDefNodeContext)
		if !strings.HasPrefix(nctx.Term().Tag(), n.tag) {
			continue
		}
		err := nctx.GetNodeSequence().Register(newContext(ctx, nctx.GetContext()))
		if err != nil {
			return err
		}
	}
	for _, c := range tags {
		nctx := c.(*termdef.TermDefNodeContext)
		if !strings.HasPrefix(nctx.Term().Tag(), n.tag) {
			continue
		}
		err := nctx.GetNodeSequence().ResolveLabels(newContext(ctx, nctx.GetContext()))
		if err != nil {
			return err
		}
	}
	return nil
}

func (n *glossarynode) ResolveValues(ctx scanner.ResolutionContext) error {
	tags := ctx.GetGlobalTags(termdef.GT_TERM)
	for _, c := range tags {
		nctx := c.(*termdef.TermDefNodeContext)
		if !strings.HasPrefix(nctx.Term().Tag(), n.tag) {
			continue
		}
		err := nctx.GetNodeSequence().ResolveValues(newContext(ctx, nctx.GetContext()))
		if err != nil {
			return err
		}
	}
	return nil
}

func (n *glossarynode) Emit(ctx scanner.ResolutionContext) error {
	glossary := Glossary{}
	caser := cases.Title(language.AmericanEnglish)

	for _, c := range ctx.GetGlobalTags(termdef.GT_TERM) {
		nctx := c.(*termdef.TermDefNodeContext)
		if !strings.HasPrefix(nctx.Term().Tag(), n.tag) {
			continue
		}
		t := caser.String(nctx.Term().Singular())
		r, _ := utf8.DecodeRuneInString(t)
		m := glossary[string(r)]
		if m == nil {
			m = map[string]*termdef.TermDefNodeContext{}
			glossary[string(r)] = m
		}
		m[nctx.Term().FormatSingular()] = nctx
	}

	header := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
outer:
	for _, k := range utils.StringMapKeys(glossary) {
		fmt.Printf("letter %s\n", k)
		for _, h := range header {
			if h == k {
				continue outer
			}
		}
		header = append(header, k)
	}
	w := ctx.Writer()

	writerBar(w, header, glossary)

	for _, h := range header {
		m := glossary[h]
		if len(m) == 0 {
			continue
		}
		fmt.Fprintf(w, "## %s\n\n", h)
		keys := utils.StringMapKeys(m)

		for _, k := range keys {
			nctx := m[k]
			link, err := ctx.DetermineLink(nctx.GetLink())
			if err != nil {
				return n.Errorf("term %s: %s", nctx.Term().Tag(), err)
			}
			tag := nctx.Term().Tag()
			if strings.HasPrefix(tag, "/") {
				tag = tag[1:]
			}
			txt := caser.String(k)
			if nctx.Term().IsFormatted() {
				txt = nctx.Term().FormatSingular()
			}
			fmt.Fprintf(w, "### [%s](%s)<a id=\"%s\"/>\n", txt, link, "glossary/"+tag)
			err = nctx.GetNodeSequence().Emit(newContext(ctx, nctx.GetContext()))
			if err != nil {
				return n.Errorf("term %s: %s", nctx.Term().Tag(), err)
			}
			fmt.Fprintf(w, "\n")
		}
	}
	return nil
}
