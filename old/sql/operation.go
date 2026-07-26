package sql

type WhereOperation int

const (
	EQUALITY_OPERATION WhereOperation = iota
	IN_OPERATOR        WhereOperation = iota
)

func ParseWhereOperation(token *Token) WhereOperation {
	switch token.Type {
	case OPERATOR_EQUALS:
		return EQUALITY_OPERATION
	case IN:
		return IN_OPERATOR
	default:
		return -1
	}
}

func (operation WhereOperation) String() string {
	switch operation {
	case EQUALITY_OPERATION:

		return "="
	case IN_OPERATOR:
		return "IN"
	default:
		return "?"
	}
}
