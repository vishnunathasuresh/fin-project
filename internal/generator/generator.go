package generator

import (
	"fmt"

	"github.com/vishnunathasuresh/fin-project/internal/ast"
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
	g.ctx.emitLine("setlocal EnableDelayedExpansion")

	var fns []*ast.FnDecl
	for _, stmt := range p.Statements {
		if fn, ok := stmt.(*ast.FnDecl); ok {
			fns = append(fns, fn)
			continue
		}
		if err := g.emitTopLevel(stmt); err != nil {
			return "", err
		}
	}

	for _, fn := range fns {
		if err := g.emitFunction(fn); err != nil {
			return "", err
		}
	}

	g.ctx.emitLine("endlocal")
	return g.ctx.String(), nil
}

func (g *BatchGenerator) emitTopLevel(stmt ast.Statement) error {
	return g.emitStmt(stmt)
}

func (g *BatchGenerator) emitFunction(fn *ast.FnDecl) error {
	return lowerFnDecl(g.ctx, fn, g.emitStmt)
}

// emitStmt lowers a statement; returns an error for unsupported nodes.
func (g *BatchGenerator) emitStmt(stmt ast.Statement) error {
	if stmt == nil {
		return errUnsupportedStmt(ast.Pos{}, stmt)
	}
	switch s := stmt.(type) {
	case *ast.EchoStmt:
		lowerEchoStmt(g.ctx, s)
	case *ast.RunStmt:
		lowerRunStmt(g.ctx, s)
	case *ast.SetStmt:
		lowerSetStmt(g.ctx, s)
	case *ast.AssignStmt:
		lowerAssignStmt(g.ctx, s)
	case *ast.IfStmt:
		return lowerIfStmt(g.ctx, s, g.emitStmt)
	case *ast.ForStmt:
		return lowerForStmt(g.ctx, s, g.emitStmt)
	case *ast.WhileStmt:
		return lowerWhileStmt(g.ctx, s, g.emitStmt)
	case *ast.CallStmt:
		lowerCallStmt(g.ctx, s)
	case *ast.ReturnStmt:
		if err := lowerReturnStmt(g.ctx, s); err != nil {
			return err
		}
	case *ast.BreakStmt:
		return lowerBreakStmt(g.ctx, s)
	case *ast.ContinueStmt:
		return lowerContinueStmt(g.ctx, s)
	case *ast.FnDecl:
		return errFunctionNotLifted(s.Pos(), s.Name)
	default:
		return errUnsupportedStmt(s.Pos(), stmt)
	}
	return nil
}

func lowerCondition(c ast.Expr) string {
	switch cond := c.(type) {
	case *ast.ExistsCond:
		return fmt.Sprintf("exist %s", lowerExpr(cond.Path))
	default:
		return lowerExpr(cond)
	}
}
