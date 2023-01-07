/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package tree

import (
	"fmt"
	"path"

	"github.com/mandelsoft/filepath/pkg/filepath"

	"github.com/mandelsoft/mdgen/labels"
	"github.com/mandelsoft/mdgen/scanner"
	utils2 "github.com/mandelsoft/mdgen/utils"
)

type info scanner.RefInfo
type resolvedRef struct {
	info
	ctx    *ResolutionContext
	anchor string
}

func NewResolvedRef(ctx *ResolutionContext, anchor string, info scanner.RefInfo) scanner.ResolvedRef {
	return &resolvedRef{
		info:   info,
		ctx:    ctx,
		anchor: anchor,
	}
}

func (r *resolvedRef) Context() scanner.ResolutionContext {
	return r.ctx
}

func (r *resolvedRef) RefPath() string {
	return r.ctx.docinfo.document.GetRefPath()
}

func (r *resolvedRef) Anchor() string {
	return r.anchor
}

var _ scanner.ResolvedRef = (*resolvedRef)(nil)

type StructInfo struct {
	document scanner.Document
	ranges   map[string]scanner.NumberRange
}

func NewStructInfo(d scanner.Document) *StructInfo {
	return &StructInfo{
		document: d,
		ranges:   map[string]scanner.NumberRange{},
	}
}

type NumberRangeInfo struct {
	scanner.NumberRange
	location *scanner.Location
	master   string
}

type RootInfo struct {
	document scanner.Document
	ranges   map[string]*NumberRangeInfo
}

func NewRootInfo(d scanner.Document) *RootInfo {
	return &RootInfo{
		document: d,
		ranges:   map[string]*NumberRangeInfo{},
	}
}

type DocumentInfo struct {
	document   scanner.Document
	structinfo *StructInfo
	rootinfo   *RootInfo
	context    *ResolutionContext

	rootnode scanner.LabeledNodeContext
}

func NewDocumentInfo(res *Resolution, d scanner.Document) *DocumentInfo {
	di := &DocumentInfo{
		document: d,
	}
	return di
}

func (i *DocumentInfo) Anchors() []string {
	return nil
}

func (i *DocumentInfo) Abbrev() string {
	if i.rootnode == nil {
		return ""
	}
	return i.rootnode.Abbrev()
}

func (i *DocumentInfo) Label() labels.Label {
	if i.rootnode == nil {
		return nil
	}
	return i.rootnode.Label()
}

func (i *DocumentInfo) Title() *string {
	t := ""
	if i.rootnode == nil {
		return &t
	}
	return i.rootnode.Title()
}

func (i *DocumentInfo) GetRefPath() string {
	return i.document.GetRefPath()
}

func (i *DocumentInfo) Source() string {
	return i.document.Source()
}

func (i *DocumentInfo) Emit(w scanner.Writer, target string) error {
	i.context.writer = w
	i.context.target = target
	return i.Walk(scanner.Resolve[scanner.Emitter](i.context))
}

func (i *DocumentInfo) IsRoot() bool {
	return i.rootinfo.document == i.document
}

func (i *DocumentInfo) Walk(f scanner.Resolver) error {
	i.context.resolution.current = i
	return i.document.Walk(f)
}

var _ scanner.RefInfo = (*DocumentInfo)(nil)

type Resolution struct {
	documents map[string]*DocumentInfo

	blocktags map[string]*DocumentInfo
	refindex  map[utils2.Link]scanner.ResolvedRef

	tagged map[string]map[string]scanner.NodeContext

	current *DocumentInfo
}

var _ scanner.LookupScope = (*Resolution)(nil)

func NewResolution(docs map[string]scanner.Document) *Resolution {
	res := &Resolution{
		documents: map[string]*DocumentInfo{},

		blocktags: map[string]*DocumentInfo{},
		refindex:  map[utils2.Link]scanner.ResolvedRef{},

		tagged: map[string]map[string]scanner.NodeContext{},
	}
	for n, d := range docs {
		di := NewDocumentInfo(res, d)
		res.documents[n] = di
		di.context = NewResolutionContext(res, di) // delay context setup to make map entry available for scope setup

		// declare standard anchorless document link
		res.refindex[utils2.NewLink(di.GetRefPath(), "")] = NewResolvedRef(di.context, "", di)
	}
	return res
}

func (r *Resolution) RegisterTag(typ string, tag string, nctx scanner.NodeContext, explicit bool) error {
	// fmt.Printf("*** registering global %s %q\n", typ, tag)
	m := r.tagged[typ]
	if m == nil {
		m = map[string]scanner.NodeContext{}
		r.tagged[typ] = m
	}
	if m[tag] != nil {
		return nctx.Errorf("%s %q already used at %s", typ, tag, m[tag].Location())
	}
	m[tag] = nctx
	return nil
}

func (r *Resolution) LookupTag(typ string, tag string) scanner.NodeContext {
	m := r.tagged[typ]
	if m == nil {
		return nil
	}
	return m[tag]
}

