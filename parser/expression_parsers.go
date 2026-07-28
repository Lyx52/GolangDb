package parser

import (
	"fmt"
	"strconv"

	"github.com/Lyx52/GolangDb/sql"
)

func ParseStatement(lexer *BaseLexer, throwOnUnknown bool) (*sql.SqlStatement, error) {
	lexer.ConsumeTokens(WHITESPACE)
	statement := sql.NewSqlStatement()

	next := lexer.PeekToken()
	if next == nil {
		return nil, fmt.Errorf("expected statement")
	}

	switch next.Type {
	case INSERT:
		statement.Type = sql.CREATE
	case SELECT:
		statement.Type = sql.READ
	case UPDATE:
		statement.Type = sql.UPDATE
	case DELETE:
		statement.Type = sql.DELETE
		//default:
		//	return nil, fmt.Errorf("unknown statement type: %v", next.Type)
	}

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
	statement, err := ParseStatement(lexer, false)
	if err != nil {
		return nil, err
	}

	return &sql.QueryExpression{
		Statement: statement,
	}, nil
}

func parseCombinatoryExpression(lexer *BaseLexer, previous sql.SqlExpression) (*sql.CombinatoryExpression, error) {
	operator := lexer.PopToken()
	if operator == nil || !IsOperatorToken(operator) {
		return nil, fmt.Errorf("expected operator token")
	}

	right, err := ParseExpression(lexer)
	if err != nil {
		return nil, err
	}

	return &sql.CombinatoryExpression{
		Left:     previous,
		Right:    right,
		Operator: operator.Type.String(),
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

func parseValueExpressions(lexer *BaseLexer) (sql.SqlExpression, error) {
	lexer.ConsumeTokens(WHITESPACE)
	next := lexer.PeekToken()
	values := make([]any, 0)
	continueValues := true
	for continueValues {
		lexer.ConsumeTokens(WHITESPACE)
		next = lexer.PeekToken()
		if next == nil || !IsTokenType(next.Type, NUMBER, STRING) {
			continueValues = false
			break
		}
		lexer.PopToken()
		if next.Type == NUMBER {
			value, err := strconv.ParseFloat(next.Value, 32)
			if err != nil {
				return nil, err
			}
			values = append(values, value)
		} else {
			values = append(values, next.Value)
		}

		lexer.ConsumeTokens(WHITESPACE)
		next = lexer.PeekToken()
		if next != nil && next.Type == COMMA {
			lexer.PopToken()
		} else {
			continueValues = false
		}
	}

	if len(values) == 0 {
		return nil, fmt.Errorf("expected values")
	}

	if len(values) == 1 {
		return &sql.ValueExpression{
			Value: values[0],
		}, nil
	}

	return &sql.ValueArrayExpression{
		Value: values,
	}, nil
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
		expr, err = parseCombinatoryExpression(lexer, previous)
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

	if IsTokenType(next.Type, NUMBER, STRING) {
		expr, err = parseValueExpressions(lexer)
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
