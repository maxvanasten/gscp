package parser_test

import (
	l "github.com/maxvanasten/gscp/lexer"
	p "github.com/maxvanasten/gscp/parser"
	"github.com/stretchr/testify/assert"
	"testing"
)

func check_equal(targets []p.Node, actual []p.Node) bool {
	if len(targets) == 0 || len(actual) == 0 {
		return false
	}
	for i, n := range targets {
		if n.Type != actual[i].Type {
			return false
		}
		if len(targets[i].Children) > 0 {

		}
	}
	return true
}

func Test_Variable_Reference(t *testing.T) {
	// =======================
	input := []l.Token{
		{Type: l.SYMBOL, Content: "test_var"},
	}
	targets := []p.Node{
		{"variable_reference", p.NodeData{VarName: "test_var"}, []p.Node{}},
	}
	// =======================

	result, _ := p.Parse(input)
	if !check_equal(targets, result) {
		t.Fatalf("target != result")
	}
}

func Test_String(t *testing.T) {
	input := []l.Token{
		{Type: l.STRING, Content: "Hello, world"},
	}
	targets := []p.Node{
		{"string", p.NodeData{Content: "Hello, world"}, []p.Node{}},
	}
	result, _ := p.Parse(input)
	if !check_equal(targets, result) {
		t.Fatalf("targets: \n%v\nActual: \n%v\n", targets, result)
	}
}

func Test_Number(t *testing.T) {
	input := []l.Token{
		{Type: l.NUMBER, Content: "23"},
	}
	targets := []p.Node{
		{"number", p.NodeData{Content: "23"}, []p.Node{}},
	}
	result, _ := p.Parse(input)
	if !check_equal(targets, result) {
		t.Fatalf("targets: \n%v\nActual: \n%v\n", targets, result)
	}
}

func Test_Simple_Expression(t *testing.T) {
	input := []l.Token{
		{Type: l.STRING, Content: "Your age is: "},
		{Type: l.OPERATOR, Content: "+"},
		{Type: l.NUMBER, Content: "23"},
	}
	targets := []p.Node{
		{"expression", p.NodeData{Operator: "+"}, []p.Node{
			{"lhs", p.NodeData{}, []p.Node{{"string", p.NodeData{Content: "Your age is: "}, []p.Node{}}}},
			{"rhs", p.NodeData{}, []p.Node{{"number", p.NodeData{Content: "23"}, []p.Node{}}}},
		}},
	}

	result, _ := p.Parse(input)
	if !check_equal(targets, result) {
		t.Fatalf("targets: \n%v\nActual: \n%v\n", targets, result)
	}
}

func Test_Complex_Expression(t *testing.T) {
	input := []l.Token{
		{Type: l.STRING, Content: "Hello, "},
		{Type: l.OPERATOR, Content: "+"},
		{Type: l.SYMBOL, Content: "name"},
		{Type: l.OPERATOR, Content: "+"},
		{Type: l.STRING, Content: ", You are: "},
		{Type: l.OPERATOR, Content: "+"},
		{Type: l.NUMBER, Content: "23"},
		{Type: l.OPERATOR, Content: "+"},
		{Type: l.STRING, Content: " years old."},
		{Type: l.TERMINATOR, Content: ";"},
	}
	targets := []p.Node{
		{"expression", p.NodeData{Operator: "+"}, []p.Node{
			{"lhs", p.NodeData{}, []p.Node{{"string", p.NodeData{Content: "Hello, "}, []p.Node{}}}},
			{"rhs", p.NodeData{}, []p.Node{{"expression", p.NodeData{Operator: "+"}, []p.Node{
				{"lhs", p.NodeData{}, []p.Node{{"variable_reference", p.NodeData{VarName: "name"}, []p.Node{}}}},
				{"rhs", p.NodeData{}, []p.Node{{"expression", p.NodeData{Operator: "+"}, []p.Node{
					{"lhs", p.NodeData{}, []p.Node{{"string", p.NodeData{Content: ", You are: "}, []p.Node{}}}},
					{"rhs", p.NodeData{}, []p.Node{{"expression", p.NodeData{Operator: "+"}, []p.Node{
						{"lhs", p.NodeData{}, []p.Node{{"number", p.NodeData{Content: "23"}, []p.Node{}}}},
						{"rhs", p.NodeData{}, []p.Node{{"string", p.NodeData{Content: " years old."}, []p.Node{}}}},
					}}}},
				}}}},
			}}}},
		}},
	}
	results, _ := p.Parse(input)

	assert.Equal(t, targets, results)
}

