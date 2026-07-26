package sql

func IsTokenType(typ TokenType, allowed ...TokenType) bool {
	for _, v := range allowed {
		if v == typ {
			return true
		}
	}

	return false
}

func IsWhereOperation(token *Token) bool {
	if token == nil {
		return false
	}

	return IsTokenType(token.Type, OPERATOR_EQUALS, OPERATOR_GREATER_THAN, OPERATOR_NOT_EQUALITY, OPERATOR_LESS_THAN, IN)
}

func IsWhereCombinatory(token *Token) bool {
	if token == nil {
		return false
	}

	return IsTokenType(token.Type, AND, OR)
}

func IsValueType(token *Token) bool {
	if token == nil {
		return false
	}

	return IsTokenType(token.Type, STRING, NUMBER)
}
