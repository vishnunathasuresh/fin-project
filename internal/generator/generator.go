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
	label := mangleFunc(fn.Name)
	// Function body
	g.ctx.emitLine("goto :eof")
	g.ctx.emitLine(":" + label)
	g.ctx.emitLine("setlocal")
	for i, p := range fn.Params {
		g.ctx.emitLine(fmt.Sprintf("set %s=%%%d", p, i+1))
	}
	g.ctx.pushIndent()
	for _, stmt := range fn.Body {
		g.emitStmt(stmt)
	}
	g.ctx.popIndent()
	g.ctx.emitLine("endlocal")
	g.ctx.emitLine("goto :eof")
}

// emitStmt lowers a statement; currently a stub to maintain compilation until lowering is implemented.
func (g *BatchGenerator) emitStmt(stmt ast.Statement) {
	switch s := stmt.(type) {
	case *ast.EchoStmt:
		g.ctx.emitLine("echo " + lowerExpr(s.Value))
	case *ast.RunStmt:
		g.ctx.emitLine(lowerExpr(s.Command))
	case *ast.SetStmt:
		lowerSetStmt(g.ctx, s)
	case *ast.IfStmt:
		cond := lowerCondition(s.Cond)
		g.ctx.emitLine(fmt.Sprintf("if %s (", cond))
		g.ctx.pushIndent()
		for _, inner := range s.Then {
			g.emitStmt(inner)
		}
		g.ctx.popIndent()
		if len(s.Else) > 0 {
			g.ctx.emitLine(") else (")
			g.ctx.pushIndent()
			for _, inner := range s.Else {
				g.emitStmt(inner)
			}
			g.ctx.popIndent()
		}
		g.ctx.emitLine(")")
	case *ast.ForStmt:
		start := lowerExpr(s.Start)
		end := lowerExpr(s.End)
		g.ctx.emitLine(fmt.Sprintf("for /L %%"+s.Var+" in (%s,1,%s) do (", start, end))
		g.ctx.pushIndent()
		for _, inner := range s.Body {
			g.emitStmt(inner)
		}
		g.ctx.popIndent()
		g.ctx.emitLine(")")
	case *ast.WhileStmt:
		id := g.ctx.NextLabel()
		start := whileStartLabel(id)
		end := whileEndLabel(id)
		g.ctx.emitLine(":" + start)
		cond := lowerCondition(s.Cond)
		g.ctx.emitLine(fmt.Sprintf("if not %s goto %s", cond, end))
		for _, inner := range s.Body {
			g.emitStmt(inner)
		}
		g.ctx.emitLine(fmt.Sprintf("goto %s", start))
		g.ctx.emitLine(":" + end)
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
