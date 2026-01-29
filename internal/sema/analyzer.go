package sema

import (
	"errors"

	"github.com/vishnunath-suresh/fin-project/internal/ast"
)

// FunctionRegistry tracks function signatures by name.
type FunctionRegistry struct {
	funcs map[string]int
}

// AnalysisResult captures scopes and errors from semantic analysis.
type AnalysisResult struct {
	Global      *Scope
	FuncScopes  map[*ast.FnDecl]*Scope
	ForScopes   map[*ast.ForStmt]*Scope
	WhileScopes map[*ast.WhileStmt]*Scope
	Errors      []error
}

// Analyzer aggregates semantic analysis results safely.
type Analyzer struct {
	prog   *ast.Program
	limit  int
	result AnalysisResult
}

// New constructs an Analyzer with no depth limit.
func New() *Analyzer {
	return &Analyzer{}
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

// AnalyzeDefinitionsWithLimit walks the AST to enforce semantic rules with an optional
// recursion depth limit. If limit <= 0, no depth check is applied.
func AnalyzeDefinitionsWithLimit(prog *ast.Program, limit int) AnalysisResult {
	res := AnalysisResult{
		Global:      NewScope(nil),
		FuncScopes:  make(map[*ast.FnDecl]*Scope),
		ForScopes:   make(map[*ast.ForStmt]*Scope),
		WhileScopes: make(map[*ast.WhileStmt]*Scope),
	}
	if prog == nil {
		return res
	}

	reg := NewFunctionRegistry()

	// Pass 1: register function declarations up front to allow forward references.
	for _, stmt := range prog.Statements {
		if fn, ok := stmt.(*ast.FnDecl); ok {
			if err := ValidateIdentifier(fn.Name, fn.P); err != nil {
				res.Errors = append(res.Errors, err)
			}
			if err := res.Global.Define(fn.Name, fn.P); err != nil {
				res.Errors = append(res.Errors, err)
			}
			if err := reg.Define(fn.Name, len(fn.Params), fn.P); err != nil {
				res.Errors = append(res.Errors, err)
			}
		}
	}

	// Pass 2: analyze statements with scopes and registered functions.
	for _, stmt := range prog.Statements {
		analyzeStmt(stmt, res.Global, reg, &res, 0, limit)
	}

	return res
}

// AnalyzeDefinitions walks the AST to enforce semantic rules such as reserved-name protection.
// It collects all errors (never panicking) and returns them for reporting.
func AnalyzeDefinitions(prog *ast.Program) AnalysisResult {
	return AnalyzeDefinitionsWithLimit(prog, 0)
}

// Analyze preserves backward-compatible API returning only errors.
func Analyze(prog *ast.Program) []error {
	return AnalyzeDefinitions(prog).Errors
}

// NewAnalyzer creates an Analyzer with an optional depth limit (<=0 means no limit).
func NewAnalyzer(prog *ast.Program, depthLimit int) *Analyzer {
	return &Analyzer{prog: prog, limit: depthLimit}
}

// Run executes the semantic analysis, collecting all errors without panicking.
func (a *Analyzer) Run() {
	a.result = AnalyzeDefinitionsWithLimit(a.prog, a.limit)
}

// Errors returns the collected semantic errors.
func (a *Analyzer) Errors() []error {
	return a.result.Errors
}

// Result returns the full analysis result (scopes plus errors).
func (a *Analyzer) Result() AnalysisResult {
	return a.result
}

// Analyze runs semantic analysis on the provided program and returns an aggregated error if any.
func (a *Analyzer) Analyze(prog *ast.Program) error {
	a.prog = prog
	a.Run()
	return aggregateErrors(a.result.Errors)
}

func aggregateErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}
	return errors.Join(errs...)
}

