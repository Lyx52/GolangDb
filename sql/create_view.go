package sql

import "github.com/Lyx52/GolangDb/backing"

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

func (statement CreateViewStatement) Execute(context *backing.ServerContext) error {
	return nil
}
