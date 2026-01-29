package generator

import (
	"fmt"

	"github.com/vishnunath-suresh/fin-project/internal/ast"
)

// Generator is the public interface for batch code generation.
type Generator interface {
	Generate(p *ast.Program) (string, error)
}

// BatchGenerator emits Windows Batch code from a validated AST.
type BatchGenerator struct {
	ctx *Context
}

// NewBatchGenerator constructs a batch generator with fresh context.
func NewBatchGenerator() *BatchGenerator {
	return &BatchGenerator{ctx: NewContext()}
}

// Generate emits batch code for the provided program.
// Assumes the AST has been semantically validated.
func (g *BatchGenerator) Generate(p *ast.Program) (string, error) {
	if p == nil {
		return "", nil
	}

	g.ctx.emitLine("@echo off")

	var fns []*ast.FnDecl
	for _, stmt := range p.Statements {
		if fn, ok := stmt.(*ast.FnDecl); ok {
			fns = append(fns, fn)
			continue
		}
		g.emitTopLevel(stmt)
	}

	for _, fn := range fns {
		g.emitFunction(fn)
	}

	return g.ctx.String(), nil
}

func (g *BatchGenerator) emitTopLevel(stmt ast.Statement) {
	g.emitStmt(stmt)
}

func (g *BatchGenerator) emitFunction(fn *ast.FnDecl) {
	lowerFnDecl(g.ctx, fn, g.emitStmt)
}

// emitStmt lowers a statement; currently a stub to maintain compilation until lowering is implemented.
func (g *BatchGenerator) emitStmt(stmt ast.Statement) {
	switch s := stmt.(type) {
	case *ast.EchoStmt:
		lowerEchoStmt(g.ctx, s)
	case *ast.RunStmt:
		lowerRunStmt(g.ctx, s)
	case *ast.SetStmt:
		lowerSetStmt(g.ctx, s)
	case *ast.IfStmt:
		lowerIfStmt(g.ctx, s, g.emitStmt)
	case *ast.ForStmt:
		lowerForStmt(g.ctx, s, g.emitStmt)
	case *ast.WhileStmt:
		lowerWhileStmt(g.ctx, s, g.emitStmt)
	case *ast.CallStmt:
		lowerCallStmt(g.ctx, s)
	default:
		// TODO: lower other statements (if/for/while/etc.)
	}
}

func lowerCondition(c ast.Expr) string {
	switch cond := c.(type) {
	case *ast.ExistsCond:
		return fmt.Sprintf("exist %s", lowerExpr(cond.Path))
	default:
		return lowerExpr(cond)
	}
}
