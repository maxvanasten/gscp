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

func joinInlineChildren(children []p.Node, sep string) string {
	parts := []string{}
	for i := 0; i < len(children); i++ {
		child := children[i]
		parts = append(parts, stripTrailingSemicolon(Generate(child)))
	}
	return strings.Join(parts, sep)
}

func stripTrailingSemicolon(value string) string {
	return strings.TrimSuffix(value, ";")
}

func indentMultiline(value string, prefix string) string {
	lines := strings.Split(value, "\n")
	for i, line := range lines {
		if line == "" {
			lines[i] = prefix
			continue
		}
		lines[i] = prefix + line
	}
	return strings.Join(lines, "\n")
}

func formatBlock(header string, scope p.Node) string {
	output := strings.Builder{}
	output.WriteString(header)
	output.WriteString("\n{")
	if scopeBody := Generate(scope); scopeBody != "" {
		output.WriteString("\n")
		output.WriteString(scopeBody)
	}
	output.WriteString("\n}")
	return output.String()
}

func Generate(node p.Node) string {
	output := strings.Builder{}

	switch node.Type {
	case "comment":
		output.WriteString(node.Data.Content)
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
	case "hash_string":
		output.WriteString("#\"")
		output.WriteString(node.Data.Content)
		output.WriteString("\"")
	case "boolean":
		output.WriteString(node.Data.Content)
	case "number":
		output.WriteString(node.Data.Content)
	case "expression":
		lhs := stripTrailingSemicolon(Generate(node.Children[0]))
		rhs := stripTrailingSemicolon(Generate(node.Children[1]))
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
			output.WriteString(stripTrailingSemicolon(Generate(node.Children[0])))
		} else if len(node.Children) > 1 {
			output.WriteString(joinInlineChildren(node.Children, ", "))
		}
		output.WriteRune(';')
	case "array_literal":
		output.WriteString("[")
		output.WriteString(joinInlineChildren(node.Children, ", "))
		output.WriteString("]")
	case "vector_literal":
		output.WriteString("(")
		output.WriteString(joinInlineChildren(node.Children, ", "))
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
		output.WriteString(joinInlineChildren(node.Children, ", "))
		output.WriteString(");")
	case "function_declaration":
		header := strings.Builder{}
		header.WriteString(node.Data.FunctionName)
		header.WriteString("(")
		if len(node.Children) > 0 {
			header.WriteString(joinInlineChildren(node.Children[0].Children, ", "))
		}
		header.WriteString(")")
		if len(node.Children) > 1 {
			output.WriteString(formatBlock(header.String(), node.Children[1]))
		} else {
			output.WriteString(header.String())
			output.WriteString("\n{\n}")
		}
	case "args":
		output.WriteString(joinInlineChildren(node.Children, ", "))
	case "scope":
		lines := []string{}
		for _, child := range node.Children {
			line := Generate(child)
			if child.Type == "function_call" && !strings.HasSuffix(line, ";") {
				line += ";"
			}
			lines = append(lines, indentMultiline(line, Indent))
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
		header := strings.Builder{}
		if init == "" && cond == "" && post == "" {
			header.WriteString("for ( ;; )")
		} else {
			header.WriteString("for (")
			if init == "" {
				header.WriteString(" ")
			} else {
				header.WriteString(init)
			}
			header.WriteString("; ")
			if cond != "" {
				header.WriteString(cond)
			}
			header.WriteString("; ")
			if post != "" {
				header.WriteString(post)
			}
			header.WriteString(")")
		}
		if len(node.Children) > 3 {
			output.WriteString(formatBlock(header.String(), node.Children[3]))
		} else {
			output.WriteString(header.String())
			output.WriteString("\n{\n}")
		}
	case "condition":
		if len(node.Children) > 0 {
			output.WriteString(stripTrailingSemicolon(Generate(node.Children[0])))
		}
	case "if_statement":
		cond := ""
		if len(node.Children) > 0 {
			cond = Generate(node.Children[0])
		}
		header := "if (" + cond + ")"
		if len(node.Children) > 1 {
			output.WriteString(formatBlock(header, node.Children[1]))
		} else {
			output.WriteString(header)
			output.WriteString("\n{\n}")
		}
	case "else_clause":
		if len(node.Children) > 0 {
			output.WriteString(formatBlock("else", node.Children[0]))
		} else {
			output.WriteString("else\n{\n}")
		}
	case "while_loop":
		cond := ""
		if len(node.Children) > 0 {
			cond = Generate(node.Children[0])
		}
		header := "while (" + cond + ")"
		if len(node.Children) > 1 {
			output.WriteString(formatBlock(header, node.Children[1]))
		} else {
			output.WriteString(header)
			output.WriteString("\n{\n}")
		}
	case "foreach_vars", "foreach_iter":
		output.WriteString(joinInlineChildren(node.Children, ", "))
	case "foreach_loop":
		vars := ""
		iter := ""
		if len(node.Children) > 0 {
			vars = Generate(node.Children[0])
		}
		if len(node.Children) > 1 {
			iter = Generate(node.Children[1])
		}
		header := "foreach (" + vars + " in " + iter + ")"
		if len(node.Children) > 2 {
			output.WriteString(formatBlock(header, node.Children[2]))
		} else {
			output.WriteString(header)
			output.WriteString("\n{\n}")
		}
	case "switch_expr":
		if len(node.Children) > 0 {
			output.WriteString(stripTrailingSemicolon(Generate(node.Children[0])))
		}
	case "case_clause":
		output.WriteString("case ")
		if len(node.Children) > 0 {
			output.WriteString(stripTrailingSemicolon(Generate(node.Children[0])))
		}
		output.WriteString(":")
	case "default_clause":
		output.WriteString("default:")
	case "switch_statement":
		switchExpr := ""
		if len(node.Children) > 0 {
			switchExpr = stripTrailingSemicolon(Generate(node.Children[0]))
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
				if child.Type == "function_call" && !strings.HasSuffix(line, ";") {
					line += ";"
				}
				if child.Type == "case_clause" || child.Type == "default_clause" {
					inCase = true
					lines = append(lines, indentMultiline(line, Indent))
					continue
				}
				if inCase {
					lines = append(lines, indentMultiline(line, Indent+Indent))
				} else {
					lines = append(lines, indentMultiline(line, Indent))
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
