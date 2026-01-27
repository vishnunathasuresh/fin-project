package parser

import (
	"fmt"

	"github.com/vishnunath-suresh/fin-project/internal/ast"
	"github.com/vishnunath-suresh/fin-project/internal/token"
)

// Parser holds token stream state for recursive-descent parsing.
type Parser struct {
	tokens []token.Token
	pos    int
	errors []error
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

// Errors returns the collected parse errors.
func (p *Parser) Errors() []error {
	return p.errors
}

// ParseProgram is the top-level entry that produces a Program AST.
// It skips NEWLINE tokens, stops at EOF, appends successfully parsed statements,
// and uses synchronization to recover from errors without panicking.
func (p *Parser) ParseProgram() *ast.Program {
	prog := &ast.Program{P: ast.Pos{Line: 1, Column: 1}}

	for !p.isAtEnd() {
		if p.check(token.NEWLINE) {
			p.next()
			continue
		}

		stmt := p.parseStatement()
		if stmt != nil {
			prog.Statements = append(prog.Statements, stmt)
			continue
		}

		p.synchronize()
	}

	return prog
}

func (p *Parser) parseStatement() ast.Statement {
	tok := p.current()
	if tok.Type == token.EOF {
		return nil
	}
	if tok.Type == token.ILLEGAL {
		p.errors = append(p.errors, fmt.Errorf("illegal token: %s", tok.Literal))
		p.next()
		return nil
	}

	// Recognize simple statement starts; build minimal AST nodes to satisfy traversal.
	switch tok.Type {
	case token.IDENT:
		p.next()
		return &ast.CallStmt{Name: tok.Literal, P: ast.Pos{Line: tok.Line, Column: tok.Column}}
	case token.ECHO:
		p.next()
		return &ast.EchoStmt{P: ast.Pos{Line: tok.Line, Column: tok.Column}}
	case token.RUN:
		p.next()
		return &ast.RunStmt{P: ast.Pos{Line: tok.Line, Column: tok.Column}}
	case token.SET:
		p.next()
		return &ast.SetStmt{Name: tok.Literal, P: ast.Pos{Line: tok.Line, Column: tok.Column}}
	case token.RETURN:
		p.next()
		return &ast.ReturnStmt{P: ast.Pos{Line: tok.Line, Column: tok.Column}}
	default:
		p.errors = append(p.errors, fmt.Errorf("unexpected token: %s", tok.Type))
		p.next()
		return nil
	}
}

// synchronize advances until after a newline or EOF to recover from an error.
func (p *Parser) synchronize() {
	for !p.isAtEnd() {
		if p.check(token.NEWLINE) {
			p.next()
			return
		}
		p.next()
	}
}
