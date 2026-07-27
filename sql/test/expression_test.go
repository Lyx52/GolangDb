package test

import (
	"fmt"
	"testing"

	"github.com/Lyx52/GolangDb/parser"
)

func TestExpression(t *testing.T) {
	s := "a.test > 100"
	lexer := parser.NewLexer(&s)
	err := lexer.Tokenize()
	if err != nil {
		t.Fatal(err)
	}
	expr, err := parser.ParseExpression(lexer)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(expr)
}
