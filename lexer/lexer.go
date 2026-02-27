package lexer

import (
	"bytes"
	"slices"
	"strings"
	"unicode"
)

type TokenType int

const (
	EOF TokenType = iota

	SYMBOL
	NUMBER
	STRING

	TERMINATOR
	COMMA
	NEWLINE

	OPEN_PAREN
	CLOSE_PAREN
	OPEN_BRACKET
	CLOSE_BRACKET
	OPEN_CURLY
	CLOSE_CURLY

	ASSIGNMENT
	OPERATOR
)

func (t TokenType) ToString() string {
	switch t {
	case EOF:
		return "EOF"
	case SYMBOL:
		return "symbol"
	case NUMBER:
		return "number"
	case STRING:
		return "string"
	case TERMINATOR:
		return "terminator"
	case COMMA:
		return "comma"
	case NEWLINE:
		return "newline"
	case OPEN_PAREN:
		return "open_paren"
	case CLOSE_PAREN:
		return "close_paren"
	case OPEN_BRACKET:
		return "open_bracket"
	case CLOSE_BRACKET:
		return "close_bracket"
	case OPEN_CURLY:
		return "open_curly"
	case CLOSE_CURLY:
		return "close_curly"
	case ASSIGNMENT:
		return "assignment"
	case OPERATOR:
		return "operator"
	default:
		return ""
	}
}

type Token struct {
	Type    TokenType
	Content string
}

type Lexer struct {
	input  []byte
	buffer []byte
	index  int
	tokens []Token
}

func TokensFromTypes(types []TokenType) []Token {
	output := []Token{}
	for _, t := range types {
		output = append(output, Token{Type: t, Content: ""})
	}
	return output
}

func TokenTypeFromChar(char byte) TokenType {
	switch char {
	case '(':
		return OPEN_PAREN
	case ')':
		return CLOSE_PAREN
	case '[':
		return OPEN_BRACKET
	case ']':
		return CLOSE_BRACKET
	case '{':
		return OPEN_CURLY
	case '}':
		return CLOSE_CURLY
	case '+', '-', '*', '/':
		return OPERATOR
	case '=':
		return ASSIGNMENT
	case ';':
		return TERMINATOR
	case '\n':
		return NEWLINE
	case ',':
		return COMMA
	case '<', '>':
		return OPERATOR
	default:
		return EOF
	}
}

func (l *Lexer) EOF() bool {
	eof := l.index >= len(l.input)
	if eof {
		l.HandleBuffer()
	}
	return eof
}

func (l *Lexer) HandleBuffer() {
	l.buffer = bytes.TrimSpace(l.buffer)
	if len(l.buffer) > 0 {
		// Todo: Check if number
		is_number := true
		has_decimal := false
		for _, c := range l.buffer {
			if !unicode.IsDigit(rune(c)) {
				if c == '.' && !has_decimal {
					has_decimal = true
				} else {
					is_number = false
				}
			}
		}

		if is_number {
			l.tokens = append(l.tokens, Token{Type: NUMBER, Content: string(l.buffer)})
		} else {
			// Check if symbol is valid
			if unicode.IsLetter(rune(l.buffer[0])) || l.buffer[0] == '#' {
				l.tokens = append(l.tokens, Token{Type: SYMBOL, Content: string(l.buffer)})
			}
		}
	}
	l.buffer = []byte{}
}

func TokensUntilAny(tokens []Token, targets []TokenType) []Token {
	token_buffer := []Token{}

	for _, t := range tokens {
		token_buffer = append(token_buffer, t)
		if slices.Contains(targets, t.Type) {
			return token_buffer
		}
	}

	return token_buffer
}

func (l *Lexer) HandleCharacter(c byte) int {
	switch c {
	case '+', '-', '*', '/':
		if l.index+1 < len(l.input) && l.input[l.index+1] == '=' {
			l.HandleBuffer()
			l.tokens = append(l.tokens, Token{Type: ASSIGNMENT, Content: string([]byte{c, '='})})
			return 2
		}
		l.HandleBuffer()
		l.tokens = append(l.tokens, Token{Type: OPERATOR, Content: strings.TrimSpace(string(c))})
		return 1
	case '<', '>':
		if l.index+1 < len(l.input) && l.input[l.index+1] == '=' {
			l.HandleBuffer()
			l.tokens = append(l.tokens, Token{Type: OPERATOR, Content: string([]byte{c, '='})})
			return 2
		}
		l.HandleBuffer()
		l.tokens = append(l.tokens, Token{Type: OPERATOR, Content: strings.TrimSpace(string(c))})
		return 1
	case '&', '|':
		if l.index+1 < len(l.input) && l.input[l.index+1] == c {
			l.HandleBuffer()
			l.tokens = append(l.tokens, Token{Type: OPERATOR, Content: string([]byte{c, c})})
			return 2
		}
		l.HandleBuffer()
		l.tokens = append(l.tokens, Token{Type: OPERATOR, Content: strings.TrimSpace(string(c))})
		return 1
	case '=':
		if l.index+1 < len(l.input) && l.input[l.index+1] == '=' {
			l.HandleBuffer()
			l.tokens = append(l.tokens, Token{Type: OPERATOR, Content: "=="})
			return 2
		}
		l.HandleBuffer()
		l.tokens = append(l.tokens, Token{Type: ASSIGNMENT, Content: "="})
		return 1
	case '!':
		if l.index+1 < len(l.input) && l.input[l.index+1] == '=' {
			l.HandleBuffer()
			l.tokens = append(l.tokens, Token{Type: OPERATOR, Content: "!="})
			return 2
		}
		l.HandleBuffer()
		l.tokens = append(l.tokens, Token{Type: OPERATOR, Content: "!"})
		return 1
	case '\n', '(', '[', '{', ')', ']', '}', ';', ',':
		l.HandleBuffer()
		token_type := TokenTypeFromChar(c)
		l.tokens = append(l.tokens, Token{Type: token_type, Content: strings.TrimSpace(string(c))})
		return 1
	case ' ', '\t':
		l.HandleBuffer()
		return 1
	case '"':
		l.HandleBuffer()

		nextIndex := bytes.Index(l.input[l.index+1:], []byte{'"'})
		if nextIndex < 0 {
			// ERROR: unterminated string
		} else {
			string_content := l.input[l.index+1 : l.index+1+nextIndex]
			l.tokens = append(l.tokens, Token{Type: STRING, Content: string(string_content)})
			return len(string_content) + 2
		}
	default:
		l.buffer = append(l.buffer, c)
		return 1
	}
	return 1
}

func (l *Lexer) Next() {
	if !l.EOF() {
		l.index += l.HandleCharacter(l.input[l.index])
	}
}

func (l *Lexer) GetTokens() []Token {
	return l.tokens
}

func NewLexer(input []byte) Lexer {
	l := Lexer{input: input}
	for !l.EOF() {
		l.Next()
	}
	return l
}
