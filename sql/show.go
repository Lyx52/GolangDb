package sql

import (
	"fmt"

	"github.com/Lyx52/GolangDb/backing"
)

type ShowObjectType int

const (
	SHOW_DATABASES ShowObjectType = iota
	SHOW_TABLES    ShowObjectType = iota
)

func (showType ShowObjectType) String() string {
	switch showType {
	case SHOW_DATABASES:
		return "DATABASES"
	case SHOW_TABLES:
		return "TABLES"
	default:
		return ""
	}
}

type ShowStatement struct {
	Type ShowObjectType
}

func NewShowStatement(objectType ShowObjectType) *ShowStatement {
	return &ShowStatement{
		Type: objectType,
	}
}

func (statement ShowStatement) String() string {
	return fmt.Sprintf("SHOW %s;", statement.Type.String())
}

func (statement ShowStatement) Execute(context *backing.ServerContext) error {
	fmt.Printf("[EXECUTE] %v\n", statement.String())
	var err error = nil
	if statement.Type == SHOW_DATABASES {
		PrintDatabases(context)
	}

	if statement.Type == SHOW_TABLES {
		err = PrintTables(context)
	}

	return err
}

func PrintDatabases(context *backing.ServerContext) {
	fmt.Println("Databases")
	fmt.Println("----------------------------")
	for _, database := range context.DatabaseStore.Databases {
		fmt.Println(database.String())
	}
	fmt.Println("----------------------------")
}

func PrintTables(context *backing.ServerContext) error {
	err := context.CheckDatabaseConnected()
	if err != nil {
		return err
	}
	err, tableStore := context.DatabaseStore.GetDatabaseTableStore(context.CurrentDatabase.Name)

	if err != nil {
		return err
	}

	fmt.Println("Tables")
	fmt.Println("----------------------------")
	for _, table := range tableStore.Tables {
		fmt.Println(table.String())
	}
	fmt.Println("----------------------------")

	return nil
}

func (statement ShowStatement) GetTableName(alias string) (*TableName, error) {
	panic("implement me")
}

func (statement ShowStatement) PushTableName(alias string, tableName *TableName) error {
	panic("implement me")
}

func (statement ShowStatement) GetBaseTable() *TableName {
	panic("implement me")
}
