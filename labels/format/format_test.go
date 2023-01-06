/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package format

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("format", func() {

	It("parses simple", func() {
		f, err := FormatFor("1")
		Expect(err).To(Succeed())
		Expect(f.Format(10)).To(Equal("10"))

		Expect(f.Separator()).To(Equal("."))
		Expect(f.Sub().Format(11)).To(Equal("11"))
		Expect(f.Sub().Separator()).To(Equal("."))
	})
	It("parses simple with sep", func() {
		f, err := FormatFor("1-")
		Expect(err).To(Succeed())
		Expect(f.Format(10)).To(Equal("10"))

		Expect(f.Separator()).To(Equal("-"))
		Expect(f.Sub().Format(11)).To(Equal("11"))
		Expect(f.Sub().Separator()).To(Equal("-"))
	})
	It("parses complex", func() {
		f, err := FormatFor("1-A")
		Expect(err).To(Succeed())
		Expect(f.Format(10)).To(Equal("10"))

		Expect(f.Separator()).To(Equal("-"))
		Expect(f.Sub().Format(11)).To(Equal("K"))
		Expect(f.Sub().Separator()).To(Equal("-"))
	})

	It("parses complex with sep", func() {
		f, err := FormatFor("1-A.")
		Expect(err).To(Succeed())
		Expect(f.Format(10)).To(Equal("10"))

		Expect(f.Separator()).To(Equal("-"))
		Expect(f.Sub().Format(11)).To(Equal("K"))
		Expect(f.Sub().Separator()).To(Equal("."))
	})
	It("parses complex without sep", func() {
		f, err := FormatFor("1A.")
		Expect(err).To(Succeed())
		Expect(f.Format(10)).To(Equal("10"))

		Expect(f.Separator()).To(Equal(""))
		Expect(f.Sub().Format(11)).To(Equal("K"))
		Expect(f.Sub().Separator()).To(Equal("."))
	})
	It("roman", func() {
		f, err := FormatFor("i")
		Expect(err).To(Succeed())
		Expect(f.Format(14)).To(Equal("xiv"))
	})
})
