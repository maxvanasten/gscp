package lexer

import (
	"bytes"
	"slices"
	"strings"
	"unicode"

	d "github.com/maxvanasten/gscp/diagnostics"
)

type TokenType int

const (
	EOF TokenType = iota

	SYMBOL
	NUMBER
	STRING

	TERMINATOR
	COMMA
	COLON
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
	case COLON:
		return "colon"
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
	Line    int
	Col     int
	EndLine int
	EndCol  int
}

type Lexer struct {
	input       []byte
	buffer      []byte
	index       int
	line        int
	col         int
	bufferLine  int
	bufferCol   int
	tokens      []Token
	diagnostics []d.Diagnostic
}

func isSymbolStart(buffer []byte) bool {
	if len(buffer) == 0 {
		return false
	}
	first := buffer[0]
	if unicode.IsLetter(rune(first)) || first == '#' || first == '_' {
		return true
	}
	if first == '.' {
		if len(buffer) > 1 {
			second := buffer[1]
			return unicode.IsLetter(rune(second)) || second == '_' || second == '#'
		}
		return false
	}
	if first == ':' && len(buffer) > 2 && buffer[1] == ':' {
		third := buffer[2]
		return unicode.IsLetter(rune(third)) || third == '_'
	}
	return false
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
	case ':':
		return COLON
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

		startLine := l.bufferLine
		startCol := l.bufferCol
		endLine := l.bufferLine
		endCol := l.bufferCol + len(l.buffer) - 1

		if is_number {
			l.tokens = append(l.tokens, Token{Type: NUMBER, Content: string(l.buffer), Line: startLine, Col: startCol, EndLine: endLine, EndCol: endCol})
		} else {
			// Check if symbol is valid
			if isSymbolStart(l.buffer) {
				l.tokens = append(l.tokens, Token{Type: SYMBOL, Content: string(l.buffer), Line: startLine, Col: startCol, EndLine: endLine, EndCol: endCol})
			} else {
				l.diagnostics = append(l.diagnostics, d.New("invalid token", startLine, startCol, endLine, endCol, "error"))
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
	startLine := l.line
	startCol := l.col
	switch c {
	case '+', '-', '*', '/':
		if c == '/' && l.index+1 < len(l.input) {
			next := l.input[l.index+1]
			if next == '/' {
				l.HandleBuffer()
				commentEnd := bytes.IndexByte(l.input[l.index+2:], '\n')
				if commentEnd < 0 {
					return len(l.input) - l.index
				}
				return commentEnd + 2
			}
			if next == '#' {
				l.HandleBuffer()
				blockEnd := bytes.Index(l.input[l.index+2:], []byte("#/"))
				if blockEnd < 0 {
					startLine := l.line
					startCol := l.col
					endLine := l.line
					endCol := l.col
					l.diagnostics = append(l.diagnostics, d.New("unterminated block comment", startLine, startCol, endLine, endCol, "error"))
					return len(l.input) - l.index
				}
				return blockEnd + 4
			}
		}
		if (c == '+' || c == '-') && l.index+1 < len(l.input) && l.input[l.index+1] == c {
			l.HandleBuffer()
			l.tokens = append(l.tokens, Token{Type: OPERATOR, Content: string([]byte{c, c}), Line: startLine, Col: startCol, EndLine: startLine, EndCol: startCol + 1})
			return 2
		}
		if l.index+1 < len(l.input) && l.input[l.index+1] == '=' {
			l.HandleBuffer()
			l.tokens = append(l.tokens, Token{Type: ASSIGNMENT, Content: string([]byte{c, '='}), Line: startLine, Col: startCol, EndLine: startLine, EndCol: startCol + 1})
			return 2
		}
		l.HandleBuffer()
		l.tokens = append(l.tokens, Token{Type: OPERATOR, Content: strings.TrimSpace(string(c)), Line: startLine, Col: startCol, EndLine: startLine, EndCol: startCol})
		return 1
	case '<', '>':
		if l.index+1 < len(l.input) && l.input[l.index+1] == '=' {
			l.HandleBuffer()
			l.tokens = append(l.tokens, Token{Type: OPERATOR, Content: string([]byte{c, '='}), Line: startLine, Col: startCol, EndLine: startLine, EndCol: startCol + 1})
			return 2
		}
		l.HandleBuffer()
		l.tokens = append(l.tokens, Token{Type: OPERATOR, Content: strings.TrimSpace(string(c)), Line: startLine, Col: startCol, EndLine: startLine, EndCol: startCol})
		return 1
	case '&', '|':
		if l.index+1 < len(l.input) && l.input[l.index+1] == c {
			l.HandleBuffer()
			l.tokens = append(l.tokens, Token{Type: OPERATOR, Content: string([]byte{c, c}), Line: startLine, Col: startCol, EndLine: startLine, EndCol: startCol + 1})
			return 2
		}
		l.HandleBuffer()
		l.tokens = append(l.tokens, Token{Type: OPERATOR, Content: strings.TrimSpace(string(c)), Line: startLine, Col: startCol, EndLine: startLine, EndCol: startCol})
		return 1
	case '=':
		if l.index+1 < len(l.input) && l.input[l.index+1] == '=' {
			l.HandleBuffer()
			l.tokens = append(l.tokens, Token{Type: OPERATOR, Content: "==", Line: startLine, Col: startCol, EndLine: startLine, EndCol: startCol + 1})
			return 2
		}
		l.HandleBuffer()
		l.tokens = append(l.tokens, Token{Type: ASSIGNMENT, Content: "=", Line: startLine, Col: startCol, EndLine: startLine, EndCol: startCol})
		return 1
	case '!':
		if l.index+1 < len(l.input) && l.input[l.index+1] == '=' {
			l.HandleBuffer()
			l.tokens = append(l.tokens, Token{Type: OPERATOR, Content: "!=", Line: startLine, Col: startCol, EndLine: startLine, EndCol: startCol + 1})
			return 2
		}
		l.HandleBuffer()
		l.tokens = append(l.tokens, Token{Type: OPERATOR, Content: "!", Line: startLine, Col: startCol, EndLine: startLine, EndCol: startCol})
		return 1
	case '\n', '(', '[', '{', ')', ']', '}', ';', ',':
		l.HandleBuffer()
		token_type := TokenTypeFromChar(c)
		l.tokens = append(l.tokens, Token{Type: token_type, Content: strings.TrimSpace(string(c)), Line: startLine, Col: startCol, EndLine: startLine, EndCol: startCol})
		return 1
	case ':':
		if l.index+1 < len(l.input) && l.input[l.index+1] == ':' {
			if len(l.buffer) == 0 {
				l.bufferLine = l.line
				l.bufferCol = l.col
			}
			l.buffer = append(l.buffer, ':', ':')
			return 2
		}
		l.HandleBuffer()
		l.tokens = append(l.tokens, Token{Type: COLON, Content: ":", Line: startLine, Col: startCol, EndLine: startLine, EndCol: startCol})
		return 1
	case ' ', '\t':
		l.HandleBuffer()
		return 1
	case '"':
		l.HandleBuffer()

		nextIndex := bytes.Index(l.input[l.index+1:], []byte{'"'})
		if nextIndex < 0 {
			endLine := l.line
			endCol := l.col
			if len(l.input) > 0 {
				endLine = l.line
				endCol = l.col
			}
			l.diagnostics = append(l.diagnostics, d.New("unterminated string literal", startLine, startCol, endLine, endCol, "error"))
			return len(l.input) - l.index
		} else {
			string_content := l.input[l.index+1 : l.index+1+nextIndex]
			endCol := startCol + nextIndex + 1
			l.tokens = append(l.tokens, Token{Type: STRING, Content: string(string_content), Line: startLine, Col: startCol, EndLine: startLine, EndCol: endCol})
			return len(string_content) + 2
		}
	default:
		if len(l.buffer) == 0 {
			l.bufferLine = l.line
			l.bufferCol = l.col
		}
		l.buffer = append(l.buffer, c)
		return 1
	}
	return 1
}

func (l *Lexer) Next() {
	if !l.EOF() {
		consumed := l.HandleCharacter(l.input[l.index])
		if l.index+consumed > len(l.input) {
			consumed = len(l.input) - l.index
		}
		l.advancePosition(l.input[l.index : l.index+consumed])
		l.index += consumed
	}
}

func (l *Lexer) GetTokens() []Token {
	return l.tokens
}

func (l *Lexer) GetDiagnostics() []d.Diagnostic {
	return l.diagnostics
}

func (l *Lexer) advancePosition(bytes []byte) {
	for _, c := range bytes {
		if c == '\n' {
			l.line++
			l.col = 1
			continue
		}
		l.col++
	}
}

func NewLexer(input []byte) Lexer {
	l := Lexer{input: input, line: 1, col: 1}
	for !l.EOF() {
		l.Next()
	}
	return l
}
