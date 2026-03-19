package parser_test

import (
	"testing"

	l "github.com/maxvanasten/gscp/lexer"
	p "github.com/maxvanasten/gscp/parser"
)

func Test_Multiline_Function_Call(t *testing.T) {
	// This is the bug: multiline function arguments should work
	input := []l.Token{
		{Type: l.SYMBOL, Content: "array"},
		{Type: l.OPEN_PAREN, Content: "("},
		{Type: l.NEWLINE, Content: ""},
		{Type: l.STRING, Content: "test"},
		{Type: l.NEWLINE, Content: ""},
		{Type: l.STRING, Content: "test2"},
		{Type: l.NEWLINE, Content: ""},
		{Type: l.CLOSE_PAREN, Content: ")"},
		{Type: l.TERMINATOR, Content: ";"},
	}
	
	result, diags := p.Parse(withPositions(input))
	
	if len(diags) > 0 {
		t.Logf("Diagnostics: %v", diags)
	}
	
	if len(result) == 0 {
		t.Fatal("Expected function_call node, got empty result")
	}
	
	if result[0].Type != "function_call" {
		t.Fatalf("Expected function_call node, got %s", result[0].Type)
	}
	
	if result[0].Data.FunctionName != "array" {
		t.Fatalf("Expected function name 'array', got %s", result[0].Data.FunctionName)
	}
	
	// Should have 2 children (arguments)
	if len(result[0].Children) != 2 {
		t.Fatalf("Expected 2 arguments, got %d", len(result[0].Children))
	}
}

func Test_Multiline_Array_Literal(t *testing.T) {
	// Test multiline array literal with brackets
	input := []l.Token{
		{Type: l.SYMBOL, Content: "x"},
		{Type: l.ASSIGNMENT, Content: "="},
		{Type: l.OPEN_BRACKET, Content: "["},
		{Type: l.NEWLINE, Content: ""},
		{Type: l.NUMBER, Content: "1"},
		{Type: l.NEWLINE, Content: ""},
		{Type: l.NUMBER, Content: "2"},
		{Type: l.NEWLINE, Content: ""},
		{Type: l.CLOSE_BRACKET, Content: "]"},
		{Type: l.TERMINATOR, Content: ";"},
	}
	
	result, diags := p.Parse(withPositions(input))
	
	if len(diags) > 0 {
		t.Logf("Diagnostics: %v", diags)
	}
	
	if len(result) == 0 {
		t.Fatal("Expected assignment node, got empty result")
	}
	
	if result[0].Type != "assignment" {
		t.Fatalf("Expected assignment node, got %s", result[0].Type)
	}
	
	// Check that we have an array_literal as a child
	if len(result[0].Children) == 0 || result[0].Children[0].Type != "array_literal" {
		t.Fatalf("Expected array_literal child, got %v", result[0].Children)
	}
	
	arrayLit := result[0].Children[0]
	if len(arrayLit.Children) != 2 {
		t.Fatalf("Expected 2 array elements, got %d", len(arrayLit.Children))
	}
}
