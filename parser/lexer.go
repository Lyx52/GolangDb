package parser

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

const CARRIAGE_NEWLINE = "\r\n"
const SQL_SINGLE_LINE_COMMENT = "--"
const SQL_BINARY_NUMBER = "0b"
const SQL_HEX_NUMBER = "0x"
const SQL_OCTAL_NUMBER = "0o"
const SQL_NOT_QUALITY = "<>"
const SQL_NOT_QUALITY_SECONDARY = "!="
const SQL_LESS_THAN_OR_EQUAL = "<="
const SQL_MORE_THAN_OR_EQUAL = ">="
const SQL_ORDER_BY = "ORDER BY"

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

func (lexer *BaseLexer) Count() int {
	return len(lexer.tokens)
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
	lexer.TokenizeCombineTokens()
	return nil
}

func (lexer *BaseLexer) TokenizeFirstPass() error {
	next := lexer.reader.PeekNext()
	for next > 0 {
		switch {
		case lexer.TryString(CARRIAGE_NEWLINE):
			pos := lexer.reader.pos
			err := lexer.reader.Consume(2)
			if err != nil {
				return err
			}
			lexer.PushToken(WHITESPACE, TOKEN_EMPTY_VALUE, pos)
		case lexer.TryString(SQL_SINGLE_LINE_COMMENT):
			err := lexer.reader.Consume(2)
			if err != nil {
				return err
			}
			lexer.tokenizeUntilNewline(COMMENT)
			lexer.PopLastToken()
		case lexer.TryString(SQL_BINARY_NUMBER):
			err := lexer.reader.Consume(2)
			if err != nil {
				return err
			}
			next = lexer.reader.PeekNext()
			builder := strings.Builder{}
			for next == '0' || next == '1' {
				builder.WriteRune(lexer.reader.Next())
				next = lexer.reader.PeekNext()
			}
			value, err := strconv.ParseInt(builder.String(), 2, 32)
			if err != nil {
				return err
			}

			lexer.PushToken(NUMBER, strconv.FormatInt(value, 10), lexer.reader.pos)
		case lexer.TryString(SQL_HEX_NUMBER):
			err := lexer.reader.Consume(2)
			if err != nil {
				return err
			}
			next = lexer.reader.PeekNext()
			builder := strings.Builder{}
			for IsLetterOrDigit(next) {
				builder.WriteRune(lexer.reader.Next())
				next = lexer.reader.PeekNext()
			}

			value, err := strconv.ParseInt(builder.String(), 16, 32)
			if err != nil {
				return err
			}

			lexer.PushToken(NUMBER, strconv.FormatInt(value, 10), lexer.reader.pos)
		case lexer.TryString(SQL_OCTAL_NUMBER):
			err := lexer.reader.Consume(2)
			if err != nil {
				return err
			}
			next = lexer.reader.PeekNext()
			builder := strings.Builder{}
			for next >= '0' && next <= '7' {
				builder.WriteRune(lexer.reader.Next())
				next = lexer.reader.PeekNext()
			}

			value, err := strconv.ParseInt(builder.String(), 8, 32)
			if err != nil {
				return err
			}

			lexer.PushToken(NUMBER, strconv.FormatInt(value, 10), lexer.reader.pos)
		case lexer.TryString(SQL_NOT_QUALITY) || lexer.TryString(SQL_NOT_QUALITY_SECONDARY):
			err := lexer.reader.Consume(2)
			if err != nil {
				return err
			}
			lexer.reader.Next()
			lexer.PushToken(OPERATOR_NOT_EQUALITY, TOKEN_EMPTY_VALUE, lexer.reader.pos)
		case lexer.TryString(SQL_MORE_THAN_OR_EQUAL):
			err := lexer.reader.Consume(2)
			if err != nil {
				return err
			}
			lexer.reader.Next()
			lexer.PushToken(OPERATOR_GREATER_THAN_OR_EQUAL, TOKEN_EMPTY_VALUE, lexer.reader.pos)
		case lexer.TryString(SQL_LESS_THAN_OR_EQUAL):
			err := lexer.reader.Consume(2)
			if err != nil {
				return err
			}
			lexer.reader.Next()
			lexer.PushToken(OPERATOR_LESS_THAN_OR_EQUAL, TOKEN_EMPTY_VALUE, lexer.reader.pos)
		case lexer.TryString(SQL_ORDER_BY):
			err := lexer.reader.Consume(8)
			if err != nil {
				return err
			}
			lexer.reader.Next()
			lexer.PushToken(ORDER_BY, TOKEN_EMPTY_VALUE, lexer.reader.pos)
		case unicode.IsSpace(next):
			lexer.reader.Next()
			lexer.PushToken(WHITESPACE, TOKEN_EMPTY_VALUE, lexer.reader.pos)
		case next == '*':
			lexer.reader.Next()
			lexer.PushToken(WILDCARD, TOKEN_EMPTY_VALUE, lexer.reader.pos)
		case next == '>':
			lexer.reader.Next()
			lexer.PushToken(OPERATOR_GREATER_THAN, TOKEN_EMPTY_VALUE, lexer.reader.pos)
		case next == '<':
			lexer.reader.Next()
			lexer.PushToken(OPERATOR_LESS_THAN, TOKEN_EMPTY_VALUE, lexer.reader.pos)
		case next == '=':
			lexer.reader.Next()
			lexer.PushToken(OPERATOR_EQUALS, TOKEN_EMPTY_VALUE, lexer.reader.pos)
		case next == ',':
			lexer.reader.Next()
			lexer.PushToken(COMMA, TOKEN_EMPTY_VALUE, lexer.reader.pos)
		case next == '(':
			lexer.reader.Next()
			lexer.PushToken(BRACKET_OPEN, TOKEN_EMPTY_VALUE, lexer.reader.pos)
		case next == ')':
			lexer.reader.Next()
			lexer.PushToken(BRACKET_CLOSE, TOKEN_EMPTY_VALUE, lexer.reader.pos)
		case next == '[':
			lexer.reader.Next()
			lexer.PushToken(SQUARE_BRACKET_OPEN, TOKEN_EMPTY_VALUE, lexer.reader.pos)
		case next == ']':
			lexer.reader.Next()
			lexer.PushToken(SQUARE_BRACKET_CLOSE, TOKEN_EMPTY_VALUE, lexer.reader.pos)
		case next == ';':
			lexer.reader.Next()
			lexer.PushToken(SEMICOLUMN, TOKEN_EMPTY_VALUE, lexer.reader.pos)
		case next == ':':
			lexer.reader.Next()
			lexer.PushToken(COLON, TOKEN_EMPTY_VALUE, lexer.reader.pos)
		case next == '-':
			lexer.reader.Next()
			lexer.PushToken(MINUS, TOKEN_EMPTY_VALUE, lexer.reader.pos)
		case next == '\'':
			err := lexer.tokenizeEnclosedString(next, STRING)
			if err != nil {
				return err
			}
		case next == '"' || next == '`':
			err := lexer.tokenizeEnclosedString(next, QUOTED_IDENTIFIER)
			if err != nil {
				return err
			}
		case next == '.':
			lexer.reader.Next()
			lexer.PushToken(DOT, TOKEN_EMPTY_VALUE, lexer.reader.pos)
		case IsDigit(next):
			allowDot := true
			allowScientific := true
			builder := strings.Builder{}
			for IsDigit(next) || (allowDot && next == '.') || (allowScientific && next == 'e') {
				builder.WriteRune(lexer.reader.Next())
				if next == '.' {
					allowDot = false
				}

				if next == 'e' {
					allowScientific = false
					if lexer.reader.PeekNext() == '-' {
						builder.WriteRune(lexer.reader.Next())
					}
				}

				next = lexer.reader.PeekNext()
			}

			lexer.PushToken(NUMBER, builder.String(), lexer.reader.pos)

		case IsLetter(next) || next == '_':
			builder := strings.Builder{}

			for IsLetterOrDigit(next) || next == '_' || next == '"' || next == '\'' {
				builder.WriteRune(lexer.reader.Next())
				next = lexer.reader.PeekNext()
			}

			stringValue := builder.String()
			lexer.PushToken(IDENTIFIER, stringValue, lexer.reader.pos)
		default:
			return fmt.Errorf("unexpected token '%s' at %v", string(next), lexer.reader.pos)
		}
		next = lexer.reader.PeekNext()
	}

	return nil
}

