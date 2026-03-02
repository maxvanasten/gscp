package parser

import (
	"strings"

	l "github.com/maxvanasten/gscp/lexer"
)

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

func lastNonTokenType(tokens []l.Token, tokenTypes ...l.TokenType) (l.Token, bool) {
	for i := len(tokens) - 1; i >= 0; i-- {
		ok := true
		for _, t := range tokenTypes {
			if tokens[i].Type == t {
				ok = false
				break
			}
		}
		if ok {
			return tokens[i], true
		}
	}
	return l.Token{}, false
}

func splitTopLevel(tokens []l.Token, delimiter l.TokenType, dropTrailingEmpty bool) [][]l.Token {
	segments := [][]l.Token{}
	buf := []l.Token{}
	depthParen := 0
	depthBracket := 0
	for _, tok := range tokens {
		switch tok.Type {
		case l.OPEN_PAREN:
			depthParen++
			buf = append(buf, tok)
		case l.CLOSE_PAREN:
			if depthParen > 0 {
				depthParen--
			}
			buf = append(buf, tok)
		case l.OPEN_BRACKET:
			depthBracket++
			buf = append(buf, tok)
		case l.CLOSE_BRACKET:
			if depthBracket > 0 {
				depthBracket--
			}
			buf = append(buf, tok)
		case delimiter:
			if depthParen == 0 && depthBracket == 0 {
				segments = append(segments, buf)
				buf = []l.Token{}
			} else {
				buf = append(buf, tok)
			}
		default:
			buf = append(buf, tok)
		}
	}
	segments = append(segments, buf)
	if dropTrailingEmpty && len(segments) > 0 && len(segments[len(segments)-1]) == 0 {
		segments = segments[:len(segments)-1]
	}
	return segments
}

func topLevelIndex(tokens []l.Token, tokenType l.TokenType) int {
	depthParen := 0
	depthBracket := 0
	for i, tok := range tokens {
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
		if depthParen == 0 && depthBracket == 0 && tok.Type == tokenType {
			return i
		}
	}
	return -1
}

func hasTopLevelToken(tokens []l.Token, tokenType l.TokenType) bool {
	return topLevelIndex(tokens, tokenType) >= 0
}
