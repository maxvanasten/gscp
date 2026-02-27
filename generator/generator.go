package generator

import (
	"strings"

	p "github.com/maxvanasten/gscp/parser"
)

var Indent = "  "

func joinChildren(children []p.Node, sep string) string {
	parts := []string{}
	for _, child := range children {
		parts = append(parts, Generate(child))
	}
	return strings.Join(parts, sep)
}

func stripTrailingSemicolon(value string) string {
	return strings.TrimSuffix(value, ";")
}

func Generate(node p.Node) string {
	output := strings.Builder{}

	switch node.Type {
	case "variable_reference":
		output.WriteString(node.Data.VarName)
		if node.Data.Index != "" {
			output.WriteString("[")
			output.WriteString(node.Data.Index)
			output.WriteString("]")
		}
	case "string":
		output.WriteRune('"')
		output.WriteString(node.Data.Content)
		output.WriteRune('"')
	case "boolean":
		output.WriteString(node.Data.Content)
	case "number":
		output.WriteString(node.Data.Content)
	case "expression":
		lhs := Generate(node.Children[0])
		rhs := Generate(node.Children[1])
		output.WriteString(lhs)
		output.WriteString(" ")
		output.WriteString(node.Data.Operator)
		output.WriteString(" ")
		output.WriteString(rhs)
	case "lhs", "rhs":
		if len(node.Children) > 0 {
			output.WriteString(Generate(node.Children[0]))
		}
	case "unary_expression":
		output.WriteString(node.Data.Operator)
		if len(node.Children) > 0 {
			output.WriteString(Generate(node.Children[0]))
		}
	case "assignment":
		output.WriteString(node.Data.VarName)
		if node.Data.Index != "" {
			output.WriteString("[")
			output.WriteString(node.Data.Index)
			output.WriteString("]")
		}
		output.WriteString(" = ")
		if len(node.Children) == 1 {
			output.WriteString(Generate(node.Children[0]))
		} else if len(node.Children) > 1 {
			output.WriteString(joinChildren(node.Children, ", "))
		}
		output.WriteRune(';')
	case "array_literal":
		output.WriteString("[")
		output.WriteString(joinChildren(node.Children, ", "))
		output.WriteString("]")
	case "vector_literal":
		output.WriteString("(")
		output.WriteString(joinChildren(node.Children, ", "))
		output.WriteString(")")
	case "include_statement":
		output.WriteString("#include ")
		output.WriteString(node.Data.Path)
		output.WriteRune(';')
	case "wait_statement":
		output.WriteString("wait ")
		output.WriteString(node.Data.Delay)
		output.WriteRune(';')
	case "break_statement":
		output.WriteString("break;")
	case "return_statement":
		output.WriteString("return")
		if len(node.Children) > 0 {
			output.WriteString(" ")
			output.WriteString(Generate(node.Children[0]))
		}
		output.WriteRune(';')
	case "function_call":
		if node.Data.Method != "" {
			output.WriteString(node.Data.Method)
			output.WriteString(" ")
		}
		if node.Data.Thread {
			output.WriteString("thread ")
		}
		if node.Data.Path != "" {
			output.WriteString(node.Data.Path)
			output.WriteString("::")
		}
		output.WriteString(node.Data.FunctionName)
		output.WriteString("(")
		output.WriteString(joinChildren(node.Children, ", "))
		output.WriteString(");")
	case "function_declaration":
		output.WriteString(node.Data.FunctionName)
		output.WriteString("(")
		if len(node.Children) > 0 {
			output.WriteString(Generate(node.Children[0]))
		}
		output.WriteString(")\n{")
		if len(node.Children) > 1 {
			scopeBody := Generate(node.Children[1])
			if scopeBody != "" {
				output.WriteString("\n")
				output.WriteString(scopeBody)
			}
		}
		output.WriteString("\n}")
	case "args":
		output.WriteString(joinChildren(node.Children, ", "))
	case "scope":
		lines := []string{}
		for _, child := range node.Children {
			lines = append(lines, Indent+Generate(child))
		}
		output.WriteString(strings.Join(lines, "\n"))
	case "for_init", "for_condition", "for_post":
		if len(node.Children) > 0 {
			output.WriteString(stripTrailingSemicolon(Generate(node.Children[0])))
		}
	case "for_loop":
		init := ""
		cond := ""
		post := ""
		if len(node.Children) > 0 {
			init = Generate(node.Children[0])
		}
		if len(node.Children) > 1 {
			cond = Generate(node.Children[1])
		}
		if len(node.Children) > 2 {
			post = Generate(node.Children[2])
		}
		if init == "" && cond == "" && post == "" {
			output.WriteString("for ( ;; )\n{")
		} else {
			output.WriteString("for (")
			if init == "" {
				output.WriteString(" ")
			} else {
				output.WriteString(init)
			}
			output.WriteString("; ")
			if cond != "" {
				output.WriteString(cond)
			}
			output.WriteString("; ")
			if post != "" {
				output.WriteString(post)
			}
			output.WriteString(")\n{")
		}
		if len(node.Children) > 3 {
			scopeBody := Generate(node.Children[3])
			if scopeBody != "" {
				output.WriteString("\n")
				output.WriteString(scopeBody)
			}
		}
		output.WriteString("\n}")
	case "condition":
		if len(node.Children) > 0 {
			output.WriteString(Generate(node.Children[0]))
		}
	case "if_statement":
		cond := ""
		if len(node.Children) > 0 {
			cond = Generate(node.Children[0])
		}
		output.WriteString("if (")
		output.WriteString(cond)
		output.WriteString(")\n{")
		if len(node.Children) > 1 {
			scopeBody := Generate(node.Children[1])
			if scopeBody != "" {
				output.WriteString("\n")
				output.WriteString(scopeBody)
			}
		}
		output.WriteString("\n}")
	case "else_clause":
		output.WriteString("else\n{")
		if len(node.Children) > 0 {
			scopeBody := Generate(node.Children[0])
			if scopeBody != "" {
				output.WriteString("\n")
				output.WriteString(scopeBody)
			}
		}
		output.WriteString("\n}")
	case "while_loop":
		cond := ""
		if len(node.Children) > 0 {
			cond = Generate(node.Children[0])
		}
		output.WriteString("while (")
		output.WriteString(cond)
		output.WriteString(")\n{")
		if len(node.Children) > 1 {
			scopeBody := Generate(node.Children[1])
			if scopeBody != "" {
				output.WriteString("\n")
				output.WriteString(scopeBody)
			}
		}
		output.WriteString("\n}")
	case "foreach_vars", "foreach_iter":
		output.WriteString(joinChildren(node.Children, ", "))
	case "foreach_loop":
		vars := ""
		iter := ""
		if len(node.Children) > 0 {
			vars = Generate(node.Children[0])
		}
		if len(node.Children) > 1 {
			iter = Generate(node.Children[1])
		}
		output.WriteString("foreach (")
		output.WriteString(vars)
		output.WriteString(" in ")
		output.WriteString(iter)
		output.WriteString(")\n{")
		if len(node.Children) > 2 {
			scopeBody := Generate(node.Children[2])
			if scopeBody != "" {
				output.WriteString("\n")
				output.WriteString(scopeBody)
			}
		}
		output.WriteString("\n}")
	case "switch_expr":
		if len(node.Children) > 0 {
			output.WriteString(Generate(node.Children[0]))
		}
	case "case_clause":
		output.WriteString("case ")
		if len(node.Children) > 0 {
			output.WriteString(Generate(node.Children[0]))
		}
		output.WriteString(":")
	case "default_clause":
		output.WriteString("default:")
	case "switch_statement":
		switchExpr := ""
		if len(node.Children) > 0 {
			switchExpr = Generate(node.Children[0])
		}
		output.WriteString("switch(")
		output.WriteString(switchExpr)
		output.WriteString(") {")
		if len(node.Children) > 1 {
			scopeNode := node.Children[1]
			lines := []string{}
			inCase := false
			for _, child := range scopeNode.Children {
				line := Generate(child)
				if child.Type == "case_clause" || child.Type == "default_clause" {
					inCase = true
					lines = append(lines, Indent+line)
					continue
				}
				if inCase {
					lines = append(lines, Indent+Indent+line)
				} else {
					lines = append(lines, Indent+line)
				}
			}
			if len(lines) > 0 {
				output.WriteString("\n")
				output.WriteString(strings.Join(lines, "\n"))
			}
		}
		output.WriteString("\n}")
	}

	return output.String()
}
