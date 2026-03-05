package parser_test

import (
	"testing"

	d "github.com/maxvanasten/gscp/diagnostics"
	l "github.com/maxvanasten/gscp/lexer"
	p "github.com/maxvanasten/gscp/parser"
	"github.com/stretchr/testify/assert"
)

func tokenLength(token l.Token) int {
	switch token.Type {
	case l.STRING:
		return len(token.Content) + 2
	case l.NEWLINE:
		return 1
	default:
		if token.Content != "" {
			return len(token.Content)
		}
		return 1
	}
}

func withPositions(tokens []l.Token) []l.Token {
	positioned := make([]l.Token, len(tokens))
	line := 1
	col := 1
	offset := 0
	for i, tok := range tokens {
		length := tokenLength(tok)
		tok.Line = line
		tok.Col = col
		tok.StartOffset = offset
		if tok.Type == l.NEWLINE {
			tok.EndLine = line
			tok.EndCol = col
			tok.EndOffset = offset
			line++
			col = 1
			offset++
			positioned[i] = tok
			continue
		}
		tok.EndLine = line
		tok.EndCol = col + length - 1
		tok.EndOffset = offset + length - 1
		col += length
		offset += length
		positioned[i] = tok
	}
	return positioned
}

type testNode struct {
	Type     string
	Data     p.NodeData
	Children []testNode
}

func toTestNodes(nodes []p.Node) []testNode {
	converted := make([]testNode, len(nodes))
	for i, node := range nodes {
		converted[i] = testNode{
			Type:     node.Type,
			Data:     node.Data,
			Children: toTestNodes(node.Children),
		}
	}
	return converted
}

func assertNodePositions(t *testing.T, nodes []p.Node) {
	t.Helper()
	for _, node := range nodes {
		assert.Greater(t, node.Line, 0)
		assert.Greater(t, node.Col, 0)
		assert.Greater(t, node.Length, 0)
		if len(node.Children) > 0 {
			assertNodePositions(t, node.Children)
		}
	}
}

func parseTokens(t *testing.T, tokens []l.Token) ([]testNode, []d.Diagnostic) {
	t.Helper()
	nodes, diags := p.Parse(withPositions(tokens))
	assertNodePositions(t, nodes)
	return toTestNodes(nodes), diags
}

func Test_Variable_Reference(t *testing.T) {
	// =======================
	input := []l.Token{
		{Type: l.SYMBOL, Content: "test_var"},
	}
	targets := []testNode{
		{"variable_reference", p.NodeData{VarName: "test_var"}, []testNode{}},
	}
	// =======================

	result, _ := parseTokens(t, input)
	assert.Equal(t, targets, result)
}

func Test_String(t *testing.T) {
	input := []l.Token{
		{Type: l.STRING, Content: "Hello, world"},
	}
	targets := []testNode{
		{"string", p.NodeData{Content: "Hello, world"}, []testNode{}},
	}
	result, _ := parseTokens(t, input)
	assert.Equal(t, targets, result)
}

func Test_Boolean(t *testing.T) {
	input := []l.Token{
		{Type: l.SYMBOL, Content: "true"},
	}
	targets := []testNode{
		{"boolean", p.NodeData{Content: "true"}, []testNode{}},
	}

	result, _ := parseTokens(t, input)
	assert.Equal(t, targets, result)
}

func Test_Number(t *testing.T) {
	input := []l.Token{
		{Type: l.NUMBER, Content: "23"},
	}
	targets := []testNode{
		{"number", p.NodeData{Content: "23"}, []testNode{}},
	}
	result, _ := parseTokens(t, input)
	assert.Equal(t, targets, result)
}

