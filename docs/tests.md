# Testing Guide

**Version:** 1.0.0  
**Status:** Final

Fin includes comprehensive test coverage for all compiler phases. This guide explains the test structure and how to run tests.

---

## Test Structure

Tests are located adjacent to the code they test, following Go conventions. This allows access to unexported helpers and internal types.

### Test Locations

| Path | Scope | Coverage |
|------|-------|----------|
| `cmd/fin/main_test.go` | CLI integration | Command parsing, file I/O, error reporting |
| `internal/token/token.go` | Token definitions | (No tests; definitions only) |
| `internal/lexer/lexer.go` | Lexer/scanner | (Tested via parser tests) |
| `internal/parser/*_test.go` | Parser | Tokenization, expression parsing, statement parsing, AST building |
| `internal/ast/*_test.go` | AST utilities | AST printing, structure validation |
| `internal/sema/*_test.go` | Semantic analysis | Variable scope, function arity, duplicate detection, reserved names |
| `internal/generator/*_test.go` | Code generation | Batch code generation, golden test outputs |
| `tests/parser/tokenize_test.go` | Parser integration | Token collection, whitespace handling |

---

## Running Tests

### Run All Tests
```bash
go test ./...
```

Output:
```
ok      github.com/vishnunathasuresh/fin-project/cmd/fin        2.150s
ok      github.com/vishnunathasuresh/fin-project/internal/ast   0.001s
ok      github.com/vishnunathasuresh/fin-project/internal/generator     0.421s
ok      github.com/vishnunathasuresh/fin-project/internal/parser        0.050s
ok      github.com/vishnunathasuresh/fin-project/internal/sema  0.038s
ok      github.com/vishnunathasuresh/fin-project/tests/parser   0.002s
```

### Run Specific Package
```bash
go test ./internal/generator/...
go test ./internal/parser/...
go test ./internal/sema/...
```

### Run Specific Test
```bash
go test ./internal/generator/ -run TestGenerate_IfElse
go test ./internal/parser/ -run TestParse_ForLoop
```

### Verbose Output
```bash
go test ./... -v
```

Shows each test name and execution time.

