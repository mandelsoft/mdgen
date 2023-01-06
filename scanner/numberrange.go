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

type HierarchyLabel interface {
	Parent() HierarchyLabel
	Nested() NumberRange
	Next() HierarchyLabel

	Type() string
	Id() labels.LabelId
	Label() labels.Label
	Name() string
	Level() int

	SetLabel(label labels.Label)
}

type hierarchyLabel struct {
	id          labels.Rule
	level       int
	parent      HierarchyLabel
	prefixlabel HierarchyLabel
	nested      NumberRange
	next        *hierarchyLabel
	label       labels.Label
}

func (l *hierarchyLabel) Nested() NumberRange {
	return l.nested
}

func (l *hierarchyLabel) Parent() HierarchyLabel {
	return l.parent
}

func (l *hierarchyLabel) Level() int {
	return l.level
}

func (l *hierarchyLabel) Next() HierarchyLabel {
	return l.next
}

func (l *hierarchyLabel) Id() labels.LabelId {
	return l.id.Id()
}

func (l *hierarchyLabel) Label() labels.Label {
	return l.label
}

func (l *hierarchyLabel) Type() string {
	return l.id.Type()
}

func NewHierarchyLabel(id labels.Rule) HierarchyLabel {
	return &hierarchyLabel{
		id: id,
	}
}

func (l *hierarchyLabel) SetLabel(label labels.Label) {
	l.label = label
}

func (l *hierarchyLabel) Name() string {
	return l.label.Name()
}

type NumberRange interface {
	Type() string
	Level() int
	Abbrev() string

	First() HierarchyLabel

	Sub() NumberRange
	Next() HierarchyLabel
	AssignableNext(lvl int) NumberRange

	// Current return the latest entry of the number range level
	Current() HierarchyLabel

	// Actual return the latest entry available for the number range
	Actual() HierarchyLabel

	SetRule(sep string, rule labels.Rule)
	CreateLabels(rule labels.Rule)
}

type numberrange struct {
	abbrev     string
	level      int
	invalid    bool
	subcreated bool
	idrule     labels.Rule
	parent     *hierarchyLabel
	first      *hierarchyLabel
	current    *hierarchyLabel

	prefixcreator func() HierarchyLabel
	prefixlabel   HierarchyLabel
	sep           string
	rule          labels.Rule
	weight        int
}

func NewNumberRange(typ string, abbrev string, prefixcreators ...func() HierarchyLabel) NumberRange {
	var p func() HierarchyLabel
	if len(prefixcreators) > 0 {
		p = prefixcreators[0]
	}

	return &numberrange{abbrev: abbrev, idrule: labels.NewVoid(typ, 1), weight: -1, prefixcreator: p}
}

func (n *numberrange) Type() string {
	return n.idrule.Type()
}

func (n *numberrange) Level() int {
	return n.level
}

func (n *numberrange) Abbrev() string {
	return n.abbrev
}

func (n *numberrange) Sub() NumberRange {
	if n.current == nil {
		panic("number range used on initial state")
	}
	if n.subcreated {
		panic(fmt.Sprintf("sub range already created for %s", n.current.Id()))
	}
	n.subcreated = true
	return &numberrange{
		abbrev:      n.abbrev,
		level:       n.level + 1,
		prefixlabel: n.prefixlabel,
		idrule:      n.idrule.Sub(),
		parent:      n.current,
		weight:      -1,
	}
}

func (n *numberrange) SetRule(sep string, rule labels.Rule) {
	n.sep = sep
	n.rule = rule
}

func (n *numberrange) Next() HierarchyLabel {
	if n.invalid {
		panic("invalid next of assiged range")
	}
	n.subcreated = false
	if n.current != nil && n.first == nil {
		n.invalid = true
		return n.current
	}
	return n.next()
}

// I love GO
func hierarchieLabel(l *hierarchyLabel) HierarchyLabel {
	if l == nil {
		return nil
	}
	return l
}

func (n *numberrange) next() *hierarchyLabel {

	if n.prefixcreator != nil {
		n.prefixlabel = n.prefixcreator()
	}
	n.idrule = n.idrule.Next()

	l := &hierarchyLabel{
		id:          n.idrule,
		level:       n.level,
		parent:      hierarchieLabel(n.parent),
		prefixlabel: n.prefixlabel,
	}
	if n.current != nil {
		n.current.next = l
	} else {
		if n.parent != nil {
			n.parent.nested = n
		}
		n.first = l
	}
	n.current = l
	return l
}

func (n *numberrange) First() HierarchyLabel {
	return hierarchieLabel(n.first)
}

func (n *numberrange) Current() HierarchyLabel {
	return hierarchieLabel(n.current)
}

func (n *numberrange) AssignableNext(lvl int) NumberRange {
	r := &numberrange{
		level:       n.level,
		prefixlabel: n.prefixlabel,
		idrule:      n.idrule,
		parent:      n.parent,
		current:     n.next(),
		weight:      lvl,
	}
	return r
}

func (n *numberrange) Actual() HierarchyLabel {
	l := n.Current()

	for l != nil && l.Nested() != nil {
		l = l.Nested().Current()
	}
	return l
}

func (n *numberrange) CreateLabels(rule labels.Rule) {
	if n.rule != nil {
		if rule != nil {
			rule = labels.NewComposeRule(rule.Parent(), n.sep, n.rule)
		} else {
			rule = n.rule
		}
	}

	if n.weight >= 0 {
		rule = rule.WithLevel(n.weight)
	}

	var p HierarchyLabel
	l := n.first
	for l != nil {
		if l.prefixlabel != p {
			rule = rule.Reset()
		}
		p = l.prefixlabel
		rule = rule.Next()
		l.label = rule
		if n.prefixlabel != nil {
			l.label = labels.NewPrefixLabel(l.prefixlabel, n.sep, rule)
		}
		fmt.Printf("   %s -> %s[%d]\n", l.Id(), rule.Name(), rule.Level())
		if l.nested != nil {
			l.nested.CreateLabels(rule.Sub())
		}
		l = l.next
	}
}
