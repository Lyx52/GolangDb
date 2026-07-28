package parser

import (
	"github.com/Lyx52/GolangDb/sql"
)

type SqlParser struct {
	statements []*sql.SqlStatement
}

func NewSqlParser() *SqlParser {
	InitClauseParsers()
	return &SqlParser{
		statements: make([]*sql.SqlStatement, 0),
	}
}

func (sp *SqlParser) PopStatement() *sql.SqlStatement {
	if len(sp.statements) > 0 {
		statement := sp.statements[len(sp.statements)-1]
		sp.statements = sp.statements[:len(sp.statements)-1]
		return statement
	}

	return nil
}

func (sp *SqlParser) Parse(sql string) error {
	lexer := NewLexer(&sql)
	err := lexer.Tokenize()
	if err != nil {
		return err
	}

	statementLexers := lexer.GetLexers()
	for _, statementLexer := range statementLexers {
		statement, err := ParseStatement(statementLexer, true)
		if err != nil {
			return err
		}

		sp.statements = append(sp.statements, statement)
	}
	return nil
}