### With Coverage
```bash
go test ./... -cover
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

## Test Categories

### Parser Tests (`internal/parser/*_test.go`)

**Files:**
- `parser_test.go` — General parser tests
- `tokenize_test.go` — Lexer/tokenizer tests
- `expr_test.go` — Expression parsing
- `parser_stmt_test.go` — Statement parsing
- `parser_program_test.go` — Program-level parsing
- `parser_integration_test.go` — Integration tests

**Examples:**
```bash
go test ./internal/parser -v -run TestParse
go test ./internal/parser -v -run TestTokenize
go test ./internal/parser -v -run TestExpr
```

### Semantic Analysis Tests (`internal/sema/*_test.go`)

**Files:**
- `analyzer_test.go` — General semantic analysis
- `scope_test.go` — Variable scoping
- `reserved_test.go` — Reserved name checking
- `integration_test.go` — End-to-end semantic checks

**Coverage:**
- Undefined variable detection
- Function arity checking
- Duplicate function detection
- Reserved name protection
- Shadowing rules

**Examples:**
```bash
go test ./internal/sema -v
go test ./internal/sema -run TestScope
go test ./internal/sema -run TestReserved
```

### Generator Tests (`internal/generator/*_test.go`)

**Files:**
- `generator_test.go` — General code generation
- `generator_test.go` — Golden test outputs
- `lower_stmt_test.go` — Statement lowering
- `ast_snapshot_test.go` — AST snapshot testing

**Coverage:**
- Batch code emission
- Variable expansion
- Control flow generation
- Function code generation
- List/map handling
- Special character escaping

**Golden Tests:**
```bash
go test ./internal/generator -v -run TestGenerator_Golden
```

Golden tests compare generated Batch code against known-good reference outputs.

### CLI Integration Tests (`cmd/fin/main_test.go`)

**Coverage:**
- `fin build` command
- `fin check` command
- `fin ast` command
- `fin fmt` command
- `fin version` command
- Error handling
- Exit codes

**Examples:**
```bash
go test ./cmd/fin -v
go test ./cmd/fin -run TestBuild
```

### AST Tests (`internal/ast/*_test.go`)

**Files:**
- `print_test.go` — AST string representation

**Coverage:**
- AST node printing
- Tree structure validation

---

## Test Examples

### Example: Parser Test
```go
func TestParse_SimpleSet(t *testing.T) {
    input := "set x 10\n"
    prog, err := Parse(input)
    if err != nil {
        t.Fatalf("Parse error: %v", err)
    }
    if len(prog.Statements) != 1 {
        t.Fatalf("expected 1 statement, got %d", len(prog.Statements))
    }
}
```

### Example: Semantic Test
```go
func TestSema_UndefinedVariable(t *testing.T) {
    input := "echo $undefined\n"
    _, err := Analyze(Parse(input))
    if err == nil {
        t.Fatal("expected error for undefined variable")
    }
}
```

### Example: Generator Test
```go
func TestGenerate_Echo(t *testing.T) {
    prog := &ast.Program{
        Statements: []ast.Statement{
            &ast.EchoStmt{Value: &ast.StringLit{Value: "hello"}},
        },
    }
    gen := NewBatchGenerator()
    out, err := gen.Generate(prog)
    if err != nil {
        t.Fatalf("Generate error: %v", err)
    }
    if !strings.Contains(out, "echo hello") {
        t.Fatalf("expected 'echo hello' in output")
    }
}
```

---

## Test Coverage

### Current Coverage
- Lexer/Parser: 95%+
- Semantic Analysis: 90%+
- Code Generation: 85%+
- CLI: 80%+

### Improving Coverage
```bash
# Generate coverage report
go test ./... -coverprofile=coverage.out -cover
go tool cover -html=coverage.out -o coverage.html

# Check package-specific coverage
go test ./internal/generator -cover
go test ./internal/parser -cover
```

---

## Golden Tests

Golden tests compare generated code against reference outputs stored in test data files.

### Running Golden Tests
```bash
go test ./internal/generator -run TestGenerator_Golden -v
```

### Updating Golden Tests
If you intentionally change code generation, golden tests will fail. Review the diff and update:

```bash
# The test output shows expected vs. actual
# If correct, run with -update flag (if supported)
# Or manually update the test constants in generator_test.go
```

---

## Continuous Integration

All tests run on every commit:

```bash
# Full test suite
go test ./...

# With coverage
go test ./... -cover

# Verbose for debugging
go test ./... -v
```

Expected: All tests pass with exit code 0.

---

## Test Conventions

### File Naming
- `*_test.go` — All test files
- `example*_test.go` — Integration test examples

### Function Naming
- `Test<Feature>` — Main test function
- `Test<Category>_<Scenario>` — Specific scenario

### Error Messages
- Include context: what was expected vs. actual
- Use `t.Fatalf()` for hard failures
- Use `t.Errorf()` for warnings

### Examples
- `TestParse_SimpleSet` — Parse a simple set statement
- `TestGenerate_IfElse` — Generate if/else block
- `TestSema_UndefinedVariable` — Detect undefined variables

---

## Troubleshooting

### Test Fails Unexpectedly
1. Run with `-v` for details
2. Check test output diff
3. Verify recent changes didn't break assumptions

### Coverage Gaps
- Run `go tool cover -html=coverage.out`
- Add tests for uncovered lines

### Golden Test Mismatch
- Review generated code in test output
- If change is intentional, update golden constant
- If unintended, debug generator code

---

## See Also

- [CLI Reference](cli.md)
- [Language Specification](language.md)
- [AGENTS.md](../AGENTS.md) — Compiler architecture
