/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package scanner

import (
	"bytes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("writer", func() {
	var buf *bytes.Buffer
	var w Writer

	BeforeEach(func() {
		buf = bytes.NewBuffer(nil)
		w = &writer{w: buf}
	})

	It("writes normal", func() {
		Expect(w.Write([]byte("a\nbc\ndef"))).To(Equal(8))
		Expect(w.Column()).To(Equal(3))
		Expect(w.Write([]byte("g\nhi\n"))).To(Equal(5))
		Expect(w.Column()).To(Equal(0))
		Expect(buf.String()).To(Equal(`a
bc
defg
hi
`))
	})

	It("writes normal with indentwriter", func() {
		iw := NewIndentWriter(w)
		Expect(iw.Write([]byte("a\nbc\ndef"))).To(Equal(8))
		Expect(iw.Column()).To(Equal(3))
		Expect(iw.Write([]byte("g\nhi\n"))).To(Equal(5))
		Expect(iw.Column()).To(Equal(0))
		Expect(buf.String()).To(Equal(`a
bc
defg
hi
`))
	})

	It("writes normal indented with indentwriter", func() {
		Expect(w.Write([]byte("- "))).To(Equal(2))
		iw := NewIndentWriter(w)
		Expect(iw.Write([]byte("a\nbc\ndef"))).To(Equal(8))
		Expect(iw.Column()).To(Equal(5))
		Expect(iw.Write([]byte("g\nhi\n"))).To(Equal(5))
		Expect(iw.Column()).To(Equal(0))
		Expect(iw.Write([]byte("jk"))).To(Equal(2))
		Expect(iw.Column()).To(Equal(4))
		Expect(buf.String()).To(Equal(`- a
  bc
  defg
  hi
  jk`))
	})
})
