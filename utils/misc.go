/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package utils

import (
	"reflect"
	"sort"

	"github.com/modern-go/reflect2"
)

func StringMapKeys[E any](m map[string]E) []string {
	if m == nil {
		return nil
	}

	keys := []string{}
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Optional returns the first optional non-zero element given as variadic argument,
// if given, or the zero element as default.
func Optional[T any](list ...T) T {
	var zero T
	for _, e := range list {
		if !reflect.DeepEqual(e, zero) {
			return e
		}
	}
	return zero
}

// OptionalDefaulted returns the first optional non-nil element given as variadic
// argument, or the given default element. For value types a given zero
// argument is excepted, also.
func OptionalDefaulted[T any](def T, list ...T) T {
	for _, e := range list {
		if !reflect2.IsNil(e) {
			return e
		}
	}
	return def
}

// OptionalDefaultedBool checks all args for true. If no true is given
// the given default is returned.
func OptionalDefaultedBool(def bool, list ...bool) bool {
	for _, e := range list {
		if e {
			return e
		}
	}
	return def
}
