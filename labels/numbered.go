/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package labels

import (
	"fmt"
)

type Numbered = *numbered

type numbered struct {
	typ    string
	id     int
	parent *numbered
	number int
	level  int
}

var _ Rule = (*numbered)(nil)

func NewNumbered(typ string, lvl int) Rule {
	return &numbered{typ: typ, level: lvl}
}

func (l *numbered) Type() string {
	return l.typ
}

func (l *numbered) Format() string {
	return "numbered"
}

func (l *numbered) WithLevel(lvl int) Rule {
	if lvl < 0 {
		return l
	}
	return &numbered{typ: l.typ, id: l.id, parent: l.parent, number: l.number, level: lvl}
}

func (l *numbered) Reset() Rule {
	return &numbered{typ: l.typ, id: l.id, parent: l.parent, level: l.level}
}

func (l *numbered) Level() int {
	return l.level
}

func (l *numbered) Name() string {
	if l.parent == nil {
		return fmt.Sprintf("%d", l.number)
	}
	return fmt.Sprintf("%s.%d", l.parent.Name(), l.number)
}

func (l *numbered) Parent() Rule {
	return l.parent
}

func (l *numbered) Id() LabelId {
	if l.parent == nil {
		return NewLabelId(l.typ, l.id)
	}
	return l.parent.Id().Sub(l.id)
}

func (l *numbered) Sub() Rule {
	return &numbered{
		typ:    l.typ,
		id:     0,
		parent: l,
		level:  l.level + 1,
		number: 0,
	}
}

func (l *numbered) Next() Rule {
	return &numbered{
		typ:    l.typ,
		id:     l.id + 1,
		parent: l.parent,
		level:  l.level,
		number: l.number + 1,
	}
}
