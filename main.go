package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/maxvanasten/gscp/diagnostics"
	"github.com/maxvanasten/gscp/generator"
	"github.com/maxvanasten/gscp/lexer"
	"github.com/maxvanasten/gscp/parser"
)

type ParseOutput struct {
	AST         []parser.Node            `json:"ast"`
	Diagnostics []diagnostics.Diagnostic `json:"diagnostics"`
}

func main() {
	parsePath := flag.String("p", "", "Parse GSC file into AST")
	generatePath := flag.String("g", "", "Generate GSC from AST JSON")
	flag.Parse()
	if (*parsePath == "" && *generatePath == "") || (*parsePath != "" && *generatePath != "") {
		fmt.Fprintln(os.Stderr, "Usage: gscp -p input_file.gsc OR gscp -g input_ast.json")
		os.Exit(1)
	}

	if *parsePath != "" {
		data, err := os.ReadFile(*parsePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}

		l := lexer.NewLexer(data)
		tokens := l.GetTokens()
		lexerDiagnostics := l.GetDiagnostics()

		ast, parseDiagnostics := parser.Parse(tokens)
		output := ParseOutput{
			AST:         ast,
			Diagnostics: append(lexerDiagnostics, parseDiagnostics...),
		}
		if err := json.NewEncoder(os.Stdout).Encode(output); err != nil {
			fmt.Fprintf(os.Stderr, "Error encoding json: %v\n", err)
			os.Exit(1)
		}
		return
	}

	data, err := os.ReadFile(*generatePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}
	var ast []parser.Node
	if err := json.Unmarshal(data, &ast); err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding json: %v\n", err)
		os.Exit(1)
	}

	lines := []string{}
	for _, node := range ast {
		lines = append(lines, generator.Generate(node))
	}
	fmt.Fprintln(os.Stdout, strings.Join(lines, "\n"))
}