func Test_Simple_Expression(t *testing.T) {
	input := []l.Token{
		{Type: l.STRING, Content: "Your age is: "},
		{Type: l.OPERATOR, Content: "+"},
		{Type: l.NUMBER, Content: "23"},
	}
	targets := []testNode{
		{"expression", p.NodeData{Operator: "+"}, []testNode{
			{"lhs", p.NodeData{}, []testNode{{"string", p.NodeData{Content: "Your age is: "}, []testNode{}}}},
			{"rhs", p.NodeData{}, []testNode{{"number", p.NodeData{Content: "23"}, []testNode{}}}},
		}},
	}

	result, _ := parseTokens(t, input)
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
	targets := []testNode{
		{"expression", p.NodeData{Operator: "+"}, []testNode{
			{"lhs", p.NodeData{}, []testNode{{"string", p.NodeData{Content: "Hello, "}, []testNode{}}}},
			{"rhs", p.NodeData{}, []testNode{{"expression", p.NodeData{Operator: "+"}, []testNode{
				{"lhs", p.NodeData{}, []testNode{{"variable_reference", p.NodeData{VarName: "name"}, []testNode{}}}},
				{"rhs", p.NodeData{}, []testNode{{"expression", p.NodeData{Operator: "+"}, []testNode{
					{"lhs", p.NodeData{}, []testNode{{"string", p.NodeData{Content: ", You are: "}, []testNode{}}}},
					{"rhs", p.NodeData{}, []testNode{{"expression", p.NodeData{Operator: "+"}, []testNode{
						{"lhs", p.NodeData{}, []testNode{{"number", p.NodeData{Content: "23"}, []testNode{}}}},
						{"rhs", p.NodeData{}, []testNode{{"string", p.NodeData{Content: " years old."}, []testNode{}}}},
					}}}},
				}}}},
			}}}},
		}},
	}
	results, _ := parseTokens(t, input)

	assert.Equal(t, targets, results)
}

func Test_Logical_Expression_Precedence(t *testing.T) {
	input := []l.Token{
		{Type: l.SYMBOL, Content: "a"},
		{Type: l.OPERATOR, Content: "&&"},
		{Type: l.SYMBOL, Content: "b"},
		{Type: l.OPERATOR, Content: "||"},
		{Type: l.SYMBOL, Content: "c"},
		{Type: l.TERMINATOR, Content: ";"},
	}
	targets := []testNode{
		{"expression", p.NodeData{Operator: "||"}, []testNode{
			{"lhs", p.NodeData{}, []testNode{
				{"expression", p.NodeData{Operator: "&&"}, []testNode{
					{"lhs", p.NodeData{}, []testNode{{"variable_reference", p.NodeData{VarName: "a"}, []testNode{}}}},
					{"rhs", p.NodeData{}, []testNode{{"variable_reference", p.NodeData{VarName: "b"}, []testNode{}}}},
				}},
			}},
			{"rhs", p.NodeData{}, []testNode{
				{"variable_reference", p.NodeData{VarName: "c"}, []testNode{}},
			}},
		}},
	}

	result, _ := parseTokens(t, input)
	assert.Equal(t, targets, result)
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

	target := []testNode{
		{"expression", p.NodeData{Operator: "+"}, []testNode{
			{"lhs", p.NodeData{}, []testNode{
				{"number", p.NodeData{Content: "5"}, []testNode{}},
			}},
			{"rhs", p.NodeData{}, []testNode{
				{"expression", p.NodeData{Operator: "*"}, []testNode{
					{"lhs", p.NodeData{}, []testNode{
						{"number", p.NodeData{Content: "6"}, []testNode{}},
					}},
					{"rhs", p.NodeData{}, []testNode{
						{"expression", p.NodeData{Operator: "+"}, []testNode{
							{"lhs", p.NodeData{}, []testNode{
								{"number", p.NodeData{Content: "7"}, []testNode{}},
							}},
							{"rhs", p.NodeData{}, []testNode{
								{"number", p.NodeData{Content: "8"}, []testNode{}},
							}},
						}},
					}},
				}},
			}},
		}},
	}

	result, _ := parseTokens(t, input)
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
	targets := []testNode{
		{"assignment", p.NodeData{VarName: "name"}, []testNode{
			{"string", p.NodeData{Content: "Max"}, []testNode{}},
		}},
		{"assignment", p.NodeData{VarName: "message"}, []testNode{
			{"expression", p.NodeData{Operator: "+"}, []testNode{
				{"lhs", p.NodeData{}, []testNode{
					{"string", p.NodeData{Content: "Hello "}, []testNode{}},
				}},
				{"rhs", p.NodeData{}, []testNode{
					{"variable_reference", p.NodeData{VarName: "name"}, []testNode{}},
				}},
			}},
		}},
	}

	result, _ := parseTokens(t, input)
	assert.Equal(t, targets, result)
}

