# Fin — Language & Compiler System Instructions

> This document is the authoritative system specification for all human and AI agents working on the **Fin** project.
> Fin is a production-grade, Fish-inspired DSL that transpiles into Windows Batch (.bat) scripts.
>
> Any code, tooling, refactors, tests, or extensions MUST follow the rules defined here unless explicitly approved.

---

## 1. Project Mission

Fin exists to:

* Provide a readable, safe, predictable shell language for Windows automation.
* Compile deterministically into Windows Batch with zero runtime dependencies.
* Serve as a serious compiler engineering project with clean architecture, strong diagnostics, and testability.

Fin is **not** a toy DSL, REPL toy, or scripting experiment. All contributions must preserve:

* Determinism
* Explicitness
* Backward compatibility
* Tooling friendliness

---

## 2. Repository Architecture

```
fin/
 ├─ cmd/fin/               # CLI entry
 ├─ internal/
 │   ├─ token/             # Token definitions
 │   ├─ lexer/             # Scanner
 │   ├─ parser/            # Parser
 │   ├─ ast/               # Syntax tree
 │   ├─ sema/              # Semantic analysis (future)
 │   ├─ generator/         # Batch codegen
 │   └─ diagnostics/       # Error formatting
 ├─ examples/
 ├─ tests/
 └─ AGENTS.md
```

Rules:

* `internal/*` packages must not import from each other cyclically.
* AST must remain parser-agnostic.
* Generator must never inspect tokens directly — only AST.
* CLI must not embed language logic.

---

## 3. Language Surface

### 3.1 File Format

* Extension: `.fin`
* UTF-8 text
* Newlines separate statements
* Blocks terminate using explicit `end`

---

### 3.2 Comments

```
# this is a comment
```

* Comments extend to end of line.
* Ignored by lexer.

---

### 3.3 Variables

```
set name "Vishnu"
set count 10
```

Rules:

* All values are strings at runtime.
* Variable references use `$name` in source.
* Shadowing allowed inside functions.

---

### 3.4 Expressions

Supported literals:

* String: `"hello"`
* Number: `123`
* Boolean: `true`, `false`
* List: `[1, 2, 3]`
* Map: `{name: "bob", age: 20}`

Operators:

```
!   unary not
**  exponentiation (right-associative, highest precedence)
* / multiplication
+ - addition
< <= > >= comparisons
== != equality
&& logical and
|| logical or
```

Indexing:

```
$nums[0]
```

Property access:

```
$user.name
```

Operator precedence is strictly defined and must not change without version bump.

---

### 3.5 Echo

```
echo "Hello $name"
```

---

### 3.6 Run Command

```
run "git status"
```

---

### 3.7 If

```
if exists "file.txt"
    echo "found"
else
    echo "missing"
end
```

---

### 3.8 For Loop

```
for i in 1..5
    echo $i
end
```

---

### 3.9 While Loop

```
while x > 0
    set x x - 1
end

While conditions accept expressions.
```

---

### 3.10 Functions

Definition:

```
fn greet name times
    echo "Hello $name"
end
```

Call:

```
greet "Bob" 3
```

Rules:

* Positional arguments only.
* Parameters are local variables.
* No return values in v1.

---

### 3.11 Return (future)

```
return expr
```

Return values are reserved for future versions.

---

## 4. Grammar (Authoritative)

```
program → { statement }

statement → setStmt
          | echoStmt
          | runStmt
          | ifStmt
          | forStmt
          | whileStmt
          | fnDecl
          | returnStmt
          | callStmt
          | NEWLINE

setStmt → "set" IDENT expr NEWLINE
echoStmt → "echo" expr NEWLINE
runStmt → "run" STRING NEWLINE

ifStmt → "if" condition NEWLINE block ["else" NEWLINE block] "end" NEWLINE
forStmt → "for" IDENT "in" expr ".." expr NEWLINE block "end" NEWLINE
whileStmt → "while" expr NEWLINE block "end" NEWLINE

fnDecl → "fn" IDENT { IDENT } NEWLINE block "end" NEWLINE
callStmt → IDENT { expr } NEWLINE
returnStmt → "return" [expr] NEWLINE

block → { statement }

condition → "exists" STRING | expr

expr → logicalOr
logicalOr → logicalAnd { "||" logicalAnd }
logicalAnd → equality { "&&" equality }
equality → comparison { ("==" | "!=") comparison }
comparison → term { ("<" | "<=" | ">" | ">=") term }
term → factor { ("+" | "-") factor }
factor → unary { ("*" | "/") unary }
unary → ("!" | "-") unary | primary
primary → NUMBER | STRING | TRUE | FALSE | IDENT | list | map | index | property | "(" expr ")"

list → "[" [expr {"," expr}] "]"
map → "{" [pair {"," pair}] "}"
pair → IDENT ":" expr
index → primary "[" expr "]"
property → primary "." IDENT
```

