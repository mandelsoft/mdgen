/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package scanner

import (
	"fmt"

	"github.com/mandelsoft/mdgen/labels"
)

type TaggedId = labels.LabelId

type TaggedNode interface {
	Node
	Id() TaggedId
	Tag() string
}

type TaggedNodeSet struct {
	tagged  map[TaggedId]TaggedNode
	anchors map[string]TaggedNode
}

func NewTaggedNodeSet() *TaggedNodeSet {
	return &TaggedNodeSet{
		tagged:  map[TaggedId]TaggedNode{},
		anchors: map[string]TaggedNode{},
	}
}

func (c *TaggedNodeSet) Register(n TaggedNode) error {
	if f := c.anchors[n.Tag()]; f != nil {
		return n.Errorf("duplicate usage of tag %q (%s)", n.Tag(), f.Location())
	}
	c.tagged[n.Id()] = n
	c.anchors[n.Tag()] = n
	return nil
}

func (c *TaggedNodeSet) GetByAnchor(anchor string) TaggedNode {
	return c.anchors[anchor]
}

func (c *TaggedNodeSet) Get(id TaggedId) TaggedNode {
	return c.tagged[id]
}

type BlockNode interface {
	TaggedNode

	Name() string
	HasParam(n string) bool

	RegisterAt(Scope) error

	Inventory() Inventory
}

type BlockNodeContext interface {
	NodeContext

	Tag() string
	Name() string
	Inventory() Inventory

	GetParameterNames() []string
	HasParam(name string) bool
	GetDefaultValue(name string) *Value
	GetNodeSequence() NodeSequence
}

type Inventory interface {
	RegisterReferencable(n TaggedNode) error
	RegisterBlock(n BlockNode) error

	GetBlockTags() map[TaggedId][]string
	GetBlockAnchors() []string
	GetBlockAnchor(anchor string) BlockNode
	GetBlock(id TaggedId) BlockNode

	//	GetTags() map[TaggedId][]string
	GetAnchors() []string
	//	GetAnchor(anchor string) *TaggedId
	//	GetRefs() map[utils2.Link]*RefInfo

	Print(gap string)
}

type inventory struct {
	unknown       int
	referencables *TaggedNodeSet
	blocks        *TaggedNodeSet
}

func NewInventory() *inventory {
	return &inventory{
		referencables: NewTaggedNodeSet(),
		blocks:        NewTaggedNodeSet(),
	}
}

func (c *inventory) RegisterReferencable(n TaggedNode) error {
	return c.referencables.Register(n)
}

func (c *inventory) RegisterBlock(n BlockNode) error {
	return c.blocks.Register(n)
}

func (c *inventory) Print(gap string) {
	fmt.Printf("%sblocks:\n", gap)
	for _, t := range c.blocks.tagged {
		t.Print(gap + "  ")
	}
	fmt.Printf("%srefererencable:\n", gap)
	for m, t := range c.referencables.tagged {
		r := "<none>"
		if t.Tag() != "" {
			r = t.Tag()
		}
		fmt.Printf("%s  %s=%s\n", gap, m, r)
	}
}

func (c *inventory) GetBlockTags() map[TaggedId][]string {
	r := map[TaggedId][]string{}
	for t, ref := range c.blocks.tagged {
		r[t] = append([]string{}, ref.Tag())
	}
	return r
}

func (c *inventory) GetTags() map[TaggedId][]string {
	r := map[TaggedId][]string{}
	for t, ref := range c.referencables.tagged {
		r[t] = append([]string{}, ref.Tag())
	}
	return r
}

func (c *inventory) GetBlockAnchors() []string {
	r := []string{}
	for ref := range c.blocks.anchors {
		r = append(r, ref)
	}
	return r
}

func (c *inventory) GetAnchors() []string {
	r := []string{}
	for ref := range c.referencables.anchors {
		r = append(r, ref)
	}
	return r
}

func (c *inventory) GetAnchor(anchor string) TaggedNode {
	return c.referencables.GetByAnchor(anchor)
}

func (c *inventory) GetBlockAnchor(anchor string) BlockNode {
	n := c.blocks.GetByAnchor(anchor)
	if n == nil {
		return nil
	}
	return n.(BlockNode)
}

func (c *inventory) GetBlock(id TaggedId) BlockNode {
	n := c.blocks.Get(id)
	if n == nil {
		return nil
	}
	return n.(BlockNode)
}
