package ir

import (
	"fmt"
	"strconv"

	"github.com/vishnunathasuresh/fin-project/internal/ast"
)

// Lowerer converts AST to IR
type Lowerer struct {
	prog      *Program
	currentFn *Function
	errors    []error
}

// Lower converts an AST program to IR
func Lower(astProg *ast.Program) (*Program, error) {
	l := &Lowerer{
		prog: &Program{
			Types:     make(map[string]*TypeDef),
			Functions: make(map[string]*Function),
			Globals:   []Var{},
		},
		errors: []error{},
	}

	err := l.lowerProgram(astProg)
	if err != nil {
		return nil, err
	}

	if len(l.errors) > 0 {
		return nil, l.errors[0]
	}

	return l.prog, nil
}

func (l *Lowerer) lowerProgram(p *ast.Program) error {
	// For now, just lower function declarations
	// v1 doesn't have type declarations or methods yet
	for _, stmt := range p.Statements {
		if fn, ok := stmt.(*ast.FnDecl); ok {
			irFn, err := l.lowerFnDecl(fn)
			if err != nil {
				l.errors = append(l.errors, err)
				continue
			}
			l.prog.Functions[fn.Name] = irFn
		}
	}

	return nil
}

func (l *Lowerer) lowerFnDecl(fn *ast.FnDecl) (*Function, error) {
	// v1 functions don't have typed params yet, so we'll use placeholder types
	params := []Param{}
	for _, paramName := range fn.Params {
		params = append(params, Param{
			Name: paramName,
			Type: &BasicType{Kind: "any"}, // Placeholder until we have type system
		})
	}

	l.currentFn = &Function{
		Name:       fn.Name,
		Params:     params,
		ReturnType: nil, // v1 doesn't have return types
		Locals:     []Var{},
		Body:       []Stmt{},
	}

	for _, stmt := range fn.Body {
		irStmt, err := l.lowerStmt(stmt)
		if err != nil {
			l.errors = append(l.errors, err)
			continue
		}
		if irStmt != nil {
			l.currentFn.Body = append(l.currentFn.Body, irStmt)
		}
	}

	result := l.currentFn
	l.currentFn = nil
	return result, nil
}

func (l *Lowerer) lowerStmt(s ast.Statement) (Stmt, error) {
	switch stmt := s.(type) {
	case *ast.SetStmt:
		// v1 set statement -> DeclStmt
		initExpr, err := l.lowerExpr(stmt.Value)
		if err != nil {
			return nil, err
		}
		return &DeclStmt{
			Name: stmt.Name,
			Type: &BasicType{Kind: "any"},
			Init: initExpr,
		}, nil

	case *ast.AssignStmt:
		valueExpr, err := l.lowerExpr(stmt.Value)
		if err != nil {
			return nil, err
		}
		return &AssignStmt{
			Name:  stmt.Name,
			Value: valueExpr,
		}, nil

	case *ast.IfStmt:
		return l.lowerIfStmt(stmt)

	case *ast.ForStmt:
		return l.lowerForStmt(stmt)

	case *ast.WhileStmt:
		return l.lowerWhileStmt(stmt)

	case *ast.ReturnStmt:
		return l.lowerReturnStmt(stmt)

	case *ast.BreakStmt:
		return &BreakStmt{}, nil

	case *ast.ContinueStmt:
		return &ContinueStmt{}, nil

	case *ast.EchoStmt:
		// v1 echo -> will eventually become run() call in v2
		// For now, skip or handle specially
		return nil, nil

	case *ast.RunStmt:
		// v1 run -> will eventually become typed run() call in v2
		// For now, skip or handle specially
		return nil, nil

	case *ast.CallStmt:
		// Function call as statement
		args := []Expr{}
		for _, arg := range stmt.Args {
			argExpr, err := l.lowerExpr(arg)
			if err != nil {
				return nil, err
			}
			args = append(args, argExpr)
		}
		callExpr := &CallExpr{
			Func: stmt.Name,
			Args: args,
			Type: &BasicType{Kind: "any"},
		}
		// Wrap in expression statement (not defined yet, but we'll handle it)
		_ = callExpr
		return nil, nil

	default:
		return nil, fmt.Errorf("unknown statement type: %T", s)
	}
}

func (l *Lowerer) lowerIfStmt(s *ast.IfStmt) (Stmt, error) {
	condExpr, err := l.lowerExpr(s.Cond)
	if err != nil {
		return nil, err
	}

	thenStmts := []Stmt{}
	for _, stmt := range s.Then {
		irStmt, err := l.lowerStmt(stmt)
		if err != nil {
			return nil, err
		}
		if irStmt != nil {
			thenStmts = append(thenStmts, irStmt)
		}
	}

	elseStmts := []Stmt{}
	for _, stmt := range s.Else {
		irStmt, err := l.lowerStmt(stmt)
		if err != nil {
			return nil, err
		}
		if irStmt != nil {
			elseStmts = append(elseStmts, irStmt)
		}
	}

	return &IfStmt{
		Cond: condExpr,
		Then: thenStmts,
		Else: elseStmts,
	}, nil
}

func (l *Lowerer) lowerForStmt(s *ast.ForStmt) (Stmt, error) {
	startExpr, err := l.lowerExpr(s.Start)
	if err != nil {
		return nil, err
	}

	endExpr, err := l.lowerExpr(s.End)
	if err != nil {
		return nil, err
	}

	bodyStmts := []Stmt{}
	for _, stmt := range s.Body {
		irStmt, err := l.lowerStmt(stmt)
		if err != nil {
			return nil, err
		}
		if irStmt != nil {
			bodyStmts = append(bodyStmts, irStmt)
		}
	}

	return &ForStmt{
		Var:   s.Var,
		Start: startExpr,
		End:   endExpr,
		Body:  bodyStmts,
	}, nil
}

