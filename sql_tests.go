package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Lyx52/GolangDb/server"
	"github.com/Lyx52/GolangDb/sql"
)

func RunAsCli(context *server.ServerContext) {
	context.RunDatabases()

	scanner := bufio.NewScanner(os.Stdin)
	parser := sql.NewSqlParser()
	var err error
	var statement sql.Statement
	for {
		fmt.Print("> ")

		if scanner.Scan() {
			input := scanner.Text()
			if input == "exit" {
				break
			}

			err = parser.Parse(input)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			statement = *parser.PopStatement()
			err = statement.Execute(context)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
		}

		err = scanner.Err()
		if err != nil {
			break
		}
	}

	close(context.Cancelled)
}

func main() {
	context := server.NewServerContext()

	parser := sql.NewSqlParser()
	err := parser.Parse("CREATE DATABASE dbs;")
	sql.RaiseIfError(err)
	statement := *parser.PopStatement()
	err = statement.Execute(context)
	if err != nil {
		return
	}

	err = parser.Parse("USE dbs;")
	sql.RaiseIfError(err)
	statement = *parser.PopStatement()
	err = statement.Execute(context)
	if err != nil {
		return
	}

	err = parser.Parse("CREATE TABLE users (id integer, test integer, test2 integer, username varchar);")
	sql.RaiseIfError(err)
	statement = *parser.PopStatement()
	err = statement.Execute(context)
	if err != nil {
		return
	}

	for i := 0; i < 100; i++ {
		err = parser.Parse("INSERT INTO dbs.users AS u (u.id, u.username) VALUES (1, 'tests');")
		sql.RaiseIfError(err)
		statement = *parser.PopStatement()
		err = statement.Execute(context)
		if err != nil {
			return
		}
	}

	err = parser.Parse("UPDATE dbs.users AS u SET u.username = 'tests123' WHERE u.id IN (1)")
	sql.RaiseIfError(err)
	statement = *parser.PopStatement()
	err = statement.Execute(context)
	if err != nil {
		return
	}

	RunAsCli(context)

	//
	//err = parser.Parse("INSERT INTO dbs.users AS u (u.id, u.username) VALUES (1 'tests');")
	//sql.RaiseIfError(err)

	//
	//statement = *parser.PopStatement()
	//err = statement.Execute(nil)
	//if err != nil {
	//	return
	//}
	//
	//err = parser.Parse("select usr.id, usr.username, (SELECT SUM(c.cost) FROM costs AS c WHERE c.user_id = usr.id) FROM users AS usr WHERE usr.id = 1;")
	//sql.RaiseIfError(err)

}
