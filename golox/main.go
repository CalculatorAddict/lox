package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

var hadError bool

func main() {
	args := os.Args[1:]
	if len(args) > 1 {
		fmt.Print("Usage: golox [script]")
		os.Exit(64)
	} else if len(args) == 1 {
		if err := runFile(args[0]); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(66)
		}
	} else {
		if err := runPrompt(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(67)
		}
	}
}

func runFile(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	run(string(bytes))
	return nil
}

func runPrompt() error {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		run(strings.TrimRight(line, "\r\n"))

		hadError = false
	}

	if hadError {
		os.Exit(65)
	}

	return nil
}

func run(source string) error {
	scanner := Scanner{source: source, tokens: make([]Token, 0), start: 0, current: 0, line: 1}
	tokens := scanner.ScanTokens()

	// for now, just print the tokens
	for _, token := range tokens {
		fmt.Println(token)
	}

	return nil
}

func reportError(line int, message string) {
	report(line, "", message)
}

func report(line int, where, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s\n", line, where, message)
	hadError = true
}
