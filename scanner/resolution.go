/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package scanner

import (
	"bytes"
	"fmt"
	"path"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/mandelsoft/mdgen/labels"
	utils2 "github.com/mandelsoft/mdgen/utils"
)

type Resolver func(n Node) error

type Register interface {
	Register(ResolutionContext) error
}

type LabelResolver interface {
	ResolveLabels(ResolutionContext) error
}

type ValueResolver interface {
	ResolveValues(ResolutionContext) error
}

type Emitter interface {
	Emit(ResolutionContext) error
}

type StaticEvaluator interface {
	EvaluateStatic(ResolutionContext) error
}

func Resolve[R any](ctx ResolutionContext) Resolver {
	return func(n Node) error {
		if _, ok := n.(R); ok {
			var i *R
			name := reflect.TypeOf(i).Elem().Method(0).Name
			m := reflect.ValueOf(n).MethodByName(name)
			res := m.Call([]reflect.Value{reflect.ValueOf(ctx)})
			err := res[0].Interface()
			if err != nil {
				return err.(error)
			}
		}
		return nil
	}
}

type NodeContext interface {
	Located
	GetNode() Node
	GetDocument() Document
}

type located = Located
type NodeContextBase[N Node] struct {
	located
	dyndoc Document
}

func NewNodeContextBase[N Node](n N, ctx ResolutionContext) NodeContextBase[N] {
	return NodeContextBase[N]{
		located: n,
		dyndoc:  ctx.GetDocument(),
	}
}

func (n *NodeContextBase[N]) GetNode() Node {
	return n.EffNode()
}

func (n *NodeContextBase[N]) EffNode() N {
	return n.located.(N)
}

func (n *NodeContextBase[N]) GetDocument() Document {
	return n.dyndoc
}

type LabeledNodeContext interface {
	NodeContext
	Abbrev() string
	Label() labels.Label
	Title() *string
	Id() TaggedId
}

type LabeledNodeContextBase[N TaggedNode] struct {
	NodeContextBase[N]
	rule       labels.Rule
	id         TaggedId
	nr         NumberRange
	explicit   bool
	titlenodes NodeSequence
	tags       []string
	hlabel     HierarchyLabel
	abbrev     string
	title      *string
}

func NewLabeledNodeContextBase[N TaggedNode](n N, ctx ResolutionContext, titlenodes NodeSequence) (*LabeledNodeContextBase[N], error) {
	typ := n.Id().Type()
	rule := ctx.NextId(typ)

	var tags []string
	var err error
	var tag string
	explicit := false
	if n.Tag() != "" {
		tag, explicit, err = EvaluateTag(ctx, n.Tag())
		if err != nil {
			return nil, n.Error(err)
		}
		tags = []string{tag, rule.Id().String()}
	} else {
		tags = []string{rule.Id().String()}
	}
	nctx := &LabeledNodeContextBase[N]{
		NodeContextBase: NewNodeContextBase[N](n, ctx),
		id:              rule.Id(),
		rule:            rule,
		titlenodes:      titlenodes,
		explicit:        explicit,
		tags:            tags,
	}

	ctx.RequestNumberRange(typ)
	_, err = ctx.RegisterReferencable(nctx, tags, explicit)
	if err != nil {
		return nil, err
	}
	if titlenodes != nil {
		err = titlenodes.Register(ctx)
		if err != nil {
			return nil, err
		}
	}
	return nctx, nil
}

func (c *LabeledNodeContextBase[N]) IdRule() labels.Rule {
	return c.rule
}

func (c *LabeledNodeContextBase[N]) NumberRange() NumberRange {
	return c.nr
}

func (c *LabeledNodeContextBase[N]) RegisterReferencable(ctx ResolutionContext) (RefInfo, error) {
	return ctx.RegisterReferencable(c, c.Tags(), c.explicit)
}

func (c *LabeledNodeContextBase[N]) ResolveLabels(ctx ResolutionContext) error {
	if c.titlenodes != nil {
		err := c.titlenodes.ResolveLabels(ctx)
		if err != nil {
			return err
		}
	}
	c.nr = ctx.GetNumberRange(c.id.Type())
	c.abbrev = c.nr.Abbrev()
	c.hlabel = c.nr.Next()
	return nil
}

func (c *LabeledNodeContextBase[N]) ResolveValues(ctx ResolutionContext) error {
	title := ""
	if c.titlenodes != nil {
		err := c.titlenodes.ResolveValues(ctx)
		if err != nil {
			return err
		}

		buf := NewBufferContext(ctx)
		err = c.titlenodes.Emit(buf)
		if err != nil {
			ctx.RegisterUnresolved(c, err)
			return nil
		}
		title = strings.TrimSpace(buf.String())
		i := strings.Index(title, "\n")
		if i >= 0 {
			return c.Errorf("resolved title contains newline")
		}
	}
	c.title = &title
	return nil
}

