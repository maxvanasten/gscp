package generator_test

import (
	"testing"

	g "github.com/maxvanasten/gscp/generator"
	p "github.com/maxvanasten/gscp/parser"
	"github.com/stretchr/testify/assert"
)

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
	input := p.Node{"variable_reference", p.NodeData{VarName: "test_var"}, []p.Node{}}

	target := "test_var"

	result := g.Generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_String(t *testing.T) {
	input := p.Node{"string", p.NodeData{Content: "hello, world"}, []p.Node{}}

	target := "\"hello, world\""

	result := g.Generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_Boolean(t *testing.T) {
	input := p.Node{"boolean", p.NodeData{Content: "true"}, []p.Node{}}

	target := "true"

	result := g.Generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_Number(t *testing.T) {
	input := p.Node{"number", p.NodeData{Content: "23"}, []p.Node{}}
	target := "23"

	result := g.Generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_Simple_Expression(t *testing.T) {
	input := p.Node{"expression", p.NodeData{Operator: "+"}, []p.Node{
		{"lhs", p.NodeData{}, []p.Node{
			{"string", p.NodeData{Content: "Hello, "}, []p.Node{}},
		}},
		{"rhs", p.NodeData{}, []p.Node{
			{"variable_reference", p.NodeData{VarName: "name"}, []p.Node{}},
		}},
	}}
	target := "\"Hello, \" + name"

	result := g.Generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_Complex_Expression(t *testing.T) {
	input := p.Node{"expression", p.NodeData{Operator: "+"}, []p.Node{
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
	}}
	target := "\"Hello, \" + name + \", You are: \" + 23 + \" years old.\""

	result := g.Generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_Assignment(t *testing.T) {
	input := p.Node{"assignment", p.NodeData{VarName: "test"}, []p.Node{
		{"string", p.NodeData{Content: "Hello, world!"}, []p.Node{}},
	}}
	target := "test = \"Hello, world!\";"

	result := g.Generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_Unary_Expression(t *testing.T) {
	input := p.Node{"unary_expression", p.NodeData{Operator: "!"}, []p.Node{
		{"boolean", p.NodeData{Content: "true"}, []p.Node{}},
	}}
	target := "!true"

	result := g.Generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_Lhs(t *testing.T) {
	input := p.Node{"lhs", p.NodeData{}, []p.Node{
		{"variable_reference", p.NodeData{VarName: "x"}, []p.Node{}},
	}}
	target := "x"

	result := g.Generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_Rhs(t *testing.T) {
	input := p.Node{"rhs", p.NodeData{}, []p.Node{
		{"number", p.NodeData{Content: "1"}, []p.Node{}},
	}}
	target := "1"

	result := g.Generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_Array_Literal(t *testing.T) {
	input := p.Node{"assignment", p.NodeData{VarName: "x"}, []p.Node{
		{"array_literal", p.NodeData{}, []p.Node{}},
	}}

	target := "x = [];"

	result := g.Generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_Array_Literal_Multiple(t *testing.T) {
	input := p.Node{"assignment", p.NodeData{VarName: "x"}, []p.Node{
		{"array_literal", p.NodeData{}, []p.Node{
			{"number", p.NodeData{Content: "1"}, []p.Node{}},
			{"string", p.NodeData{Content: "a"}, []p.Node{}},
			{"variable_reference", p.NodeData{VarName: "y"}, []p.Node{}},
		}},
	}}
	target := "x = [1, \"a\", y];"

	result := g.Generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_ArrayLiteral(t *testing.T) {
	input := p.Node{"array_literal", p.NodeData{}, []p.Node{
		{"number", p.NodeData{Content: "1"}, []p.Node{}},
		{"string", p.NodeData{Content: "a"}, []p.Node{}},
		{"variable_reference", p.NodeData{VarName: "y"}, []p.Node{}},
	}}
	target := "[1, \"a\", y]"

	result := g.Generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_Array_Indexing(t *testing.T) {
	input := p.Node{"variable_reference", p.NodeData{VarName: "arr", Index: "0"}, []p.Node{}}
	target := "arr[0]"

	result := g.Generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_Array_Indexing_Assignment(t *testing.T) {
	input := p.Node{"assignment", p.NodeData{VarName: "arr", Index: "1"}, []p.Node{
		{"string", p.NodeData{Content: "x"}, []p.Node{}},
	}}
	target := "arr[1] = \"x\";"

	result := g.Generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_IncludeStatement(t *testing.T) {
	input := p.Node{"include_statement", p.NodeData{Path: "common_scripts\\utility"}, []p.Node{}}

	target := "#include common_scripts\\utility;"

	result := g.Generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_WaitStatement(t *testing.T) {
	input := p.Node{"wait_statement", p.NodeData{Delay: "0.05"}, []p.Node{}}

	target := "wait 0.05;"

	result := g.Generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_BreakStatement(t *testing.T) {
	input := p.Node{"break_statement", p.NodeData{}, []p.Node{}}

	target := "break;"

	result := g.Generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_ReturnStatement(t *testing.T) {
	input := p.Node{"return_statement", p.NodeData{}, []p.Node{
		{"string", p.NodeData{Content: "value"}, []p.Node{}},
	}}

	target := "return \"value\";"

	result := g.Generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_FunctionCall(t *testing.T) {
	input := p.Node{"function_call", p.NodeData{FunctionName: "do_thing"}, []p.Node{}}

	target := "do_thing();"

	result := g.Generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_FunctionCall_WithArgs(t *testing.T) {
	input := p.Node{"function_call", p.NodeData{FunctionName: "do_thing"}, []p.Node{
		{"variable_reference", p.NodeData{VarName: "a"}, []p.Node{}},
		{"number", p.NodeData{Content: "1"}, []p.Node{}},
		{"string", p.NodeData{Content: "x"}, []p.Node{}},
	}}

	target := "do_thing(a, 1, \"x\");"

	result := g.Generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_FunctionCall_Method(t *testing.T) {
	input := p.Node{"function_call", p.NodeData{FunctionName: "do_thing", Method: "self"}, []p.Node{}}

	target := "self do_thing();"

	result := g.Generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_FunctionCall_Thread(t *testing.T) {
	input := p.Node{"function_call", p.NodeData{FunctionName: "do_thing", Thread: true}, []p.Node{}}

	target := "thread do_thing();"

	result := g.Generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_FunctionCall_MethodThread(t *testing.T) {
	input := p.Node{"function_call", p.NodeData{FunctionName: "do_thing", Method: "self", Thread: true}, []p.Node{}}

	target := "self thread do_thing();"

	result := g.Generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_FunctionCall_Path(t *testing.T) {
	input := p.Node{"function_call", p.NodeData{FunctionName: "specific_powerup_drop", Path: "maps\\mp\\zombies\\_zm_powerups"}, []p.Node{}}

	target := "maps\\mp\\zombies\\_zm_powerups::specific_powerup_drop();"

	result := g.Generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_FunctionDeclaration(t *testing.T) {
	input := p.Node{"function_declaration", p.NodeData{FunctionName: "init"}, []p.Node{
		{"args", p.NodeData{}, []p.Node{
			{"variable_reference", p.NodeData{VarName: "a"}, []p.Node{}},
			{"variable_reference", p.NodeData{VarName: "b"}, []p.Node{}},
		}},
		{"scope", p.NodeData{}, []p.Node{
			{"wait_statement", p.NodeData{Delay: "0.05"}, []p.Node{}},
		}},
	}}

	target := "init(a, b)\n{\n" + indentLines("wait 0.05;") + "\n}"

	result := g.Generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_Args(t *testing.T) {
	input := p.Node{"args", p.NodeData{}, []p.Node{
		{"variable_reference", p.NodeData{VarName: "a"}, []p.Node{}},
		{"variable_reference", p.NodeData{VarName: "b"}, []p.Node{}},
	}}

	target := "a, b"

	result := g.Generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_Scope(t *testing.T) {
	input := p.Node{"scope", p.NodeData{}, []p.Node{
		{"wait_statement", p.NodeData{Delay: "0.05"}, []p.Node{}},
		{"break_statement", p.NodeData{}, []p.Node{}},
	}}

	target := indentLines("wait 0.05;", "break;")

	result := g.Generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_VectorLiteral(t *testing.T) {
	input := p.Node{"vector_literal", p.NodeData{}, []p.Node{
		{"number", p.NodeData{Content: "0"}, []p.Node{}},
		{"number", p.NodeData{Content: "1"}, []p.Node{}},
		{"number", p.NodeData{Content: "2"}, []p.Node{}},
	}}

	target := "(0, 1, 2)"

	result := g.Generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_ForInit(t *testing.T) {
	input := p.Node{"for_init", p.NodeData{}, []p.Node{
		{"assignment", p.NodeData{VarName: "i"}, []p.Node{
			{"number", p.NodeData{Content: "0"}, []p.Node{}},
		}},
	}}

	target := "i = 0"

	result := g.Generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_ForCondition(t *testing.T) {
	input := p.Node{"for_condition", p.NodeData{}, []p.Node{
		{"expression", p.NodeData{Operator: "<"}, []p.Node{
			{"lhs", p.NodeData{}, []p.Node{
				{"variable_reference", p.NodeData{VarName: "i"}, []p.Node{}},
			}},
			{"rhs", p.NodeData{}, []p.Node{
				{"number", p.NodeData{Content: "10"}, []p.Node{}},
			}},
		}},
	}}

	target := "i < 10"

	result := g.Generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_ForPost(t *testing.T) {
	input := p.Node{"for_post", p.NodeData{}, []p.Node{
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
	}}

	target := "i = i + 1"

	result := g.Generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_ForLoop(t *testing.T) {
	input := p.Node{"for_loop", p.NodeData{}, []p.Node{
		{"for_init", p.NodeData{}, []p.Node{}},
		{"for_condition", p.NodeData{}, []p.Node{}},
		{"for_post", p.NodeData{}, []p.Node{}},
		{"scope", p.NodeData{}, []p.Node{
			{"wait_statement", p.NodeData{Delay: "0.05"}, []p.Node{}},
		}},
	}}

	target := "for ( ;; )\n{\n" + indentLines("wait 0.05;") + "\n}"

	result := g.Generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_Condition(t *testing.T) {
	input := p.Node{"condition", p.NodeData{}, []p.Node{
		{"variable_reference", p.NodeData{VarName: "cond"}, []p.Node{}},
	}}

	target := "cond"

	result := g.Generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_IfStatement(t *testing.T) {
	input := p.Node{"if_statement", p.NodeData{}, []p.Node{
		{"condition", p.NodeData{}, []p.Node{
			{"variable_reference", p.NodeData{VarName: "cond"}, []p.Node{}},
		}},
		{"scope", p.NodeData{}, []p.Node{
			{"wait_statement", p.NodeData{Delay: "0.05"}, []p.Node{}},
		}},
	}}

	target := "if (cond)\n{\n" + indentLines("wait 0.05;") + "\n}"

	result := g.Generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_ElseClause(t *testing.T) {
	input := p.Node{"else_clause", p.NodeData{}, []p.Node{
		{"scope", p.NodeData{}, []p.Node{
			{"wait_statement", p.NodeData{Delay: "0.05"}, []p.Node{}},
		}},
	}}

	target := "else\n{\n" + indentLines("wait 0.05;") + "\n}"

	result := g.Generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_WhileLoop(t *testing.T) {
	input := p.Node{"while_loop", p.NodeData{}, []p.Node{
		{"condition", p.NodeData{}, []p.Node{
			{"boolean", p.NodeData{Content: "true"}, []p.Node{}},
		}},
		{"scope", p.NodeData{}, []p.Node{
			{"wait_statement", p.NodeData{Delay: "0.05"}, []p.Node{}},
		}},
	}}

	target := "while (true)\n{\n" + indentLines("wait 0.05;") + "\n}"

	result := g.Generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_ForeachVars(t *testing.T) {
	input := p.Node{"foreach_vars", p.NodeData{}, []p.Node{
		{"variable_reference", p.NodeData{VarName: "item"}, []p.Node{}},
	}}

	target := "item"

	result := g.Generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_ForeachIter(t *testing.T) {
	input := p.Node{"foreach_iter", p.NodeData{}, []p.Node{
		{"variable_reference", p.NodeData{VarName: "items"}, []p.Node{}},
	}}

	target := "items"

	result := g.Generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_ForeachLoop(t *testing.T) {
	input := p.Node{"foreach_loop", p.NodeData{}, []p.Node{
		{"foreach_vars", p.NodeData{}, []p.Node{
			{"variable_reference", p.NodeData{VarName: "item"}, []p.Node{}},
		}},
		{"foreach_iter", p.NodeData{}, []p.Node{
			{"variable_reference", p.NodeData{VarName: "items"}, []p.Node{}},
		}},
		{"scope", p.NodeData{}, []p.Node{
			{"wait_statement", p.NodeData{Delay: "0.05"}, []p.Node{}},
		}},
	}}

	target := "foreach (item in items)\n{\n" + indentLines("wait 0.05;") + "\n}"

	result := g.Generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_SwitchExpr(t *testing.T) {
	input := p.Node{"switch_expr", p.NodeData{}, []p.Node{
		{"variable_reference", p.NodeData{VarName: "x"}, []p.Node{}},
	}}

	target := "x"

	result := g.Generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_CaseClause(t *testing.T) {
	input := p.Node{"case_clause", p.NodeData{}, []p.Node{
		{"string", p.NodeData{Content: "a"}, []p.Node{}},
	}}

	target := "case \"a\":"

	result := g.Generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_DefaultClause(t *testing.T) {
	input := p.Node{"default_clause", p.NodeData{}, []p.Node{}}

	target := "default:"

	result := g.Generate(input)
	assert.Equal(t, target, result)
}

func Test_Generate_SwitchStatement(t *testing.T) {
	input := p.Node{"switch_statement", p.NodeData{}, []p.Node{
		{"switch_expr", p.NodeData{}, []p.Node{
			{"variable_reference", p.NodeData{VarName: "x"}, []p.Node{}},
		}},
		{"scope", p.NodeData{}, []p.Node{
			{"case_clause", p.NodeData{}, []p.Node{
				{"string", p.NodeData{Content: "a"}, []p.Node{}},
			}},
			{"wait_statement", p.NodeData{Delay: "0.05"}, []p.Node{}},
			{"break_statement", p.NodeData{}, []p.Node{}},
			{"default_clause", p.NodeData{}, []p.Node{}},
		}},
	}}

	inner := indentLines("case \"a\":", g.Indent+"wait 0.05;", g.Indent+"break;", "default:")
	target := "switch(x) {\n" + inner + "\n}"

	result := g.Generate(input)
	assert.Equal(t, target, result)
}
