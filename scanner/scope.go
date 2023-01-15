/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package scanner

import (
	"fmt"
	"path"

	utils2 "github.com/mandelsoft/mdgen/utils"
)

type LookupScope interface {
	GetNamespace() string

	LookupTag(typ string, tag string) NodeContext
	LookupReferencable(link utils2.Link) RefInfo
	LookupBlock(link utils2.Link) (BlockNodeContext, Scope)
	LookupValue(name string) *Value

	RegisterTag(typ string, tag string, nctx NodeContext, explicit bool) error
	RegisterReferencable(nctx LabeledNodeContext, tags []string, explicit bool) (RefInfo, error)
	RegisterBlock(anchor string, nctx BlockNodeContext) error
}

type Scope interface {
	LookupScope

	GetContext() ResolutionContext
	GetScope() Scope
	GetParentScope() LookupScope
	GetStaticParentScope() LookupScope

	GetName() string

	GetSubScope(name string) Scope
	NextSubScopeName(name string, extend bool) (string, error)
	AddSubScope(name string, scope Scope) error

	GetNodeContext(n Node) NodeContext
	SetNodeContext(n Node, nctx NodeContext)

	GetIdsForType(typ string) []TaggedId
	GetReferencable(id TaggedId) RefInfo

	GetBlock(anchor string) BlockNodeContext
	RegisterBlocks() error

	SetValue(name string, v *Value)
}

type nodecontexts = NodeContexts
type scope struct {
	nodecontexts
	context       ResolutionContext
	parent        LookupScope
	static        LookupScope
	name          string
	namespace     string
	inventory     Inventory
	referencables *taggedSet
	scopes        map[string]Scope
	tagged        map[string]map[string]NodeContext
	blocks        map[string]BlockNodeContext
	values        map[string]*Value
	counts        map[string]int
}

func NewScope(parent, static LookupScope, ctx ResolutionContext, inv Inventory, name string) Scope {
	ns := name
	if parent != nil {
		ns = path.Join(parent.GetNamespace(), name)
	}
	s := &scope{
		nodecontexts:  NodeContexts{},
		context:       ctx,
		parent:        parent,
		static:        static,
		name:          name,
		namespace:     ns,
		inventory:     inv,
		referencables: NewTaggedSet(),
		blocks:        map[string]BlockNodeContext{},
		tagged:        map[string]map[string]NodeContext{},
		counts:        map[string]int{},
		scopes:        map[string]Scope{},
		values:        map[string]*Value{},
	}
	s.RegisterBlocks()
	return s
}

func (s *scope) GetScope() Scope {
	return s
}

func (s *scope) GetContext() ResolutionContext {
	return s.context
}

func (s *scope) GetParentScope() LookupScope {
	return s.parent
}

func (s *scope) GetStaticParentScope() LookupScope {
	return s.static
}

func (s *scope) GetName() string {
	return s.name
}

func (s *scope) GetNamespace() string {
	return s.namespace
}

func (s *scope) GetSubScope(name string) Scope {
	return s.scopes[name]
}

func (s *scope) NextSubScopeName(name string, extend bool) (string, error) {
	if extend {
		c := s.counts[name]
		c++
		s.counts[name] = c
		return fmt.Sprintf("%s-%d", name, c), nil
	}
	if s.counts[name] != 0 {
		return "", fmt.Errorf("scope name %q already used", name)
	}
	s.counts[name] = 1
	return name, nil
}

func (s *scope) AddSubScope(name string, scope Scope) error {
	if old := s.scopes[name]; old != nil {
		return fmt.Errorf("scope name %q already used", name)
	}
	s.scopes[name] = scope
	return nil
}

func (s *scope) GetIdsForType(typ string) []TaggedId {
	return s.referencables.GetIdsForType(typ)
}

func (s *scope) GetReferencable(id TaggedId) RefInfo {
	return s.referencables.Get(id)
}

func (s *scope) RegisterTag(typ string, tag string, nctx NodeContext, explicit bool) error {
	m := s.tagged[typ]
	if m == nil {
		m = map[string]NodeContext{}
		s.tagged[typ] = m
	}
	if m[tag] != nil {
		return nctx.Errorf("%s %q already used at %s", typ, tag, m[tag].Location())
	}
	m[tag] = nctx
	if s.parent == nil {
		return nil
	}
	if !explicit {
		tag = path.Join(s.name, tag)
	}
	return s.parent.RegisterTag(typ, tag, nctx, explicit)
}

func (s *scope) RegisterReferencable(nctx LabeledNodeContext, tags []string, explicit bool) (RefInfo, error) {
	ti, err := s.referencables.Register(nctx, tags)
	if err != nil {
		return nil, err
	}
	if s.parent != nil {
		forward := tags
		if !explicit {
			prefix := "/" + s.GetName()
			if s.namespace == "" {
				prefix = nctx.GetDocument().GetRefPath()
			}
			forward = make([]string, len(tags))
			for k, v := range tags {
				if k == len(tags)-1 || explicit {
					forward[k] = v
				} else {
					if path.IsAbs(v) {
						forward[k] = path.Join(prefix + v)
					} else {
						if s.name == "" {
							forward[k] = v
						} else {
							forward[k] = path.Join(s.name, v)
						}
					}
				}
			}
		}
		i, err := s.parent.RegisterReferencable(nctx, forward, explicit)
		if err != nil {
			return nil, err
		}
		ti.anchors = i.Anchors()
	} else {
		ti.anchors = tags
	}
	return ti, nil
}

