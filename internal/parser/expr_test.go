package parser

import (
	"testing"

	"github.com/vishnunath-suresh/fin-project/internal/ast"
	"github.com/vishnunath-suresh/fin-project/internal/lexer"
)

func parseExpr(t *testing.T, src string) ast.Expr {
	t.Helper()
	l := lexer.New(src)
	toks := CollectTokens(l)
	p := New(toks)
	expr := p.parseExpression(0)
	return expr
}

func requireBinary(t *testing.T, expr ast.Expr) *ast.BinaryExpr {
	t.Helper()
	bin, ok := expr.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expr not BinaryExpr: %T", expr)
	}
	return bin
}

func TestParseExpression_ArithmeticPrecedence(t *testing.T) {
	expr := parseExpr(t, "1 + 2 * 3")

	bin := requireBinary(t, expr)
	if bin.Op != "+" {
		t.Fatalf("root op = %q, want +", bin.Op)
	}
	leftNum := bin.Left.(*ast.NumberLit)
	if leftNum.Value != "1" {
		t.Fatalf("left value = %s, want 1", leftNum.Value)
	}
	right, ok := bin.Right.(*ast.BinaryExpr)
	if !ok || right.Op != "*" {
		t.Fatalf("right not BinaryExpr *: %T %q", bin.Right, right.Op)
	}
	if right.Left.(*ast.NumberLit).Value != "2" || right.Right.(*ast.NumberLit).Value != "3" {
		t.Fatalf("multiplication operands incorrect")
	}
}

func TestParseExpression_BooleanPrecedence(t *testing.T) {
	expr := parseExpr(t, "true || false && true")
	orNode := requireBinary(t, expr)
	if orNode.Op != "||" {
		t.Fatalf("root op = %q, want ||", orNode.Op)
	}
	leftBool := orNode.Left.(*ast.BoolLit)
	if !leftBool.Value {
		t.Fatalf("left bool should be true")
	}
	andNode, ok := orNode.Right.(*ast.BinaryExpr)
	if !ok || andNode.Op != "&&" {
		t.Fatalf("right not &&: %T %q", orNode.Right, andNode.Op)
	}
}

func TestParseExpression_MixedWithGrouping(t *testing.T) {
	expr := parseExpr(t, "(1 + 2) * 3")
	mul := requireBinary(t, expr)
	if mul.Op != "*" {
		t.Fatalf("root not *: %T %q", expr, mul.Op)
	}
	add, ok := mul.Left.(*ast.BinaryExpr)
	if !ok || add.Op != "+" {
		t.Fatalf("left not +: %T %q", mul.Left, add.Op)
	}
	if add.Left.(*ast.NumberLit).Value != "1" || add.Right.(*ast.NumberLit).Value != "2" {
		t.Fatalf("addition operands incorrect")
	}
}

func TestParseExpression_IndexAndProperty(t *testing.T) {
	expr := parseExpr(t, "arr[0].name")
	prop, ok := expr.(*ast.PropertyExpr)
	if !ok {
		t.Fatalf("root not PropertyExpr: %T", expr)
	}
	if prop.Field != "name" {
		t.Fatalf("field = %q, want name", prop.Field)
	}
	idx, ok := prop.Object.(*ast.IndexExpr)
	if !ok {
		t.Fatalf("object not IndexExpr: %T", prop.Object)
	}
	if idx.Left.(*ast.IdentExpr).Name != "arr" {
		t.Fatalf("index left ident wrong")
	}
	if idx.Index.(*ast.NumberLit).Value != "0" {
		t.Fatalf("index value wrong")
	}
}

func TestParseExpression_UnaryBinding(t *testing.T) {
	expr := parseExpr(t, "-1 * 2")
	mul := requireBinary(t, expr)
	if mul.Op != "*" {
		t.Fatalf("root op = %q, want *", mul.Op)
	}
	unary, ok := mul.Left.(*ast.UnaryExpr)
	if !ok || unary.Op != "-" {
		t.Fatalf("left not unary -: %T %q", mul.Left, unary.Op)
	}
}

func TestParseExpression_ExponentPrecedence(t *testing.T) {
	expr := parseExpr(t, "2 ** 3 * 4")
	root := requireBinary(t, expr)
	if root.Op != "*" {
		t.Fatalf("root op = %q, want *", root.Op)
	}
	pow := requireBinary(t, root.Left)
	if pow.Op != "**" {
		t.Fatalf("left op = %q, want **", pow.Op)
	}
}

func TestParseExpression_ExponentRightAssociative(t *testing.T) {
	expr := parseExpr(t, "2 ** 3 ** 2")
	root := requireBinary(t, expr)
	if root.Op != "**" {
		t.Fatalf("root op = %q, want **", root.Op)
	}
	// Right-associative: left is 2, right is (3 ** 2)
	if _, ok := root.Left.(*ast.NumberLit); !ok {
		t.Fatalf("left not number")
	}
	rightPow := requireBinary(t, root.Right)
	if rightPow.Op != "**" {
		t.Fatalf("right op = %q, want **", rightPow.Op)
	}
}

func TestParseExpression_Malformed_NoPanic(t *testing.T) {
	l := lexer.New("!")
	toks := CollectTokens(l)
	p := New(toks)
	_ = p.parseExpression(0)
	// No panic; errors may be recorded
}
