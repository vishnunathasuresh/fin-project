package generator

import (
	"fmt"
	"strings"

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

func lowerAssignStmt(ctx *Context, s *ast.AssignStmt) {
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
func lowerIfStmt(ctx *Context, s *ast.IfStmt, emit func(ast.Statement) error) error {
	cond := lowerCondition(s.Cond)
	ctx.emitLine(fmt.Sprintf("if %s (", cond))
	ctx.pushIndent()
	for _, inner := range s.Then {
		if err := emit(inner); err != nil {
			return err
		}
	}
	ctx.popIndent()
	if len(s.Else) > 0 {
		ctx.emitLine(") else (")
		ctx.pushIndent()
		for _, inner := range s.Else {
			if err := emit(inner); err != nil {
				return err
			}
		}
		ctx.popIndent()
	}
	ctx.emitLine(")")
	return nil
}

// lowerForStmt lowers a numeric range loop using labels to support break/continue.
func lowerForStmt(ctx *Context, s *ast.ForStmt, emit func(ast.Statement) error) error {
	startVal := lowerExpr(s.Start)
	endVal := lowerExpr(s.End)
	id := ctx.NextLabel()
	startLbl := loopContinueLabel(id)
	endLbl := loopBreakLabel(id)
	ctx.emitLine(fmt.Sprintf("set %s=%s", s.Var, startVal))
	ctx.emitLine(":" + startLbl)
	ctx.emitLine(fmt.Sprintf("if %%%s%% GTR %s goto %s", s.Var, endVal, endLbl))
	ctx.pushLoop(endLbl, startLbl)
	ctx.pushIndent()
	for _, inner := range s.Body {
		if err := emit(inner); err != nil {
			ctx.popIndent()
			ctx.popLoop()
			return err
		}
	}
	ctx.popIndent()
	ctx.popLoop()
	ctx.emitLine(fmt.Sprintf("set /a %s=%%%s%%+1", s.Var, s.Var))
	ctx.emitLine(fmt.Sprintf("goto %s", startLbl))
	ctx.emitLine(":" + endLbl)
	return nil
}

// lowerWhileStmt lowers a while loop using labels and conditional jumps.
func lowerWhileStmt(ctx *Context, s *ast.WhileStmt, emit func(ast.Statement) error) error {
	id := ctx.NextLabel()
	start := whileStartLabel(id)
	end := whileEndLabel(id)
	ctx.emitLine(":" + start)
	switch c := s.Cond.(type) {
	case *ast.ExistsCond:
		cond := lowerCondition(c)
		ctx.emitLine(fmt.Sprintf("if not %s goto %s", cond, end))
	default:
		raw := lowerExpr(s.Cond)
		arith := strings.ReplaceAll(raw, "%", "")
		temp := mangleTemp("cond", ctx.NextLabel())
		ctx.emitLine(fmt.Sprintf("set /a %s=%s", temp, arith))
		ctx.emitLine(fmt.Sprintf("if %%%s%%==0 goto %s", temp, end))
	}
	ctx.pushLoop(end, start)
	for _, inner := range s.Body {
		if err := emit(inner); err != nil {
			ctx.popLoop()
			return err
		}
	}
	ctx.popLoop()
	ctx.emitLine(fmt.Sprintf("goto %s", start))
	ctx.emitLine(":" + end)
	return nil
}

// lowerFnDecl lowers a function declaration to a batch label with parameter mapping.
func lowerFnDecl(ctx *Context, fn *ast.FnDecl, emit func(ast.Statement) error) error {
	label := mangleFunc(fn.Name)
	retLabel := fnReturnLabel(fn.Name)
	retTemp := mangleTemp("ret_"+fn.Name, ctx.NextLabel())
	outVar := fmt.Sprintf("%s_ret", mangleFunc(fn.Name))
	ret := returnTarget{label: retLabel, tempVar: retTemp, outVar: outVar}
	ctx.emitLine("goto :eof")
	ctx.emitLine(":" + label)
	ctx.emitLine("setlocal")
	for i, p := range fn.Params {
		ctx.emitLine(fmt.Sprintf("set %s=%%%d", p, i+1))
	}
	ctx.emitLine(fmt.Sprintf("set %s=", retTemp))
	ctx.pushReturn(ret.label, ret.tempVar, ret.outVar)
	ctx.pushIndent()
	for _, stmt := range fn.Body {
		if err := emit(stmt); err != nil {
			ctx.popReturn()
			ctx.popIndent()
			return err
		}
	}
	ctx.popIndent()
	ctx.popReturn()
	ctx.emitLine(":" + ret.label)
	ctx.emitLine(fmt.Sprintf("endlocal & set %s=%%%s%%", ret.outVar, ret.tempVar))
	ctx.emitLine("goto :eof")
	return nil
}

// lowerCallStmt lowers a function call to a batch call label.
func lowerCallStmt(ctx *Context, s *ast.CallStmt) {
	label := mangleFunc(s.Name)
	var b strings.Builder
	for i, arg := range s.Args {
		if i > 0 {
			b.WriteString(" ")
		}
		lowered := lowerExpr(arg)
		b.WriteString(escapeCallArg(lowered))
	}
	ctx.emitLine(fmt.Sprintf("call :%s %s", label, b.String()))
}

// lowerReturnStmt currently emits a stub; return values are not supported.
func lowerReturnStmt(ctx *Context, s *ast.ReturnStmt) error {
	if s.Value != nil {
		if ret, ok := ctx.currentReturn(); ok {
			ctx.emitLine(fmt.Sprintf("set %s=%s", ret.tempVar, lowerExpr(s.Value)))
			ctx.emitLine("goto " + ret.label)
			return nil
		}
		return errUnsupportedStmt(s.Pos(), s)
	}
	if ret, ok := ctx.currentReturn(); ok {
		ctx.emitLine("goto " + ret.label)
		return nil
	}
	return errUnsupportedStmt(s.Pos(), s)
}

func lowerBreakStmt(ctx *Context, s *ast.BreakStmt) error {
	if labels, ok := ctx.currentLoop(); ok {
		ctx.emitLine("goto " + labels.breakLabel)
		return nil
	}
	return errUnsupportedStmt(s.Pos(), s)
}

func lowerContinueStmt(ctx *Context, s *ast.ContinueStmt) error {
	if labels, ok := ctx.currentLoop(); ok {
		ctx.emitLine("goto " + labels.continueLabel)
		return nil
	}
	return errUnsupportedStmt(s.Pos(), s)
}

// escapeCallArg escapes batch specials and quotes when needed.
func escapeCallArg(arg string) string {
	specials := "^&|><()!"
	needQuote := false
	var b strings.Builder
	for i := 0; i < len(arg); i++ {
		ch := arg[i]
		if ch == ' ' || ch == '\t' {
			needQuote = true
		}
		if strings.ContainsRune(specials, rune(ch)) || ch == '^' {
			b.WriteByte('^')
			needQuote = true
		}
		b.WriteByte(ch)
	}
	res := b.String()
	if needQuote {
		return "\"" + res + "\""
	}
	return res
}
