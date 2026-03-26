package diagnostics_test

import (
	"testing"

	d "github.com/maxvanasten/gscp/diagnostics"
	l "github.com/maxvanasten/gscp/lexer"
	p "github.com/maxvanasten/gscp/parser"
	"github.com/stretchr/testify/assert"
)

func tokAt(tokenType l.TokenType, content string, col int) l.Token {
	endCol := col
	if len(content) > 0 {
		endCol = col + len(content) - 1
	}
	return l.Token{Type: tokenType, Content: content, Line: 1, Col: col, EndLine: 1, EndCol: endCol}
}

func assertHasDiagnostic(t *testing.T, diags []d.Diagnostic, expected d.Diagnostic) {
	t.Helper()
	for _, diag := range diags {
		if diag == expected {
			return
		}
	}
	assert.Failf(t, "expected diagnostic", "expected %+v, got %+v", expected, diags)
}

func tokenContents(tokens []l.Token) []string {
	contents := []string{}
	for _, tok := range tokens {
		contents = append(contents, tok.Content)
	}
	return contents
}

func TestDiagnosticsLexerUnterminatedString(t *testing.T) {
	input := []byte("\"hello")
	lexer := l.NewLexer(input)
	actual := lexer.GetDiagnostics()

	expected := d.New("unterminated string literal", 1, 1, 1, 1, "error")
	assertHasDiagnostic(t, actual, expected)
}

func TestDiagnosticsLexerInvalidToken(t *testing.T) {
	input := []byte("1a")
	lexer := l.NewLexer(input)
	actual := lexer.GetDiagnostics()

	expected := d.New("invalid token", 1, 1, 1, 2, "error")
	assertHasDiagnostic(t, actual, expected)
}

func TestDiagnosticsLexerUnterminatedBlockComment(t *testing.T) {
	input := []byte("/# unterminated")
	lexer := l.NewLexer(input)
	actual := lexer.GetDiagnostics()

	expected := d.New("unterminated block comment", 1, 1, 1, 1, "error")
	assertHasDiagnostic(t, actual, expected)
}

func TestDiagnosticsLexerSkipsLineComment(t *testing.T) {
	input := []byte("a = 1; // comment\nb = 2;")
	lexer := l.NewLexer(input)
	assert.Len(t, lexer.GetDiagnostics(), 0)
	assert.Contains(t, tokenContents(lexer.GetTokens()), "a")
	assert.Contains(t, tokenContents(lexer.GetTokens()), "b")
}

func TestDiagnosticsLexerSkipsBlockComment(t *testing.T) {
	input := []byte("/# block\ncomment #/\na = 1;")
	lexer := l.NewLexer(input)
	assert.Len(t, lexer.GetDiagnostics(), 0)
	assert.Contains(t, tokenContents(lexer.GetTokens()), "a")
}

func TestDiagnosticsLexerSymbolStarts(t *testing.T) {
	input := []byte("::init_sidequest points[i].target _private")
	lexer := l.NewLexer(input)
	assert.Len(t, lexer.GetDiagnostics(), 0)
	// :: is now tokenized separately as FUNCTION_POINTER, followed by the function name
	assert.Contains(t, tokenContents(lexer.GetTokens()), "::")
	assert.Contains(t, tokenContents(lexer.GetTokens()), "init_sidequest")
	assert.Contains(t, tokenContents(lexer.GetTokens()), "points")
	assert.Contains(t, tokenContents(lexer.GetTokens()), ".target")
	assert.Contains(t, tokenContents(lexer.GetTokens()), "_private")
}

func TestDiagnosticsWaitFunctionCall(t *testing.T) {
	lexer := l.NewLexer([]byte("wait( 0.05 );"))
	_, diags := p.Parse(lexer.GetTokens())
	assert.Len(t, diags, 0)
}

func TestDiagnosticsIncrementOperator(t *testing.T) {
	lexer := l.NewLexer([]byte("x++;"))
	_, diags := p.Parse(lexer.GetTokens())
	assert.Len(t, diags, 0)
}

func TestDiagnosticsNestedBlockComment(t *testing.T) {
	input := []byte("/# outer /# inner #/ still outer #/ x = 1;")
	lexer := l.NewLexer(input)
	_, diags := p.Parse(lexer.GetTokens())
	assert.Len(t, diags, 0)
}

func TestDiagnosticsDoubleNegation(t *testing.T) {
	input := []byte("if ( !!isdefined( x ) ) { }")
	lexer := l.NewLexer(input)
	_, diags := p.Parse(lexer.GetTokens())
	assert.Len(t, diags, 0)
}

func TestDiagnosticsPercentPrefix(t *testing.T) {
	input := []byte("anim = %o_riot_stand_deploy;")
	lexer := l.NewLexer(input)
	_, diags := p.Parse(lexer.GetTokens())
	assert.Len(t, diags, 0)
}

func TestDiagnosticsMissingIncludePath(t *testing.T) {
	input := []l.Token{
		tokAt(l.SYMBOL, "#include", 1),
	}

	_, diags := p.Parse(input)
	expected := d.New("missing include path", 1, 1, 1, 8, "error")
	assertHasDiagnostic(t, diags, expected)
}

