package sema

import (
    "github.com/vishnunath-suresh/fin-project/internal/ast"
    "github.com/vishnunath-suresh/fin-project/internal/token"
)

// reservedNames contains keywords and builtins that cannot be used as identifiers.
var reservedNames map[string]struct{}

func init() {
    reservedNames = make(map[string]struct{}, len(token.Keywords))
    for k := range token.Keywords {
        reservedNames[k] = struct{}{}
    }
}

// IsReserved reports whether the given identifier is reserved.
func IsReserved(name string) bool {
    _, ok := reservedNames[name]
    return ok
}

// ValidateIdentifier ensures the provided name is not reserved. Returns ReservedNameError on violation.
func ValidateIdentifier(name string, pos ast.Pos) error {
    if IsReserved(name) {
        return ReservedNameError{Name: name, P: pos}
    }
    return nil
}
