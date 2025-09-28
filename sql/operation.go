package sql

type WhereOperation int

const (
	EQUALITY_OPERATION WhereOperation = iota
)

func ParseWhereOperation(token *Token) WhereOperation {
	switch token.Type {
	case OPERATOR_EQUALS:
		return EQUALITY_OPERATION
	default:
		return -1
	}
}

func (operation WhereOperation) String() string {
	switch operation {
	case EQUALITY_OPERATION:
		return "="
	default:
		return "?"
	}
}
