package parser_test

import (
	"testing"

	"github.com/vishnunathasuresh/fin-project/internal/lexer"
	"github.com/vishnunathasuresh/fin-project/internal/parser"
	"github.com/vishnunathasuresh/fin-project/internal/token"
)

func TestCollectTokens_SimpleProgram(t *testing.T) {
	src := "set name \"Alice\"\necho $name\n"
	l := lexer.New(src)
	toks := parser.CollectTokens(l)

	// Expected tokens: SET, IDENT(name), STRING("Alice"), NEWLINE, ECHO, IDENT(name), NEWLINE, EOF
	const expectedLen = 8
	if len(toks) != expectedLen {
		t.Fatalf("expected %d tokens, got %d", expectedLen, len(toks))
	}

	eofCount := 0
	for i, tok := range toks {
		if tok.Type == token.EOF {
			eofCount++
			if i != len(toks)-1 {
				t.Errorf("EOF token at index %d is not last", i)
			}
		}
	}
	if eofCount != 1 {
		t.Fatalf("expected exactly one EOF token, got %d", eofCount)
	}

	expected := []struct {
		typ token.Type
		lit string
	}{
		{token.SET, "set"},
		{token.IDENT, "name"},
		{token.STRING, "Alice"},
		{token.NEWLINE, "\n"},
		{token.ECHO, "echo"},
		{token.IDENT, "name"},
		{token.NEWLINE, "\n"},
		{token.EOF, ""},
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
