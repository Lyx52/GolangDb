package parser

import (
	"testing"
)

func TestFail(t *testing.T) {
	sq := `'john doe' "O'Connor"`
	tokenizer := NewLexer(&sq)
	err := tokenizer.Tokenize()
	if err == nil {
		panic(err)
	}
}

func TestBaseLexer_IdentifierStringNumberResolution(t *testing.T) {
	// We define a small struct to hold our expected token results.
	// Value is now explicitly a string.
	type ExpectedToken struct {
		Type  TokenType
		Value string
	}

	tests := []struct {
		name        string
		sql         string
		expectError bool
		expected    []ExpectedToken
	}{
		{
			name:        "Standard Identifiers",
			sql:         "users first_name table123",
			expectError: false,
			expected: []ExpectedToken{
				{Type: IDENTIFIER, Value: "users"},
				{Type: WHITESPACE, Value: ""},
				{Type: IDENTIFIER, Value: "first_name"},
				{Type: WHITESPACE, Value: ""},
				{Type: IDENTIFIER, Value: "table123"},
			},
		},
		{
			name:        "String Literals (Single and Double Quotes)",
			sql:         `'john doe' "O'Connor"`,
			expectError: false,
			expected: []ExpectedToken{
				{Type: STRING, Value: "john doe"}, // Assuming your lexer strips the outer quotes
				{Type: WHITESPACE, Value: ""},
				{Type: STRING, Value: "O'Connor"},
			},
		},
		{
			name:        "Positive Numbers",
			sql:         "42 3.14 0.99",
			expectError: false,
			expected: []ExpectedToken{
				{Type: NUMBER, Value: "42"},
				{Type: WHITESPACE, Value: ""},
				{Type: NUMBER, Value: "3.14"},
				{Type: WHITESPACE, Value: ""},
				{Type: NUMBER, Value: "0.99"},
			},
		},
		{
			name:        "Negative Numbers",
			sql:         "-42 -3.14",
			expectError: false,
			expected: []ExpectedToken{
				{Type: NUMBER, Value: "-42"},
				{Type: WHITESPACE, Value: ""},
				{Type: NUMBER, Value: "-3.14"},
			},
		},
		{
			name:        "Mixed Query (Ensuring no minus-trap)",
			sql:         "score = -50 AND name = 'alice'",
			expectError: false,
			expected: []ExpectedToken{
				{Type: IDENTIFIER, Value: "score"},
				{Type: WHITESPACE, Value: ""},
				{Type: OPERATOR_EQUALS, Value: ""},
				{Type: WHITESPACE, Value: ""},
				{Type: NUMBER, Value: "-50"},
				{Type: WHITESPACE, Value: ""},
				{Type: IDENTIFIER, Value: "AND"}, // Keyword pass will later change this to the AND TokenType
				{Type: WHITESPACE, Value: ""},
				{Type: IDENTIFIER, Value: "name"},
				{Type: WHITESPACE, Value: ""},
				{Type: OPERATOR_EQUALS, Value: ""},
				{Type: WHITESPACE, Value: ""},
				{Type: STRING, Value: "alice"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(&tt.sql)
			err := lexer.TokenizeFirstPass()

			if tt.expectError {
				if err == nil {
					t.Fatalf("Expected error for sql %q, got nil", tt.sql)
				}
				return // Stop checking tokens if we expected an error
			}

			if err != nil {
				t.Fatalf("Expected no error for sql %q, got %v", tt.sql, err)
			}

			if len(lexer.tokens) != len(tt.expected) {
				t.Fatalf("Token count mismatch. Expected %d, got %d.\nTokens: %v", len(tt.expected), len(lexer.tokens), lexer.tokens)
			}

			for i, expectedTok := range tt.expected {
				actualTok := lexer.tokens[i]

				if actualTok.Type != expectedTok.Type {
					t.Errorf("Token [%d] Type mismatch. Expected %v, got %v (Value: %q)", i, expectedTok.Type, actualTok.Type, actualTok.Value)
				}

				// Only check string values for tokens that actually capture strings (numbers, strings, identifiers).
				// We skip checking exact string matches for whitespace/operators just in case your lexer
				// assigns them the literal string (e.g. " ") instead of an empty string ("").
				if expectedTok.Value != "" {
					if actualTok.Value != expectedTok.Value {
						t.Errorf("Token [%d] Value mismatch. Expected %q, got %q", i, expectedTok.Value, actualTok.Value)
					}
				}
			}
		})
	}
}
