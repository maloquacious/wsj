// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package main implements a WSJ script runner
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/maloquacious/semver"
	"github.com/maloquacious/wsj/ast"
	"github.com/maloquacious/wsj/parser"
	"github.com/maloquacious/wsj/interpreter"
)

var (
	version = semver.Version{Minor: 10, Patch: 5, PreRelease: "alpha", Build: semver.Commit()}
)

func main() {
	var (
		showVersion   = flag.Bool("version", false, "show version and exit")
		showBuildInfo = flag.Bool("build-info", false, "show build information and exit")
		debugFlag     = flag.Bool("debug", false, "enable debug mode")
	)

	flag.Parse()

	// Handle flags that exit immediately
	if *showVersion {
		fmt.Println(version.Short())
		return
	}

	if *showBuildInfo {
		fmt.Println(version.String())
		return
	}

	args := flag.Args()

	// Handle different execution modes
	switch {
	case len(args) == 0:
		// REPL mode
		if err := runREPL(*debugFlag); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

	case len(args) == 1 && strings.HasSuffix(args[0], ".wsj"):
		// Execute as a script file
		if err := runScriptFile(args[0], *debugFlag); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

	default:
		// Evaluate the arguments as program
		program := strings.Join(args, " ")
		if err := runProgram(nil, program, *debugFlag); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}
}

func runScriptFile(filename string, debug bool) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file %q: %w", filename, err)
	}

	// Handle shebang lines by converting "#!" to "//" to preserve line numbers
	if bytes.HasPrefix(data, []byte{'#', '!'}) {
		data[0], data[1] = '/', '/'
	}

	input := string(data)
	return runProgram(nil, input, debug)
}

func runProgram(interp *interpreter.Interpreter, input string, debug bool) error {
	result, err := parser.Parse("", []byte(input))
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	prog, ok := result.(*ast.Program)
	if !ok {
		return fmt.Errorf("unexpected AST type: %T", result)
	}

	if debug {
		spew.Dump(prog)
		fmt.Println("Parse successful!")
		for i, stmt := range prog.Statements {
			fmt.Printf("Statement %d: %#v\n", i+1, stmt)
		}
	}

	// TODO: Execute the program using the interpreter
	// For now, just indicate successful parsing
	fmt.Println("Program parsed successfully")

	return nil
}