func Test_Compound_Assignment(t *testing.T) {
	input := []l.Token{
		{Type: l.SYMBOL, Content: "x"},
		{Type: l.ASSIGNMENT, Content: "+="},
		{Type: l.NUMBER, Content: "1"},
		{Type: l.TERMINATOR, Content: ";"},
	}
	targets := []testNode{
		{"assignment", p.NodeData{VarName: "x"}, []testNode{
			{"expression", p.NodeData{Operator: "+"}, []testNode{
				{"lhs", p.NodeData{}, []testNode{
					{"variable_reference", p.NodeData{VarName: "x"}, []testNode{}},
				}},
				{"rhs", p.NodeData{}, []testNode{
					{"number", p.NodeData{Content: "1"}, []testNode{}},
				}},
			}},
		}},
	}

	result, _ := parseTokens(t, input)
	assert.Equal(t, targets, result)
}

func Test_Unary_Expression(t *testing.T) {
	input := []l.Token{
		{Type: l.OPERATOR, Content: "!"},
		{Type: l.SYMBOL, Content: "true"},
	}
	targets := []testNode{
		{"unary_expression", p.NodeData{Operator: "!"}, []testNode{
			{"boolean", p.NodeData{Content: "true"}, []testNode{}},
		}},
	}

	result, _ := parseTokens(t, input)
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
	targets := []testNode{
		{"assignment", p.NodeData{VarName: "x"}, []testNode{
			{"array_literal", p.NodeData{}, []testNode{}},
		}},
	}

	result, _ := parseTokens(t, input)
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
	targets := []testNode{
		{"assignment", p.NodeData{VarName: "x"}, []testNode{
			{"array_literal", p.NodeData{}, []testNode{
				{"number", p.NodeData{Content: "1"}, []testNode{}},
				{"string", p.NodeData{Content: "a"}, []testNode{}},
				{"variable_reference", p.NodeData{VarName: "y"}, []testNode{}},
			}},
		}},
	}

	result, _ := parseTokens(t, input)
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
	targets := []testNode{
		{"assignment", p.NodeData{VarName: "x"}, []testNode{
			{"variable_reference", p.NodeData{VarName: "arr", Index: "0"}, []testNode{}},
		}},
	}

	result, _ := parseTokens(t, input)
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
	targets := []testNode{
		{"assignment", p.NodeData{VarName: "arr", Index: "1"}, []testNode{
			{"string", p.NodeData{Content: "x"}, []testNode{}},
		}},
	}

	result, _ := parseTokens(t, input)
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
	targets := []testNode{
		{"assignment", p.NodeData{VarName: "arr", Index: "i"}, []testNode{
			{"expression", p.NodeData{Operator: "+"}, []testNode{
				{"lhs", p.NodeData{}, []testNode{
					{"variable_reference", p.NodeData{VarName: "arr", Index: "i"}, []testNode{}},
				}},
				{"rhs", p.NodeData{}, []testNode{
					{"number", p.NodeData{Content: "1"}, []testNode{}},
				}},
			}},
		}},
	}

	result, _ := parseTokens(t, input)
	assert.Equal(t, targets, result)
}

func Test_Vector_Literal(t *testing.T) {
	input := []l.Token{
		{Type: l.SYMBOL, Content: "pos"},
		{Type: l.ASSIGNMENT, Content: "="},
		{Type: l.OPEN_PAREN, Content: "("},
		{Type: l.NUMBER, Content: "0"},
		{Type: l.COMMA, Content: ","},
		{Type: l.NUMBER, Content: "1"},
		{Type: l.COMMA, Content: ","},
		{Type: l.NUMBER, Content: "2"},
		{Type: l.CLOSE_PAREN, Content: ")"},
		{Type: l.TERMINATOR, Content: ";"},
	}
	targets := []testNode{
		{"assignment", p.NodeData{VarName: "pos"}, []testNode{
			{"vector_literal", p.NodeData{}, []testNode{
				{"number", p.NodeData{Content: "0"}, []testNode{}},
				{"number", p.NodeData{Content: "1"}, []testNode{}},
				{"number", p.NodeData{Content: "2"}, []testNode{}},
			}},
		}},
	}

	result, _ := parseTokens(t, input)
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
	target := []testNode{
		{"function_call", p.NodeData{FunctionName: "init"}, []testNode{
			{"variable_reference", p.NodeData{VarName: "arg1"}, []testNode{}},
			{"expression", p.NodeData{Operator: "+"}, []testNode{
				{"lhs", p.NodeData{}, []testNode{
					{"string", p.NodeData{Content: "Hello "}, []testNode{}},
				}},
				{"rhs", p.NodeData{}, []testNode{
					{"variable_reference", p.NodeData{VarName: "name"}, []testNode{}},
				}},
			}},
		}},
	}

	result, _ := parseTokens(t, input)
	assert.Equal(t, target, result)
}

