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
	BASH     Type = "BASH"
	BAT      Type = "BAT"
	PS1      Type = "PS1"

	// punctuation
	DOT Type = "."

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
	EQ      Type = "=="
	NEQ     Type = "!="
	ARROW   Type = "->"
	OR      Type = "||"
	AND     Type = "&&"
	PLUS    Type = "+"
	MINUS   Type = "-"
	STAR    Type = "*"
	POWER   Type = "**"
	SLASH   Type = "/"
	BANG    Type = "!"

	// command literal delimiters
	CMD_START Type = "CMD_START"
	CMD_TEXT  Type = "CMD_TEXT"
	CMD_END   Type = "CMD_END"
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
	"and":      AND,
	"or":       OR,
	"not":      BANG,
	"bash":     BASH,
	"bat":      BAT,
	"ps1":      PS1,
}

func LookupIdent(ident string) Type {
	if tok, ok := Keywords[ident]; ok {
		return tok
	}
	return IDENT
}
