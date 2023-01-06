/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package utils

import (
	"sort"

	"golang.org/x/exp/constraints"
)

type Set[T comparable] map[T]struct{}

func (s Set[T]) Add(elems ...T) Set[T] {
	for _, e := range elems {
		s[e] = struct{}{}
	}
	return s
}

func (s Set[T]) Remove(elems ...T) Set[T] {
	for _, e := range elems {
		delete(s, e)
	}
	return s
}

func (s Set[T]) Has(e T) bool {
	_, ok := s[e]
	return ok
}

////////////////////////////////////////////////////////////////////////////////

func SortedMapKeys[K constraints.Ordered, E any](m map[K]E) []K {
	var result []K
	for k := range m {
		result = append(result, k)
	}

	sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })
	return result
}
