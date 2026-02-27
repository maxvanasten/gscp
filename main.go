package main

import (
	"os"
	"fmt"
	"github.com/maxvanasten/gscp/lexer"
	"github.com/maxvanasten/gscp/parser"
	"encoding/json"
)

func main() {
	args := os.Args[1:]

	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Usage: gscp input_file.gsc")
		os.Exit(1)
	}

	data, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
	}

	l := lexer.NewLexer(data)
	tokens := l.GetTokens()

	ast, _ := parser.Parse(tokens)
	if err := json.NewEncoder(os.Stdout).Encode(ast); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding json: %v\n", err)
		os.Exit(1)
	}
}
