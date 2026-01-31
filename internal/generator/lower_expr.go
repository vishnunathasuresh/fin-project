package generator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/vishnunath-suresh/fin-project/internal/ast"
)

// lowerExpr converts an expression into a batch-safe string fragment.
// It performs no evaluation; it only maps AST nodes to batch syntax.
func lowerExpr(expr ast.Expr) string {
	return lowerExprWithContext(expr, false)
}

// lowerExprArithmetic lowers an expression for use in set /a context.
// Variables in set /a don't need expansion markers.
func lowerExprArithmetic(expr ast.Expr) string {
	return lowerExprWithContext(expr, true)
}

func lowerExprWithContext(expr ast.Expr, arithmetic bool) string {
	switch e := expr.(type) {
	case *ast.StringLit:
		return interpolateString(e.Value)
	case *ast.NumberLit:
		return e.Value
	case *ast.BoolLit:
		if e.Value {
			return "true"
		}
		return "false"
	case *ast.IdentExpr:
		if arithmetic {
			return e.Name
		}
		return fmt.Sprintf("!%s!", e.Name)
	case *ast.PropertyExpr:
		base := trimPercents(lowerExprWithContext(e.Object, arithmetic))
		if arithmetic {
			return fmt.Sprintf("%s_%s", base, e.Field)
		}
		return fmt.Sprintf("!%s_%s!", base, e.Field)
	case *ast.IndexExpr:
		left := trimPercents(lowerExprWithContext(e.Left, false))
		idx := trimPercents(lowerExprWithContext(e.Index, false))
		return fmt.Sprintf("!%s_!%s!!", left, idx)
	case *ast.BinaryExpr:
		left := lowerExprWithContext(e.Left, arithmetic)
		right := lowerExprWithContext(e.Right, arithmetic)
		return fmt.Sprintf("%s %s %s", left, e.Op, right)
	case *ast.UnaryExpr:
		return fmt.Sprintf("%s%s", e.Op, lowerExprWithContext(e.Right, arithmetic))
	case *ast.ListLit:
		// Lists lower as comma-separated literal elements.
		out := ""
		for i, el := range e.Elements {
			if i > 0 {
				out += ","
			}
			out += lowerExpr(el)
		}
		return out
	case *ast.MapLit:
		// Maps lower as key=value pairs comma-separated.
		out := ""
		for i, p := range e.Pairs {
			if i > 0 {
				out += ","
			}
			out += fmt.Sprintf("%s=%s", p.Key, lowerExpr(p.Value))
		}
		return out
	case *ast.ExistsCond:
		return lowerExpr(e.Path)
	default:
		return ""
	}
}

// trimPercents removes leading/trailing % used for identifier expansion.
func trimPercents(s string) string {
	if len(s) >= 2 {
		switch {
		case s[0] == '%' && s[len(s)-1] == '%':
			return s[1 : len(s)-1]
		case s[0] == '!' && s[len(s)-1] == '!':
			return s[1 : len(s)-1]
		}
	}
	return s
}

var identPlaceholder = regexp.MustCompile(`\$[A-Za-z_][A-Za-z0-9_]*`)

// interpolateString replaces $ident with %ident% and supports escaping $$ -> $.
func interpolateString(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); {
		if s[i] == '$' {
			// Escaped dollar
			if i+1 < len(s) && s[i+1] == '$' {
				b.WriteByte('$')
				i += 2
				continue
			}
			// Identifier interpolation
			j := i + 1
			if j < len(s) && isIdentStart(s[j]) {
				j++
				for j < len(s) && isIdentPart(s[j]) {
					j++
				}
				name := s[i+1 : j]
				b.WriteString("%" + name + "%")
				i = j
				continue
			}
		}
		b.WriteByte(s[i])
		i++
	}
	return b.String()
}

func isIdentStart(b byte) bool {
	return (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') || b == '_'
}

func isIdentPart(b byte) bool {
	return isIdentStart(b) || (b >= '0' && b <= '9')
}
