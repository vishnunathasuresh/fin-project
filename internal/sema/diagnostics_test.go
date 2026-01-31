package sema

import (
	"testing"

	"github.com/vishnunathasuresh/fin-project/internal/diagnostics"
	"github.com/vishnunathasuresh/fin-project/internal/lexer"
	"github.com/vishnunathasuresh/fin-project/internal/parser"
)

func TestAnalyzeReportsDiagnostics(t *testing.T) {
	src := "echo a\n"
	l := lexer.New(src)
	toks := parser.CollectTokens(l)
	p := parser.New(toks)
	prog := p.ParseProgram()

	reporter := diagnostics.NewReporter("test.fin", src)
	res := AnalyzeDefinitionsWithReporter(prog, reporter, 0)
	if len(res.Errors) == 0 {
		t.Fatalf("expected semantic errors")
	}
	if !reporter.HasErrors() {
		t.Fatalf("expected reporter errors")
	}
	diags := reporter.Diagnostics()
	if len(diags) == 0 {
		t.Fatalf("expected diagnostics entries")
	}
	if diags[0].Code != diagnostics.ErrUndeclaredVar {
		t.Fatalf("expected code %s, got %s", diagnostics.ErrUndeclaredVar, diags[0].Code)
	}
}
