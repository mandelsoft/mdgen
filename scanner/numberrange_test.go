/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package scanner

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mandelsoft/mdgen/labels"
	"github.com/mandelsoft/mdgen/labels/format"
)

var _ = Describe("numberrange", func() {
	var nr NumberRange
	var nr1 NumberRange
	var nr2 NumberRange
	var l1 HierarchyLabel
	var l11 HierarchyLabel
	var l12 HierarchyLabel
	var l13 HierarchyLabel
	var l2 HierarchyLabel
	var l21 HierarchyLabel
	var l22 HierarchyLabel

	BeforeEach(func() {
		nr = NewNumberRange("test", "")

		l1 = nr.Next()
		nr1 = nr.Sub()
		l11 = nr1.Next()
		l12 = nr1.Next()
		l13 = nr1.Next()
		l2 = nr.Next()
		nr2 = nr.Sub()
		l21 = nr2.Next()
		l22 = nr2.Next()
	})

	It("creates", func() {
		Expect(l1.Id().String()).To(Equal("test-1"))
		Expect(l11.Id().String()).To(Equal("test-1-1"))
		Expect(l12.Id().String()).To(Equal("test-1-2"))
		Expect(l13.Id().String()).To(Equal("test-1-3"))
		Expect(l2.Id().String()).To(Equal("test-2"))
		Expect(l21.Id().String()).To(Equal("test-2-1"))
		Expect(l22.Id().String()).To(Equal("test-2-2"))
	})

	It("numbers", func() {
		nr.CreateLabels(labels.NewFreeForm("test", format.MustFormatFor("1."), 1))
		Expect(l1.Label().Name()).To(Equal("1"))
		Expect(l11.Label().Name()).To(Equal("1.1"))
		Expect(l12.Label().Name()).To(Equal("1.2"))
		Expect(l13.Label().Name()).To(Equal("1.3"))
		Expect(l2.Label().Name()).To(Equal("2"))
		Expect(l21.Label().Name()).To(Equal("2.1"))
		Expect(l22.Label().Name()).To(Equal("2.2"))
	})
	It("freeform", func() {
		nr.CreateLabels(labels.NewFreeForm("test", format.MustFormatFor("A-1."), 1))
		Expect(l1.Label().Name()).To(Equal("A"))
		Expect(l11.Label().Name()).To(Equal("A-1"))
		Expect(l12.Label().Name()).To(Equal("A-2"))
		Expect(l13.Label().Name()).To(Equal("A-3"))
		Expect(l2.Label().Name()).To(Equal("B"))
		Expect(l21.Label().Name()).To(Equal("B-1"))
		Expect(l22.Label().Name()).To(Equal("B-2"))
	})
	It("composes", func() {
		nr1.SetRule("", labels.NewFreeForm("test", format.MustFormatFor("a.1"), 1))
		nr.CreateLabels(labels.NewFreeForm("test", format.MustFormatFor("1."), 1))
		Expect(l1.Label().Name()).To(Equal("1"))
		Expect(l11.Label().Name()).To(Equal("1a"))
		Expect(l12.Label().Name()).To(Equal("1b"))
		Expect(l13.Label().Name()).To(Equal("1c"))
		Expect(l2.Label().Name()).To(Equal("2"))
		Expect(l21.Label().Name()).To(Equal("2.1"))
		Expect(l22.Label().Name()).To(Equal("2.2"))
	})

	Context("derived ranges", func() {
		var dep NumberRange

		BeforeEach(func() {
			dep = NewNumberRange("dep", "", func() HierarchyLabel { return nr.Current() })
			dep.SetRule("-", nil)
		})

		It("extends labels", func() {
			h1 := dep.Next()

			dep.CreateLabels(labels.NewFreeForm("dep", format.MustFormatFor("a"), 1))
			nr.CreateLabels(labels.NewFreeForm("test", format.MustFormatFor("1."), 1))

			Expect(h1.Label().Name()).To(Equal("2-a"))
		})

		It("restarts numbering for change context", func() {
			h1 := dep.Next()
			h2 := dep.Next()

			nr.Next()

			h3 := dep.Next()

			nr.CreateLabels(labels.NewFreeForm("test", format.MustFormatFor("1."), 1))
			dep.CreateLabels(labels.NewFreeForm("test", format.MustFormatFor("a"), 1))

			Expect(h1.Label().Name()).To(Equal("2-a"))
			Expect(h2.Label().Name()).To(Equal("2-b"))

			Expect(h3.Label().Name()).To(Equal("3-a"))
		})
	})

	Context("derived parent ranges", func() {
		var nr NumberRange
		var dep NumberRange
		var ctx NumberRange

		provider := func() HierarchyLabel {
			a := ctx.Current()
			for a.Level() > 1 {
				a = a.Parent()
			}
			return a
		}

		BeforeEach(func() {
			nr = NewNumberRange("test", "")
			nr.Next()
			ctx = nr
			dep = NewNumberRange("dep", "", provider)
			dep.SetRule("-", nil)
		})

		It("extends labels", func() {
			h1 := dep.Next()

			dep.CreateLabels(labels.NewFreeForm("dep", format.MustFormatFor("a"), 1))
			nr.CreateLabels(labels.NewFreeForm("test", format.MustFormatFor("1."), 1))

			Expect(h1.Label().Name()).To(Equal("1-a"))
		})

		It("restarts numbering for change context", func() {
			h1 := dep.Next()
			h2 := dep.Next()

			ctx.Next()

			h3 := dep.Next()

			nr.CreateLabels(labels.NewFreeForm("test", format.MustFormatFor("1."), 1))
			dep.CreateLabels(labels.NewFreeForm("test", format.MustFormatFor("a"), 1))

			Expect(h1.Label().Name()).To(Equal("1-a"))
			Expect(h2.Label().Name()).To(Equal("1-b"))

			Expect(h3.Label().Name()).To(Equal("2-a"))
		})

		It("limits prefix level", func() {
			h1 := dep.Next()
			h2 := dep.Next()

			ctx = ctx.Sub()
			ctx.Next()

			h3 := dep.Next()

			ctx = ctx.Sub()
			ctx.Next()

			h4 := dep.Next()
			h5 := dep.Next()

			ctx.Next()

			h6 := dep.Next()

			nr.CreateLabels(labels.NewFreeForm("test", format.MustFormatFor("1."), 1))
			dep.CreateLabels(labels.NewFreeForm("test", format.MustFormatFor("a"), 1))

			Expect(h1.Label().Name()).To(Equal("1-a"))
			Expect(h2.Label().Name()).To(Equal("1-b"))

			Expect(h3.Label().Name()).To(Equal("1.1-a"))
			Expect(h4.Label().Name()).To(Equal("1.1-b"))
			Expect(h5.Label().Name()).To(Equal("1.1-c"))
			Expect(h6.Label().Name()).To(Equal("1.1-d"))
		})
	})
})
