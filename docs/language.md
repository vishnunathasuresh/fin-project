# Fin Language Specification

**Version:** 1.0.0  
**Status:** Final  
**Last Updated:** February 1, 2026

Fin is a Fish-inspired DSL that compiles deterministically to Windows Batch scripts. This specification defines the complete language surface, grammar, and semantics.

---

## 1. File Format

| Property | Value |
|----------|-------|
| Extension | `.fin` |
| Encoding | UTF-8 |
| Line Ending | LF or CRLF |
| Statement Separator | Newline (`\n` or `\r\n`) |
| Block Terminator | `end` keyword |

**Example:**
```fin
# This is a comment
set name "Alice"
echo "Hello $name"
```

---

## 2. Lexical Elements

### Comments
```fin
# This is a line comment
# Comments extend to end of line
```
- Prefix: `#`
- No block comments
- Ignored by lexer

### Identifiers
```fin
set myVar 42
fn calculate x y
```
- Start with letter or underscore: `[A-Za-z_]`
- Continue with letters, digits, underscores: `[A-Za-z0-9_]*`
- Case-sensitive
- Maximum length: unlimited
- Reserved words: `set`, `echo`, `run`, `if`, `else`, `end`, `for`, `while`, `fn`, `return`, `in`, `exists`, `true`, `false`

### Strings
```fin
set msg "Hello, World!"
set empty ""
set escaped "Say \"hi\""
set dollar "Price: $$50"
```
- Delimited by double quotes (`"`)
- Support escape sequences:
  - `\"` → literal quote
  - `$$` → literal dollar sign
  - `\n` → interpreted by batch (literal `\n` in output)
- Variable interpolation: `$name`, `$obj.prop`, `$arr[i]`

### Numbers
```fin
set count 42
set negative -7
set decimal 3.14
set hex 0xFF
```
- Integer or float literals
- Signed (with `-` prefix)
- No type distinction at runtime (all stored as strings)

### Keywords & Literals
```fin
set flag true
set empty false
```
- `true`, `false` — Boolean literals
- `in`, `exists` — Keywords for loop and condition syntax

### Operators & Delimiters
| Operator | Meaning |
|----------|---------|
| `+` | Addition |
| `-` | Subtraction (or unary negation) |
| `*` | Multiplication |
| `/` | Division |
| `%` | Modulo |
| `**` | Exponentiation (right-associative) |
| `<` | Less than |
| `<=` | Less than or equal |
| `>` | Greater than |
| `>=` | Greater than or equal |
| `==` | Equality |
| `!=` | Inequality |
| `&&` | Logical AND |
| `\|\|` | Logical OR |
| `!` | Logical NOT |
| `..` | Range operator |
| `.` | Property access |
| `[` `]` | Index/subscript |
| `(` `)` | Grouping |
| `{` `}` | Map/object literal |
| `,` | List/map separator |
| `:` | Map key-value separator |
| `=` | Assignment |

---

## 3. Data Types

All values are **strings at runtime**. Type checking happens at compile-time (semantic analysis).

### String Literal
```fin
set name "Alice"
echo $name
```
Runtime value: `Alice`

### Number Literal
```fin
set count 42
set sum $count + 8
```
- Stored as string `"42"`
- Arithmetic operations work on numeric interpretation
- Result is string

### Boolean Literal
```fin
set flag true
if $flag
    echo "yes"
end
```
- `true` → string `"true"`
- `false` → string `"false"`
- Comparable in conditionals

### List Literal
```fin
set colors [red, green, blue]
set nums [1, 2, 3]
```
- Syntax: `[expr, expr, ...]`
- Compiled to indexed environment variables:
  - `colors_0` = `red`
  - `colors_1` = `green`
  - `colors_2` = `blue`
  - `colors_len` = `3`
- Access: `$colors[0]`, `$colors[$i]`

### Map Literal
```fin
set person {name: "Bob", age: 30}
```
- Syntax: `{key: expr, key: expr, ...}`
- Compiled to named environment variables:
  - `person_name` = `Bob`
  - `person_age` = `30`
- Access: `$person.name`, `$person.age`
- Keys must be identifiers (not arbitrary strings)

---

## 4. Expressions

### Precedence (High to Low)

| Precedence | Operator | Associativity |
|-----------|----------|---------------|
| 10 | `(expr)` (grouping) | N/A |
| 9 | `!` (unary not), `-` (unary minus) | Right |
| 8 | `**` (exponentiation) | Right |
| 7 | `*`, `/`, `%` | Left |
| 6 | `+`, `-` (binary) | Left |
| 5 | `<`, `<=`, `>`, `>=` | Left |
| 4 | `==`, `!=` | Left |
| 3 | `&&` | Left |
| 2 | `\|\|` | Left |
| 1 | `,` (separator) | N/A |

### Primary Expressions

#### Literal
```fin
echo 42
echo "text"
echo true
echo [a, b, c]
echo {x: 1, y: 2}
```

#### Variable Reference
```fin
set x 10
echo $x
```
- Prefix: `$`
- Must be a valid identifier

