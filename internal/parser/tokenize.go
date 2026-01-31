package parser

import (
	"fmt"

	"github.com/vishnunathasuresh/fin-project/internal/lexer"
	"github.com/vishnunathasuresh/fin-project/internal/token"
)

// CollectTokens drains the lexer into a slice of tokens, preserving order and positions.
// It stops after reading the first EOF and validates that exactly one EOF exists and it is last.
// NEWLINE tokens are preserved. Panics if the stream is invalid.
func CollectTokens(l *lexer.Lexer) []token.Token {
	var tokens []token.Token
	eofCount := 0

	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)

		if tok.Type == token.EOF {
			eofCount++
			break
		}
	}

	if eofCount != 1 {
		panic(fmt.Sprintf("CollectTokens: expected exactly one EOF, got %d", eofCount))
	}
	if len(tokens) == 0 || tokens[len(tokens)-1].Type != token.EOF {
		panic("CollectTokens: EOF token is not the last token")
	}

	return tokens
}
