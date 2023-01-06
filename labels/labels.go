/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package labels

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////

type Label interface {
	Id() LabelId
	Type() string
	Name() string
	Level() int
}

////////////////////////////////////////////////////////////////////////////////

type LabelId struct {
	typ string
	id  string
}

func NewLabelId(typ string, id int) LabelId {
	return LabelId{
		typ: typ,
		id:  fmt.Sprintf("%d", id),
	}
}

func (l LabelId) IsPrefix(o LabelId) bool {
	return l.typ == o.typ && strings.HasPrefix(o.id, l.id+"-")
}

func (l LabelId) Type() string {
	return l.typ
}

func (l LabelId) Id() string {
	return l.id
}

func (l LabelId) Sub(id int) LabelId {
	l.id = fmt.Sprintf("%s-%d", l.id, id)
	return l
}

func (l LabelId) String() string {
	return fmt.Sprintf("%s-%s", l.typ, l.id)
}

////////////////////////////////////////////////////////////////////////////////

type StructureInfo interface {
	Format() string
	Type() string
	Level() int
	Id() LabelId
}

type StructureComponent interface {
	Comparable() string
}

type IntComponent int

func (c IntComponent) Comparable() string {
	return fmt.Sprintf("%010d", c)
}

type StringComponent string

func (c StringComponent) Comparable() string {
	return string(c)
}

type Structure struct {
	id    LabelId
	comps []StructureComponent
}

func (s *Structure) Id() LabelId {
	return s.id
}

func (s *Structure) Level() int {
	return len(s.comps)
}

func (s *Structure) Type() string {
	return s.id.typ
}

func Structured(id LabelId) *Structure {
	comps := strings.Split(id.id, "-")
	var result []StructureComponent
	for _, c := range comps {
		v, err := strconv.ParseInt(c, 10, 64)
		if err == nil {
			result = append(result, IntComponent(v))
		} else {
			result = append(result, StringComponent(c))
		}
	}
	return &Structure{
		id:    id,
		comps: result,
	}
}

type StructureList []*Structure

func (l StructureList) Sort() {
	sort.Slice(l, func(i, j int) bool {
		a := l[i]
		b := l[j]
		c := strings.Compare(a.id.typ, b.id.typ)
		if c != 0 {
			return c < 0
		}
		min := len(a.comps)
		if len(b.comps) < min {
			min = len(b.comps)
		}

		for k, ca := range a.comps[:min] {
			c := strings.Compare(ca.Comparable(), b.comps[k].Comparable())
			if c != 0 {
				return c < 0
			}
		}
		return len(b.comps) != min
	})
}

////////////////////////////////////////////////////////////////////////////////

type Rule interface {
	StructureInfo
	WithLevel(lvl int) Rule

	Name() string
	Parent() Rule

	Sub() Rule
	Next() Rule

	Reset() Rule
}

////////////////////////////////////////////////////////////////////////////////
