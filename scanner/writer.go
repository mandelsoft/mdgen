/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package scanner

import (
	"fmt"
	"io"
)

type Writer interface {
	io.Writer
	Column() int
}

////////////////////////////////////////////////////////////////////////////////

type writer struct {
	w   io.Writer
	col int
}

func NewWriter(w io.Writer, defcol ...int) Writer {
	col := 0
	if len(defcol) > 0 {
		col = defcol[0]
	}
	return &writer{w: w, col: col}
}

func (w *writer) Write(p []byte) (n int, err error) {
	for _, c := range string(p) {
		w.col++
		if c == '\n' {
			w.col = 0
		}
	}
	return w.w.Write(p)
}

func (w *writer) Column() int {
	return w.col
}

////////////////////////////////////////////////////////////////////////////////

type indentwriter struct {
	Writer
	gap     string
	pending bool
}

func NewIndentWriter(w Writer) Writer {
	return &indentwriter{
		Writer: w,
		gap:    fmt.Sprintf("%*s", w.Column(), ""),
	}
}

func (w *indentwriter) Write(p []byte) (int, error) {
	var err error
	var n int

	written := int(0)
	last := 0
	for i, c := range string(p) {
		if c == '\n' {
			n, err = w.Writer.Write(p[last : i+1])
			written += n
			if err != nil {
				return written, err
			}
			last = i + 1
			w.pending = true
		} else {
			if w.pending {
				w.pending = false
				_, err := w.Writer.Write([]byte(w.gap))
				if err != nil {
					return written, err
				}
			}
		}
	}
	if last < len(p) {
		n, err = w.Writer.Write(p[last:])
		written += n
	}
	return written, err
}
