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

	"github.com/mandelsoft/mdgen/testutils"
)

var _ = Describe("tokenizer", func() {
	Context("just tokens", func() {
		It("parses token", func() {
			t := NewTokenizer("data", bytes.NewBufferString("{{test}}"))
			e := testutils.Must(t.NextElement())
			Expect(e.Token()).To(Equal("test"))
		})

		It("parses token with simple arg", func() {
			t := NewTokenizer("data", bytes.NewBufferString("{{test arg}}"))
			e := testutils.Must(t.NextElement())
			Expect(e.Token()).To(Equal("test"))
			Expect(e.Tags()).To(Equal([]string{"arg"}))
		})

		It("parses token with simple arg with surrounding spaces", func() {
			t := NewTokenizer("data", bytes.NewBufferString("{{test   arg  }}"))
			e := testutils.Must(t.NextElement())
			Expect(e.Token()).To(Equal("test"))
			Expect(e.Tags()).To(Equal([]string{"arg"}))
		})

		It("parses token with multiple simple args", func() {
			t := NewTokenizer("data", bytes.NewBufferString("{{test   arg1  arg2 arg3}}"))
			e := testutils.Must(t.NextElement())
			Expect(e.Token()).To(Equal("test"))
			Expect(e.Tags()).To(Equal([]string{"arg1", "arg2", "arg3"}))
		})

		It("parses token with arg with escaped chars", func() {
			t := NewTokenizer("data", bytes.NewBufferString("{{test  arg\\\"\\\\\\  }}"))
			e := testutils.Must(t.NextElement())
			Expect(e.Token()).To(Equal("test"))
			Expect(e.Tags()).To(Equal([]string{"arg\"\\ "}))
		})

		It("parses token with quoted arg", func() {
			t := NewTokenizer("data", bytes.NewBufferString("{{test  \"a test\"  }}"))
			e := testutils.Must(t.NextElement())
			Expect(e.Token()).To(Equal("test"))
			Expect(e.Tags()).To(Equal([]string{"a test"}))
		})

		It("parses token with quoted arg with quote", func() {
			t := NewTokenizer("data", bytes.NewBufferString("{{test  \"a \\\"test\\\"\"  }}"))
			e := testutils.Must(t.NextElement())
			Expect(e.Token()).To(Equal("test"))
			Expect(e.Tags()).To(Equal([]string{"a \"test\""}))
		})

		It("parses token after comment", func() {
			t := NewTokenizer("data", bytes.NewBufferString("/#xxx\n{{test}}"))
			e := testutils.Must(t.NextElement())
			Expect(e.Token()).To(Equal("test"))
			Expect(e.Location().line).To(Equal(2))
		})
	})
})
