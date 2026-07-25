package parser

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

func (lexer *BaseLexer) PushToken(token TokenType, value string, position int) {
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

func (lexer *BaseLexer) PopLastToken() *Token {
	if len(lexer.tokens) == 0 {
		return nil
	}

	res := lexer.tokens[len(lexer.tokens)-1]
	lexer.tokens = lexer.tokens[:len(lexer.tokens)-1]
	return &res
}

func (lexer *BaseLexer) PeekLastToken() *Token {
	if len(lexer.tokens) == 0 {
		return nil
	}

	return &lexer.tokens[len(lexer.tokens)-1]
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
			lexer.PushToken(WHITESPACE, TOKEN_EMPTY_VALUE, lexer.reader.pos)
		case next == '*':
			lexer.reader.Next()
			lexer.PushToken(WILDCARD, TOKEN_EMPTY_VALUE, lexer.reader.pos)
		case next == ',':
			lexer.reader.Next()
			lexer.PushToken(COMMA, TOKEN_EMPTY_VALUE, lexer.reader.pos)
		case next == '=':
			lexer.reader.Next()
			lexer.PushToken(OPERATOR_EQUALS, TOKEN_EMPTY_VALUE, lexer.reader.pos)
		case next == '(':
			lexer.reader.Next()
			lexer.PushToken(BRACKET_OPEN, TOKEN_EMPTY_VALUE, lexer.reader.pos)
		case next == ')':
			lexer.reader.Next()
			lexer.PushToken(BRACKET_CLOSE, TOKEN_EMPTY_VALUE, lexer.reader.pos)
		case next == ';':
			lexer.reader.Next()
			lexer.PushToken(SEMICOLUMN, TOKEN_EMPTY_VALUE, lexer.reader.pos)
		case next == '-':
			lexer.reader.Next()
			lexer.PushToken(MINUS, TOKEN_EMPTY_VALUE, lexer.reader.pos)
		case next == '\'':
			builder := strings.Builder{}

			for next != '\'' {
				builder.WriteRune(lexer.reader.Next())
				next = lexer.reader.PeekNext()
			}

		case IsLetterOrDigit(next) || next == '"' || next == '\'':
			builder := strings.Builder{}

			number := IsDigit(next)
			for IsLetterOrDigit(next) || next == '_' || next == '.' || next == '"' || next == '\'' || next == ' ' {
				builder.WriteRune(lexer.reader.Next())
				number = number && (IsDigit(next) || next == '.')
				next = lexer.reader.PeekNext()
			}

			stringValue := builder.String()
			if number {
				_, err := strconv.ParseFloat(stringValue, 64)
				if err != nil {
					return err
				}
				previous := lexer.PeekLastToken()
				if previous != nil && previous.Type == MINUS {
					lexer.PopLastToken()
					stringValue = "-" + stringValue
				}
				lexer.PushToken(NUMBER, stringValue, lexer.reader.pos)
			} else {
				lastIndex := len(stringValue) - 1

				if stringValue[0] == '\'' && stringValue[lastIndex] == '\'' {
					lexer.PushToken(STRING, stringValue[1:lastIndex], lexer.reader.pos)
				} else if stringValue[0] == '"' && stringValue[lastIndex] == '"' {
					lexer.PushToken(QUOTED_IDENTIFIER, stringValue[1:lastIndex], lexer.reader.pos)
				} else {
					lexer.PushToken(IDENTIFIER, stringValue, lexer.reader.pos)
				}
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
		if token.Type != IDENTIFIER {
			continue
		}

		tokenUpper := strings.ToUpper(token.Value)

		if keywordType, exists := StringToKeyword[tokenUpper]; exists {
			lexer.tokens[i].Type = keywordType
			lexer.tokens[i].Value = TOKEN_EMPTY_VALUE
		}
	}
}
