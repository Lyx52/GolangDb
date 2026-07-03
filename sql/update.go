package sql

import (
	"fmt"
	"strings"

	"github.com/Lyx52/GolangDb/models"
	"github.com/Lyx52/GolangDb/server"
)

type UpdateStatement struct {
	TableName *TableName
	Fields    []FieldValue
	Where     WhereStatement
	Sources   map[string]*TableName
}

func (statement UpdateStatement) String() string {
	//UPDATE users AS u SET u.username = 'tests', u.val = 123 WHERE u.id = 1 AND u.username = 'test' OR ((u.id = 123 AND u.username = 'tests') AND u.id = 232);
	sets := make([]string, 0)
	for _, field := range statement.Fields {
		sets = append(sets, field.SetFieldValueString())
	}
	if statement.Where != nil {
		return fmt.Sprintf("UPDATE %s SET %s WHERE %s;", statement.TableName.String(), strings.Join(sets, ", "), statement.Where.String())
	}

	return fmt.Sprintf("UPDATE %s SET %s;", statement.TableName.String(), strings.Join(sets, ", "))
}

func NewUpdateStatement() *UpdateStatement {
	return &UpdateStatement{
		Sources: make(map[string]*TableName),
		Fields:  make([]FieldValue, 0),
	}
}

func (statement UpdateStatement) Execute(context *server.ServerContext) error {
	fmt.Printf("[EXECUTE] %v\n", statement.String())
	err, database := context.DatabaseStore.GetDatabase(statement.TableName.Schema)
	if err != nil {
		return err
	}

	err, table := database.GetTable(statement.TableName.Name)
	if err != nil {
		return err
	}

	sets := make([]func(row *models.DataRow), len(statement.Fields))
	for i, source := range statement.Fields {
		field := table.Fields[source.Field.Name]
		err = field.ValidateValue(source.Value.Data)
		if err != nil {
			return err
		}

		sets[i] = func(row *models.DataRow) {
			row.Values[field.Index] = source.Value.Data
		}
	}

	//row := models.DataRow{
	//	Values: make([]models.ColumnValue, len(table.Fields)),
	//}
	//for _, source := range statement.Fields {
	//	field := table.Fields[source.Field.Name]
	//

	//
	//	row.Values[field.Index] = source.Value.Data
	//}
	//
	go func() {
		for _, row := range table.Datastore.Values {
			for _, set := range sets {
				set(row)
			}
		}
	}()
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