func (r *Resolution) GetGlobalTags(typ string) []scanner.NodeContext {
	m := r.tagged[typ]
	if m == nil {
		return nil
	}

	var result []scanner.NodeContext
	for _, v := range m {
		result = append(result, v)
	}
	return result
}

// scanner.LookupScope

func (r *Resolution) link(index map[string]*DocumentInfo, link utils2.Link) (*DocumentInfo, string) {
	var di *DocumentInfo
	var anchor string
	if link.IsTag() {
		di = index[link.Tag()]
		anchor = link.Tag()
	} else {
		if link.Path() == "" {
			return nil, ""
		}
		di = r.documents[link.Path()]
		anchor = link.Anchor()
	}
	return di, anchor
}

func (r *Resolution) GetNamespace() string {
	return ""
}

func (r *Resolution) LookupReferencable(link utils2.Link) scanner.RefInfo {
	ri := r.refindex[link]
	if ri == nil {
		return nil
	}
	return ri
}

func (r *Resolution) LookupBlock(link utils2.Link) (scanner.BlockNodeContext, scanner.Scope) {
	di, anchor := r.link(r.blocktags, link)
	if di == nil {
		return nil, nil
	}
	return di.context.GetBlock(anchor), di.context.GetScope()
}

func (r *Resolution) LookupValue(name string) *scanner.Value {
	return nil
}

func (r *Resolution) RegisterReferencable(nctx scanner.LabeledNodeContext, tags []string, explicit bool) (scanner.RefInfo, error) {
	ti := scanner.NewRefInfo(nctx, tags)
	ctx := r.documents[nctx.GetDocument().GetRefPath()].context
	if ctx.docinfo.rootnode == nil && nctx.Id().Type() == scanner.SECTION_TYPE {
		ctx.docinfo.rootnode = nctx
	}
	for _, t := range tags {
		var l utils2.Link
		if path.IsAbs(t) {
			l = utils2.NewTagLink(t)
			if ri := r.refindex[l]; ri != nil {
				return nil, fmt.Errorf("duplicate definition of tag %q: %s and %s", t, ri.Context().GetDocument().Source(), nctx.GetDocument().Source())
			}
			fmt.Printf("%s: found absolute tag %q\n", nctx.GetDocument().Source(), t)
		} else {
			l = utils2.NewLink(ctx.GetDocument().GetRefPath(), t)
			fmt.Printf("%s: found document tag %q\n", nctx.GetDocument().Source(), l)
		}
		r.refindex[l] = NewResolvedRef(ctx, t, ti)
	}
	return ti, nil
}

func (r *Resolution) RegisterBlock(anchor string, nctx scanner.BlockNodeContext) error {
	if path.IsAbs(anchor) {
		if di := r.blocktags[anchor]; di != nil {
			return fmt.Errorf("duplicate definition of block %q: %s and %s", anchor, di.document.Source(), nctx.GetDocument().Source())
		}
		r.blocktags[anchor] = r.documents[nctx.GetDocument().GetRefPath()]
		fmt.Printf("%s: found absolute block %q\n", nctx.GetDocument().Source(), anchor)
	}
	return nil
}

type ids = scanner.Ids
type scope = scanner.Scope
type unresolved struct {
	scanner.NodeContext
	err error
}

type ResolutionContext struct {
	scope
	ids

	resolution *Resolution
	callstack  scanner.CallStack

	docinfo *DocumentInfo

	docrefs      map[utils2.Link]scanner.Document
	numberranges utils2.Set[string]
	writer       scanner.Writer
	target       string

	unresolved []unresolved
}

var _ scanner.ResolutionContext = (*ResolutionContext)(nil)

func NewResolutionContext(res *Resolution, di *DocumentInfo) *ResolutionContext {
	ctx := &ResolutionContext{
		resolution:   res,
		docinfo:      di,
		ids:          scanner.Ids{},
		callstack:    scanner.NewCallStack(),
		docrefs:      map[utils2.Link]scanner.Document{},
		numberranges: utils2.Set[string]{},
	}
	ctx.scope = scanner.NewScope(res, res, ctx, di.document, "")
	return ctx
}

func (r *ResolutionContext) Info(key string) interface{} {
	return nil
}

func (r *ResolutionContext) CallStack() scanner.CallStack {
	return r.callstack
}

func (r *ResolutionContext) GetContextNodeContext() scanner.NodeContext {
	return r.docinfo.document
}

func (r *ResolutionContext) Parent() scanner.ResolutionContext {
	return nil
}

func (r *ResolutionContext) GetDocument() scanner.Document {
	return r.docinfo.document
}

func (r *ResolutionContext) GetParentDocument() scanner.Document {
	p := r.docinfo.structinfo
	if p == nil {
		return nil
	}
	return p.document
}

