package generator

import (
	"fmt"

	"github.com/vishnunath-suresh/fin-project/internal/ast"
)

// lowerSetStmt handles lowering of set statements, including lists and maps.
func lowerSetStmt(ctx *Context, s *ast.SetStmt) {
	switch v := s.Value.(type) {
	case *ast.ListLit:
		for i, el := range v.Elements {
			ctx.emitLine(fmt.Sprintf("set %s_%d=%s", s.Name, i, lowerExpr(el)))
		}
		ctx.emitLine(fmt.Sprintf("set %s_len=%d", s.Name, len(v.Elements)))
	case *ast.MapLit:
		for _, p := range v.Pairs {
			ctx.emitLine(fmt.Sprintf("set %s_%s=%s", s.Name, p.Key, lowerExpr(p.Value)))
		}
	default:
		ctx.emitLine(fmt.Sprintf("set %s=%s", s.Name, lowerExpr(s.Value)))
	}
}

// lowerEchoStmt emits an echo with expression lowering for interpolation.
func lowerEchoStmt(ctx *Context, s *ast.EchoStmt) {
	ctx.emitLine("echo " + lowerExpr(s.Value))
}

// lowerRunStmt emits a command invocation with expression lowering.
func lowerRunStmt(ctx *Context, s *ast.RunStmt) {
	ctx.emitLine(lowerExpr(s.Command))
}

// lowerIfStmt lowers an if/else statement with proper indentation.
// emit is used to recursively lower nested statements.
func lowerIfStmt(ctx *Context, s *ast.IfStmt, emit func(ast.Statement)) {
	cond := lowerCondition(s.Cond)
	ctx.emitLine(fmt.Sprintf("if %s (", cond))
	ctx.pushIndent()
	for _, inner := range s.Then {
		emit(inner)
	}
	ctx.popIndent()
	if len(s.Else) > 0 {
		ctx.emitLine(") else (")
		ctx.pushIndent()
		for _, inner := range s.Else {
			emit(inner)
		}
		ctx.popIndent()
	}
	ctx.emitLine(")")
}

// lowerForStmt lowers a numeric range loop.
func lowerForStmt(ctx *Context, s *ast.ForStmt, emit func(ast.Statement)) {
	start := lowerExpr(s.Start)
	end := lowerExpr(s.End)
	ctx.emitLine(fmt.Sprintf("for /L %%"+s.Var+" in (%s,1,%s) do (", start, end))
	ctx.pushIndent()
	for _, inner := range s.Body {
		emit(inner)
	}
	ctx.popIndent()
	ctx.emitLine(")")
}

// lowerWhileStmt lowers a while loop using labels and conditional jumps.
func lowerWhileStmt(ctx *Context, s *ast.WhileStmt, emit func(ast.Statement)) {
	id := ctx.NextLabel()
	start := whileStartLabel(id)
	end := whileEndLabel(id)
	ctx.emitLine(":" + start)
	cond := lowerCondition(s.Cond)
	ctx.emitLine(fmt.Sprintf("if not %s goto %s", cond, end))
	for _, inner := range s.Body {
		emit(inner)
	}
	ctx.emitLine(fmt.Sprintf("goto %s", start))
	ctx.emitLine(":" + end)
}
