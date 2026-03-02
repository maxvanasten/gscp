package parser

import (
	d "github.com/maxvanasten/gscp/diagnostics"
	l "github.com/maxvanasten/gscp/lexer"
)

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
