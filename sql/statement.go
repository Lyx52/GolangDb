package sql

type SqlStatement struct {
	Commands []*SqlCommand
}

func NewSqlStatement() *SqlStatement {
	return &SqlStatement{
		Commands: make([]*SqlCommand, 0),
	}
}
