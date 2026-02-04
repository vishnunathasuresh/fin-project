package parser

import (
	"errors"
	"fmt"

	"github.com/vishnunathasuresh/fin-project/internal/ast"
	"github.com/vishnunathasuresh/fin-project/internal/diagnostics"
	"github.com/vishnunathasuresh/fin-project/internal/token"
)

// Parser holds token stream state for recursive-descent parsing.
type Parser struct {
	tokens   []token.Token
	pos      int
	errors   []error
	reporter *diagnostics.Reporter
}

// New creates a parser from a token slice.
func New(tokens []token.Token) *Parser {
	return &Parser{tokens: tokens, pos: 0}
}

// NewWithReporter creates a parser that reports diagnostics while parsing.
func NewWithReporter(tokens []token.Token, reporter *diagnostics.Reporter) *Parser {
	return &Parser{tokens: tokens, pos: 0, reporter: reporter}
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

func (p *Parser) tokenPos(tok token.Token) ast.Pos {
	return ast.Pos{Line: tok.Line, Column: tok.Column}
}

func (p *Parser) currentPos() ast.Pos {
	return p.tokenPos(p.current())
}

func (p *Parser) reportError(pos ast.Pos, code, message string) {
	p.errors = append(p.errors, errors.New(message))
	if p.reporter != nil {
		p.reporter.Error(pos, code, message)
	}
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
		p.reportError(p.tokenPos(tok), diagnostics.ErrSyntax, fmt.Sprintf("illegal token: %s", tok.Literal))
		p.next()
		return nil
	}

	switch tok.Type {
	case token.DEF:
		return p.parseFn()
	case token.RETURN:
		return p.parseReturn()
	case token.IF:
		return p.parseIf()
	case token.FOR:
		return p.parseFor()
	case token.WHILE:
		return p.parseWhile()
	case token.BREAK:
		return p.parseBreak()
	case token.CONTINUE:
		return p.parseContinue()
	case token.IDENT:
		// declaration or assignment
		if next := p.peek(); next.Type == token.DECLARE {
			return p.parseDecl()
		}
		if next := p.peek(); next.Type == token.ASSIGN {
			return p.parseAssign()
		}
		return p.parseCall()
	case token.LPAREN:
		// Might be tuple unpacking: (x, y) := ... or (x, y) = ...
		// Peek ahead to check if it's a tuple pattern
		if p.isTuplePattern() {
			if p.peekAheadFor(token.DECLARE) {
				return p.parseDecl()
			}
			if p.peekAheadFor(token.ASSIGN) {
				return p.parseAssign()
			}
		}
		fallthrough
	default:
		p.reportError(p.tokenPos(tok), diagnostics.ErrUnexpectedToken, fmt.Sprintf("unexpected token: %s", tok.Type))
		p.next()
		return nil
	}
}

func (p *Parser) peek() token.Token {
	if p.pos+1 >= len(p.tokens) {
		return token.Token{Type: token.EOF}
	}
	return p.tokens[p.pos+1]
}

// isTuplePattern checks if the current position starts a tuple pattern like (x, y, z)
func (p *Parser) isTuplePattern() bool {
	if !p.check(token.LPAREN) {
		return false
	}
	// Scan forward to check if it looks like (ident, ident, ...) without assignments
	i := p.pos + 1
	for i < len(p.tokens) {
		if p.tokens[i].Type == token.RPAREN {
			return true // Found closing paren, looks like a tuple
		}
		if p.tokens[i].Type == token.IDENT {
			i++
			// After ident, expect COMMA or RPAREN
			if i < len(p.tokens) {
				if p.tokens[i].Type == token.COMMA {
					i++
					continue
				}
				if p.tokens[i].Type == token.RPAREN {
					return true
				}
			}
		}
		// If we see anything else, it's not a tuple pattern
		return false
	}
	return false
}

// peekAheadFor looks ahead in the token stream to find a specific token type
// Useful for checking if there's a := or = after a tuple pattern
func (p *Parser) peekAheadFor(t token.Type) bool {
	for i := p.pos; i < len(p.tokens); i++ {
		if p.tokens[i].Type == t {
			return true
		}
		// Stop if we hit newline or other statement-ending tokens
		if p.tokens[i].Type == token.NEWLINE || p.tokens[i].Type == token.EOF {
			return false
		}
	}
	return false
}

