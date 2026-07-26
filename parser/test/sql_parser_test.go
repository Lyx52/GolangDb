package test

import (
	"testing"

	"github.com/Lyx52/GolangDb/parser"
)

func TestSqlParser(test *testing.T) {
	test.Run("Simple SELECT", func(test *testing.T) {
		p := parser.NewSqlParser()
		err := p.Parse("SELECT test.a, test.b, test.c FROM abcd AS test;    SELECT xxxx.a AS row, xxxx.b AS qweqweq, xxxx.c AS cdg FROM ffff AS xxxx")
		if err != nil {
			test.Fatal(err)
		}
	})
}
