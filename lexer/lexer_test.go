package lexer_test

import (
	"testing"

	"github.com/maxvanasten/gscp/lexer"
	"github.com/stretchr/testify/assert"
)

func tok(tokenType lexer.TokenType, content string) lexer.Token {
	return lexer.Token{Type: tokenType, Content: content}
}

func assertTokens(t *testing.T, input []byte, targets []lexer.Token) {
	t.Helper()

	l := lexer.NewLexer(input)
	tokens := l.GetTokens()

	actual := []lexer.Token{}
	for _, tok := range tokens {
		actual = append(actual, lexer.Token{Type: tok.Type, Content: tok.Content})
	}
	assert.Equal(t, targets, actual)
}

func TestLexerSymbol(t *testing.T) {
	input := []byte("alpha")

	targets := []lexer.Token{
		tok(lexer.SYMBOL, "alpha"),
	}

	assertTokens(t, input, targets)
}

func TestLexerNumber(t *testing.T) {
	input := []byte("23.5")

	targets := []lexer.Token{
		tok(lexer.NUMBER, "23.5"),
	}

	assertTokens(t, input, targets)
}

func TestLexerString(t *testing.T) {
	input := []byte("\"hello\"")

	targets := []lexer.Token{
		tok(lexer.STRING, "hello"),
	}

	assertTokens(t, input, targets)
}

func TestLexerTerminator(t *testing.T) {
	input := []byte(";")

	targets := []lexer.Token{
		tok(lexer.TERMINATOR, ";"),
	}

	assertTokens(t, input, targets)
}

func TestLexerComma(t *testing.T) {
	input := []byte(",")

	targets := []lexer.Token{
		tok(lexer.COMMA, ","),
	}

	assertTokens(t, input, targets)
}

func TestLexerNewline(t *testing.T) {
	input := []byte("\n")

	targets := []lexer.Token{
		tok(lexer.NEWLINE, ""),
	}

	assertTokens(t, input, targets)
}

func TestLexerOpenParen(t *testing.T) {
	input := []byte("(")

	targets := []lexer.Token{
		tok(lexer.OPEN_PAREN, "("),
	}

	assertTokens(t, input, targets)
}

func TestLexerCloseParen(t *testing.T) {
	input := []byte(")")

	targets := []lexer.Token{
		tok(lexer.CLOSE_PAREN, ")"),
	}

	assertTokens(t, input, targets)
}

func TestLexerOpenBracket(t *testing.T) {
	input := []byte("[")

	targets := []lexer.Token{
		tok(lexer.OPEN_BRACKET, "["),
	}

	assertTokens(t, input, targets)
}

func TestLexerCloseBracket(t *testing.T) {
	input := []byte("]")

	targets := []lexer.Token{
		tok(lexer.CLOSE_BRACKET, "]"),
	}

	assertTokens(t, input, targets)
}

func TestLexerOpenCurly(t *testing.T) {
	input := []byte("{")

	targets := []lexer.Token{
		tok(lexer.OPEN_CURLY, "{"),
	}

	assertTokens(t, input, targets)
}

func TestLexerCloseCurly(t *testing.T) {
	input := []byte("}")

	targets := []lexer.Token{
		tok(lexer.CLOSE_CURLY, "}"),
	}

	assertTokens(t, input, targets)
}

func TestLexerAssignment(t *testing.T) {
	input := []byte("=")

	targets := []lexer.Token{
		tok(lexer.ASSIGNMENT, "="),
	}

	assertTokens(t, input, targets)
}

func TestLexerCompoundAssignment(t *testing.T) {
	input := []byte("a += 1; a -= 1; a *= 1; a /= 1;")

	targets := []lexer.Token{
		tok(lexer.SYMBOL, "a"),
		tok(lexer.ASSIGNMENT, "+="),
		tok(lexer.NUMBER, "1"),
		tok(lexer.TERMINATOR, ";"),
		tok(lexer.SYMBOL, "a"),
		tok(lexer.ASSIGNMENT, "-="),
		tok(lexer.NUMBER, "1"),
		tok(lexer.TERMINATOR, ";"),
		tok(lexer.SYMBOL, "a"),
		tok(lexer.ASSIGNMENT, "*="),
		tok(lexer.NUMBER, "1"),
		tok(lexer.TERMINATOR, ";"),
		tok(lexer.SYMBOL, "a"),
		tok(lexer.ASSIGNMENT, "/="),
		tok(lexer.NUMBER, "1"),
		tok(lexer.TERMINATOR, ";"),
	}

	assertTokens(t, input, targets)
}

func TestLexerArithmeticOperators(t *testing.T) {
	input := []byte("a + b - c * d / e;")

	targets := []lexer.Token{
		tok(lexer.SYMBOL, "a"),
		tok(lexer.OPERATOR, "+"),
		tok(lexer.SYMBOL, "b"),
		tok(lexer.OPERATOR, "-"),
		tok(lexer.SYMBOL, "c"),
		tok(lexer.OPERATOR, "*"),
		tok(lexer.SYMBOL, "d"),
		tok(lexer.OPERATOR, "/"),
		tok(lexer.SYMBOL, "e"),
		tok(lexer.TERMINATOR, ";"),
	}

	assertTokens(t, input, targets)
}

func TestLexerComparisonOperators(t *testing.T) {
	input := []byte("a < b; a > b; a <= b; a >= b; a == b; a != b; !a;")

	targets := []lexer.Token{
		tok(lexer.SYMBOL, "a"),
		tok(lexer.OPERATOR, "<"),
		tok(lexer.SYMBOL, "b"),
		tok(lexer.TERMINATOR, ";"),
		tok(lexer.SYMBOL, "a"),
		tok(lexer.OPERATOR, ">"),
		tok(lexer.SYMBOL, "b"),
		tok(lexer.TERMINATOR, ";"),
		tok(lexer.SYMBOL, "a"),
		tok(lexer.OPERATOR, "<="),
		tok(lexer.SYMBOL, "b"),
		tok(lexer.TERMINATOR, ";"),
		tok(lexer.SYMBOL, "a"),
		tok(lexer.OPERATOR, ">="),
		tok(lexer.SYMBOL, "b"),
		tok(lexer.TERMINATOR, ";"),
		tok(lexer.SYMBOL, "a"),
		tok(lexer.OPERATOR, "=="),
		tok(lexer.SYMBOL, "b"),
		tok(lexer.TERMINATOR, ";"),
		tok(lexer.SYMBOL, "a"),
		tok(lexer.OPERATOR, "!="),
		tok(lexer.SYMBOL, "b"),
		tok(lexer.TERMINATOR, ";"),
		tok(lexer.OPERATOR, "!"),
		tok(lexer.SYMBOL, "a"),
		tok(lexer.TERMINATOR, ";"),
	}

	assertTokens(t, input, targets)
}

func TestLexerLogicalOperators(t *testing.T) {
	input := []byte("a && b || c;")

	targets := []lexer.Token{
		tok(lexer.SYMBOL, "a"),
		tok(lexer.OPERATOR, "&&"),
		tok(lexer.SYMBOL, "b"),
		tok(lexer.OPERATOR, "||"),
		tok(lexer.SYMBOL, "c"),
		tok(lexer.TERMINATOR, ";"),
	}

	assertTokens(t, input, targets)
}
