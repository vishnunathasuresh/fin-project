package ast

import (
	"fmt"
	"strings"
)

// Format returns a human-readable, indented representation of the AST node.
// It includes source positions and is intended for debugging only.
func Format(node Node) string {
	var b strings.Builder
	p := printer{buf: &b}
	p.printNode(node, 0, "")
	return b.String()
}

type printer struct {
	buf *strings.Builder
}

func (p *printer) indent(level int) {
	for i := 0; i < level; i++ {
		p.buf.WriteString("  ")
	}
}

func (p *printer) printNode(n Node, level int, label string) {
	if n == nil {
		p.indent(level)
		if label != "" {
			p.buf.WriteString(label + ": ")
		}
		p.buf.WriteString("<nil>\n")
		return
	}

	p.indent(level)
	if label != "" {
		p.buf.WriteString(label + ": ")
	}

	switch node := n.(type) {
	case *Program:
		fmt.Fprintf(p.buf, "Program @%d:%d\n", node.P.Line, node.P.Column)
		for _, s := range node.Statements {
			p.printNode(s, level+1, "")
		}
	case *AssignStmt:
		fmt.Fprintf(p.buf, "AssignStmt name=%s @%d:%d\n", node.Name, node.P.Line, node.P.Column)
		p.printNode(node.Value, level+1, "value")
	case *CallStmt:
		fmt.Fprintf(p.buf, "CallStmt name=%s @%d:%d\n", node.Name, node.P.Line, node.P.Column)
		for i, arg := range node.Args {
			p.printNode(arg, level+1, fmt.Sprintf("arg[%d]", i))
		}
	case *FnDecl:
		fmt.Fprintf(p.buf, "FnDecl name=%s params=%v @%d:%d\n", node.Name, node.Params, node.P.Line, node.P.Column)
		for _, s := range node.Body {
			p.printNode(s, level+1, "body")
		}
	case *IfStmt:
		fmt.Fprintf(p.buf, "IfStmt @%d:%d\n", node.P.Line, node.P.Column)
		p.printNode(node.Cond, level+1, "cond")
		p.indent(level + 1)
		p.buf.WriteString("then:\n")
		for _, s := range node.Then {
			p.printNode(s, level+2, "")
		}
		if len(node.Else) > 0 {
			p.indent(level + 1)
			p.buf.WriteString("else:\n")
			for _, s := range node.Else {
				p.printNode(s, level+2, "")
			}
		}
	case *ForStmt:
		fmt.Fprintf(p.buf, "ForStmt var=%s @%d:%d\n", node.Var, node.P.Line, node.P.Column)
		p.printNode(node.Start, level+1, "start")
		p.printNode(node.End, level+1, "end")
		for _, s := range node.Body {
			p.printNode(s, level+1, "body")
		}
	case *WhileStmt:
		fmt.Fprintf(p.buf, "WhileStmt @%d:%d\n", node.P.Line, node.P.Column)
		p.printNode(node.Cond, level+1, "cond")
		for _, s := range node.Body {
			p.printNode(s, level+1, "body")
		}
	case *ReturnStmt:
		fmt.Fprintf(p.buf, "ReturnStmt @%d:%d\n", node.P.Line, node.P.Column)
		p.printNode(node.Value, level+1, "value")
	case *BreakStmt:
		fmt.Fprintf(p.buf, "BreakStmt @%d:%d\n", node.P.Line, node.P.Column)
	case *ContinueStmt:
		fmt.Fprintf(p.buf, "ContinueStmt @%d:%d\n", node.P.Line, node.P.Column)
	case *ExistsCond:
		fmt.Fprintf(p.buf, "ExistsCond @%d:%d\n", node.P.Line, node.P.Column)
		p.printNode(node.Path, level+1, "path")
	case *IdentExpr:
		fmt.Fprintf(p.buf, "IdentExpr %s @%d:%d\n", node.Name, node.P.Line, node.P.Column)
	case *StringLit:
		fmt.Fprintf(p.buf, "StringLit %q @%d:%d\n", node.Value, node.P.Line, node.P.Column)
	case *NumberLit:
		fmt.Fprintf(p.buf, "NumberLit %s @%d:%d\n", node.Value, node.P.Line, node.P.Column)
	case *BoolLit:
		fmt.Fprintf(p.buf, "BoolLit %t @%d:%d\n", node.Value, node.P.Line, node.P.Column)
	case *ListLit:
		fmt.Fprintf(p.buf, "ListLit @%d:%d\n", node.P.Line, node.P.Column)
		for i, el := range node.Elements {
			p.printNode(el, level+1, fmt.Sprintf("elem[%d]", i))
		}
	case *MapLit:
		fmt.Fprintf(p.buf, "MapLit @%d:%d\n", node.P.Line, node.P.Column)
		for i := range node.Pairs {
			pair := node.Pairs[i]
			p.indent(level + 1)
			fmt.Fprintf(p.buf, "pair[%d] key=%s @%d:%d\n", i, pair.Key, pair.P.Line, pair.P.Column)
			p.printNode(pair.Value, level+2, "value")
		}
	case *IndexExpr:
		fmt.Fprintf(p.buf, "IndexExpr @%d:%d\n", node.P.Line, node.P.Column)
		p.printNode(node.Left, level+1, "left")
		p.printNode(node.Index, level+1, "index")
	case *PropertyExpr:
		fmt.Fprintf(p.buf, "PropertyExpr field=%s @%d:%d\n", node.Field, node.P.Line, node.P.Column)
		p.printNode(node.Object, level+1, "object")
	case *BinaryExpr:
		fmt.Fprintf(p.buf, "BinaryExpr op=%s @%d:%d\n", node.Op, node.P.Line, node.P.Column)
		p.printNode(node.Left, level+1, "left")
		p.printNode(node.Right, level+1, "right")
	case *UnaryExpr:
		fmt.Fprintf(p.buf, "UnaryExpr op=%s @%d:%d\n", node.Op, node.P.Line, node.P.Column)
		p.printNode(node.Right, level+1, "right")
	default:
		fmt.Fprintf(p.buf, "%T @%d:%d\n", n, n.Pos().Line, n.Pos().Column)
	}
}
