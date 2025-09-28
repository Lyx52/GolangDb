package main

import "github.com/Lyx52/GolangDb/sql"

func main() {
	parser := sql.NewSqlParser()
	err := parser.Parse("INSERT INTO dbs.users AS u (u.id, u.username) VALUES (1 'tests');")
	sql.RaiseIfError(err)
	err = parser.Parse("UPDATE users AS u SET u.username = 'tests', u.val = 123 WHERE u.id = 1 AND u.username = 'test' OR ((u.id = 123 AND u.username = 'tests') AND u.id = 232);")
	sql.RaiseIfError(err)
	statement := *parser.PopStatement()
	err = statement.Execute(nil)
	if err != nil {
		return
	}

	err = parser.Parse("select usr.id, usr.username, (SELECT SUM(c.cost) FROM costs AS c WHERE c.user_id = usr.id) FROM users AS usr WHERE usr.id = 1;")
	sql.RaiseIfError(err)

}