func (c *LabeledNodeContextBase[N]) Id() TaggedId {
	return c.id
}

func (c *LabeledNodeContextBase[N]) Tags() []string {
	return c.tags
}

func (c *LabeledNodeContextBase[N]) Label() labels.Label {
	return c.hlabel.Label()
}

func (c *LabeledNodeContextBase[N]) Abbrev() string {
	return c.abbrev
}

func (c *LabeledNodeContextBase[N]) Title() *string {
	return c.title
}

func (c *LabeledNodeContextBase[N]) EmitAnchors(ctx ResolutionContext) {
	info := ctx.GetReferencable(c.Id())
	w := ctx.Writer()
	fmt.Fprintf(w, "\n")
	if anchors := info.Anchors(); len(anchors) > 0 {
		fmt.Fprintf(w, "<a/>") // without at least two <a> github cannot render sections anymore
		for _, a := range anchors {
			fmt.Fprintf(w, "<a id=\"%s\"/>", a)
		}
		fmt.Fprintf(w, "\n")
	}
}

func (c *LabeledNodeContextBase[N]) EmitTitle(ctx ResolutionContext) {
	w := ctx.Writer()
	if *c.Title() == "" {
		return
	}

	info := ctx.GetReferencable(c.Id())
	abbrev := c.Abbrev()
	if abbrev != "" {
		r, i := utf8.DecodeRuneInString(abbrev)
		abbrev = string(unicode.ToTitle(r)) + abbrev[i:] + " "
	}

	label := info.Label().Name()
	if label != "" {
		label = abbrev + label + ": "
	}
	fmt.Fprintf(w, " %s%s\n</br></br>\n", label, *c.Title())
}

type NodeContexts map[Node]NodeContext

func (i NodeContexts) GetNodeContext(n Node) NodeContext {
	// fmt.Printf("%p: get node context for %T[%p](%s) -> %p\n", i, n, n, n.Location(), i[n])
	return i[n]
}

func (i NodeContexts) SetNodeContext(n Node, v NodeContext) {
	// fmt.Printf("%p: set node context for %T[%p](%s) to %p\n", i, n, n, n.Location(), v)
	i[n] = v
}

type LinkingNode interface {
	Node
	GetLink() utils2.Link
}

type LinkContextInfoNode[N Node] struct {
	NodeContextBase[N]
	link utils2.Link
	RefInfo
}

func NewLinkContextInfoNode[N LinkingNode](n N, ctx ResolutionContext) (LinkContextInfoNode[N], error) {
	var err error
	c := LinkContextInfoNode[N]{
		NodeContextBase: NewNodeContextBase(n, ctx),
	}
	c.link, err = n.GetLink().Abs(ctx.GetDocument().GetRefPath(), false)
	if err != nil {
		return c, c.Errorf("%s", err)
	}
	return c, nil
}

func (c *LinkContextInfoNode[N]) GetLink() utils2.Link {
	return c.link
}

func (c *LinkContextInfoNode[N]) Resolve(ctx ResolutionContext) error {
	ri := ctx.LookupReferencable(c.link)
	if ri == nil {
		return c.Errorf("cannot resolve link %q", c.link)
	}
	c.RefInfo = ri
	return nil
}

func (c *LinkContextInfoNode[N]) Link(ctx ResolutionContext) (string, error) {
	link, err := ctx.DetermineLink(LinkFor(c, c.link.Anchor()))
	if err != nil {
		return link, c.Errorf("%s", err)
	}
	return link, nil
}

type TreeLabelInfo interface {
	RefInfo
	Context() ResolutionContext
	LabelId() labels.LabelId
	Link() utils2.Link
}

func NewTreeLabelInfo(info RefInfo, ctx ResolutionContext) TreeLabelInfo {
	return &treeLabelInfo{
		RefInfo: info,
		ctx:     ctx,
	}
}

type treeLabelInfo struct {
	RefInfo
	ctx ResolutionContext
}

func (i *treeLabelInfo) Context() ResolutionContext {
	return i.ctx
}

func (i *treeLabelInfo) LabelId() labels.LabelId {
	return i.Label().Id()
}

func (i *treeLabelInfo) Link() utils2.Link {
	return NewLink(i.ctx, i.Anchors()...)
}

type RefInfo interface {
	GetRefPath() string
	Anchors() []string
	Label() labels.Label
	Abbrev() string
	Title() *string
}

