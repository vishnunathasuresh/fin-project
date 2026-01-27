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

	switch tok.Type {
	case token.SET:
		return p.parseSet()
	case token.ECHO:
		return p.parseEcho()
	case token.RUN:
		return p.parseRun()
	case token.RETURN:
		return p.parseReturn()
	case token.IF:
		return p.parseIf()
	case token.FOR:
		return p.parseFor()
	case token.WHILE:
		return p.parseWhile()
	case token.FN:
		return p.parseFn()
	case token.BREAK:
		return p.parseBreak()
	case token.CONTINUE:
		return p.parseContinue()
	case token.IDENT:
		return p.parseCall()
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

// --- statement parsers ---

func (p *Parser) consumeNewlineIfPresent() {
	if p.check(token.NEWLINE) {
		p.next()
	}
}

func (p *Parser) parseSet() ast.Statement {
	setTok := p.next() // consume 'set'
	nameTok, ok := p.expect(token.IDENT)
	if !ok {
		p.errors = append(p.errors, fmt.Errorf("expected identifier after set"))
		return nil
	}
	val := p.parseExpression(0)
	p.consumeNewlineIfPresent()
	return &ast.SetStmt{Name: nameTok.Literal, Value: val, P: ast.Pos{Line: setTok.Line, Column: setTok.Column}}
}

func (p *Parser) parseEcho() ast.Statement {
	echoTok := p.next() // consume 'echo'
	var val ast.Expr
	if !p.check(token.NEWLINE) && !p.isAtEnd() {
		val = p.parseExpression(0)
	}
	p.consumeNewlineIfPresent()
	return &ast.EchoStmt{Value: val, P: ast.Pos{Line: echoTok.Line, Column: echoTok.Column}}
}

func (p *Parser) parseRun() ast.Statement {
	runTok := p.next() // consume 'run'
	cmdTok, ok := p.expect(token.STRING)
	if !ok {
		p.errors = append(p.errors, fmt.Errorf("expected string after run"))
		return nil
	}
	p.consumeNewlineIfPresent()
	return &ast.RunStmt{Command: &ast.StringLit{Value: cmdTok.Literal, P: ast.Pos{Line: cmdTok.Line, Column: cmdTok.Column}}, P: ast.Pos{Line: runTok.Line, Column: runTok.Column}}
}

func (p *Parser) parseReturn() ast.Statement {
	retTok := p.next() // consume 'return'
	var val ast.Expr
	if !p.check(token.NEWLINE) && !p.isAtEnd() {
		val = p.parseExpression(0)
	}
	p.consumeNewlineIfPresent()
	return &ast.ReturnStmt{Value: val, P: ast.Pos{Line: retTok.Line, Column: retTok.Column}}
}

func (p *Parser) parseCall() ast.Statement {
	nameTok := p.next() // consume ident
	var args []ast.Expr
	for !p.check(token.NEWLINE) && !p.isAtEnd() {
		args = append(args, p.parseExpression(0))
		if p.check(token.NEWLINE) {
			break
		}
	}
	p.consumeNewlineIfPresent()
	return &ast.CallStmt{Name: nameTok.Literal, Args: args, P: ast.Pos{Line: nameTok.Line, Column: nameTok.Column}}
}

func (p *Parser) parseIf() ast.Statement {
	ifTok := p.next() // consume 'if'
	cond := p.parseCondition()
	if !p.check(token.NEWLINE) {
		p.errors = append(p.errors, fmt.Errorf("expected newline after if condition"))
	}
	p.consumeNewlineIfPresent()
	thenBlock := p.parseBlock(token.ELSE, token.END)
	var elseBlock []ast.Statement
	if p.check(token.ELSE) {
		p.next()
		p.consumeNewlineIfPresent()
		elseBlock = p.parseBlock(token.END)
	}
	if !p.check(token.END) {
		p.errors = append(p.errors, fmt.Errorf("expected end to close if"))
	} else {
		p.next() // consume end
	}
	p.consumeNewlineIfPresent()
	return &ast.IfStmt{Cond: cond, Then: thenBlock, Else: elseBlock, P: ast.Pos{Line: ifTok.Line, Column: ifTok.Column}}
}

func (p *Parser) parseFor() ast.Statement {
	forTok := p.next() // consume 'for'
	iterTok, ok := p.expect(token.IDENT)
	if !ok {
		p.errors = append(p.errors, fmt.Errorf("expected identifier after for"))
		return nil
	}
	if !p.check(token.IN) {
		p.errors = append(p.errors, fmt.Errorf("expected in after for variable"))
		return nil
	}
	p.next() // consume 'in'
	start := p.parseExpression(0)
	if !p.check(token.DOTDOT) {
		p.errors = append(p.errors, fmt.Errorf("expected .. in for range"))
		return nil
	}
	p.next() // consume '..'
	end := p.parseExpression(0)
	if !p.check(token.NEWLINE) {
		p.errors = append(p.errors, fmt.Errorf("expected newline after for header"))
	}
	p.consumeNewlineIfPresent()
	body := p.parseBlock(token.END)
	if !p.check(token.END) {
		p.errors = append(p.errors, fmt.Errorf("expected end to close for"))
	} else {
		p.next()
	}
	p.consumeNewlineIfPresent()
	return &ast.ForStmt{Var: iterTok.Literal, Start: start, End: end, Body: body, P: ast.Pos{Line: forTok.Line, Column: forTok.Column}}
}

func (p *Parser) parseWhile() ast.Statement {
	whileTok := p.next() // consume 'while'
	cond := p.parseExpression(0)
	if !p.check(token.NEWLINE) {
		p.errors = append(p.errors, fmt.Errorf("expected newline after while condition"))
	}
	p.consumeNewlineIfPresent()
	body := p.parseBlock(token.END)
	if !p.check(token.END) {
		p.errors = append(p.errors, fmt.Errorf("expected end to close while"))
	} else {
		p.next()
	}
	p.consumeNewlineIfPresent()
	return &ast.WhileStmt{Cond: cond, Body: body, P: ast.Pos{Line: whileTok.Line, Column: whileTok.Column}}
}

func (p *Parser) parseFn() ast.Statement {
	fnTok := p.next() // consume 'fn'
	nameTok, ok := p.expect(token.IDENT)
	if !ok {
		p.errors = append(p.errors, fmt.Errorf("expected function name"))
		return nil
	}
	var params []string
	for p.check(token.IDENT) {
		tok := p.next()
		params = append(params, tok.Literal)
	}
	if !p.check(token.NEWLINE) {
		p.errors = append(p.errors, fmt.Errorf("expected newline after fn signature"))
	}
	p.consumeNewlineIfPresent()
	body := p.parseBlock(token.END)
	if !p.check(token.END) {
		p.errors = append(p.errors, fmt.Errorf("expected end to close function"))
	} else {
		p.next()
	}
	p.consumeNewlineIfPresent()
	return &ast.FnDecl{Name: nameTok.Literal, Params: params, Body: body, P: ast.Pos{Line: fnTok.Line, Column: fnTok.Column}}
}

func (p *Parser) parseBreak() ast.Statement {
	brTok := p.next() // consume 'break'
	p.consumeNewlineIfPresent()
	return &ast.BreakStmt{P: ast.Pos{Line: brTok.Line, Column: brTok.Column}}
}

func (p *Parser) parseContinue() ast.Statement {
	ctTok := p.next() // consume 'continue'
	p.consumeNewlineIfPresent()
	return &ast.ContinueStmt{P: ast.Pos{Line: ctTok.Line, Column: ctTok.Column}}
}

func (p *Parser) parseCondition() ast.Condition {
	if !p.check(token.EXISTS) {
		p.errors = append(p.errors, fmt.Errorf("expected exists condition"))
		return nil
	}
	condTok := p.next() // consume exists
	path := p.parseExpression(0)
	return &ast.ExistsCond{Path: path, P: ast.Pos{Line: condTok.Line, Column: condTok.Column}}
}

func (p *Parser) parseBlock(terminators ...token.Type) []ast.Statement {
	var stmts []ast.Statement
	for !p.isAtEnd() {
		if p.check(token.NEWLINE) {
			p.next()
			continue
		}
		for _, term := range terminators {
			if p.check(term) {
				return stmts
			}
		}
		s := p.parseStatement()
		if s != nil {
			stmts = append(stmts, s)
		} else {
			p.synchronize()
		}
	}
	return stmts
}
