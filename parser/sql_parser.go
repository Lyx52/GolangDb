package parser

import (
	"fmt"

	"github.com/Lyx52/GolangDb/sql"
)

type SqlParser struct {
	statements []*sql.SqlStatement
}

func NewSqlParser() *SqlParser {
	return &SqlParser{
		statements: make([]*sql.SqlStatement, 0),
	}
}

func (sp *SqlParser) Parse(sql string) error {
	lexer := NewLexer(&sql)
	err := lexer.Tokenize()
	if err != nil {
		return err
	}

	statementLexers := lexer.GetLexers()
	for _, statementLexer := range statementLexers {
		err = sp.parseStatement(statementLexer)
		if err != nil {
			return err
		}
	}
	return nil
}

func getCommandType(lexer *BaseLexer) sql.SqlCommandType {
	next := lexer.PopToken()
	if next == nil {
		return sql.INVALID
	}
	switch next.Type {
	case SELECT:
		return sql.SELECT
	case CREATE:
		return sql.CREATE
	default:
		return sql.INVALID
	}
}

func parseCommaSeperatedList(lexer *BaseLexer, allowedTypes ...TokenType) ([]*Token, error) {
	res := make([]*Token, 0)
	continueValues := true
	for continueValues {
		lexer.ConsumeTokens(WHITESPACE)
		next := lexer.PopToken()
		if next == nil || !IsTokenType(next.Type, allowedTypes...) {
			return nil, fmt.Errorf("expected value types")
		}

		res = append(res, next)
		lexer.ConsumeTokens(WHITESPACE)

		next = lexer.PeekToken()
		if next == nil || next.Type != COMMA {
			continueValues = false
		} else {
			lexer.PopToken()
		}
	}
	return res, nil
}

func parseSelectCommand(lexer *BaseLexer, statement *sql.SqlCommand) error {
	fields, err := parseCommaSeperatedList(lexer, IDENTIFIER)
	if err != nil {
		return err
	}
	lexer.ConsumeTokens(WHITESPACE)
	next := lexer.PopToken()
	if next == nil || next.Type != FROM {
		return fmt.Errorf("expected FROM")
	}
	fmt.Println(fields)
	return nil
}

type CommandParserFunc = func(lexer *BaseLexer, statement *sql.SqlCommand) error

var CommandParser = map[sql.SqlCommandType]CommandParserFunc{
	sql.SELECT: parseSelectCommand,
}

func parseCommand(lexer *BaseLexer) (*sql.SqlCommand, error) {
	commandType := getCommandType(lexer)
	if commandType == sql.INVALID {
		return nil, nil
	}
	command := sql.NewSqlCommand(commandType)
	parser, ok := CommandParser[commandType]
	if !ok {
		return nil, fmt.Errorf("unknown command type: %v", commandType)
	}
	err := parser(lexer, command)
	if err != nil {
		return nil, err
	}

	return command, nil
}

func (sp *SqlParser) parseStatement(lexer *BaseLexer) error {
	lexer.ConsumeTokens(WHITESPACE)
	statement := sql.NewSqlStatement()
	command, err := parseCommand(lexer)
	if err != nil {
		return err
	}

	for command != nil {
		statement.Commands = append(statement.Commands, command)
		command, err = parseCommand(lexer)
		if err != nil {
			return err
		}
	}

	if len(statement.Commands) == 0 {
		return fmt.Errorf("statements is empty")
	}

	sp.statements = append(sp.statements, statement)

	return nil
}
