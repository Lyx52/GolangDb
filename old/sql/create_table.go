package sql

import (
	"fmt"
	"strings"

	"github.com/Lyx52/GolangDb/schema"
	"github.com/Lyx52/GolangDb/server"
)

type FieldDefinition struct {
	Name string
	Type schema.FieldType
}

func (definition *FieldDefinition) String() string {
	return fmt.Sprintf("%s %s", definition.Name, definition.Type.String())
}

type CreateTableStatement struct {
	Name   string
	Fields []FieldDefinition
}

func NewCreateTableStatement(name string) *CreateTableStatement {
	return &CreateTableStatement{
		Name:   name,
		Fields: make([]FieldDefinition, 0),
	}
}

func (statement *CreateTableStatement) AddField(name string, fieldType string) error {
	err, parsedType := schema.TypeFromString(fieldType)
	if err != nil {
		return err
	}

	statement.Fields = append(statement.Fields, FieldDefinition{
		Name: name,
		Type: parsedType,
	})

	return nil
}

func (statement *CreateTableStatement) String() string {
	fields := make([]string, len(statement.Fields))
	for i, field := range statement.Fields {
		fields[i] = field.String()
	}

	return fmt.Sprintf("CREATE TABLE %s (%s)", statement.Name, strings.Join(fields, ", "))
}

func (statement *CreateTableStatement) Execute(context *server.ServerContext) error {
	err, table := context.CreateTable(statement.Name)
	if err != nil {
		return err
	}

	for _, field := range statement.Fields {
		err = table.CreateField(field.Name, field.Type, make([]*schema.Constraint, 0))
		if err != nil {
			return err
		}
	}

	fmt.Printf("CREATED TABLE %v\n", table.Name)
	return err
}
