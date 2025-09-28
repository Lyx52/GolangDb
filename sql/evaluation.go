package sql

import "fmt"

// WhereEvaluation or WhereEvaluationCombination
type WhereStatement interface {
}

type WhereEvaluation struct {
	Field     *Field
	Operation WhereOperation
}

type WhereEvaluationCombination struct {
	Left        WhereStatement
	Right       WhereStatement
	Combinatory WhereCombination
	Depth       int
}

func ParseWhereEvaluation(lexer *BaseLexer, commandStatement Statement) (*WhereEvaluation, error) {
	lexer.ConsumeTokens(WHITESPACE)
	token := lexer.PeekToken()

	var field *Token = nil
	var operation WhereOperation
	var value *Token = nil

	for token != nil && IsTokenType(token.Type, WHITESPACE, OPERATOR_EQUALS, STRING, NUMBER) {
		if token.Type == WHITESPACE {
			lexer.PopToken()
			token = lexer.PeekToken()
			continue
		}

		if field == nil && IsTokenType(token.Type, STRING, NUMBER) {
			field = token
			lexer.PopToken()
			token = lexer.PeekToken()
			continue
		}

		if IsWhereOperation(token) {
			operation = ParseWhereOperation(token)
			lexer.PopToken()
			token = lexer.PeekToken()
			continue
		}

		if value == nil && IsTokenType(token.Type, STRING, NUMBER) {
			value = token
			lexer.PopToken()
			token = lexer.PeekToken()
			continue
		}

		token = lexer.PeekToken()
	}

	if value == nil {
		return nil, fmt.Errorf("expected value")
	}

	if field == nil {
		return nil, fmt.Errorf("expected field name")
	}

	evaluationField, err := ParseFieldName(field, commandStatement)
	if err != nil {
		return nil, err
	}

	evaluationField.Value = value.Value

	return &WhereEvaluation{
		Field:     evaluationField,
		Operation: operation,
	}, nil
}
