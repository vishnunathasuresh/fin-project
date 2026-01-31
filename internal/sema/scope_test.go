package sema

import (
	"testing"

	"github.com/vishnunathasuresh/fin-project/internal/ast"
)

func TestScope_NestedShadowIsError(t *testing.T) {
	root := NewScope(nil)
	if err := root.Define("a", ast.Pos{Line: 1, Column: 1}); err != nil {
		t.Fatalf("expected to define a in root: %v", err)
	}
	child := NewScope(root)
	if _, ok := child.Lookup("a"); !ok {
		t.Fatalf("expected to find a in parent")
	}
	if err := child.Define("a", ast.Pos{Line: 2, Column: 1}); err == nil {
		t.Fatalf("expected error when shadowing a from parent")
	}
}

func TestScope_UndefinedLookup(t *testing.T) {
	s := NewScope(nil)
	if _, ok := s.Lookup("missing"); ok {
		t.Fatalf("expected missing to be undefined")
	}
}

func TestScope_DuplicateInSameScope(t *testing.T) {
	s := NewScope(nil)
	if err := s.Define("x", ast.Pos{Line: 1, Column: 1}); err != nil {
		t.Fatalf("unexpected error defining x: %v", err)
	}
	if err := s.Define("x", ast.Pos{Line: 1, Column: 1}); err == nil {
		t.Fatalf("expected error on duplicate definition in same scope")
	}
}
