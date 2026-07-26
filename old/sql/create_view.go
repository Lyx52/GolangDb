package sql

import "github.com/Lyx52/GolangDb/server"

type CreateViewStatement struct {
	Name string
}

func NewCreateViewStatement(name string) *CreateViewStatement {
	return &CreateViewStatement{
		Name: name,
	}
}

func (statement CreateViewStatement) String() string {
	return statement.Name
}

func (statement CreateViewStatement) Execute(context *server.ServerContext) error {
	return nil
}
