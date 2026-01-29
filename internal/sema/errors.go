package sema

import (
	"fmt"

	"github.com/vishnunath-suresh/fin-project/internal/ast"
)

// UndefinedVariableError is raised when a variable is referenced before declaration.
type UndefinedVariableError struct {
	Name string
	P    ast.Pos
}

func (e UndefinedVariableError) Error() string {
	return fmt.Sprintf("undefined variable %q at %d:%d — referenced before declaration", e.Name, e.P.Line, e.P.Column)
}

// DuplicateFunctionError is raised when a function name is declared more than once.
type DuplicateFunctionError struct {
	Name string
	P    ast.Pos
}

func (e DuplicateFunctionError) Error() string {
	return fmt.Sprintf("duplicate function %q at %d:%d — function names must be unique", e.Name, e.P.Line, e.P.Column)
}

// InvalidArityError is raised when a function is called with an unexpected number of arguments.
type InvalidArityError struct {
	Name     string
	Expected int
	Got      int
	P        ast.Pos
}

func (e InvalidArityError) Error() string {
	return fmt.Sprintf("invalid arity for %q at %d:%d — expected %d args, got %d", e.Name, e.P.Line, e.P.Column, e.Expected, e.Got)
}

// ReservedNameError is raised when a reserved identifier is used illegally.
type ReservedNameError struct {
	Name string
	P    ast.Pos
}

func (e ReservedNameError) Error() string {
	return fmt.Sprintf("reserved name %q at %d:%d — choose a different identifier", e.Name, e.P.Line, e.P.Column)
}

// ShadowingError is raised when a name is redefined in an enclosing scope.
// P is the shadowing position; Def is the original definition position.
type ShadowingError struct {
	Name string
	P    ast.Pos
	Def  ast.Pos
}

func (e ShadowingError) Error() string {
	return fmt.Sprintf("name %q already defined in an enclosing scope at %d:%d (original at %d:%d)", e.Name, e.P.Line, e.P.Column, e.Def.Line, e.Def.Column)
}
