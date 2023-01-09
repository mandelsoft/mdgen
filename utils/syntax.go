/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package utils

import (
	"fmt"
	"strings"
)

func Syntax(s string) (string, error) {
	result := ""
	srt := 0
	sym := ""
	txt := ""
	msk := false
	ins := false
	brk := []string{}

	for i, c := range s {
		p := string(c)
		if ins {
			if c == '>' {
				if sym == "" {
					return "", fmt.Errorf("empty symbol at %d", srt)
				}
				if result != "" {
					result += " "
				}
				result = fmt.Sprintf("%s&lt;*%s*&gt;", result, sym)
				ins = false
			} else {
				sym += p
			}
			continue
		}
		if msk {
			msk = false
			txt += p
			continue
		}
		switch c {
		case '<':
			result = textSyntax(result, &txt)
			srt = i + 1
			ins = true
			sym = ""
			txt = ""
		case '\\':
			msk = true
		case '|':
			result = textSyntax(result, &txt) + " |"
		case '[', '{', '(':
			result = textSyntax(result, &txt)
			brk = append(brk, p)
			result += " " + p
		case ']', '}', ')':
			o := "{"
			if c == ']' {
				o = "["
			} else if c == ')' {
				o = "("
			}
			if len(brk) == 0 {
				return "", fmt.Errorf("no open %s at %d", o, i+1)
			}
			if len(brk) == 0 || brk[len(brk)-1] != o {
				return "", fmt.Errorf("missing %s at %d", brk[len(brk)-1], i+1)
			}
			result = textSyntax(result, &txt)
			result += " " + p
			brk = brk[:len(brk)-1]
		default:
			if txt == "" {
				srt = i + 1
			}
			txt += p
		}
	}
	if ins {
		return "", fmt.Errorf("unfinished symbol at %d", srt)
	}
	if len(brk) > 0 {
		return "", fmt.Errorf("still open brackets: %s", strings.Join(brk, ""))
	}
	return textSyntax(result, &txt), nil
}

func textSyntax(pre string, s *string) string {
	if *s == "" {
		return pre
	}
	result := ""
	if pre != "" {
		result = " "
	}
	sc := 0
	q := false

	for _, c := range *s {
		if c == ' ' {
			sc++
		} else {
			p := string(c)
			if sc > 0 {
				if sc > 1 {
					if sc == 2 {
						if q {
							result += "`' "
							q = false
						}
						result += "{'` `'} "
					} else {
						if q {
							result += " `' {'` `'} "
							q = false
						} else {
							result += "'` `' {'` `'} "
						}
					}
				} else {
					p = " " + p
				}
			}
			if !q {
				result += "'`"
				q = true
			}
			sc = 0
			result += p
		}
	}
	switch sc {
	case 0:
		result += "`'"
	case 1:
		if q {
			result += " `'"
		} else {
			result += "'` `'"
		}
	case 2:
		if q {
			result += "`' "
		}
		result += "{'` `'}"
	case 3:
		if q {
			result += " `' {'` `'}"
		} else {
			result += "'` `' {'` `'}"
		}
	}
	*s = ""
	return pre + result
}
