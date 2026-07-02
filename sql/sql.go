package sql

import (
	"fmt"
	"os"
	"strings"

	"github.com/Lyx52/GolangDb/backing"
)

type Statement interface {
	GetTableName(alias string) (*TableName, error)
	PushTableName(alias string, tableName *TableName) error
	GetBaseTable() *TableName
	Execute(context *backing.ServerContext) error
}

type SqlParser struct {
	statements []Statement
}

func (parser *SqlParser) PushStatement(statement Statement) {
	parser.statements = append(parser.statements, statement)
}

func (parser *SqlParser) PopStatement() *Statement {
	if len(parser.statements) == 0 {
		return nil
	}

	res := parser.statements[0]
	parser.statements = parser.statements[1:]
	return &res
}

type TableName struct {
	Name   string
	Alias  string
	Schema string
}

func (tableName TableName) String() string {
	name := tableName.Name
	if tableName.Schema != "" {
		name = tableName.Schema + "." + name
	}

	if tableName.Alias != "" {
		name = name + " AS " + tableName.Alias
	}

	return name
}

func NewSqlParser() *SqlParser {
	return &SqlParser{}
}

func RaiseIfError(err error) {
	if err != nil {
		fmt.Printf("SqlParser Error: %s\n", err.Error())
		os.Exit(1)
	}
}

func ParseList(lexer *BaseLexer, allowedTypes ...TokenType) ([][]Token, error) {
	token := lexer.PeekToken()
	res := make([][]Token, 0)
	current := make([]Token, 0)
	bracketDepth := 0
	for token != nil && IsTokenType(token.Type, allowedTypes...) {
		if token.Type == WHITESPACE {
			lexer.PopToken()
			token = lexer.PeekToken()
			continue
		}

		if token.Type == BRACKET_OPEN {
			bracketDepth++
			current = append(current, *token)
			lexer.PopToken()
			token = lexer.PeekToken()
			continue
		}

		if token.Type == BRACKET_CLOSE {
			bracketDepth--
			current = append(current, *token)
			lexer.PopToken()
			token = lexer.PeekToken()
			continue
		}

		if token.Type == COMMA && bracketDepth == 0 {
			lexer.PopToken()
			token = lexer.PeekToken()
			res = append(res, current)
			current = make([]Token, 0)
			continue
		}

		current = append(current, *token)
		lexer.PopToken()
		token = lexer.PeekToken()
	}

	if len(current) > 0 {
		res = append(res, current)
	}

	if bracketDepth > 0 {
		return nil, fmt.Errorf("uneven opening or closing brackets")
	}

	return res, nil
}

func ParseValueList(lexer *BaseLexer, mandatoryComma bool, allowNumbers bool) ([]*Token, error) {
	token := lexer.PeekToken()
	res := make([]*Token, 0)
	expectComma := false
	expectSpace := false
	for token != nil && (token.Type == STRING || token.Type == COMMA || token.Type == WHITESPACE || (allowNumbers && token.Type == NUMBER)) {
		if token.Type == WHITESPACE {
			lexer.PopToken()
			token = lexer.PeekToken()
			expectSpace = false
			continue
		} else if expectSpace {
			return res, fmt.Errorf("expected space position at %v", token.Position)
		}

		if expectComma && token.Type != COMMA {
			return res, fmt.Errorf("expected comma position at %v", token.Position)
		} else if token.Type == COMMA {
			lexer.PopToken()
			expectComma = false
		}

		if token.Type == STRING || token.Type == NUMBER {
			lexer.PopToken()
			res = append(res, token)
			token = lexer.PeekToken()
			expectComma = mandatoryComma
			expectSpace = !mandatoryComma
			continue
		}

		token = lexer.PeekToken()
	}

	return res, nil
}

func ParseTableName(lexer *BaseLexer, allowAlias bool) (*TableName, error) {
	token := lexer.PopToken()
	if token == nil || token.Type != STRING {
		return nil, fmt.Errorf("expected table_name string")
	}
	tableName := fmt.Sprint(token.Value)
	schema := ""
	alias := ""
	parts := strings.Split(tableName, ".")

	if len(parts) > 1 {
		schema = parts[0]
		tableName = parts[1]
	}

	lexer.ConsumeTokens(WHITESPACE)
	token = lexer.PeekToken()
	if token != nil && token.Type == AS {
		if !allowAlias {
			return nil, fmt.Errorf("cannot use alias here")
		}
		lexer.PopToken()
		lexer.ConsumeTokens(WHITESPACE)
		token = lexer.PopToken()
		if token == nil || token.Type != STRING {
			return nil, fmt.Errorf("expected table_name alias string")
		}
		alias = fmt.Sprint(token.Value)
	}

	return &TableName{
		Name:   tableName,
		Alias:  alias,
		Schema: schema,
	}, nil
}

