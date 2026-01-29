package generator

import (
    "fmt"

    "github.com/vishnunath-suresh/fin-project/internal/ast"
)

// lowerExpr converts an expression into a batch-safe string fragment.
// It performs no evaluation; it only maps AST nodes to batch syntax.
func lowerExpr(expr ast.Expr) string {
    switch e := expr.(type) {
    case *ast.StringLit:
        return e.Value
    case *ast.NumberLit:
        return e.Value
    case *ast.BoolLit:
        if e.Value {
            return "true"
        }
        return "false"
    case *ast.IdentExpr:
        return fmt.Sprintf("%%%s%%", e.Name)
    case *ast.PropertyExpr:
        base := lowerExpr(e.Object)
        // Property accesses lower to base_field; assume base is an identifier expansion.
        return fmt.Sprintf("%s_%s", trimPercents(base), e.Field)
    case *ast.IndexExpr:
        left := lowerExpr(e.Left)
        idx := lowerExpr(e.Index)
        return fmt.Sprintf("%s_%s", trimPercents(left), idx)
    case *ast.BinaryExpr:
        left := lowerExpr(e.Left)
        right := lowerExpr(e.Right)
        return fmt.Sprintf("%s %s %s", left, e.Op, right)
    case *ast.UnaryExpr:
        return fmt.Sprintf("%s%s", e.Op, lowerExpr(e.Right))
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
    if len(s) >= 2 && s[0] == '%' && s[len(s)-1] == '%' {
        return s[1 : len(s)-1]
    }
    return s
}
