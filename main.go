package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/maxvanasten/gscp/diagnostics"
	"github.com/maxvanasten/gscp/lexer"
	"github.com/maxvanasten/gscp/parser"
)

type ParseOutput struct {
	AST         []parser.Node            `json:"ast"`
	Tokens      []lexer.Token            `json:"tokens"`
	Diagnostics []diagnostics.Diagnostic `json:"diagnostics"`
}

func main() {
	args := os.Args[1:]
	var data []byte
	var err error
	if len(args) == 0 {
		// Get input from stdin
		data, err = io.ReadAll(os.Stdin)
	} else if len(args) == 1 {
		// Get input from file
		data, err = os.ReadFile(args[0])
	} else {
		// Incorrect usage
		fmt.Fprintln(os.Stderr, "Usage: gscp input_file.gsc OR echo \"some_var = 20;\" | gscp")
		os.Exit(1)
	}

	// Handle data
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting input: %v\n", err)
		os.Exit(1)
	}

	l := lexer.NewLexer(data)
	tokens := l.GetTokens()
	lexerDiagnostics := l.GetDiagnostics()

	ast, parseDiagnostics := parser.Parse(tokens)
	output := ParseOutput{
		AST:         ast,
		Tokens:      tokens,
		Diagnostics: append(lexerDiagnostics, parseDiagnostics...),
	}
	if err := json.NewEncoder(os.Stdout).Encode(output); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding json: %v\n", err)
		os.Exit(1)
	}
}
