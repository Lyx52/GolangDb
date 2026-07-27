package sql

import (
	"fmt"
	"strings"
)

type SqlClause interface {
	String() string
}

type ProjectionClause struct {
	Fields []SqlExpression
}

func (p *ProjectionClause) String() string {
	items := make([]string, len(p.Fields))
	for i, field := range p.Fields {
		items[i] = field.String()
	}

	return fmt.Sprintf("SELECT %s", strings.Join(items, ", "))
}

type TableReferenceClause struct {
	Tables []SqlExpression
}

func (p *TableReferenceClause) String() string {
	items := make([]string, len(p.Tables))
	for i, table := range p.Tables {
		items[i] = table.String()
	}

	return fmt.Sprintf("FROM %s", strings.Join(items, ", "))
}

type SearchConditionClause struct {
	Condition SqlExpression
}

func (p *SearchConditionClause) String() string {
	return fmt.Sprintf("WHERE %s", p.Condition)
}
