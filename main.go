package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/maxvanasten/gscp/generator"
	"github.com/maxvanasten/gscp/lexer"
	"github.com/maxvanasten/gscp/parser"
)

func main() {
	mode := flag.String("mode", "lex-parse", "Mode: lex, parse, lex-parse, generate, lex-parse-generate")
	input := flag.String("input", "", "Input file path")
	pretty := flag.Bool("pretty", false, "Pretty-print JSON output")
	flag.Parse()

	inputPath := *input
	if inputPath == "" && flag.NArg() > 0 {
		inputPath = flag.Arg(0)
	}

	if inputPath == "" {
		fmt.Fprintln(os.Stderr, "Usage: gscp --mode <mode> --input <file> [file]")
		os.Exit(1)
	}

	data, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	l := lexer.NewLexer(data)
	tokens := l.GetTokens()

	modeValue := strings.ToLower(*mode)
	validModes := map[string]bool{
		"lex":                true,
		"parse":              true,
		"lex-parse":          true,
		"generate":           true,
		"lex-parse-generate": true,
	}
	if !validModes[modeValue] {
		fmt.Fprintf(os.Stderr, "Invalid mode: %s\n", modeValue)
		os.Exit(1)
	}

	var jsonTokens []lexer.Token
	needsParse := modeValue == "parse" || modeValue == "lex-parse" || modeValue == "generate" || modeValue == "lex-parse-generate"
	if needsParse {
		encoded, err := json.Marshal(tokens)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error encoding tokens: %v\n", err)
			os.Exit(1)
		}
		if err := json.Unmarshal(encoded, &jsonTokens); err != nil {
			fmt.Fprintf(os.Stderr, "Error decoding tokens: %v\n", err)
			os.Exit(1)
		}
	} else {
		jsonTokens = tokens
	}

	var ast []parser.Node
	if needsParse {
		parsed, _ := parser.Parse(jsonTokens)
		ast = parsed
	}

	output := map[string]interface{}{}
	writeJSON := func(value interface{}) {
		encoder := json.NewEncoder(os.Stdout)
		if *pretty {
			encoder.SetIndent("", "  ")
		}
		if err := encoder.Encode(value); err != nil {
			fmt.Fprintf(os.Stderr, "Error encoding json: %v\n", err)
			os.Exit(1)
		}
	}

	switch modeValue {
	case "lex":
		writeJSON(tokens)
		return
	case "parse":
		writeJSON(ast)
		return
	case "lex-parse":
		output["tokens"] = jsonTokens
		output["ast"] = ast
	case "generate":
		lines := []string{}
		for _, node := range ast {
			lines = append(lines, generator.Generate(node))
		}
		fmt.Fprintln(os.Stdout, strings.Join(lines, "\n"))
		return
	case "lex-parse-generate":
		lines := []string{}
		for _, node := range ast {
			lines = append(lines, generator.Generate(node))
		}
		output["tokens"] = jsonTokens
		output["ast"] = ast
		output["generated"] = strings.Join(lines, "\n")
	}

	writeJSON(output)
}
