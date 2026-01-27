package ast

import (
	"strings"
	"testing"
)

func TestFormat_ProgramWithStatements(t *testing.T) {
	prog := &Program{
		Statements: []Statement{
			&SetStmt{Name: "x", Value: &NumberLit{Value: "1", P: Pos{Line: 1, Column: 7}}, P: Pos{Line: 1, Column: 1}},
			&EchoStmt{Value: &IdentExpr{Name: "x", P: Pos{Line: 2, Column: 6}}, P: Pos{Line: 2, Column: 1}},
		},
		P: Pos{Line: 1, Column: 1},
	}

	out := Format(prog)
	if !containsAll(out, []string{"Program @1:1", "SetStmt name=x @1:1", "NumberLit 1 @1:7", "EchoStmt @2:1", "IdentExpr x @2:6"}) {
		t.Fatalf("format output missing expected substrings:\n%s", out)
	}
}

func containsAll(s string, subs []string) bool {
	for _, sub := range subs {
		if !strings.Contains(s, sub) {
			return false
		}
	}
	return true
}
