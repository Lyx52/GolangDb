package sql

type OrderDirection int

func (o OrderDirection) String() string {
	switch o {
	case ORDER_DESC:
		return "DESC"
	case ORDER_ASC:
		return "ASC"
	default:
		return "?"
	}
}

const (
	ORDER_DESC OrderDirection = iota
	ORDER_ASC  OrderDirection = iota
)