func (lexer *BaseLexer) TokenizeCombineTokens() {
	tokens := make([]Token, 0)
	current := lexer.PopToken()
	next := lexer.PeekToken()
	for current != nil {
		// COMBINE - AND NUMBER into negative numbers
		if current.Type == MINUS && next != nil && next.Type == NUMBER {
			numeric := *lexer.PopToken()
			numeric.Value = "-" + numeric.Value
			tokens = append(tokens, numeric)

			// COMBINE WHITESPACES into ONE whitespace
		} else if current.Type == WHITESPACE {
			lexer.ConsumeTokens(WHITESPACE)
			tokens = append(tokens, *current)
		} else {
			tokens = append(tokens, *current)
		}

		current = lexer.PopToken()
		next = lexer.PeekToken()
	}

	lexer.tokens = tokens
}

func (lexer *BaseLexer) GetStatements() [][]Token {
	statements := make([][]Token, 0)
	current := make([]Token, 0)
	for _, token := range lexer.tokens {
		if len(current) != 0 && token.Type == SEMICOLUMN {
			statements = append(statements, current)
			current = make([]Token, 0)
			continue
		}

		current = append(current, token)
	}

	if len(current) != 0 {
		statements = append(statements, current)
	}

	return statements
}

func (lexer *BaseLexer) GetLexers() []*BaseLexer {
	statements := lexer.GetStatements()
	lexers := make([]*BaseLexer, len(statements))
	for i, statement := range statements {
		lexers[i] = &BaseLexer{
			tokens: statement,
		}
	}

	return lexers
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
			continue
		}

		if keywordType, exists := StringToTypeName[tokenUpper]; exists {
			lexer.tokens[i].Type = keywordType
			lexer.tokens[i].Value = tokenUpper
			continue
		}
	}
}

