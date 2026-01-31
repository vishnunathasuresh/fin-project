package diagnostics

import (
	"fmt"
	"strings"

	"github.com/vishnunath-suresh/fin-project/internal/ast"
)

// Severity levels for diagnostics
type Severity int

const (
	SeverityError Severity = iota
	SeverityWarning
	SeverityInfo
)

func (s Severity) String() string {
	switch s {
	case SeverityError:
		return "error"
	case SeverityWarning:
		return "warning"
	case SeverityInfo:
		return "info"
	default:
		return "unknown"
	}
}

// Diagnostic represents a compiler diagnostic (error, warning, or info)
type Diagnostic struct {
	Severity Severity
	Pos      ast.Pos
	Message  string
	Code     string // Error code like "E001", "W005"
}

// Reporter collects diagnostics
type Reporter struct {
	diagnostics []Diagnostic
	source      string
	filename    string
	ErrorCount  int
	WarnCount   int
}

// NewReporter creates a new diagnostic reporter
func NewReporter(filename, source string) *Reporter {
	return &Reporter{
		diagnostics: []Diagnostic{},
		source:      source,
		filename:    filename,
	}
}

// Error adds an error diagnostic
func (r *Reporter) Error(pos ast.Pos, code, message string) {
	r.diagnostics = append(r.diagnostics, Diagnostic{
		Severity: SeverityError,
		Pos:      pos,
		Message:  message,
		Code:     code,
	})
	r.ErrorCount++
}

// Warning adds a warning diagnostic
func (r *Reporter) Warning(pos ast.Pos, code, message string) {
	r.diagnostics = append(r.diagnostics, Diagnostic{
		Severity: SeverityWarning,
		Pos:      pos,
		Message:  message,
		Code:     code,
	})
	r.WarnCount++
}

// Info adds an info diagnostic
func (r *Reporter) Info(pos ast.Pos, code, message string) {
	r.diagnostics = append(r.diagnostics, Diagnostic{
		Severity: SeverityInfo,
		Pos:      pos,
		Message:  message,
		Code:     code,
	})
}

// HasErrors returns true if any errors were reported
func (r *Reporter) HasErrors() bool {
	return r.ErrorCount > 0
}

// HasWarnings returns true if any warnings were reported
func (r *Reporter) HasWarnings() bool {
	return r.WarnCount > 0
}

// Diagnostics returns all diagnostics
func (r *Reporter) Diagnostics() []Diagnostic {
	return r.diagnostics
}

// Format returns a formatted string of all diagnostics
func (r *Reporter) Format() string {
	var sb strings.Builder

	for _, diag := range r.diagnostics {
		sb.WriteString(r.FormatDiagnostic(&diag))
		sb.WriteString("\n")
	}

	return sb.String()
}

// FormatDiagnostic formats a single diagnostic
func (r *Reporter) FormatDiagnostic(diag *Diagnostic) string {
	var sb strings.Builder

	// Format: filename:line:col: severity[code]: message
	sb.WriteString(fmt.Sprintf("%s:%d:%d: %s[%s]: %s\n",
		r.filename,
		diag.Pos.Line,
		diag.Pos.Column,
		diag.Severity.String(),
		diag.Code,
		diag.Message,
	))

	// Add source line with caret
	if r.source != "" {
		sourceLine := r.getSourceLine(diag.Pos.Line)
		if sourceLine != "" {
			sb.WriteString(sourceLine)
			sb.WriteString("\n")
			sb.WriteString(strings.Repeat(" ", diag.Pos.Column-1))
			sb.WriteString("^\n")
		}
	}

	return sb.String()
}

// getSourceLine extracts a specific line from source
func (r *Reporter) getSourceLine(lineNum int) string {
	lines := strings.Split(r.source, "\n")
	if lineNum < 1 || lineNum > len(lines) {
		return ""
	}
	return lines[lineNum-1]
}

// Clear resets all diagnostics
func (r *Reporter) Clear() {
	r.diagnostics = []Diagnostic{}
	r.ErrorCount = 0
	r.WarnCount = 0
}
