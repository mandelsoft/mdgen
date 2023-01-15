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

type Document = *document

type LabelRuleInfo struct {
	Location  *Location
	Abbrev    string
	Rule      labels.Rule
	Separator string
	Level     int

	Master string
	Limit  int
}

func (l *LabelRuleInfo) String() string {
	if l.Rule == nil {
		if l.Level < 0 {
			return fmt.Sprintf("requested")
		}
		return fmt.Sprintf("request level %d", l.Level)
	}
	return fmt.Sprintf("%s[%s]", l.Rule.Format(), l.Rule.Id())
}

type LabelRules map[string]*LabelRuleInfo

func (d LabelRules) RequestNumberRanges(ctx ResolutionContext) {
	for t := range d {
		ctx.RequestNumberRange(t)
	}
}
func (d LabelRules) GetLabelRule(typ string) *LabelRuleInfo {
	return d[typ]
}

func (d LabelRules) SetLabelRule(loc *Location, typ string, abbrev, sep string, rule labels.Rule, lvl int) error {
	old := d[typ]

	if old != nil {
		if old.Location == nil {
			old.Location = loc
		}
		if rule != nil && old.Rule != nil {
			return fmt.Errorf("label type already set")
		}
		if lvl >= 0 {
			old.Level = lvl
		}
		if rule != nil {
			old.Rule = rule
		}
		if abbrev != "" {
			old.Abbrev = abbrev
		}
		if sep != "" {
			old.Separator = sep
		}
	} else {
		old = &LabelRuleInfo{Location: loc, Rule: rule, Level: lvl, Separator: sep, Abbrev: abbrev}
	}
	d[typ] = old
	return nil
}

func (d LabelRules) SetLabelMaster(typ string, master string, sep string, lvl int) error {
	old := d[typ]

	if old != nil {
		if old.Master != "" && old.Master != master {
			return fmt.Errorf("label master already set for %s", typ)
		}
		old.Master = master
		old.Limit = lvl
		old.Separator = sep
	} else {
		old = &LabelRuleInfo{Master: master, Limit: lvl, Separator: sep}
	}
	d[typ] = old
	return nil
}

// //////////////////////////////////////////////////////////////////////////////
type document struct {
	NodeContainerBase
	inventory documentInventory

	template  bool
	targetref string

	refpath    string
	references map[string]Node
}

type documentInventory struct {
	*inventory
	LabelRules
}

func NewDocument(source, refpath string) Document {
	d := &document{
		inventory:  documentInventory{NewInventory(), LabelRules{}},
		refpath:    refpath,
		targetref:  refpath,
		references: map[string]Node{},
	}
	d.NodeContainerBase = NewContainerBase("document", d, NewLocation(source, 0), d.inventory)
	return d
}

func (d *document) GetNode() Node {
	return d
}

func (d *document) IsTemplate() bool {
	return d.template
}

func (d *document) GetTargetRefPath() string {
	return d.targetref
}

func (d *document) Print(gap string) {
	if d.template {
		fmt.Printf("%s* template: %s\n", gap, d.Source())
	} else {
		fmt.Printf("%s* document: %s\n", gap, d.Source())
	}
	fmt.Printf("%s  refpath: %s\n", gap, d.GetRefPath())
	gap += "  "
	for k, l := range d.inventory.LabelRules {
		fmt.Printf("%slabelrule %s: %s\n", gap, k, l)
	}
	d.inventory.Print(gap)
	d.NodeContainerBase.Print(gap + "  ")
}

func (d *document) RequestNumberRanges(ctx ResolutionContext) {
	d.inventory.RequestNumberRanges(ctx)
}

func (d *document) GetInventory() Inventory {
	return d.inventory
}

func (d *document) GetRefPath() string {
	return d.refpath
}

func (d *document) GetLabelRules() LabelRules {
	return d.inventory.LabelRules
}
