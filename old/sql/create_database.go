package sql

import (
	"fmt"

	"github.com/Lyx52/GolangDb/server"
)

type CreateDatabaseStatement struct {
	Name string
}

func NewCreateDatabaseStatement(name string) *CreateDatabaseStatement {
	return &CreateDatabaseStatement{
		Name: name,
	}
}

func (statement CreateDatabaseStatement) String() string {
	return statement.Name
}

func (statement CreateDatabaseStatement) Execute(context *server.ServerContext) error {
	err, database := context.CreateDatabase(statement.Name)
	if err != nil {
		return err
	}
	context.SetCurrentDatabase(database)
	fmt.Printf("CREATED DATABASE %v\n", database.Name)
	return nil
}
