package parser

import (
	"strings"

	d "github.com/maxvanasten/gscp/diagnostics"
	l "github.com/maxvanasten/gscp/lexer"
)

func variableReferenceWithIndex(node Node) string {
	name := node.Data.VarName
	if node.Data.Index != "" {
		name += "[" + node.Data.Index + "]"
	}
	return name
}

func Parse(tokens []l.Token) ([]Node, []d.Diagnostic) {
	output := []Node{}
	diagnostics := []d.Diagnostic{}

	index := 0
	for index < len(tokens) {
		switch tokens[index].Type {
		case l.LINE_COMMENT, l.BLOCK_COMMENT:
			output = append(output, nodeWithSpan("comment", NodeData{Content: tokens[index].Content}, []Node{}, spanFromToken(tokens[index])))
		case l.STRING:
			output = append(output, nodeWithSpan("string", NodeData{Content: tokens[index].Content}, []Node{}, spanFromToken(tokens[index])))
		case l.SYMBOL:
			normalized := strings.TrimSuffix(tokens[index].Content, ":")
			switch normalized {
			case "#include":
				if index+1 < len(tokens) && tokens[index+1].Type == l.SYMBOL {
					span := mergeSpan(spanFromToken(tokens[index]), spanFromToken(tokens[index+1]))
					output = append(output, nodeWithSpan("include_statement", NodeData{Path: tokens[index+1].Content}, []Node{}, span))
					index++
				} else {
					diagnostics = append(diagnostics, diagnosticAtIndex("missing include path", tokens, index, "error"))
				}
			case "wait":
				if index+1 < len(tokens) {
					next := tokens[index+1]
					if next.Type == l.OPEN_PAREN {
						output = append(output, nodeWithSpan("variable_reference", NodeData{VarName: tokens[index].Content}, []Node{}, spanFromToken(tokens[index])))
						break
					}
					if next.Type == l.NUMBER || next.Type == l.SYMBOL {
						span := mergeSpan(spanFromToken(tokens[index]), spanFromToken(next))
						output = append(output, nodeWithSpan("wait_statement", NodeData{Delay: next.Content}, []Node{}, span))
						index++
						break
					}
				}
				diagnostics = append(diagnostics, diagnosticAtIndex("missing wait duration", tokens, index, "error"))
			case "thread":
				output = append(output, nodeWithSpan("thread_keyword", NodeData{}, []Node{}, spanFromToken(tokens[index])))
			case "true", "false":
				output = append(output, nodeWithSpan("boolean", NodeData{Content: normalized}, []Node{}, spanFromToken(tokens[index])))
			case "break":
				output = append(output, nodeWithSpan("break_statement", NodeData{}, []Node{}, spanFromToken(tokens[index])))
			case "return":
				ret_tokens := l.TokensUntilAny(tokens[index+1:], []l.TokenType{l.NEWLINE, l.TERMINATOR})
				ret_children, diags := Parse(ret_tokens)
				endToken, ok := lastNonTokenType(ret_tokens, l.NEWLINE, l.TERMINATOR)
				endSpan := spanFromToken(tokens[index])
				if ok {
					endSpan = spanFromToken(endToken)
				}
				span := mergeSpan(spanFromToken(tokens[index]), endSpan)
				output = append(output, nodeWithSpan("return_statement", NodeData{}, ret_children, span))
				diagnostics = append(diagnostics, diags...)
				index += len(ret_tokens)
			case "case":
				case_tokens := l.TokensUntilAny(tokens[index+1:], []l.TokenType{l.NEWLINE, l.TERMINATOR, l.COLON})
				case_children, diags := Parse(case_tokens)
				endToken, ok := lastNonTokenType(case_tokens, l.NEWLINE, l.TERMINATOR, l.COLON)
				endSpan := spanFromToken(tokens[index])
				if ok {
					endSpan = spanFromToken(endToken)
				}
				span := mergeSpan(spanFromToken(tokens[index]), endSpan)
				output = append(output, nodeWithSpan("case_clause", NodeData{}, case_children, span))
				diagnostics = append(diagnostics, diags...)
				index += len(case_tokens)
			case "default":
				output = append(output, nodeWithSpan("default_clause", NodeData{}, []Node{}, spanFromToken(tokens[index])))
			case "else":
				output = append(output, nodeWithSpan("else_header", NodeData{}, []Node{}, spanFromToken(tokens[index])))
			case "do":
				output = append(output, nodeWithSpan("do_header", NodeData{}, []Node{}, spanFromToken(tokens[index])))
			default:
				output = append(output, nodeWithSpan("variable_reference", NodeData{VarName: tokens[index].Content}, []Node{}, spanFromToken(tokens[index])))
			}

		case l.NUMBER:
			output = append(output, nodeWithSpan("number", NodeData{Content: tokens[index].Content}, []Node{}, spanFromToken(tokens[index])))
		case l.OPERATOR:
			updatedOutput, opDiags, newIndex, handled := parseOperatorToken(tokens, index, output)
			if handled {
				output = updatedOutput
				diagnostics = append(diagnostics, opDiags...)
				index = newIndex
				break
			}
		case l.ASSIGNMENT:
			if index <= 0 {
				diagnostics = append(diagnostics, diagnosticAtIndex("assignment missing left-hand side", tokens, index, "error"))
				break
			}
			if len(output) == 0 {
				diagnostics = append(diagnostics, diagnosticAtIndex("assignment missing left-hand side", tokens, index, "error"))
				break
			}
			// Check if previous node is a variable_reference
			previous_node := output[len(output)-1]
			if previous_node.Type != "variable_reference" {
				diagnostics = append(diagnostics, diagnosticAtIndex("assignment target must be a variable", tokens, index, "error"))
				output = append(output, nodeWithSpan("assignment", NodeData{Content: tokens[index].Content}, []Node{}, spanFromToken(tokens[index])))
				break
			}
			output = output[:len(output)-1]
			// Get all tokens from ASSIGNMENT until END, NEWLINE or TERMINATOR (at depth 0 for multiline support)
			ass_tokens := l.TokensUntilAnyBalanced(tokens[index+1:], []l.TokenType{l.NEWLINE, l.TERMINATOR})

			ass_children, diags := Parse(ass_tokens)
			assignment_data := NodeData{VarName: previous_node.Data.VarName, Index: previous_node.Data.Index}
			endToken, ok := lastNonTokenType(ass_tokens, l.TERMINATOR, l.NEWLINE)
			endSpan := spanFromToken(tokens[index])
			if ok {
				endSpan = spanFromToken(endToken)
			}
			assignmentSpan := mergeSpan(spanFromNode(previous_node), endSpan)
			if tokens[index].Content == "+=" || tokens[index].Content == "-=" || tokens[index].Content == "*=" || tokens[index].Content == "/=" {
				operator := tokens[index].Content[:1]
				lhs := nodeWithSpan("lhs", NodeData{}, []Node{previous_node}, spanFromNode(previous_node))
				rhsSpan := spanFromNodes(ass_children)
				rhs := nodeWithSpan("rhs", NodeData{}, ass_children, rhsSpan)
				exprSpan := mergeSpan(spanFromNode(previous_node), rhsSpan)
				if !exprSpan.valid {
					exprSpan = assignmentSpan
				}
				expr := nodeWithSpan("expression", NodeData{Operator: operator}, []Node{lhs, rhs}, exprSpan)
				output = append(output, nodeWithSpan("assignment", assignment_data, []Node{expr}, assignmentSpan))
			} else {
				output = append(output, nodeWithSpan("assignment", assignment_data, ass_children, assignmentSpan))
			}
			diagnostics = append(diagnostics, diags...)
			index += len(ass_tokens)
		case l.OPEN_PAREN:
			arg_tokens, foundClose := tokensUntilMatchingClose(tokens[index+1:], l.OPEN_PAREN, l.CLOSE_PAREN)
			if !foundClose {
				diagnostics = append(diagnostics, diagnosticAtIndex("missing closing )", tokens, index, "error"))
			}
			trimmedArgTokens := trimTrailingToken(arg_tokens, l.CLOSE_PAREN)
			if len(output)-1 >= 0 {
				previous_node := output[len(output)-1]
				if previous_node.Type == "variable_reference" && previous_node.Data.VarName == "for" {
					header_tokens := trimmedArgTokens
					headerSpan := mergeSpan(spanFromNode(previous_node), spanFromTokens(arg_tokens))

					segments := splitTopLevel(header_tokens, l.TERMINATOR, false)
					for len(segments) < 3 {
						segments = append(segments, []l.Token{})
					}

					init_children, init_diags := Parse(segments[0])
					cond_children, cond_diags := Parse(segments[1])
					post_children, post_diags := Parse(segments[2])
					initNode := nodeWithSpan("for_init", NodeData{}, init_children, headerSpan)
					condNode := nodeWithSpan("for_condition", NodeData{}, cond_children, headerSpan)
					postNode := nodeWithSpan("for_post", NodeData{}, post_children, headerSpan)
					for_header := nodeWithSpan("for_header", NodeData{}, []Node{initNode, condNode, postNode}, headerSpan)

					diagnostics = append(diagnostics, init_diags...)
					diagnostics = append(diagnostics, cond_diags...)
					diagnostics = append(diagnostics, post_diags...)

					output = output[:len(output)-1]
					output = append(output, for_header)
					index += len(arg_tokens)
					break
				}
				if previous_node.Type == "variable_reference" && (previous_node.Data.VarName == "if" || previous_node.Data.VarName == "while" || previous_node.Data.VarName == "foreach" || previous_node.Data.VarName == "switch") {
					header_tokens := trimmedArgTokens
					headerSpan := mergeSpan(spanFromNode(previous_node), spanFromTokens(arg_tokens))
					output = output[:len(output)-1]
					switch previous_node.Data.VarName {
					case "if":
						cond_children, diags := Parse(header_tokens)
						condNode := nodeWithSpan("condition", NodeData{}, cond_children, headerSpan)
						output = append(output, nodeWithSpan("if_header", NodeData{}, []Node{condNode}, headerSpan))
						diagnostics = append(diagnostics, diags...)
					case "while":
						cond_children, diags := Parse(header_tokens)
						condNode := nodeWithSpan("condition", NodeData{}, cond_children, headerSpan)
						output = append(output, nodeWithSpan("while_header", NodeData{}, []Node{condNode}, headerSpan))
						diagnostics = append(diagnostics, diags...)
					case "switch":
						expr_children, diags := Parse(header_tokens)
						exprNode := nodeWithSpan("switch_expr", NodeData{}, expr_children, headerSpan)
						output = append(output, nodeWithSpan("switch_header", NodeData{}, []Node{exprNode}, headerSpan))
						diagnostics = append(diagnostics, diags...)
					case "foreach":
						left := []l.Token{}
						right := []l.Token{}
						depthParen := 0
						depthBracket := 0
						foundIn := false
						for _, tok := range header_tokens {
							if !foundIn && depthParen == 0 && depthBracket == 0 && tok.Type == l.SYMBOL && tok.Content == "in" {
								foundIn = true
								continue
							}
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
							if foundIn {
								right = append(right, tok)
							} else {
								left = append(left, tok)
							}
						}
						left_children, left_diags := Parse(left)
						right_children, right_diags := Parse(right)
						leftNode := nodeWithSpan("foreach_vars", NodeData{}, left_children, headerSpan)
						rightNode := nodeWithSpan("foreach_iter", NodeData{}, right_children, headerSpan)
						output = append(output, nodeWithSpan("foreach_header", NodeData{}, []Node{leftNode, rightNode}, headerSpan))
						diagnostics = append(diagnostics, left_diags...)
						diagnostics = append(diagnostics, right_diags...)
					}
					index += len(arg_tokens)
					break
				}
			}

			if len(output)-1 < 0 || output[len(output)-1].Type != "variable_reference" {
				inner_tokens := trimmedArgTokens
				hasComma := hasTopLevelToken(inner_tokens, l.COMMA)
				if hasComma {
					elem_slices := splitTopLevel(inner_tokens, l.COMMA, true)
					vec_children := []Node{}
					for _, elem_tokens := range elem_slices {
						children, diags := Parse(elem_tokens)
						vec_children = append(vec_children, children...)
						diagnostics = append(diagnostics, diags...)
					}
					startSpan := spanFromToken(tokens[index])
					endSpan := startSpan
					if len(arg_tokens) > 0 {
						endSpan = spanFromToken(arg_tokens[len(arg_tokens)-1])
					}
					output = append(output, nodeWithSpan("vector_literal", NodeData{}, vec_children, mergeSpan(startSpan, endSpan)))
					index += len(arg_tokens)
					break
				}
				group_children, diags := Parse(inner_tokens)
				output = append(output, group_children...)
				diagnostics = append(diagnostics, diags...)
				index += len(arg_tokens)
				break
			}

			arg_token_slices := [][]l.Token{}
			// Split arg_tokens into slices at top-level commas only
			arg_token_slices = splitTopLevel(trimmedArgTokens, l.COMMA, false)

			arg_children := []Node{}
			argDiagnostics := []d.Diagnostic{}
			for _, argument_tokens := range arg_token_slices {
				children, diags := Parse(argument_tokens)
				arg_children = append(arg_children, children...)
				argDiagnostics = append(argDiagnostics, diags...)
			}

			data := NodeData{}
			if len(output)-1 >= 0 {
				c := 0
				callStartSpan := nodeSpan{}
				if output[len(output)-1].Type == "variable_reference" {
					c = 1
					name := output[len(output)-1].Data.VarName
					if strings.Contains(name, "::") {
						parts := strings.Split(name, "::")
						if len(parts) > 1 {
							data.Path = strings.Join(parts[:len(parts)-1], "::")
							name = parts[len(parts)-1]
						}
					} else if strings.Contains(name, ".") {
						parts := strings.Split(name, ".")
						if len(parts) > 1 {
							data.Method = strings.Join(parts[:len(parts)-1], ".")
							name = parts[len(parts)-1]
						}
					}
					data.FunctionName = name
					if len(output)-2 >= 0 {
						if output[len(output)-2].Type == "thread_keyword" {
							c = 2
							data.Thread = true
							if len(output)-3 >= 0 {
								if output[len(output)-3].Type == "variable_reference" {
									c = 3
									data.Method = variableReferenceWithIndex(output[len(output)-3])
								}
							}
						}
						if output[len(output)-2].Type == "variable_reference" {
							c = 2
							data.Method = variableReferenceWithIndex(output[len(output)-2])
						}
					}
				}
				if c > 0 {
					callStartSpan = spanFromNodes(output[len(output)-c:])
				}
				output = output[:len(output)-c]
				callEndSpan := spanFromToken(tokens[index])
				if len(arg_tokens) > 0 {
					callEndSpan = spanFromToken(arg_tokens[len(arg_tokens)-1])
				}
				callSpan := mergeSpan(callStartSpan, callEndSpan)
				output = append(output, nodeWithSpan("function_call", data, arg_children, callSpan))
			}
			diagnostics = append(diagnostics, argDiagnostics...)

			index += len(arg_tokens)
		case l.OPEN_BRACKET:
			bracket_tokens, foundClose := tokensUntilMatchingClose(tokens[index+1:], l.OPEN_BRACKET, l.CLOSE_BRACKET)
			if !foundClose {
				diagnostics = append(diagnostics, diagnosticAtIndex("missing closing ]", tokens, index, "error"))
			}
			if len(output)-1 >= 0 {
				previous_node := output[len(output)-1]
				if previous_node.Type == "variable_reference" {
					index_tokens := bracket_tokens
					if len(index_tokens) > 0 && index_tokens[len(index_tokens)-1].Type == l.CLOSE_BRACKET {
						index_tokens = index_tokens[:len(index_tokens)-1]
					}
					index_content := tokensToString(index_tokens)
					startSpan := spanFromNode(previous_node)
					endSpan := startSpan
					if len(bracket_tokens) > 0 {
						endSpan = spanFromToken(bracket_tokens[len(bracket_tokens)-1])
					}
					indexed_node := nodeWithSpan("variable_reference", NodeData{VarName: previous_node.Data.VarName, Index: index_content}, []Node{}, mergeSpan(startSpan, endSpan))
					output = output[:len(output)-1]
					output = append(output, indexed_node)
					index += len(bracket_tokens)
					break
				}
			}

			array_tokens := bracket_tokens
			if len(array_tokens) > 0 && array_tokens[len(array_tokens)-1].Type == l.CLOSE_BRACKET {
				array_tokens = array_tokens[:len(array_tokens)-1]
			}
			// Filter out newlines from array tokens before parsing elements
			filtered_tokens := []l.Token{}
			for _, t := range array_tokens {
				if t.Type != l.NEWLINE {
					filtered_tokens = append(filtered_tokens, t)
				}
			}
			elem_slices := splitTopLevel(filtered_tokens, l.COMMA, true)

			array_children := []Node{}
			for _, elem_tokens := range elem_slices {
				children, diags := Parse(elem_tokens)
				array_children = append(array_children, children...)
				diagnostics = append(diagnostics, diags...)
			}
			startSpan := spanFromToken(tokens[index])
			endSpan := startSpan
			if len(bracket_tokens) > 0 {
				endSpan = spanFromToken(bracket_tokens[len(bracket_tokens)-1])
			}
			output = append(output, nodeWithSpan("array_literal", NodeData{}, array_children, mergeSpan(startSpan, endSpan)))
			index += len(bracket_tokens)
		case l.OPEN_CURLY:
			if len(output)-1 < 0 {
				diagnostics = append(diagnostics, diagnosticAtIndex("unexpected {", tokens, index, "error"))
				break
			}
			previous_node := output[len(output)-1]
			switch previous_node.Type {
			case "function_call":
				output = output[:len(output)-1]
				scope_node, scope_diags, rawLen := parseScope(tokens, index)
				arg_node := nodeWithSpan("args", NodeData{}, previous_node.Children, spanFromNode(previous_node))
				// Add function declararation node
				declSpan := mergeSpan(spanFromNode(previous_node), spanFromNode(scope_node))
				output = append(output, nodeWithSpan("function_declaration", NodeData{FunctionName: previous_node.Data.FunctionName}, []Node{arg_node, scope_node}, declSpan))
				diagnostics = append(diagnostics, scope_diags...)
				index += rawLen
			case "for_header":
				output = output[:len(output)-1]
				scope_node, scope_diags, rawLen := parseScope(tokens, index)
				for_children := append(previous_node.Children, scope_node)
				loopSpan := mergeSpan(spanFromNode(previous_node), spanFromNode(scope_node))
				output = append(output, nodeWithSpan("for_loop", NodeData{}, for_children, loopSpan))
				diagnostics = append(diagnostics, scope_diags...)
				index += rawLen
			case "if_header":
				output = output[:len(output)-1]
				scope_node, scope_diags, rawLen := parseScope(tokens, index)
				if_children := append(previous_node.Children, scope_node)
				ifSpan := mergeSpan(spanFromNode(previous_node), spanFromNode(scope_node))
				output = append(output, nodeWithSpan("if_statement", NodeData{}, if_children, ifSpan))
				diagnostics = append(diagnostics, scope_diags...)
				index += rawLen
			case "while_header":
				output = output[:len(output)-1]
				scope_node, scope_diags, rawLen := parseScope(tokens, index)
				while_children := append(previous_node.Children, scope_node)
				whileSpan := mergeSpan(spanFromNode(previous_node), spanFromNode(scope_node))
				output = append(output, nodeWithSpan("while_loop", NodeData{}, while_children, whileSpan))
				diagnostics = append(diagnostics, scope_diags...)
				index += rawLen
			case "foreach_header":
				output = output[:len(output)-1]
				scope_node, scope_diags, rawLen := parseScope(tokens, index)
				foreach_children := append(previous_node.Children, scope_node)
				foreachSpan := mergeSpan(spanFromNode(previous_node), spanFromNode(scope_node))
				output = append(output, nodeWithSpan("foreach_loop", NodeData{}, foreach_children, foreachSpan))
				diagnostics = append(diagnostics, scope_diags...)
				index += rawLen
			case "switch_header":
				output = output[:len(output)-1]
				scope_node, scope_diags, rawLen := parseScope(tokens, index)
				switch_children := append(previous_node.Children, scope_node)
				switchSpan := mergeSpan(spanFromNode(previous_node), spanFromNode(scope_node))
				output = append(output, nodeWithSpan("switch_statement", NodeData{}, switch_children, switchSpan))
				diagnostics = append(diagnostics, scope_diags...)
				index += rawLen
			case "else_header":
				output = output[:len(output)-1]
				scope_node, scope_diags, rawLen := parseScope(tokens, index)
				elseSpan := mergeSpan(spanFromNode(previous_node), spanFromNode(scope_node))
				output = append(output, nodeWithSpan("else_clause", NodeData{}, []Node{scope_node}, elseSpan))
				diagnostics = append(diagnostics, scope_diags...)
				index += rawLen
			case "do_header":
				output = output[:len(output)-1]
				scope_node, scope_diags, rawLen := parseScope(tokens, index)
				diagnostics = append(diagnostics, scope_diags...)

				condition_children := []Node{}
				consumedIndex := index + rawLen
				afterClose := consumedIndex + 1
				if afterClose < len(tokens) && tokens[afterClose].Type == l.SYMBOL && tokens[afterClose].Content == "while" {
					if afterClose+1 < len(tokens) && tokens[afterClose+1].Type == l.OPEN_PAREN {
						condRaw, _ := tokensUntilMatchingClose(tokens[afterClose+2:], l.OPEN_PAREN, l.CLOSE_PAREN)
						condTokens := trimTrailingToken(condRaw, l.CLOSE_PAREN)
						condChildren, condDiags := Parse(condTokens)
						diagnostics = append(diagnostics, condDiags...)
						condition_children = condChildren
						consumedIndex = afterClose + 1 + len(condRaw)
						if consumedIndex+1 < len(tokens) && tokens[consumedIndex+1].Type == l.TERMINATOR {
							consumedIndex++
						}
					}
				}

				condSpan := spanFromNodes(condition_children)
				if afterClose < len(tokens) {
					condSpan = mergeSpan(spanFromToken(tokens[afterClose]), condSpan)
				}
				condNode := nodeWithSpan("condition", NodeData{}, condition_children, condSpan)
				endSpan := spanFromNode(scope_node)
				if consumedIndex < len(tokens) {
					endSpan = spanFromToken(tokens[consumedIndex])
				}
				doSpan := mergeSpan(spanFromNode(previous_node), endSpan)
				do_node := nodeWithSpan("do_while_loop", NodeData{}, []Node{condNode, scope_node}, doSpan)
				output = append(output, do_node)
				index = consumedIndex
			default:
				diagnostics = append(diagnostics, diagnosticAtIndex("unexpected {", tokens, index, "error"))
				output = append(output, nodeWithSpan("open_curly", NodeData{}, []Node{}, spanFromToken(tokens[index])))
				break
			}
		case l.CLOSE_PAREN:
			diagnostics = append(diagnostics, diagnosticAtIndex("unexpected )", tokens, index, "error"))
		case l.CLOSE_BRACKET:
			diagnostics = append(diagnostics, diagnosticAtIndex("unexpected ]", tokens, index, "error"))
		case l.CLOSE_CURLY:
			diagnostics = append(diagnostics, diagnosticAtIndex("unexpected }", tokens, index, "error"))
		case l.COLON:
			// Ignore colons (used in case/default labels)
		default:
		}
		index++
	}

	return output, diagnostics
}
