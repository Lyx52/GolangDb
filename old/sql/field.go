package sql

import (
	"fmt"
	"strings"
)

type ValueType int

const (
	ARRAY_VALUE  ValueType = iota
	STRING_VALUE ValueType = iota
	NUMBER_VALUE ValueType = iota
)

type Value struct {
	Type ValueType
	Data any
}

func (value Value) String() string {
	switch value.Type {
	case ARRAY_VALUE:
		result := make([]string, 0)
		for _, value := range value.Data.([]any) {
			result = append(result, fmt.Sprint(value))
		}

		return fmt.Sprintf("(%v)", strings.Join(result, ", "))
	case STRING_VALUE:
		return fmt.Sprintf("'%s'", value.Data.(string))
	}

	return fmt.Sprintf("%v", value.Data)
}

type Field struct {
	Name   string
	Alias  string
	Source *TableName
}

func (field *Field) String() string {
	if field.Alias != "" {
		return field.Alias + "." + field.Name
	}

	return field.Name
}

func ParseFieldName(token *Token, commandStatement Statement) (*Field, error) {
	parts := strings.Split(fmt.Sprint(token.Value), ".")
	if len(parts) > 1 {
		alias := parts[0]
		var source = commandStatement.GetBaseTable()
		var err error

		if alias != "" {
			source, err = commandStatement.GetTableName(alias)
			if err != nil {
				return nil, err
			}
		}

		return &Field{
			Name:   parts[1],
			Source: source,
			Alias:  alias,
		}, nil
	}

	return &Field{
		Name:   fmt.Sprint(token.Value),
		Alias:  "",
		Source: commandStatement.GetBaseTable(),
	}, nil
}

type FieldValue struct {
	Field *Field
	Value *Value
}

func (fieldValue *FieldValue) SetFieldValueString() string {
	return fmt.Sprintf("%s = %s", fieldValue.Field.String(), fieldValue.Field.String())
}
