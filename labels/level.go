/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package labels

type Level = *level

type level struct {
	typ    string
	id     int
	parent *level
	number int
	level  int
}

var _ Rule = (*level)(nil)

func NewVoid(typ string, lvl int) Rule {
	if lvl < 0 {
		lvl = 0
	}
	return &level{typ: typ, level: lvl}
}

func (l *level) Type() string {
	return l.typ
}

func (l *level) Format() string {
	return "void"
}

func (l *level) WithLevel(lvl int) Rule {
	if lvl < 0 {
		return l
	}
	return &level{typ: l.typ, id: l.id, parent: l.parent, number: l.number, level: lvl}
}

func (l *level) Reset() Rule {
	return &level{typ: l.typ, id: l.id, parent: l.parent, level: l.level}
}

func (l *level) Level() int {
	return l.level
}

func (l *level) Name() string {
	return ""
}

func (l *level) Parent() Rule {
	return l.parent
}

func (l *level) Id() LabelId {
	if l.parent == nil {
		return NewLabelId(l.typ, l.id)
	}
	return l.parent.Id().Sub(l.id)
}

func (l *level) Sub() Rule {
	return &level{
		typ:    l.typ,
		id:     0,
		parent: l,
		level:  l.level + 1,
		number: 0,
	}
}

func (l *level) Next() Rule {
	return &level{
		typ:    l.typ,
		id:     l.id + 1,
		parent: l.parent,
		level:  l.level,
		number: l.number + 1,
	}
}
