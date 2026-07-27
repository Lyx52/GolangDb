package test

import (
	"testing"

	"github.com/Lyx52/GolangDb/parser"
)

func TestSqlParser(test *testing.T) {
	test.Run("Simple SELECT", func(test *testing.T) {
		p := parser.NewSqlParser()
		//err := p.Parse("SELECT test.c AS test AS test AS test FROM abcd AS test")
		err := p.Parse("SELECT test.* FROM (SELECT * FROM abcd) AS test WHERE 1=1")
		if err != nil {
			test.Fatal(err)
		}
		statement := p.PopStatement()
		test.Log(statement)
	})
}
