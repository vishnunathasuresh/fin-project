package ir

// Program represents the entire IR
type Program struct {
	Types     map[string]*TypeDef
	Functions map[string]*Function
	Globals   []Var
}

// TypeDef is a user-defined type
type TypeDef struct {
	Name   string
	Fields []Field
}

// Function is an IR function
type Function struct {
	Name       string
	Params     []Param
	ReturnType Type
	Locals     []Var
	Body       []Stmt
}

// Stmt represents IR statements
type Stmt interface {
	irStmt()
}

// Expression types
type Expr interface {
	irExpr()
}

// Specific statement types
type DeclStmt struct {
	Name string
	Type Type
	Init Expr
}

func (s *DeclStmt) irStmt() {}

type AssignStmt struct {
	Name  string
	Value Expr
}

func (s *AssignStmt) irStmt() {}

type IfStmt struct {
	Cond Expr
	Then []Stmt
	Else []Stmt
}

func (s *IfStmt) irStmt() {}

type ForStmt struct {
	Var   string
	Start Expr
	End   Expr
	Body  []Stmt
}

func (s *ForStmt) irStmt() {}

type WhileStmt struct {
	Cond Expr
	Body []Stmt
}

func (s *WhileStmt) irStmt() {}

type RunStmt struct {
	Platform string // bash, fish, bat, ps1
	Cmd      Expr   // command expression
	OutVar   string // variable for stdout
	ErrVar   string // variable for error
}

func (s *RunStmt) irStmt() {}

type ReturnStmt struct {
	Value Expr
}

func (s *ReturnStmt) irStmt() {}

type BreakStmt struct{}

func (s *BreakStmt) irStmt() {}

type ContinueStmt struct{}

func (s *ContinueStmt) irStmt() {}

// Expression types
type IntLit struct {
	Value int
}

func (e *IntLit) irExpr() {}

type FloatLit struct {
	Value float64
}

func (e *FloatLit) irExpr() {}

type StringLit struct {
	Value string
}

func (e *StringLit) irExpr() {}

type BoolLit struct {
	Value bool
}

func (e *BoolLit) irExpr() {}

type Ident struct {
	Name string
	Type Type
}

func (e *Ident) irExpr() {}

type BinaryOp struct {
	Op    string
	Left  Expr
	Right Expr
	Type  Type
}

func (e *BinaryOp) irExpr() {}

type UnaryOp struct {
	Op   string
	Expr Expr
	Type Type
}

func (e *UnaryOp) irExpr() {}

type CallExpr struct {
	Func string
	Args []Expr
	Type Type
}

func (e *CallExpr) irExpr() {}

type CommandLit struct {
	Command string
}

func (e *CommandLit) irExpr() {}

type ListLit struct {
	Elements []Expr
}

func (e *ListLit) irExpr() {}

type MapLit struct {
	Keys   []Expr
	Values []Expr
}

func (e *MapLit) irExpr() {}

type IndexExpr struct {
	Object Expr
	Index  Expr
}

func (e *IndexExpr) irExpr() {}

type PropertyExpr struct {
	Object   Expr
	Property string
}

func (e *PropertyExpr) irExpr() {}

// Type represents IR types
type Type interface {
	TypeString() string
}

type BasicType struct {
	Kind string // "int", "float", "bool", "str"
}

func (t *BasicType) TypeString() string {
	return t.Kind
}

type ListType struct {
	ElemType Type
}

func (t *ListType) TypeString() string {
	return "list[" + t.ElemType.TypeString() + "]"
}

type MapType struct {
	KeyType   Type
	ValueType Type
}

func (t *MapType) TypeString() string {
	return "map[" + t.KeyType.TypeString() + ", " + t.ValueType.TypeString() + "]"
}

type CommandType struct{}

func (t *CommandType) TypeString() string {
	return "command"
}

type ErrorType struct{}

func (t *ErrorType) TypeString() string {
	return "error"
}

type Var struct {
	Name string
	Type Type
}

type Param struct {
	Name string
	Type Type
}

type Field struct {
	Name string
	Type Type
}
