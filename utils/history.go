/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package utils

import (
	"strings"
)

type History []string

func (h History) Contains(s string) bool {
	for _, e := range h {
		if e == s {
			return true
		}
	}
	return false
}

func (h History) Add(s string) (History, History) {
	for i, e := range h {
		if e == s {
			return nil, append(h[i:], s)
		}
	}
	return append(append(h[:0:0], h...), s), nil
}

func (h History) String() string {
	return strings.Join(h, "->")
}
