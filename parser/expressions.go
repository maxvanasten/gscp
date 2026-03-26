package parser

import (
	d "github.com/maxvanasten/gscp/diagnostics"
	l "github.com/maxvanasten/gscp/lexer"
)

func parseOperatorToken(tokens []l.Token, index int, output []Node) ([]Node, []d.Diagnostic, int, bool) {
	if index < 0 || index >= len(tokens) {
		return output, nil, index, false
	}
	if tokens[index].Type != l.OPERATOR {
		return output, nil, index, false
	}
	updated := output
	diagnostics := []d.Diagnostic{}
	newIndex := index

	// Unary operator handling
	if tokens[index].Content == "!" || tokens[index].Content == "!!" || tokens[index].Content == "-" || tokens[index].Content == "&" || tokens[index].Content == "~" || tokens[index].Content == "%" {
		isUnary := index == 0
		if index > 0 {
			prev := tokens[index-1].Type
			if prev == l.OPERATOR || prev == l.ASSIGNMENT || prev == l.OPEN_PAREN || prev == l.OPEN_BRACKET || prev == l.COMMA || prev == l.NEWLINE || prev == l.TERMINATOR {
				isUnary = true
			}
		}
		if isUnary {
			operand_tokens := []l.Token{}
			depthParen := 0
			depthBracket := 0
			for i := index + 1; i < len(tokens); i++ {
				curr := tokens[i]
				if depthParen == 0 && depthBracket == 0 {
					if curr.Type == l.OPERATOR || curr.Type == l.TERMINATOR || curr.Type == l.COMMA || curr.Type == l.NEWLINE || curr.Type == l.CLOSE_PAREN || curr.Type == l.CLOSE_BRACKET || curr.Type == l.CLOSE_CURLY {
						break
					}
				}
				switch curr.Type {
				case l.OPEN_PAREN:
					depthParen++
				case l.CLOSE_PAREN:
					if depthParen > 0 {
						depthParen--
					}
				case l.OPEN_BRACKET:
					depthBracket++
				case l.CLOSE_BRACKET:
					if depthBracket > 0 {
						depthBracket--
					}
				}
				operand_tokens = append(operand_tokens, curr)
			}
			if len(operand_tokens) > 0 {
				operand_children, diags := Parse(operand_tokens)
				if len(operand_children) > 0 {
					operand := operand_children[0]
					span := mergeSpan(spanFromToken(tokens[index]), spanFromNode(operand))
					updated = append(updated, nodeWithSpan("unary_expression", NodeData{Operator: tokens[index].Content}, []Node{operand}, span))
					diagnostics = append(diagnostics, diags...)
					newIndex += len(operand_tokens)
					return updated, diagnostics, newIndex, true
				}
				diagnostics = append(diagnostics, diags...)
			}
			diagnostics = append(diagnostics, diagnosticAtIndex("missing unary operand", tokens, index, "error"))
			return updated, diagnostics, newIndex, true
		}
	}

	if tokens[index].Content == "++" || tokens[index].Content == "--" {
		if len(updated) > 0 {
			previous_node := updated[len(updated)-1]
			if previous_node.Type == "variable_reference" {
				updated = updated[:len(updated)-1]
				operator := "+"
				if tokens[index].Content == "--" {
					operator = "-"
				}
				lhs := nodeWithSpan("lhs", NodeData{}, []Node{previous_node}, spanFromNode(previous_node))
				numberSpan := spanFromToken(tokens[index])
				rhsValue := nodeWithSpan("number", NodeData{Content: "1"}, []Node{}, numberSpan)
				rhs := nodeWithSpan("rhs", NodeData{}, []Node{rhsValue}, spanFromNode(rhsValue))
				exprSpan := mergeSpan(spanFromNode(previous_node), numberSpan)
				expr := nodeWithSpan("expression", NodeData{Operator: operator}, []Node{lhs, rhs}, exprSpan)
				assignment_data := NodeData{VarName: previous_node.Data.VarName, Index: previous_node.Data.Index}
				assignmentSpan := mergeSpan(spanFromNode(previous_node), numberSpan)
				updated = append(updated, nodeWithSpan("assignment", assignment_data, []Node{expr}, assignmentSpan))
				return updated, diagnostics, newIndex, true
			}
		}
		operand_tokens := []l.Token{}
		depthParen := 0
		depthBracket := 0
		for i := index + 1; i < len(tokens); i++ {
			curr := tokens[i]
			if depthParen == 0 && depthBracket == 0 {
				if curr.Type == l.OPERATOR || curr.Type == l.TERMINATOR || curr.Type == l.COMMA || curr.Type == l.NEWLINE || curr.Type == l.CLOSE_PAREN || curr.Type == l.CLOSE_BRACKET || curr.Type == l.CLOSE_CURLY {
					break
				}
			}
			switch curr.Type {
			case l.OPEN_PAREN:
				depthParen++
			case l.CLOSE_PAREN:
				if depthParen > 0 {
					depthParen--
				}
			case l.OPEN_BRACKET:
				depthBracket++
			case l.CLOSE_BRACKET:
				if depthBracket > 0 {
					depthBracket--
				}
			}
			operand_tokens = append(operand_tokens, curr)
		}
		if len(operand_tokens) > 0 {
			operand_children, diags := Parse(operand_tokens)
			if len(operand_children) > 0 {
				operand := operand_children[0]
				if operand.Type == "variable_reference" {
					operator := "+"
					if tokens[index].Content == "--" {
						operator = "-"
					}
					lhs := nodeWithSpan("lhs", NodeData{}, []Node{operand}, spanFromNode(operand))
					numberSpan := spanFromToken(tokens[index])
					rhsValue := nodeWithSpan("number", NodeData{Content: "1"}, []Node{}, numberSpan)
					rhs := nodeWithSpan("rhs", NodeData{}, []Node{rhsValue}, spanFromNode(rhsValue))
					exprSpan := mergeSpan(spanFromToken(tokens[index]), spanFromNode(operand))
					expr := nodeWithSpan("expression", NodeData{Operator: operator}, []Node{lhs, rhs}, exprSpan)
					assignment_data := NodeData{VarName: operand.Data.VarName, Index: operand.Data.Index}
					assignmentSpan := mergeSpan(spanFromToken(tokens[index]), spanFromNode(operand))
					updated = append(updated, nodeWithSpan("assignment", assignment_data, []Node{expr}, assignmentSpan))
					diagnostics = append(diagnostics, diags...)
					newIndex += len(operand_tokens)
					return updated, diagnostics, newIndex, true
				}
			}
			diagnostics = append(diagnostics, diags...)
		}
		return updated, diagnostics, newIndex, true
	}

	if tokens[index].Content == "?" {
		if len(updated) == 0 {
			diagnostics = append(diagnostics, diagnosticAtIndex("operator missing left-hand operand", tokens, index, "error"))
			return updated, diagnostics, newIndex, true
		}
		condition := updated[len(updated)-1]
		updated = updated[:len(updated)-1]
		expr_tokens := l.TokensUntilAny(tokens[index+1:], []l.TokenType{l.NEWLINE, l.TERMINATOR})
		splitIndex := topLevelIndex(expr_tokens, l.COLON)
		true_tokens := []l.Token{}
		false_tokens := []l.Token{}
		if splitIndex >= 0 {
			true_tokens = expr_tokens[:splitIndex]
			false_tokens = expr_tokens[splitIndex+1:]
			false_tokens = trimTrailingAny(false_tokens, l.TERMINATOR, l.NEWLINE)
		} else {
			true_tokens = expr_tokens
		}
		true_tokens = trimTrailingAny(true_tokens, l.TERMINATOR, l.NEWLINE)
		true_children, true_diags := Parse(true_tokens)
		false_children, false_diags := Parse(false_tokens)
		conditionNode := nodeWithSpan("condition", NodeData{}, []Node{condition}, spanFromNode(condition))
		trueSpan := spanFromNodes(true_children)
		falseSpan := spanFromNodes(false_children)
		trueNode := nodeWithSpan("true_expr", NodeData{}, true_children, trueSpan)
		falseNode := nodeWithSpan("false_expr", NodeData{}, false_children, falseSpan)
		endTokens := trimTrailingAny(expr_tokens, l.TERMINATOR, l.NEWLINE)
		ternarySpan := mergeSpan(spanFromNode(condition), spanFromTokens(endTokens))
		ternary := nodeWithSpan("ternary_expression", NodeData{}, []Node{conditionNode, trueNode, falseNode}, ternarySpan)
		updated = append(updated, ternary)
		diagnostics = append(diagnostics, true_diags...)
		diagnostics = append(diagnostics, false_diags...)
		newIndex += len(expr_tokens)
		return updated, diagnostics, newIndex, true
	}

	if index <= 0 {
		diagnostics = append(diagnostics, diagnosticAtIndex("operator missing left-hand operand", tokens, index, "error"))
		return updated, diagnostics, newIndex, true
	}
	if len(updated) == 0 {
		diagnostics = append(diagnostics, diagnosticAtIndex("operator missing left-hand operand", tokens, index, "error"))
		return updated, diagnostics, newIndex, true
	}
	// Check if previous node is either a string, variable_reference, number, unary/function_call or expression
	previous_node := updated[len(updated)-1]
	switch previous_node.Type {
	case "string", "variable_reference", "number", "unary_expression", "function_call", "expression":
		// Set LHS to previous node
		lhs := nodeWithSpan("lhs", NodeData{}, []Node{previous_node}, spanFromNode(previous_node))
		// Delete previous node
		updated = updated[:len(updated)-1]
		// Get all tokens from OPERATOR until TERMINATOR or CLOSE_PAREN at depth 0
		// Using TokensForExpression to handle multiline expressions
		expr_tokens := l.TokensForExpression(tokens[index+1:], []l.TokenType{l.TERMINATOR, l.CLOSE_PAREN})
		if tokens[index].Content == "&&" {
			depthParen := 0
			depthBracket := 0
			cutIndex := -1
			for i, tok := range expr_tokens {
				switch tok.Type {
				case l.OPEN_PAREN:
					depthParen++
				case l.CLOSE_PAREN:
					if depthParen > 0 {
						depthParen--
					}
				case l.OPEN_BRACKET:
					depthBracket++
				case l.CLOSE_BRACKET:
					if depthBracket > 0 {
						depthBracket--
					}
				}
				if depthParen == 0 && depthBracket == 0 && tok.Type == l.OPERATOR && tok.Content == "||" {
					cutIndex = i
					break
				}
			}
			if cutIndex >= 0 {
				expr_tokens = expr_tokens[:cutIndex]
			}
		}
		// Parse those tokens into RHS
		rhs_children, diags := Parse(expr_tokens)
		rhsSpan := spanFromNodes(rhs_children)
		rhs := nodeWithSpan("rhs", NodeData{}, rhs_children, rhsSpan)
		if len(rhs_children) == 0 {
			diagnostics = append(diagnostics, diagnosticAtIndex("operator missing right-hand operand", tokens, index, "error"))
		}
		// Add Expression node to output
		exprSpan := mergeSpan(spanFromNode(previous_node), rhsSpan)
		if !exprSpan.valid {
			exprSpan = mergeSpan(spanFromNode(previous_node), spanFromToken(tokens[index]))
		}
		updated = append(updated, nodeWithSpan("expression", NodeData{Operator: tokens[index].Content}, []Node{lhs, rhs}, exprSpan))
		diagnostics = append(diagnostics, diags...)
		newIndex += len(expr_tokens)
	default:
		updated = append(updated, nodeWithSpan("operator", NodeData{Content: tokens[index].Content}, []Node{}, spanFromToken(tokens[index])))
	}
	return updated, diagnostics, newIndex, true
}
