// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package main implements a WSJ script runner
package main

import (
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/maloquacious/semver"
	"github.com/maloquacious/wsj/ast"
	"github.com/maloquacious/wsj/parser"
)

var (
	version = semver.Version{Minor: 2, PreRelease: "alpha", Build: semver.Commit()}
)

func main() {
	input := `
	// This is a line comment
	let x = 5; /* inline comment */
	
	/* Multi-line
	   block comment */
	print(x);
	`

	result, err := parser.Parse("", []byte(input))
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
		os.Exit(1)
	}

	prog, ok := result.(*ast.Program)
	if !ok {
		fmt.Fprintf(os.Stderr, "unexpected AST type: %T\n", result)
		os.Exit(1)
	}
	spew.Dump(prog)

	fmt.Println("Parse successful!")
	for i, stmt := range prog.Statements {
		fmt.Printf("Statement %d: %#v\n", i+1, stmt)
	}
}