func (l *Lowerer) lowerWhileStmt(s *ast.WhileStmt) (Stmt, error) {
	condExpr, err := l.lowerExpr(s.Cond)
	if err != nil {
		return nil, err
	}

	bodyStmts := []Stmt{}
	for _, stmt := range s.Body {
		irStmt, err := l.lowerStmt(stmt)
		if err != nil {
			return nil, err
		}
		if irStmt != nil {
			bodyStmts = append(bodyStmts, irStmt)
		}
	}

	return &WhileStmt{
		Cond: condExpr,
		Body: bodyStmts,
	}, nil
}

func (l *Lowerer) lowerReturnStmt(s *ast.ReturnStmt) (Stmt, error) {
	if s.Value == nil {
		return &ReturnStmt{Value: nil}, nil
	}

	valueExpr, err := l.lowerExpr(s.Value)
	if err != nil {
		return nil, err
	}

	return &ReturnStmt{Value: valueExpr}, nil
}

func (l *Lowerer) lowerExpr(e ast.Expr) (Expr, error) {
	switch expr := e.(type) {
	case *ast.NumberLit:
		// Determine if int or float
		if containsDecimal(expr.Value) {
			f, err := strconv.ParseFloat(expr.Value, 64)
			if err != nil {
				return nil, err
			}
			return &FloatLit{Value: f}, nil
		}
		i, err := strconv.Atoi(expr.Value)
		if err != nil {
			return nil, err
		}
		return &IntLit{Value: i}, nil

	case *ast.StringLit:
		return &StringLit{Value: expr.Value}, nil

	case *ast.BoolLit:
		return &BoolLit{Value: expr.Value}, nil

	case *ast.IdentExpr:
		return &Ident{
			Name: expr.Name,
			Type: &BasicType{Kind: "any"},
		}, nil

	case *ast.BinaryExpr:
		return l.lowerBinaryExpr(expr)

	case *ast.UnaryExpr:
		operandExpr, err := l.lowerExpr(expr.Right)
		if err != nil {
			return nil, err
		}
		return &UnaryOp{
			Op:   expr.Op,
			Expr: operandExpr,
			Type: &BasicType{Kind: "any"},
		}, nil

	case *ast.ListLit:
		return l.lowerListLit(expr)

	case *ast.MapLit:
		return l.lowerMapLit(expr)

	case *ast.IndexExpr:
		return l.lowerIndexExpr(expr)

	case *ast.PropertyExpr:
		return l.lowerPropertyExpr(expr)

	default:
		return nil, fmt.Errorf("unknown expression type: %T", e)
	}
}

func (l *Lowerer) lowerBinaryExpr(e *ast.BinaryExpr) (Expr, error) {
	leftExpr, err := l.lowerExpr(e.Left)
	if err != nil {
		return nil, err
	}

	rightExpr, err := l.lowerExpr(e.Right)
	if err != nil {
		return nil, err
	}

	return &BinaryOp{
		Op:    e.Op,
		Left:  leftExpr,
		Right: rightExpr,
		Type:  &BasicType{Kind: "any"},
	}, nil
}

func (l *Lowerer) lowerUnaryExpr(e *ast.UnaryExpr) (Expr, error) {
	operandExpr, err := l.lowerExpr(e.Right)
	if err != nil {
		return nil, err
	}

	return &UnaryOp{
		Op:   e.Op,
		Expr: operandExpr,
		Type: &BasicType{Kind: "any"},
	}, nil
}

func (l *Lowerer) lowerListLit(e *ast.ListLit) (Expr, error) {
	elements := []Expr{}
	for _, elem := range e.Elements {
		elemExpr, err := l.lowerExpr(elem)
		if err != nil {
			return nil, err
		}
		elements = append(elements, elemExpr)
	}

	return &ListLit{Elements: elements}, nil
}

func (l *Lowerer) lowerMapLit(e *ast.MapLit) (Expr, error) {
	keys := []Expr{}
	values := []Expr{}

	for _, pair := range e.Pairs {
		keyExpr := &StringLit{Value: pair.Key}
		valueExpr, err := l.lowerExpr(pair.Value)
		if err != nil {
			return nil, err
		}
		keys = append(keys, keyExpr)
		values = append(values, valueExpr)
	}

	return &MapLit{Keys: keys, Values: values}, nil
}

func (l *Lowerer) lowerIndexExpr(e *ast.IndexExpr) (Expr, error) {
	objectExpr, err := l.lowerExpr(e.Left)
	if err != nil {
		return nil, err
	}

	indexExpr, err := l.lowerExpr(e.Index)
	if err != nil {
		return nil, err
	}

	return &IndexExpr{
		Object: objectExpr,
		Index:  indexExpr,
	}, nil
}

func (l *Lowerer) lowerPropertyExpr(e *ast.PropertyExpr) (Expr, error) {
	objectExpr, err := l.lowerExpr(e.Object)
	if err != nil {
		return nil, err
	}

	return &PropertyExpr{
		Object:   objectExpr,
		Property: e.Field,
	}, nil
}

func containsDecimal(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	if err == nil {
		_, intErr := strconv.ParseInt(s, 10, 64)
		return intErr != nil
	}
	return false
}
