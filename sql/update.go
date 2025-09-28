package sql

import "fmt"

type UpdateStatement struct {
	TableName *TableName
	Fields    []Field
	Where     *WhereStatement
	Sources   map[string]*TableName
}

func NewUpdateStatement() *UpdateStatement {
	return &UpdateStatement{
		Sources: make(map[string]*TableName),
		Fields:  make([]Field, 0),
	}
}

func (statement UpdateStatement) Execute(connection *Connection) error {
	return nil
}

func (statement UpdateStatement) GetTableName(alias string) (*TableName, error) {
	tableName, ok := statement.Sources[alias]
	if !ok {
		return nil, fmt.Errorf("alias %s not found", alias)
	}
	return tableName, nil
}

func (statement UpdateStatement) PushTableName(alias string, tableName *TableName) error {
	name, ok := statement.Sources[alias]
	if ok {
		return fmt.Errorf("alias %s already exists for table %s", alias, name.Name)
	}
	statement.Sources[alias] = tableName

	return nil
}

func (statement UpdateStatement) GetBaseTable() *TableName {
	return statement.TableName
}
