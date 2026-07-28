package parser

type TokenType int

const TOKEN_EMPTY_VALUE = ""
const (
	WHITESPACE                     TokenType = iota
	WILDCARD                       TokenType = iota
	DOT                            TokenType = iota
	IDENTIFIER                     TokenType = iota
	QUOTED_IDENTIFIER              TokenType = iota
	COMMA                          TokenType = iota
	OPERATOR_EQUALS                TokenType = iota
	OPERATOR_NOT_EQUALITY          TokenType = iota
	OPERATOR_LESS_THAN             TokenType = iota
	OPERATOR_GREATER_THAN          TokenType = iota
	OPERATOR_LESS_THAN_OR_EQUAL    TokenType = iota
	OPERATOR_GREATER_THAN_OR_EQUAL TokenType = iota
	SELECT                         TokenType = iota
	INSERT                         TokenType = iota
	UPDATE                         TokenType = iota
	DELETE                         TokenType = iota
	CREATE                         TokenType = iota
	DATABASE                       TokenType = iota
	TABLE                          TokenType = iota
	DATABASES                      TokenType = iota
	TABLES                         TokenType = iota
	VIEW                           TokenType = iota
	SHOW                           TokenType = iota
	FROM                           TokenType = iota
	INTO                           TokenType = iota
	AS                             TokenType = iota
	SET                            TokenType = iota
	VALUES                         TokenType = iota
	AND                            TokenType = iota
	OR                             TokenType = iota
	IN                             TokenType = iota
	USE                            TokenType = iota
	WHERE                          TokenType = iota
	STRING                         TokenType = iota
	NUMBER                         TokenType = iota
	FIELD                          TokenType = iota
	BRACKET_OPEN                   TokenType = iota
	BRACKET_CLOSE                  TokenType = iota
	SQUARE_BRACKET_OPEN            TokenType = iota
	SQUARE_BRACKET_CLOSE           TokenType = iota
	SEMICOLUMN                     TokenType = iota
	COLON                          TokenType = iota
	MINUS                          TokenType = iota
	COMMENT                        TokenType = iota
	ALTER                          TokenType = iota
	DROP                           TokenType = iota
	CAST                           TokenType = iota
	TYPE_NAME                      TokenType = iota
	ORDER_BY                       TokenType = iota
	ORDER                          TokenType = iota
	BY                             TokenType = iota
	DESC                           TokenType = iota
	ASC                            TokenType = iota
)

func (tokenType TokenType) String() string {
	switch tokenType {
	case WHITESPACE:
		return "WHITESPACE"
	case WILDCARD:
		return "*"
	case COMMA:
		return ","
	case DOT:
		return "."
	case IDENTIFIER:
		return "IDENTIFIER"
	case QUOTED_IDENTIFIER:
		return "QUOTED_IDENTIFIER"
	case OPERATOR_EQUALS:
		return "="
	case OPERATOR_NOT_EQUALITY:
		return "<>"
	case OPERATOR_LESS_THAN:
		return "<"
	case OPERATOR_GREATER_THAN:
		return ">"
	case OPERATOR_LESS_THAN_OR_EQUAL:
		return "<="
	case OPERATOR_GREATER_THAN_OR_EQUAL:
		return ">="
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
		return "("
	case BRACKET_CLOSE:
		return ")"
	case SQUARE_BRACKET_OPEN:
		return "["
	case SQUARE_BRACKET_CLOSE:
		return "]"
	case SEMICOLUMN:
		return ";"
	case COLON:
		return ":"
	case AS:
		return "AS"
	case SET:
		return "SET"
	case OR:
		return "OR"
	case WHERE:
		return "WHERE"
	case ORDER_BY:
		return "ORDER BY"
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
	case CAST:
		return "CAST"
	case ORDER:
		return "ORDER"
	case BY:
		return "BY"
	case TYPE_NAME:
		return "TYPE_NAME"
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
	"CAST":      CAST,
	"DESC":      DESC,
	"ASC":       ASC,
	"ORDER":     ORDER,
	"BY":        BY,
}

var StringToTypeName = map[string]TokenType{
	"NUMERIC": TYPE_NAME,
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

var FunctionTokenType = []TokenType{
	CAST,
}

var OperatorTokenType = []TokenType{
	OPERATOR_GREATER_THAN_OR_EQUAL,
	OPERATOR_LESS_THAN_OR_EQUAL,
	OPERATOR_GREATER_THAN,
	OPERATOR_LESS_THAN,
	OPERATOR_EQUALS,
	OPERATOR_NOT_EQUALITY,
	AND,
	OR,
	IN,
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

func IsFunctionToken(token *Token) bool {
	for _, tokenType := range FunctionTokenType {
		if tokenType == token.Type {
			return true
		}
	}

	return false
}

func IsOperatorToken(token *Token) bool {
	for _, tokenType := range OperatorTokenType {
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
