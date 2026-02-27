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

func Parse(tokens []l.Token) ([]Node, []d.Diagnostic) {
	output := []Node{}
	diagnostics := []d.Diagnostic{}

	index := 0
	for index < len(tokens) {
		switch tokens[index].Type {
		case l.STRING:
			output = append(output, Node{"string", NodeData{Content: tokens[index].Content}, []Node{}})
		case l.SYMBOL:
			switch tokens[index].Content {
			case "#include":
				if index+1 < len(tokens) {
					output = append(output, Node{"include_statement", NodeData{Path: tokens[index+1].Content}, []Node{}})
					index++
				}
			case "wait":
				if index+1 < len(tokens) && (tokens[index+1].Type == l.NUMBER || tokens[index+1].Type == l.SYMBOL) {
					output = append(output, Node{"wait_statement", NodeData{Delay: tokens[index+1].Content}, []Node{}})
					index++
				}
			case "thread":
				output = append(output, Node{"thread_keyword", NodeData{}, []Node{}})
			default:
				output = append(output, Node{"variable_reference", NodeData{VarName: tokens[index].Content}, []Node{}})
			}

		case l.NUMBER:
			output = append(output, Node{"number", NodeData{Content: tokens[index].Content}, []Node{}})
		case l.OPERATOR:
			if index <= 0 {
				break
			}
			if len(output) == 0 {
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
				// Parse those tokens into RHS
				rhs_children, diags := Parse(expr_tokens)
				rhs := Node{"rhs", NodeData{}, rhs_children}
				// Add Expression node to output
				output = append(output, Node{"expression", NodeData{Operator: tokens[index].Content}, []Node{lhs, rhs}})
				diagnostics = append(diagnostics, diags...)
				index += len(expr_tokens)
			default:
				output = append(output, Node{"operator", NodeData{Content: tokens[index].Content}, []Node{}})
			}
		case l.ASSIGNMENT:
			if index <= 0 {
				break
			}
			// Check if previous node is a variable_reference
			previous_node := output[len(output)-1]
			if previous_node.Type != "variable_reference" {
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
			if index <= 0 {
				break
			}

			arg_tokens := l.TokensUntilAny(tokens[index+1:], []l.TokenType{l.CLOSE_PAREN})
			if len(output)-1 >= 0 {
				previous_node := output[len(output)-1]
				if previous_node.Type == "variable_reference" && previous_node.Data.VarName == "for" {
					header_tokens := arg_tokens
					if len(header_tokens) > 0 && header_tokens[len(header_tokens)-1].Type == l.CLOSE_PAREN {
						header_tokens = header_tokens[:len(header_tokens)-1]
					}

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
			}

			arg_token_slices := [][]l.Token{}
			// Split arg_tokens into slices at top-level commas only
			buf := []l.Token{}
			depth := 0
			for _, at := range arg_tokens {
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
			diagnostics := []d.Diagnostic{}
			for _, argument_tokens := range arg_token_slices {
				children, diags := Parse(argument_tokens)
				arg_children = append(arg_children, children...)
				diagnostics = append(diagnostics, diags...)
			}

			data := NodeData{}
			if len(output)-1 >= 0 {
				c := 0
				if output[len(output)-1].Type == "variable_reference" {
					c = 1
					data.FunctionName = output[len(output)-1].Data.VarName
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

			index += len(arg_tokens)
		case l.OPEN_BRACKET:
			bracket_tokens := l.TokensUntilAny(tokens[index+1:], []l.TokenType{l.CLOSE_BRACKET})
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
				break
			}
			previous_node := output[len(output)-1]
			switch previous_node.Type {
			case "function_call":
				output = output[:len(output)-1]
				// Get all tokens from OPEN_CURLY until CLOSE_CURLY
				scope_tokens := l.TokensUntilAny(tokens[index+1:], []l.TokenType{l.CLOSE_CURLY})
				// Parse those tokens into scope node
				arg_node := Node{"args", NodeData{}, previous_node.Children}
				scope_children, diags := Parse(scope_tokens)
				scope_node := Node{"scope", NodeData{}, scope_children}
				// Add function declararation node
				output = append(output, Node{"function_declaration", NodeData{FunctionName: previous_node.Data.FunctionName}, []Node{arg_node, scope_node}})
				diagnostics = append(diagnostics, diags...)
				index += len(scope_tokens)
			case "for_header":
				output = output[:len(output)-1]
				scope_tokens := l.TokensUntilAny(tokens[index+1:], []l.TokenType{l.CLOSE_CURLY})
				scope_children, diags := Parse(scope_tokens)
				scope_node := Node{"scope", NodeData{}, scope_children}
				for_children := append(previous_node.Children, scope_node)
				output = append(output, Node{"for_loop", NodeData{}, for_children})
				diagnostics = append(diagnostics, diags...)
				index += len(scope_tokens)
			default:
				output = append(output, Node{"open_curly", NodeData{}, []Node{}})
				break
			}
		default:
		}
		index++
	}

	return output, diagnostics
}
