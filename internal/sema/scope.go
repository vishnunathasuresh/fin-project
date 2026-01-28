package sema

import "fmt"

// Scope represents a lexical scope with an optional parent and a table of names.
type Scope struct {
	Parent *Scope
	vars   map[string]struct{}
}

// NewScope creates a new scope with the given parent.
func NewScope(parent *Scope) *Scope {
	return &Scope{Parent: parent, vars: make(map[string]struct{})}
}

// Define adds a name to the current scope. Shadowing across scopes is allowed, but
// redeclaring in the same scope returns an error.
func (s *Scope) Define(name string) error {
	if _, exists := s.vars[name]; exists {
		return fmt.Errorf("name %q already defined in this scope", name)
	}
	s.vars[name] = struct{}{}
	return nil
}

// Lookup searches for name starting at this scope and walking parents.
// Returns true if found and the scope where it was found.
func (s *Scope) Lookup(name string) (*Scope, bool) {
	for sc := s; sc != nil; sc = sc.Parent {
		if _, ok := sc.vars[name]; ok {
			return sc, true
		}
	}
	return nil, false
}
