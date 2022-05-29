// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package codegenerator was copied and created from singlechecker.
// The difference is that instead of os.Exit(), do panic().
package codegenerator

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/internal/analysisflags"
	"golang.org/x/tools/go/analysis/internal/checker"
	"golang.org/x/tools/go/analysis/unitchecker"
)

// Main is the main function for a checker command for a single analysis.
func Main(a *analysis.Analyzer) {
	log.SetFlags(0)
	log.SetPrefix(a.Name + ": ")

	analyzers := []*analysis.Analyzer{a}

	if err := analysis.Validate(analyzers); err != nil {
		panic(err)
	}

	checker.RegisterFlags()

	flag.Usage = func() {
		paras := strings.Split(a.Doc, "\n\n")
		fmt.Fprintf(os.Stderr, "%s: %s\n\n", a.Name, paras[0])
		fmt.Fprintf(os.Stderr, "Usage: %s [-flag] [package]\n\n", a.Name)
		if len(paras) > 1 {
			fmt.Fprintln(os.Stderr, strings.Join(paras[1:], "\n\n"))
		}
		fmt.Fprintln(os.Stderr, "\nFlags:")
		flag.PrintDefaults()
	}

	analyzers = analysisflags.Parse(analyzers, false)

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		panic("no args")
	}

	if len(args) == 1 && strings.HasSuffix(args[0], ".cfg") {
		unitchecker.Run(args[0], analyzers)
		panic("unreachable")
	}

	checker.Run(args, analyzers)
}
