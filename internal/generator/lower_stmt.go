package generator

import (
	"fmt"
	"strings"

	"github.com/vishnunathasuresh/fin-project/internal/ast"
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
	case *ast.IndexExpr:
		// Index access depends on whether index is literal or variable
		left, ok := v.Left.(*ast.IdentExpr)
		if !ok {
			// Fallback for complex expressions
			base := trimPercents(lowerExpr(v.Left))
			idx := trimPercents(lowerExpr(v.Index))
			ctx.emitLine(fmt.Sprintf("call set %s=%%%%!%s!_!%s!%%%%", s.Name, base, idx))
		} else {
			base := left.Name
			switch idxExpr := v.Index.(type) {
			case *ast.NumberLit:
				// Literal index: direct access with delayed expansion
				ctx.emitLine(fmt.Sprintf("set %s=!%s_%s!", s.Name, base, idxExpr.Value))
			case *ast.IdentExpr:
				// Variable index: need call set for double delayed expansion
				ctx.emitLine(fmt.Sprintf("call set %s=%%%%!%s!_!%s!%%%%", s.Name, base, idxExpr.Name))
			default:
				// Complex index expression
				idx := trimPercents(lowerExpr(v.Index))
				ctx.emitLine(fmt.Sprintf("call set %s=%%%%!%s!_!%s!%%%%", s.Name, base, idx))
			}
		}
	default:
		if isArithmeticExpr(s.Value) {
			ctx.emitLine(fmt.Sprintf("set /a %s=%s", s.Name, lowerExprArithmetic(s.Value)))
		} else {
			ctx.emitLine(fmt.Sprintf("set %s=%s", s.Name, lowerExpr(s.Value)))
		}
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
	case *ast.IndexExpr:
		// Index access depends on whether index is literal or variable
		left, ok := v.Left.(*ast.IdentExpr)
		if !ok {
			// Fallback for complex expressions
			base := trimPercents(lowerExpr(v.Left))
			idx := trimPercents(lowerExpr(v.Index))
			ctx.emitLine(fmt.Sprintf("call set %s=%%%%!%s!_!%s!%%%%", s.Name, base, idx))
		} else {
			base := left.Name
			switch idxExpr := v.Index.(type) {
			case *ast.NumberLit:
				// Literal index: direct access with delayed expansion
				ctx.emitLine(fmt.Sprintf("set %s=!%s_%s!", s.Name, base, idxExpr.Value))
			case *ast.IdentExpr:
				// Variable index: need call set for double delayed expansion
				ctx.emitLine(fmt.Sprintf("call set %s=%%%%!%s!_!%s!%%%%", s.Name, base, idxExpr.Name))
			default:
				// Complex index expression
				idx := trimPercents(lowerExpr(v.Index))
				ctx.emitLine(fmt.Sprintf("call set %s=%%%%!%s!_!%s!%%%%", s.Name, base, idx))
			}
		}
	default:
		if isArithmeticExpr(s.Value) {
			ctx.emitLine(fmt.Sprintf("set /a %s=%s", s.Name, lowerExprArithmetic(s.Value)))
		} else {
			ctx.emitLine(fmt.Sprintf("set %s=%s", s.Name, lowerExpr(s.Value)))
		}
	}
}

func isArithmeticExpr(e ast.Expr) bool {
	switch v := e.(type) {
	case *ast.BinaryExpr:
		switch v.Op {
		case "+", "-", "*", "/", "**":
			return true
		}
	case *ast.UnaryExpr:
		if v.Op == "-" {
			return true
		}
	}
	return false
}

// lowerEchoStmt emits an echo with expression lowering for interpolation.
func lowerEchoStmt(ctx *Context, s *ast.EchoStmt) {
	val := lowerExpr(s.Value)
	// Escape batch special characters in echo output
	val = escapeBatchSpecials(val)
	ctx.emitLine("echo " + val)
}