type ResolvedRef interface {
	RefInfo
	Anchor() string
	Context() ResolutionContext
}

func LinkFor(ri RefInfo, preferred ...string) utils2.Link {
	anchor := ""
	if len(ri.Anchors()) > 0 {
		anchor = ri.Anchors()[0]
	}
	if rr, ok := ri.(ResolvedRef); ok {
		anchor = rr.Anchor()
	} else {
		if len(preferred) > 0 {
			for _, a := range ri.Anchors() {
				if a == preferred[0] || strings.HasSuffix(a, "/"+preferred[0]) {
					anchor = a
					break
				}
			}
		}
	}
	if path.IsAbs(anchor) {
		return utils2.NewTagLink(anchor)
	}
	return utils2.NewLink(ri.GetRefPath(), anchor)
}

type DocumentInfo interface {
	GetRefPath() string
	GetParentDocument() DocumentInfo
}

type Unscoped interface {
	GetDocument() Document
	GetParentDocument() DocumentInfo
	GetRootContext() ResolutionContext
	GetDocumentForLink(l utils2.Link) Document

	NextId(typ string) labels.Rule
	RequestDocument(link utils2.Link, location Location) error
	RequestNumberRange(typ string)
	GetNumberRange(typ string) NumberRange
	SetNumberRangeFor(d Document, id TaggedId, typ string, nr NumberRange) HierarchyLabel
	//GetLabelInfosForType(typ string) map[labels.LabelId]TreeLabelInfo
	GetIdsForTypeInTree(typ string) map[labels.LabelId]TreeLabelInfo
	DetermineLinkPath(src, rp string) (string, error)
	DetermineLink(l utils2.Link) (string, error)
	GetLinkInfo(l utils2.Link) ResolvedRef
	GetGlobalTags(typ string) []NodeContext

	GetContextNodeContext() NodeContext
	CallStack() CallStack
	Info(key string) interface{}
	Writer() Writer
	Target() string
	RegisterUnresolved(nctx NodeContext, err error) error
}

type CallStack interface {
	History() utils2.History
	Locations() []Location

	Add(name string, loc Location) (CallStack, utils2.History)
	String() string
}

type callstack struct {
	history   utils2.History
	locations []Location
}

func NewCallStack() CallStack {
	return &callstack{}
}

func (c *callstack) History() utils2.History {
	return c.history
}

func (c *callstack) Locations() []Location {
	return c.locations
}

func (c *callstack) String() string {
	var h []string

	for _, e := range c.locations {
		h = append(h, e.String())
	}
	return strings.Join(h, " -> ")
}

func (c *callstack) Add(name string, loc Location) (CallStack, utils2.History) {
	var cycle utils2.History

	n := &callstack{}

	n.history, cycle = c.history.Add(name)
	if cycle != nil {
		return nil, cycle
	}
	n.locations = append(append(c.locations[:0:0], c.locations...), loc)
	return n, nil
}

type ResolutionContext interface {
	Parent() ResolutionContext

	Scope
	Unscoped
}

func GetNodeContext[C NodeContext](ctx ResolutionContext, n Node) C {
	return ctx.GetNodeContext(n).(C)
}

var ContextAttrs = map[string]bool{
	"scope":     true,
	"namespace": true,
	"docpath":   true,
	"docname":   true,
	"docdir":    true,
}

func GetContextAttr(name string, ctx ResolutionContext) string {
	switch name {
	case "scope":
		return ctx.GetScope().GetName()
	case "namespace":
		return ctx.GetScope().GetNamespace()
	case "docpath":
		return ctx.GetDocument().GetRefPath()
	case "docname":
		return path.Base(ctx.GetDocument().GetRefPath())
	case "docdir":
		return path.Dir(ctx.GetDocument().GetRefPath())
	}
	return ""
}

func EvaluateTag(ctx ResolutionContext, tag string) (string, bool, error) {
	explicit := ctx.GetScope().GetName() == ""
	use := true
	n := ""
	r := ""
	for _, c := range tag {
		if use {
			if c == '{' {
				if !use {
					return "", explicit, fmt.Errorf("invalid use of '{' in tag substitution")
				}
				n = ""
				use = false
				if ctx.GetScope().GetName() == "" {
					return "", explicit, fmt.Errorf("no anchor composition for document scope")
				}
			} else {
				r += string(c)
			}
		} else {
			if c == '}' {
				if n == "" {
					return "", explicit, fmt.Errorf("empty tag substitution")
				}
				if ContextAttrs[n] {
					r += GetContextAttr(n, ctx)
				} else {
					if v := ctx.LookupValue(n); v != nil {
						bctx := NewBufferContext(ctx)
						err := v.EvaluateStatic(bctx)
						if err != nil {
							return "", explicit, err
						}
						t := strings.TrimSpace(bctx.String())
						if strings.Contains(t, "\n") {
							return "", explicit, fmt.Errorf("tag substitution contains a newline")
						}
						r += strings.TrimSpace(bctx.String())
					} else {
						return "", explicit, fmt.Errorf("invalid tag substitution %q", n)
					}
				}
				explicit = true
			} else {
				n += string(c)
			}
		}
	}
	return r, explicit, nil
}

