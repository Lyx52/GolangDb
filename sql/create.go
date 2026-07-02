package sql

import (
	"fmt"

	"github.com/Lyx52/GolangDb/backing"
)

type CreateObjectType int
type CreateObjectStatement interface {
	String() string
	Execute(context *backing.ServerContext) error
}

const (
	CREATE_DATABASE CreateObjectType = iota
	CREATE_TABLE    CreateObjectType = iota
	CREATE_VIEW     CreateObjectType = iota
)

func (createType CreateObjectType) String() string {
	switch createType {
	case CREATE_DATABASE:
		return "DATABASE"
	case CREATE_TABLE:
		return "TABLE"
	case CREATE_VIEW:
		return "VIEW"
	default:
		return ""
	}
}

type CreateStatement struct {
	Type         CreateObjectType
	SubStatement *CreateObjectStatement
}

func NewCreateStatement(createType CreateObjectType) *CreateStatement {
	return &CreateStatement{
		Type: createType,
	}
}

func (statement CreateStatement) String() string {
	return fmt.Sprintf("CREATE %s %s;", statement.Type.String(), (*statement.SubStatement).String())
}

func (statement CreateStatement) Execute(context *backing.ServerContext) error {
	fmt.Printf("[EXECUTE] %v\n", statement.String())
	return (*statement.SubStatement).Execute(context)
}

func (statement CreateStatement) GetTableName(alias string) (*TableName, error) {
	panic("implement me")
}

func (statement CreateStatement) PushTableName(alias string, tableName *TableName) error {
	panic("implement me")
}

func (statement CreateStatement) GetBaseTable() *TableName {
	panic("implement me")
}
