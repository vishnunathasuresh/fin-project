package sema

import (
	"github.com/vishnunathasuresh/fin-project/internal/ast"
	"github.com/vishnunathasuresh/fin-project/internal/diagnostics"
)

// AnalyzeDefinitionsWithReporter runs analysis and reports diagnostics to the reporter.
func AnalyzeDefinitionsWithReporter(prog *ast.Program, reporter *diagnostics.Reporter, limit int) AnalysisResult {
	res := AnalyzeDefinitionsWithLimit(prog, limit)
	if reporter != nil {
		ReportDiagnostics(reporter, res.Errors)
	}
	return res
}

// ReportDiagnostics maps semantic errors into diagnostics.
func ReportDiagnostics(reporter *diagnostics.Reporter, errs []error) {
	for _, err := range errs {
		if err == nil {
			continue
		}
		if diagErr, ok := err.(DiagnosticError); ok {
			reporter.Error(diagErr.Pos(), diagErr.DiagnosticCode(), diagErr.DiagnosticMessage())
			continue
		}
		if posErr, ok := err.(interface{ Pos() ast.Pos }); ok {
			reporter.Error(posErr.Pos(), diagnostics.ErrSyntax, err.Error())
			continue
		}
		reporter.Error(ast.Pos{Line: 1, Column: 1}, diagnostics.ErrSyntax, err.Error())
	}
}