func TestDiagnosticsMissingWaitDuration(t *testing.T) {
	input := []l.Token{
		tokAt(l.SYMBOL, "wait", 5),
	}

	_, diags := p.Parse(input)
	expected := d.New("missing wait duration", 1, 5, 1, 8, "error")
	assertHasDiagnostic(t, diags, expected)
}

func TestDiagnosticsMissingUnaryOperand(t *testing.T) {
	input := []l.Token{
		tokAt(l.OPERATOR, "!", 3),
	}

	_, diags := p.Parse(input)
	expected := d.New("missing unary operand", 1, 3, 1, 3, "error")
	assertHasDiagnostic(t, diags, expected)
}

func TestDiagnosticsOperatorMissingLeftOperand(t *testing.T) {
	input := []l.Token{
		tokAt(l.OPERATOR, "+", 2),
		tokAt(l.NUMBER, "1", 4),
	}

	_, diags := p.Parse(input)
	expected := d.New("operator missing left-hand operand", 1, 2, 1, 2, "error")
	assertHasDiagnostic(t, diags, expected)
}

func TestDiagnosticsOperatorMissingRightOperand(t *testing.T) {
	input := []l.Token{
		tokAt(l.NUMBER, "1", 1),
		tokAt(l.OPERATOR, "+", 3),
		tokAt(l.TERMINATOR, ";", 4),
	}

	_, diags := p.Parse(input)
	expected := d.New("operator missing right-hand operand", 1, 3, 1, 3, "error")
	assertHasDiagnostic(t, diags, expected)
}

func TestDiagnosticsAssignmentMissingLeftSide(t *testing.T) {
	input := []l.Token{
		tokAt(l.ASSIGNMENT, "=", 6),
	}

	_, diags := p.Parse(input)
	expected := d.New("assignment missing left-hand side", 1, 6, 1, 6, "error")
	assertHasDiagnostic(t, diags, expected)
}

func TestDiagnosticsAssignmentTargetMustBeVariable(t *testing.T) {
	input := []l.Token{
		tokAt(l.NUMBER, "1", 1),
		tokAt(l.ASSIGNMENT, "=", 3),
		tokAt(l.NUMBER, "2", 5),
	}

	_, diags := p.Parse(input)
	expected := d.New("assignment target must be a variable", 1, 3, 1, 3, "error")
	assertHasDiagnostic(t, diags, expected)
}

func TestDiagnosticsMissingClosingParen(t *testing.T) {
	input := []l.Token{
		tokAt(l.SYMBOL, "if", 1),
		tokAt(l.OPEN_PAREN, "(", 4),
		tokAt(l.SYMBOL, "x", 5),
	}

	_, diags := p.Parse(input)
	expected := d.New("missing closing )", 1, 4, 1, 4, "error")
	assertHasDiagnostic(t, diags, expected)
}

func TestDiagnosticsMissingClosingBracket(t *testing.T) {
	input := []l.Token{
		tokAt(l.SYMBOL, "arr", 1),
		tokAt(l.OPEN_BRACKET, "[", 4),
		tokAt(l.NUMBER, "1", 5),
	}

	_, diags := p.Parse(input)
	expected := d.New("missing closing ]", 1, 4, 1, 4, "error")
	assertHasDiagnostic(t, diags, expected)
}

func TestDiagnosticsMissingClosingCurly(t *testing.T) {
	input := []l.Token{
		tokAt(l.SYMBOL, "if", 1),
		tokAt(l.OPEN_PAREN, "(", 4),
		tokAt(l.SYMBOL, "x", 5),
		tokAt(l.CLOSE_PAREN, ")", 6),
		tokAt(l.OPEN_CURLY, "{", 8),
	}

	_, diags := p.Parse(input)
	expected := d.New("missing closing }", 1, 8, 1, 8, "error")
	assertHasDiagnostic(t, diags, expected)
}

func TestDiagnosticsUnexpectedOpenCurly(t *testing.T) {
	input := []l.Token{
		tokAt(l.OPEN_CURLY, "{", 1),
	}

	_, diags := p.Parse(input)
	// Parser now handles bare blocks gracefully for decompiled code compatibility
	// Expect "missing closing }" instead of "unexpected {"
	expected := d.New("missing closing }", 1, 1, 1, 1, "error")
	assertHasDiagnostic(t, diags, expected)
}

func TestDiagnosticsUnexpectedCloseParen(t *testing.T) {
	input := []l.Token{
		tokAt(l.CLOSE_PAREN, ")", 2),
	}

	_, diags := p.Parse(input)
	expected := d.New("unexpected )", 1, 2, 1, 2, "error")
	assertHasDiagnostic(t, diags, expected)
}

func TestDiagnosticsUnexpectedCloseBracket(t *testing.T) {
	input := []l.Token{
		tokAt(l.CLOSE_BRACKET, "]", 3),
	}

	_, diags := p.Parse(input)
	expected := d.New("unexpected ]", 1, 3, 1, 3, "error")
	assertHasDiagnostic(t, diags, expected)
}

func TestDiagnosticsUnexpectedCloseCurly(t *testing.T) {
	input := []l.Token{
		tokAt(l.CLOSE_CURLY, "}", 4),
	}

	_, diags := p.Parse(input)
	// Parser now handles stray closing braces gracefully for decompiled code compatibility
	// No error should be reported
	if len(diags) != 0 {
		t.Errorf("Expected no diagnostics for stray closing brace, got %v", diags)
	}
}
