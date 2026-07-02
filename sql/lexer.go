package sql

import (
	"fmt"
	"strconv"
	"strings"
)

func IsDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func IsLetter(r rune) bool {
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z'
}

func IsLetterOrDigit(r rune) bool {
	return IsLetter(r) || IsDigit(r)
}

type BaseLexer struct {
	tokens []Token
	reader *StringReader
}

func NewLexer(sql *string) *BaseLexer {
	return &BaseLexer{
		tokens: make([]Token, 0),
		reader: NewStringReader(sql),
	}
}

func (lexer *BaseLexer) TryString(text string) bool {
	res := lexer.reader.Peek(len(text))
	if res == nil {
		return false
	}

	return strings.EqualFold(text, *res)
}

func (lexer *BaseLexer) TryConsumeString(text string) bool {
	if lexer.TryString(text) {
		err := lexer.reader.Consume(len(text))
		return err == nil
	}

	return false
}

func (lexer *BaseLexer) PushToken(token TokenType, value any, position int) {
	lexer.tokens = append(lexer.tokens, Token{Type: token, Value: value, Position: position})
}

func (lexer *BaseLexer) PopToken() *Token {
	if len(lexer.tokens) == 0 {
		return nil
	}

	res := lexer.tokens[0]
	lexer.tokens = lexer.tokens[1:]
	return &res
}

func (lexer *BaseLexer) PeekToken() *Token {
	if len(lexer.tokens) == 0 {
		return nil
	}

	return &lexer.tokens[0]
}

func (lexer *BaseLexer) ConsumeTokens(typ TokenType) int {
	token := lexer.PeekToken()
	count := 0
	for token != nil && token.Type == typ {
		lexer.PopToken()
		token = lexer.PeekToken()
		count++
	}

	return count
}

func (lexer *BaseLexer) ExpectToken(typ TokenType) error {
	token := lexer.PeekToken()

	if token == nil {
		return fmt.Errorf("expected token %s but got nil", typ.String())
	}

	if token.Type != typ {
		return fmt.Errorf("expected token %s but got %s", typ.String(), token.Type.String())
	}

	return nil
}

func (lexer *BaseLexer) ConsumeExpectToken(typ TokenType) error {
	err := lexer.ExpectToken(typ)
	if err != nil {
		return err
	}

	lexer.PopToken()
	return nil
}

func (lexer *BaseLexer) Tokenize() error {
	next := lexer.reader.PeekNext()
	for next > 0 {
		switch {
		case next == ' ' || next == '\t' || next == '\n':
			lexer.reader.Next()
			lexer.PushToken(WHITESPACE, nil, lexer.reader.pos)
		case next == '*':
			lexer.reader.Next()
			lexer.PushToken(WILDCARD, nil, lexer.reader.pos)
		case next == ',':
			lexer.reader.Next()
			lexer.PushToken(COMMA, nil, lexer.reader.pos)
		case next == '=':
			lexer.reader.Next()
			lexer.PushToken(OPERATOR_EQUALS, nil, lexer.reader.pos)
		case next == '(':
			lexer.reader.Next()
			lexer.PushToken(BRACKET_OPEN, nil, lexer.reader.pos)
		case next == ')':
			lexer.reader.Next()
			lexer.PushToken(BRACKET_CLOSE, nil, lexer.reader.pos)
		case next == ';':
			lexer.reader.Next()
			lexer.PushToken(SEMICOLUMN, nil, lexer.reader.pos)
		case lexer.TryConsumeString(SelectTokenString):
			lexer.PushToken(SELECT, nil, lexer.reader.pos)
		case lexer.TryConsumeString(DeleteTokenString):
			lexer.PushToken(DELETE, nil, lexer.reader.pos)
		case lexer.TryConsumeString(UpdateTokenString):
			lexer.PushToken(UPDATE, nil, lexer.reader.pos)
		case lexer.TryConsumeString(InsertTokenString):
			lexer.PushToken(INSERT, nil, lexer.reader.pos)
		case lexer.TryConsumeString(CreateTokenString):
			lexer.PushToken(CREATE, nil, lexer.reader.pos)

		case lexer.TryConsumeString(DatabasesTokenString):
			lexer.PushToken(DATABASES, nil, lexer.reader.pos)
		case lexer.TryConsumeString(DatabaseTokenString):
			lexer.PushToken(DATABASE, nil, lexer.reader.pos)
		case lexer.TryString(TablesTokenString):
			lexer.PushToken(TABLES, nil, lexer.reader.pos)
		case lexer.TryString(TableTokenString):
			lexer.PushToken(TABLE, nil, lexer.reader.pos)
		case lexer.TryConsumeString(ViewTokenString):
			lexer.PushToken(VIEW, nil, lexer.reader.pos)
		case lexer.TryConsumeString(IntoTokenString):
			lexer.PushToken(INTO, nil, lexer.reader.pos)
		case lexer.TryConsumeString(AsTokenString):
			lexer.PushToken(AS, nil, lexer.reader.pos)
		case lexer.TryConsumeString(SetTokenString):
			lexer.PushToken(SET, nil, lexer.reader.pos)
		case lexer.TryConsumeString(ValuesTokenString):
			lexer.PushToken(VALUES, nil, lexer.reader.pos)
		case lexer.TryConsumeString(WhereTokenString):
			lexer.PushToken(WHERE, nil, lexer.reader.pos)
		case lexer.TryConsumeString(AndTokenString):
			lexer.PushToken(AND, nil, lexer.reader.pos)
		case lexer.TryConsumeString(OrTokenString):
			lexer.PushToken(OR, nil, lexer.reader.pos)
		case lexer.TryConsumeString(InTokenString):
			lexer.PushToken(IN, nil, lexer.reader.pos)
		case lexer.TryConsumeString(UseTokenString):
			lexer.PushToken(USE, nil, lexer.reader.pos)
		case lexer.TryConsumeString(ShowTokenString):
			lexer.PushToken(SHOW, nil, lexer.reader.pos)
		case IsLetterOrDigit(next) || next == '-' || next == '"' || next == '\'':
			res := make([]rune, 0)
			negative := false
			if next == '-' {
				negative = true
				res = append(res, lexer.reader.Next())
				next = lexer.reader.PeekNext()
			}

			number := IsDigit(next)
			for IsLetterOrDigit(next) || next == '_' || next == '.' || next == '"' || next == '\'' {
				res = append(res, lexer.reader.Next())
				number = number && (IsDigit(next) || next == '.')
				next = lexer.reader.PeekNext()
			}

			if negative && !number {
				return fmt.Errorf("expected a number but got '%s'", string(res))
			}

			if number {
				value, err := strconv.ParseFloat(string(res), 64)
				if err != nil {
					return err
				}
				lexer.PushToken(NUMBER, value, lexer.reader.pos)
			} else {
				lexer.PushToken(STRING, string(res), lexer.reader.pos)
			}
		default:
			return fmt.Errorf("unexpected token '%s' at %v", string(next), lexer.reader.pos)
		}
		next = lexer.reader.PeekNext()
	}

	return nil
}
