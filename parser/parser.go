package parser

import (
	"strings"

	d "github.com/maxvanasten/gscp/diagnostics"
	l "github.com/maxvanasten/gscp/lexer"
)

type NodeData struct {
	VarName      string `json:"variable_name,omitempty"`
	FunctionName string `json:"function_name,omitempty"`
	Path         string `json:"path,omitempty"`
	Operator     string `json:"operator,omitempty"`
	Delay        string `json:"delay,omitempty"`
	Thread       bool   `json:"thread,omitempty"`
	Method       string `json:"method,omitempty"`
	Index        string `json:"index,omitempty"`
	Content      string `json:"content,omitempty"`
}

type Node struct {
	Type     string   `json:"type"`
	Data     NodeData `json:"data"`
	Children []Node   `json:"children,omitempty"`
}

func tokensToString(tokens []l.Token) string {
	var builder strings.Builder
	for _, token := range tokens {
		switch token.Type {
		case l.STRING:
			builder.WriteString("\"")
			builder.WriteString(token.Content)
			builder.WriteString("\"")
		default:
			builder.WriteString(token.Content)
		}
	}

	return builder.String()
}

func tokensUntilMatchingClose(tokens []l.Token, openType l.TokenType, closeType l.TokenType) ([]l.Token, bool) {
	depth := 0
	for i, token := range tokens {
		switch token.Type {
		case openType:
			depth++
		case closeType:
			if depth == 0 {
				return tokens[:i+1], true
			}
			depth--
		}
	}

	return tokens, false
}

func trimTrailingToken(tokens []l.Token, tokenType l.TokenType) []l.Token {
	if len(tokens) == 0 {
		return tokens
	}
	if tokens[len(tokens)-1].Type == tokenType {
		return tokens[:len(tokens)-1]
	}
	return tokens
}

func trimTrailingAny(tokens []l.Token, tokenTypes ...l.TokenType) []l.Token {
	if len(tokens) == 0 {
		return tokens
	}
	last := tokens[len(tokens)-1].Type
	for _, t := range tokenTypes {
		if last == t {
			return tokens[:len(tokens)-1]
		}
	}
	return tokens
}

func diagnosticAtToken(message string, token l.Token, severity string) d.Diagnostic {
	line := token.Line
	col := token.Col
	endLine := token.EndLine
	endCol := token.EndCol
	return d.New(message, line, col, endLine, endCol, severity)
}

func diagnosticAtIndex(message string, tokens []l.Token, index int, severity string) d.Diagnostic {
	if index < 0 || index >= len(tokens) {
		return d.New(message, 0, 0, 0, 0, severity)
	}
	return diagnosticAtToken(message, tokens[index], severity)
}

