package parser

import "github.com/Lyx52/GolangDb/sql"

var ClauseParsers map[TokenType]ClauseParseFunc = nil

func InitClauseParsers() {
	if ClauseParsers != nil {
		return
	}

	ClauseParsers = map[TokenType]ClauseParseFunc{
		SELECT:   parseProjectionClause,
		FROM:     parseTableReferenceClause,
		WHERE:    parseSearchConditionClause,
		ORDER_BY: parseOrderingClause,
		VALUES:   parseValuesClause,
	}
}

type ClauseParseFunc func(lexer *BaseLexer) (sql.SqlClause, error)

func parseValuesClause(lexer *BaseLexer) (sql.SqlClause, error) {
	lexer.ConsumeTokens(WHITESPACE)
	values, err := parseExpressionList(lexer)
	if err != nil {
		return nil, err
	}
	return &sql.ValuesClause{
		Values: values,
	}, nil
}

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

func parseOrderingClause(lexer *BaseLexer) (sql.SqlClause, error) {
	lexer.ConsumeTokens(WHITESPACE)
	next := lexer.PeekToken()
	orders := make([]sql.SqlExpression, 0)
	continueOrders := true
	for continueOrders {
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
			continueOrders = false
		}

		orders = append(orders, expr)
		next = lexer.PeekToken()
	}

	lexer.ConsumeTokens(WHITESPACE)
	next = lexer.PeekToken()
	direction := sql.ORDER_ASC
	if next != nil && next.Type == DESC {
		lexer.PopToken()
		direction = sql.ORDER_DESC
	}

	if next != nil && next.Type == ASC {
		lexer.PopToken()
		direction = sql.ORDER_ASC
	}

	return &sql.OrderingClause{
		Orders:    orders,
		Direction: direction,
	}, nil
}

func parseSearchConditionClause(lexer *BaseLexer) (sql.SqlClause, error) {
	lexer.ConsumeTokens(WHITESPACE)
	condition, err := ParseExpression(lexer)
	if err != nil {
		return nil, err
	}
	return &sql.SearchConditionClause{
		Condition: condition,
	}, nil
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
