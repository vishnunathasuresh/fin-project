package lexer

import (
	"unicode"

	"github.com/vishnunathasuresh/fin-project/internal/token"
)

type Lexer struct {
	input []rune
	pos   int
	line  int
	col   int
	inCmd bool

	indentStack []int         // stack of indent widths; starts with 0
	pending     []token.Token // buffered tokens (INDENT/DEDENT/EOF)
	atLineStart bool          // true if next token is at start of line
}

// readIndentWidth reads spaces/tabs at line start and returns the indent width (tab=4 spaces).
// Stops at first non-space/tab or newline. Returns ok=false if an unexpected character is seen.
func (l *Lexer) readIndentWidth() (int, bool) {
	width := 0
	for {
		ch := l.peek()
		switch ch {
		case ' ':
			l.next()
			width++
		case '\t':
			l.next()
			width += 4
		case '\n', 0:
			return width, true
		default:
			return width, true
		}
	}
}

func New(input string) *Lexer {
	return &Lexer{
		input:       []rune(input),
		pos:         0,
		line:        1,
		col:         1,
		indentStack: []int{0},
		pending:     nil,
		atLineStart: true,
	}
}

func (l *Lexer) NextToken() token.Token {
	// Flush any pending tokens first (e.g., INDENT/DEDENT/EOF)
	if len(l.pending) > 0 {
		tok := l.pending[0]
		l.pending = l.pending[1:]
		return tok
	}

	// Handle indentation at the start of a line (outside command literals)
	if !l.inCmd && l.atLineStart {
		l.atLineStart = false
		indentWidth, ok := l.readIndentWidth()
		if !ok {
			return token.New(token.ILLEGAL, "invalid indentation", l.line, l.col)
		}
		prev := l.indentStack[len(l.indentStack)-1]
		switch {
		case indentWidth > prev:
			l.indentStack = append(l.indentStack, indentWidth)
			l.pending = append(l.pending, token.New(token.INDENT, "", l.line, l.col))
			return l.NextToken()
		case indentWidth < prev:
			for len(l.indentStack) > 0 && indentWidth < l.indentStack[len(l.indentStack)-1] {
				l.indentStack = l.indentStack[:len(l.indentStack)-1]
				l.pending = append(l.pending, token.New(token.DEDENT, "", l.line, l.col))
			}
			if len(l.indentStack) == 0 || l.indentStack[len(l.indentStack)-1] != indentWidth {
				return token.New(token.ILLEGAL, "inconsistent indentation", l.line, l.col)
			}
			if len(l.pending) > 0 {
				return l.NextToken()
			}
		}
	}

	if !l.inCmd {
		l.skipWhitespaceExceptNewline()
	}

	startLine := l.line
	startCol := l.col

	ch := l.peek()

	// Inside command literal: only emit CMD_TEXT or CMD_END.
	if l.inCmd {
		switch ch {
		case 0:
			return token.New(token.ILLEGAL, "unterminated command", startLine, startCol)
		case '>':
			l.inCmd = false
			l.next()
			return token.New(token.CMD_END, ">", startLine, startCol)
		default:
			text := l.readCommandText()
			return token.New(token.CMD_TEXT, text, startLine, startCol)
		}
	}

	switch {
	case ch == 0:
		// Emit any remaining dedents before EOF
		if len(l.indentStack) > 1 {
			l.indentStack = l.indentStack[:len(l.indentStack)-1]
			return token.New(token.DEDENT, "", startLine, startCol)
		}
		return token.New(token.EOF, "", startLine, startCol)

	case ch == '\n':
		l.next()
		l.atLineStart = true
		return token.New(token.NEWLINE, "\n", startLine, startCol)

	case ch == '#':
		l.skipComment()
		return l.NextToken()

	case isLetter(ch):
		literal := l.readIdentifier()
		typ := token.LookupIdent(literal)
		return token.New(typ, literal, startLine, startCol)

	case isDigit(ch):
		return token.New(token.NUMBER, l.readNumber(), startLine, startCol)

	case ch == '"':
		str, ok := l.readString()
		if !ok {
			return token.New(token.ILLEGAL, str, startLine, startCol)
		}
		return token.New(token.STRING, str, startLine, startCol)

	case ch == '.':
		l.next()
		return token.New(token.DOT, ".", startLine, startCol)

	case ch == '|':
		if l.peekNext() == '|' {
			l.next()
			l.next()
			return token.New(token.OR, "||", startLine, startCol)
		}
		l.next()
		return token.New(token.ILLEGAL, "|", startLine, startCol)

	case ch == '&':
		if l.peekNext() == '&' {
			l.next()
			l.next()
			return token.New(token.AND, "&&", startLine, startCol)
		}
		l.next()
		return token.New(token.ILLEGAL, "&", startLine, startCol)

	case ch == '[':
		l.next()
		return token.New(token.LBRACKET, "[", startLine, startCol)

	case ch == ']':
		l.next()
		return token.New(token.RBRACKET, "]", startLine, startCol)

	case ch == '{':
		l.next()
		return token.New(token.LBRACE, "{", startLine, startCol)

	case ch == '}':
		l.next()
		return token.New(token.RBRACE, "}", startLine, startCol)

	case ch == '(':
		l.next()
		return token.New(token.LPAREN, "(", startLine, startCol)

	case ch == ')':
		l.next()
		return token.New(token.RPAREN, ")", startLine, startCol)

	case ch == ',':
		l.next()
		return token.New(token.COMMA, ",", startLine, startCol)

	case ch == ':':
		if l.peekNext() == '=' {
			l.next()
			l.next()
			return token.New(token.DECLARE, ":=", startLine, startCol)
		}
		l.next()
		return token.New(token.COLON, ":", startLine, startCol)

	case ch == '$':
		// Variable reference: $name
		l.next() // consume '$'
		if !isLetter(l.peek()) {
			return token.New(token.ILLEGAL, string(ch), startLine, startCol)
		}
		ident := l.readIdentifier()
		return token.New(token.IDENT, ident, startLine, startCol)
	case ch == '+':
		l.next()
		return token.New(token.PLUS, "+", startLine, startCol)

	case ch == '-':
		if l.peekNext() == '>' {
			l.next()
			l.next()
			return token.New(token.ARROW, "->", startLine, startCol)
		}
		l.next()
		return token.New(token.MINUS, "-", startLine, startCol)

	case ch == '*':
		if l.peekNext() == '*' {
			l.next()
			l.next()
			return token.New(token.POWER, "**", startLine, startCol)
		}
		l.next()
		return token.New(token.STAR, "*", startLine, startCol)

	case ch == '/':
		l.next()
		return token.New(token.SLASH, "/", startLine, startCol)

	case ch == '=':
		if l.peekNext() == '=' {
			l.next()
			l.next()
			return token.New(token.EQ, "==", startLine, startCol)
		}
		l.next()
		return token.New(token.ASSIGN, "=", startLine, startCol)

	case ch == '!':
		if l.peekNext() == '=' {
			l.next()
			l.next()
			return token.New(token.NEQ, "!=", startLine, startCol)
		}
		l.next()
		return token.New(token.BANG, "!", startLine, startCol)

	case ch == '<':
		// command literal start
		l.inCmd = true
		l.next()
		return token.New(token.CMD_START, "<", startLine, startCol)

	case ch == '>':
		// standalone '>' outside command is ILLEGAL under Fin v2
		l.next()
		return token.New(token.ILLEGAL, ">", startLine, startCol)

	default:
		l.next()
		return token.New(token.ILLEGAL, string(ch), startLine, startCol)
	}
}

