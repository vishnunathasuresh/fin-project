package generator

import (
	"fmt"

	"github.com/vishnunath-suresh/fin-project/internal/ast"
)

// GeneratorError is a typed error for generator failures.
type GeneratorError struct {
	Msg string
	Pos ast.Pos
}

func (e *GeneratorError) Error() string {
	if e.Pos.Line > 0 {
		return fmt.Sprintf("generator error at %d:%d: %s", e.Pos.Line, e.Pos.Column, e.Msg)
	}
	return fmt.Sprintf("generator error: %s", e.Msg)
}

func errUnsupportedStmt(pos ast.Pos, stmt ast.Statement) error {
	return &GeneratorError{Msg: fmt.Sprintf("unsupported statement type %T", stmt), Pos: pos}
}

func errFunctionNotLifted(pos ast.Pos, name string) error {
	return &GeneratorError{Msg: fmt.Sprintf("function declaration '%s' should be lifted before lowering", name), Pos: pos}
}
