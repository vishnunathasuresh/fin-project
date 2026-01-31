package generator

import (
	"fmt"
	"strings"

	"github.com/vishnunathasuresh/fin-project/internal/ir"
)

// IRBatchGenerator emits Windows Batch code from validated IR.
type IRBatchGenerator struct {
	ctx *Context
}

// NewIRBatchGenerator constructs an IR-based batch generator with fresh context.
func NewIRBatchGenerator() *IRBatchGenerator {
	return &IRBatchGenerator{ctx: NewContext()}
}

// Generate emits batch code for the provided IR program.
func (g *IRBatchGenerator) Generate(p *ir.Program) (string, error) {
	if p == nil {
		return "", nil
	}

	g.ctx.emitLine("@echo off")
	g.ctx.emitLine("setlocal EnableDelayedExpansion")

	// Emit functions
	for _, fn := range p.Functions {
		if err := g.emitFunction(fn); err != nil {
			return "", err
		}
	}

	g.ctx.emitLine("endlocal")
	return g.ctx.String(), nil
}

func (g *IRBatchGenerator) emitFunction(fn *ir.Function) error {
	if fn == nil {
		return fmt.Errorf("nil function")
	}

	// Main entry point has no label
	if fn.Name != "main" && fn.Name != "" {
		g.ctx.emitLine("")
		g.ctx.emitLine(fmt.Sprintf(":%s", fn.Name))
	}

	for _, stmt := range fn.Body {
		if err := g.emitStmt(stmt); err != nil {
			return err
		}
	}

	// Non-main functions need explicit return
	if fn.Name != "main" && fn.Name != "" {
		g.ctx.emitLine("goto :eof")
	}

	return nil
}

func (g *IRBatchGenerator) emitStmt(stmt ir.Stmt) error {
	if stmt == nil {
		return nil
	}

	switch s := stmt.(type) {
	case *ir.DeclStmt:
		return g.emitDeclStmt(s)
	case *ir.AssignStmt:
		return g.emitAssignStmt(s)
	case *ir.IfStmt:
		return g.emitIfStmt(s)
	case *ir.ForStmt:
		return g.emitForStmt(s)
	case *ir.WhileStmt:
		return g.emitWhileStmt(s)
	case *ir.RunStmt:
		return g.emitRunStmt(s)
	case *ir.ReturnStmt:
		return g.emitReturnStmt(s)
	case *ir.BreakStmt:
		g.ctx.emitLine("goto :break")
		return nil
	case *ir.ContinueStmt:
		g.ctx.emitLine("goto :continue")
		return nil
	default:
		return fmt.Errorf("unsupported IR statement type: %T", stmt)
	}
}

func (g *IRBatchGenerator) emitDeclStmt(s *ir.DeclStmt) error {
	if s.Init == nil {
		return nil
	}

	switch v := s.Init.(type) {
	case *ir.ListLit:
		for i, el := range v.Elements {
			val := g.emitExpr(el)
			g.ctx.emitLine(fmt.Sprintf("set %s_%d=%s", s.Name, i, val))
		}
		g.ctx.emitLine(fmt.Sprintf("set %s_len=%d", s.Name, len(v.Elements)))
	case *ir.MapLit:
		for i, key := range v.Keys {
			keyStr := g.emitExpr(key)
			valStr := g.emitExpr(v.Values[i])
			g.ctx.emitLine(fmt.Sprintf("set %s_%s=%s", s.Name, trimQuotes(keyStr), valStr))
		}
	default:
		val := g.emitExpr(s.Init)
		if isArithmeticIRExpr(s.Init) {
			g.ctx.emitLine(fmt.Sprintf("set /a %s=%s", s.Name, val))
		} else {
			g.ctx.emitLine(fmt.Sprintf("set %s=%s", s.Name, val))
		}
	}
	return nil
}

func (g *IRBatchGenerator) emitAssignStmt(s *ir.AssignStmt) error {
	if s.Value == nil {
		return nil
	}

	val := g.emitExpr(s.Value)
	if isArithmeticIRExpr(s.Value) {
		g.ctx.emitLine(fmt.Sprintf("set /a %s=%s", s.Name, val))
	} else {
		g.ctx.emitLine(fmt.Sprintf("set %s=%s", s.Name, val))
	}
	return nil
}

func (g *IRBatchGenerator) emitIfStmt(s *ir.IfStmt) error {
	cond := g.emitCondition(s.Cond)
	g.ctx.emitLine(fmt.Sprintf("if %s (", cond))
	g.ctx.indent++

	for _, stmt := range s.Then {
		if err := g.emitStmt(stmt); err != nil {
			return err
		}
	}

	g.ctx.indent--
	if len(s.Else) > 0 {
		g.ctx.emitLine(") else (")
		g.ctx.indent++

		for _, stmt := range s.Else {
			if err := g.emitStmt(stmt); err != nil {
				return err
			}
		}

		g.ctx.indent--
	}
	g.ctx.emitLine(")")
	return nil
}

