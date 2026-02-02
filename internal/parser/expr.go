package parser

import (
	"fmt"

	"github.com/vishnunathasuresh/fin-project/internal/ast"
	"github.com/vishnunathasuresh/fin-project/internal/diagnostics"
	"github.com/vishnunathasuresh/fin-project/internal/token"
)

// Pratt parser implementation for expressions.
// parseExpression accepts a precedence threshold and returns the parsed expression.

// precedences maps infix token types to their binding power.
var precedences = map[token.Type]int{
	// lowest to highest
	token.OR:   1,
	token.AND:  2,
	token.EQ:   3,
	token.NEQ:  3,
	token.PLUS: 4, token.MINUS: 4,
	token.STAR: 5, token.SLASH: 5,
	token.POWER:    6,
	token.DOT:      7,
	token.LBRACKET: 7, // index has high precedence
	token.LPAREN:   8, // function call has highest precedence
}

var prefixParseFns map[token.Type]prefixParseFn

func init() {
	prefixParseFns = map[token.Type]prefixParseFn{
		token.IDENT:     parseIdent,
		token.STRING:    parseString,
		token.NUMBER:    parseNumber,
		token.TRUE:      parseBool,
		token.FALSE:     parseBool,
		token.MINUS:     parseUnary,
		token.BANG:      parseUnary,
		token.LPAREN:    parseGrouped,
		token.LBRACKET:  parseList,
		token.LBRACE:    parseMap,
		token.CMD_START: parseCommand,
	}
}

type prefixParseFn func(*Parser) ast.Expr
type infixParseFn func(*Parser, ast.Expr) ast.Expr

