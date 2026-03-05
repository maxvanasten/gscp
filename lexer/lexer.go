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
	LINE_COMMENT
	BLOCK_COMMENT
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
	case LINE_COMMENT:
		return "line_comment"
	case BLOCK_COMMENT:
		return "block_comment"
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

var twoCharTokens = map[string]TokenType{
	"++": OPERATOR,
	"--": OPERATOR,
	"==": OPERATOR,
	"!=": OPERATOR,
	"<=": OPERATOR,
	">=": OPERATOR,
	"&&": OPERATOR,
	"||": OPERATOR,
	"!!": OPERATOR,
	"<<": OPERATOR,
	">>": OPERATOR,
	"+=": ASSIGNMENT,
	"-=": ASSIGNMENT,
	"*=": ASSIGNMENT,
	"/=": ASSIGNMENT,
	"%=": ASSIGNMENT,
	"&=": ASSIGNMENT,
	"|=": ASSIGNMENT,
	"?=": ASSIGNMENT,
	"^=": ASSIGNMENT,
	"~=": ASSIGNMENT,
}

type Token struct {
	Type        TokenType
	Content     string
	Line        int
	Col         int
	EndLine     int
	EndCol      int
	StartOffset int
	EndOffset   int
}

type Lexer struct {
	input       []byte
	buffer      []byte
	index       int
	line        int
	col         int
	bufferLine  int
	bufferCol   int
	bufferIndex int
	tokens      []Token
	diagnostics []d.Diagnostic
}

func (l *Lexer) emitToken(tokenType TokenType, content string, startLine int, startCol int, endLine int, endCol int, startOffset int, endOffset int) {
	l.tokens = append(l.tokens, Token{Type: tokenType, Content: content, Line: startLine, Col: startCol, EndLine: endLine, EndCol: endCol, StartOffset: startOffset, EndOffset: endOffset})
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
	case '+', '-', '*', '/', '^', '~':
		return OPERATOR
	case '%':
		return OPERATOR
	case '?':
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
		startOffset := l.bufferIndex
		endOffset := l.bufferIndex + len(l.buffer) - 1

		if is_number {
			l.tokens = append(l.tokens, Token{Type: NUMBER, Content: string(l.buffer), Line: startLine, Col: startCol, EndLine: endLine, EndCol: endCol, StartOffset: startOffset, EndOffset: endOffset})
		} else {
			// Check if symbol is valid
			if isSymbolStart(l.buffer) {
				l.tokens = append(l.tokens, Token{Type: SYMBOL, Content: string(l.buffer), Line: startLine, Col: startCol, EndLine: endLine, EndCol: endCol, StartOffset: startOffset, EndOffset: endOffset})
			} else {
				l.diagnostics = append(l.diagnostics, d.New("invalid token", startLine, startCol, endLine, endCol, "error"))
			}
		}
	}
	l.buffer = []byte{}
}

func endPositionForConsumed(startLine int, startCol int, consumed []byte) (int, int) {
	if len(consumed) == 0 {
		return startLine, startCol
	}

	line := startLine
	col := startCol
	for i, c := range consumed {
		if i == len(consumed)-1 {
			return line, col
		}
		if c == '\n' {
			line++
			col = 1
			continue
		}
		col++
	}

	return line, col
}

func (l *Lexer) emitCommentToken(tokenType TokenType, startLine int, startCol int, consumed []byte) {
	if len(consumed) == 0 {
		return
	}
	endLine, endCol := endPositionForConsumed(startLine, startCol, consumed)
	l.emitToken(tokenType, string(consumed), startLine, startCol, endLine, endCol, l.index, l.index+len(consumed)-1)
}

func (l *Lexer) consumeLineComment(startLine int, startCol int) int {
	commentEnd := bytes.IndexByte(l.input[l.index:], '\n')
	if commentEnd < 0 {
		commentEnd = len(l.input) - l.index
	}
	consumed := l.input[l.index : l.index+commentEnd]
	l.emitCommentToken(LINE_COMMENT, startLine, startCol, consumed)
	return commentEnd
}

func (l *Lexer) consumeSlashHashBlockComment(startLine int, startCol int) int {
	depth := 1
	i := l.index + 2
	for i < len(l.input)-1 {
		if l.input[i] == '/' && l.input[i+1] == '#' {
			depth++
			i += 2
			continue
		}
		if l.input[i] == '#' && l.input[i+1] == '/' {
			depth--
			i += 2
			if depth == 0 {
				consumed := l.input[l.index:i]
				l.emitCommentToken(BLOCK_COMMENT, startLine, startCol, consumed)
				return i - l.index
			}
			continue
		}
		i++
	}
	l.diagnostics = append(l.diagnostics, d.New("unterminated block comment", startLine, startCol, startLine, startCol, "error"))
	return len(l.input) - l.index
}

