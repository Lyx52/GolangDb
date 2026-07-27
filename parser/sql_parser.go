package parser

import (
	"fmt"
	"strconv"

	"github.com/Lyx52/GolangDb/sql"
)

var ClauseParsers map[TokenType]ClauseParseFunc

type SqlParser struct {
	statements []*sql.SqlStatement
}

func NewSqlParser() *SqlParser {
	ClauseParsers = map[TokenType]ClauseParseFunc{
		SELECT: parseProjectionClause,
		FROM:   parseTableReferenceClause,
		WHERE:  parseSearchConditionClause,
	}
	return &SqlParser{
		statements: make([]*sql.SqlStatement, 0),
	}
}

func (sp *SqlParser) PopStatement() *sql.SqlStatement {
	if len(sp.statements) > 0 {
		return sp.statements[len(sp.statements)-1]
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
		statement, err := parseStatement(statementLexer, true)
		if err != nil {
			return err
		}

		sp.statements = append(sp.statements, statement)
	}
	return nil
}

type ClauseParseFunc func(lexer *BaseLexer) (sql.SqlClause, error)

func parseProjectionClause(lexer *BaseLexer) (sql.SqlClause, error) {
	lexer.ConsumeTokens(WHITESPACE)
	next := lexer.PeekToken()
	fields := make([]sql.SqlExpression, 0)
	continueFields := true
	for continueFields {
		lexer.ConsumeTokens(WHITESPACE)
		expr, err := ParseExpression(lexer)
		if err != nil {
			return nil, err
		}
		lexer.ConsumeTokens(WHITESPACE)
		next = lexer.PeekToken()
		if next != nil && next.Type == COMMA {
			lexer.PopToken()
		} else {
			continueFields = false
		}

		fields = append(fields, expr)
		next = lexer.PeekToken()
	}

	return &sql.ProjectionClause{Fields: fields}, nil
}

func parseSearchConditionClause(lexer *BaseLexer) (sql.SqlClause, error) {

	return &sql.SearchConditionClause{}, nil
}

func parseTableReferenceClause(lexer *BaseLexer) (sql.SqlClause, error) {
	lexer.ConsumeTokens(WHITESPACE)
	next := lexer.PeekToken()
	tableRefs := make([]sql.SqlExpression, 0)
	continueTableRefs := true
	for continueTableRefs {
		lexer.ConsumeTokens(WHITESPACE)
		expr, err := ParseExpression(lexer)
		if err != nil {
			return nil, err
		}
		lexer.ConsumeTokens(WHITESPACE)
		next = lexer.PeekToken()
		if next != nil && next.Type == COMMA {
			lexer.PopToken()
		} else {
			continueTableRefs = false
		}
		next = lexer.PeekToken()
		tableRefs = append(tableRefs, expr)
	}

	return &sql.TableReferenceClause{Tables: tableRefs}, nil
}

func parseStatement(lexer *BaseLexer, throwOnUnknown bool) (*sql.SqlStatement, error) {
	lexer.ConsumeTokens(WHITESPACE)
	statement := sql.NewSqlStatement()

	next := lexer.PeekToken()
	for next != nil {
		lexer.ConsumeTokens(WHITESPACE)
		parser, ok := ClauseParsers[next.Type]
		if !ok && throwOnUnknown {
			return nil, fmt.Errorf("unexpected clause %v", next.Type)
		} else if !ok {
			break
		}

		lexer.PopToken()
		clause, err := parser(lexer)
		if err != nil {
			return nil, err
		}
		statement.Clauses = append(statement.Clauses, clause)
		lexer.ConsumeTokens(WHITESPACE)
		next = lexer.PeekToken()
	}

	return statement, nil
}

func parseSubExpression(lexer *BaseLexer) (*sql.SubExpression, error) {
	lexer.ConsumeTokens(WHITESPACE)
	next := lexer.PopToken()
	if next == nil || next.Type != BRACKET_OPEN {
		return nil, fmt.Errorf("expected bracket open")
	}

	expr, err := ParseExpression(lexer)
	if err != nil {
		return nil, err
	}

	lexer.ConsumeTokens(WHITESPACE)
	next = lexer.PopToken()
	if next == nil || next.Type != BRACKET_CLOSE {
		return nil, fmt.Errorf("expected bracket close")
	}

	return &sql.SubExpression{
		Expression: expr,
	}, nil
}

