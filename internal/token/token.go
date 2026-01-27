package token

type Type string

const (
	EOF     Type = "EOF"
	ILLEGAL Type = "ILLEGAL"

	IDENT   Type = "IDENT"
	NUMBER  Type = "NUMBER"
	STRING  Type = "STRING"
	NEWLINE Type = "NEWLINE"

	SET    Type = "SET"
	ECHO   Type = "ECHO"
	RUN    Type = "RUN"
	IF     Type = "IF"
	ELSE   Type = "ELSE"
	END    Type = "END"
	FOR    Type = "FOR"
	IN     Type = "IN"
	EXISTS Type = "EXISTS"
	FN     Type = "FN"

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
	WHILE    Type = "WHILE"
	RETURN   Type = "RETURN"
	TRUE     Type = "TRUE"
	FALSE    Type = "FALSE"
	PLUS     Type = "+"
	MINUS    Type = "-"
	STAR     Type = "*"
	SLASH    Type = "/"
	BANG     Type = "!"
	POW      Type = "**"
	EQEQ     Type = "=="
	NOTEQ    Type = "!="
	LT       Type = "<"
	LTE      Type = "<="
	GT       Type = ">"
	GTE      Type = ">="
	AND      Type = "&&"
	OR       Type = "||"
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
	"set":    SET,
	"echo":   ECHO,
	"run":    RUN,
	"if":     IF,
	"else":   ELSE,
	"end":    END,
	"for":    FOR,
	"while":  WHILE,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
	"in":     IN,
	"exists": EXISTS,
	"fn":     FN,
}

func LookupIdent(ident string) Type {
	if tok, ok := Keywords[ident]; ok {
		return tok
	}
	return IDENT
}
