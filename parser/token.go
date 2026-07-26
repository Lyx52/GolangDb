package parser

type TokenType int

const TOKEN_EMPTY_VALUE = ""
const (
	WHITESPACE            TokenType = iota
	WILDCARD              TokenType = iota
	DOT                   TokenType = iota
	IDENTIFIER            TokenType = iota
	QUOTED_IDENTIFIER     TokenType = iota
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
	MINUS                 TokenType = iota
	COMMENT               TokenType = iota
	ALTER                 TokenType = iota
	DROP                  TokenType = iota
)

func (tokenType TokenType) String() string {
	switch tokenType {
	case WHITESPACE:
		return "WHITESPACE"
	case WILDCARD:
		return "WILDCARD"
	case COMMA:
		return "COMMA"
	case DOT:
		return "DOT"
	case IDENTIFIER:
		return "IDENTIFIER"
	case QUOTED_IDENTIFIER:
		return "QUOTED_IDENTIFIER"
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
	case ALTER:
		return "ALTER"
	case DROP:
		return "DROP"
	default:
		return "UNKNOWN"
	}
}

var StringToKeyword = map[string]TokenType{
	"SELECT":    SELECT,
	"INSERT":    INSERT,
	"UPDATE":    UPDATE,
	"DELETE":    DELETE,
	"CREATE":    CREATE,
	"DATABASE":  DATABASE,
	"TABLE":     TABLE,
	"DATABASES": DATABASES,
	"TABLES":    TABLES,
	"VIEW":      VIEW,
	"FROM":      FROM,
	"INTO":      INTO,
	"AS":        AS,
	"SET":       SET,
	"VALUES":    VALUES,
	"WHERE":     WHERE,
	"AND":       AND,
	"OR":        OR,
	"IN":        IN,
	"USE":       USE,
	"SHOW":      SHOW,
	"ALTER":     ALTER,
	"DROP":      DROP,
}
var CommandTokenType = []TokenType{
	SELECT,
	ALTER,
	UPDATE,
	CREATE,
	DELETE,
	DROP,
	SET,
	INSERT,
}

type Token struct {
	Type     TokenType
	Value    string
	Position int
}

func IsCommandToken(token *Token) bool {
	for _, tokenType := range CommandTokenType {
		if tokenType == token.Type {
			return true
		}
	}

	return false
}

func IsTokenType(typ TokenType, allowed ...TokenType) bool {
	for _, v := range allowed {
		if v == typ {
			return true
		}
	}

	return false
}
