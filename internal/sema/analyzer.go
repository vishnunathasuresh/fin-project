package sema

import "github.com/vishnunath-suresh/fin-project/internal/ast"

// FunctionRegistry tracks function signatures by name.
type FunctionRegistry struct {
	funcs map[string]int
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

// Analyze walks the AST to enforce semantic rules such as reserved-name protection.
// It collects all errors (never panicking) and returns them for reporting.
func Analyze(prog *ast.Program) []error {
	if prog == nil {
		return nil
	}

	global := NewScope(nil)
	reg := NewFunctionRegistry()
	var errs []error

	// Pass 1: register function declarations up front to allow forward references.
	for _, stmt := range prog.Statements {
		if fn, ok := stmt.(*ast.FnDecl); ok {
			if err := ValidateIdentifier(fn.Name, fn.P); err != nil {
				errs = append(errs, err)
			}
			if err := global.Define(fn.Name); err != nil {
				errs = append(errs, err)
			}
			if err := reg.Define(fn.Name, len(fn.Params), fn.P); err != nil {
				errs = append(errs, err)
			}
		}
	}

	// Pass 2: analyze statements with scopes and registered functions.
	for _, stmt := range prog.Statements {
		analyzeStmt(stmt, global, reg, &errs)
	}

	return errs
}

func analyzeStmt(stmt ast.Statement, scope *Scope, reg *FunctionRegistry, errs *[]error) {
	switch s := stmt.(type) {
	case *ast.SetStmt:
		if err := ValidateIdentifier(s.Name, s.P); err != nil {
			*errs = append(*errs, err)
		}
		if err := scope.Define(s.Name); err != nil {
			*errs = append(*errs, err)
		}
		analyzeExpr(s.Value, scope, errs)
	case *ast.FnDecl:
		// Name already validated/registered in pass 1; still validate params and body.
		fnScope := NewScope(scope)
		for _, param := range s.Params {
			if err := ValidateIdentifier(param, s.P); err != nil {
				*errs = append(*errs, err)
			}
			if err := fnScope.Define(param); err != nil {
				*errs = append(*errs, err)
			}
		}
		for _, inner := range s.Body {
			analyzeStmt(inner, fnScope, reg, errs)
		}
	case *ast.IfStmt:
		analyzeExpr(s.Cond, scope, errs)
		thenScope := NewScope(scope)
		for _, inner := range s.Then {
			analyzeStmt(inner, thenScope, reg, errs)
		}
		if len(s.Else) > 0 {
			elseScope := NewScope(scope)
			for _, inner := range s.Else {
				analyzeStmt(inner, elseScope, reg, errs)
			}
		}
	case *ast.ForStmt:
		loopScope := NewScope(scope)
		if err := ValidateIdentifier(s.Var, s.P); err != nil {
			*errs = append(*errs, err)
		}
		if err := loopScope.Define(s.Var); err != nil {
			*errs = append(*errs, err)
		}
		analyzeExpr(s.Start, scope, errs)
		analyzeExpr(s.End, scope, errs)
		for _, inner := range s.Body {
			analyzeStmt(inner, loopScope, reg, errs)
		}
	case *ast.WhileStmt:
		analyzeExpr(s.Cond, scope, errs)
		bodyScope := NewScope(scope)
		for _, inner := range s.Body {
			analyzeStmt(inner, bodyScope, reg, errs)
		}
	case *ast.CallStmt:
		if arity, ok := reg.Lookup(s.Name); !ok {
			*errs = append(*errs, UndefinedVariableError{Name: s.Name, P: s.P})
		} else if arity != len(s.Args) {
			*errs = append(*errs, InvalidArityError{Name: s.Name, Expected: arity, Got: len(s.Args), P: s.P})
		}
		for _, arg := range s.Args {
			analyzeExpr(arg, scope, errs)
		}
	case *ast.EchoStmt:
		analyzeExpr(s.Value, scope, errs)
	case *ast.RunStmt:
		analyzeExpr(s.Command, scope, errs)
	case *ast.ReturnStmt:
		if s.Value != nil {
			analyzeExpr(s.Value, scope, errs)
		}
	case *ast.BreakStmt, *ast.ContinueStmt:
		// nothing to validate
	}
}

func analyzeExpr(expr ast.Expr, scope *Scope, errs *[]error) {
	switch e := expr.(type) {
	case *ast.IdentExpr:
		if _, ok := scope.Lookup(e.Name); !ok {
			*errs = append(*errs, UndefinedVariableError{Name: e.Name, P: e.P})
		}
	case *ast.IndexExpr:
		analyzeExpr(e.Left, scope, errs)
		analyzeExpr(e.Index, scope, errs)
	case *ast.PropertyExpr:
		analyzeExpr(e.Object, scope, errs)
	case *ast.BinaryExpr:
		analyzeExpr(e.Left, scope, errs)
		analyzeExpr(e.Right, scope, errs)
	case *ast.UnaryExpr:
		analyzeExpr(e.Right, scope, errs)
	case *ast.ListLit:
		for _, el := range e.Elements {
			analyzeExpr(el, scope, errs)
		}
	case *ast.MapLit:
		for _, p := range e.Pairs {
			analyzeExpr(p.Value, scope, errs)
		}
	case *ast.ExistsCond:
		analyzeExpr(e.Path, scope, errs)
	case *ast.StringLit, *ast.NumberLit, *ast.BoolLit:
		return
	}
}
