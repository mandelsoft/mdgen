/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/mandelsoft/mdgen/scanner"
	"github.com/mandelsoft/mdgen/tree"
	"github.com/mandelsoft/mdgen/version"
)

func Tree() {
	print := false
	args := os.Args[1:]
	if len(args) > 0 {
		if args[0] == "--version" {
			info := version.Get()

			fmt.Printf("mdgen version %s.%s.%s (%s) [%s %s]\n", info.Major, info.Minor, info.Patch, info.PreRelease, info.GitTreeState, info.GitCommit)
			os.Exit(0)
		}

		if args[0] == "--doc" {
			print = true
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

	if print {
		t.Print("")
	}

	err = t.Resolve()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("--------------\n")

	buf := bytes.NewBuffer(nil)

	err = t.GetDocument("/README").Emit(scanner.NewWriter(buf), "tmp/out")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
	//	fmt.Printf("%s\n", buf.String())

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
