# Tests

Fin keeps tests adjacent to their packages to allow access to unexported helpers. Key locations:

- CLI integration tests: `cmd/fin/main_test.go`
- Generator tests: `internal/generator/*_test.go`
- Parser tests: `internal/parser/*_test.go`
- Semantic analyzer tests: `internal/sema/*_test.go`
- AST utils tests: `internal/ast/print_test.go`
- Additional parser tokenization: `tests/parser/tokenize_test.go`

Rationale: Go’s visibility rules require tests that touch unexported helpers to live in the same package directory. A separate `tests/` tree would force exporting internals or adding shims; we avoid that to keep the internal surface tight.
