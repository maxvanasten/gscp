package parser_test

import (
	"testing"

	l "github.com/maxvanasten/gscp/lexer"
	p "github.com/maxvanasten/gscp/parser"
)

// Test for Issue #3: Parser fails on multiline array arguments in switch/case blocks
func Test_Multiline_Array_In_Return_Within_Switch(t *testing.T) {
	// This reproduces the exact issue from GitHub issue #3
	input := []l.Token{
		// get_weapons_list()
		{Type: l.SYMBOL, Content: "get_weapons_list"},
		{Type: l.OPEN_PAREN, Content: "("},
		{Type: l.CLOSE_PAREN, Content: ")"},
		{Type: l.OPEN_CURLY, Content: "{"},
		{Type: l.NEWLINE, Content: ""},
		// switch(level.mapname)
		{Type: l.SYMBOL, Content: "switch"},
		{Type: l.OPEN_PAREN, Content: "("},
		{Type: l.SYMBOL, Content: "level"},
		{Type: l.OPERATOR, Content: "."},
		{Type: l.SYMBOL, Content: "mapname"},
		{Type: l.CLOSE_PAREN, Content: ")"},
		{Type: l.OPEN_CURLY, Content: "{"},
		{Type: l.NEWLINE, Content: ""},
		// case "zm_tomb":
		{Type: l.SYMBOL, Content: "case"},
		{Type: l.STRING, Content: "zm_tomb"},
		{Type: l.COLON, Content: ":"},
		{Type: l.NEWLINE, Content: ""},
		// return array(
		{Type: l.SYMBOL, Content: "return"},
		{Type: l.SYMBOL, Content: "array"},
		{Type: l.OPEN_PAREN, Content: "("},
		{Type: l.NEWLINE, Content: ""},
		// "870mcs_zm",
		{Type: l.STRING, Content: "870mcs_zm"},
		{Type: l.COMMA, Content: ","},
		{Type: l.NEWLINE, Content: ""},
		// "ak74u_extclip_zm",
		{Type: l.STRING, Content: "ak74u_extclip_zm"},
		{Type: l.COMMA, Content: ","},
		{Type: l.NEWLINE, Content: ""},
		// "ballista_zm"
		{Type: l.STRING, Content: "ballista_zm"},
		{Type: l.NEWLINE, Content: ""},
		// );
		{Type: l.CLOSE_PAREN, Content: ")"},
		{Type: l.TERMINATOR, Content: ";"},
		{Type: l.NEWLINE, Content: ""},
		// closing braces
		{Type: l.CLOSE_CURLY, Content: "}"},
		{Type: l.NEWLINE, Content: ""},
		{Type: l.CLOSE_CURLY, Content: "}"},
	}

	result, diags := p.Parse(withPositions(input))

	// Should have no diagnostics
	if len(diags) > 0 {
		t.Logf("Unexpected diagnostics: %v", diags)
		for _, d := range diags {
			t.Errorf("Unexpected diagnostic: %s at line %d, col %d", d.Message, d.Line, d.Col)
		}
	}

	// Should have at least one node
	if len(result) == 0 {
		t.Fatal("Expected function declaration, got empty result")
	}

	// Check that we got a function declaration
	if result[0].Type != "function_declaration" {
		t.Fatalf("Expected function_declaration, got %s", result[0].Type)
	}
}

// Test simpler multiline return with array
func Test_Multiline_Return_Array(t *testing.T) {
	input := []l.Token{
		{Type: l.SYMBOL, Content: "test_func"},
		{Type: l.OPEN_PAREN, Content: "("},
		{Type: l.CLOSE_PAREN, Content: ")"},
		{Type: l.OPEN_CURLY, Content: "{"},
		{Type: l.NEWLINE, Content: ""},
		{Type: l.SYMBOL, Content: "return"},
		{Type: l.SYMBOL, Content: "array"},
		{Type: l.OPEN_PAREN, Content: "("},
		{Type: l.NEWLINE, Content: ""},
		{Type: l.STRING, Content: "a"},
		{Type: l.COMMA, Content: ","},
		{Type: l.NEWLINE, Content: ""},
		{Type: l.STRING, Content: "b"},
		{Type: l.NEWLINE, Content: ""},
		{Type: l.CLOSE_PAREN, Content: ")"},
		{Type: l.TERMINATOR, Content: ";"},
		{Type: l.NEWLINE, Content: ""},
		{Type: l.CLOSE_CURLY, Content: "}"},
	}

	result, diags := p.Parse(withPositions(input))

	// Should have no error diagnostics
	hasErrors := false
	for _, d := range diags {
		if d.Severity == "error" {
			hasErrors = true
			t.Errorf("Unexpected error: %s at line %d, col %d", d.Message, d.Line, d.Col)
		}
	}

	if hasErrors {
		t.Fatal("Got error diagnostics for valid multiline return statement")
	}

	// Verify structure
	if len(result) == 0 || result[0].Type != "function_declaration" {
		t.Fatal("Expected function_declaration")
	}

	funcDecl := result[0]
	if len(funcDecl.Children) < 2 {
		t.Fatal("Expected function declaration to have args and scope children")
	}

	scope := funcDecl.Children[1]
	if scope.Type != "scope" {
		t.Fatalf("Expected scope node, got %s", scope.Type)
	}

	// Check for return_statement in scope
	foundReturn := false
	for _, child := range scope.Children {
		if child.Type == "return_statement" {
			foundReturn = true
			// The return should have the function_call as a child
			if len(child.Children) == 0 {
				t.Fatal("Return statement should have children")
			}
			if child.Children[0].Type != "function_call" {
				t.Fatalf("Expected function_call as return child, got %s", child.Children[0].Type)
			}
			if child.Children[0].Data.FunctionName != "array" {
				t.Fatalf("Expected function name 'array', got %s", child.Children[0].Data.FunctionName)
			}
			break
		}
	}

	if !foundReturn {
		t.Fatal("Expected to find return_statement in scope")
	}
}