func parseFunctionExpression(lexer *BaseLexer) (*sql.FunctionExpression, error) {
	lexer.ConsumeTokens(WHITESPACE)
	functionToken := lexer.PopToken()
	if functionToken == nil || !IsFunctionToken(functionToken) {
		return nil, fmt.Errorf("expected function token")
	}

	next := lexer.PeekToken()
	var parameters []sql.SqlExpression
	var err error
	if next != nil && next.Type == BRACKET_OPEN {
		parameters, err = parseExpressionList(lexer)
		if err != nil {
			return nil, err
		}
	} else {
		parameters = make([]sql.SqlExpression, 0)
	}

	return &sql.FunctionExpression{
		FunctionName: functionToken.Type.String(),
		Parameters:   parameters,
	}, nil
}

func parseQueryExpression(lexer *BaseLexer) (*sql.QueryExpression, error) {
	lexer.ConsumeTokens(WHITESPACE)
	statement, err := parseStatement(lexer, false)
	if err != nil {
		return nil, err
	}

	return &sql.QueryExpression{
		Statement: statement,
	}, nil
}

func parseExpressionList(lexer *BaseLexer) ([]sql.SqlExpression, error) {

	lexer.ConsumeTokens(WHITESPACE)
	next := lexer.PopToken()
	if next == nil || next.Type != BRACKET_OPEN {
		return nil, fmt.Errorf("expected bracket open")
	}
	expressions := make([]sql.SqlExpression, 0)
	continueExpressions := true
	for continueExpressions {
		lexer.ConsumeTokens(WHITESPACE)
		expr, err := ParseExpression(lexer)
		if err != nil {
			return nil, err
		}

		lexer.ConsumeTokens(WHITESPACE)
		next = lexer.PeekToken()
		if next == nil || next.Type == BRACKET_CLOSE {
			continueExpressions = false
		}

		expressions = append(expressions, expr)
	}

	lexer.ConsumeTokens(WHITESPACE)
	next = lexer.PopToken()
	if next == nil || next.Type != BRACKET_CLOSE {
		return nil, fmt.Errorf("expected bracket close")
	}

	return expressions, nil
}

func parseFieldExpression(lexer *BaseLexer, previous sql.SqlExpression) (*sql.FieldExpression, error) {
	lexer.ConsumeTokens(WHITESPACE)
	next := lexer.PopToken()
	if next == nil || next.Type != DOT {
		return nil, fmt.Errorf("expected dot")
	}

	lexer.ConsumeTokens(WHITESPACE)
	fieldName := lexer.PopToken()
	if fieldName == nil {
		return nil, fmt.Errorf("expected field name")
	}

	return &sql.FieldExpression{
		FieldName:  fieldName.Value,
		Expression: previous,
	}, nil
}

func parseSubscriptExpression(lexer *BaseLexer, previous sql.SqlExpression) (sql.SqlExpression, error) {
	lexer.ConsumeTokens(WHITESPACE)
	next := lexer.PopToken()
	if next == nil || next.Type != SQUARE_BRACKET_OPEN {
		return nil, fmt.Errorf("expected open colon")
	}

	lexer.ConsumeTokens(WHITESPACE)
	subscript := lexer.PopToken()
	if subscript == nil || subscript.Type != NUMBER {
		return nil, fmt.Errorf("expected subscript")
	}
	from, err := strconv.Atoi(subscript.Value)
	if err != nil {
		return nil, err
	}

	lexer.ConsumeTokens(WHITESPACE)
	next = lexer.PopToken()
	if next == nil || next.Type != COLON {
		next = lexer.PopToken()
		if next == nil || next.Type != SQUARE_BRACKET_CLOSE {
			return nil, fmt.Errorf("expected close colon")
		}

		return &sql.SubscriptExpression{
			Subscript:  from,
			Expression: previous,
		}, nil
	}

	upperSubscript := lexer.PopToken()
	if upperSubscript == nil || upperSubscript.Type != NUMBER {
		return nil, fmt.Errorf("expected upper subscript")
	}
	till, err := strconv.Atoi(upperSubscript.Value)
	if err != nil {
		return nil, err
	}

	next = lexer.PopToken()
	if next == nil || next.Type != SQUARE_BRACKET_CLOSE {
		return nil, fmt.Errorf("expected close colon")
	}

	return &sql.SubscriptSliceExpression{
		LowerSubscript: from,
		UpperSubscript: till,
		Expression:     previous,
	}, nil
}