func Test_Namespace_Function_Call(t *testing.T) {
	input := []l.Token{
		{Type: l.SYMBOL, Content: "maps\\mp\\zombies\\_zm_powerups::specific_powerup_drop"},
		{Type: l.OPEN_PAREN, Content: "("},
		{Type: l.CLOSE_PAREN, Content: ")"},
	}
	target := []testNode{
		{"function_call", p.NodeData{FunctionName: "specific_powerup_drop", Path: "maps\\mp\\zombies\\_zm_powerups"}, []testNode{}},
	}

	result, _ := parseTokens(t, input)
	assert.Equal(t, target, result)
}

func Test_Method_Function_Call(t *testing.T) {
	input := []l.Token{
		{Type: l.SYMBOL, Content: "self.method"},
		{Type: l.OPEN_PAREN, Content: "("},
		{Type: l.CLOSE_PAREN, Content: ")"},
	}
	target := []testNode{
		{"function_call", p.NodeData{FunctionName: "method", Method: "self"}, []testNode{}},
	}

	result, _ := parseTokens(t, input)
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
	target := []testNode{
		{"function_call", p.NodeData{FunctionName: "some_func", Thread: true}, []testNode{}},
	}

	result, _ := parseTokens(t, input)
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
	target := []testNode{
		{"function_declaration", p.NodeData{FunctionName: "init"}, []testNode{
			{"args", p.NodeData{}, []testNode{
				{"variable_reference", p.NodeData{VarName: "name"}, []testNode{}},
				{"variable_reference", p.NodeData{VarName: "age"}, []testNode{}},
			}},
			{"scope", p.NodeData{}, []testNode{
				{"function_call", p.NodeData{FunctionName: "print"}, []testNode{
					{"variable_reference", p.NodeData{VarName: "name"}, []testNode{}},
				}},
			}},
		}},
	}

	result, _ := parseTokens(t, input)
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
	target := []testNode{
		{"for_loop", p.NodeData{}, []testNode{
			{"for_init", p.NodeData{}, []testNode{}},
			{"for_condition", p.NodeData{}, []testNode{}},
			{"for_post", p.NodeData{}, []testNode{}},
			{"scope", p.NodeData{}, []testNode{}},
		}},
	}

	result, _ := parseTokens(t, input)
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
	target := []testNode{
		{"for_loop", p.NodeData{}, []testNode{
			{"for_init", p.NodeData{}, []testNode{
				{"assignment", p.NodeData{VarName: "i"}, []testNode{
					{"number", p.NodeData{Content: "0"}, []testNode{}},
				}},
			}},
			{"for_condition", p.NodeData{}, []testNode{
				{"expression", p.NodeData{Operator: "<"}, []testNode{
					{"lhs", p.NodeData{}, []testNode{
						{"variable_reference", p.NodeData{VarName: "i"}, []testNode{}},
					}},
					{"rhs", p.NodeData{}, []testNode{
						{"number", p.NodeData{Content: "10"}, []testNode{}},
					}},
				}},
			}},
			{"for_post", p.NodeData{}, []testNode{
				{"assignment", p.NodeData{VarName: "i"}, []testNode{
					{"expression", p.NodeData{Operator: "+"}, []testNode{
						{"lhs", p.NodeData{}, []testNode{
							{"variable_reference", p.NodeData{VarName: "i"}, []testNode{}},
						}},
						{"rhs", p.NodeData{}, []testNode{
							{"number", p.NodeData{Content: "1"}, []testNode{}},
						}},
					}},
				}},
			}},
			{"scope", p.NodeData{}, []testNode{}},
		}},
	}

	result, _ := parseTokens(t, input)
	assert.Equal(t, target, result)
}

