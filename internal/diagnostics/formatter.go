package diagnostics

import (
	"fmt"
	"strings"

	"github.com/vishnunath-suresh/fin-project/internal/ast"
)

// Common error codes and formatters

// Error codes
const (
	ErrSyntax          = "E001"
	ErrUnexpectedToken = "E002"
	ErrUndeclaredVar   = "E003"
	ErrRedeclared      = "E004"
	ErrTypeMismatch    = "E005"
	ErrInvalidType     = "E006"
	ErrTooFewArgs      = "E007"
	ErrTooManyArgs     = "E008"
	ErrBreakOutside    = "E009"
	ErrContinueOutside = "E010"
	ErrReturnOutside   = "E011"
	ErrDivByZero       = "E012"
)

// Warning codes
const (
	WarnUnusedVar = "W001"
	WarnUnusedFn  = "W002"
	WarnShadowing = "W003"
)

// SyntaxError creates a syntax error diagnostic
func SyntaxError(pos ast.Pos, message string) Diagnostic {
	return Diagnostic{
		Severity: SeverityError,
		Pos:      pos,
		Code:     ErrSyntax,
		Message:  message,
	}
}

// UnexpectedTokenError creates an unexpected token error
func UnexpectedTokenError(pos ast.Pos, expected, got string) Diagnostic {
	return Diagnostic{
		Severity: SeverityError,
		Pos:      pos,
		Code:     ErrUnexpectedToken,
		Message:  fmt.Sprintf("expected %s, got %s", expected, got),
	}
}

// UndeclaredVarError creates an undeclared variable error
func UndeclaredVarError(pos ast.Pos, name string) Diagnostic {
	return Diagnostic{
		Severity: SeverityError,
		Pos:      pos,
		Code:     ErrUndeclaredVar,
		Message:  fmt.Sprintf("undeclared variable: %s", name),
	}
}

// RedeclaredError creates a redeclaration error
func RedeclaredError(pos ast.Pos, name string) Diagnostic {
	return Diagnostic{
		Severity: SeverityError,
		Pos:      pos,
		Code:     ErrRedeclared,
		Message:  fmt.Sprintf("variable already declared: %s", name),
	}
}

// TypeMismatchError creates a type mismatch error
func TypeMismatchError(pos ast.Pos, expected, got string) Diagnostic {
	return Diagnostic{
		Severity: SeverityError,
		Pos:      pos,
		Code:     ErrTypeMismatch,
		Message:  fmt.Sprintf("type mismatch: expected %s, got %s", expected, got),
	}
}

// InvalidTypeError creates an invalid type error
func InvalidTypeError(pos ast.Pos, typeName string) Diagnostic {
	return Diagnostic{
		Severity: SeverityError,
		Pos:      pos,
		Code:     ErrInvalidType,
		Message:  fmt.Sprintf("invalid type: %s", typeName),
	}
}

// TooFewArgsError creates a too few arguments error
func TooFewArgsError(pos ast.Pos, fn string, expected, got int) Diagnostic {
	return Diagnostic{
		Severity: SeverityError,
		Pos:      pos,
		Code:     ErrTooFewArgs,
		Message:  fmt.Sprintf("%s expects %d arguments, got %d", fn, expected, got),
	}
}

// TooManyArgsError creates a too many arguments error
func TooManyArgsError(pos ast.Pos, fn string, expected, got int) Diagnostic {
	return Diagnostic{
		Severity: SeverityError,
		Pos:      pos,
		Code:     ErrTooManyArgs,
		Message:  fmt.Sprintf("%s expects %d arguments, got %d", fn, expected, got),
	}
}

// BreakOutsideLoopError creates a break outside loop error
func BreakOutsideLoopError(pos ast.Pos) Diagnostic {
	return Diagnostic{
		Severity: SeverityError,
		Pos:      pos,
		Code:     ErrBreakOutside,
		Message:  "break statement outside loop",
	}
}

// ContinueOutsideLoopError creates a continue outside loop error
func ContinueOutsideLoopError(pos ast.Pos) Diagnostic {
	return Diagnostic{
		Severity: SeverityError,
		Pos:      pos,
		Code:     ErrContinueOutside,
		Message:  "continue statement outside loop",
	}
}

// ReturnOutsideFnError creates a return outside function error
func ReturnOutsideFnError(pos ast.Pos) Diagnostic {
	return Diagnostic{
		Severity: SeverityError,
		Pos:      pos,
		Code:     ErrReturnOutside,
		Message:  "return statement outside function",
	}
}

// UnusedVarWarning creates an unused variable warning
func UnusedVarWarning(pos ast.Pos, name string) Diagnostic {
	return Diagnostic{
		Severity: SeverityWarning,
		Pos:      pos,
		Code:     WarnUnusedVar,
		Message:  fmt.Sprintf("unused variable: %s", name),
	}
}

// ShadowingWarning creates a variable shadowing warning
func ShadowingWarning(pos ast.Pos, name string) Diagnostic {
	return Diagnostic{
		Severity: SeverityWarning,
		Pos:      pos,
		Code:     WarnShadowing,
		Message:  fmt.Sprintf("variable %s shadows declaration in outer scope", name),
	}
}

// FormatMultiple formats multiple diagnostics as a single string
func FormatMultiple(filename, source string, diagnostics []Diagnostic) string {
	if len(diagnostics) == 0 {
		return ""
	}

	r := NewReporter(filename, source)
	for _, diag := range diagnostics {
		r.diagnostics = append(r.diagnostics, diag)
		switch diag.Severity {
		case SeverityError:
			r.ErrorCount++
		case SeverityWarning:
			r.WarnCount++
		}
	}

	var sb strings.Builder
	sb.WriteString(r.Format())

	// Summary line
	if r.ErrorCount > 0 || r.WarnCount > 0 {
		sb.WriteString(fmt.Sprintf("Found %d error(s) and %d warning(s)\n", r.ErrorCount, r.WarnCount))
	}

	return sb.String()
}