func parseAliasExpression(lexer *BaseLexer, previous sql.SqlExpression) (sql.SqlExpression, error) {
	if previous == nil {
		return nil, fmt.Errorf("expected alias expression")
	}

	lexer.ConsumeTokens(WHITESPACE)
	next := lexer.PopToken()
	if next == nil || next.Type != AS {
		return nil, fmt.Errorf("expected AS")
	}
	lexer.ConsumeTokens(WHITESPACE)
	alias := lexer.PopToken()
	if alias == nil {
		return nil, fmt.Errorf("expected identifier or quoted identifier or typename")
	}

	if IsTokenType(alias.Type, QUOTED_IDENTIFIER, IDENTIFIER) {
		return &sql.AliasExpression{
			Alias:      alias.Value,
			Quoted:     alias.Type == QUOTED_IDENTIFIER,
			Expression: previous,
		}, nil
	} else if IsTokenType(alias.Type, TYPE_NAME) {
		return &sql.TypeAliasExpression{
			TypeName:   alias.Value,
			Expression: previous,
		}, nil
	}
	return nil, fmt.Errorf("expected identifier or quoted identifier or typename")
}

func parseIdentifierExpression(lexer *BaseLexer) (sql.SqlExpression, error) {
	lexer.ConsumeTokens(WHITESPACE)
	first := lexer.PopToken()
	if first == nil || first.Type != IDENTIFIER {
		return nil, fmt.Errorf("expected identifier")
	}
	lexer.ConsumeTokens(WHITESPACE)
	next := lexer.PeekToken()

	// correlation.columname
	if next != nil && next.Type == DOT {
		lexer.PopToken()
		second := lexer.PopToken()

		if second != nil && second.Type == WILDCARD {
			return &sql.WildcardExpression{
				Correlation: first.Value,
			}, nil
		} else if second != nil && second.Type == IDENTIFIER {
			return &sql.ReferenceExpression{
				Correlation: first.Value,
				Name:        second.Value,
			}, nil
		}
		return nil, fmt.Errorf("expected identifier after dot")
	}

	return &sql.ReferenceExpression{
		Name: first.Value,
	}, nil
}

func parseSubExpressions(lexer *BaseLexer, previous sql.SqlExpression) (sql.SqlExpression, error) {
	lexer.ConsumeTokens(WHITESPACE)
	next := lexer.PeekToken()
	var err error
	var expr sql.SqlExpression

	if next != nil && next.Type == AS {
		expr, err = parseAliasExpression(lexer, previous)
	}

	if next != nil && next.Type == DOT {
		expr, err = parseFieldExpression(lexer, previous)
	}

	if next != nil && next.Type == SQUARE_BRACKET_OPEN {
		expr, err = parseSubscriptExpression(lexer, previous)
	}

	if next != nil && IsOperatorToken(next) {
		expr, err = parse(lexer, previous)
	}

	if err != nil {
		return nil, err
	}

	if expr == nil {
		return previous, nil
	}

	expr, err = parseSubExpressions(lexer, expr)
	if err != nil {
		return nil, err
	}

	return expr, nil
}

func ParseExpression(lexer *BaseLexer) (sql.SqlExpression, error) {
	lexer.ConsumeTokens(WHITESPACE)
	next := lexer.PeekToken()
	if next == nil {
		return nil, fmt.Errorf("expected expression")
	}

	var err error
	var expr sql.SqlExpression
	if IsFunctionToken(next) {
		expr, err = parseFunctionExpression(lexer)
	}

	if next.Type == WILDCARD {
		lexer.PopToken()
		expr = &sql.WildcardExpression{}
	}

	if next.Type == BRACKET_OPEN {
		expr, err = parseSubExpression(lexer)
	}

	if next.Type == IDENTIFIER {
		expr, err = parseIdentifierExpression(lexer)
	}

	// Sub query
	if next.Type == SELECT {
		expr, err = parseQueryExpression(lexer)
	}

	if err != nil {
		return nil, err
	}

	expr, err = parseSubExpressions(lexer, expr)
	if err != nil {
		return nil, err
	}

	return expr, nil
}