func (r *ResolutionContext) GetDocumentForLink(l utils2.Link) scanner.Document {
	if l.IsTag() {
		ri := r.resolution.refindex[l]
		if ri != nil {
			return ri.Context().GetDocument()
		}
	} else {
		di := r.resolution.documents[l.Path()]
		if di != nil {
			return di.document
		}
	}
	return nil
}

func (c *ResolutionContext) RequestNumberRange(typ string) {
	c.numberranges.Add(typ)
}

func (c *ResolutionContext) RequestDocument(link utils2.Link, d scanner.Document) error {
	r := c.docrefs[link]
	if r != nil {
		return fmt.Errorf("document ref %q already requested by %s", link, d.GetRefPath())
	}
	c.docrefs[link] = d
	return nil
}

func (r *ResolutionContext) LookupTag(typ string, tag string) scanner.NodeContext {
	return r.resolution.LookupTag(typ, tag)
}

func (r *ResolutionContext) GetGlobalTags(typ string) []scanner.NodeContext {
	return r.resolution.GetGlobalTags(typ)
}

func (r *ResolutionContext) GetNumberRange(typ string) scanner.NumberRange {
	di := r.docinfo
	if di.structinfo != nil && di.structinfo.ranges[typ] != nil {
		// handle bubble down label for structural documents
		return di.structinfo.ranges[typ]
	}
	return di.rootinfo.ranges[typ]
}

func (r *ResolutionContext) SetNumberRangeFor(d scanner.Document, id scanner.TaggedId, typ string, nr scanner.NumberRange) scanner.HierarchyLabel {
	lvl := -1
	di := r.resolution.documents[d.GetRefPath()]
	if l := di.document.GetLabelRules()[typ]; l != nil {
		if l.Rule == nil && l.Level >= 0 {
			lvl = l.Level
		}
	}
	next := nr.AssignableNext(lvl)
	di.structinfo.ranges[typ] = next
	return next.Current()
}

func (r *ResolutionContext) GetLabelInfosForType(typ string) map[labels.LabelId]scanner.TreeLabelInfo {
	result := map[labels.LabelId]scanner.TreeLabelInfo{}
	for _, id := range r.scope.GetIdsForType(typ) {
		info := r.scope.GetReferencable(id)
		result[info.Label().Id()] = scanner.NewTreeLabelInfo(info, r)
	}
	return result
}

func (r *ResolutionContext) GetIdsForTypeInTree(typ string) map[labels.LabelId]scanner.TreeLabelInfo {
	di := r.docinfo
	/*
		for di.structinfo != nil {
			di = r.resolution.documents[di.structinfo.document.GetRefPath()]
		}
	*/
	result := map[labels.LabelId]scanner.TreeLabelInfo{}
	di.context.appendIdsForType(typ, result)
	return result
}

func (r *ResolutionContext) appendIdsForType(typ string, result map[labels.LabelId]scanner.TreeLabelInfo) {
	for id, info := range r.GetLabelInfosForType(typ) {
		result[id] = info
	}
	for _, di := range r.resolution.documents {
		if di.structinfo != nil && di.structinfo.document == r.docinfo.document {
			di.context.appendIdsForType(typ, result)
		}
	}
}

func (r *ResolutionContext) DetermineLinkPath(src, rp string) (string, error) {
	target := filepath.Dir(r.Target())

	if !path.IsAbs(rp) && rp != "" {
		rp = filepath.Join(filepath.Dir(src), rp)
	}
	if rp != "" {
		cp, err := filepath.Canonical(rp, false)
		if err != nil {
			return "", fmt.Errorf("cannot determine relative file path from %s to %s : %w", target, rp, err)
		}
		rel, err := filepath.Rel(target, cp)
		if err != nil {
			return "", fmt.Errorf("cannot determine relative file path from %s to %s : %w", target, rp, err)
		}
		return rel, nil
	}
	return "", nil
}

func (r *ResolutionContext) DetermineLink(l utils2.Link) (string, error) {
	var err error

	resolved := r.resolution.refindex[l]

	if resolved == nil {
		return "", fmt.Errorf("cannot resolve link %s", l)
	}
	refpath := r.docinfo.document.GetRefPath()

	rel := ""
	rp := resolved.GetRefPath()
	if rp != "" && refpath != rp {
		rel, err = filepath.Rel(filepath.Dir(refpath), rp+".md")
		if err != nil {
			return "", fmt.Errorf("cannot determine relative file path for %s: %w", l, err)
		}
	}
	if resolved.Anchor() != "" {
		return rel + "#" + resolved.Anchor(), nil
	}
	return rel, nil
}

func (r *ResolutionContext) GetLinkInfo(l utils2.Link) scanner.ResolvedRef {
	return r.resolution.refindex[l]
}

func (r *ResolutionContext) Writer() scanner.Writer {
	return r.writer
}

func (r *ResolutionContext) Target() string {
	return r.target
}

func (r *ResolutionContext) RegisterUnresolved(nctx scanner.NodeContext, err error) error {
	r.unresolved = append(r.unresolved, unresolved{nctx, err})
	return err
}
