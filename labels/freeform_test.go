/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package labels_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mandelsoft/mdgen/labels"
	"github.com/mandelsoft/mdgen/labels/format"
)

var _ = Describe("freeform", func() {

	It("simple void freeform", func() {
		f, err := format.FormatFor("V*1")
		Expect(err).To(Succeed())
		l := labels.NewFreeForm("test", f, 0)
		Expect(l.Next().Name()).To(Equal(""))
		Expect(l.Next().Next().Name()).To(Equal(""))
		Expect(l.Next().Sub().Next().Name()).To(Equal("1"))
		Expect(l.Next().Sub().Next().Next().Sub().Next().Name()).To(Equal("2*1"))
	})

	It("simple freeform", func() {
		f, err := format.FormatFor("1")
		Expect(err).To(Succeed())
		l := labels.NewFreeForm("test", f, 0)
		Expect(l.Next().Name()).To(Equal("1"))
		Expect(l.Next().Next().Name()).To(Equal("2"))
		Expect(l.Next().Sub().Next().Name()).To(Equal("1.1"))
		Expect(l.Next().Sub().Next().Next().Sub().Next().Name()).To(Equal("1.2.1"))
	})

	It("freeform without sep", func() {
		f, err := format.FormatFor("1a")
		Expect(err).To(Succeed())
		l := labels.NewFreeForm("test", f, 0)
		Expect(l.Next().Name()).To(Equal("1"))
		Expect(l.Next().Next().Name()).To(Equal("2"))
		Expect(l.Next().Sub().Next().Name()).To(Equal("1a"))
		Expect(l.Next().Sub().Next().Next().Sub().Next().Name()).To(Equal("1ba"))
	})

	It("freeform with end sep", func() {
		f, err := format.FormatFor("1i.")
		Expect(err).To(Succeed())
		l := labels.NewFreeForm("test", f, 0)
		Expect(l.Next().Name()).To(Equal("1"))
		Expect(l.Next().Next().Name()).To(Equal("2"))
		Expect(l.Next().Sub().Next().Name()).To(Equal("1i"))
		Expect(l.Next().Sub().Next().Next().Sub().Next().Name()).To(Equal("1ii.i"))
	})

})
