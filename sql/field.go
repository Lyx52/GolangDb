package sql

import (
	"fmt"
	"strings"
)

type Field struct {
	Name   string
	Alias  string
	Source *TableName
	Value  any
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