func (l *Lexer) peek() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	return l.input[l.pos]
}

func (l *Lexer) peekNext() rune {
	if l.pos+1 >= len(l.input) {
		return 0
	}
	return l.input[l.pos+1]
}

func (l *Lexer) next() rune {
	ch := l.peek()
	l.pos++

	if ch == '\n' {
		l.line++
		l.col = 1
	} else {
		l.col++
	}

	return ch
}

func (l *Lexer) skipWhitespaceExceptNewline() {
	for {
		ch := l.peek()
		if ch == ' ' || ch == '\t' || ch == '\r' {
			l.next()
			continue
		}
		break
	}
}

func (l *Lexer) skipComment() {
	for {
		ch := l.peek()
		if ch == '\n' || ch == 0 {
			return
		}
		l.next()
	}
}

func (l *Lexer) readIdentifier() string {
	start := l.pos
	for isLetter(l.peek()) || isDigit(l.peek()) {
		l.next()
	}
	return string(l.input[start:l.pos])
}

func (l *Lexer) readNumber() string {
	start := l.pos
	for isDigit(l.peek()) {
		l.next()
	}
	return string(l.input[start:l.pos])
}

// readCommandText reads raw text until the next '>' or EOF without consuming the '>' delimiter.
func (l *Lexer) readCommandText() string {
	start := l.pos
	for {
		ch := l.peek()
		if ch == 0 || ch == '>' {
			break
		}
		l.next()
	}
	return string(l.input[start:l.pos])
}

func (l *Lexer) readString() (string, bool) {
	l.next() // consume opening quote

	start := l.pos
	var out []rune
	for {
		ch := l.peek()
		if ch == 0 {
			// Unterminated string; return what we have, mark not ok.
			return string(l.input[start:l.pos]), false
		}
		if ch == '"' {
			break
		}
		if ch == '\\' {
			l.next()
			esc := l.peek()
			switch esc {
			case '"':
				out = append(out, '"')
			case '\\':
				out = append(out, '\\')
			case 'n':
				out = append(out, '\n')
			case 't':
				out = append(out, '\t')
			default:
				// Unknown escape, treat as literal char.
				out = append(out, esc)
			}
			l.next()
			continue
		}
		out = append(out, ch)
		l.next()
	}

	l.next() // closing quote

	// If no escapes were encountered, slice directly for efficiency.
	if len(out) == 0 {
		return string(l.input[start : l.pos-1]), true
	}
	return string(out), true
}

func isLetter(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}