func (parser *SqlParser) ParseInsert(lexer *BaseLexer) error {
	statement := NewInsertStatement()
	// 'SELECT INTO table_name'
	if lexer.ConsumeTokens(WHITESPACE) <= 0 {
		return fmt.Errorf("expected WHITESPACE after SELECT")
	}

	err := lexer.ConsumeExpectToken(INTO)
	if err != nil {
		return err
	}

	if lexer.ConsumeTokens(WHITESPACE) <= 0 {
		return fmt.Errorf("expected WHITESPACE after INTO")
	}
	tableName, err := ParseTableName(lexer, true)
	if err != nil {
		return err
	}

	if tableName.Alias != "" {
		err = statement.PushTableName(tableName.Alias, tableName)
		if err != nil {
			return err
		}
	}

	statement.TableName = tableName

	// ' (field1, field2, ...) '
	lexer.ConsumeTokens(WHITESPACE)
	err = lexer.ConsumeExpectToken(BRACKET_OPEN)
	if err != nil {
		return err
	}

	fields, err := ParseValueList(lexer, true, false)
	if err != nil {
		return err
	}

	err = lexer.ConsumeExpectToken(BRACKET_CLOSE)
	if err != nil {
		return err
	}

	// ' VALUES '
	lexer.ConsumeTokens(WHITESPACE)
	err = lexer.ConsumeExpectToken(VALUES)
	if err != nil {
		return err
	}

	// ' (value1 value2 value) '
	lexer.ConsumeTokens(WHITESPACE)
	err = lexer.ConsumeExpectToken(BRACKET_OPEN)
	if err != nil {
		return err
	}

	values, err := ParseValueList(lexer, false, true)
	if err != nil {
		return err
	}

	if len(fields) != len(values) {
		return fmt.Errorf("expected equal fields for values, got %d fields and %d values", len(fields), len(values))
	}

	for i, field := range fields {
		var statementField *Field
		statementField, err = ParseFieldName(field, statement)
		if err != nil {
			return err
		}

		statement.Fields = append(statement.Fields, FieldValue{
			Field: statementField,
			Value: &Value{
				Type: ParseTokenValueType(*values[i]),
				Data: values[i].Value,
			},
		})
	}

	lexer.ConsumeTokens(WHITESPACE)
	err = lexer.ConsumeExpectToken(BRACKET_CLOSE)
	if err != nil {
		return err
	}

	lexer.ConsumeTokens(WHITESPACE)
	lexer.ConsumeTokens(SEMICOLUMN)

	parser.PushStatement(statement)
	return nil
}

func ParseWhereEvaluationOrStatement(lexer *BaseLexer, depth int, commandStatement Statement) (WhereStatement, error) {
	var statement WhereStatement
	var err error
	lexer.ConsumeTokens(WHITESPACE)
	token := lexer.PeekToken()

	if token != nil && token.Type == BRACKET_OPEN {
		lexer.PopToken()
		statement, err = ParseWhereStatement(lexer, nil, depth+1, commandStatement)
		if err != nil {
			return nil, err
		}
		token = lexer.PopToken()
		if token == nil || token.Type != BRACKET_CLOSE {
			return nil, fmt.Errorf("expected BRACKET_CLOSE")
		}
	} else {
		evaluation, err := ParseWhereEvaluation(lexer, commandStatement)
		if err != nil {
			return nil, err
		}

		statement = *evaluation
	}

	return statement, nil
}

func ParseWhereStatement(lexer *BaseLexer, left WhereStatement, depth int, commandStatement Statement) (WhereStatement, error) {
	var statement WhereStatement
	var right WhereStatement
	var err error
	lexer.ConsumeTokens(WHITESPACE)
	statement = left

	if statement == nil {
		statement, err = ParseWhereEvaluationOrStatement(lexer, depth, commandStatement)
		if err != nil {
			return nil, err
		}
	}

	lexer.ConsumeTokens(WHITESPACE)
	token := lexer.PeekToken()
	if IsWhereCombinatory(token) {
		lexer.PopToken()
		combination := ParseWhereCombination(token)
		right, err = ParseWhereEvaluationOrStatement(lexer, depth, commandStatement)
		if err != nil {
			return nil, err
		}

		statement = WhereEvaluationCombination{
			Left:        statement,
			Right:       right,
			Combinatory: combination,
			Depth:       depth,
		}
	}

	lexer.ConsumeTokens(WHITESPACE)
	token = lexer.PeekToken()
	if IsWhereCombinatory(token) {
		return ParseWhereStatement(lexer, statement, depth, commandStatement)
	}

	return statement, nil
}

