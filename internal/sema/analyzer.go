package sema

import "github.com/vishnunath-suresh/fin-project/internal/ast"

// FunctionRegistry tracks function signatures by name.
type FunctionRegistry struct {
	funcs map[string]int
}

// NewFunctionRegistry creates an empty registry.
func NewFunctionRegistry() *FunctionRegistry {
	return &FunctionRegistry{funcs: make(map[string]int)}
}

// Define registers a function name and its parameter count.
// It returns an error if the name already exists. The provided pos is used for diagnostics.
func (r *FunctionRegistry) Define(name string, arity int, pos ast.Pos) error {
	if _, exists := r.funcs[name]; exists {
		return DuplicateFunctionError{Name: name, P: pos}
	}
	r.funcs[name] = arity
	return nil
}

// Lookup returns the arity for a function and whether it was found.
func (r *FunctionRegistry) Lookup(name string) (int, bool) {
	arity, ok := r.funcs[name]
	return arity, ok
}
