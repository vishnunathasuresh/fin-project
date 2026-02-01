package parser_test

import (
	"testing"

	"github.com/vishnunathasuresh/fin-project/internal/lexer"
	"github.com/vishnunathasuresh/fin-project/internal/parser"
	"github.com/vishnunathasuresh/fin-project/internal/token"
)

func TestCollectTokens_FinV2(t *testing.T) {
	src := "" +
		"x := 1\n" +
		"def foo(a: int) -> int\n\n" +
		"cmd := <grep \"abc\" file.txt>\n"

	l := lexer.New(src)
	toks := parser.CollectTokens(l)

	expected := []struct {
		typ token.Type
		lit string
	}{
		{token.IDENT, "x"},
		{token.DECLARE, ":="},
		{token.NUMBER, "1"},
		{token.NEWLINE, "\n"},
		{token.DEF, "def"},
		{token.IDENT, "foo"},
		{token.LPAREN, "("},
		{token.IDENT, "a"},
		{token.COLON, ":"},
		{token.IDENT, "int"},
		{token.RPAREN, ")"},
		{token.ARROW, "->"},
		{token.IDENT, "int"},
		{token.NEWLINE, "\n"},
		{token.NEWLINE, "\n"},
		{token.IDENT, "cmd"},
		{token.DECLARE, ":="},
		{token.CMD_START, "<"},
		{token.CMD_TEXT, "grep \"abc\" file.txt"},
		{token.CMD_END, ">"},
		{token.NEWLINE, "\n"},
		{token.EOF, ""},
	}

	if len(toks) != len(expected) {
		t.Fatalf("expected %d tokens, got %d", len(expected), len(toks))
	}
	for i, exp := range expected {
		if toks[i].Type != exp.typ {
			t.Errorf("token %d: expected type %s, got %s", i, exp.typ, toks[i].Type)
		}
		if toks[i].Literal != exp.lit {
			t.Errorf("token %d: expected literal %q, got %q", i, exp.lit, toks[i].Literal)
		}
	}
}

func TestCollectTokens_EmptyInput(t *testing.T) {
	l := lexer.New("")
	toks := parser.CollectTokens(l)

	if len(toks) != 1 {
		t.Fatalf("expected 1 token (EOF) for empty input, got %d", len(toks))
	}
	if toks[0].Type != token.EOF {
		t.Fatalf("expected EOF token for empty input, got %s", toks[0].Type)
	}
}

func TestCollectTokens_PreservesNewlines(t *testing.T) {
	l := lexer.New("\n\n")
	toks := parser.CollectTokens(l)

	if len(toks) != 3 {
		t.Fatalf("expected 3 tokens, got %d", len(toks))
	}
	if toks[0].Type != token.NEWLINE || toks[1].Type != token.NEWLINE {
		t.Fatalf("NEWLINE tokens not preserved: %v, %v", toks[0].Type, toks[1].Type)
	}
	if toks[2].Type != token.EOF {
		t.Fatalf("final token is not EOF: %v", toks[2].Type)
	}
}