func Test_Variable_Assignment(t *testing.T) {
	input := []l.Token{
		{Type: l.SYMBOL, Content: "name"},
		{Type: l.ASSIGNMENT, Content: "="},
		{Type: l.STRING, Content: "Max"},
		{Type: l.TERMINATOR, Content: ";"},
		{Type: l.NEWLINE, Content: ""},
		{Type: l.SYMBOL, Content: "message"},
		{Type: l.ASSIGNMENT, Content: "="},
		{Type: l.STRING, Content: "Hello "},
		{Type: l.OPERATOR, Content: "+"},
		{Type: l.SYMBOL, Content: "name"},
		{Type: l.TERMINATOR, Content: ";"},
	}
	targets := []p.Node{
		{"variable_assignment", p.NodeData{VarName: "name"}, []p.Node{
			{"string", p.NodeData{Content: "Max"}, []p.Node{}},
		}},
		{"variable_assignment", p.NodeData{VarName: "message"}, []p.Node{
			{"expression", p.NodeData{Operator: "+"}, []p.Node{
				{"lhs", p.NodeData{}, []p.Node{
					{"string", p.NodeData{Content: "Hello "}, []p.Node{}},
				}},
				{"rhs", p.NodeData{}, []p.Node{
					{"variable_reference", p.NodeData{VarName: "name"}, []p.Node{}},
				}},
			}},
		}},
	}

	result, _ := p.Parse(input)
	assert.Equal(t, targets, result)
}

func Test_Function_Call(t *testing.T) {
	input := []l.Token{
		{Type: l.SYMBOL, Content: "init"},
		{Type: l.OPEN_PAREN, Content: "("},
		{Type: l.SYMBOL, Content: "arg1"},
		{Type: l.COMMA, Content: ","},
		{Type: l.STRING, Content: "Hello "},
		{Type: l.OPERATOR, Content: "+"},
		{Type: l.SYMBOL, Content: "name"},
		{Type: l.CLOSE_PAREN, Content: ")"},
		{Type: l.TERMINATOR, Content: ";"},
	}
	target := []p.Node{
		{"function_call", p.NodeData{FunctionName: "init"}, []p.Node{
			{"variable_reference", p.NodeData{VarName: "arg1"}, []p.Node{}},
			{"expression", p.NodeData{Operator: "+"}, []p.Node{
				{"lhs", p.NodeData{}, []p.Node{
					{"string", p.NodeData{Content: "Hello "}, []p.Node{}},
				}},
				{"rhs", p.NodeData{}, []p.Node{
					{"variable_reference", p.NodeData{VarName: "name"}, []p.Node{}},
				}},
			}},
		}},
	}

	result, _ := p.Parse(input)
	assert.Equal(t, target, result)
}

func Test_Threaded_Function_Call(t *testing.T) {
	input := []l.Token{
		{Type: l.SYMBOL, Content: "thread"},
		{Type: l.SYMBOL, Content: "some_func"},
		{Type: l.OPEN_PAREN, Content: "("},
		{Type: l.CLOSE_PAREN, Content: ")"},
		{Type: l.TERMINATOR, Content: ";"},
	}
	target := []p.Node{
		{"function_call", p.NodeData{FunctionName: "some_func", Thread: true}, []p.Node{}},
	}

	result, _ := p.Parse(input)
	assert.Equal(t, target, result)
}

func Test_Function_Declaration(t *testing.T) {
	input := []l.Token{
		{Type: l.SYMBOL, Content: "init"},
		{Type: l.OPEN_PAREN, Content: "("},
		{Type: l.SYMBOL, Content: "name"},
		{Type: l.COMMA, Content: ","},
		{Type: l.SYMBOL, Content: "age"},
		{Type: l.CLOSE_PAREN, Content: ")"},
		{Type: l.OPEN_CURLY, Content: "{"},
		{Type: l.NEWLINE, Content: ""},
		{Type: l.SYMBOL, Content: "print"},
		{Type: l.OPEN_PAREN, Content: "("},
		{Type: l.SYMBOL, Content: "name"},
		{Type: l.CLOSE_PAREN, Content: ")"},
		{Type: l.TERMINATOR, Content: ";"},
		{Type: l.NEWLINE, Content: ""},
		{Type: l.CLOSE_CURLY, Content: "}"},
	}
	target := []p.Node{
		{"function_declaration", p.NodeData{FunctionName: "init"}, []p.Node{
			{"args", p.NodeData{}, []p.Node{
				{"variable_reference", p.NodeData{VarName: "name"}, []p.Node{}},
				{"variable_reference", p.NodeData{VarName: "age"}, []p.Node{}},
			}},
			{"scope", p.NodeData{}, []p.Node{
				{"function_call", p.NodeData{FunctionName: "print"}, []p.Node{
					{"variable_reference", p.NodeData{VarName: "name"}, []p.Node{}},
				}},
			}},
		}},
	}

	result, _ := p.Parse(input)
	assert.Equal(t, target, result)
}

func Test_IncludeStatement(t *testing.T) {
	input := []l.Token{
		{Type: l.SYMBOL, Content: "#include"},
		{Type: l.SYMBOL, Content: "path\\to\\file"},
		{Type: l.TERMINATOR, Content: ";"},
	}
	target := []p.Node{
		{"include_statement", p.NodeData{Path: "path\\to\\file"}, []p.Node{}},
	}

	result, _ := p.Parse(input)
	assert.Equal(t, target, result)
}

func Test_WaitStatement(t *testing.T) {
	input := []l.Token{
		{Type: l.SYMBOL, Content: "wait"},
		{Type: l.NUMBER, Content: "0.05"},
		{Type: l.TERMINATOR, Content: ";"},
	}
	target := []p.Node{
		{"wait_statement", p.NodeData{Delay: "0.05"}, []p.Node{}},
	}

	result, _ := p.Parse(input)
	assert.Equal(t, target, result)
}
