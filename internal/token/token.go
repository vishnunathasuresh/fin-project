package token

type Type string

const (
	EOF     Type = "EOF"
	ILLEGAL Type = "ILLEGAL"

	IDENT   Type = "IDENT"
	NUMBER  Type = "NUMBER"
	STRING  Type = "STRING"
	NEWLINE Type = "NEWLINE"

	// keywords
	DEF      Type = "DEF"
	TYPE     Type = "TYPE"
	IMPORT   Type = "IMPORT"
	NIL      Type = "NIL"
	IF       Type = "IF"
	ELSE     Type = "ELSE"
	FOR      Type = "FOR"
	WHILE    Type = "WHILE"
	RETURN   Type = "RETURN"
	BREAK    Type = "BREAK"
	CONTINUE Type = "CONTINUE"
	TRUE     Type = "TRUE"
	FALSE    Type = "FALSE"

	// punctuation
	DOTDOT Type = ".."
	DOT    Type = "."

	LBRACKET Type = "["
	RBRACKET Type = "]"
	LBRACE   Type = "{"
	RBRACE   Type = "}"
	LPAREN   Type = "("
	RPAREN   Type = ")"
	COMMA    Type = ","
	COLON    Type = ":"

	// operators
	DECLARE Type = ":="
	ASSIGN  Type = "="
	ARROW   Type = "->"
	PLUS    Type = "+"
	MINUS   Type = "-"
	STAR    Type = "*"
	SLASH   Type = "/"

	// command literal delimiters
	LANGLE Type = "<"
	RANGLE Type = ">"
)

type Token struct {
	Type    Type
	Literal string
	Line    int
	Column  int
}

func New(t Type, lit string, line, col int) Token {
	return Token{
		Type:    t,
		Literal: lit,
		Line:    line,
		Column:  col,
	}
}

var Keywords = map[string]Type{
	"def":      DEF,
	"type":     TYPE,
	"import":   IMPORT,
	"nil":      NIL,
	"if":       IF,
	"else":     ELSE,
	"for":      FOR,
	"while":    WHILE,
	"return":   RETURN,
	"break":    BREAK,
	"continue": CONTINUE,
	"true":     TRUE,
	"false":    FALSE,
}

func LookupIdent(ident string) Type {
	if tok, ok := Keywords[ident]; ok {
		return tok
	}
	return IDENT
}
