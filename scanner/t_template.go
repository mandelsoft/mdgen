/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package scanner

func init() {
	Tokens.Register("template", ParseTemplate)
}

func ParseTemplate(p Parser, e Element) (Element, error) {
	if e.HasTags() {
		return nil, e.Errorf("tag not possible")
	}
	p.doc.template = true
	return p.tokenizer.NextElement()
}
