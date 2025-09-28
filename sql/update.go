package sql

import (
	"fmt"
	"strings"
)

type UpdateStatement struct {
	TableName *TableName
	Fields    []Field
	Where     WhereStatement
	Sources   map[string]*TableName
}

func (statement UpdateStatement) String() string {
	//UPDATE users AS u SET u.username = 'tests', u.val = 123 WHERE u.id = 1 AND u.username = 'test' OR ((u.id = 123 AND u.username = 'tests') AND u.id = 232);
	sets := make([]string, 0)
	for _, field := range statement.Fields {
		sets = append(sets, fmt.Sprintf("%s = %s", field.String(), fmt.Sprint(field.Value)))
	}
	if statement.Where != nil {
		return fmt.Sprintf("UPDATE %s SET %s WHERE %s;", statement.TableName.String(), strings.Join(sets, ", "), statement.Where.String())
	}

	return fmt.Sprintf("UPDATE %s SET %s;", statement.TableName.String(), strings.Join(sets, ", "))
}

func NewUpdateStatement() *UpdateStatement {
	return &UpdateStatement{
		Sources: make(map[string]*TableName),
		Fields:  make([]Field, 0),
	}
}

func (statement UpdateStatement) Execute(connection *Connection) error {
	fmt.Printf("[EXECUTE] %v\n", statement.String())
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
