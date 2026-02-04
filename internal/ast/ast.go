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

// TypeRef represents a resolved or annotated type name.
type TypeRef struct {
	Name string
}

//
// ---- Statements ----
//

// DeclStmt represents a declaration with optional type annotation.
type DeclStmt struct {
	Names []string // Can be a single name or multiple names for tuple unpacking
	Value Expr
	Type  *TypeRef
	P     Pos
}

func (s *DeclStmt) Pos() Pos { return s.P }
func (*DeclStmt) node()      {}
func (*DeclStmt) stmt()      {}

type AssignStmt struct {
	Names []string // Can be a single name or multiple names for tuple unpacking
	Value Expr
	Type  *TypeRef
	P     Pos
}

func (s *AssignStmt) Pos() Pos { return s.P }
func (*AssignStmt) node()      {}
func (*AssignStmt) stmt()      {}

type CallStmt struct {
	Name string
	Args []Expr
	Type *TypeRef
	P    Pos
}

func (s *CallStmt) Pos() Pos { return s.P }
func (*CallStmt) node()      {}
func (*CallStmt) stmt()      {}

type FnDecl struct {
	Name   string
	Params []Param
	Return *TypeRef
	Body   []Statement
	Type   *TypeRef
	P      Pos
}

func (s *FnDecl) Pos() Pos { return s.P }
func (*FnDecl) node()      {}
func (*FnDecl) stmt()      {}

type IfStmt struct {
	Cond Expr
	Then []Statement
	Else []Statement
	Type *TypeRef
	P    Pos
}

func (s *IfStmt) Pos() Pos { return s.P }
func (*IfStmt) node()      {}
func (*IfStmt) stmt()      {}

type ForStmt struct {
	Var      string
	Iterable Expr
	Body     []Statement
	Else     []Statement // optional else branch executed if loop not exited early (fin-v2)
	Type     *TypeRef
	P        Pos
}

func (s *ForStmt) Pos() Pos { return s.P }
func (*ForStmt) node()      {}
func (*ForStmt) stmt()      {}

type WhileStmt struct {
	Cond Expr
	Body []Statement
	Type *TypeRef
	P    Pos
}

func (s *WhileStmt) Pos() Pos { return s.P }
func (*WhileStmt) node()      {}
func (*WhileStmt) stmt()      {}

type ReturnStmt struct {
	Value Expr // optional; nil means bare return
	Type  *TypeRef
	P     Pos
}

func (s *ReturnStmt) Pos() Pos { return s.P }
func (*ReturnStmt) node()      {}
func (*ReturnStmt) stmt()      {}

type BreakStmt struct {
	Type *TypeRef
	P    Pos
}

func (s *BreakStmt) Pos() Pos { return s.P }
func (*BreakStmt) node()      {}
func (*BreakStmt) stmt()      {}

type ContinueStmt struct {
	Type *TypeRef
	P    Pos
}

// TypeDef represents a type declaration.
type TypeDef struct {
	Name   string
	Fields []Field
	Type   *TypeRef
	P      Pos
}

func (s *TypeDef) Pos() Pos { return s.P }
func (*TypeDef) node()      {}
func (*TypeDef) stmt()      {}

// MethodDecl represents a method with a receiver.
type MethodDecl struct {
	Receiver Param
	Name     string
	Params   []Param
	Return   *TypeRef
	Body     []Statement
	Type     *TypeRef
	P        Pos
}

func (s *MethodDecl) Pos() Pos { return s.P }
func (*MethodDecl) node()      {}
func (*MethodDecl) stmt()      {}

// Param is a named parameter with type.
type Param struct {
	Name string
	Type *TypeRef
	P    Pos
}

// Field is a named field with type.
type Field struct {
	Name string
	Type *TypeRef
	P    Pos
}

func (s *ContinueStmt) Pos() Pos { return s.P }
func (*ContinueStmt) node()      {}
func (*ContinueStmt) stmt()      {}

//
// ---- Conditions ----
//

type ExistsCond struct {
	Path Expr
	Type *TypeRef
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
	Type *TypeRef
	P    Pos
}

func (e *IdentExpr) Pos() Pos { return e.P }
func (*IdentExpr) node()      {}
func (*IdentExpr) expr()      {}

type StringLit struct {
	Value string
	Type  *TypeRef
	P     Pos
}

func (e *StringLit) Pos() Pos { return e.P }
func (*StringLit) node()      {}
func (*StringLit) expr()      {}

type NumberLit struct {
	Value string
	Type  *TypeRef
	P     Pos
}

func (e *NumberLit) Pos() Pos { return e.P }
func (*NumberLit) node()      {}
func (*NumberLit) expr()      {}

type ListLit struct {
	Elements []Expr
	Type     *TypeRef
	P        Pos
}

func (e *ListLit) Pos() Pos { return e.P }
func (*ListLit) node()      {}
func (*ListLit) expr()      {}

type MapPair struct {
	Key   string
	Value Expr
	Type  *TypeRef
	P     Pos
}

func (p *MapPair) Pos() Pos { return p.P }
func (*MapPair) node()      {}

type MapLit struct {
	Pairs []MapPair
	Type  *TypeRef
	P     Pos
}

func (e *MapLit) Pos() Pos { return e.P }
func (*MapLit) node()      {}
func (*MapLit) expr()      {}

type IndexExpr struct {
	Left  Expr
	Index Expr
	Type  *TypeRef
	P     Pos
}

func (e *IndexExpr) Pos() Pos { return e.P }
func (*IndexExpr) node()      {}
func (*IndexExpr) expr()      {}

type PropertyExpr struct {
	Object Expr
	Field  string
	Type   *TypeRef
	P      Pos
}

func (e *PropertyExpr) Pos() Pos { return e.P }
func (*PropertyExpr) node()      {}
func (*PropertyExpr) expr()      {}

type BinaryExpr struct {
	Left  Expr
	Op    string
	Right Expr
	Type  *TypeRef
	P     Pos
}

func (e *BinaryExpr) Pos() Pos { return e.P }
func (*BinaryExpr) node()      {}
func (*BinaryExpr) expr()      {}

type UnaryExpr struct {
	Op    string
	Right Expr
	Type  *TypeRef
	P     Pos
}

func (e *UnaryExpr) Pos() Pos { return e.P }
func (*UnaryExpr) node()      {}
func (*UnaryExpr) expr()      {}

type BoolLit struct {
	Value bool
	Type  *TypeRef
	P     Pos
}

// CommandLit captures raw command text.
type CommandLit struct {
	Text string
	Type *TypeRef
	P    Pos
}

func (e *CommandLit) Pos() Pos { return e.P }
func (*CommandLit) node()      {}
func (*CommandLit) expr()      {}

// NamedArg represents a named argument in a function call: name=value
type NamedArg struct {
	Name  string
	Value Expr
	P     Pos
}

func (e *NamedArg) Pos() Pos { return e.P }
func (*NamedArg) node()      {}
func (*NamedArg) expr()      {}

// CallExpr represents a function/method call used as an expression.
type CallExpr struct {
	Callee    Expr
	Args      []Expr     // positional arguments
	NamedArgs []NamedArg // named arguments (e.g., platform=bash, cmd=cmd)
	Type      *TypeRef
	P         Pos
}

func (e *CallExpr) Pos() Pos { return e.P }
func (*CallExpr) node()      {}
func (*CallExpr) expr()      {}

func (e *BoolLit) Pos() Pos { return e.P }
func (*BoolLit) node()      {}
func (*BoolLit) expr()      {}
