/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"fmt"
	"os"

	"github.com/mandelsoft/mdgen/tree"
	"github.com/mandelsoft/mdgen/version"
)

func Tree() {
	print := false
	copy := false
	args := os.Args[1:]
	if len(args) > 0 {
		if args[0] == "--version" {
			info := version.Get()

			fmt.Printf("mdgen version %s.%s.%s (%s) [%s %s]\n", info.Major, info.Minor, info.Patch, info.PreRelease, info.GitTreeState, info.GitCommit)
			os.Exit(0)
		}

		if args[0] == "--help" {
			fmt.Printf("mdgen [--doc] [--copy] [<source dir> [<target dir>]]\n")
			fmt.Printf(`
Flags:
  --doc   print doc graph
  --copy  copy used resources into target tree

mdgen generated GitHub consistently interlinked markdown files for a tree of mdg
source files (see https://github.com/mandelsoft/mdgen).
`)
			os.Exit(0)
		}
		if args[0] == "--doc" {
			print = true
			args = args[1:]
		}
		if args[0] == "--copy" {
			copy = true
			args = args[1:]
		}
	}
	if len(args) > 2 {
		fmt.Printf("use mdgen [--doc] [<source> [<target>]]")
		os.Exit(1)
	}
	src := "."
	if len(args) > 0 {
		src = args[0]
	}
	dst := "doc"
	if len(args) > 1 {
		dst = args[1]
	}
	t, err := tree.ForFolder(src)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
	t.SetCopyMode(copy)
	if print {
		t.Print("")
	}

	err = t.Resolve()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	tw, err := tree.NewFileTreeWriter(dst)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
	tw.Close()
	err = t.Emit(tw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func main() {
	Tree()
}
