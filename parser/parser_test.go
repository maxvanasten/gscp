package parser_test

import (
	"testing"

	l "github.com/maxvanasten/gscp/lexer"
	p "github.com/maxvanasten/gscp/parser"
	"github.com/stretchr/testify/assert"
)

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
	assert.Equal(t, targets, result)
}

func Test_String(t *testing.T) {
	input := []l.Token{
		{Type: l.STRING, Content: "Hello, world"},
	}
	targets := []p.Node{
		{"string", p.NodeData{Content: "Hello, world"}, []p.Node{}},
	}
	result, _ := p.Parse(input)
	assert.Equal(t, targets, result)
}

func Test_Number(t *testing.T) {
	input := []l.Token{
		{Type: l.NUMBER, Content: "23"},
	}
	targets := []p.Node{
		{"number", p.NodeData{Content: "23"}, []p.Node{}},
	}
	result, _ := p.Parse(input)
	assert.Equal(t, targets, result)
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
	assert.Equal(t, targets, result)
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

func Test_Complex_Math_Expression(t *testing.T) {
	input := []l.Token{
		{Type: l.NUMBER, Content: "5"},
		{Type: l.OPERATOR, Content: "+"},
		{Type: l.OPEN_PAREN, Content: "("},
		{Type: l.NUMBER, Content: "6"},
		{Type: l.OPERATOR, Content: "*"},
		{Type: l.OPEN_PAREN, Content: "("},
		{Type: l.NUMBER, Content: "7"},
		{Type: l.OPERATOR, Content: "+"},
		{Type: l.NUMBER, Content: "8"},
		{Type: l.CLOSE_PAREN, Content: ")"},
		{Type: l.CLOSE_PAREN, Content: ")"},
		{Type: l.TERMINATOR, Content: ";"},
	}

	target := []p.Node{
		{"expression", p.NodeData{Operator: "+"}, []p.Node{
			{"lhs", p.NodeData{}, []p.Node{
				{"number", p.NodeData{Content: "5"}, []p.Node{}},
			}},
			{"rhs", p.NodeData{}, []p.Node{
				{"expression", p.NodeData{Operator: "*"}, []p.Node{
					{"lhs", p.NodeData{}, []p.Node{
						{"number", p.NodeData{Content: "6"}, []p.Node{}},
					}},
					{"rhs", p.NodeData{}, []p.Node{
						{"expression", p.NodeData{Operator: "+"}, []p.Node{
							{"lhs", p.NodeData{}, []p.Node{
								{"number", p.NodeData{Content: "7"}, []p.Node{}},
							}},
							{"rhs", p.NodeData{}, []p.Node{
								{"number", p.NodeData{Content: "8"}, []p.Node{}},
							}},
						}},
					}},
				}},
			}},
		}},
	}

	result, _ := p.Parse(input)
	assert.Equal(t, target, result)
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
		{"assignment", p.NodeData{VarName: "name"}, []p.Node{
			{"string", p.NodeData{Content: "Max"}, []p.Node{}},
		}},
		{"assignment", p.NodeData{VarName: "message"}, []p.Node{
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

func Test_Compound_Assignment(t *testing.T) {
	input := []l.Token{
		{Type: l.SYMBOL, Content: "x"},
		{Type: l.ASSIGNMENT, Content: "+="},
		{Type: l.NUMBER, Content: "1"},
		{Type: l.TERMINATOR, Content: ";"},
	}
	targets := []p.Node{
		{"assignment", p.NodeData{VarName: "x"}, []p.Node{
			{"expression", p.NodeData{Operator: "+"}, []p.Node{
				{"lhs", p.NodeData{}, []p.Node{
					{"variable_reference", p.NodeData{VarName: "x"}, []p.Node{}},
				}},
				{"rhs", p.NodeData{}, []p.Node{
					{"number", p.NodeData{Content: "1"}, []p.Node{}},
				}},
			}},
		}},
	}

	result, _ := p.Parse(input)
	assert.Equal(t, targets, result)
}

func Test_Array_Literal_Empty(t *testing.T) {
	input := []l.Token{
		{Type: l.SYMBOL, Content: "x"},
		{Type: l.ASSIGNMENT, Content: "="},
		{Type: l.OPEN_BRACKET, Content: "["},
		{Type: l.CLOSE_BRACKET, Content: "]"},
		{Type: l.TERMINATOR, Content: ";"},
	}
	targets := []p.Node{
		{"assignment", p.NodeData{VarName: "x"}, []p.Node{
			{"array_literal", p.NodeData{}, []p.Node{}},
		}},
	}

	result, _ := p.Parse(input)
	assert.Equal(t, targets, result)
}

func Test_Array_Literal_Multiple(t *testing.T) {
	input := []l.Token{
		{Type: l.SYMBOL, Content: "x"},
		{Type: l.ASSIGNMENT, Content: "="},
		{Type: l.OPEN_BRACKET, Content: "["},
		{Type: l.NUMBER, Content: "1"},
		{Type: l.COMMA, Content: ","},
		{Type: l.STRING, Content: "a"},
		{Type: l.COMMA, Content: ","},
		{Type: l.SYMBOL, Content: "y"},
		{Type: l.CLOSE_BRACKET, Content: "]"},
		{Type: l.TERMINATOR, Content: ";"},
	}
	targets := []p.Node{
		{"assignment", p.NodeData{VarName: "x"}, []p.Node{
			{"array_literal", p.NodeData{}, []p.Node{
				{"number", p.NodeData{Content: "1"}, []p.Node{}},
				{"string", p.NodeData{Content: "a"}, []p.Node{}},
				{"variable_reference", p.NodeData{VarName: "y"}, []p.Node{}},
			}},
		}},
	}

	result, _ := p.Parse(input)
	assert.Equal(t, targets, result)
}

func Test_Array_Indexing(t *testing.T) {
	input := []l.Token{
		{Type: l.SYMBOL, Content: "x"},
		{Type: l.ASSIGNMENT, Content: "="},
		{Type: l.SYMBOL, Content: "arr"},
		{Type: l.OPEN_BRACKET, Content: "["},
		{Type: l.NUMBER, Content: "0"},
		{Type: l.CLOSE_BRACKET, Content: "]"},
		{Type: l.TERMINATOR, Content: ";"},
	}
	targets := []p.Node{
		{"assignment", p.NodeData{VarName: "x"}, []p.Node{
			{"variable_reference", p.NodeData{VarName: "arr", Index: "0"}, []p.Node{}},
		}},
	}

	result, _ := p.Parse(input)
	assert.Equal(t, targets, result)
}

func Test_Array_Index_Assignment(t *testing.T) {
	input := []l.Token{
		{Type: l.SYMBOL, Content: "arr"},
		{Type: l.OPEN_BRACKET, Content: "["},
		{Type: l.NUMBER, Content: "1"},
		{Type: l.CLOSE_BRACKET, Content: "]"},
		{Type: l.ASSIGNMENT, Content: "="},
		{Type: l.STRING, Content: "x"},
		{Type: l.TERMINATOR, Content: ";"},
	}
	targets := []p.Node{
		{"assignment", p.NodeData{VarName: "arr", Index: "1"}, []p.Node{
			{"string", p.NodeData{Content: "x"}, []p.Node{}},
		}},
	}

	result, _ := p.Parse(input)
	assert.Equal(t, targets, result)
}

func Test_Array_Index_Compound_Assignment(t *testing.T) {
	input := []l.Token{
		{Type: l.SYMBOL, Content: "arr"},
		{Type: l.OPEN_BRACKET, Content: "["},
		{Type: l.SYMBOL, Content: "i"},
		{Type: l.CLOSE_BRACKET, Content: "]"},
		{Type: l.ASSIGNMENT, Content: "+="},
		{Type: l.NUMBER, Content: "1"},
		{Type: l.TERMINATOR, Content: ";"},
	}
	targets := []p.Node{
		{"assignment", p.NodeData{VarName: "arr", Index: "i"}, []p.Node{
			{"expression", p.NodeData{Operator: "+"}, []p.Node{
				{"lhs", p.NodeData{}, []p.Node{
					{"variable_reference", p.NodeData{VarName: "arr", Index: "i"}, []p.Node{}},
				}},
				{"rhs", p.NodeData{}, []p.Node{
					{"number", p.NodeData{Content: "1"}, []p.Node{}},
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

func Test_For_Loop_Infinite(t *testing.T) {
	input := []l.Token{
		{Type: l.SYMBOL, Content: "for"},
		{Type: l.OPEN_PAREN, Content: "("},
		{Type: l.TERMINATOR, Content: ";"},
		{Type: l.TERMINATOR, Content: ";"},
		{Type: l.CLOSE_PAREN, Content: ")"},
		{Type: l.OPEN_CURLY, Content: "{"},
		{Type: l.CLOSE_CURLY, Content: "}"},
	}
	target := []p.Node{
		{"for_loop", p.NodeData{}, []p.Node{
			{"for_init", p.NodeData{}, []p.Node{}},
			{"for_condition", p.NodeData{}, []p.Node{}},
			{"for_post", p.NodeData{}, []p.Node{}},
			{"scope", p.NodeData{}, []p.Node{}},
		}},
	}

	result, _ := p.Parse(input)
	assert.Equal(t, target, result)
}

func Test_For_Loop_Common(t *testing.T) {
	input := []l.Token{
		{Type: l.SYMBOL, Content: "for"},
		{Type: l.OPEN_PAREN, Content: "("},
		{Type: l.SYMBOL, Content: "i"},
		{Type: l.ASSIGNMENT, Content: "="},
		{Type: l.NUMBER, Content: "0"},
		{Type: l.TERMINATOR, Content: ";"},
		{Type: l.SYMBOL, Content: "i"},
		{Type: l.OPERATOR, Content: "<"},
		{Type: l.NUMBER, Content: "10"},
		{Type: l.TERMINATOR, Content: ";"},
		{Type: l.SYMBOL, Content: "i"},
		{Type: l.ASSIGNMENT, Content: "+="},
		{Type: l.NUMBER, Content: "1"},
		{Type: l.CLOSE_PAREN, Content: ")"},
		{Type: l.OPEN_CURLY, Content: "{"},
		{Type: l.CLOSE_CURLY, Content: "}"},
	}
	target := []p.Node{
		{"for_loop", p.NodeData{}, []p.Node{
			{"for_init", p.NodeData{}, []p.Node{
				{"assignment", p.NodeData{VarName: "i"}, []p.Node{
					{"number", p.NodeData{Content: "0"}, []p.Node{}},
				}},
			}},
			{"for_condition", p.NodeData{}, []p.Node{
				{"expression", p.NodeData{Operator: "<"}, []p.Node{
					{"lhs", p.NodeData{}, []p.Node{
						{"variable_reference", p.NodeData{VarName: "i"}, []p.Node{}},
					}},
					{"rhs", p.NodeData{}, []p.Node{
						{"number", p.NodeData{Content: "10"}, []p.Node{}},
					}},
				}},
			}},
			{"for_post", p.NodeData{}, []p.Node{
				{"assignment", p.NodeData{VarName: "i"}, []p.Node{
					{"expression", p.NodeData{Operator: "+"}, []p.Node{
						{"lhs", p.NodeData{}, []p.Node{
							{"variable_reference", p.NodeData{VarName: "i"}, []p.Node{}},
						}},
						{"rhs", p.NodeData{}, []p.Node{
							{"number", p.NodeData{Content: "1"}, []p.Node{}},
						}},
					}},
				}},
			}},
			{"scope", p.NodeData{}, []p.Node{}},
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

func Test_Function_Calls(t *testing.T) {
	input := []l.Token{
		{Type: l.SYMBOL, Content: "level"},
		{Type: l.SYMBOL, Content: "thread"},
		{Type: l.SYMBOL, Content: "somefunc"},
		{Type: l.OPEN_PAREN, Content: "("},
		{Type: l.CLOSE_PAREN, Content: ")"},
		{Type: l.TERMINATOR, Content: ";"},
		{Type: l.NEWLINE, Content: ""},

		{Type: l.SYMBOL, Content: "thread"},
		{Type: l.SYMBOL, Content: "somefunc"},
		{Type: l.OPEN_PAREN, Content: "("},
		{Type: l.CLOSE_PAREN, Content: ")"},
		{Type: l.TERMINATOR, Content: ";"},
		{Type: l.NEWLINE, Content: ""},

		{Type: l.SYMBOL, Content: "level"},
		{Type: l.SYMBOL, Content: "somefunc"},
		{Type: l.OPEN_PAREN, Content: "("},
		{Type: l.CLOSE_PAREN, Content: ")"},
		{Type: l.TERMINATOR, Content: ";"},
		{Type: l.NEWLINE, Content: ""},

		{Type: l.SYMBOL, Content: "somefunc"},
		{Type: l.OPEN_PAREN, Content: "("},
		{Type: l.CLOSE_PAREN, Content: ")"},
		{Type: l.TERMINATOR, Content: ";"},
		{Type: l.NEWLINE, Content: ""},
	}

	target := []p.Node{
		{"function_call", p.NodeData{FunctionName: "somefunc", Thread: true, Method: "level"}, []p.Node{}},
		{"function_call", p.NodeData{FunctionName: "somefunc", Thread: true}, []p.Node{}},
		{"function_call", p.NodeData{FunctionName: "somefunc", Method: "level"}, []p.Node{}},
		{"function_call", p.NodeData{FunctionName: "somefunc"}, []p.Node{}},
	}

	result, _ := p.Parse(input)
	assert.Equal(t, target, result)
}

func Test_Function_Call_Complex_Args(t *testing.T) {
	input := []l.Token{
		{Type: l.SYMBOL, Content: "somefunc"},
		{Type: l.OPEN_PAREN, Content: "("},
		{Type: l.SYMBOL, Content: "somefunc_arg"},
		{Type: l.COMMA, Content: ","},
		{Type: l.SYMBOL, Content: "child_func"},
		{Type: l.OPEN_PAREN, Content: "("},
		{Type: l.SYMBOL, Content: "child_arg1"},
		{Type: l.COMMA, Content: ","},
		{Type: l.SYMBOL, Content: "child_arg2"},
		{Type: l.CLOSE_PAREN, Content: ")"},
		{Type: l.CLOSE_PAREN, Content: ")"},
	}

	target := []p.Node{
		{"function_call", p.NodeData{FunctionName: "somefunc"}, []p.Node{
			{"variable_reference", p.NodeData{VarName: "somefunc_arg"}, []p.Node{}},
			{"function_call", p.NodeData{FunctionName: "child_func"}, []p.Node{
				{"variable_reference", p.NodeData{VarName: "child_arg1"}, []p.Node{}},
				{"variable_reference", p.NodeData{VarName: "child_arg2"}, []p.Node{}},
			}},
		}},
	}

	result, _ := p.Parse(input)
	assert.Equal(t, target, result)
}
