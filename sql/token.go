package sql

import "fmt"

type TokenType int

const (
	WHITESPACE            TokenType = iota
	WILDCARD              TokenType = iota
	DOT                   TokenType = iota
	COMMA                 TokenType = iota
	OPERATOR_EQUALS       TokenType = iota
	OPERATOR_NOT_EQUALITY TokenType = iota
	OPERATOR_LESS_THAN    TokenType = iota
	OPERATOR_GREATER_THAN TokenType = iota
	SELECT                TokenType = iota
	INSERT                TokenType = iota
	UPDATE                TokenType = iota
	DELETE                TokenType = iota
	CREATE                TokenType = iota
	DATABASE              TokenType = iota
	TABLE                 TokenType = iota
	DATABASES             TokenType = iota
	TABLES                TokenType = iota
	VIEW                  TokenType = iota
	SHOW                  TokenType = iota
	FROM                  TokenType = iota
	INTO                  TokenType = iota
	AS                    TokenType = iota
	SET                   TokenType = iota
	VALUES                TokenType = iota
	AND                   TokenType = iota
	OR                    TokenType = iota
	IN                    TokenType = iota
	USE                   TokenType = iota
	WHERE                 TokenType = iota
	STRING                TokenType = iota
	NUMBER                TokenType = iota
	FIELD                 TokenType = iota
	BRACKET_OPEN          TokenType = iota
	BRACKET_CLOSE         TokenType = iota
	SEMICOLUMN            TokenType = iota
)

func (tokenType TokenType) String() string {
	switch tokenType {
	case WHITESPACE:
		return "WHITESPACE"
	case WILDCARD:
		return "WILDCARD"
	case COMMA:
		return "COMMA"
	case OPERATOR_EQUALS:
		return "OPERATOR_EQUALITY"
	case OPERATOR_NOT_EQUALITY:
		return "OPERATOR_NOT_EQUALITY"
	case OPERATOR_LESS_THAN:
		return "OPERATOR_LESS_THAN"
	case OPERATOR_GREATER_THAN:
		return "OPERATOR_GREATER_THAN"
	case SELECT:
		return "SELECT"
	case INSERT:
		return "INSERT"
	case UPDATE:
		return "UPDATE"
	case DELETE:
		return "DELETE"
	case CREATE:
		return "CREATE"
	case DATABASE:
		return "DATABASE"
	case TABLE:
		return "TABLE"
	case DATABASES:
		return "DATABASES"
	case TABLES:
		return "TABLES"
	case VIEW:
		return "VIEW"
	case FROM:
		return "FROM"
	case INTO:
		return "INTO"
	case VALUES:
		return "VALUES"
	case STRING:
		return "STRING"
	case NUMBER:
		return "NUMBER"
	case FIELD:
		return "FIELD"
	case BRACKET_OPEN:
		return "BRACKET_OPEN"
	case BRACKET_CLOSE:
		return "BRACKET_CLOSE"
	case SEMICOLUMN:
		return "SEMICOLUMN"
	case AS:
		return "AS"
	case SET:
		return "SET"
	case OR:
		return "OR"
	case WHERE:
		return "WHERE"
	case AND:
		return "AND"
	case IN:
		return "IN"
	case USE:
		return "USE"
	case SHOW:
		return "SHOW"
	default:
		return "UNKNOWN"
	}
}

const (
	SelectTokenString    string = "SELECT"
	InsertTokenString    string = "INSERT"
	UpdateTokenString    string = "UPDATE"
	DeleteTokenString    string = "DELETE"
	CreateTokenString    string = "CREATE"
	DatabaseTokenString  string = "DATABASE"
	TableTokenString     string = "TABLE"
	DatabasesTokenString string = "DATABASES"
	TablesTokenString    string = "TABLES"
	ViewTokenString      string = "VIEW"
	FromTokenString      string = "FROM"
	IntoTokenString      string = "INTO"
	AsTokenString        string = "AS"
	SetTokenString       string = "SET"
	ValuesTokenString    string = "VALUES"
	WhereTokenString     string = "WHERE"
	AndTokenString       string = "AND"
	OrTokenString        string = "OR"
	InTokenString        string = "IN"
	UseTokenString       string = "USE"
	ShowTokenString      string = "SHOW"
)

type Token struct {
	Type     TokenType
	Value    any
	Position int
}

func (token Token) StringValue() string {
	return fmt.Sprint(token.Value)
}
