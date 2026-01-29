package sema

import (
    "testing"

    "github.com/vishnunath-suresh/fin-project/internal/ast"
)

func TestValidateIdentifier_RejectsKeywords(t *testing.T) {
    pos := ast.Pos{Line: 1, Column: 2}
    if err := ValidateIdentifier("if", pos); err == nil {
        t.Fatalf("expected reserved error for keyword")
    }
}

func TestValidateIdentifier_AllowsNonReserved(t *testing.T) {
    pos := ast.Pos{Line: 1, Column: 2}
    if err := ValidateIdentifier("myVar", pos); err != nil {
        t.Fatalf("unexpected error for non-reserved: %v", err)
    }
}
