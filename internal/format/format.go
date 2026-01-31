package format

import (
	"fmt"
	"strings"

	"github.com/vishnunathasuresh/fin-project/internal/ast"
)

// Format program into canonical Fin source (deterministic, minimal spacing).
func Format(prog *ast.Program) string {
	if prog == nil {
		return ""
	}
	var b strings.Builder
	var prev ast.Statement
	first := true
	for _, stmt := range prog.Statements {
		if stmt == nil {
			continue
		}
		if !first {
			if isFnDecl(prev) && isFnDecl(stmt) {
				b.WriteString("\n\n")
			} else {
				b.WriteByte('\n')
			}
		}
		writeStmt(&b, stmt, 0)
		prev = stmt
		first = false
	}
	return b.String()
}

func isFnDecl(stmt ast.Statement) bool {
	_, ok := stmt.(*ast.FnDecl)
	return ok
}

func writeStmt(b *strings.Builder, stmt ast.Statement, indent int) {
	ind := strings.Repeat("    ", indent)
	switch s := stmt.(type) {
	case *ast.SetStmt:
		fmt.Fprintf(b, "%sset %s %s", ind, s.Name, formatExpr(s.Value))
	case *ast.EchoStmt:
		fmt.Fprintf(b, "%secho %s", ind, formatExpr(s.Value))
	case *ast.RunStmt:
		fmt.Fprintf(b, "%srun %s", ind, formatExpr(s.Command))
	case *ast.CallStmt:
		fmt.Fprintf(b, "%s%s", ind, s.Name)
		for _, a := range s.Args {
			fmt.Fprintf(b, " %s", formatExpr(a))
		}
	case *ast.ReturnStmt:
		if s.Value != nil {
			fmt.Fprintf(b, "%sreturn %s", ind, formatExpr(s.Value))
		} else {
			fmt.Fprintf(b, "%sreturn", ind)
		}
	case *ast.IfStmt:
		fmt.Fprintf(b, "%sif %s\n", ind, formatExpr(s.Cond))
		for i, inner := range s.Then {
			writeStmt(b, inner, indent+1)
			b.WriteByte('\n')
			if i == len(s.Then)-1 && len(s.Else) == 0 {
				// no extra
			}
		}
		if len(s.Else) > 0 {
			fmt.Fprintf(b, "%selse\n", ind)
			for _, inner := range s.Else {
				writeStmt(b, inner, indent+1)
				b.WriteByte('\n')
			}
		}
		fmt.Fprintf(b, "%send", ind)
	case *ast.ForStmt:
		fmt.Fprintf(b, "%sfor %s in %s .. %s\n", ind, s.Var, formatExpr(s.Start), formatExpr(s.End))
		for _, inner := range s.Body {
			writeStmt(b, inner, indent+1)
			b.WriteByte('\n')
		}
		fmt.Fprintf(b, "%send", ind)
	case *ast.WhileStmt:
		fmt.Fprintf(b, "%swhile %s\n", ind, formatExpr(s.Cond))
		for _, inner := range s.Body {
			writeStmt(b, inner, indent+1)
			b.WriteByte('\n')
		}
		fmt.Fprintf(b, "%send", ind)
	case *ast.FnDecl:
		fmt.Fprintf(b, "%sfn %s", ind, s.Name)
		for _, p := range s.Params {
			fmt.Fprintf(b, " %s", p)
		}
		b.WriteByte('\n')
		for _, inner := range s.Body {
			writeStmt(b, inner, indent+1)
			b.WriteByte('\n')
		}
		fmt.Fprintf(b, "%send", ind)
	default:
		fmt.Fprintf(b, "%s# unsupported stmt %T", ind, stmt)
	}
}

func formatExpr(e ast.Expr) string {
	if e == nil {
		return ""
	}
	switch v := e.(type) {
	case *ast.StringLit:
		return v.Value
	case *ast.NumberLit:
		return v.Value
	case *ast.BoolLit:
		if v.Value {
			return "true"
		}
		return "false"
	case *ast.IdentExpr:
		return "$" + v.Name
	case *ast.ListLit:
		parts := make([]string, 0, len(v.Elements))
		for _, el := range v.Elements {
			parts = append(parts, formatExpr(el))
		}
		return "[" + strings.Join(parts, ", ") + "]"
	case *ast.MapLit:
		parts := make([]string, 0, len(v.Pairs))
		for _, p := range v.Pairs {
			parts = append(parts, fmt.Sprintf("%s: %s", p.Key, formatExpr(p.Value)))
		}
		return "{" + strings.Join(parts, ", ") + "}"
	case *ast.IndexExpr:
		return fmt.Sprintf("%s[%s]", formatExpr(v.Left), formatExpr(v.Index))
	case *ast.PropertyExpr:
		return fmt.Sprintf("%s.%s", formatExpr(v.Object), v.Field)
	case *ast.UnaryExpr:
		return fmt.Sprintf("%s%s", v.Op, formatExpr(v.Right))
	case *ast.BinaryExpr:
		return fmt.Sprintf("(%s %s %s)", formatExpr(v.Left), v.Op, formatExpr(v.Right))
	case *ast.ExistsCond:
		return "exists " + formatExpr(v.Path)
	default:
		return "" // fallback
	}
}
