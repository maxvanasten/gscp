package lexer_test

import (
	"testing"

	"github.com/maxvanasten/gscp/lexer"
	"github.com/stretchr/testify/assert"
)

func assertTokens(t *testing.T, input []byte, targets []lexer.Token) {
	t.Helper()

	l := lexer.NewLexer(input)
	tokens := l.GetTokens()

	assert.Equal(t, targets, tokens)
}

func TestLexerSymbol(t *testing.T) {
	input := []byte("alpha")

	targets := []lexer.Token{
		{lexer.SYMBOL, "alpha"},
	}

	assertTokens(t, input, targets)
}

func TestLexerNumber(t *testing.T) {
	input := []byte("23.5")

	targets := []lexer.Token{
		{lexer.NUMBER, "23.5"},
	}

	assertTokens(t, input, targets)
}

func TestLexerString(t *testing.T) {
	input := []byte("\"hello\"")

	targets := []lexer.Token{
		{lexer.STRING, "hello"},
	}

	assertTokens(t, input, targets)
}

func TestLexerTerminator(t *testing.T) {
	input := []byte(";")

	targets := []lexer.Token{
		{lexer.TERMINATOR, ";"},
	}

	assertTokens(t, input, targets)
}

func TestLexerComma(t *testing.T) {
	input := []byte(",")

	targets := []lexer.Token{
		{lexer.COMMA, ","},
	}

	assertTokens(t, input, targets)
}

func TestLexerNewline(t *testing.T) {
	input := []byte("\n")

	targets := []lexer.Token{
		{lexer.NEWLINE, ""},
	}

	assertTokens(t, input, targets)
}

func TestLexerOpenParen(t *testing.T) {
	input := []byte("(")

	targets := []lexer.Token{
		{lexer.OPEN_PAREN, "("},
	}

	assertTokens(t, input, targets)
}

func TestLexerCloseParen(t *testing.T) {
	input := []byte(")")

	targets := []lexer.Token{
		{lexer.CLOSE_PAREN, ")"},
	}

	assertTokens(t, input, targets)
}

func TestLexerOpenBracket(t *testing.T) {
	input := []byte("[")

	targets := []lexer.Token{
		{lexer.OPEN_BRACKET, "["},
	}

	assertTokens(t, input, targets)
}

func TestLexerCloseBracket(t *testing.T) {
	input := []byte("]")

	targets := []lexer.Token{
		{lexer.CLOSE_BRACKET, "]"},
	}

	assertTokens(t, input, targets)
}

func TestLexerOpenCurly(t *testing.T) {
	input := []byte("{")

	targets := []lexer.Token{
		{lexer.OPEN_CURLY, "{"},
	}

	assertTokens(t, input, targets)
}

func TestLexerCloseCurly(t *testing.T) {
	input := []byte("}")

	targets := []lexer.Token{
		{lexer.CLOSE_CURLY, "}"},
	}

	assertTokens(t, input, targets)
}

func TestLexerAssignment(t *testing.T) {
	input := []byte("=")

	targets := []lexer.Token{
		{lexer.ASSIGNMENT, "="},
	}

	assertTokens(t, input, targets)
}

func TestLexerCompoundAssignment(t *testing.T) {
	input := []byte("a += 1; a -= 1; a *= 1; a /= 1;")

	targets := []lexer.Token{
		{lexer.SYMBOL, "a"},
		{lexer.ASSIGNMENT, "+="},
		{lexer.NUMBER, "1"},
		{lexer.TERMINATOR, ";"},
		{lexer.SYMBOL, "a"},
		{lexer.ASSIGNMENT, "-="},
		{lexer.NUMBER, "1"},
		{lexer.TERMINATOR, ";"},
		{lexer.SYMBOL, "a"},
		{lexer.ASSIGNMENT, "*="},
		{lexer.NUMBER, "1"},
		{lexer.TERMINATOR, ";"},
		{lexer.SYMBOL, "a"},
		{lexer.ASSIGNMENT, "/="},
		{lexer.NUMBER, "1"},
		{lexer.TERMINATOR, ";"},
	}

	assertTokens(t, input, targets)
}

func TestLexerArithmeticOperators(t *testing.T) {
	input := []byte("a + b - c * d / e;")

	targets := []lexer.Token{
		{lexer.SYMBOL, "a"},
		{lexer.OPERATOR, "+"},
		{lexer.SYMBOL, "b"},
		{lexer.OPERATOR, "-"},
		{lexer.SYMBOL, "c"},
		{lexer.OPERATOR, "*"},
		{lexer.SYMBOL, "d"},
		{lexer.OPERATOR, "/"},
		{lexer.SYMBOL, "e"},
		{lexer.TERMINATOR, ";"},
	}

	assertTokens(t, input, targets)
}

func TestLexerComparisonOperators(t *testing.T) {
	input := []byte("a < b; a > b; a <= b; a >= b; a == b; a != b; !a;")

	targets := []lexer.Token{
		{lexer.SYMBOL, "a"},
		{lexer.OPERATOR, "<"},
		{lexer.SYMBOL, "b"},
		{lexer.TERMINATOR, ";"},
		{lexer.SYMBOL, "a"},
		{lexer.OPERATOR, ">"},
		{lexer.SYMBOL, "b"},
		{lexer.TERMINATOR, ";"},
		{lexer.SYMBOL, "a"},
		{lexer.OPERATOR, "<="},
		{lexer.SYMBOL, "b"},
		{lexer.TERMINATOR, ";"},
		{lexer.SYMBOL, "a"},
		{lexer.OPERATOR, ">="},
		{lexer.SYMBOL, "b"},
		{lexer.TERMINATOR, ";"},
		{lexer.SYMBOL, "a"},
		{lexer.OPERATOR, "=="},
		{lexer.SYMBOL, "b"},
		{lexer.TERMINATOR, ";"},
		{lexer.SYMBOL, "a"},
		{lexer.OPERATOR, "!="},
		{lexer.SYMBOL, "b"},
		{lexer.TERMINATOR, ";"},
		{lexer.OPERATOR, "!"},
		{lexer.SYMBOL, "a"},
		{lexer.TERMINATOR, ";"},
	}

	assertTokens(t, input, targets)
}

func TestLexerLogicalOperators(t *testing.T) {
	input := []byte("a && b || c;")

	targets := []lexer.Token{
		{lexer.SYMBOL, "a"},
		{lexer.OPERATOR, "&&"},
		{lexer.SYMBOL, "b"},
		{lexer.OPERATOR, "||"},
		{lexer.SYMBOL, "c"},
		{lexer.TERMINATOR, ";"},
	}

	assertTokens(t, input, targets)
}
