package sema

import "github.com/vishnunath-suresh/fin-project/internal/ast"

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

// AnalyzeDefinitions walks the AST to enforce semantic rules such as reserved-name protection.
// It collects all errors (never panicking) and returns them for reporting.
func AnalyzeDefinitions(prog *ast.Program) AnalysisResult {
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
		analyzeStmt(stmt, res.Global, reg, &res)
	}

	return res
}

// Analyze preserves backward-compatible API returning only errors.
func Analyze(prog *ast.Program) []error {
	return AnalyzeDefinitions(prog).Errors
}

func analyzeStmt(stmt ast.Statement, scope *Scope, reg *FunctionRegistry, res *AnalysisResult) {
	switch s := stmt.(type) {
	case *ast.SetStmt:
		if err := ValidateIdentifier(s.Name, s.P); err != nil {
			res.Errors = append(res.Errors, err)
		}
		if err := scope.Define(s.Name, s.P); err != nil {
			res.Errors = append(res.Errors, err)
		}
		analyzeExpr(s.Value, scope, res)
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
			analyzeStmt(inner, fnScope, reg, res)
		}
	case *ast.IfStmt:
		analyzeExpr(s.Cond, scope, res)
		thenScope := NewScope(scope)
		for _, inner := range s.Then {
			analyzeStmt(inner, thenScope, reg, res)
		}
		if len(s.Else) > 0 {
			elseScope := NewScope(scope)
			for _, inner := range s.Else {
				analyzeStmt(inner, elseScope, reg, res)
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
		analyzeExpr(s.Start, scope, res)
		analyzeExpr(s.End, scope, res)
		for _, inner := range s.Body {
			analyzeStmt(inner, loopScope, reg, res)
		}
	case *ast.WhileStmt:
		analyzeExpr(s.Cond, scope, res)
		bodyScope := NewScope(scope)
		res.WhileScopes[s] = bodyScope
		for _, inner := range s.Body {
			analyzeStmt(inner, bodyScope, reg, res)
		}
	case *ast.CallStmt:
		if arity, ok := reg.Lookup(s.Name); !ok {
			res.Errors = append(res.Errors, UndefinedVariableError{Name: s.Name, P: s.P})
		} else if arity != len(s.Args) {
			res.Errors = append(res.Errors, InvalidArityError{Name: s.Name, Expected: arity, Got: len(s.Args), P: s.P})
		}
		for _, arg := range s.Args {
			analyzeExpr(arg, scope, res)
		}
	case *ast.EchoStmt:
		analyzeExpr(s.Value, scope, res)
	case *ast.RunStmt:
		analyzeExpr(s.Command, scope, res)
	case *ast.ReturnStmt:
		if s.Value != nil {
			analyzeExpr(s.Value, scope, res)
		}
	case *ast.BreakStmt, *ast.ContinueStmt:
		// nothing to validate
	}
}

func analyzeExpr(expr ast.Expr, scope *Scope, res *AnalysisResult) {
	switch e := expr.(type) {
	case *ast.IdentExpr:
		if IsReserved(e.Name) {
			return
		}
		if _, ok := scope.Lookup(e.Name); !ok {
			res.Errors = append(res.Errors, UndefinedVariableError{Name: e.Name, P: e.P})
		}
	case *ast.IndexExpr:
		analyzeExpr(e.Left, scope, res)
		analyzeExpr(e.Index, scope, res)
	case *ast.PropertyExpr:
		analyzeExpr(e.Object, scope, res)
	case *ast.BinaryExpr:
		analyzeExpr(e.Left, scope, res)
		analyzeExpr(e.Right, scope, res)
	case *ast.UnaryExpr:
		analyzeExpr(e.Right, scope, res)
	case *ast.ListLit:
		for _, el := range e.Elements {
			analyzeExpr(el, scope, res)
		}
	case *ast.MapLit:
		for _, p := range e.Pairs {
			analyzeExpr(p.Value, scope, res)
		}
	case *ast.ExistsCond:
		analyzeExpr(e.Path, scope, res)
	case *ast.StringLit, *ast.NumberLit, *ast.BoolLit:
		return
	}
}
