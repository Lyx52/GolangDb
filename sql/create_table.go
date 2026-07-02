package sql

import (
	"fmt"

	"github.com/Lyx52/GolangDb/backing"
)

type CreateTableStatement struct {
	Name string
}

func NewCreateTableStatement(name string) *CreateTableStatement {
	return &CreateTableStatement{
		Name: name,
	}
}

func (statement CreateTableStatement) String() string {
	return statement.Name
}

func (statement CreateTableStatement) Execute(context *backing.ServerContext) error {
	err, table := context.CreateTable(statement.Name)
	fmt.Printf("CREATED TABLE %v\n", table.Name)
	return err
}
