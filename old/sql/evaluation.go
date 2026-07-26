package sql

import "fmt"

// WhereEvaluation or WhereEvaluationCombination
type WhereStatement interface {
	String() string
}

type WhereEvaluation struct {
	Field     *FieldValue
	Operation WhereOperation
}

func (evaluation WhereEvaluation) String() string {
	return fmt.Sprintf("%s %s %s", evaluation.Field.Field.String(), evaluation.Operation.String(), fmt.Sprint(evaluation.Field.Value.String()))
}

type WhereEvaluationCombination struct {
	Left        WhereStatement
	Right       WhereStatement
	Combinatory WhereCombination
	Depth       int
}

func (evaluationCombination WhereEvaluationCombination) String() string {
	if evaluationCombination.Depth > 0 {
		return fmt.Sprintf("(%s %s %s)", evaluationCombination.Left.String(), evaluationCombination.Combinatory.String(), evaluationCombination.Right.String())
	}

	return fmt.Sprintf("%s %s %s", evaluationCombination.Left.String(), evaluationCombination.Combinatory.String(), evaluationCombination.Right.String())
}

func ParseWhereEvaluation(lexer *BaseLexer, commandStatement Statement) (*WhereEvaluation, error) {
	lexer.ConsumeTokens(WHITESPACE)
	token := lexer.PeekToken()

	var field *Token = nil
	var operation WhereOperation
	var value = make([]Token, 0)
	var bracketOpen bool = false
	for token != nil && (IsTokenType(token.Type, WHITESPACE, BRACKET_OPEN, BRACKET_CLOSE, COMMA, STRING, NUMBER) || IsWhereOperation(token)) {
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

		if !bracketOpen && len(value) <= 0 && token.Type == BRACKET_OPEN {
			bracketOpen = true
			value = append(value, *token)
			lexer.PopToken()
			token = lexer.PeekToken()
			continue
		}

		if bracketOpen && token.Type == BRACKET_CLOSE {
			bracketOpen = false
			value = append(value, *token)
			lexer.PopToken()
			token = lexer.PeekToken()
			continue
		} else if token.Type == BRACKET_CLOSE {
			break
		}

		value = append(value, *token)
		lexer.PopToken()
		token = lexer.PeekToken()
	}

	if bracketOpen {
		return nil, fmt.Errorf("expected a closing bracket")
	}

	if len(value) <= 0 {
		return nil, fmt.Errorf("expected value")
	}

	if field == nil {
		return nil, fmt.Errorf("expected field name")
	}

	evaluationField, err := ParseFieldName(field, commandStatement)
	if err != nil {
		return nil, err
	}
	var evaluationFieldValue *Value
	evaluationFieldValue, err = ParseValue(value...)
	if err != nil {
		return nil, err
	}

	return &WhereEvaluation{
		Field:     &FieldValue{Field: evaluationField, Value: evaluationFieldValue},
		Operation: operation,
	}, nil
}