// escapeBatchSpecials escapes characters that have special meaning in batch commands.
// This includes < > | & which need to be prefixed with ^ to be printed literally.
// Also escapes ! when it appears in a != sequence (not inside variable expansion).
func escapeBatchSpecials(s string) string {
	var b strings.Builder
	inExpand := false
	expandChar := byte(0)

	for i := 0; i < len(s); i++ {
		c := s[i]

		// Track if we're inside a variable expansion
		if c == '!' || c == '%' {
			if inExpand && c == expandChar {
				// End of expansion
				inExpand = false
				expandChar = 0
				b.WriteByte(c)
				continue
			} else if !inExpand {
				// Check if this is the start of a variable expansion (!name!)
				// or a standalone ! character (like in !=)
				if c == '!' {
					// Look ahead to see if this is a variable pattern
					hasClosing := false
					for j := i + 1; j < len(s); j++ {
						if s[j] == '!' {
							hasClosing = true
							break
						}
						// If we hit a space or special char before closing !, not a var
						if s[j] == ' ' || s[j] == '=' || s[j] == '<' || s[j] == '>' {
							break
						}
					}
					if hasClosing && i+1 < len(s) && isIdentStartByte(s[i+1]) {
						// This is a variable expansion
						inExpand = true
						expandChar = c
						b.WriteByte(c)
						continue
					} else {
						// Standalone !, escape it
						b.WriteString("^^!")
						continue
					}
				}
				// Start of % expansion
				inExpand = true
				expandChar = c
			}
			b.WriteByte(c)
			continue
		}

		// Only escape special chars outside of variable expansions
		if !inExpand {
			switch c {
			case '<', '>', '|', '&':
				b.WriteByte('^')
			}
		}
		b.WriteByte(c)
	}
	return b.String()
}

func isIdentStartByte(c byte) bool {
	return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || c == '_'
}

// lowerRunStmt emits a command invocation with expression lowering.
func lowerRunStmt(ctx *Context, s *ast.RunStmt) {
	cmd := lowerExpr(s.Command)
	cmd = strings.TrimSpace(cmd)
	cmd = strings.Trim(cmd, "\"")
	ctx.emitLine(cmd)
}

