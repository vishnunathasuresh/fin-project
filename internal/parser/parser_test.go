package parser

import (
	"testing"

	"github.com/vishnunathasuresh/fin-project/internal/token"
)

// Table-driven coverage of parser helpers: current, next, match, check, expect, isAtEnd.
// Focus on EOF/off-by-one safety and non-panicking behavior.

func TestHelpers_CurrentAndIsAtEnd(t *testing.T) {
	tests := []struct {
		name   string
		tokens []token.Token
		pos    int
		want   token.Type
		end    bool
	}{
		{name: "empty slice", tokens: nil, pos: 0, want: token.EOF, end: true},
		{name: "at start", tokens: []token.Token{{Type: token.IDENT}}, pos: 0, want: token.IDENT, end: false},
		{name: "past end returns last", tokens: []token.Token{{Type: token.IDENT}, {Type: token.EOF}}, pos: 5, want: token.EOF, end: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{tokens: tt.tokens, pos: tt.pos}
			if got := p.current().Type; got != tt.want {
				t.Fatalf("current Type = %s, want %s", got, tt.want)
			}
			if end := p.isAtEnd(); end != tt.end {
				t.Fatalf("isAtEnd = %v, want %v", end, tt.end)
			}
		})
	}
}

func TestHelpers_Next(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []token.Token
		pos      int
		wantType token.Type
		wantPos  int
	}{
		{name: "advance normal", tokens: []token.Token{{Type: token.IDENT}, {Type: token.NUMBER}, {Type: token.EOF}}, pos: 0, wantType: token.IDENT, wantPos: 1},
		{name: "stop at EOF", tokens: []token.Token{{Type: token.EOF}}, pos: 0, wantType: token.EOF, wantPos: 0},
		{name: "past end stays", tokens: []token.Token{{Type: token.IDENT}, {Type: token.EOF}}, pos: 5, wantType: token.EOF, wantPos: 5},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{tokens: tt.tokens, pos: tt.pos}
			got := p.next()
			if got.Type != tt.wantType {
				t.Fatalf("next Type = %s, want %s", got.Type, tt.wantType)
			}
			if p.pos != tt.wantPos {
				t.Fatalf("pos = %d, want %d", p.pos, tt.wantPos)
			}
		})
	}
}

func TestHelpers_CheckMatchExpect(t *testing.T) {
	tests := []struct {
		name       string
		tokens     []token.Token
		wantChecks []bool
		matchTypes []token.Type
		expectType token.Type
		expectOK   bool
		finalPos   int
		finalEOF   bool
	}{
		{
			name:       "match consumes and expect succeeds then fails",
			tokens:     []token.Token{{Type: token.IDENT}, {Type: token.NUMBER}, {Type: token.EOF}},
			wantChecks: []bool{true, false},
			matchTypes: []token.Type{token.NUMBER, token.IDENT}, // IDENT matches first
			expectType: token.NUMBER,
			expectOK:   true,
			finalPos:   2,
			finalEOF:   true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{tokens: tt.tokens, pos: 0}

			if got := p.check(token.IDENT); got != tt.wantChecks[0] {
				t.Fatalf("check IDENT = %v, want %v", got, tt.wantChecks[0])
			}
			if got := p.check(token.NUMBER); got != tt.wantChecks[1] {
				t.Fatalf("check NUMBER = %v, want %v", got, tt.wantChecks[1])
			}

			if !p.match(tt.matchTypes...) {
				t.Fatalf("match should consume IDENT")
			}
			if p.pos != 1 {
				t.Fatalf("pos after match = %d, want 1", p.pos)
			}

			tok, ok := p.expect(tt.expectType)
			if ok != tt.expectOK {
				t.Fatalf("expect ok = %v, want %v", ok, tt.expectOK)
			}
			if tok.Type != tt.expectType {
				t.Fatalf("expect type = %s, want %s", tok.Type, tt.expectType)
			}
			if p.pos != tt.finalPos {
				t.Fatalf("pos after expect = %d, want %d", p.pos, tt.finalPos)
			}

			// expect failure does not advance
			failTok, failOK := p.expect(token.IDENT)
			if failOK {
				t.Fatalf("expect should fail at EOF")
			}
			if failTok.Type != "" {
				t.Fatalf("fail token Type = %s, want empty", failTok.Type)
			}
			if p.pos != tt.finalPos {
				t.Fatalf("pos after failed expect = %d, want %d", p.pos, tt.finalPos)
			}

			if atEnd := p.isAtEnd(); atEnd != tt.finalEOF {
				t.Fatalf("isAtEnd = %v, want %v", atEnd, tt.finalEOF)
			}
		})
	}
}
