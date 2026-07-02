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
	err := lexer.TokenizeFirstPass()
	if err != nil {
		return err
	}

	lexer.TokenizeKeywordPass()
	return nil
}

func (lexer *BaseLexer) TokenizeFirstPass() error {
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

func (lexer *BaseLexer) TokenizeKeywordPass() {
	for i, token := range lexer.tokens {
		if token.Type != STRING {
			continue
		}

		switch token.StringValue() {
		case SelectTokenString:
			lexer.tokens[i].Type = SELECT
			lexer.tokens[i].Value = nil
			break
		case DeleteTokenString:
			lexer.tokens[i].Type = DELETE
			lexer.tokens[i].Value = nil
			break
		case UpdateTokenString:
			lexer.tokens[i].Type = UPDATE
			lexer.tokens[i].Value = nil
			break
		case InsertTokenString:
			lexer.tokens[i].Type = INSERT
			lexer.tokens[i].Value = nil
			break
		case CreateTokenString:
			lexer.tokens[i].Type = CREATE
			lexer.tokens[i].Value = nil
			break
		case DatabasesTokenString:
			lexer.tokens[i].Type = DATABASES
			lexer.tokens[i].Value = nil
			break
		case DatabaseTokenString:
			lexer.tokens[i].Type = DATABASE
			lexer.tokens[i].Value = nil
			break
		case TablesTokenString:
			lexer.tokens[i].Type = TABLES
			lexer.tokens[i].Value = nil
			break
		case TableTokenString:
			lexer.tokens[i].Type = TABLE
			lexer.tokens[i].Value = nil
			break
		case ViewTokenString:
			lexer.tokens[i].Type = VIEW
			lexer.tokens[i].Value = nil
			break
		case IntoTokenString:
			lexer.tokens[i].Type = INTO
			lexer.tokens[i].Value = nil
			break
		case AsTokenString:
			lexer.tokens[i].Type = AS
			lexer.tokens[i].Value = nil
			break
		case SetTokenString:
			lexer.tokens[i].Type = SET
			lexer.tokens[i].Value = nil
			break
		case ValuesTokenString:
			lexer.tokens[i].Type = VALUES
			lexer.tokens[i].Value = nil
			break
		case WhereTokenString:
			lexer.tokens[i].Type = WHERE
			lexer.tokens[i].Value = nil
			break
		case AndTokenString:
			lexer.tokens[i].Type = AND
			lexer.tokens[i].Value = nil
			break
		case OrTokenString:
			lexer.tokens[i].Type = OR
			lexer.tokens[i].Value = nil
			break
		case InTokenString:
			lexer.tokens[i].Type = IN
			lexer.tokens[i].Value = nil
			break
		case UseTokenString:
			lexer.tokens[i].Type = USE
			lexer.tokens[i].Value = nil
			break
		case ShowTokenString:
			lexer.tokens[i].Type = SHOW
			lexer.tokens[i].Value = nil
			break
		}
	}
}
