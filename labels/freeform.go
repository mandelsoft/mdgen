/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package labels

import (
	"fmt"

	"github.com/mandelsoft/mdgen/labels/format"
)

type FreeForm = *freeform

type freeform struct {
	typ    string
	id     int
	parent *freeform
	nf     format.NumberFormat
	number int
	level  int
}

var _ Rule = (*freeform)(nil)

func NewFreeForm(typ string, nf format.NumberFormat, lvl int) Rule {
	return &freeform{typ: typ, nf: nf, level: lvl}
}

func (l *freeform) Type() string {
	return l.typ
}

func (l *freeform) Format() string {
	return "freeform"
}

func (l *freeform) WithLevel(lvl int) Rule {
	if lvl < 0 {
		return l
	}
	return &freeform{typ: l.typ, id: l.id, nf: l.nf, parent: l.parent, number: l.number, level: lvl}
}

func (l *freeform) Reset() Rule {
	return &freeform{typ: l.typ, id: l.id, nf: l.nf, parent: l.parent, level: l.level}
}

func (l *freeform) Level() int {
	return l.level
}

func (l *freeform) Name() string {
	p := ""
	if l.parent != nil {
		p = l.parent.Name()
	}
	if p == "" {
		return l.nf.Format(l.number)
	}
	return fmt.Sprintf("%s%s%s", p, l.parent.nf.Separator(), l.nf.Format(l.number))
}

func (l *freeform) Parent() Rule {
	return l.parent
}

func (l *freeform) Id() LabelId {
	if l.parent == nil {
		return NewLabelId(l.typ, l.id)
	}
	return l.parent.Id().Sub(l.id)
}

func (l *freeform) Sub() Rule {
	return &freeform{
		typ:    l.typ,
		id:     0,
		parent: l,
		level:  l.level + 1,
		nf:     l.nf.Sub(),
		number: 0,
	}
}

func (l *freeform) Next() Rule {
	return &freeform{
		typ:    l.typ,
		id:     l.id + 1,
		parent: l.parent,
		nf:     l.nf,
		level:  l.level,
		number: l.number + 1,
	}
}
