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

type OrderingClause struct {
	Orders    []SqlExpression
	Direction OrderDirection
}

func (p *OrderingClause) String() string {
	items := make([]string, len(p.Orders))
	for i, order := range p.Orders {
		items[i] = order.String()
	}

	return fmt.Sprintf("ORDER BY %s %s", strings.Join(items, ", "), p.Direction)
}

type ValuesClause struct {
	Values []SqlExpression
}

func (p *ValuesClause) String() string {
	items := make([]string, len(p.Values))
	for i, value := range p.Values {
		items[i] = value.String()
	}

	return fmt.Sprintf("VALUES (%s)", strings.Join(items, ", "))
}

type InsertTargetClause struct {
	TableName SqlExpression
	Targets   []SqlExpression
}

func (p *InsertTargetClause) String() string {
	items := make([]string, len(p.Targets))
	for i, target := range p.Targets {
		items[i] = target.String()
	}

	return fmt.Sprintf("INSERT INTO %s (%s)", p.TableName, strings.Join(items, ", "))
}

type ReturningClause struct {
	Returning []SqlExpression
}

func (p *ReturningClause) String() string {
	items := make([]string, len(p.Returning))
	for i, returnValue := range p.Returning {
		items[i] = returnValue.String()
	}

	return fmt.Sprintf("RETURNING %s", strings.Join(items, ", "))
}

type DeleteTargetClause struct {
	TableName SqlExpression
}

func (p *DeleteTargetClause) String() string {
	return fmt.Sprintf("DELETE FROM %s", p.TableName)
}

type TruncateTargetClause struct {
	Truncations []SqlExpression
	Cascade     bool
	Restart     bool
}

func (p *TruncateTargetClause) String() string {
	items := make([]string, len(p.Truncations))
	for i, trunc := range p.Truncations {
		items[i] = trunc.String()
	}
	parameters := make([]string, 0)
	if p.Restart {
		parameters = append(parameters, "RESTART IDENTITY")
	}

	if p.Cascade {
		parameters = append(parameters, "CASCADE")
	}

	if len(parameters) > 0 {
		return fmt.Sprintf("TRUNCATE TABLE %s %s", strings.Join(items, ", "), strings.Join(parameters, " "))
	}

	return fmt.Sprintf("TRUNCATE TABLE %s", strings.Join(items, ", "))
}
