/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package scanner

import (
	"path"
)

func init() {
	Tokens.Register("target", ParseTarget)
}

func ParseTarget(p Parser, e Element) (Element, error) {
	tag, err := e.Tag("target")
	if err != nil {
		return nil, err
	}
	tag = path.Clean(tag)
	if !path.IsAbs(tag) {
		tag = "/" + tag
	}
	p.doc.targetref = tag
	return p.tokenizer.NextElement()
}