func ParseTokenValueType(token Token) ValueType {
	switch token.Type {
	case NUMBER:
		return NUMBER_VALUE
	case STRING:
		return STRING_VALUE
	default:
		return ARRAY_VALUE
	}
}

func ParseValue(tokens ...Token) (*Value, error) {
	if len(tokens) < 1 {
		return nil, fmt.Errorf("expected at least one token")
	}
	var result []any = nil
	var commas int = 0
	for _, token := range tokens {
		if token.Type == BRACKET_OPEN {
			result = make([]any, 0)
			continue
		}

		if token.Type == BRACKET_CLOSE {
			if (len(result) - 1) != commas {
				return nil, fmt.Errorf("expected commas between values")
			}
			return &Value{
				Type: ARRAY_VALUE,
				Data: result,
			}, nil
		}

		if token.Type == COMMA {
			commas++
			continue
		}

		// Single value return early
		if result == nil {
			if token.Type == STRING {
				value := fmt.Sprint(token.Value)
				value = value[1 : len(value)-1]
				return &Value{
					Type: ParseTokenValueType(token),
					Data: value,
				}, nil
			}

			return &Value{
				Type: ParseTokenValueType(token),
				Data: token.Value,
			}, nil
		}

		// Array path append value
		if token.Type == STRING {
			stringValue := fmt.Sprint(token.Value)
			stringValue = stringValue[1 : len(stringValue)-1]
			result = append(result, stringValue)
			continue
		}

		result = append(result, token.Value)
	}

	return nil, fmt.Errorf("array with uneven brackets")
}

func (parser *SqlParser) ParseUpdate(lexer *BaseLexer) error {
	statement := NewUpdateStatement()

	if lexer.ConsumeTokens(WHITESPACE) <= 0 {
		return fmt.Errorf("expected WHITESPACE after UPDATE")
	}

	tableName, err := ParseTableName(lexer, true)
	if err != nil {
		return err
	}

	if tableName.Alias != "" {
		err = statement.PushTableName(tableName.Alias, tableName)
		if err != nil {
			return err
		}
	}

	statement.TableName = tableName
	lexer.ConsumeTokens(WHITESPACE)

	err = lexer.ConsumeExpectToken(SET)
	if err != nil {
		return err
	}
	lexer.ConsumeTokens(WHITESPACE)
	sets, err := ParseList(lexer, WHITESPACE, COMMA, OPERATOR_EQUALS, BRACKET_OPEN, BRACKET_CLOSE, STRING, NUMBER)
	for _, set := range sets {
		if len(set) < 3 {
			return fmt.Errorf("expected set to have 3 tokens, got %d", len(set))
		}
		var statementField *Field
		statementField, err = ParseFieldName(&set[0], statement)

		if err != nil {
			return err
		}

		var fieldValue *Value
		fieldValue, err = ParseValue(set[2:]...)
		if err != nil {
			return err
		}

		statement.Fields = append(statement.Fields, FieldValue{
			Field: statementField,
			Value: fieldValue,
		})
	}
	lexer.ConsumeTokens(WHITESPACE)
	token := lexer.PopToken()
	if token != nil && token.Type == WHERE {
		whereStatement, err := ParseWhereStatement(lexer, nil, 0, statement)
		if err != nil {
			return err
		}
		statement.Where = whereStatement
	}
	lexer.ConsumeTokens(WHITESPACE)
	lexer.ConsumeTokens(SEMICOLUMN)

	parser.PushStatement(statement)
	return nil
}

func (parser *SqlParser) ParseSelect(lexer *BaseLexer) error {
	if lexer.ConsumeTokens(WHITESPACE) <= 0 {
		return fmt.Errorf("expected WHITESPACE after UPDATE")
	}
	return nil
}

func (parser *SqlParser) ParseDelete(lexer *BaseLexer) error {
	if lexer.ConsumeTokens(WHITESPACE) <= 0 {
		return fmt.Errorf("expected WHITESPACE after DELETE")
	}
	return nil
}