func (s *scope) LookupTag(typ string, tag string) NodeContext {
	set := s.tagged[typ]
	if set == nil {
		return nil
	}

	if nctx := set[tag]; nctx != nil {
		return nctx
	}
	if s.static != nil {
		return s.static.LookupTag(typ, tag)
	}
	return nil
}

func (s *scope) LookupReferencable(link utils2.Link) RefInfo {
	if anchor := s.anchor(link); anchor != "" {
		if ti := s.referencables.GetAnchor(anchor); ti != nil {
			return ti
		}
	}
	if s.static != nil {
		return s.static.LookupReferencable(link)
	}
	return nil
}

func (s *scope) LookupBlock(link utils2.Link) (BlockNodeContext, Scope) {
	if anchor := s.anchor(link); anchor != "" {
		if b := s.blocks[anchor]; b != nil {
			return b, s
		}
	}
	if s.static != nil {
		return s.static.LookupBlock(link)
	}
	return nil, nil
}

func (s *scope) SetValue(name string, v *Value) {
	s.values[name] = v
}

func (s *scope) LookupValue(name string) *Value {
	if v := s.values[name]; v != nil {
		return v
	}
	if s.static != nil {
		return s.static.LookupValue(name)
	}
	return nil
}

func (s *scope) GetBlock(anchor string) BlockNodeContext {
	return s.blocks[anchor]
}

func (s *scope) RegisterBlock(anchor string, nctx BlockNodeContext) error {
	if old := s.blocks[anchor]; old != nil {
		return nctx.Errorf("block %q already defined at %s", anchor, old.Location())
	}
	s.blocks[anchor] = nctx
	if s.parent != nil && path.IsAbs(anchor) {
		s.parent.RegisterBlock(anchor, nctx)
	}
	return nil
}

func (s *scope) RegisterBlocks() error {
	for t := range s.inventory.GetBlockTags() {
		b := s.inventory.GetBlock(t)
		err := b.RegisterAt(s)
		if err != nil {
			return err
		}
		nctx := s.GetNodeContext(b)
		if nctx == nil {
			panic("block context not found")
		}
		err = s.RegisterBlock(b.Tag(), nctx.(BlockNodeContext))
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *scope) anchor(link utils2.Link) string {
	anchor := ""
	if link.IsTag() {
		anchor = link.Tag()
	} else {
		if link.Path() == "" || link.Path() == s.context.GetDocument().GetRefPath() {
			anchor = link.Anchor()
		}
	}
	return anchor
}

////////////////////////////////////////////////////////////////////////////////

type nctx LabeledNodeContext
type refinfo struct {
	nctx
	anchors []string
}

func NewRefInfo(nctx LabeledNodeContext, tags []string) RefInfo {
	return &refinfo{nctx, tags}
}

func (t *refinfo) GetRefPath() string {
	return t.GetDocument().GetRefPath()
}

func (t *refinfo) GetTargetRefPath() string {
	return t.GetDocument().GetTargetRefPath()
}

func (t *refinfo) Anchors() []string {
	return t.anchors
}

var _ RefInfo = (*refinfo)(nil)

type taggedSet struct {
	bytype       map[string][]TaggedId
	tagged       map[TaggedId]*refinfo
	localanchors map[string]*refinfo
}

func NewTaggedSet() *taggedSet {
	return &taggedSet{
		bytype:       map[string][]TaggedId{},
		tagged:       map[TaggedId]*refinfo{},
		localanchors: map[string]*refinfo{},
	}
}

func refInfo(ri *refinfo) RefInfo {
	if ri == nil {
		return nil
	}
	return ri
}

func (c *taggedSet) Get(id TaggedId) RefInfo {
	return refInfo(c.tagged[id])
}

func (c *taggedSet) GetAnchorsFor(id TaggedId) []string {
	t := c.tagged[id]
	if t != nil {
		return t.anchors
	}
	return nil
}

func (c *taggedSet) GetAnchor(anchor string) RefInfo {
	return refInfo(c.localanchors[anchor])
}

func (c *taggedSet) GetIdsForType(typ string) []TaggedId {
	return c.bytype[typ]
}

func (c *taggedSet) Register(nctx LabeledNodeContext, tags []string) (*refinfo, error) {
	for _, t := range tags {
		if f := c.localanchors[t]; f != nil {
			return nil, nctx.Errorf("duplicate usage of tag %q (%s)", t, f.Location())
		}
	}
	id := nctx.Id()
	if c.tagged[id] != nil {
		return nil, nctx.Errorf("duplicate id %s", id)
	}
	ti := &refinfo{
		nctx:    nctx,
		anchors: nil,
	}
	c.bytype[id.Type()] = append(c.bytype[id.Type()], id)
	c.tagged[id] = ti
	for _, t := range tags {
		c.localanchors[t] = ti
	}
	return ti, nil
}
