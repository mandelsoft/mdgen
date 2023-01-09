/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package scanner

import (
	"fmt"
	"strings"
)

type TextNode = *textnode

type textnode struct {
	NodeBase
	text string
}

func NewTextNode(d Document, location Location, txt string) TextNode {
	return &textnode{
		NodeBase: NewNodeBase(d, location),
		text:     txt,
	}
}

func (t *textnode) Print(gap string) {
	txt := t.text
	fmt.Printf("%sTEXT: %d\n%s>%s\n", gap, len(t.text), gap, strings.ReplaceAll(txt, "\n", "\n"+gap+">"))
}

func (n *textnode) Emit(ctx ResolutionContext) error {
	w := ctx.Writer()
	fmt.Fprintf(w, "%s", n.text)
	return nil
}

func (n *textnode) EvaluateStatic(ctx ResolutionContext) error {
	return n.Emit(ctx)
}
