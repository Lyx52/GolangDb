package sql

type WhereCombination int

const (
	AND_COMBINATION WhereCombination = iota
	OR_COMBINATION  WhereCombination = iota
)

func ParseWhereCombination(token *Token) WhereCombination {
	switch token.Type {
	case AND:
		return AND_COMBINATION
	case OR:
		return OR_COMBINATION
	default:
		return -1
	}
}
