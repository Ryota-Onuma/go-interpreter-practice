package scanner

import (
	"fmt"
	"testing"
)

func TestScanner(t *testing.T) {

	inputWithNewLines := `var a = 1;
	if hoge {
		var c = 10.21;
	}
	a++;
	// hogehoge
	var b = false;
	`

	testCases := []struct {
		input      string
		wantTokens []Token
		wantErrors []error
	}{
		{
			input:      "(",
			wantTokens: []Token{{Type: LEFT_PAREN, RawToken: "(", Line: 1}, {Type: EOF}},
			wantErrors: []error{},
		},
		{
			input:      ")",
			wantTokens: []Token{{Type: RIGHT_PAREN, RawToken: ")", Line: 1}, {Type: EOF}},
			wantErrors: []error{},
		},
		{
			input:      "{",
			wantTokens: []Token{{Type: LEFT_BRACE, RawToken: "{", Line: 1}, {Type: EOF}},
			wantErrors: []error{},
		},
		{
			input:      "}",
			wantTokens: []Token{{Type: RIGHT_BRACE, RawToken: "}", Line: 1}, {Type: EOF}},
			wantErrors: []error{},
		},
		{
			input:      ",",
			wantTokens: []Token{{Type: COMMA, RawToken: ",", Line: 1}, {Type: EOF}},
			wantErrors: []error{},
		},
		{
			input:      ".",
			wantTokens: []Token{{Type: DOT, RawToken: ".", Line: 1}, {Type: EOF}},
			wantErrors: []error{},
		},
		{
			input:      "-",
			wantTokens: []Token{{Type: MINUS, RawToken: "-", Line: 1}, {Type: EOF}},
			wantErrors: []error{},
		},
		{
			input:      "+",
			wantTokens: []Token{{Type: PLUS, RawToken: "+", Line: 1}, {Type: EOF}},
			wantErrors: []error{},
		},
		{
			input:      ";",
			wantTokens: []Token{{Type: SEMICOLON, RawToken: ";", Line: 1}, {Type: EOF}},
			wantErrors: []error{},
		},
		{
			input:      "*",
			wantTokens: []Token{{Type: STAR, RawToken: "*", Line: 1}, {Type: EOF}},
			wantErrors: []error{},
		},
		{
			input:      "!=",
			wantTokens: []Token{{Type: BANG_EQUAL, RawToken: "!=", Line: 1}, {Type: EOF}},
			wantErrors: []error{},
		},
		{
			input:      "=",
			wantTokens: []Token{{Type: EQUAL, RawToken: "=", Line: 1}, {Type: EOF}},
			wantErrors: []error{},
		},
		{
			input:      "<=",
			wantTokens: []Token{{Type: LESS_EQUAL, RawToken: "<=", Line: 1}, {Type: EOF}},
			wantErrors: []error{},
		},
		{
			input:      ">=",
			wantTokens: []Token{{Type: GREATER_EQUAL, RawToken: ">=", Line: 1}, {Type: EOF}},
			wantErrors: []error{},
		},
		{
			input:      "/",
			wantTokens: []Token{{Type: SLASH, RawToken: "/", Line: 1}, {Type: EOF}},
			wantErrors: []error{},
		},
		{
			input:      "\n",
			wantTokens: []Token{{Type: EOF}},
			wantErrors: []error{},
		},
		{
			input:      "123",
			wantTokens: []Token{{Type: INTEGER, RawToken: "123", Literal: 123, Line: 1}, {Type: EOF}},
			wantErrors: []error{},
		},
		{
			input:      "123.456",
			wantTokens: []Token{{Type: FLOAT, RawToken: "123.456", Literal: 123.456, Line: 1}, {Type: EOF}},
			wantErrors: []error{},
		},
		{
			input:      "abc",
			wantTokens: []Token{{Type: IDENTIFIER, RawToken: "abc", Literal: "abc", Line: 1}, {Type: EOF}},
			wantErrors: []error{},
		},
		{
			input:      "@",
			wantTokens: []Token{{Type: EOF}},
			wantErrors: []error{fmt.Errorf("line 1: Unexpected character: @")},
		},
		{
			input: "var a = 1;",
			wantTokens: []Token{
				{Type: VAR, RawToken: "var", Line: 1},
				{Type: IDENTIFIER, RawToken: "a", Literal: "a", Line: 1},
				{Type: EQUAL, RawToken: "=", Line: 1},
				{Type: INTEGER, RawToken: "1", Literal: 1, Line: 1},
				{Type: SEMICOLON, RawToken: ";", Line: 1},
				{Type: EOF},
			},
			wantErrors: []error{},
		},
		{
			input: inputWithNewLines,
			wantTokens: []Token{
				{Type: VAR, RawToken: "var", Literal: nil, Line: 1},
				{Type: IDENTIFIER, RawToken: "a", Literal: "a", Line: 1},
				{Type: EQUAL, RawToken: "=", Literal: nil, Line: 1},
				{Type: INTEGER, RawToken: "1", Literal: 1, Line: 1},
				{Type: SEMICOLON, RawToken: ";", Literal: nil, Line: 1},
				{Type: IF, RawToken: "if", Literal: nil, Line: 2},
				{Type: IDENTIFIER, RawToken: "hoge", Literal: "hoge", Line: 2},
				{Type: LEFT_BRACE, RawToken: "{", Literal: nil, Line: 2},
				{Type: VAR, RawToken: "var", Literal: nil, Line: 3},
				{Type: IDENTIFIER, RawToken: "c", Literal: "c", Line: 3},
				{Type: EQUAL, RawToken: "=", Literal: nil, Line: 3},
				{Type: FLOAT, RawToken: "10.21", Literal: 10.21, Line: 3},
				{Type: SEMICOLON, RawToken: ";", Literal: nil, Line: 3},
				{Type: RIGHT_BRACE, RawToken: "}", Literal: nil, Line: 4},
				{Type: IDENTIFIER, RawToken: "a", Literal: "a", Line: 5},
				{Type: PLUS, RawToken: "+", Literal: nil, Line: 5},
				{Type: PLUS, RawToken: "+", Literal: nil, Line: 5},
				{Type: SEMICOLON, RawToken: ";", Literal: nil, Line: 5},
				{Type: VAR, RawToken: "var", Literal: nil, Line: 7},
				{Type: IDENTIFIER, RawToken: "b", Literal: "b", Line: 7},
				{Type: EQUAL, RawToken: "=", Literal: nil, Line: 7},
				{Type: TRUE, RawToken: "false", Literal: false, Line: 7},
				{Type: SEMICOLON, RawToken: ";", Literal: nil, Line: 7},
				{Type: EOF, RawToken: "", Literal: nil, Line: 0},
			},
			wantErrors: []error{},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("入力: %s", tc.input), func(t *testing.T) {
			s := NewScanner(tc.input)
			s.ScanTokens()
			for i, token := range s.Tokens() {
				if token.Type != tc.wantTokens[i].Type {
					t.Errorf("expected %v, but got %v", tc.wantTokens[i].Type, token.Type)
				}
				if token.RawToken != tc.wantTokens[i].RawToken {
					t.Errorf("expected %v, but got %v", tc.wantTokens[i].RawToken, token.RawToken)
				}
				if token.Literal != tc.wantTokens[i].Literal {
					t.Errorf("expected %v, but got %v", tc.wantTokens[i].Literal, token.Literal)
				}
				if token.Line != tc.wantTokens[i].Line {
					t.Errorf("expected %v, but got %v", tc.wantTokens[i].Line, token.Line)
				}
			}

			for i, err := range s.GetErrors() {
				if err.Error() != tc.wantErrors[i].Error() {
					t.Errorf("expected %v, but got %v", tc.wantErrors[i].Error(), err.Error())
				}
			}
		})
	}
}