func (g *IRBatchGenerator) emitForStmt(s *ir.ForStmt) error {
	start := g.emitExpr(s.Start)
	end := g.emitExpr(s.End)

	g.ctx.emitLine(fmt.Sprintf("for /L %%%s in (%s,1,%s) do (", s.Var, start, end))
	g.ctx.indent++

	for _, stmt := range s.Body {
		if err := g.emitStmt(stmt); err != nil {
			return err
		}
	}

	g.ctx.indent--
	g.ctx.emitLine(")")
	return nil
}

func (g *IRBatchGenerator) emitWhileStmt(s *ir.WhileStmt) error {
	g.ctx.emitLine(":while_loop")
	cond := g.emitCondition(s.Cond)
	g.ctx.emitLine(fmt.Sprintf("if not %s goto :break", cond))

	for _, stmt := range s.Body {
		if err := g.emitStmt(stmt); err != nil {
			return err
		}
	}

	g.ctx.emitLine("goto :while_loop")
	g.ctx.emitLine(":break")
	return nil
}

func (g *IRBatchGenerator) emitRunStmt(s *ir.RunStmt) error {
	cmd := g.emitExpr(s.Cmd)
	cmd = trimQuotes(cmd)
	g.ctx.emitLine(cmd)
	return nil
}

func (g *IRBatchGenerator) emitReturnStmt(s *ir.ReturnStmt) error {
	if s.Value != nil {
		val := g.emitExpr(s.Value)
		g.ctx.emitLine(fmt.Sprintf("set __retval=%s", val))
	}
	g.ctx.emitLine("goto :eof")
	return nil
}

func (g *IRBatchGenerator) emitExpr(expr ir.Expr) string {
	return g.emitExprWithContext(expr, false)
}

func (g *IRBatchGenerator) emitExprArithmetic(expr ir.Expr) string {
	return g.emitExprWithContext(expr, true)
}

func (g *IRBatchGenerator) emitExprWithContext(expr ir.Expr, arithmetic bool) string {
	if expr == nil {
		return ""
	}

	switch e := expr.(type) {
	case *ir.IntLit:
		return fmt.Sprintf("%d", e.Value)
	case *ir.FloatLit:
		return fmt.Sprintf("%f", e.Value)
	case *ir.StringLit:
		return interpolateIRString(e.Value)
	case *ir.BoolLit:
		if e.Value {
			return "true"
		}
		return "false"
	case *ir.Ident:
		if arithmetic {
			return e.Name
		}
		return fmt.Sprintf("!%s!", e.Name)
	case *ir.BinaryOp:
		left := g.emitExprWithContext(e.Left, arithmetic)
		right := g.emitExprWithContext(e.Right, arithmetic)
		return fmt.Sprintf("%s %s %s", left, e.Op, right)
	case *ir.UnaryOp:
		operand := g.emitExprWithContext(e.Expr, arithmetic)
		return fmt.Sprintf("%s%s", e.Op, operand)
	case *ir.CallExpr:
		return fmt.Sprintf("call :%s", e.Func)
	case *ir.CommandLit:
		return e.Command
	case *ir.ListLit:
		var parts []string
		for _, el := range e.Elements {
			parts = append(parts, g.emitExpr(el))
		}
		return strings.Join(parts, ",")
	case *ir.MapLit:
		var parts []string
		for i, key := range e.Keys {
			keyStr := g.emitExpr(key)
			valStr := g.emitExpr(e.Values[i])
			parts = append(parts, fmt.Sprintf("%s=%s", trimQuotes(keyStr), valStr))
		}
		return strings.Join(parts, ",")
	case *ir.IndexExpr:
		base := trimPercentMarks(g.emitExpr(e.Object))
		idx := trimPercentMarks(g.emitExpr(e.Index))
		return fmt.Sprintf("!%s_!%s!!", base, idx)
	case *ir.PropertyExpr:
		base := trimPercentMarks(g.emitExpr(e.Object))
		if arithmetic {
			return fmt.Sprintf("%s_%s", base, e.Property)
		}
		return fmt.Sprintf("!%s_%s!", base, e.Property)
	default:
		return ""
	}
}

func (g *IRBatchGenerator) emitCondition(expr ir.Expr) string {
	// For now, simple condition handling
	return g.emitExpr(expr)
}

func isArithmeticIRExpr(e ir.Expr) bool {
	switch v := e.(type) {
	case *ir.IntLit:
		return true
	case *ir.BinaryOp:
		switch v.Op {
		case "+", "-", "*", "/", "**":
			return true
		}
	case *ir.UnaryOp:
		if v.Op == "-" {
			return true
		}
	}
	return false
}

func trimQuotes(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}

func trimPercentMarks(s string) string {
	if len(s) >= 2 {
		if (s[0] == '%' && s[len(s)-1] == '%') || (s[0] == '!' && s[len(s)-1] == '!') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func interpolateIRString(s string) string {
	// Simple string interpolation - just return quoted string for now
	// In full implementation, would handle $var expansions
	return fmt.Sprintf("\"%s\"", s)
}