#### Property Access
```fin
set user {name: "Alice", id: 123}
echo $user.name
```
- Syntax: `$obj.property`
- Property must be identifier (no expressions)

#### Index Access
```fin
set colors [red, green, blue]
echo $colors[0]
echo $colors[$i]
```
- Syntax: `$arr[expr]`
- Index can be literal number or variable
- Out-of-bounds: no error (generates empty or variable)

#### Existence Condition
```fin
if exists "file.txt"
    echo "found"
end
```
- Syntax: `exists "path"`
- Returns true/false
- Only valid in `if` conditions

#### Grouped Expression
```fin
set result $(2 + 3) * 4
```
- Syntax: `(expr)`

### Binary Operations

#### Arithmetic
```fin
set sum $a + $b
set diff $a - $b
set prod $a * $b
set quot $a / $b
set rem $a % $b
set power $a ** $b
```
- Operands coerced to numeric
- Division: integer division (batch behavior)

#### Comparison
```fin
if $x < 10
if $y >= 5
if $a == $b
if $x != $y
```
- Returns true/false strings
- Works with numeric or string values

#### Logical
```fin
if $a && $b
if $x || $y
if !$flag
```
- Short-circuit evaluation (in generated code)
- Operands: `true`/`false`

#### Range
```fin
for i in 1 .. 10
    echo $i
end
```
- Syntax: `start .. end`
- Inclusive on both ends
- Used in `for` loops

### String Interpolation

Inside string literals:

```fin
set name "World"
echo "Hello, $name!"          # → Hello, World!
echo "$user.email"             # → alice@example.com (if user.email = alice@example.com)
echo "$items[0]"               # → first (if items[0] = first)
echo "Price: $$100"            # → Price: $100 (literal $)
```

**Rules:**
- `$identifier` → replaced with variable value
- `$ident.property` → replaced with `ident_property`
- `$ident[expr]` → replaced with `ident_N` (if expr is literal) or dynamic access
- `$$` → literal `$`
- Outside strings, `$` is literal

---

## 5. Statements

### Set Statement
```fin
set name "Alice"
set x 42
set items [1, 2, 3]
set person {id: 1, name: "Bob"}
```
- Syntax: `set IDENT expr NEWLINE`
- Defines new variable
- If already defined: overwrites
- Lists/maps expand to multiple variables

### Assignment Statement
```fin
set x 10
x = 20                         # x now 20
```
- Syntax: `IDENT = expr NEWLINE`
- Updates existing variable
- Must be previously defined (semantic error if not)
- Preferred over `set` when updating

### Echo Statement
```fin
echo "Hello"
echo $name
echo "Count: $count"
```
- Syntax: `echo expr NEWLINE`
- Outputs expression to console
- String interpolation applied
- Special characters escaped for batch

### Run Statement
```fin
run "git status"
run "dir C:\temp"
```
- Syntax: `run STRING NEWLINE`
- Executes shell command
- Command must be string literal (no interpolation)
- Used for side effects (build, deployment)

### If Statement
```fin
if $x > 5
    echo "large"
else
    echo "small"
end
```
- Syntax: `if COND ... [else ...] end NEWLINE`
- Condition: expression or `exists "path"`
- Else branch optional
- Supports numeric comparisons (`<`, `>`, `<=`, `>=`)
- Supports equality (`==`, `!=`)
- Proper batch delayed expansion handling

### For Loop
```fin
for i in 1 .. 5
    echo "Iteration $i"
end
```
- Syntax: `for IDENT in expr .. expr ... end NEWLINE`
- Loop variable: local to body
- Range: inclusive, e.g., `1 .. 5` = [1,2,3,4,5]
- Numeric ranges only
- No break/continue (reserved for future)

### While Loop
```fin
set n 5
while $n > 0
    echo $n
    set n $n - 1
end
```
- Syntax: `while COND ... end NEWLINE`
- Condition: expression or `exists "path"`
- Supports numeric comparisons
- No break/continue (reserved for future)

### Function Declaration
```fin
fn greet name
    echo "Hello, $name!"
end

fn add a b
    set sum $a + $b
    echo "Sum: $sum"
end

greet "Alice"
add 2 3
```
- Syntax: `fn IDENT [IDENT ...] ... end NEWLINE`
- Parameters: positional only, no defaults
- Body: list of statements
- Variables: locals shadow globals
- Recursion: fully supported
- Return values: v1.0 does not support explicit return values
- Return: `return` statement jumps to function end

### Return Statement
```fin
fn check x
    if $x < 0
        echo "negative"
        return
    end
    echo "non-negative"
end
```
- Syntax: `return [expr] NEWLINE`
- Only valid inside functions
- Jumps to function end
- Expression: reserved for future (currently ignored)

### Function Call
```fin
greet "Alice"
add 1 2
```
- Syntax: `IDENT [expr ...] NEWLINE`
- Arguments: positional, separated by space
- Arity must match declaration
- Recursive calls supported

---

## 6. Grammar (Canonical)