---

## 5. Lexer Rules

* Single-pass, rune-based scanner.
* No regex usage.
* O(n) complexity.
* Emits NEWLINE tokens.
* Tracks line and column.
* Skips comments.
* Variable references: `$name` lexes as IDENT with literal `name` (leading `$` is not kept). `$` not followed by a letter yields an ILLEGAL token.

Tokens include:

* Keywords: set, echo, run, if, else, end, for, while, fn, return, in, exists, true, false
* Literals: IDENT, STRING, NUMBER
* Operators: + - * / == != < <= > >= && || ! .. .
* Delimiters: [ ] { } , :

Lexer must never panic.

---

## 6. AST Rules

* AST nodes must contain source position.
* No tokens inside AST.
* Interfaces: Node, Statement, Expr, Condition.
* AST must remain immutable after parsing.

---

## 7. Parser Rules

* Recursive descent with Pratt expression parsing.
* Error recovery via synchronization.
* Must continue parsing after error.
* Must not panic.
* Collects all errors.

---

## 8. Semantic Rules

* Undefined variables produce errors.
* Function arity must match.
* Shadowing allowed inside function.
* Reserved names protected.
* Duplicate function names forbidden.

PART I — REAL GENERATOR ARCHITECTURE
🎯 Generator’s Responsibility (Authoritative)

The generator must:

Convert a semantically validated AST into deterministic, readable, correct Windows Batch code.

It must not:

re-parse

re-check semantics

inspect tokens

guess intent

🧱 High-Level Architecture
Validated AST
   ↓
Lowering (AST → IR-like ops)
   ↓
Batch Emission (IR → text)
   ↓
Final .bat


We separate concerns to keep this maintainable.

📁 Directory Structure
internal/generator/
 ├── generator.go      # public API
 ├── context.go        # state (labels, vars, scopes)
 ├── lower_stmt.go     # statements → ops
 ├── lower_expr.go     # expressions → batch fragments
 ├── runtime.go        # helpers (labels, setlocal)
 ├── emit.go           # string emission
 ├── names.go          # deterministic name mangling
 └── errors.go         # generator errors


This mirrors real compiler backends.

🧠 Core Concepts
1. Generator Context (Critical)

Holds all mutable state:

type Context struct {
    labelCounter int
    indent       int
    out          *strings.Builder

    varMap   map[string]string
    funcMap  map[string]string
}


Rules:

No globals

No randomness

Deterministic naming

2. Lowering Phase (AST → Ops)

We lower high-level constructs into batch-compatible structures.

Examples:

Fin construct	Batch lowering
function	label + CALL
while	labels + IF + GOTO
list	env var expansion
map	prefixed vars
if	IF (...) ELSE (...)

This is where language meaning is preserved.

3. Emission Phase (Ops → Text)

Only string formatting:

indentation

newlines

escaping

No logic.

🧾 Generator Public API

📄 generator.go

package generator

import "fin/internal/ast"

type Generator interface {
    Generate(p *ast.Program) (string, error)
}


Concrete implementation:

type BatchGenerator struct {
    ctx *Context
}

---

## 9. Batch Code Generation Rules

### Variables

```
set name "bob" → set name=bob
$name → %name%
```

---

### Lists

```
set nums [10,20]
```

Generates:

```
set nums_0=10
set nums_1=20
set nums_len=2
```

Index:

```
$nums[1] → %nums_1%
```

---

### Maps

```
set user {name:"bob"}
```

Generates:

```
set user_name=bob
```

Property:

```
$user.name → %user_name%
```

---

### Functions

```
fn greet name
    echo "Hi $name"
end

greet "Bob"
```

Generates:

```
call :greet Bob
goto :eof

:greet
setlocal
set name=%1
echo Hi %name%
endlocal
goto :eof
```

---

### If

```
if exists "file.txt"
```

Generates:

```
if exist file.txt (
```

---

### For

```
for i in 1..5
```

Generates:

```
for /L %%i in (1,1,5) do (
```

---

### While

While loops must be lowered using labels and goto safely.

---

## 10. Tooling Expectations

* CLI: deterministic output
* Formatting stable
* No hidden magic
* Errors human readable

---

## 11. Compatibility Rules

* Grammar changes require version bump.
* Breaking changes must be documented.
* Backward compatibility preferred.

---

## 12. Agent Responsibilities

Agents must:

* Follow this specification strictly.
* Never introduce silent breaking behavior.
* Prefer correctness over cleverness.
* Add tests for every change.
* Keep code readable and documented.