func analyzeStmt(stmt ast.Statement, scope *Scope, reg *FunctionRegistry, res *AnalysisResult, depth, limit int) {
	if exceeded := checkDepth(stmt.Pos(), depth, limit); exceeded != nil {
		res.Errors = append(res.Errors, exceeded)
		return
	}
	switch s := stmt.(type) {
	case *ast.SetStmt:
		if err := ValidateIdentifier(s.Name, s.P); err != nil {
			res.Errors = append(res.Errors, err)
		}
		if err := scope.Define(s.Name, s.P); err != nil {
			res.Errors = append(res.Errors, err)
		}
		analyzeExpr(s.Value, scope, res, depth+1, limit)
	case *ast.FnDecl:
		// Name already validated/registered in pass 1; still validate params and body.
		fnScope := NewScope(scope)
		for _, param := range s.Params {
			if err := ValidateIdentifier(param, s.P); err != nil {
				res.Errors = append(res.Errors, err)
			}
			if err := fnScope.Define(param, s.P); err != nil {
				res.Errors = append(res.Errors, err)
			}
		}
		res.FuncScopes[s] = fnScope
		for _, inner := range s.Body {
			analyzeStmt(inner, fnScope, reg, res, depth+1, limit)
		}
	case *ast.IfStmt:
		analyzeExpr(s.Cond, scope, res, depth+1, limit)
		thenScope := NewScope(scope)
		for _, inner := range s.Then {
			analyzeStmt(inner, thenScope, reg, res, depth+1, limit)
		}
		if len(s.Else) > 0 {
			elseScope := NewScope(scope)
			for _, inner := range s.Else {
				analyzeStmt(inner, elseScope, reg, res, depth+1, limit)
			}
		}
	case *ast.ForStmt:
		loopScope := NewScope(scope)
		if err := ValidateIdentifier(s.Var, s.P); err != nil {
			res.Errors = append(res.Errors, err)
		}
		if err := loopScope.Define(s.Var, s.P); err != nil {
			res.Errors = append(res.Errors, err)
		}
		res.ForScopes[s] = loopScope
		analyzeExpr(s.Start, scope, res, depth+1, limit)
		analyzeExpr(s.End, scope, res, depth+1, limit)
		for _, inner := range s.Body {
			analyzeStmt(inner, loopScope, reg, res, depth+1, limit)
		}
	case *ast.WhileStmt:
		analyzeExpr(s.Cond, scope, res, depth+1, limit)
		bodyScope := NewScope(scope)
		res.WhileScopes[s] = bodyScope
		for _, inner := range s.Body {
			analyzeStmt(inner, bodyScope, reg, res, depth+1, limit)
		}
	case *ast.CallStmt:
		if arity, ok := reg.Lookup(s.Name); !ok {
			res.Errors = append(res.Errors, UndefinedVariableError{Name: s.Name, P: s.P})
		} else if arity != len(s.Args) {
			res.Errors = append(res.Errors, InvalidArityError{Name: s.Name, Expected: arity, Got: len(s.Args), P: s.P})
		}
		for _, arg := range s.Args {
			analyzeExpr(arg, scope, res, depth+1, limit)
		}
	case *ast.EchoStmt:
		analyzeExpr(s.Value, scope, res, depth+1, limit)
	case *ast.RunStmt:
		analyzeExpr(s.Command, scope, res, depth+1, limit)
	case *ast.ReturnStmt:
		if s.Value != nil {
			analyzeExpr(s.Value, scope, res, depth+1, limit)
		}
	case *ast.BreakStmt, *ast.ContinueStmt:
		// nothing to validate
	}
}

func analyzeExpr(expr ast.Expr, scope *Scope, res *AnalysisResult, depth, limit int) {
	if expr == nil {
		return
	}
	if exceeded := checkDepth(expr.Pos(), depth, limit); exceeded != nil {
		res.Errors = append(res.Errors, exceeded)
		return
	}
	switch e := expr.(type) {
	case *ast.IdentExpr:
		if IsReserved(e.Name) {
			return
		}
		if _, ok := scope.Lookup(e.Name); !ok {
			res.Errors = append(res.Errors, UndefinedVariableError{Name: e.Name, P: e.P})
		}
	case *ast.IndexExpr:
		analyzeExpr(e.Left, scope, res, depth+1, limit)
		analyzeExpr(e.Index, scope, res, depth+1, limit)
	case *ast.PropertyExpr:
		analyzeExpr(e.Object, scope, res, depth+1, limit)
	case *ast.BinaryExpr:
		analyzeExpr(e.Left, scope, res, depth+1, limit)
		analyzeExpr(e.Right, scope, res, depth+1, limit)
	case *ast.UnaryExpr:
		analyzeExpr(e.Right, scope, res, depth+1, limit)
	case *ast.ListLit:
		for _, el := range e.Elements {
			analyzeExpr(el, scope, res, depth+1, limit)
		}
	case *ast.MapLit:
		for _, p := range e.Pairs {
			analyzeExpr(p.Value, scope, res, depth+1, limit)
		}
	case *ast.ExistsCond:
		analyzeExpr(e.Path, scope, res, depth+1, limit)
	case *ast.StringLit, *ast.NumberLit, *ast.BoolLit:
		return
	}
}

func checkDepth(pos ast.Pos, depth, limit int) error {
	if limit > 0 && depth > limit {
		return DepthExceededError{Limit: limit, P: pos}
	}
	return nil
}