```
program           → { statement }

statement         → setStmt
                  | assignStmt
                  | echoStmt
                  | runStmt
                  | ifStmt
                  | forStmt
                  | whileStmt
                  | fnDecl
                  | returnStmt
                  | callStmt
                  | NEWLINE

setStmt           → "set" IDENT expr NEWLINE
assignStmt        → IDENT "=" expr NEWLINE
echoStmt          → "echo" expr NEWLINE
runStmt           → "run" STRING NEWLINE

ifStmt            → "if" condition NEWLINE block
                    ["else" NEWLINE block] "end" NEWLINE

condition         → "exists" STRING | expr

forStmt           → "for" IDENT "in" expr ".." expr NEWLINE
                    block "end" NEWLINE

whileStmt         → "while" expr NEWLINE block "end" NEWLINE

fnDecl            → "fn" IDENT [IDENT ...] NEWLINE
                    block "end" NEWLINE

returnStmt        → "return" [expr] NEWLINE

callStmt          → IDENT [expr ...] NEWLINE

block             → { statement }

expr              → logicalOr

logicalOr         → logicalAnd { "||" logicalAnd }

logicalAnd        → equality { "&&" equality }

equality          → comparison { ("==" | "!=") comparison }

comparison        → term { ("<" | "<=" | ">" | ">=") term }

term              → factor { ("+" | "-") factor }

factor            → unary { ("*" | "/" | "%") unary }

unary             → ("!" | "-") unary | exponent

exponent          → primary { "**" primary }

primary           → NUMBER
                  | STRING
                  | TRUE
                  | FALSE
                  | IDENT
                  | list
                  | map
                  | index
                  | property
                  | exists
                  | "(" expr ")"

list              → "[" [expr {"," expr}] "]"

map               → "{" [pair {"," pair}] "}"

pair              → IDENT ":" expr

index             → primary "[" expr "]"

property          → primary "." IDENT

exists            → "exists" STRING
```

---

## 7. Semantic Rules

### Variable Scope
- **Global scope:** Variables defined at top level
- **Function scope:** Parameters and locals shadow globals
- **Local lifetime:** Function parameters and local variables exist only during execution
- **Shadowing:** Allowed inside functions; outer scope restored after return

### Function Rules
- **Parameters:** Positional only, no type annotations
- **Arity:** Must match exactly (no defaults, no varargs)
- **Recursion:** Fully supported with proper local scoping
- **Return:** Implicit at end; `return` jumps to end early
- **Forward references:** Functions can call functions defined later

### Type Coercion
- **Arithmetic:** Operands coerced to numbers; result is string
- **Comparison:** Values compared numerically if both are numeric-looking
- **Boolean:** Truthy = non-empty string; falsy = `false` keyword
- **String interpolation:** Variables interpolated only inside strings

### Errors (Compile-Time)
- Undefined variable reference
- Function call arity mismatch
- Duplicate function definition
- Reserved name used as variable
- Return outside function
- Invalid syntax

---

## 8. Canonical Formatting

`fin fmt` produces deterministic output:

```fin
# Comment style
set x 1

# Variables
set name "Alice"
set count 42

# Lists
set items [a, b, c]

# Maps
set person {name: "Bob", age: 30}

# Control flow
if $x > 5
    echo "large"
else
    echo "small"
end

for i in 1 .. 10
    echo $i
end

while $n > 0
    echo $n
    set n $n - 1
end

# Functions
fn greet name
    echo "Hello, $name!"
end

fn add a b
    echo $a + $b
end

# Expressions
set sum $a + $b
set product $x * $y
set expr $(a + b) * c
```

**Rules:**
- 4-space indentation
- Blank line before top-level functions
- Spaces around binary operators
- Space after keywords (`if`, `for`, `while`, `fn`)
- Space around `..` in ranges
- Maps/lists formatted as `{k: v}`, `[a, b]`
- Commas followed by space
- Colons in maps followed by space

---

## 9. Batch Code Generation

When compiled, Fin generates Windows Batch code with:

- `setlocal EnableDelayedExpansion` for variable expansion
- `!var!` syntax for delayed expansion
- `set /a var=expr` for arithmetic
- `if ... ( ) else ( )` for conditionals
- `for /L` for numeric loops
- `:label` and `goto` for function calls and loop control
- `call set` for indirect variable access

Example:
```fin
set x 10
if $x > 5
    echo "large"
end
```

Generates:
```batch
@echo off
setlocal EnableDelayedExpansion
set x=10
if !x! GTR 5 (
    echo large
)
endlocal
```

---

## 10. Version & Compatibility

| Version | Date | Status | Notes |
|---------|------|--------|-------|
| 1.0.0 | 2026-02-01 | Final | All core features. No return values yet. |

**Stability Promise:**
- Minor syntax changes require major version bump
- New features in minor versions (backward compatible)
- Bug fixes in patch versions

---

## 11. Limitations

### Not Supported
- Return values from functions (v1.1)
- Module/import system (v1.1)
- Closures
- Variadic functions
- Type annotations
- Generics
- Async/await
- Exception handling

### Windows-Only
- Targets Windows Batch only
- Requires Windows to run generated `.bat` files
- No POSIX shell output

### Batch Limitations (Inherited)
- String length limits (batch constraint)
- No true 64-bit arithmetic
- Environment variable restrictions apply