func Parse(tokens []l.Token) ([]Node, []d.Diagnostic) {
	output := []Node{}
	diagnostics := []d.Diagnostic{}

	index := 0
	for index < len(tokens) {
		switch tokens[index].Type {
		case l.STRING:
			output = append(output, Node{"string", NodeData{Content: tokens[index].Content}, []Node{}})
		case l.SYMBOL:
			normalized := strings.TrimSuffix(tokens[index].Content, ":")
			switch normalized {
			case "#include":
				if index+1 < len(tokens) && tokens[index+1].Type == l.SYMBOL {
					output = append(output, Node{"include_statement", NodeData{Path: tokens[index+1].Content}, []Node{}})
					index++
				} else {
					diagnostics = append(diagnostics, diagnosticAtIndex("missing include path", tokens, index, "error"))
				}
			case "wait":
				if index+1 < len(tokens) {
					next := tokens[index+1]
					if next.Type == l.OPEN_PAREN {
						output = append(output, Node{"variable_reference", NodeData{VarName: tokens[index].Content}, []Node{}})
						break
					}
					if next.Type == l.NUMBER || next.Type == l.SYMBOL {
						output = append(output, Node{"wait_statement", NodeData{Delay: next.Content}, []Node{}})
						index++
						break
					}
				}
				diagnostics = append(diagnostics, diagnosticAtIndex("missing wait duration", tokens, index, "error"))
			case "thread":
				output = append(output, Node{"thread_keyword", NodeData{}, []Node{}})
			case "true", "false":
				output = append(output, Node{"boolean", NodeData{Content: normalized}, []Node{}})
			case "break":
				output = append(output, Node{"break_statement", NodeData{}, []Node{}})
			case "return":
				ret_tokens := l.TokensUntilAny(tokens[index+1:], []l.TokenType{l.NEWLINE, l.TERMINATOR})
				ret_children, diags := Parse(ret_tokens)
				output = append(output, Node{"return_statement", NodeData{}, ret_children})
				diagnostics = append(diagnostics, diags...)
				index += len(ret_tokens)
			case "case":
				case_tokens := l.TokensUntilAny(tokens[index+1:], []l.TokenType{l.NEWLINE, l.TERMINATOR, l.COLON})
				case_children, diags := Parse(case_tokens)
				output = append(output, Node{"case_clause", NodeData{}, case_children})
				diagnostics = append(diagnostics, diags...)
				index += len(case_tokens)
			case "default":
				output = append(output, Node{"default_clause", NodeData{}, []Node{}})
			case "else":
				output = append(output, Node{"else_header", NodeData{}, []Node{}})
			case "do":
				output = append(output, Node{"do_header", NodeData{}, []Node{}})
			default:
				output = append(output, Node{"variable_reference", NodeData{VarName: tokens[index].Content}, []Node{}})
			}

		case l.NUMBER:
			output = append(output, Node{"number", NodeData{Content: tokens[index].Content}, []Node{}})
		case l.OPERATOR:
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
							output = append(output, Node{"unary_expression", NodeData{Operator: tokens[index].Content}, []Node{operand}})
							diagnostics = append(diagnostics, diags...)
							index += len(operand_tokens)
							break
						}
						diagnostics = append(diagnostics, diags...)
					}
					diagnostics = append(diagnostics, diagnosticAtIndex("missing unary operand", tokens, index, "error"))
				}
			}
			if tokens[index].Content == "++" || tokens[index].Content == "--" {
				if len(output) > 0 {
					previous_node := output[len(output)-1]
					if previous_node.Type == "variable_reference" {
						output = output[:len(output)-1]
						operator := "+"
						if tokens[index].Content == "--" {
							operator = "-"
						}
						lhs := Node{"lhs", NodeData{}, []Node{previous_node}}
						rhs := Node{"rhs", NodeData{}, []Node{{"number", NodeData{Content: "1"}, []Node{}}}}
						expr := Node{"expression", NodeData{Operator: operator}, []Node{lhs, rhs}}
						assignment_data := NodeData{VarName: previous_node.Data.VarName, Index: previous_node.Data.Index}
						output = append(output, Node{"assignment", assignment_data, []Node{expr}})
						break
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
							lhs := Node{"lhs", NodeData{}, []Node{operand}}
							rhs := Node{"rhs", NodeData{}, []Node{{"number", NodeData{Content: "1"}, []Node{}}}}
							expr := Node{"expression", NodeData{Operator: operator}, []Node{lhs, rhs}}
							assignment_data := NodeData{VarName: operand.Data.VarName, Index: operand.Data.Index}
							output = append(output, Node{"assignment", assignment_data, []Node{expr}})
							diagnostics = append(diagnostics, diags...)
							index += len(operand_tokens)
							break
						}
					}
					diagnostics = append(diagnostics, diags...)
				}
			}
			if tokens[index].Content == "?" {
				if len(output) == 0 {
					diagnostics = append(diagnostics, diagnosticAtIndex("operator missing left-hand operand", tokens, index, "error"))
					break
				}
				condition := output[len(output)-1]
				output = output[:len(output)-1]
				expr_tokens := l.TokensUntilAny(tokens[index+1:], []l.TokenType{l.NEWLINE, l.TERMINATOR})
				depthParen := 0
				depthBracket := 0
				splitIndex := -1
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
					if depthParen == 0 && depthBracket == 0 && tok.Type == l.COLON {
						splitIndex = i
						break
					}
				}
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
				ternary := Node{"ternary_expression", NodeData{}, []Node{
					{"condition", NodeData{}, []Node{condition}},
					{"true_expr", NodeData{}, true_children},
					{"false_expr", NodeData{}, false_children},
				}}
				output = append(output, ternary)
				diagnostics = append(diagnostics, true_diags...)
				diagnostics = append(diagnostics, false_diags...)
				index += len(expr_tokens)
				break
			}
			if index <= 0 {
				diagnostics = append(diagnostics, diagnosticAtIndex("operator missing left-hand operand", tokens, index, "error"))
				break
			}
			if len(output) == 0 {
				diagnostics = append(diagnostics, diagnosticAtIndex("operator missing left-hand operand", tokens, index, "error"))
				break
			}
			// Check if previous node is either a string, variable_reference, number, function_call or another expression
			previous_node := output[len(output)-1]
			switch previous_node.Type {
			case "string", "variable_reference", "number", "function_call", "expression":
				// Set LHS to previous node
				lhs := Node{"lhs", NodeData{}, []Node{previous_node}}
				// Delete previous node
				output = output[:len(output)-1]
				// Get all tokens from OPERATOR until END, NEWLINE or TERMINATOR
				expr_tokens := l.TokensUntilAny(tokens[index+1:], []l.TokenType{l.NEWLINE, l.TERMINATOR})
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
				rhs := Node{"rhs", NodeData{}, rhs_children}
				if len(rhs_children) == 0 {
					diagnostics = append(diagnostics, diagnosticAtIndex("operator missing right-hand operand", tokens, index, "error"))
				}
				// Add Expression node to output
				output = append(output, Node{"expression", NodeData{Operator: tokens[index].Content}, []Node{lhs, rhs}})
				diagnostics = append(diagnostics, diags...)
				index += len(expr_tokens)
			default:
				output = append(output, Node{"operator", NodeData{Content: tokens[index].Content}, []Node{}})
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
				output = append(output, Node{"assignment", NodeData{Content: tokens[index].Content}, []Node{}})
				break
			}
			output = output[:len(output)-1]
			// Get all tokens from ASSIGNMENT until END, NEWLINE or TERMINATOR
			ass_tokens := l.TokensUntilAny(tokens[index+1:], []l.TokenType{l.NEWLINE, l.TERMINATOR})

			ass_children, diags := Parse(ass_tokens)
			assignment_data := NodeData{VarName: previous_node.Data.VarName, Index: previous_node.Data.Index}
			if tokens[index].Content == "+=" || tokens[index].Content == "-=" || tokens[index].Content == "*=" || tokens[index].Content == "/=" {
				operator := tokens[index].Content[:1]
				lhs := Node{"lhs", NodeData{}, []Node{previous_node}}
				rhs := Node{"rhs", NodeData{}, ass_children}
				expr := Node{"expression", NodeData{Operator: operator}, []Node{lhs, rhs}}
				output = append(output, Node{"assignment", assignment_data, []Node{expr}})
			} else {
				output = append(output, Node{"assignment", assignment_data, ass_children})
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

					segments := [][]l.Token{}
					buf := []l.Token{}
					depth := 0
					for _, ht := range header_tokens {
						switch ht.Type {
						case l.OPEN_PAREN:
							depth++
							buf = append(buf, ht)
						case l.CLOSE_PAREN:
							if depth > 0 {
								depth--
							}
							buf = append(buf, ht)
						case l.TERMINATOR:
							if depth == 0 {
								segments = append(segments, buf)
								buf = []l.Token{}
							} else {
								buf = append(buf, ht)
							}
						default:
							buf = append(buf, ht)
						}
					}
					segments = append(segments, buf)
					for len(segments) < 3 {
						segments = append(segments, []l.Token{})
					}

					init_children, init_diags := Parse(segments[0])
					cond_children, cond_diags := Parse(segments[1])
					post_children, post_diags := Parse(segments[2])
					for_header := Node{"for_header", NodeData{}, []Node{
						{"for_init", NodeData{}, init_children},
						{"for_condition", NodeData{}, cond_children},
						{"for_post", NodeData{}, post_children},
					}}

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
					output = output[:len(output)-1]
					switch previous_node.Data.VarName {
					case "if":
						cond_children, diags := Parse(header_tokens)
						output = append(output, Node{"if_header", NodeData{}, []Node{{"condition", NodeData{}, cond_children}}})
						diagnostics = append(diagnostics, diags...)
					case "while":
						cond_children, diags := Parse(header_tokens)
						output = append(output, Node{"while_header", NodeData{}, []Node{{"condition", NodeData{}, cond_children}}})
						diagnostics = append(diagnostics, diags...)
					case "switch":
						expr_children, diags := Parse(header_tokens)
						output = append(output, Node{"switch_header", NodeData{}, []Node{{"switch_expr", NodeData{}, expr_children}}})
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
						output = append(output, Node{"foreach_header", NodeData{}, []Node{
							{"foreach_vars", NodeData{}, left_children},
							{"foreach_iter", NodeData{}, right_children},
						}})
						diagnostics = append(diagnostics, left_diags...)
						diagnostics = append(diagnostics, right_diags...)
					}
					index += len(arg_tokens)
					break
				}
			}

			if len(output)-1 < 0 || output[len(output)-1].Type != "variable_reference" {
				inner_tokens := trimmedArgTokens
				depth := 0
				hasComma := false
				for _, tok := range inner_tokens {
					switch tok.Type {
					case l.OPEN_PAREN:
						depth++
					case l.CLOSE_PAREN:
						if depth > 0 {
							depth--
						}
					case l.COMMA:
						if depth == 0 {
							hasComma = true
							break
						}
					}
					if hasComma {
						break
					}
				}
				if hasComma {
					elem_slices := [][]l.Token{}
					buf := []l.Token{}
					depthParen := 0
					depthBracket := 0
					for _, at := range inner_tokens {
						switch at.Type {
						case l.OPEN_PAREN:
							depthParen++
							buf = append(buf, at)
						case l.CLOSE_PAREN:
							if depthParen > 0 {
								depthParen--
							}
							buf = append(buf, at)
						case l.OPEN_BRACKET:
							depthBracket++
							buf = append(buf, at)
						case l.CLOSE_BRACKET:
							if depthBracket > 0 {
								depthBracket--
							}
							buf = append(buf, at)
						case l.COMMA:
							if depthParen == 0 && depthBracket == 0 {
								elem_slices = append(elem_slices, buf)
								buf = []l.Token{}
							} else {
								buf = append(buf, at)
							}
						default:
							buf = append(buf, at)
						}
					}
					if len(buf) > 0 {
						elem_slices = append(elem_slices, buf)
					}
					vec_children := []Node{}
					for _, elem_tokens := range elem_slices {
						children, diags := Parse(elem_tokens)
						vec_children = append(vec_children, children...)
						diagnostics = append(diagnostics, diags...)
					}
					output = append(output, Node{"vector_literal", NodeData{}, vec_children})
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
			buf := []l.Token{}
			depth := 0
			for _, at := range trimmedArgTokens {
				switch at.Type {
				case l.OPEN_PAREN:
					depth++
					buf = append(buf, at)
				case l.CLOSE_PAREN:
					if depth > 0 {
						depth--
					}
					buf = append(buf, at)
				case l.COMMA:
					if depth == 0 {
						arg_token_slices = append(arg_token_slices, buf)
						buf = []l.Token{}
					} else {
						buf = append(buf, at)
					}
				default:
					buf = append(buf, at)
				}
			}
			arg_token_slices = append(arg_token_slices, buf)

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
									data.Method = output[len(output)-3].Data.VarName
								}
							}
						}
						if output[len(output)-2].Type == "variable_reference" {
							c = 2
							data.Method = output[len(output)-2].Data.VarName
						}
					}
				}
				output = output[:len(output)-c]
				output = append(output, Node{"function_call", data, arg_children})
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
					indexed_node := Node{"variable_reference", NodeData{VarName: previous_node.Data.VarName, Index: index_content}, []Node{}}
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
			elem_slices := [][]l.Token{}
			buf := []l.Token{}
			depthParen := 0
			depthBracket := 0
			for _, at := range array_tokens {
				switch at.Type {
				case l.OPEN_PAREN:
					depthParen++
					buf = append(buf, at)
				case l.CLOSE_PAREN:
					if depthParen > 0 {
						depthParen--
					}
					buf = append(buf, at)
				case l.OPEN_BRACKET:
					depthBracket++
					buf = append(buf, at)
				case l.CLOSE_BRACKET:
					if depthBracket > 0 {
						depthBracket--
					}
					buf = append(buf, at)
				case l.COMMA:
					if depthParen == 0 && depthBracket == 0 {
						elem_slices = append(elem_slices, buf)
						buf = []l.Token{}
					} else {
						buf = append(buf, at)
					}
				default:
					buf = append(buf, at)
				}
			}
			if len(buf) > 0 {
				elem_slices = append(elem_slices, buf)
			}

			array_children := []Node{}
			for _, elem_tokens := range elem_slices {
				children, diags := Parse(elem_tokens)
				array_children = append(array_children, children...)
				diagnostics = append(diagnostics, diags...)
			}
			output = append(output, Node{"array_literal", NodeData{}, array_children})
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
				// Get all tokens from OPEN_CURLY until matching CLOSE_CURLY
				rawScopeTokens, foundClose := tokensUntilMatchingClose(tokens[index+1:], l.OPEN_CURLY, l.CLOSE_CURLY)
				if !foundClose {
					diagnostics = append(diagnostics, diagnosticAtIndex("missing closing }", tokens, index, "error"))
				}
				scope_tokens := trimTrailingToken(rawScopeTokens, l.CLOSE_CURLY)
				// Parse those tokens into scope node
				arg_node := Node{"args", NodeData{}, previous_node.Children}
				scope_children, diags := Parse(scope_tokens)
				scope_node := Node{"scope", NodeData{}, scope_children}
				// Add function declararation node
				output = append(output, Node{"function_declaration", NodeData{FunctionName: previous_node.Data.FunctionName}, []Node{arg_node, scope_node}})
				diagnostics = append(diagnostics, diags...)
				index += len(rawScopeTokens)
			case "for_header":
				output = output[:len(output)-1]
				rawScopeTokens, foundClose := tokensUntilMatchingClose(tokens[index+1:], l.OPEN_CURLY, l.CLOSE_CURLY)
				if !foundClose {
					diagnostics = append(diagnostics, diagnosticAtIndex("missing closing }", tokens, index, "error"))
				}
				scope_tokens := trimTrailingToken(rawScopeTokens, l.CLOSE_CURLY)
				scope_children, diags := Parse(scope_tokens)
				scope_node := Node{"scope", NodeData{}, scope_children}
				for_children := append(previous_node.Children, scope_node)
				output = append(output, Node{"for_loop", NodeData{}, for_children})
				diagnostics = append(diagnostics, diags...)
				index += len(rawScopeTokens)
			case "if_header":
				output = output[:len(output)-1]
				rawScopeTokens, foundClose := tokensUntilMatchingClose(tokens[index+1:], l.OPEN_CURLY, l.CLOSE_CURLY)
				if !foundClose {
					diagnostics = append(diagnostics, diagnosticAtIndex("missing closing }", tokens, index, "error"))
				}
				scope_tokens := trimTrailingToken(rawScopeTokens, l.CLOSE_CURLY)
				scope_children, diags := Parse(scope_tokens)
				scope_node := Node{"scope", NodeData{}, scope_children}
				if_children := append(previous_node.Children, scope_node)
				output = append(output, Node{"if_statement", NodeData{}, if_children})
				diagnostics = append(diagnostics, diags...)
				index += len(rawScopeTokens)
			case "while_header":
				output = output[:len(output)-1]
				rawScopeTokens, foundClose := tokensUntilMatchingClose(tokens[index+1:], l.OPEN_CURLY, l.CLOSE_CURLY)
				if !foundClose {
					diagnostics = append(diagnostics, diagnosticAtIndex("missing closing }", tokens, index, "error"))
				}
				scope_tokens := trimTrailingToken(rawScopeTokens, l.CLOSE_CURLY)
				scope_children, diags := Parse(scope_tokens)
				scope_node := Node{"scope", NodeData{}, scope_children}
				while_children := append(previous_node.Children, scope_node)
				output = append(output, Node{"while_loop", NodeData{}, while_children})
				diagnostics = append(diagnostics, diags...)
				index += len(rawScopeTokens)
			case "foreach_header":
				output = output[:len(output)-1]
				rawScopeTokens, foundClose := tokensUntilMatchingClose(tokens[index+1:], l.OPEN_CURLY, l.CLOSE_CURLY)
				if !foundClose {
					diagnostics = append(diagnostics, diagnosticAtIndex("missing closing }", tokens, index, "error"))
				}
				scope_tokens := trimTrailingToken(rawScopeTokens, l.CLOSE_CURLY)
				scope_children, diags := Parse(scope_tokens)
				scope_node := Node{"scope", NodeData{}, scope_children}
				foreach_children := append(previous_node.Children, scope_node)
				output = append(output, Node{"foreach_loop", NodeData{}, foreach_children})
				diagnostics = append(diagnostics, diags...)
				index += len(rawScopeTokens)
			case "switch_header":
				output = output[:len(output)-1]
				rawScopeTokens, foundClose := tokensUntilMatchingClose(tokens[index+1:], l.OPEN_CURLY, l.CLOSE_CURLY)
				if !foundClose {
					diagnostics = append(diagnostics, diagnosticAtIndex("missing closing }", tokens, index, "error"))
				}
				scope_tokens := trimTrailingToken(rawScopeTokens, l.CLOSE_CURLY)
				scope_children, diags := Parse(scope_tokens)
				scope_node := Node{"scope", NodeData{}, scope_children}
				switch_children := append(previous_node.Children, scope_node)
				output = append(output, Node{"switch_statement", NodeData{}, switch_children})
				diagnostics = append(diagnostics, diags...)
				index += len(rawScopeTokens)
			case "else_header":
				output = output[:len(output)-1]
				rawScopeTokens, foundClose := tokensUntilMatchingClose(tokens[index+1:], l.OPEN_CURLY, l.CLOSE_CURLY)
				if !foundClose {
					diagnostics = append(diagnostics, diagnosticAtIndex("missing closing }", tokens, index, "error"))
				}
				scope_tokens := trimTrailingToken(rawScopeTokens, l.CLOSE_CURLY)
				scope_children, diags := Parse(scope_tokens)
				scope_node := Node{"scope", NodeData{}, scope_children}
				output = append(output, Node{"else_clause", NodeData{}, []Node{scope_node}})
				diagnostics = append(diagnostics, diags...)
				index += len(rawScopeTokens)
			case "do_header":
				output = output[:len(output)-1]
				rawScopeTokens, foundClose := tokensUntilMatchingClose(tokens[index+1:], l.OPEN_CURLY, l.CLOSE_CURLY)
				if !foundClose {
					diagnostics = append(diagnostics, diagnosticAtIndex("missing closing }", tokens, index, "error"))
				}
				scope_tokens := trimTrailingToken(rawScopeTokens, l.CLOSE_CURLY)
				scope_children, scope_diags := Parse(scope_tokens)
				diagnostics = append(diagnostics, scope_diags...)
				scope_node := Node{"scope", NodeData{}, scope_children}

				condition_children := []Node{}
				consumedIndex := index + len(rawScopeTokens)
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

				do_node := Node{"do_while_loop", NodeData{}, []Node{
					{"condition", NodeData{}, condition_children},
					scope_node,
				}}
				output = append(output, do_node)
				index = consumedIndex
			default:
				diagnostics = append(diagnostics, diagnosticAtIndex("unexpected {", tokens, index, "error"))
				output = append(output, Node{"open_curly", NodeData{}, []Node{}})
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
