package test

import (
	"testing"

	"github.com/Lyx52/GolangDb/parser"
)

type ValueExpectedCount struct {
	Value    string
	Expected int
}

type ValueExpected struct {
	Value    string
	Expected string
}

type ValueExpectedType struct {
	Value    string
	Expected parser.TokenType
}

func TestLexer(test *testing.T) {
	test.Run("Lexer test numbers", func(test *testing.T) {

		values := []ValueExpected{
			{"123", "123"},
			{"123.123", "123.123"},
			{"-123.123", "-123.123"},
			{"0b0101", "5"},
			{"0B0111", "7"},
			{"0XF", "15"},
			{"0XFA", "250"},
			{"-0o123", "-83"},
		}
		var err error
		for _, value := range values {
			lexer := parser.NewLexer(&value.Value)
			err = lexer.Tokenize()
			if err != nil {
				test.Fatal(err)
			}

			if lexer.Count() != 1 {
				test.Fatalf("Expected 1 token, got %d", lexer.Count())
			}

			token := lexer.PopToken()
			if token.Value != value.Expected {
				test.Fatalf("Expected %s, got %s", value.Expected, token.Value)
			}
		}
	})

	test.Run("Lexer test whitespace trim", func(test *testing.T) {

		values := []ValueExpectedCount{
			{"131312     12313          2131   1223   123", 4},
			{"131312 123 123 123  123    ' 12312 1231 1223 21 2131'", 5},
		}
		var err error
		for _, value := range values {
			lexer := parser.NewLexer(&value.Value)
			err = lexer.Tokenize()
			if err != nil {
				test.Fatal(err)
			}
			count := 0
			next := lexer.PopToken()
			for next != nil {
				if next.Type == parser.WHITESPACE {
					count++
				}
				next = lexer.PopToken()
			}

			if count != value.Expected {
				test.Fatalf("Expected %d, got %d", value.Expected, count)
			}
		}
	})

	test.Run("Lexer test number, string, identifier, quoted identifier", func(test *testing.T) {
		values := []ValueExpectedType{
			{"'asddad'", parser.STRING},
			{`"sadasda"`, parser.QUOTED_IDENTIFIER},
			{"`aasdasda`", parser.QUOTED_IDENTIFIER},
			{"asdadasdasdsadasd", parser.IDENTIFIER},
			{"_asdadasdasdsadasd", parser.IDENTIFIER},
		}
		var err error
		for _, value := range values {
			lexer := parser.NewLexer(&value.Value)
			err = lexer.Tokenize()
			if err != nil {
				test.Fatal(err)
			}

			if lexer.Count() != 1 {
				test.Fatalf("Expected 1 token, got %d", lexer.Count())
			}

			token := lexer.PopToken()
			if token.Type != value.Expected {
				test.Fatalf("Expected %s, got %s", value.Expected, token.Value)
			}
		}
	})
}