func (parser *SqlParser) ParseCreate(lexer *BaseLexer) error {
	if lexer.ConsumeTokens(WHITESPACE) <= 0 {
		return fmt.Errorf("expected WHITESPACE after CREATE")
	}
	token := lexer.PopToken()
	if token == nil {
		return fmt.Errorf("expected object type token after CREATE")
	}

	var objectType CreateObjectType
	var err error
	var subStatement CreateObjectStatement
	switch token.Type {
	case DATABASE:
		objectType = CREATE_DATABASE
		err, subStatement = parser.ParseCreateDatabase(lexer)
	case TABLE:
		objectType = CREATE_TABLE
		err, subStatement = parser.ParseCreateTable(lexer)
	case VIEW:
		objectType = CREATE_VIEW
		err, subStatement = parser.ParseCreateView(lexer)
	default:
		return fmt.Errorf("expected CREATE object type token")
	}

	if err != nil {
		return err
	}

	statement := NewCreateStatement(objectType)
	statement.SubStatement = &subStatement
	parser.PushStatement(statement)
	return nil
}

func (parser *SqlParser) ParseCreateDatabase(lexer *BaseLexer) (error, *CreateDatabaseStatement) {
	if lexer.ConsumeTokens(WHITESPACE) <= 0 {
		return fmt.Errorf("expected WHITESPACE after CREATE DATABASE"), nil
	}
	databaseNameToken := lexer.PopToken()
	if databaseNameToken == nil || databaseNameToken.Type != STRING {
		return fmt.Errorf("expected database name token"), nil
	}

	statement := NewCreateDatabaseStatement(databaseNameToken.StringValue())

	return nil, statement
}

func (parser *SqlParser) ParseCreateTable(lexer *BaseLexer) (error, *CreateTableStatement) {
	if lexer.ConsumeTokens(WHITESPACE) <= 0 {
		return fmt.Errorf("expected WHITESPACE after CREATE TABLE"), nil
	}
	tableNameToken := lexer.PopToken()
	if tableNameToken == nil || tableNameToken.Type != STRING {
		return fmt.Errorf("expected table name token"), nil
	}

	statement := NewCreateTableStatement(tableNameToken.StringValue())

	return nil, statement
}

func (parser *SqlParser) ParseCreateView(lexer *BaseLexer) (error, *CreateViewStatement) {
	if lexer.ConsumeTokens(WHITESPACE) <= 0 {
		return fmt.Errorf("expected WHITESPACE after CREATE VIEW"), nil
	}
	viewNameToken := lexer.PopToken()
	if viewNameToken == nil || viewNameToken.Type != STRING {
		return fmt.Errorf("expected view name token"), nil
	}

	statement := NewCreateViewStatement(viewNameToken.StringValue())

	return nil, statement
}

func (parser *SqlParser) ParseUse(lexer *BaseLexer) error {
	if lexer.ConsumeTokens(WHITESPACE) <= 0 {
		return fmt.Errorf("expected WHITESPACE after USE")
	}

	useDatabaseNameToken := lexer.PopToken()
	if useDatabaseNameToken == nil || useDatabaseNameToken.Type != STRING {
		return fmt.Errorf("expected database name token")
	}

	statement := NewUseStatement(useDatabaseNameToken.StringValue())
	parser.PushStatement(statement)

	return nil
}

func (parser *SqlParser) ParseShow(lexer *BaseLexer) error {
	if lexer.ConsumeTokens(WHITESPACE) <= 0 {
		return fmt.Errorf("expected WHITESPACE after SHOW")
	}
	objectType := lexer.PopToken()
	if objectType == nil {
		return fmt.Errorf("expected object type token")
	}
	var showObjectType ShowObjectType
	switch objectType.Type {
	case DATABASES:
		showObjectType = SHOW_DATABASES
	case TABLES:
		showObjectType = SHOW_TABLES
	default:
		return fmt.Errorf("expected SHOW object type token")
	}
	statement := NewShowStatement(showObjectType)
	parser.PushStatement(statement)

	return nil
}

func (parser *SqlParser) Parse(sql string) error {
	lexer := NewLexer(&sql)
	err := lexer.Tokenize()
	if err != nil {
		return err
	}

	lexer.ConsumeTokens(WHITESPACE)
	token := lexer.PopToken()
	if token == nil {
		return fmt.Errorf("SqlParser Error: empty string")
	}

	switch token.Type {
	case DELETE:
		return parser.ParseDelete(lexer)
	case INSERT:
		return parser.ParseInsert(lexer)
	case SELECT:
		return parser.ParseSelect(lexer)
	case UPDATE:
		return parser.ParseUpdate(lexer)
	case CREATE:
		return parser.ParseCreate(lexer)
	case USE:
		return parser.ParseUse(lexer)
	case SHOW:
		return parser.ParseShow(lexer)
	default:
		return fmt.Errorf("SqlParser: Unexpected token type: %s", token.Type)
	}
}