func WalkContext[C ResolutionContext](r ResolutionContext, f func(ctx C) bool) bool {
	for r != nil {
		if c, ok := r.(C); ok {
			if f(c) {
				return true
			}
		}
		r = r.Parent()
	}
	return false
}

func LookupContext(ctx ResolutionContext, c func(ctx ResolutionContext) bool) ResolutionContext {
	for ctx != nil {
		if c(ctx) {
			return ctx
		}
		ctx = ctx.Parent()
	}
	return nil
}

func LookupNodeContext[C NodeContext, N Node](ctx ResolutionContext) C {
	var zero C
	for ctx != nil {
		if nctx, ok := ctx.GetContextNodeContext().(C); ok {
			if _, ok := nctx.GetNode().(N); ok {
				return nctx
			}
		}
		ctx = ctx.Parent()
	}
	return zero
}

////////////////////////////////////////////////////////////////////////////////

type writerContext struct {
	ResolutionContext
	writer Writer
}

func (w *writerContext) Writer() Writer {
	return w.writer
}

func NewWriterContext(w Writer, ctx ResolutionContext) ResolutionContext {
	return &writerContext{ctx, w}
}

////////////////////////////////////////////////////////////////////////////////

type BufferContext struct {
	ResolutionContext
	buffer bytes.Buffer
}

func NewBufferContext(ctx ResolutionContext) *BufferContext {
	c := &BufferContext{}
	c.ResolutionContext = NewWriterContext(NewWriter(&c.buffer), ctx)
	return c
}

func (c *BufferContext) String() string {
	return c.buffer.String()
}

////////////////////////////////////////////////////////////////////////////////

type subContext struct {
	ResolutionContext
	nodeContext NodeContext
}

func (c *subContext) GetContextNodeContext() NodeContext {
	return c.nodeContext
}

func (c *subContext) Parent() ResolutionContext {
	return c.ResolutionContext
}

type SubNumberRangeContext = *subNumberRangeContext

type subNumberRangeContext struct {
	subContext
	typ string
	ids labels.NumberRange
	nr  NumberRange

	deprnr labels.NumberRange
}

func NewSubNumberRangeContext(typ string, ctx ResolutionContext, id labels.Rule, nctx NodeContext) SubNumberRangeContext {
	return &subNumberRangeContext{
		subContext: subContext{ctx, nctx},
		typ:        typ,
		ids:        labels.NewNumberRange(id.Sub()),
	}
}

func (s *subNumberRangeContext) SetNumberRange(nr NumberRange) {
	s.nr = nr
}

func (c *subNumberRangeContext) NextId(typ string) labels.Rule {
	if typ == c.typ {
		return c.ids.Next()
	}
	return c.Parent().NextId(typ)
}

func (c *subNumberRangeContext) GetNumberRange(typ string) NumberRange {
	if typ == c.typ {
		return c.nr
	}
	return c.ResolutionContext.GetNumberRange(typ)
}

////////////////////////////////////////////////////////////////////////////////

type Value struct {
	NodeSequence
	ctx ResolutionContext
}

func NewValue(ctx ResolutionContext, v NodeSequence) *Value {
	return &Value{v, ctx}
}

func (v *Value) GetContext() ResolutionContext {
	return v.ctx
}

type scoped = Scope
type unscoped = Unscoped

type StaticContext struct {
	scoped
	unscoped
}

func NewStaticContext(scope Scope, unscoped Unscoped) *StaticContext {
	return &StaticContext{
		scoped:   scope,
		unscoped: unscoped,
	}
}

func (c *StaticContext) SetScope(s Scope) {
	c.scoped = s
}

func (c *StaticContext) Parent() ResolutionContext {
	return c.unscoped.(ResolutionContext)
}

func NewLink(ctx ResolutionContext, anchors ...string) utils2.Link {
	anchor := ""
	for _, a := range anchors {
		if anchor == "" {
			anchor = a
		}
		if path.IsAbs(a) {
			return utils2.NewTagLink(anchor)
		}
	}
	return utils2.NewLink(ctx.GetDocument().GetRefPath(), anchor)
}
