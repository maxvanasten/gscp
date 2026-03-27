package generator_test

import (
	"testing"

	g "github.com/maxvanasten/gscp/generator"
	p "github.com/maxvanasten/gscp/parser"
	"github.com/stretchr/testify/assert"
)

type testNode struct {
	Type     string
	Data     p.NodeData
	Children []testNode
}

func toParserNode(node testNode) p.Node {
	children := make([]p.Node, len(node.Children))
	for i, child := range node.Children {
		children[i] = toParserNode(child)
	}
	return p.Node{Type: node.Type, Data: node.Data, Children: children}
}

func generate(node testNode) string {
	return g.Generate(toParserNode(node))
}

func indentLines(lines ...string) string {
	output := ""
	for i, line := range lines {
		if i > 0 {
			output += "\n"
		}
		output += g.Indent + line
	}
	return output
}

func Test_Generate_VariableReference(t *testing.T) {
	input := testNode{"variable_reference", p.NodeData{VarName: "test_var"}, []testNode{}}

	target := "test_var"

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_String(t *testing.T) {
	input := testNode{"string", p.NodeData{Content: "hello, world"}, []testNode{}}

	target := "\"hello, world\""

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_Boolean(t *testing.T) {
	input := testNode{"boolean", p.NodeData{Content: "true"}, []testNode{}}

	target := "true"

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_Number(t *testing.T) {
	input := testNode{"number", p.NodeData{Content: "23"}, []testNode{}}
	target := "23"

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_Simple_Expression(t *testing.T) {
	input := testNode{"expression", p.NodeData{Operator: "+"}, []testNode{
		{"lhs", p.NodeData{}, []testNode{
			{"string", p.NodeData{Content: "Hello, "}, []testNode{}},
		}},
		{"rhs", p.NodeData{}, []testNode{
			{"variable_reference", p.NodeData{VarName: "name"}, []testNode{}},
		}},
	}}
	target := "\"Hello, \" + name"

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_Complex_Expression(t *testing.T) {
	input := testNode{"expression", p.NodeData{Operator: "+"}, []testNode{
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
	}}
	target := "\"Hello, \" + name + \", You are: \" + 23 + \" years old.\""

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_Assignment(t *testing.T) {
	input := testNode{"assignment", p.NodeData{VarName: "test"}, []testNode{
		{"string", p.NodeData{Content: "Hello, world!"}, []testNode{}},
	}}
	target := "test = \"Hello, world!\";"

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_Unary_Expression(t *testing.T) {
	input := testNode{"unary_expression", p.NodeData{Operator: "!"}, []testNode{
		{"boolean", p.NodeData{Content: "true"}, []testNode{}},
	}}
	target := "!true"

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_AssignmentWithUnaryExpressionBinaryLHS(t *testing.T) {
	input := testNode{"assignment", p.NodeData{VarName: "x"}, []testNode{
		{"expression", p.NodeData{Operator: "-"}, []testNode{
			{"lhs", p.NodeData{}, []testNode{
				{"unary_expression", p.NodeData{Operator: "-"}, []testNode{
					{"number", p.NodeData{Content: "130"}, []testNode{}},
				}},
			}},
			{"rhs", p.NodeData{}, []testNode{
				{"expression", p.NodeData{Operator: "*"}, []testNode{
					{"lhs", p.NodeData{}, []testNode{
						{"variable_reference", p.NodeData{VarName: "i"}, []testNode{}},
					}},
					{"rhs", p.NodeData{}, []testNode{
						{"number", p.NodeData{Content: "10"}, []testNode{}},
					}},
				}},
			}},
		}},
	}}
	target := "x = -130 - i * 10;"

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_Lhs(t *testing.T) {
	input := testNode{"lhs", p.NodeData{}, []testNode{
		{"variable_reference", p.NodeData{VarName: "x"}, []testNode{}},
	}}
	target := "x"

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_Rhs(t *testing.T) {
	input := testNode{"rhs", p.NodeData{}, []testNode{
		{"number", p.NodeData{Content: "1"}, []testNode{}},
	}}
	target := "1"

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_Array_Literal(t *testing.T) {
	input := testNode{"assignment", p.NodeData{VarName: "x"}, []testNode{
		{"array_literal", p.NodeData{}, []testNode{}},
	}}

	target := "x = [];"

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_Array_Literal_Multiple(t *testing.T) {
	input := testNode{"assignment", p.NodeData{VarName: "x"}, []testNode{
		{"array_literal", p.NodeData{}, []testNode{
			{"number", p.NodeData{Content: "1"}, []testNode{}},
			{"string", p.NodeData{Content: "a"}, []testNode{}},
			{"variable_reference", p.NodeData{VarName: "y"}, []testNode{}},
		}},
	}}
	target := "x = [1, \"a\", y];"

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_ArrayLiteral(t *testing.T) {
	input := testNode{"array_literal", p.NodeData{}, []testNode{
		{"number", p.NodeData{Content: "1"}, []testNode{}},
		{"string", p.NodeData{Content: "a"}, []testNode{}},
		{"variable_reference", p.NodeData{VarName: "y"}, []testNode{}},
	}}
	target := "[1, \"a\", y]"

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_Array_Indexing(t *testing.T) {
	input := testNode{"variable_reference", p.NodeData{VarName: "arr", Index: "0"}, []testNode{}}
	target := "arr[0]"

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_Array_Indexing_Assignment(t *testing.T) {
	input := testNode{"assignment", p.NodeData{VarName: "arr", Index: "1"}, []testNode{
		{"string", p.NodeData{Content: "x"}, []testNode{}},
	}}
	target := "arr[1] = \"x\";"

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_IncludeStatement(t *testing.T) {
	input := testNode{"include_statement", p.NodeData{Path: "common_scripts\\utility"}, []testNode{}}

	target := "#include common_scripts\\utility;"

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_WaitStatement(t *testing.T) {
	input := testNode{"wait_statement", p.NodeData{Delay: "0.05"}, []testNode{}}

	target := "wait 0.05;"

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_BreakStatement(t *testing.T) {
	input := testNode{"break_statement", p.NodeData{}, []testNode{}}

	target := "break;"

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_ReturnStatement(t *testing.T) {
	input := testNode{"return_statement", p.NodeData{}, []testNode{
		{"string", p.NodeData{Content: "value"}, []testNode{}},
	}}

	target := "return \"value\";"

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_FunctionCall(t *testing.T) {
	input := testNode{"function_call", p.NodeData{FunctionName: "do_thing"}, []testNode{}}

	target := "do_thing();"

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_FunctionCall_WithArgs(t *testing.T) {
	input := testNode{"function_call", p.NodeData{FunctionName: "do_thing"}, []testNode{
		{"variable_reference", p.NodeData{VarName: "a"}, []testNode{}},
		{"number", p.NodeData{Content: "1"}, []testNode{}},
		{"string", p.NodeData{Content: "x"}, []testNode{}},
	}}

	target := "do_thing(a, 1, \"x\");"

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_FunctionCall_Method(t *testing.T) {
	input := testNode{"function_call", p.NodeData{FunctionName: "do_thing", Method: "self"}, []testNode{}}

	target := "self do_thing();"

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_FunctionCall_MethodWithIndex(t *testing.T) {
	input := testNode{"function_call", p.NodeData{FunctionName: "ml_update_text", Method: "self.hud_perks[i]"}, []testNode{
		{"variable_reference", p.NodeData{VarName: "perk"}, []testNode{}},
	}}

	target := "self.hud_perks[i] ml_update_text(perk);"

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_FunctionCall_Thread(t *testing.T) {
	input := testNode{"function_call", p.NodeData{FunctionName: "do_thing", Thread: true}, []testNode{}}

	target := "thread do_thing();"

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_FunctionCall_MethodThread(t *testing.T) {
	input := testNode{"function_call", p.NodeData{FunctionName: "do_thing", Method: "self", Thread: true}, []testNode{}}

	target := "self thread do_thing();"

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_FunctionCall_Path(t *testing.T) {
	input := testNode{"function_call", p.NodeData{FunctionName: "specific_powerup_drop", Path: "maps\\mp\\zombies\\_zm_powerups"}, []testNode{}}

	target := "maps\\mp\\zombies\\_zm_powerups::specific_powerup_drop();"

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_FunctionDeclaration(t *testing.T) {
	input := testNode{"function_declaration", p.NodeData{FunctionName: "init"}, []testNode{
		{"args", p.NodeData{}, []testNode{
			{"variable_reference", p.NodeData{VarName: "a"}, []testNode{}},
			{"variable_reference", p.NodeData{VarName: "b"}, []testNode{}},
		}},
		{"scope", p.NodeData{}, []testNode{
			{"wait_statement", p.NodeData{Delay: "0.05"}, []testNode{}},
		}},
	}}

	target := "init(a, b)\n{\n" + indentLines("wait 0.05;") + "\n}"

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_Args(t *testing.T) {
	input := testNode{"args", p.NodeData{}, []testNode{
		{"variable_reference", p.NodeData{VarName: "a"}, []testNode{}},
		{"variable_reference", p.NodeData{VarName: "b"}, []testNode{}},
	}}

	target := "a, b"

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_Scope(t *testing.T) {
	input := testNode{"scope", p.NodeData{}, []testNode{
		{"wait_statement", p.NodeData{Delay: "0.05"}, []testNode{}},
		{"comment", p.NodeData{Content: "// keep"}, []testNode{}},
		{"break_statement", p.NodeData{}, []testNode{}},
	}}

	target := indentLines("wait 0.05;", "// keep", "break;")

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_Comment(t *testing.T) {
	input := testNode{"comment", p.NodeData{Content: "/* keep me */"}, []testNode{}}
	target := "/* keep me */"

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_VectorLiteral(t *testing.T) {
	input := testNode{"vector_literal", p.NodeData{}, []testNode{
		{"number", p.NodeData{Content: "0"}, []testNode{}},
		{"number", p.NodeData{Content: "1"}, []testNode{}},
		{"number", p.NodeData{Content: "2"}, []testNode{}},
	}}

	target := "(0, 1, 2)"

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_ForInit(t *testing.T) {
	input := testNode{"for_init", p.NodeData{}, []testNode{
		{"assignment", p.NodeData{VarName: "i"}, []testNode{
			{"number", p.NodeData{Content: "0"}, []testNode{}},
		}},
	}}

	target := "i = 0"

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_ForCondition(t *testing.T) {
	input := testNode{"for_condition", p.NodeData{}, []testNode{
		{"expression", p.NodeData{Operator: "<"}, []testNode{
			{"lhs", p.NodeData{}, []testNode{
				{"variable_reference", p.NodeData{VarName: "i"}, []testNode{}},
			}},
			{"rhs", p.NodeData{}, []testNode{
				{"number", p.NodeData{Content: "10"}, []testNode{}},
			}},
		}},
	}}

	target := "i < 10"

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_ForPost(t *testing.T) {
	input := testNode{"for_post", p.NodeData{}, []testNode{
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
	}}

	target := "i = i + 1"

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_ForLoop(t *testing.T) {
	input := testNode{"for_loop", p.NodeData{}, []testNode{
		{"for_init", p.NodeData{}, []testNode{}},
		{"for_condition", p.NodeData{}, []testNode{}},
		{"for_post", p.NodeData{}, []testNode{}},
		{"scope", p.NodeData{}, []testNode{
			{"wait_statement", p.NodeData{Delay: "0.05"}, []testNode{}},
		}},
	}}

	target := "for ( ;; )\n{\n" + indentLines("wait 0.05;") + "\n}"

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_Condition(t *testing.T) {
	input := testNode{"condition", p.NodeData{}, []testNode{
		{"variable_reference", p.NodeData{VarName: "cond"}, []testNode{}},
	}}

	target := "cond"

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_IfStatement(t *testing.T) {
	input := testNode{"if_statement", p.NodeData{}, []testNode{
		{"condition", p.NodeData{}, []testNode{
			{"variable_reference", p.NodeData{VarName: "cond"}, []testNode{}},
		}},
		{"scope", p.NodeData{}, []testNode{
			{"wait_statement", p.NodeData{Delay: "0.05"}, []testNode{}},
		}},
	}}

	target := "if (cond)\n{\n" + indentLines("wait 0.05;") + "\n}"

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_ElseClause(t *testing.T) {
	input := testNode{"else_clause", p.NodeData{}, []testNode{
		{"scope", p.NodeData{}, []testNode{
			{"wait_statement", p.NodeData{Delay: "0.05"}, []testNode{}},
		}},
	}}

	target := "else\n{\n" + indentLines("wait 0.05;") + "\n}"

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_WhileLoop(t *testing.T) {
	input := testNode{"while_loop", p.NodeData{}, []testNode{
		{"condition", p.NodeData{}, []testNode{
			{"boolean", p.NodeData{Content: "true"}, []testNode{}},
		}},
		{"scope", p.NodeData{}, []testNode{
			{"wait_statement", p.NodeData{Delay: "0.05"}, []testNode{}},
		}},
	}}

	target := "while (true)\n{\n" + indentLines("wait 0.05;") + "\n}"

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_ForeachVars(t *testing.T) {
	input := testNode{"foreach_vars", p.NodeData{}, []testNode{
		{"variable_reference", p.NodeData{VarName: "item"}, []testNode{}},
	}}

	target := "item"

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_ForeachIter(t *testing.T) {
	input := testNode{"foreach_iter", p.NodeData{}, []testNode{
		{"variable_reference", p.NodeData{VarName: "items"}, []testNode{}},
	}}

	target := "items"

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_ForeachLoop(t *testing.T) {
	input := testNode{"foreach_loop", p.NodeData{}, []testNode{
		{"foreach_vars", p.NodeData{}, []testNode{
			{"variable_reference", p.NodeData{VarName: "item"}, []testNode{}},
		}},
		{"foreach_iter", p.NodeData{}, []testNode{
			{"variable_reference", p.NodeData{VarName: "items"}, []testNode{}},
		}},
		{"scope", p.NodeData{}, []testNode{
			{"wait_statement", p.NodeData{Delay: "0.05"}, []testNode{}},
		}},
	}}

	target := "foreach (item in items)\n{\n" + indentLines("wait 0.05;") + "\n}"

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_SwitchExpr(t *testing.T) {
	input := testNode{"switch_expr", p.NodeData{}, []testNode{
		{"variable_reference", p.NodeData{VarName: "x"}, []testNode{}},
	}}

	target := "x"

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_CaseClause(t *testing.T) {
	input := testNode{"case_clause", p.NodeData{}, []testNode{
		{"string", p.NodeData{Content: "a"}, []testNode{}},
	}}

	target := "case \"a\":"

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_DefaultClause(t *testing.T) {
	input := testNode{"default_clause", p.NodeData{}, []testNode{}}

	target := "default:"

	result := generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_SwitchStatement(t *testing.T) {
	input := testNode{"switch_statement", p.NodeData{}, []testNode{
		{"switch_expr", p.NodeData{}, []testNode{
			{"variable_reference", p.NodeData{VarName: "x"}, []testNode{}},
		}},
		{"scope", p.NodeData{}, []testNode{
			{"case_clause", p.NodeData{}, []testNode{
				{"string", p.NodeData{Content: "a"}, []testNode{}},
			}},
			{"wait_statement", p.NodeData{Delay: "0.05"}, []testNode{}},
			{"break_statement", p.NodeData{}, []testNode{}},
			{"default_clause", p.NodeData{}, []testNode{}},
		}},
	}}

	inner := indentLines("case \"a\":", g.Indent+"wait 0.05;", g.Indent+"break;", "default:")
	target := "switch(x) {\n" + inner + "\n}"

	result := generate(input)
	assert.Equal(t, target, result)
}
