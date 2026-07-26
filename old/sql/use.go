package sql

import (
	"fmt"

	"github.com/Lyx52/GolangDb/server"
)

type UseStatement struct {
	Database string
}

func NewUseStatement(database string) *UseStatement {
	return &UseStatement{
		Database: database,
	}
}

func (statement UseStatement) String() string {
	return fmt.Sprintf("USE %s;", statement.Database)
}

func (statement UseStatement) Execute(context *server.ServerContext) error {
	fmt.Printf("[EXECUTE] %v\n", statement.String())
	err, database := context.DatabaseStore.GetDatabase(statement.Database)
	if err != nil {
		return err
	}

	context.SetCurrentDatabase(database)
	return nil
}

func (statement UseStatement) GetTableName(alias string) (*TableName, error) {
	panic("implement me")
}

func (statement UseStatement) PushTableName(alias string, tableName *TableName) error {
	panic("implement me")
}

func (statement UseStatement) GetBaseTable() *TableName {
	panic("implement me")
}
