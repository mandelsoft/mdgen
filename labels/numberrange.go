/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package labels

////////////////////////////////////////////////////////////////////////////////

type NumberRange interface {
	Sub() NumberRange
	Next() Rule
	Dup(lvl int) NumberRange
	Current() Rule
}

type numberrange struct {
	label Rule
}

func NewNumberRange(label Rule) NumberRange {
	return &numberrange{label}
}

func (n *numberrange) Sub() NumberRange {
	return NewNumberRange(n.label.Sub())
}

func (n *numberrange) Next() Rule {
	n.label = n.label.Next()
	return n.label
}

func (n *numberrange) Current() Rule {
	return n.label
}

func (n *numberrange) Dup(lvl int) NumberRange {
	if lvl < 0 {
		return NewNumberRange(n.label)
	}
	return NewNumberRange(n.label.WithLevel(lvl))
}
