package main

import (
	"fmt"

	"github.com/vishnunath-suresh/fin-project/internal/ast"
	"github.com/vishnunath-suresh/fin-project/internal/lexer"
	"github.com/vishnunath-suresh/fin-project/internal/parser"
)

func main() {
	input := `
fn greet name
    set nums [1,2,3]
    set user {name: "bob", age: 20}
    echo "Hello $name"
end
`
	l := lexer.New(input)
	toks := parser.CollectTokens(l)
	p := parser.New(toks)
	prog := p.ParseProgram()

	if errs := p.Errors(); len(errs) > 0 {
		fmt.Println("Parse errors:")
		for _, err := range errs {
			fmt.Println(" -", err)
		}
	}

	fmt.Println(ast.Format(prog))
}
