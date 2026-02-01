package ast

// Pos represents a source position.
type Pos struct {
	Line   int
	Column int
}

//
// ---- Core Node Interfaces ----
//

type Node interface {
	Pos() Pos
	node()
}

type Statement interface {
	Node
	stmt()
}

type Expr interface {
	Node
	expr()
}

//
// ---- Program Root ----
//

type Program struct {
	Statements []Statement
	P          Pos
}

func (p *Program) Pos() Pos { return p.P }
func (*Program) node()      {}

//
// ---- Statements ----
//

type AssignStmt struct {
	Name  string
	Value Expr
	P     Pos
}

func (s *AssignStmt) Pos() Pos { return s.P }
func (*AssignStmt) node()      {}
func (*AssignStmt) stmt()      {}

type CallStmt struct {
	Name string
	Args []Expr
	P    Pos
}

func (s *CallStmt) Pos() Pos { return s.P }
func (*CallStmt) node()      {}
func (*CallStmt) stmt()      {}

type FnDecl struct {
	Name   string
	Params []string
	Body   []Statement
	P      Pos
}

func (s *FnDecl) Pos() Pos { return s.P }
func (*FnDecl) node()      {}
func (*FnDecl) stmt()      {}

type IfStmt struct {
	Cond Expr
	Then []Statement
	Else []Statement
	P    Pos
}

func (s *IfStmt) Pos() Pos { return s.P }
func (*IfStmt) node()      {}
func (*IfStmt) stmt()      {}

type ForStmt struct {
	Var   string
	Start Expr
	End   Expr
	Body  []Statement
	P     Pos
}

func (s *ForStmt) Pos() Pos { return s.P }
func (*ForStmt) node()      {}
func (*ForStmt) stmt()      {}

type WhileStmt struct {
	Cond Expr
	Body []Statement
	P    Pos
}

func (s *WhileStmt) Pos() Pos { return s.P }
func (*WhileStmt) node()      {}
func (*WhileStmt) stmt()      {}

type ReturnStmt struct {
	Value Expr // optional; nil means bare return
	P     Pos
}

func (s *ReturnStmt) Pos() Pos { return s.P }
func (*ReturnStmt) node()      {}
func (*ReturnStmt) stmt()      {}

type BreakStmt struct {
	P Pos
}

func (s *BreakStmt) Pos() Pos { return s.P }
func (*BreakStmt) node()      {}
func (*BreakStmt) stmt()      {}

type ContinueStmt struct {
	P Pos
}

func (s *ContinueStmt) Pos() Pos { return s.P }
func (*ContinueStmt) node()      {}
func (*ContinueStmt) stmt()      {}

//
// ---- Conditions ----
//

type ExistsCond struct {
	Path Expr
	P    Pos
}

func (c *ExistsCond) Pos() Pos { return c.P }
func (*ExistsCond) node()      {}
func (*ExistsCond) expr()      {}

//
// ---- Expressions ----
//

type IdentExpr struct {
	Name string
	P    Pos
}

func (e *IdentExpr) Pos() Pos { return e.P }
func (*IdentExpr) node()      {}
func (*IdentExpr) expr()      {}

type StringLit struct {
	Value string
	P     Pos
}

func (e *StringLit) Pos() Pos { return e.P }
func (*StringLit) node()      {}
func (*StringLit) expr()      {}

type NumberLit struct {
	Value string
	P     Pos
}

func (e *NumberLit) Pos() Pos { return e.P }
func (*NumberLit) node()      {}
func (*NumberLit) expr()      {}

type ListLit struct {
	Elements []Expr
	P        Pos
}

func (e *ListLit) Pos() Pos { return e.P }
func (*ListLit) node()      {}
func (*ListLit) expr()      {}

type MapPair struct {
	Key   string
	Value Expr
	P     Pos
}

func (p *MapPair) Pos() Pos { return p.P }
func (*MapPair) node()      {}

type MapLit struct {
	Pairs []MapPair
	P     Pos
}

func (e *MapLit) Pos() Pos { return e.P }
func (*MapLit) node()      {}
func (*MapLit) expr()      {}

type IndexExpr struct {
	Left  Expr
	Index Expr
	P     Pos
}

func (e *IndexExpr) Pos() Pos { return e.P }
func (*IndexExpr) node()      {}
func (*IndexExpr) expr()      {}

type PropertyExpr struct {
	Object Expr
	Field  string
	P      Pos
}

func (e *PropertyExpr) Pos() Pos { return e.P }
func (*PropertyExpr) node()      {}
func (*PropertyExpr) expr()      {}

type BinaryExpr struct {
	Left  Expr
	Op    string
	Right Expr
	P     Pos
}

func (e *BinaryExpr) Pos() Pos { return e.P }
func (*BinaryExpr) node()      {}
func (*BinaryExpr) expr()      {}

type UnaryExpr struct {
	Op    string
	Right Expr
	P     Pos
}

func (e *UnaryExpr) Pos() Pos { return e.P }
func (*UnaryExpr) node()      {}
func (*UnaryExpr) expr()      {}

type BoolLit struct {
	Value bool
	P     Pos
}

func (e *BoolLit) Pos() Pos { return e.P }
func (*BoolLit) node()      {}
func (*BoolLit) expr()      {}