func Test_If_Else(t *testing.T) {
	input := []l.Token{
		{Type: l.SYMBOL, Content: "if"},
		{Type: l.OPEN_PAREN, Content: "("},
		{Type: l.SYMBOL, Content: "cond"},
		{Type: l.CLOSE_PAREN, Content: ")"},
		{Type: l.OPEN_CURLY, Content: "{"},
		{Type: l.CLOSE_CURLY, Content: "}"},
		{Type: l.SYMBOL, Content: "else"},
		{Type: l.OPEN_CURLY, Content: "{"},
		{Type: l.CLOSE_CURLY, Content: "}"},
	}
	target := []testNode{
		{"if_statement", p.NodeData{}, []testNode{
			{"condition", p.NodeData{}, []testNode{
				{"variable_reference", p.NodeData{VarName: "cond"}, []testNode{}},
			}},
			{"scope", p.NodeData{}, []testNode{}},
		}},
		{"else_clause", p.NodeData{}, []testNode{
			{"scope", p.NodeData{}, []testNode{}},
		}},
	}

	result, _ := parseTokens(t, input)
	assert.Equal(t, target, result)
}

func Test_While_Loop(t *testing.T) {
	input := []l.Token{
		{Type: l.SYMBOL, Content: "while"},
		{Type: l.OPEN_PAREN, Content: "("},
		{Type: l.SYMBOL, Content: "running"},
		{Type: l.CLOSE_PAREN, Content: ")"},
		{Type: l.OPEN_CURLY, Content: "{"},
		{Type: l.CLOSE_CURLY, Content: "}"},
	}
	target := []testNode{
		{"while_loop", p.NodeData{}, []testNode{
			{"condition", p.NodeData{}, []testNode{{"variable_reference", p.NodeData{VarName: "running"}, []testNode{}}}},
			{"scope", p.NodeData{}, []testNode{}},
		}},
	}

	result, _ := parseTokens(t, input)
	assert.Equal(t, target, result)
}

func Test_Foreach_Loop(t *testing.T) {
	input := []l.Token{
		{Type: l.SYMBOL, Content: "foreach"},
		{Type: l.OPEN_PAREN, Content: "("},
		{Type: l.SYMBOL, Content: "item"},
		{Type: l.SYMBOL, Content: "in"},
		{Type: l.SYMBOL, Content: "items"},
		{Type: l.CLOSE_PAREN, Content: ")"},
		{Type: l.OPEN_CURLY, Content: "{"},
		{Type: l.CLOSE_CURLY, Content: "}"},
	}
	target := []testNode{
		{"foreach_loop", p.NodeData{}, []testNode{
			{"foreach_vars", p.NodeData{}, []testNode{{"variable_reference", p.NodeData{VarName: "item"}, []testNode{}}}},
			{"foreach_iter", p.NodeData{}, []testNode{{"variable_reference", p.NodeData{VarName: "items"}, []testNode{}}}},
			{"scope", p.NodeData{}, []testNode{}},
		}},
	}

	result, _ := parseTokens(t, input)
	assert.Equal(t, target, result)
}

