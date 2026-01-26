package main

import (
	"fmt"
	"github.com/vishnunath-suresh/fin-project/internal/lexer"
)

func main() {
	input := `
fn greet name
    set nums [1,2,3]
    set user {name: "bob", age: 20}
    echo "Hello $name
end
`

	l := lexer.New(input)

	for {
		tok := l.NextToken()
		fmt.Printf("%+v\n", tok)
		if tok.Type == "EOF" {
			break
		}
	}
}
