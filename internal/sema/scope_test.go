package sema

import "testing"

func TestScope_NestedAndShadow(t *testing.T) {
	root := NewScope(nil)
	if err := root.Define("a"); err != nil {
		t.Fatalf("expected to define a in root: %v", err)
	}
	child := NewScope(root)
	if _, ok := child.Lookup("a"); !ok {
		t.Fatalf("expected to find a in parent")
	}
	// shadow
	if err := child.Define("a"); err != nil {
		t.Fatalf("expected to define shadow a in child: %v", err)
	}
	foundScope, ok := child.Lookup("a")
	if !ok || foundScope != child {
		t.Fatalf("expected shadowed a in child scope")
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
	if err := s.Define("x"); err != nil {
		t.Fatalf("unexpected error defining x: %v", err)
	}
	if err := s.Define("x"); err == nil {
		t.Fatalf("expected error on duplicate definition in same scope")
	}
}