func (l *Lexer) consumeCBlockComment(startLine int, startCol int) int {
	i := l.index + 2
	for i < len(l.input)-1 {
		if l.input[i] == '*' && l.input[i+1] == '/' {
			consumed := l.input[l.index : i+2]
			l.emitCommentToken(BLOCK_COMMENT, startLine, startCol, consumed)
			return i - l.index + 2
		}
		i++
	}
	l.diagnostics = append(l.diagnostics, d.New("unterminated block comment", startLine, startCol, startLine, startCol, "error"))
	return len(l.input) - l.index
}

func (l *Lexer) consumeString(startLine int, startCol int) int {
	i := l.index + 1
	escaped := false
	for i < len(l.input) {
		c := l.input[i]
		if escaped {
			escaped = false
			i++
			continue
		}
		if c == '\\' {
			escaped = true
			i++
			continue
		}
		if c == '"' {
			stringContent := l.input[l.index+1 : i]
			endCol := startCol + (i - l.index)
			l.tokens = append(l.tokens, Token{Type: STRING, Content: string(stringContent), Line: startLine, Col: startCol, EndLine: startLine, EndCol: endCol, StartOffset: l.index, EndOffset: i})
			return i - l.index + 1
		}
		i++
	}
	l.diagnostics = append(l.diagnostics, d.New("unterminated string literal", startLine, startCol, startLine, startCol, "error"))
	return len(l.input) - l.index
}

func (l *Lexer) handleOperatorToken(c byte, startLine int, startCol int) int {
	startOffset := l.index
	if l.index+1 < len(l.input) {
		next := l.input[l.index+1]
		if tokenType, ok := twoCharTokens[string([]byte{c, next})]; ok {
			l.HandleBuffer()
			l.emitToken(tokenType, string([]byte{c, next}), startLine, startCol, startLine, startCol+1, startOffset, startOffset+1)
			return 2
		}
	}

	l.HandleBuffer()
	if c == '=' {
		l.emitToken(ASSIGNMENT, "=", startLine, startCol, startLine, startCol, startOffset, startOffset)
		return 1
	}
	l.emitToken(OPERATOR, strings.TrimSpace(string(c)), startLine, startCol, startLine, startCol, startOffset, startOffset)
	return 1
}

func TokensUntilAny(tokens []Token, targets []TokenType) []Token {
	token_buffer := []Token{}

	for _, t := range tokens {
		if t.Type == LINE_COMMENT || t.Type == BLOCK_COMMENT {
			return token_buffer
		}
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
	case '+', '-', '*', '/', '%', '?', '^', '~', '<', '>', '&', '|', '=', '!':
		if c == '/' && l.index+1 < len(l.input) {
			next := l.input[l.index+1]
			if next == '/' {
				l.HandleBuffer()
				return l.consumeLineComment(startLine, startCol)
			}
			if next == '#' {
				l.HandleBuffer()
				return l.consumeSlashHashBlockComment(startLine, startCol)
			}
			if next == '*' {
				l.HandleBuffer()
				return l.consumeCBlockComment(startLine, startCol)
			}
		}
		return l.handleOperatorToken(c, startLine, startCol)
	case '#':
		if l.index+1 < len(l.input) && l.input[l.index+1] == '/' {
			l.HandleBuffer()
			return l.consumeLineComment(startLine, startCol)
		}
		if len(l.buffer) == 0 {
			l.bufferLine = l.line
			l.bufferCol = l.col
			l.bufferIndex = l.index
		}
		l.buffer = append(l.buffer, c)
		return 1
	case '\n', '(', '[', '{', ')', ']', '}', ';', ',':
		l.HandleBuffer()
		token_type := TokenTypeFromChar(c)
		l.tokens = append(l.tokens, Token{Type: token_type, Content: strings.TrimSpace(string(c)), Line: startLine, Col: startCol, EndLine: startLine, EndCol: startCol, StartOffset: l.index, EndOffset: l.index})
		return 1
	case ':':
		if l.index+1 < len(l.input) && l.input[l.index+1] == ':' {
			if len(l.buffer) == 0 {
				l.bufferLine = l.line
				l.bufferCol = l.col
				l.bufferIndex = l.index
			}
			l.buffer = append(l.buffer, ':', ':')
			return 2
		}
		l.HandleBuffer()
		l.tokens = append(l.tokens, Token{Type: COLON, Content: ":", Line: startLine, Col: startCol, EndLine: startLine, EndCol: startCol, StartOffset: l.index, EndOffset: l.index})
		return 1
	case ' ', '\t':
		l.HandleBuffer()
		return 1
	case '"':
		l.HandleBuffer()
		return l.consumeString(startLine, startCol)
	default:
		if len(l.buffer) == 0 {
			l.bufferLine = l.line
			l.bufferCol = l.col
			l.bufferIndex = l.index
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
