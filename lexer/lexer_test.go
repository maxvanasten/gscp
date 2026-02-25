package lexer_test

import (
	"github.com/maxvanasten/gscp/lexer"
	"testing"
)

func TestLexer(t *testing.T) {
	// =================
	input := []byte("#include path\\to\\file;\ninit(arg1, arg2) {\n\tname = \"Max\";\n\tage = 23.5;\n\tprint(\"Hello \" + name + \". You are \" + age + \" years old.\");\n}")

	targets := []lexer.Token{
		{lexer.SYMBOL, "#include"},
		{lexer.SYMBOL, "path\\to\\file"},
		{lexer.TERMINATOR, ";"},
		{lexer.NEWLINE, ""},
		{lexer.SYMBOL, "init"},
		{lexer.OPEN_PAREN, "("},
		{lexer.SYMBOL, "arg1"},
		{lexer.COMMA, ","},
		{lexer.SYMBOL, "arg2"},
		{lexer.CLOSE_PAREN, ")"},
		{lexer.OPEN_CURLY, "{"},
		{lexer.NEWLINE, ""},

		{lexer.SYMBOL, "name"},
		{lexer.ASSIGNMENT, "="},
		{lexer.STRING, "Max"},
		{lexer.TERMINATOR, ";"},
		{lexer.NEWLINE, ""},

		{lexer.SYMBOL, "age"},
		{lexer.ASSIGNMENT, "="},
		{lexer.NUMBER, "23.5"},
		{lexer.TERMINATOR, ";"},
		{lexer.NEWLINE, ""},

		{lexer.SYMBOL, "print"},
		{lexer.OPEN_PAREN, "("},
		{lexer.STRING, "Hello "},
		{lexer.OPERATOR, "+"},
		{lexer.SYMBOL, "name"},
		{lexer.OPERATOR, "+"},
		{lexer.STRING, ". You are "},
		{lexer.OPERATOR, "+"},
		{lexer.SYMBOL, "age"},
		{lexer.OPERATOR, "+"},
		{lexer.STRING, " years old."},
		{lexer.CLOSE_PAREN, ")"},
		{lexer.TERMINATOR, ";"},
		{lexer.NEWLINE, ""},

		{lexer.CLOSE_CURLY, "}"},
	}
	// =================

	l := lexer.NewLexer(input)

	tokens := l.GetTokens()
	for i, token := range tokens {
		t.Logf("[%v] Type: %v, Content: %v\n", i, token.Type.ToString(), token.Content)
	}

	if len(targets) != len(tokens) {
		t.Fatalf("len(targets)(%v) != len(tokens)(%v)\n", len(targets), len(tokens))
	}

	for i := range targets {
		if targets[i] != tokens[i] {
			t.Fatalf("targets[%v][%v](%v) != tokens[%v][%v](%v)\n", i, targets[i].Type.ToString(), targets[i].Content, i, tokens[i].Type.ToString(), tokens[i].Content)
		}
	}
}
