/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package labels

type composeRule struct {
	base    Rule
	sep     string
	current Rule
}

var _ Rule = (*composeRule)(nil)

func NewComposeRule(prefix Rule, sep string, label Rule) Rule {
	lvl := prefix.Level() + 1
	if label.Level() >= 0 {
		lvl = label.Level()
	}
	return &composeRule{base: prefix, sep: sep, current: label.WithLevel(lvl)}
}

func (l *composeRule) Format() string {
	return "composed"
}

func (l *composeRule) Type() string {
	return l.current.Type()
}

func (l *composeRule) Level() int {
	return l.current.Level()
}

func (l *composeRule) Id() LabelId {
	base := l.base.Id()
	cur := l.current.Id()
	base.id += "-" + cur.id
	return base
}

func (l *composeRule) WithLevel(lvl int) Rule {
	return &composeRule{base: l.base, sep: l.sep, current: l.current.WithLevel(lvl)}
}

func (l *composeRule) Reset() Rule {
	return &composeRule{base: l.base, sep: l.sep, current: l.current.Reset()}
}

func (l *composeRule) Name() string {
	base := l.base.Name()
	cur := l.current.Name()
	if base != "" && cur != "" {
		return base + l.sep + cur
	}
	return base + cur
}

func (l *composeRule) Parent() Rule {
	p := l.current.Parent()
	if p == nil {
		p = l.base
	}
	return p
}

func (l *composeRule) Sub() Rule {
	return &composeRule{base: l.base, sep: l.sep, current: l.current.Sub()}
}

func (l composeRule) Next() Rule {
	return &composeRule{base: l.base, sep: l.sep, current: l.current.Next()}
}
