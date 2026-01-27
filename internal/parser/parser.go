package parser

import "github.com/vishnunath-suresh/fin-project/internal/token"

// Parser holds token stream state for recursive-descent parsing.
type Parser struct {
    tokens []token.Token
    pos    int
}

// New creates a parser from a token slice.
func New(tokens []token.Token) *Parser {
    return &Parser{tokens: tokens, pos: 0}
}

// current returns the token at the current position safely (EOF if out of bounds).
func (p *Parser) current() token.Token {
    if len(p.tokens) == 0 {
        return token.Token{Type: token.EOF}
    }
    if p.pos >= len(p.tokens) {
        return p.tokens[len(p.tokens)-1]
    }
    return p.tokens[p.pos]
}

// next advances the parser if not at EOF and returns the token that was current before advancing.
func (p *Parser) next() token.Token {
    tok := p.current()
    if tok.Type != token.EOF && p.pos < len(p.tokens) {
        p.pos++
    }
    return tok
}

// match checks whether the current token is one of the given types; if so, it consumes it and returns true.
func (p *Parser) match(types ...token.Type) bool {
    for _, t := range types {
        if p.check(t) {
            p.next()
            return true
        }
    }
    return false
}

// check reports whether the current token is of the given type.
func (p *Parser) check(t token.Type) bool {
    return p.current().Type == t
}

// expect ensures the current token matches the given type, consuming it on success.
// Returns (token, true) on success, or (zero, false) on failure without advancing.
func (p *Parser) expect(t token.Type) (token.Token, bool) {
    if p.check(t) {
        return p.next(), true
    }
    return token.Token{}, false
}

// isAtEnd reports whether the parser has reached EOF.
func (p *Parser) isAtEnd() bool {
    return p.current().Type == token.EOF
}
