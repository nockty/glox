package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"

	"github.com/nockty/glox/internal/lox"
)

func main() {
	args := os.Args
	if len(args) > 2 {
		println("Usage: glox [script]")
		os.Exit(64)
	} else if len(args) == 2 {
		runFile(args[1])
	} else {
		runPrompt()
	}
}

func runFile(path string) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	run(string(bytes))
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)
	for {
		print("> ")
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			os.Exit(0)
		}
		if err != nil {
			panic(err)
		}
		run(line)
	}
}

func run(source string) {
	// TODO when running files: exit 65 for static errors, exit 70 for runtime errors
	scanner := lox.NewScanner(source)
	tokens := scanner.ScanTokens()
	parser := lox.NewParser(tokens)
	statements := parser.Parse()
	if parser.HadErrors() {
		return
	}
	lox.NewInterpreter().Interpret(statements)
}
