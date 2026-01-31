package sema

import (
	"fmt"

	"github.com/vishnunathasuresh/fin-project/internal/ast"
	"github.com/vishnunathasuresh/fin-project/internal/diagnostics"
)

// DiagnosticError exposes structured information for diagnostics reporting.
type DiagnosticError interface {
	error
	Pos() ast.Pos
	DiagnosticCode() string
	DiagnosticMessage() string
}

// UndefinedVariableError is raised when a variable is referenced before declaration.
type UndefinedVariableError struct {
	Name string
	P    ast.Pos
}

func (e UndefinedVariableError) Error() string {
	return fmt.Sprintf("undefined variable %q at %d:%d — referenced before declaration", e.Name, e.P.Line, e.P.Column)
}

func (e UndefinedVariableError) Pos() ast.Pos {
	return e.P
}

func (e UndefinedVariableError) DiagnosticCode() string {
	return diagnostics.ErrUndeclaredVar
}

func (e UndefinedVariableError) DiagnosticMessage() string {
	return fmt.Sprintf("undefined variable %q", e.Name)
}

// DuplicateFunctionError is raised when a function name is declared more than once.
type DuplicateFunctionError struct {
	Name string
	P    ast.Pos
}

func (e DuplicateFunctionError) Error() string {
	return fmt.Sprintf("duplicate function %q at %d:%d — function names must be unique", e.Name, e.P.Line, e.P.Column)
}

func (e DuplicateFunctionError) Pos() ast.Pos {
	return e.P
}

func (e DuplicateFunctionError) DiagnosticCode() string {
	return diagnostics.ErrRedeclared
}

func (e DuplicateFunctionError) DiagnosticMessage() string {
	return fmt.Sprintf("duplicate function %q", e.Name)
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

func (e InvalidArityError) Pos() ast.Pos {
	return e.P
}

func (e InvalidArityError) DiagnosticCode() string {
	if e.Got < e.Expected {
		return diagnostics.ErrTooFewArgs
	}
	return diagnostics.ErrTooManyArgs
}

func (e InvalidArityError) DiagnosticMessage() string {
	return fmt.Sprintf("invalid arity for %q: expected %d args, got %d", e.Name, e.Expected, e.Got)
}

// ReservedNameError is raised when a reserved identifier is used illegally.
type ReservedNameError struct {
	Name string
	P    ast.Pos
}

func (e ReservedNameError) Error() string {
	return fmt.Sprintf("reserved name %q at %d:%d — choose a different identifier", e.Name, e.P.Line, e.P.Column)
}

func (e ReservedNameError) Pos() ast.Pos {
	return e.P
}

func (e ReservedNameError) DiagnosticCode() string {
	return diagnostics.ErrInvalidType
}

func (e ReservedNameError) DiagnosticMessage() string {
	return fmt.Sprintf("reserved name %q", e.Name)
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

func (e ShadowingError) Pos() ast.Pos {
	return e.P
}

func (e ShadowingError) DiagnosticCode() string {
	return diagnostics.ErrRedeclared
}

func (e ShadowingError) DiagnosticMessage() string {
	return fmt.Sprintf("name %q already defined in an enclosing scope", e.Name)
}

// DepthExceededError is raised when traversal exceeds the configured recursion limit.
type DepthExceededError struct {
	Limit int
	P     ast.Pos
}

func (e DepthExceededError) Error() string {
	return fmt.Sprintf("traversal depth exceeded limit %d at %d:%d", e.Limit, e.P.Line, e.P.Column)
}

func (e DepthExceededError) Pos() ast.Pos {
	return e.P
}

func (e DepthExceededError) DiagnosticCode() string {
	return diagnostics.ErrSyntax
}

func (e DepthExceededError) DiagnosticMessage() string {
	return fmt.Sprintf("traversal depth exceeded limit %d", e.Limit)
}

// ReturnOutsideFunctionError is raised when a return is used outside a function body.
type ReturnOutsideFunctionError struct {
	P ast.Pos
}

func (e ReturnOutsideFunctionError) Error() string {
	return fmt.Sprintf("return used outside function at %d:%d", e.P.Line, e.P.Column)
}

func (e ReturnOutsideFunctionError) Pos() ast.Pos {
	return e.P
}

func (e ReturnOutsideFunctionError) DiagnosticCode() string {
	return diagnostics.ErrReturnOutside
}

func (e ReturnOutsideFunctionError) DiagnosticMessage() string {
	return "return used outside function"
}
