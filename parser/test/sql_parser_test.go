package test

import (
	"fmt"
	"testing"

	"github.com/Lyx52/GolangDb/parser"
)

func TestSqlParser(test *testing.T) {
	test.Run("Simple SELECT", func(test *testing.T) {
		p := parser.NewSqlParser()
		s := `SELECT * FROM asgsg AS test 
         WHERE (a.test > 100) AND (a.b < 50) OR (aahhh.test != '12313') 
         ORDER BY ahhh.test, ahhh.asd, ahhh.zzz ASC`
		a := fmt.Sprintf(`SELECT * FROM asgsg AS test, (%s) AS GSGS`, s)
		b := fmt.Sprintf(`SELECT * FROM asgsg AS test, (%s) AS GSGS`, a)
		c := fmt.Sprintf(`SELECT * FROM asgsg AS "asdadsad", (%s) AS GSGS`, b)
		d := fmt.Sprintf(`SELECT *, sss[1:2], sdsf[0] FROM asgsg AS test, (%s) AS GSGS WHERE aaaaa >= 123`, c)
		err := p.Parse(d)
		if err != nil {
			test.Fatal(err)
		}
		statement := p.PopStatement()
		test.Log(statement)
	})
}
