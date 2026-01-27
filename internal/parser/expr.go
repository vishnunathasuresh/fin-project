package parser

import (
	"fmt"

	"github.com/vishnunath-suresh/fin-project/internal/ast"
	"github.com/vishnunath-suresh/fin-project/internal/token"
)

// Pratt parser implementation for expressions.
// parseExpression accepts a precedence threshold and returns the parsed expression.

// precedences maps infix token types to their binding power.
var precedences = map[token.Type]int{
	token.OR:   1,
	token.AND:  2,
	token.EQEQ: 3, token.NOTEQ: 3,
	token.LT: 4, token.LTE: 4, token.GT: 4, token.GTE: 4,
	token.PLUS: 5, token.MINUS: 5,
	token.STAR: 6, token.SLASH: 6,
	token.DOT:      7,
	token.LBRACKET: 7, // index has high precedence
	token.POW:      8, // highest, right-associative
}

type prefixParseFn func(*Parser) ast.Expr
type infixParseFn func(*Parser, ast.Expr) ast.Expr

// parseExpression implements Pratt parsing using prefix/infix functions.
func (p *Parser) parseExpression(precedence int) ast.Expr {
	prefix := p.prefixFn(p.current().Type)
	if prefix == nil {
		p.errors = append(p.errors, fmt.Errorf("no prefix parse function for %s", p.current().Type))
		return nil
	}

	left := prefix(p)

	for !p.isAtEnd() {
		currPrec := p.currentPrecedence()
		if precedence >= currPrec {
			break
		}

		infix := p.infixFn(p.current().Type)
		if infix == nil {
			break
		}

		left = infix(p, left)
	}

	return left
}

func (p *Parser) prefixFn(t token.Type) prefixParseFn {
	switch t {
	case token.IDENT:
		return parseIdent
	case token.NUMBER:
		return parseNumber
	case token.STRING:
		return parseString
	case token.TRUE, token.FALSE:
		return parseBool
	case token.MINUS, token.BANG:
		return parseUnary
	case token.LBRACKET:
		return parseList
	case token.LBRACE:
		return parseMap
	case token.LPAREN:
		return parseGrouped
	default:
		return nil
	}
}

func (p *Parser) infixFn(t token.Type) infixParseFn {
	switch t {
	case token.PLUS, token.MINUS, token.STAR, token.SLASH,
		token.EQEQ, token.NOTEQ,
		token.LT, token.LTE, token.GT, token.GTE,
		token.AND, token.OR,
		token.POW:
		return parseBinary
	case token.LBRACKET:
		return parseIndex
	case token.DOT:
		return parseProperty
	default:
		return nil
	}
}

func (p *Parser) currentPrecedence() int {
	if prec, ok := precedences[p.current().Type]; ok {
		return prec
	}
	return 0
}

// ---- prefix parse functions ----

func parseIdent(p *Parser) ast.Expr {
	tok := p.next()
	return &ast.IdentExpr{Name: tok.Literal, P: ast.Pos{Line: tok.Line, Column: tok.Column}}
}

func parseNumber(p *Parser) ast.Expr {
	tok := p.next()
	return &ast.NumberLit{Value: tok.Literal, P: ast.Pos{Line: tok.Line, Column: tok.Column}}
}

func parseString(p *Parser) ast.Expr {
	tok := p.next()
	return &ast.StringLit{Value: tok.Literal, P: ast.Pos{Line: tok.Line, Column: tok.Column}}
}

func parseBool(p *Parser) ast.Expr {
	tok := p.next()
	val := tok.Type == token.TRUE
	return &ast.BoolLit{Value: val, P: ast.Pos{Line: tok.Line, Column: tok.Column}}
}

func parseUnary(p *Parser) ast.Expr {
	tok := p.next()
	const prefixPrecedence = 7 // higher than multiplicative to bind unary tightly
	right := p.parseExpression(prefixPrecedence)
	return &ast.UnaryExpr{Op: tok.Literal, Right: right, P: ast.Pos{Line: tok.Line, Column: tok.Column}}
}

