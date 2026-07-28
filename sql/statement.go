package sql

import "strings"

type SqlStatementType int

const (
	UNKNOWN SqlStatementType = iota
	CREATE  SqlStatementType = iota
	READ    SqlStatementType = iota
	UPDATE  SqlStatementType = iota
	DELETE  SqlStatementType = iota
)

type SqlStatement struct {
	Clauses []SqlClause
	Type    SqlStatementType
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