// parseExpression implements Pratt parsing using prefix/infix functions.
func (p *Parser) parseExpression(precedence int) ast.Expr {
	prefix := prefixParseFns[p.current().Type]
	if prefix == nil {
		p.reportError(p.currentPos(), diagnostics.ErrUnexpectedToken, fmt.Sprintf("no prefix parse function for %s", p.current().Type))
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

func (p *Parser) infixFn(t token.Type) infixParseFn {
	switch t {
	case token.PLUS, token.MINUS, token.STAR, token.SLASH:
		return parseBinary
	case token.POWER:
		return parseBinary
	case token.EQ, token.NEQ:
		return parseBinary
	case token.OR, token.AND:
		return parseBinary
	case token.LBRACKET:
		return parseIndex
	case token.DOT:
		return parseProperty
	case token.LPAREN:
		return parseCallExpr
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

func parseExists(p *Parser) ast.Expr {
	// no longer supported in Fin v2
	return nil
}

func parseUnary(p *Parser) ast.Expr {
	tok := p.next()
	const prefixPrecedence = 9 // higher than power and multiplicative to bind unary tightly
	right := p.parseExpression(prefixPrecedence)
	return &ast.UnaryExpr{Op: tok.Literal, Right: right, P: ast.Pos{Line: tok.Line, Column: tok.Column}}
}

func parseGrouped(p *Parser) ast.Expr {
	p.next() // consume '('
	expr := p.parseExpression(0)
	if !p.check(token.RPAREN) {
		p.reportError(p.currentPos(), diagnostics.ErrSyntax, "expected )")
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
		p.reportError(p.currentPos(), diagnostics.ErrSyntax, "expected , or ] in list")
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
			p.reportError(p.currentPos(), diagnostics.ErrSyntax, "expected map key ident")
			break
		}
		keyTok := p.next()
		if !p.check(token.COLON) {
			p.reportError(p.currentPos(), diagnostics.ErrSyntax, "expected : after map key")
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
		p.reportError(p.currentPos(), diagnostics.ErrSyntax, "expected , or } in map")
		break
	}
	return &ast.MapLit{Pairs: pairs, P: ast.Pos{Line: mTok.Line, Column: mTok.Column}}
}

// ---- infix parse functions ----

func parseBinary(p *Parser, left ast.Expr) ast.Expr {
	opTok := p.current()
	opPrec := p.currentPrecedence()
	p.next() // consume operator
	// Exponentiation is right-associative; other operators are left-associative.
	rightPrec := opPrec
	if opTok.Type == token.POWER {
		rightPrec = opPrec - 1
	}
	right := p.parseExpression(rightPrec)
	return &ast.BinaryExpr{Left: left, Op: opTok.Literal, Right: right, P: ast.Pos{Line: opTok.Line, Column: opTok.Column}}
}

func parseIndex(p *Parser, left ast.Expr) ast.Expr {
	lTok := p.current()
	p.next() // consume '['
	index := p.parseExpression(0)
	if !p.check(token.RBRACKET) {
		p.reportError(p.currentPos(), diagnostics.ErrSyntax, "expected ] after index expression")
	} else {
		p.next()
	}
	return &ast.IndexExpr{Left: left, Index: index, P: ast.Pos{Line: lTok.Line, Column: lTok.Column}}
}

func parseProperty(p *Parser, left ast.Expr) ast.Expr {
	dotTok := p.current()
	p.next() // consume '.'
	if !p.check(token.IDENT) {
		p.reportError(p.currentPos(), diagnostics.ErrSyntax, "expected property name after .")
		return left
	}
	nameTok := p.next()
	return &ast.PropertyExpr{Object: left, Field: nameTok.Literal, P: ast.Pos{Line: dotTok.Line, Column: dotTok.Column}}
}

// parseCommand parses a command literal <...> as CommandLit.
func parseCommand(p *Parser) ast.Expr {
	startTok := p.next() // consume CMD_START
	if !p.check(token.CMD_TEXT) {
		p.reportError(p.currentPos(), diagnostics.ErrSyntax, "expected command text")
		return &ast.CommandLit{Text: "", P: ast.Pos{Line: startTok.Line, Column: startTok.Column}}
	}
	textTok := p.next()
	if !p.check(token.CMD_END) {
		p.reportError(p.currentPos(), diagnostics.ErrSyntax, "expected '>' to close command literal")
	} else {
		p.next()
	}
	return &ast.CommandLit{Text: textTok.Literal, P: ast.Pos{Line: startTok.Line, Column: startTok.Column}}
}

// parseCallExpr parses a function call expression: name(args, key=value, ...)
func parseCallExpr(p *Parser, callee ast.Expr) ast.Expr {
	lpTok := p.current()
	p.next() // consume '('

	var args []ast.Expr
	var namedArgs []ast.NamedArg

	// Parse arguments (both positional and named)
	for !p.check(token.RPAREN) && !p.isAtEnd() {
		// Check if this is a named argument by looking ahead: ident = value
		if p.check(token.IDENT) {
			// Peek ahead to see if there's an '=' after the identifier
			nextPos := p.pos + 1
			if nextPos < len(p.tokens) && p.tokens[nextPos].Type == token.ASSIGN {
				// This is a named argument
				nameTok := p.next() // consume identifier
				p.next()            // consume '='
				val := p.parseExpression(0)
				namedArgs = append(namedArgs, ast.NamedArg{
					Name:  nameTok.Literal,
					Value: val,
					P:     ast.Pos{Line: nameTok.Line, Column: nameTok.Column},
				})
				// After a named argument, require comma or closing paren.
				if p.check(token.COMMA) {
					p.next()
					continue
				}
				if !p.check(token.RPAREN) {
					p.reportError(p.currentPos(), diagnostics.ErrSyntax, "expected , or ) after named argument")
				}
				continue
			}
		}

		// Positional argument path
		args = append(args, p.parseExpression(0))

		if p.check(token.COMMA) {
			p.next() // consume ','
		} else if !p.check(token.RPAREN) {
			p.reportError(p.currentPos(), diagnostics.ErrSyntax, "expected , or ) in function call arguments")
			break
		}
	}

	if !p.check(token.RPAREN) {
		p.reportError(p.currentPos(), diagnostics.ErrSyntax, "expected ) after function call arguments")
	} else {
		p.next() // consume ')'
	}

	return &ast.CallExpr{
		Callee:    callee,
		Args:      args,
		NamedArgs: namedArgs,
		P:         ast.Pos{Line: lpTok.Line, Column: lpTok.Column},
	}
}