func (lexer *BaseLexer) tokenizeEnclosedString(enclosedWith rune, tokenType TokenType) error {
	builder := strings.Builder{}
	startPos := lexer.reader.pos
	next := lexer.reader.PeekNext()
	if next != enclosedWith {
		return fmt.Errorf("expected string to start with ' got %c", next)
	}

	builder.WriteRune(lexer.reader.Next())
	next = lexer.reader.PeekNext()

	for next != enclosedWith && next != -1 {
		builder.WriteRune(lexer.reader.Next())
		next = lexer.reader.PeekNext()
	}

	if next == -1 {
		return fmt.Errorf("expected string to end with ' got EOF at %v", startPos)
	}

	builder.WriteRune(lexer.reader.Next())
	stringValue := builder.String()
	lexer.PushToken(tokenType, stringValue[1:len(stringValue)-1], startPos)

	return nil
}

func (lexer *BaseLexer) tokenizeUntilNewline(tokenType TokenType) {
	builder := strings.Builder{}
	startPos := lexer.reader.pos

	for !lexer.nextIsNewline() {
		builder.WriteRune(lexer.reader.Next())
	}

	stringValue := builder.String()
	lexer.PushToken(tokenType, stringValue, startPos)
}

func (lexer *BaseLexer) nextIsNewline() bool {
	return lexer.reader.PeekNext() == '\n' || *lexer.reader.Peek(2) == CARRIAGE_NEWLINE
}
