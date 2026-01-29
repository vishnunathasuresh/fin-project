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
