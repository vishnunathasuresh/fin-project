package parser

import (
	"fmt"

	"github.com/vishnunath-suresh/fin-project/internal/ast"
	"github.com/vishnunath-suresh/fin-project/internal/token"
)

type Parser struct {
	tokens []token.Token
	pos    int
	errors []error
}

func New(tokens []token.Token) *Parser {
	return &Parser{
		tokens: tokens,
		pos:    0,
	}
}

//
// ---- Public API ----
//

// Parse parses the full token stream into an AST program.
// It never panics. Errors are accumulated in p.errors.
func (p *Parser) Parse() *ast.Program {
	prog := &ast.Program{
		Statements: make([]ast.Statement, 0),
		P:          p.currentPos(),
	}

	for !p.isAtEnd() {
		if p.match(token.NEWLINE) {
			continue
		}

		stmt := p.parseStatement()
		if stmt != nil {
			prog.Statements = append(prog.Statements, stmt)
		} else {
			p.sync()
		}
	}

	return prog
}

func (p *Parser) Errors() []error {
	return p.errors
}

//
// ---- Core Helpers ----
//

func (p *Parser) current() token.Token {
	if p.pos >= len(p.tokens) {
		return p.tokens[len(p.tokens)-1]
	}
	return p.tokens[p.pos]
}

func (p *Parser) previous() token.Token {
	return p.tokens[p.pos-1]
}

func (p *Parser) next() token.Token {
	if !p.isAtEnd() {
		p.pos++
	}
	return p.previous()
}

func (p *Parser) match(t token.Type) bool {
	if p.check(t) {
		p.next()
		return true
	}
	return false
}

func (p *Parser) check(t token.Type) bool {
	if p.isAtEnd() {
		return false
	}
	return p.current().Type == t
}

func (p *Parser) expect(t token.Type, msg string) token.Token {
	if p.check(t) {
		return p.next()
	}

	tok := p.current()
	p.errorAt(tok, msg)
	return tok
}

func (p *Parser) isAtEnd() bool {
	return p.current().Type == token.EOF
}

func (p *Parser) currentPos() ast.Pos {
	tok := p.current()
	return ast.Pos{
		Line:   tok.Line,
		Column: tok.Column,
	}
}

//
// ---- Error Handling ----
//

func (p *Parser) errorAt(tok token.Token, msg string) {
	err := fmt.Errorf(
		"%d:%d: %s (got %q)",
		tok.Line,
		tok.Column,
		msg,
		tok.Literal,
	)
	p.errors = append(p.errors, err)
}

// sync advances tokens until a reasonable recovery point.
// This prevents cascading errors.
func (p *Parser) sync() {
	p.next()

	for !p.isAtEnd() {
		if p.previous().Type == token.NEWLINE {
			return
		}

		switch p.current().Type {
		case token.SET,
			token.ECHO,
			token.RUN,
			token.FN,
			token.IF,
			token.FOR,
			token.WHILE,
			token.RETURN:
			return
		}

		p.next()
	}
}

//
// ---- Statement Dispatch ----
//

func (p *Parser) parseStatement() ast.Statement {
	switch p.current().Type {

	case token.SET:
		return p.parseSet()

	case token.ECHO:
		return p.parseEcho()

	case token.RUN:
		return p.parseRun()

	case token.FN:
		return p.parseFn()

	case token.IF:
		return p.parseIf()

	case token.FOR:
		return p.parseFor()

	case token.WHILE:
		return p.parseWhile()

	case token.RETURN:
		return p.parseReturn()

	case token.IDENT:
		return p.parseCall()

	default:
		tok := p.current()
		p.errorAt(tok, "unexpected token")
		return nil
	}
}

//
// ---- Stub Parsers (to be implemented next) ----
//

func (p *Parser) parseSet() ast.Statement {
	p.errorAt(p.current(), "parseSet not implemented")
	return nil
}

func (p *Parser) parseEcho() ast.Statement {
	p.errorAt(p.current(), "parseEcho not implemented")
	return nil
}

func (p *Parser) parseRun() ast.Statement {
	p.errorAt(p.current(), "parseRun not implemented")
	return nil
}

func (p *Parser) parseFn() ast.Statement {
	p.errorAt(p.current(), "parseFn not implemented")
	return nil
}

func (p *Parser) parseIf() ast.Statement {
	p.errorAt(p.current(), "parseIf not implemented")
	return nil
}

func (p *Parser) parseFor() ast.Statement {
	p.errorAt(p.current(), "parseFor not implemented")
	return nil
}

func (p *Parser) parseWhile() ast.Statement {
	p.errorAt(p.current(), "parseWhile not implemented")
	return nil
}

func (p *Parser) parseReturn() ast.Statement {
	p.errorAt(p.current(), "parseReturn not implemented")
	return nil
}

func (p *Parser) parseCall() ast.Statement {
	p.errorAt(p.current(), "parseCall not implemented")
	return nil
}
