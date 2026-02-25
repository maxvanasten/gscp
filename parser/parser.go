package parser

import (
	l "github.com/maxvanasten/gscp/lexer"
)

type NodeData struct {
	VarName      string `json:"variable_name,omitempty"`
	FunctionName string `json:"function_name,omitempty"`
	Path         string `json:"path,omitempty"`
	Operator     string `json:"operator,omitempty"`
	Content      string `json:"content,omitempty"`
}

type Node struct {
	Type     string   `json:"type"`
	Data     NodeData `json:"data"`
	Children []Node   `json:"children,omitempty"`
}

func Parse(tokens []l.Token) []Node {
	output := []Node{}

	index := 0
	for index < len(tokens) {
		switch tokens[index].Type {
		case l.STRING:
			output = append(output, Node{"string", NodeData{Content: tokens[index].Content}, []Node{}})
		case l.SYMBOL:
			if tokens[index].Content == "#include" {
				if index+1 < len(tokens) {
					output = append(output, Node{"include_statement", NodeData{Path: tokens[index+1].Content}, []Node{}})
					index++
					break
				}
			}

			output = append(output, Node{"variable_reference", NodeData{VarName: tokens[index].Content}, []Node{}})
		case l.NUMBER:
			output = append(output, Node{"number", NodeData{Content: tokens[index].Content}, []Node{}})
		case l.OPERATOR:
			if index <= 0 {
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
				rhs := Node{"rhs", NodeData{}, Parse(expr_tokens)}
				// Add Expression node to output
				output = append(output, Node{"expression", NodeData{Operator: tokens[index].Content}, []Node{lhs, rhs}})
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

			output = append(output, Node{"variable_assignment", NodeData{VarName: previous_node.Data.VarName}, Parse(ass_tokens)})
			index += len(ass_tokens)
		case l.OPEN_PAREN:
			if index <= 0 {
				break
			}
			// Check if previous node is a variable_reference
			previous_node := output[len(output)-1]
			if previous_node.Type != "variable_reference" {
				break
			}
			output = output[:len(output)-1]
			// Get all tokens from OPEN_PAREN until CLOSE_PAREN
			arg_tokens := l.TokensUntilAny(tokens[index+1:], []l.TokenType{l.CLOSE_PAREN})
			// Add function call node
			output = append(output, Node{"function_call", NodeData{FunctionName: previous_node.Data.VarName}, Parse(arg_tokens)})
			index += len(arg_tokens)
		case l.OPEN_CURLY:
			// Check if previous node is a function_call
			previous_node := output[len(output)-1]
			if previous_node.Type != "function_call" {
				output = append(output, Node{"open_curly", NodeData{}, []Node{}})
				break
			}
			output = output[:len(output)-1]
			// Get all tokens from OPEN_CURLY until CLOSE_CURLY
			scope_tokens := l.TokensUntilAny(tokens[index+1:], []l.TokenType{l.CLOSE_CURLY})
			// Parse those tokens into scope node
			arg_node := Node{"args", NodeData{}, previous_node.Children}
			scope_node := Node{"scope", NodeData{}, Parse(scope_tokens)}
			// Add function declararation node
			output = append(output, Node{"function_declaration", NodeData{FunctionName: previous_node.Data.FunctionName}, []Node{arg_node, scope_node}})
			index += len(scope_tokens)
		default:
		}
		index++
	}

	return output
}