// lowerIfStmt lowers an if/else statement with proper indentation.
func lowerIfStmt(ctx *Context, s *ast.IfStmt, emit func(ast.Statement) error) error {
	if b, ok := s.Cond.(*ast.BinaryExpr); ok {
		leftVal := lowerExpr(b.Left)
		rightVal := lowerExpr(b.Right)

		// Check if this is a numeric comparison operator (<, >, <=, >=)
		if isNumericComparisonOp(b.Op) {
			return lowerIfComparison(ctx, b, s.Then, s.Else, emit)
		}

		// Format operand: if it's already a variable expansion (!x!), use as-is
		// otherwise treat as literal
		formatOperand := func(val string, expr ast.Expr) string {
			switch expr.(type) {
			case *ast.IdentExpr, *ast.PropertyExpr, *ast.IndexExpr:
				// It's a variable reference - lowerExpr already added !...!
				return fmt.Sprintf("\"%s\"", val)
			default:
				// It's a literal value
				return fmt.Sprintf("\"%s\"", val)
			}
		}

		left := formatOperand(leftVal, b.Left)
		right := formatOperand(rightVal, b.Right)

		var header string
		switch b.Op {
		case "==":
			header = fmt.Sprintf("if %s==%s (", left, right)
		case "!=":
			header = fmt.Sprintf("if %s NEQ %s (", left, right)
		}
		if header != "" {
			ctx.emitLine(header)
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
	}
	cond := lowerExpr(s.Cond)
	ctx.emitLine(fmt.Sprintf("if \"%s\"==\"true\" (", cond))
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
	ctx.emitLine(fmt.Sprintf("set /a %s=%s", s.Var, startVal))
	ctx.emitRawLine(":" + startLbl)
	ctx.emitLine(fmt.Sprintf("if !%s! GTR %s goto %s", s.Var, endVal, endLbl))
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
	ctx.emitLine(fmt.Sprintf("set /a %s=%s+1", s.Var, s.Var))
	ctx.emitLine(fmt.Sprintf("goto %s", startLbl))
	ctx.emitRawLine(":" + endLbl)
	return nil
}

// lowerWhileStmt lowers a while loop using labels and conditional jumps.

func lowerWhileStmt(ctx *Context, s *ast.WhileStmt, emit func(ast.Statement) error) error {
	id := ctx.NextLabel()
	start := whileStartLabel(id)
	end := whileEndLabel(id)
	ctx.emitRawLine(":" + start)
	switch c := s.Cond.(type) {
	case *ast.ExistsCond:
		cond := lowerCondition(c)
		ctx.emitLine(fmt.Sprintf("if not %s goto %s", cond, end))
	case *ast.BinaryExpr:
		// Handle comparison operators specially since set /a doesn't support them
		if isComparisonOp(c.Op) {
			lowerComparisonCondition(ctx, c, end)
		} else if isBooleanOp(c.Op) {
			// For && and ||, we need more complex handling
			// For now, treat as a general expression that evaluates to true/false
			arith := lowerExprArithmetic(s.Cond)
			temp := mangleTemp("cond", ctx.NextLabel())
			ctx.emitLine(fmt.Sprintf("set /a %s=(%s)", temp, arith))
			ctx.emitLine(fmt.Sprintf("if !%s! equ 0 goto %s", temp, end))
		} else {
			// Arithmetic expression
			arith := lowerExprArithmetic(s.Cond)
			temp := mangleTemp("cond", ctx.NextLabel())
			ctx.emitLine(fmt.Sprintf("set /a %s=(%s)", temp, arith))
			ctx.emitLine(fmt.Sprintf("if !%s! equ 0 goto %s", temp, end))
		}
	case *ast.BoolLit:
		if !c.Value {
			// while false -> immediately exit
			ctx.emitLine(fmt.Sprintf("goto %s", end))
		}
		// while true -> no condition check needed, infinite loop
	default:
		// General expression - try to evaluate as arithmetic
		arith := lowerExprArithmetic(s.Cond)
		temp := mangleTemp("cond", ctx.NextLabel())
		ctx.emitLine(fmt.Sprintf("set /a %s=(%s)", temp, arith))
		ctx.emitLine(fmt.Sprintf("if !%s! equ 0 goto %s", temp, end))
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
	ctx.emitRawLine(":" + end)
	return nil
}

func isComparisonOp(op string) bool {
	switch op {
	case "<", "<=", ">", ">=", "==", "!=":
		return true
	}
	return false
}

func isNumericComparisonOp(op string) bool {
	switch op {
	case "<", "<=", ">", ">=":
		return true
	}
	return false
}

func isBooleanOp(op string) bool {
	return op == "&&" || op == "||"
}

// lowerComparisonCondition handles comparison expressions for while/if conditions
func lowerComparisonCondition(ctx *Context, c *ast.BinaryExpr, endLabel string) {
	left := lowerExprArithmetic(c.Left)
	right := lowerExprArithmetic(c.Right)

	// We need to compute left and right if they're complex expressions
	leftTemp := ""
	rightTemp := ""
	leftExpr := c.Left
	rightExpr := c.Right

	// Check if left or right contain arithmetic that needs pre-computation
	if needsPreCompute(c.Left) {
		leftTemp = mangleTemp("left", ctx.NextLabel())
		ctx.emitLine(fmt.Sprintf("set /a %s=%s", leftTemp, left))
		left = leftTemp
		leftExpr = nil // Mark as temp variable
	}
	if needsPreCompute(c.Right) {
		rightTemp = mangleTemp("right", ctx.NextLabel())
		ctx.emitLine(fmt.Sprintf("set /a %s=%s", rightTemp, right))
		right = rightTemp
		rightExpr = nil // Mark as temp variable
	}

	// Wrap variables in !...! but not literals
	formatForCompare := func(val string, expr ast.Expr) string {
		// If expr is nil, it's a temp variable we created
		if expr == nil {
			return fmt.Sprintf("!%s!", val)
		}
		switch expr.(type) {
		case *ast.IdentExpr, *ast.PropertyExpr, *ast.IndexExpr:
			return fmt.Sprintf("!%s!", val)
		default:
			return val
		}
	}

	leftCmp := formatForCompare(left, leftExpr)
	rightCmp := formatForCompare(right, rightExpr)

	// Generate the comparison using if command
	// Note: We jump to end if condition is FALSE (to exit loop)
	var cmp string
	switch c.Op {
	case "<":
		// if NOT (left < right) goto end  =>  if left >= right goto end  =>  if left GEQ right goto end
		cmp = fmt.Sprintf("if %s GEQ %s goto %s", leftCmp, rightCmp, endLabel)
	case "<=":
		cmp = fmt.Sprintf("if %s GTR %s goto %s", leftCmp, rightCmp, endLabel)
	case ">":
		cmp = fmt.Sprintf("if %s LEQ %s goto %s", leftCmp, rightCmp, endLabel)
	case ">=":
		cmp = fmt.Sprintf("if %s LSS %s goto %s", leftCmp, rightCmp, endLabel)
	case "==":
		cmp = fmt.Sprintf("if %s NEQ %s goto %s", leftCmp, rightCmp, endLabel)
	case "!=":
		cmp = fmt.Sprintf("if %s EQU %s goto %s", leftCmp, rightCmp, endLabel)
	}
	ctx.emitLine(cmp)
}

func needsPreCompute(e ast.Expr) bool {
	switch v := e.(type) {
	case *ast.NumberLit, *ast.IdentExpr:
		return false
	case *ast.BinaryExpr:
		return true
	case *ast.UnaryExpr:
		return needsPreCompute(v.Right)
	default:
		return true
	}
}

// lowerIfComparison handles if statements with comparison operators (<, >, <=, >=, ==, !=).
func lowerIfComparison(ctx *Context, c *ast.BinaryExpr, thenBlock, elseBlock []ast.Statement, emit func(ast.Statement) error) error {
	left := lowerExprArithmetic(c.Left)
	right := lowerExprArithmetic(c.Right)

	// Pre-compute complex expressions
	if needsPreCompute(c.Left) {
		leftTemp := mangleTemp("left", ctx.NextLabel())
		ctx.emitLine(fmt.Sprintf("set /a %s=%s", leftTemp, left))
		left = leftTemp
	}
	if needsPreCompute(c.Right) {
		rightTemp := mangleTemp("right", ctx.NextLabel())
		ctx.emitLine(fmt.Sprintf("set /a %s=%s", rightTemp, right))
		right = rightTemp
	}

	// Map Fin operators to batch comparison operators
	var batchOp string
	switch c.Op {
	case "<":
		batchOp = "LSS"
	case "<=":
		batchOp = "LEQ"
	case ">":
		batchOp = "GTR"
	case ">=":
		batchOp = "GEQ"
	case "==":
		batchOp = "EQU"
	case "!=":
		batchOp = "NEQ"
	}

	// Wrap variables in !...! but not literals
	formatForCompare := func(val string, expr ast.Expr) string {
		switch expr.(type) {
		case *ast.IdentExpr, *ast.PropertyExpr, *ast.IndexExpr:
			return fmt.Sprintf("!%s!", val)
		default:
			return val
		}
	}

	leftCmp := formatForCompare(left, c.Left)
	rightCmp := formatForCompare(right, c.Right)

	ctx.emitLine(fmt.Sprintf("if %s %s %s (", leftCmp, batchOp, rightCmp))
	ctx.pushIndent()
	for _, stmt := range thenBlock {
		if err := emit(stmt); err != nil {
			return err
		}
	}
	ctx.popIndent()
	if len(elseBlock) > 0 {
		ctx.emitLine(") else (")
		ctx.pushIndent()
		for _, stmt := range elseBlock {
			if err := emit(stmt); err != nil {
				return err
			}
		}
		ctx.popIndent()
	}
	ctx.emitLine(")")
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
	ctx.emitLine("setlocal EnableDelayedExpansion")
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
	specials := "^&|><()\""
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
