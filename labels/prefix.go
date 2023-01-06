/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package labels

type NameProvider interface {
	Name() string
}

type prefixLabel struct {
	prefix NameProvider
	sep    string
	label  Label
}

func NewPrefixLabel(prefix NameProvider, sep string, label Label) Label {
	return &prefixLabel{
		prefix: prefix,
		sep:    sep,
		label:  label,
	}
}

func (p *prefixLabel) Id() LabelId {
	return p.label.Id()
}

func (p *prefixLabel) Type() string {
	return p.label.Type()
}

func (p *prefixLabel) Name() string {
	t := p.prefix.Name()
	if t == "" {
		return p.label.Name()
	}
	return t + p.sep + p.label.Name()
}

func (p *prefixLabel) Level() int {
	return p.label.Level()
}
