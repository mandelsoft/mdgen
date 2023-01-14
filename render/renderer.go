/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package render

import (
	"fmt"

	"github.com/mandelsoft/mdgen/scanner"
)

type Renderer = *renderer
type renderer struct {
	mdlink bool
}

func (r *renderer) Link(ctx scanner.ResolutionContext, link string, content func(ctx2 scanner.ResolutionContext) error) error {
	w := ctx.Writer()
	if r.mdlink {
		fmt.Fprintf(w, "[")
		err := content(ctx)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "](%s)", link)
	} else {
		fmt.Fprintf(w, "<a href=\"%s\">", link)
		err := content(ctx)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "</a>")
	}
	return nil
}

var Current Renderer = &renderer{mdlink: false}
