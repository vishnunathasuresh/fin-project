package ir

import "fmt"

// Validator checks IR for correctness
type Validator struct {
	errors []error
}

// Validate checks IR program validity
func Validate(p *Program) error {
	v := &Validator{errors: []error{}}
	v.validateProgram(p)

	if len(v.errors) > 0 {
		return v.errors[0] // Return first error
	}
	return nil
}

func (v *Validator) validateProgram(p *Program) {
	// Check all type definitions are valid
	for name, td := range p.Types {
		if td.Name != name {
			v.errors = append(v.errors, fmt.Errorf("type name mismatch: %s vs %s", name, td.Name))
		}
		v.validateTypeDefFields(td)
	}

	// Check all functions are valid
	for name, fn := range p.Functions {
		if fn.Name != name {
			v.errors = append(v.errors, fmt.Errorf("function name mismatch: %s vs %s", name, fn.Name))
		}
		v.validateFunction(fn)
	}
}

func (v *Validator) validateTypeDefFields(td *TypeDef) {
	seen := make(map[string]bool)
	for _, field := range td.Fields {
		if seen[field.Name] {
			v.errors = append(v.errors, fmt.Errorf("duplicate field name: %s in type %s", field.Name, td.Name))
		}
		seen[field.Name] = true

		if field.Type == nil {
			v.errors = append(v.errors, fmt.Errorf("field %s in type %s has nil type", field.Name, td.Name))
		}
	}
}

func (v *Validator) validateFunction(fn *Function) {
	// Check parameters don't have duplicate names
	seen := make(map[string]bool)
	for _, param := range fn.Params {
		if seen[param.Name] {
			v.errors = append(v.errors, fmt.Errorf("duplicate parameter name: %s in function %s", param.Name, fn.Name))
		}
		seen[param.Name] = true

		if param.Type == nil {
			v.errors = append(v.errors, fmt.Errorf("parameter %s in function %s has nil type", param.Name, fn.Name))
		}
	}

	// Check function body
	for _, stmt := range fn.Body {
		v.validateStmt(stmt)
	}
}

func (v *Validator) validateStmt(stmt Stmt) {
	switch s := stmt.(type) {
	case *DeclStmt:
		if s.Type == nil {
			v.errors = append(v.errors, fmt.Errorf("declaration %s has nil type", s.Name))
		}
		if s.Init != nil {
			v.validateExpr(s.Init)
		}

	case *AssignStmt:
		if s.Value != nil {
			v.validateExpr(s.Value)
		}

	case *IfStmt:
		if s.Cond != nil {
			v.validateExpr(s.Cond)
		}
		for _, thenStmt := range s.Then {
			v.validateStmt(thenStmt)
		}
		for _, elseStmt := range s.Else {
			v.validateStmt(elseStmt)
		}

	case *ForStmt:
		if s.Start != nil {
			v.validateExpr(s.Start)
		}
		if s.End != nil {
			v.validateExpr(s.End)
		}
		for _, bodyStmt := range s.Body {
			v.validateStmt(bodyStmt)
		}

	case *WhileStmt:
		if s.Cond != nil {
			v.validateExpr(s.Cond)
		}
		for _, bodyStmt := range s.Body {
			v.validateStmt(bodyStmt)
		}

	case *ReturnStmt:
		if s.Value != nil {
			v.validateExpr(s.Value)
		}

	case *RunStmt:
		if s.Cmd != nil {
			v.validateExpr(s.Cmd)
		}
	}
}

func (v *Validator) validateExpr(expr Expr) {
	switch e := expr.(type) {
	case *BinaryOp:
		if e.Left != nil {
			v.validateExpr(e.Left)
		}
		if e.Right != nil {
			v.validateExpr(e.Right)
		}

	case *UnaryOp:
		if e.Expr != nil {
			v.validateExpr(e.Expr)
		}

	case *CallExpr:
		for _, arg := range e.Args {
			v.validateExpr(arg)
		}

	case *ListLit:
		for _, elem := range e.Elements {
			v.validateExpr(elem)
		}

	case *MapLit:
		for _, key := range e.Keys {
			v.validateExpr(key)
		}
		for _, val := range e.Values {
			v.validateExpr(val)
		}

	case *IndexExpr:
		if e.Object != nil {
			v.validateExpr(e.Object)
		}
		if e.Index != nil {
			v.validateExpr(e.Index)
		}

	case *PropertyExpr:
		if e.Object != nil {
			v.validateExpr(e.Object)
		}
	}
}
