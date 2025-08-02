// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/maloquacious/wsj/interpreter"
)

type replEnv struct {
	debug bool
}

func runREPL(debug bool) error {
	// todo: don't run the REPL if we're not connected to a terminal

	renv := &replEnv{debug: debug}

	rl, err := readline.NewEx(&readline.Config{
		Prompt:            "> ",
		HistoryFile:       "/tmp/wsh.repl.history", // todo: replace with ~/.wsh.history
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		HistorySearchFold: true,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to initialize readline: %v\n", err))
	}
	defer rl.Close()
	rl.CaptureExitSignal()
	log.SetOutput(rl.Stderr())

	interp, err := interpreter.New()
	if err != nil {
		return err
	}

	println("WSJ REPL - type `$exit` to quit, `$help` for help\n")

	var lines []string
	for {
		line, err := rl.Readline()
		if err == readline.ErrInterrupt {
			if len(lines) > 0 {
				lines = nil
				continue
			}
			break
		} else if err == io.EOF {
			break
		}

		if strings.TrimSpace(line) == "" {
			continue
		} else if strings.HasPrefix(strings.TrimSpace(line), "$") {
			handleReplCommand(interp, renv, strings.TrimSpace(line))
			continue
		}

		lines = append(lines, line)
		if blockComplete(lines) {
			input := strings.Join(lines, "\n")
			lines = nil

			// Change prompt back to single line
			rl.SetPrompt("> ")

			runProgram(interp, input, renv.debug)
		} else {
			rl.SetPrompt(". ")
		}
	}
	fmt.Printf("\n\n")
	return nil
}

// A simple heuristic to know when the user is done typing a block:
// 📌 Note: This is crude, but good enough for early usage. Eventually you can:
// * Track open control blocks more reliably
// * Use the parser to detect incomplete inputs (e.g., recoverable errors)
func blockComplete(lines []string) bool {
	text := strings.Join(lines, "\n")
	open := strings.Count(text, "if") + strings.Count(text, "for")
	close := strings.Count(text, "end")
	return close >= open
}

func handleReplCommand(interp *interpreter.Interpreter, env *replEnv, line string) {
	// drop any leading spaces and the '$' that signifies repl commands
	line = strings.TrimPrefix(strings.TrimSpace(line), "$")
	args := strings.Fields(line)
	if len(args) == 0 {
		return
	}
	switch args[0] {
	case "cwd":
		wd, err := os.Getwd()
		if err != nil {
			println(err)
			return
		}
		println(wd)
		return
	case "debug":
		if len(args) > 1 && args[1] == "on" {
			env.debug = true
			fmt.Println("Debug mode now enabled")
		} else if len(args) > 1 && args[1] == "off" {
			env.debug = false
			fmt.Println("Debug mode now disabled")
		} else if env.debug {
			fmt.Println("Debug mode is enabled")
		} else {
			fmt.Println("Debug mode is disabled")
		}
		return
	case "exit":
		os.Exit(0)
	case "hexes":
		fmt.Printf("$hexes is not implemented yet\n")
		//	for i, h := range vm.Root().Hexes {
		//		fmt.Printf("hexes[%d] = %s\n", i, h.Terrain)
		//	}
		return
	case "vars":
		fmt.Printf("$vars is not implemented yet\n")
		//	for k := range vm.Vars() {
		//		fmt.Println(k)
		//	}
		return
	case "version":
		println(fmt.Sprintf("repl %s", version.String()))
		return

	default:
		fmt.Printf("Unknown REPL command: %s\n", args[0])
	}
}