func (p *Parser) parseAssign() ast.Statement {
	// Parse names: could be "x" or "(x, y, z)"
	var names []string

	if p.check(token.LPAREN) {
		// Tuple unpacking: (x, y, z) = ...
		p.next() // consume '('
		for !p.check(token.RPAREN) && !p.isAtEnd() {
			nameTok, ok := p.expect(token.IDENT)
			if !ok {
				p.reportError(p.currentPos(), diagnostics.ErrSyntax, "expected identifier in tuple")
				return nil
			}
			names = append(names, nameTok.Literal)
			if p.check(token.COMMA) {
				p.next() // consume ','
			} else if !p.check(token.RPAREN) {
				p.reportError(p.currentPos(), diagnostics.ErrSyntax, "expected , or ) in tuple")
				return nil
			}
		}
		if !p.check(token.RPAREN) {
			p.reportError(p.currentPos(), diagnostics.ErrSyntax, "expected ) after tuple")
			return nil
		}
		p.next() // consume ')'
	} else {
		// Single name: x = ...
		nameTok := p.next() // ident
		names = append(names, nameTok.Literal)
	}

	assignTok, ok := p.expect(token.ASSIGN)
	if !ok {
		p.reportError(p.currentPos(), diagnostics.ErrSyntax, "expected '=' after identifier")
		return nil
	}
	val := p.parseExpression(0)
	p.consumeNewlineIfPresent()
	return &ast.AssignStmt{Names: names, Value: val, P: ast.Pos{Line: assignTok.Line, Column: assignTok.Column}}
}

func (p *Parser) parseDecl() ast.Statement {
	// Parse names: could be "x" or "(x, y, z)"
	var names []string

	if p.check(token.LPAREN) {
		// Tuple unpacking: (x, y, z) := ...
		p.next() // consume '('
		for !p.check(token.RPAREN) && !p.isAtEnd() {
			nameTok, ok := p.expect(token.IDENT)
			if !ok {
				p.reportError(p.currentPos(), diagnostics.ErrSyntax, "expected identifier in tuple")
				return nil
			}
			names = append(names, nameTok.Literal)
			if p.check(token.COMMA) {
				p.next() // consume ','
			} else if !p.check(token.RPAREN) {
				p.reportError(p.currentPos(), diagnostics.ErrSyntax, "expected , or ) in tuple")
				return nil
			}
		}
		if !p.check(token.RPAREN) {
			p.reportError(p.currentPos(), diagnostics.ErrSyntax, "expected ) after tuple")
			return nil
		}
		p.next() // consume ')'
	} else {
		// Single name: x := ...
		nameTok := p.next() // ident
		names = append(names, nameTok.Literal)
	}

	declTok, ok := p.expect(token.DECLARE)
	if !ok {
		p.reportError(p.currentPos(), diagnostics.ErrSyntax, "expected ':=' after identifier")
		return nil
	}
	val := p.parseExpression(0)
	p.consumeNewlineIfPresent()
	return &ast.DeclStmt{Names: names, Value: val, P: ast.Pos{Line: declTok.Line, Column: declTok.Column}}
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
	cond := p.parseExpression(0)
	if !p.check(token.NEWLINE) {
		p.reportError(p.currentPos(), diagnostics.ErrSyntax, "expected newline after if condition")
	}
	p.consumeNewlineIfPresent()
	thenBlock := p.parseBlock(token.ELSE, token.EOF)
	var elseBlock []ast.Statement
	if p.check(token.ELSE) {
		p.next()
		p.consumeNewlineIfPresent()
		elseBlock = p.parseBlock(token.EOF)
	}
	p.consumeNewlineIfPresent()
	return &ast.IfStmt{Cond: cond, Then: thenBlock, Else: elseBlock, P: ast.Pos{Line: ifTok.Line, Column: ifTok.Column}}
}

