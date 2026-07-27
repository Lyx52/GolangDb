package sql

import "strings"

type SqlStatement struct {
	Clauses []SqlClause
}

func NewSqlStatement() *SqlStatement {
	return &SqlStatement{
		Clauses: make([]SqlClause, 0),
	}
}

func (statement *SqlStatement) String() string {
	items := make([]string, len(statement.Clauses))
	for i, clause := range statement.Clauses {
		items[i] = clause.String()
	}

	return strings.Join(items, " ")
}