func Test_Switch_Case_Default(t *testing.T) {
	input := []l.Token{
		{Type: l.SYMBOL, Content: "switch"},
		{Type: l.OPEN_PAREN, Content: "("},
		{Type: l.SYMBOL, Content: "x"},
		{Type: l.CLOSE_PAREN, Content: ")"},
		{Type: l.OPEN_CURLY, Content: "{"},
		{Type: l.SYMBOL, Content: "case"},
		{Type: l.NUMBER, Content: "1"},
		{Type: l.TERMINATOR, Content: ";"},
		{Type: l.SYMBOL, Content: "break"},
		{Type: l.TERMINATOR, Content: ";"},
		{Type: l.SYMBOL, Content: "default"},
		{Type: l.CLOSE_CURLY, Content: "}"},
	}
	target := []testNode{
		{"switch_statement", p.NodeData{}, []testNode{
			{"switch_expr", p.NodeData{}, []testNode{{"variable_reference", p.NodeData{VarName: "x"}, []testNode{}}}},
			{"scope", p.NodeData{}, []testNode{
				{"case_clause", p.NodeData{}, []testNode{{"number", p.NodeData{Content: "1"}, []testNode{}}}},
				{"break_statement", p.NodeData{}, []testNode{}},
				{"default_clause", p.NodeData{}, []testNode{}},
			}},
		}},
	}

	result, _ := parseTokens(t, input)
	assert.Equal(t, target, result)
}

func Test_Return_Statement(t *testing.T) {
	input := []l.Token{
		{Type: l.SYMBOL, Content: "return"},
		{Type: l.SYMBOL, Content: "value"},
		{Type: l.TERMINATOR, Content: ";"},
	}
	target := []testNode{
		{"return_statement", p.NodeData{}, []testNode{{"variable_reference", p.NodeData{VarName: "value"}, []testNode{}}}},
	}

	result, _ := parseTokens(t, input)
	assert.Equal(t, target, result)
}

func Test_IncludeStatement(t *testing.T) {
	input := []l.Token{
		{Type: l.SYMBOL, Content: "#include"},
		{Type: l.SYMBOL, Content: "path\\to\\file"},
		{Type: l.TERMINATOR, Content: ";"},
	}
	target := []testNode{
		{"include_statement", p.NodeData{Path: "path\\to\\file"}, []testNode{}},
	}

	result, _ := parseTokens(t, input)
	assert.Equal(t, target, result)
}

func Test_WaitStatement(t *testing.T) {
	input := []l.Token{
		{Type: l.SYMBOL, Content: "wait"},
		{Type: l.NUMBER, Content: "0.05"},
		{Type: l.TERMINATOR, Content: ";"},
	}
	target := []testNode{
		{"wait_statement", p.NodeData{Delay: "0.05"}, []testNode{}},
	}

	result, _ := parseTokens(t, input)
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

	target := []testNode{
		{"function_call", p.NodeData{FunctionName: "somefunc", Thread: true, Method: "level"}, []testNode{}},
		{"function_call", p.NodeData{FunctionName: "somefunc", Thread: true}, []testNode{}},
		{"function_call", p.NodeData{FunctionName: "somefunc", Method: "level"}, []testNode{}},
		{"function_call", p.NodeData{FunctionName: "somefunc"}, []testNode{}},
	}

	result, _ := parseTokens(t, input)
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

	target := []testNode{
		{"function_call", p.NodeData{FunctionName: "somefunc"}, []testNode{
			{"variable_reference", p.NodeData{VarName: "somefunc_arg"}, []testNode{}},
			{"function_call", p.NodeData{FunctionName: "child_func"}, []testNode{
				{"variable_reference", p.NodeData{VarName: "child_arg1"}, []testNode{}},
				{"variable_reference", p.NodeData{VarName: "child_arg2"}, []testNode{}},
			}},
		}},
	}

	result, _ := parseTokens(t, input)
	assert.Equal(t, target, result)
}

func Test_Comments_AreParsedAsNodes(t *testing.T) {
	input := []l.Token{
		{Type: l.SYMBOL, Content: "a"},
		{Type: l.ASSIGNMENT, Content: "="},
		{Type: l.NUMBER, Content: "1"},
		{Type: l.TERMINATOR, Content: ";"},
		{Type: l.LINE_COMMENT, Content: "// keep me"},
		{Type: l.NEWLINE, Content: ""},
		{Type: l.BLOCK_COMMENT, Content: "/* and me */"},
	}

	target := []testNode{
		{"assignment", p.NodeData{VarName: "a"}, []testNode{{"number", p.NodeData{Content: "1"}, []testNode{}}}},
		{"comment", p.NodeData{Content: "// keep me"}, []testNode{}},
		{"comment", p.NodeData{Content: "/* and me */"}, []testNode{}},
	}

	result, _ := parseTokens(t, input)
	assert.Equal(t, target, result)
}
