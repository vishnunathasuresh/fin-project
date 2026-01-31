package diagnostics

import (
	"strings"
	"testing"

	"github.com/vishnunathasuresh/fin-project/internal/ast"
)

func TestReporter(t *testing.T) {
	source := `a := 2
b := 3
c := a + b
`
	r := NewReporter("test.fin", source)

	if r.HasErrors() {
		t.Error("new reporter should have no errors")
	}

	r.Error(ast.Pos{Line: 1, Column: 1}, ErrSyntax, "test error")

	if !r.HasErrors() {
		t.Error("reporter should have errors after Error()")
	}

	if r.ErrorCount != 1 {
		t.Errorf("expected 1 error, got %d", r.ErrorCount)
	}
}

func TestFormatDiagnostic(t *testing.T) {
	source := `a := 2
b := 3
`
	r := NewReporter("test.fin", source)
	r.Error(ast.Pos{Line: 1, Column: 1}, ErrSyntax, "syntax error here")

	output := r.Format()

	if !strings.Contains(output, "test.fin:1:1") {
		t.Error("output should contain filename and position")
	}

	if !strings.Contains(output, "error[E001]") {
		t.Error("output should contain error severity and code")
	}

	if !strings.Contains(output, "syntax error here") {
		t.Error("output should contain message")
	}

	if !strings.Contains(output, "a := 2") {
		t.Error("output should contain source line")
	}

	if !strings.Contains(output, "^") {
		t.Error("output should contain caret indicator")
	}
}

func TestMultipleDiagnostics(t *testing.T) {
	source := `a := 2
b := 3
c := a + b
`
	r := NewReporter("test.fin", source)
	r.Error(ast.Pos{Line: 1, Column: 1}, ErrSyntax, "first error")
	r.Warning(ast.Pos{Line: 2, Column: 1}, WarnUnusedVar, "unused variable")
	r.Error(ast.Pos{Line: 3, Column: 5}, ErrTypeMismatch, "type error")

	if r.ErrorCount != 2 {
		t.Errorf("expected 2 errors, got %d", r.ErrorCount)
	}

	if r.WarnCount != 1 {
		t.Errorf("expected 1 warning, got %d", r.WarnCount)
	}

	diagnostics := r.Diagnostics()
	if len(diagnostics) != 3 {
		t.Errorf("expected 3 diagnostics, got %d", len(diagnostics))
	}
}

func TestFormatterFunctions(t *testing.T) {
	pos := ast.Pos{Line: 1, Column: 5}

	tests := []struct {
		name     string
		diag     Diagnostic
		wantCode string
		wantMsg  string
	}{
		{
			name:     "syntax error",
			diag:     SyntaxError(pos, "invalid syntax"),
			wantCode: ErrSyntax,
			wantMsg:  "invalid syntax",
		},
		{
			name:     "unexpected token",
			diag:     UnexpectedTokenError(pos, "IDENT", "NUMBER"),
			wantCode: ErrUnexpectedToken,
			wantMsg:  "expected IDENT, got NUMBER",
		},
		{
			name:     "undeclared var",
			diag:     UndeclaredVarError(pos, "x"),
			wantCode: ErrUndeclaredVar,
			wantMsg:  "undeclared variable: x",
		},
		{
			name:     "redeclared",
			diag:     RedeclaredError(pos, "y"),
			wantCode: ErrRedeclared,
			wantMsg:  "variable already declared: y",
		},
		{
			name:     "type mismatch",
			diag:     TypeMismatchError(pos, "int", "str"),
			wantCode: ErrTypeMismatch,
			wantMsg:  "type mismatch: expected int, got str",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.diag.Code != tt.wantCode {
				t.Errorf("expected code %s, got %s", tt.wantCode, tt.diag.Code)
			}
			if tt.diag.Message != tt.wantMsg {
				t.Errorf("expected message '%s', got '%s'", tt.wantMsg, tt.diag.Message)
			}
		})
	}
}
