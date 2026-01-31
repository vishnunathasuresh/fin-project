package parser

import (
	"testing"

	"github.com/vishnunathasuresh/fin-project/internal/diagnostics"
	"github.com/vishnunathasuresh/fin-project/internal/lexer"
)

func TestParserReportsDiagnostics(t *testing.T) {
	src := "!\n"
	l := lexer.New(src)
	toks := CollectTokens(l)
	reporter := diagnostics.NewReporter("test.fin", src)
	p := NewWithReporter(toks, reporter)
	_ = p.ParseProgram()

	if len(p.Errors()) == 0 {
		t.Fatalf("expected parser errors")
	}
	if !reporter.HasErrors() {
		t.Fatalf("expected diagnostics errors")
	}
	diags := reporter.Diagnostics()
	if len(diags) == 0 {
		t.Fatalf("expected at least one diagnostic")
	}
	if diags[0].Code != diagnostics.ErrUnexpectedToken {
		t.Fatalf("expected code %s, got %s", diagnostics.ErrUnexpectedToken, diags[0].Code)
	}
}
