/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package utils

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("syntax", func() {
	It("simple text", func() {
		Expect(Syntax("text")).To(Equal("'`text`'"))
	})
	It("simple symbol", func() {
		Expect(Syntax("<text>")).To(Equal("&lt;*text*&gt;"))
	})

	It("simple sequence", func() {
		Expect(Syntax("text<text>other")).To(Equal("'`text`' &lt;*text*&gt; '`other`'"))
	})
	It("simple spaces", func() {
		Expect(Syntax(" ")).To(Equal("'` `'"))
		Expect(Syntax("  ")).To(Equal("{'` `'}"))
		Expect(Syntax("   ")).To(Equal("'` `' {'` `'}"))
	})
	It("intermediate spaces", func() {
		Expect(Syntax("a b")).To(Equal("'`a b`'"))
		Expect(Syntax("a  b")).To(Equal("'`a`' {'` `'} '`b`'"))
		Expect(Syntax("a   b")).To(Equal("'`a `' {'` `'} '`b`'"))
	})
	It("leading spaces", func() {
		Expect(Syntax(" b")).To(Equal("'` b`'"))
		Expect(Syntax("  b")).To(Equal("{'` `'} '`b`'"))
		Expect(Syntax("   b")).To(Equal("'` `' {'` `'} '`b`'"))
	})
	It("trailing spaces", func() {
		Expect(Syntax("a ")).To(Equal("'`a `'"))
		Expect(Syntax("a  ")).To(Equal("'`a`' {'` `'}"))
		Expect(Syntax("a   ")).To(Equal("'`a `' {'` `'}"))
	})
	It("escape", func() {
		Expect(Syntax("\\{\\{")).To(Equal("'`{{`'"))
	})
	It("some complex thing", func() {
		Expect(Syntax("\\{\\{  [*]<keyword>   <arg>{   <arg>}  \\}\\}")).To(Equal("'`{{`' {'` `'} [ '`*`' ] &lt;*keyword*&gt; '` `' {'` `'} &lt;*arg*&gt; { '` `' {'` `'} &lt;*arg*&gt; } {'` `'} '`}}`'"))
	})

})
