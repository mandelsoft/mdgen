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
		if plus {
			Expect(Syntax("   ")).To(Equal("{'` `'}+"))
		} else {
			Expect(Syntax("   ")).To(Equal("'` `' {'` `'}"))
		}
	})
	It("intermediate spaces", func() {
		Expect(Syntax("a b")).To(Equal("'`a b`'"))
		Expect(Syntax("a  b")).To(Equal("'`a`' {'` `'} '`b`'"))
		if plus {
			Expect(Syntax("a   b")).To(Equal("'`a`' {'` `'}+ '`b`'"))
		} else {
			Expect(Syntax("a   b")).To(Equal("'`a `' {'` `'} '`b`'"))
		}
	})
	It("leading spaces", func() {
		Expect(Syntax(" b")).To(Equal("'` b`'"))
		Expect(Syntax("  b")).To(Equal("{'` `'} '`b`'"))
		if plus {
			Expect(Syntax("   b")).To(Equal("{'` `'}+ '`b`'"))
		} else {
			Expect(Syntax("   b")).To(Equal("'` `' {'` `'} '`b`'"))
		}
	})
	It("trailing spaces", func() {
		Expect(Syntax("a ")).To(Equal("'`a `'"))
		Expect(Syntax("a  ")).To(Equal("'`a`' {'` `'}"))
		if plus {
			Expect(Syntax("a   ")).To(Equal("'`a`' {'` `'}+"))
		} else {
			Expect(Syntax("a   ")).To(Equal("'`a `' {'` `'}"))
		}
	})
	It("escape", func() {
		Expect(Syntax("\\{\\{")).To(Equal("'`{{`'"))
	})
	It("some complex thing", func() {
		if plus {
			Expect(Syntax("\\{\\{  [*]<keyword>   <arg>{   <arg>}  \\}\\}")).To(Equal("'`{{`' {'` `'} [ '`*`' ] &lt;*keyword*&gt; {'` `'}+ &lt;*arg*&gt; { {'` `'}+ &lt;*arg*&gt; } {'` `'} '`}}`'"))
		} else {
			Expect(Syntax("\\{\\{  [*]<keyword>   <arg>{   <arg>}  \\}\\}")).To(Equal("'`{{`' {'` `'} [ '`*`' ] &lt;*keyword*&gt; '` `' {'` `'} &lt;*arg*&gt; { '` `' {'` `'} &lt;*arg*&gt; } {'` `'} '`}}`'"))
		}
	})

	It("plus", func() {
		Expect(Syntax("{a}+b")).To(Equal("{ '`a`' }+ '`b`'"))
		Expect(Syntax("{a}\\+b")).To(Equal("{ '`a`' } '`+b`'"))
		Expect(Syntax("{a}\\++b")).To(Equal("{ '`a`' } '`++b`'"))
		Expect(Syntax("[a]+b")).To(Equal("[ '`a`' ] '`+b`'"))
	})

	It("rule", func() {
		Expect(Syntax("<rule>=a")).To(Equal("&lt;*rule*&gt; = '`a`'"))
	})

})
