/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package format

import (
	"fmt"
	"strings"
)

type NumberFormat interface {
	Sub() NumberFormat
	Format(int) string
	Separator() string
}

const Separators = "-.+*~/#_°^"

////////////////////////////////////////////////////////////////////////////////

var numberformats = map[string]func(NumberFormat, string) NumberFormat{
	"V": NewVoid,
	"1": NewNumber,
	"A": NewUpperCase,
	"a": NewLowerCase,
	"I": NewUpperRoman,
	"i": NewLowerRoman,
}

func MustFormatFor(f string) NumberFormat {
	nf, err := FormatFor(f)
	if err != nil {
		panic(err)
	}
	return nf
}

func FormatFor(f string) (NumberFormat, error) {
	if f == "" {
		return nil, fmt.Errorf("number format string missing")
	}
	var sep []string
	var typ []string
	for _, c := range f {
		if numberformats[string(c)] != nil {
			if len(typ) > len(sep) {
				sep = append(sep, "")
			}
			typ = append(typ, string(c))
		} else {
			if len(typ) == 0 || !strings.Contains(Separators, string(c)) || len(sep) >= len(typ) {
				return nil, fmt.Errorf("invalid number format %q", string(c))
			}
			sep = append(sep, string(c))
		}
	}
	if len(typ) > len(sep) {
		c := "."
		if len(sep) > 0 {
			c = sep[len(sep)-1]
		}
		sep = append(sep, c)
	}

	var n NumberFormat
	for i := range typ {
		index := len(typ) - i - 1
		n = numberformats[typ[index]](n, sep[index])
	}
	return n, nil
}

////////////////////////////////////////////////////////////////////////////////

type Format struct {
	sub    NumberFormat
	sep    string
	format func(int) string
}

func NewNumberFormat(format func(int) string, next NumberFormat, sep string) NumberFormat {
	f := &Format{sep: sep, format: format}
	if next == nil {
		next = f
	}
	f.sub = next
	f.sep = sep
	return f
}

func (f *Format) Sub() NumberFormat {
	return f.sub
}

func (f *Format) Separator() string {
	return f.sep
}

func (f *Format) Format(i int) string {
	return f.format(i)
}

////////////////////////////////////////////////////////////////////////////////

func NewVoid(next NumberFormat, sep string) NumberFormat {
	return NewNumberFormat(func(int) string { return "" }, next, sep)
}

////////////////////////////////////////////////////////////////////////////////

func NewNumber(next NumberFormat, sep string) NumberFormat {
	return NewNumberFormat(number, next, sep)
}

func number(i int) string {
	return fmt.Sprintf("%d", i)
}

////////////////////////////////////////////////////////////////////////////////

func NewLowerCase(next NumberFormat, sep string) NumberFormat {
	return NewNumberFormat(lower, next, sep)
}

func lower(i int) string {
	return fmt.Sprintf("%s", string(rune(int(rune('a'))-1+i)))
}

////////////////////////////////////////////////////////////////////////////////

func NewUpperCase(next NumberFormat, sep string) NumberFormat {
	return NewNumberFormat(upper, next, sep)
}

func upper(i int) string {
	return fmt.Sprintf("%s", string(rune(int(rune('A'))-1+i)))
}

////////////////////////////////////////////////////////////////////////////////

func NewUpperRoman(next NumberFormat, sep string) NumberFormat {
	return NewNumberFormat(Roman, next, sep)
}

func NewLowerRoman(next NumberFormat, sep string) NumberFormat {
	return NewNumberFormat(func(i int) string { return strings.ToLower(Roman(i)) }, next, sep)
}

func Roman(number int) string {
	conversions := []struct {
		value int
		digit string
	}{
		{1000, "M"},
		{900, "CM"},
		{500, "D"},
		{400, "CD"},
		{100, "C"},
		{90, "XC"},
		{50, "L"},
		{40, "XL"},
		{10, "X"},
		{9, "IX"},
		{5, "V"},
		{4, "IV"},
		{1, "I"},
	}

	roman := ""
	for _, conversion := range conversions {
		for number >= conversion.value {
			roman += conversion.digit
			number -= conversion.value
		}
	}
	return roman
}
