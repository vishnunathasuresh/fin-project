package lexer

import (
	"testing"

	"github.com/vishnunathasuresh/fin-project/internal/token"
)

// collectTokens drains the lexer for assertions.
func collectTokens(l *Lexer) []token.Token {
	var toks []token.Token
	for {
		tok := l.NextToken()
		toks = append(toks, tok)
		if tok.Type == token.EOF || tok.Type == token.ILLEGAL {
			break
		}
	}
	return toks
}

func assertTokenSeq(t *testing.T, toks []token.Token, want []token.Type) {
	t.Helper()
	if len(toks) != len(want) {
		t.Fatalf("len=%d want %d tokens: %+v", len(toks), len(want), toks)
	}
	for i, tok := range toks {
		if tok.Type != want[i] {
			t.Fatalf("tok %d: got %s want %s (lit=%q)", i, tok.Type, want[i], tok.Literal)
		}
	}
}

func TestLexDeclare(t *testing.T) {
	l := New(":=\n")
	toks := collectTokens(l)
	assertTokenSeq(t, toks, []token.Type{token.DECLARE, token.NEWLINE, token.EOF})
}

func TestLexArrow(t *testing.T) {
	l := New("->")
	toks := collectTokens(l)
	assertTokenSeq(t, toks, []token.Type{token.ARROW, token.EOF})
}

func TestLexCommandLiteral(t *testing.T) {
	l := New("<grep \"abc\" file.txt>")
	toks := collectTokens(l)
	assertTokenSeq(t, toks, []token.Type{token.CMD_START, token.CMD_TEXT, token.CMD_END, token.EOF})
	if toks[1].Literal != "grep \"abc\" file.txt" {
		t.Fatalf("cmd text literal = %q", toks[1].Literal)
	}
}

func TestLexLogicalKeywordAliases(t *testing.T) {
	l := New("and or not\n")
	toks := collectTokens(l)
	assertTokenSeq(t, toks, []token.Type{token.AND, token.OR, token.BANG, token.NEWLINE, token.EOF})
}

func TestLexPlatformKeywords(t *testing.T) {
	l := New("bash bat ps1\n")
	toks := collectTokens(l)
	assertTokenSeq(t, toks, []token.Type{token.BASH, token.BAT, token.PS1, token.NEWLINE, token.EOF})
}

func TestLexElifKeyword(t *testing.T) {
	l := New("elif\n")
	toks := collectTokens(l)
	assertTokenSeq(t, toks, []token.Type{token.ELIF, token.NEWLINE, token.EOF})
}
