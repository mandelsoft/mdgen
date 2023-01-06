/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package utils

import pluralize "github.com/gertd/go-pluralize"

var client = pluralize.NewClient()

func Plural(s string) string {
	return client.Plural(s)
}
