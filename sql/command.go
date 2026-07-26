package sql

type SqlCommandType int

const (
	INVALID SqlCommandType = iota
	SELECT  SqlCommandType = iota
	CREATE  SqlCommandType = iota
)

type SqlCommand struct {
	Type    SqlCommandType
	Clauses []*SqlClause
}

func NewSqlCommand(commandType SqlCommandType) *SqlCommand {
	return &SqlCommand{
		Type:    commandType,
		Clauses: make([]*SqlClause, 0),
	}
}