func (p *Parser) parseFor() ast.Statement {
	forTok := p.next() // consume 'for'
	iterTok, ok := p.expect(token.IDENT)
	if !ok {
		p.reportError(p.currentPos(), diagnostics.ErrSyntax, "expected identifier after for")
		return nil
	}
	if !p.check(token.IN) {
		p.reportError(p.currentPos(), diagnostics.ErrSyntax, "expected 'in' after loop variable")
		return nil
	}
	p.next() // consume 'in'
	iterable := p.parseExpression(0)
	if !p.check(token.NEWLINE) {
		p.reportError(p.currentPos(), diagnostics.ErrSyntax, "expected newline after for header")
	}
	p.consumeNewlineIfPresent()
	body := p.parseBlock(token.EOF)
	p.consumeNewlineIfPresent()
	return &ast.ForStmt{Var: iterTok.Literal, Iterable: iterable, Body: body, P: ast.Pos{Line: forTok.Line, Column: forTok.Column}}
}

func (p *Parser) parseWhile() ast.Statement {
	whileTok := p.next() // consume 'while'
	cond := p.parseExpression(0)
	if !p.check(token.NEWLINE) {
		p.reportError(p.currentPos(), diagnostics.ErrSyntax, "expected newline after while condition")
	}
	p.consumeNewlineIfPresent()
	body := p.parseBlock(token.EOF)
	p.consumeNewlineIfPresent()
	return &ast.WhileStmt{Cond: cond, Body: body, P: ast.Pos{Line: whileTok.Line, Column: whileTok.Column}}
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

// parseFn parses: def name(param: type, ...) -> return_type:
//
//	    body...
//	end
func (p *Parser) parseFn() ast.Statement {
	defTok := p.next() // consume 'def'

	// Parse function name
	nameTok, ok := p.expect(token.IDENT)
	if !ok {
		p.reportError(p.currentPos(), diagnostics.ErrSyntax, "expected function name after def")
		return nil
	}

	// Parse parameter list: (param: type, ...)
	if !p.check(token.LPAREN) {
		p.reportError(p.currentPos(), diagnostics.ErrSyntax, "expected ( after function name")
		return nil
	}
	p.next() // consume '('

	params := []ast.Param{}
	for !p.check(token.RPAREN) && !p.isAtEnd() {
		// Parse parameter: name: type
		paramTok, ok := p.expect(token.IDENT)
		if !ok {
			p.reportError(p.currentPos(), diagnostics.ErrSyntax, "expected parameter name")
			return nil
		}

		if !p.check(token.COLON) {
			p.reportError(p.currentPos(), diagnostics.ErrSyntax, "expected : after parameter name")
			return nil
		}
		p.next() // consume ':'

		// Parse parameter type
		typeTok, ok := p.expect(token.IDENT)
		if !ok {
			p.reportError(p.currentPos(), diagnostics.ErrSyntax, "expected parameter type")
			return nil
		}

		params = append(params, ast.Param{
			Name: paramTok.Literal,
			Type: &ast.TypeRef{Name: typeTok.Literal},
			P:    ast.Pos{Line: paramTok.Line, Column: paramTok.Column},
		})

		// Check for comma or end of parameters
		if p.check(token.COMMA) {
			p.next() // consume ','
		} else if !p.check(token.RPAREN) {
			p.reportError(p.currentPos(), diagnostics.ErrSyntax, "expected , or ) in parameter list")
			return nil
		}
	}

	if !p.check(token.RPAREN) {
		p.reportError(p.currentPos(), diagnostics.ErrSyntax, "expected ) after parameters")
		return nil
	}
	p.next() // consume ')'

	// Parse return type: -> return_type
	if !p.check(token.ARROW) {
		p.reportError(p.currentPos(), diagnostics.ErrSyntax, "expected -> after parameters")
		return nil
	}
	p.next() // consume '->'

	returnTok, ok := p.expect(token.IDENT)
	if !ok {
		p.reportError(p.currentPos(), diagnostics.ErrSyntax, "expected return type")
		return nil
	}
	returnType := &ast.TypeRef{Name: returnTok.Literal}

	// Expect ':' and newline
	if !p.check(token.COLON) {
		p.reportError(p.currentPos(), diagnostics.ErrSyntax, "expected : after return type")
		return nil
	}
	p.next() // consume ':'
	p.consumeNewlineIfPresent()

	// Parse function body
	body := p.parseBlock(token.EOF)

	return &ast.FnDecl{
		Name:   nameTok.Literal,
		Params: params,
		Return: returnType,
		Body:   body,
		P:      ast.Pos{Line: defTok.Line, Column: defTok.Column},
	}
}

func (p *Parser) parseBlock(until token.Type, others ...token.Type) []ast.Statement {
	terminators := append([]token.Type{until}, others...)
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
