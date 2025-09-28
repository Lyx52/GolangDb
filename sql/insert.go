package sql

import (
	"fmt"
	"strings"
)

type InsertStatement struct {
	TableName *TableName
	Fields    []Field
	Sources   map[string]*TableName
}

func NewInsertStatement() *InsertStatement {
	return &InsertStatement{
		Sources: make(map[string]*TableName),
		Fields:  make([]Field, 0),
	}
}

func (statement InsertStatement) String() string {

	fieldNames := make([]string, 0)
	values := make([]string, 0)
	for _, field := range statement.Fields {
		fieldNames = append(fieldNames, field.String())

		values = append(values, fmt.Sprint(field.Value))
	}

	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);", statement.TableName.String(), strings.Join(fieldNames, ", "), strings.Join(values, ", "))
}

func (statement InsertStatement) Execute(connection *Connection) error {
	fmt.Printf("[EXECUTE] %v\n", statement.String())
	return nil
}

func (statement InsertStatement) GetTableName(alias string) (*TableName, error) {
	tableName, ok := statement.Sources[alias]
	if !ok {
		return nil, fmt.Errorf("alias %s not found", alias)
	}
	return tableName, nil
}

func (statement InsertStatement) PushTableName(alias string, tableName *TableName) error {
	tableName, ok := statement.Sources[alias]
	if ok {
		return fmt.Errorf("alias %s already exists for table %s", alias, tableName.Name)
	}
	statement.Sources[alias] = tableName

	return nil
}

func (statement InsertStatement) GetBaseTable() *TableName {
	return statement.TableName
}
