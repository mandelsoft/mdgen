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

var _ = Describe("compose", func() {
	It("composes labels", func() {
		prefix := labels.NewNumbered("test", 1).Next().Sub().Next()
		suffix := labels.NewFreeForm("test", format.MustFormatFor("A-"), 1)
		l := labels.NewComposeRule(prefix, "~", suffix)

		Expect(l.Next().Sub().Next().Next().Name()).To(Equal("1.1~A-B"))
		Expect(l.Next().Sub().Next().Next().Id().String()).To(Equal("test-1-1-1-2"))
	})
})
