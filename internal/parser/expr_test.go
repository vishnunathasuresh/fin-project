package parser

import (
	"testing"

	"github.com/vishnunathasuresh/fin-project/internal/ast"
	"github.com/vishnunathasuresh/fin-project/internal/lexer"
)

func parseExpr(t *testing.T, src string) ast.Expr {
	t.Helper()
	l := lexer.New(src)
	toks := CollectTokens(l)
	p := New(toks)
	expr := p.parseExpression(0)
	return expr
}

func parseExprWithParser(t *testing.T, src string) (ast.Expr, *Parser) {
	t.Helper()
	l := lexer.New(src)
	toks := CollectTokens(l)
	p := New(toks)
	return p.parseExpression(0), p
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

func TestParseExpression_EmptyListAndMap(t *testing.T) {
	listExpr := parseExpr(t, "[]")
	if l, ok := listExpr.(*ast.ListLit); !ok || len(l.Elements) != 0 {
		t.Fatalf("expected empty list, got %T with len=%d", listExpr, len(l.Elements))
	}

	mapExpr := parseExpr(t, "{}")
	if m, ok := mapExpr.(*ast.MapLit); !ok || len(m.Pairs) != 0 {
		t.Fatalf("expected empty map, got %T with len=%d", mapExpr, len(m.Pairs))
	}
}

func TestParseExpression_UnmatchedDelimiters_NoPanic(t *testing.T) {
	_, p1 := parseExprWithParser(t, "[1, 2")
	if len(p1.Errors()) == 0 {
		t.Fatalf("expected errors for unmatched list delimiter")
	}
	_, p2 := parseExprWithParser(t, "{a:1")
	if len(p2.Errors()) == 0 {
		t.Fatalf("expected errors for unmatched map delimiter")
	}
}

func TestParseExpression_DeepChaining(t *testing.T) {
	expr := parseExpr(t, "a.b[0].c")
	propC, ok := expr.(*ast.PropertyExpr)
	if !ok || propC.Field != "c" {
		t.Fatalf("root not property .c: %T", expr)
	}
	idx, ok := propC.Object.(*ast.IndexExpr)
	if !ok {
		t.Fatalf("object not IndexExpr: %T", propC.Object)
	}
	if _, ok := idx.Index.(*ast.NumberLit); !ok {
		t.Fatalf("index not number")
	}
	propB, ok := idx.Left.(*ast.PropertyExpr)
	if !ok || propB.Field != "b" {
		t.Fatalf("left not property .b: %T", idx.Left)
	}
	if ident, ok := propB.Object.(*ast.IdentExpr); !ok || ident.Name != "a" {
		t.Fatalf("base not ident a: %T", propB.Object)
	}
}

func TestParseExpression_InvalidAccessRecovery(t *testing.T) {
	_, p := parseExprWithParser(t, "a.1")
	if len(p.Errors()) == 0 {
		t.Fatalf("expected errors for invalid property access")
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

func TestParseExpression_UnaryBangChaining(t *testing.T) {
	expr := parseExpr(t, "!!true")
	un1, ok := expr.(*ast.UnaryExpr)
	if !ok || un1.Op != "!" {
		t.Fatalf("outer not unary !: %T", expr)
	}
	un2, ok := un1.Right.(*ast.UnaryExpr)
	if !ok || un2.Op != "!" {
		t.Fatalf("inner not unary !: %T", un1.Right)
	}
	if b, ok := un2.Right.(*ast.BoolLit); !ok || b.Value != true {
		t.Fatalf("expected bool true inside unary chain")
	}
}

func TestParseExpression_UnaryPrecedenceWithMul(t *testing.T) {
	expr := parseExpr(t, "-5 * 3")
	mul := requireBinary(t, expr)
	if mul.Op != "*" {
		t.Fatalf("root not *: %q", mul.Op)
	}
	if un, ok := mul.Left.(*ast.UnaryExpr); !ok || un.Op != "-" {
		t.Fatalf("left not unary -: %T", mul.Left)
	}
	if num, ok := mul.Right.(*ast.NumberLit); !ok || num.Value != "3" {
		t.Fatalf("right not number 3: %T", mul.Right)
	}
}

func TestParseExpression_UnaryBinding(t *testing.T) {
	expr := parseExpr(t, "-a * b")
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

func TestParseExpression_SingleLiteral(t *testing.T) {
	expr := parseExpr(t, "\"hi\"")
	str, ok := expr.(*ast.StringLit)
	if !ok || str.Value != "hi" {
		t.Fatalf("expected string literal 'hi', got %T %v", expr, str)
	}
}

func TestParseExpression_NestedList(t *testing.T) {
	expr := parseExpr(t, "[1, [2, 3]]")
	list, ok := expr.(*ast.ListLit)
	if !ok {
		t.Fatalf("root not ListLit: %T", expr)
	}
	if len(list.Elements) != 2 {
		t.Fatalf("outer list len = %d, want 2", len(list.Elements))
	}
	inner, ok := list.Elements[1].(*ast.ListLit)
	if !ok || len(inner.Elements) != 2 {
		t.Fatalf("inner list invalid: %T len=%d", list.Elements[1], len(inner.Elements))
	}
}

func TestParseExpression_MapMultipleKeys(t *testing.T) {
	expr := parseExpr(t, "{a:1, b:2}")
	mapLit, ok := expr.(*ast.MapLit)
	if !ok {
		t.Fatalf("root not MapLit: %T", expr)
	}
	if len(mapLit.Pairs) != 2 {
		t.Fatalf("map pairs len = %d, want 2", len(mapLit.Pairs))
	}
	if mapLit.Pairs[0].Key != "a" || mapLit.Pairs[1].Key != "b" {
		t.Fatalf("map keys wrong: %q %q", mapLit.Pairs[0].Key, mapLit.Pairs[1].Key)
	}
}
