package sql

import (
	"fmt"
	"strings"

	"github.com/Lyx52/GolangDb/models"
	"github.com/Lyx52/GolangDb/server"
)

type InsertStatement struct {
	TableName *TableName
	Fields    []FieldValue
	Sources   map[string]*TableName
}

func NewInsertStatement() *InsertStatement {
	return &InsertStatement{
		Sources: make(map[string]*TableName),
		Fields:  make([]FieldValue, 0),
	}
}

func (statement InsertStatement) String() string {

	fieldNames := make([]string, 0)
	values := make([]string, 0)
	for _, field := range statement.Fields {
		fieldNames = append(fieldNames, field.Field.String())

		values = append(values, field.Value.String())
	}

	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);", statement.TableName.String(), strings.Join(fieldNames, ", "), strings.Join(values, ", "))
}

func (statement InsertStatement) Execute(context *server.ServerContext) error {
	fmt.Printf("[EXECUTE] %v\n", statement.String())
	err, database := context.DatabaseStore.GetDatabase(statement.TableName.Schema)
	if err != nil {
		return err
	}

	err, table := database.GetTable(statement.TableName.Name)
	if err != nil {
		return err
	}
	row := models.DataRow{
		Values: make([]models.ColumnValue, len(table.Fields)),
	}
	for _, source := range statement.Fields {
		field := table.Fields[source.Field.Name]

		err = field.ValidateValue(source.Value.Data)
		if err != nil {
			return err
		}

		row.Values[field.Index] = source.Value.Data
	}

	go func() {
		table.Datastore.AppendChannel <- row
	}()

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
