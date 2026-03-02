package parser

import (
	d "github.com/maxvanasten/gscp/diagnostics"
	l "github.com/maxvanasten/gscp/lexer"
)

func parseScope(tokens []l.Token, index int) (Node, []d.Diagnostic, int) {
	rawScopeTokens, foundClose := tokensUntilMatchingClose(tokens[index+1:], l.OPEN_CURLY, l.CLOSE_CURLY)
	diagnostics := []d.Diagnostic{}
	if !foundClose {
		diagnostics = append(diagnostics, diagnosticAtIndex("missing closing }", tokens, index, "error"))
	}
	scopeTokens := trimTrailingToken(rawScopeTokens, l.CLOSE_CURLY)
	scopeChildren, scopeDiags := Parse(scopeTokens)
	diagnostics = append(diagnostics, scopeDiags...)
	startSpan := spanFromToken(tokens[index])
	endSpan := startSpan
	if len(rawScopeTokens) > 0 {
		endSpan = spanFromToken(rawScopeTokens[len(rawScopeTokens)-1])
	}
	scopeNode := nodeWithSpan("scope", NodeData{}, scopeChildren, mergeSpan(startSpan, endSpan))
	return scopeNode, diagnostics, len(rawScopeTokens)
}
