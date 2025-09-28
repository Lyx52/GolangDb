package sql

func IsTokenType(typ TokenType, allowed ...TokenType) bool {
	for _, v := range allowed {
		if v == typ {
			return true
		}
	}

	return false
}
