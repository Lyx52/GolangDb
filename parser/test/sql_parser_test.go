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
		fmt.Println(d)
		//h := fmt.Sprintf(`SELECT * FROM asgsg AS "asdadsad" WHERE asass IN test.c;SELECT * FROM asgsg AS "asdadsad" WHERE asass IN test.c`)
		cc := `CREATE TABLE films (
					code        char(5),
					title       varchar(40),
					did         integer,
					date_prod   date,
					kind        varchar(10),
					len         interval hour to minute,
					CONSTRAINT production UNIQUE(date_prod)
				);`
		err := p.Parse(cc)
		if err != nil {
			test.Fatal(err)
		}
		statement := p.PopStatement()
		test.Log(statement)
	})
	test.Run("Select", func(test *testing.T) {
		parser.InitClauseParsers()
		//s := `INSERT INTO test (a,b,c,d) VALUES (1,2,3,'12313') RETURNING a AS sfsafsaf, c, d`
		// s := `DELETE FROM tasks WHERE status = 'DONE' RETURNING *;`

		s := `TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;TRUNCATE bigtable, fattable RESTART IDENTITY CASCADE;`
		lexer := parser.NewLexer(&s)
		err := lexer.Tokenize()

		if err != nil {
			test.Fatal(err)
		}

		expr, err := parser.ParseStatement(lexer, false)
		if err != nil {
			test.Fatal(err)
		}
		test.Log(expr)
	})
}