func parseGrouped(p *Parser) ast.Expr {
	p.next() // consume '('
	expr := p.parseExpression(0)
	if !p.check(token.RPAREN) {
		p.errors = append(p.errors, fmt.Errorf("expected )"))
		return expr
	}
	p.next() // consume ')'
	return expr
}

func parseList(p *Parser) ast.Expr {
	lTok := p.next() // consume '['
	var elems []ast.Expr
	if p.check(token.RBRACKET) {
		p.next()
		return &ast.ListLit{Elements: elems, P: ast.Pos{Line: lTok.Line, Column: lTok.Column}}
	}
	for {
		elem := p.parseExpression(0)
		elems = append(elems, elem)
		if p.check(token.RBRACKET) {
			p.next()
			break
		}
		if p.check(token.COMMA) {
			p.next()
			continue
		}
		p.errors = append(p.errors, fmt.Errorf("expected , or ] in list"))
		break
	}
	return &ast.ListLit{Elements: elems, P: ast.Pos{Line: lTok.Line, Column: lTok.Column}}
}

func parseMap(p *Parser) ast.Expr {
	mTok := p.next() // consume '{'
	var pairs []ast.MapPair
	if p.check(token.RBRACE) {
		p.next()
		return &ast.MapLit{Pairs: pairs, P: ast.Pos{Line: mTok.Line, Column: mTok.Column}}
	}
	for {
		if !p.check(token.IDENT) {
			p.errors = append(p.errors, fmt.Errorf("expected map key ident"))
			break
		}
		keyTok := p.next()
		if !p.check(token.COLON) {
			p.errors = append(p.errors, fmt.Errorf("expected : after map key"))
			break
		}
		p.next() // consume ':'
		val := p.parseExpression(0)
		pairs = append(pairs, ast.MapPair{Key: keyTok.Literal, Value: val, P: ast.Pos{Line: keyTok.Line, Column: keyTok.Column}})

		if p.check(token.RBRACE) {
			p.next()
			break
		}
		if p.check(token.COMMA) {
			p.next()
			continue
		}
		p.errors = append(p.errors, fmt.Errorf("expected , or } in map"))
		break
	}
	return &ast.MapLit{Pairs: pairs, P: ast.Pos{Line: mTok.Line, Column: mTok.Column}}
}

// ---- infix parse functions ----

func parseBinary(p *Parser, left ast.Expr) ast.Expr {
	opTok := p.current()
	opPrec := p.currentPrecedence()
	p.next() // consume operator
	// Exponentiation is right-associative; reduce precedence for RHS to bind right.
	nextPrec := opPrec
	if opTok.Type == token.POW {
		nextPrec = opPrec - 1
	}
	right := p.parseExpression(nextPrec)
	return &ast.BinaryExpr{Left: left, Op: opTok.Literal, Right: right, P: ast.Pos{Line: opTok.Line, Column: opTok.Column}}
}

func parseIndex(p *Parser, left ast.Expr) ast.Expr {
	lTok := p.current()
	p.next() // consume '['
	index := p.parseExpression(0)
	if !p.check(token.RBRACKET) {
		p.errors = append(p.errors, fmt.Errorf("expected ] after index expression"))
	} else {
		p.next()
	}
	return &ast.IndexExpr{Left: left, Index: index, P: ast.Pos{Line: lTok.Line, Column: lTok.Column}}
}

func parseProperty(p *Parser, left ast.Expr) ast.Expr {
	dotTok := p.current()
	p.next() // consume '.'
	if !p.check(token.IDENT) {
		p.errors = append(p.errors, fmt.Errorf("expected property name after ."))
		return left
	}
	nameTok := p.next()
	return &ast.PropertyExpr{Object: left, Field: nameTok.Literal, P: ast.Pos{Line: dotTok.Line, Column: dotTok.Column}}
}
