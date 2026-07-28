package parser

import (
	"fmt"

	"github.com/Lyx52/GolangDb/sql"
)

var ClauseParsers map[TokenType]ClauseParseFunc = nil

func InitClauseParsers() {
	if ClauseParsers != nil {
		return
	}

	ClauseParsers = map[TokenType]ClauseParseFunc{
		SELECT:    parseProjectionClause,
		FROM:      parseTableReferenceClause,
		WHERE:     parseSearchConditionClause,
		ORDER_BY:  parseOrderingClause,
		VALUES:    parseValuesClause,
		INSERT:    parseInsertTargetClause,
		RETURNING: parseReturningClause,
		DELETE:    parseDeleteTargetClause,
		TRUNCATE:  parseTruncateTargetClause,
		CREATE:    parseCreateTargetClause,
	}
}

type ClauseParseFunc func(lexer *BaseLexer) (sql.SqlClause, error)

func parseCreateTargetClause(lexer *BaseLexer) (sql.SqlClause, error) {
	lexer.ConsumeTokens(WHITESPACE)
	return nil, nil
}

func parseTruncateTargetClause(lexer *BaseLexer) (sql.SqlClause, error) {
	lexer.ConsumeTokens(WHITESPACE)
	next := lexer.PeekToken()
	if next != nil && next.Type == TABLE {
		lexer.PopToken()
		next = lexer.PeekToken()
	}

	truncatedTables := make([]sql.SqlExpression, 0)
	continueTruncatedTables := true
	for continueTruncatedTables {
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
			continueTruncatedTables = false
		}

		truncatedTables = append(truncatedTables, expr)
		next = lexer.PeekToken()
	}

	cascade := false
	restart := false
	next = lexer.PopToken()
	for next != nil {
		lexer.ConsumeTokens(WHITESPACE)
		switch next.Type {
		case RESTART_IDENTITY:
			restart = true
		case CONTINUE_IDENTITY:
			restart = false
		case CASCADE:
			cascade = true
		case RESTRICT:
			cascade = false
		}

		next = lexer.PopToken()
	}

	return &sql.TruncateTargetClause{
		Truncations: truncatedTables,
		Restart:     restart,
		Cascade:     cascade,
	}, nil
}

func parseReturningClause(lexer *BaseLexer) (sql.SqlClause, error) {
	lexer.ConsumeTokens(WHITESPACE)
	next := lexer.PeekToken()
	returningFields := make([]sql.SqlExpression, 0)
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

		returningFields = append(returningFields, expr)
		next = lexer.PeekToken()
	}

	return &sql.ReturningClause{
		Returning: returningFields,
	}, nil
}

func parseDeleteTargetClause(lexer *BaseLexer) (sql.SqlClause, error) {
	lexer.ConsumeTokens(WHITESPACE)
	next := lexer.PopToken()
	if next == nil || next.Type != FROM {
		return nil, fmt.Errorf("expected FROM but got %v", next)
	}

	tableNameExpr, err := ParseExpression(lexer)
	if err != nil {
		return nil, err
	}

	return &sql.DeleteTargetClause{
		TableName: tableNameExpr,
	}, nil
}

func parseInsertTargetClause(lexer *BaseLexer) (sql.SqlClause, error) {
	lexer.ConsumeTokens(WHITESPACE)
	next := lexer.PopToken()
	if next == nil || next.Type != INTO {
		return nil, fmt.Errorf("expected INTO but got %v", next)
	}

	tableNameExpr, err := ParseExpression(lexer)
	if err != nil {
		return nil, err
	}

	targetList, err := ParseExpressionList(lexer)

	return &sql.InsertTargetClause{
		TableName: tableNameExpr,
		Targets:   targetList,
	}, nil
}

func parseValuesClause(lexer *BaseLexer) (sql.SqlClause, error) {
	lexer.ConsumeTokens(WHITESPACE)
	values, err := ParseExpressionList(lexer)
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
