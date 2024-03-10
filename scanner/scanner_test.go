package scanner

import (
	"fmt"
	"go-interpreter-practice/token"
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
		wantTokens []token.Token
		wantErrors []error
	}{
		{
			input:      "(",
			wantTokens: []token.Token{{Type: token.LEFT_PAREN, RawToken: "(", Line: 1}, {Type: token.EOF}},
			wantErrors: []error{},
		},
		{
			input:      ")",
			wantTokens: []token.Token{{Type: token.RIGHT_PAREN, RawToken: ")", Line: 1}, {Type: token.EOF}},
			wantErrors: []error{},
		},
		{
			input:      "{",
			wantTokens: []token.Token{{Type: token.LEFT_BRACE, RawToken: "{", Line: 1}, {Type: token.EOF}},
			wantErrors: []error{},
		},
		{
			input:      "}",
			wantTokens: []token.Token{{Type: token.RIGHT_BRACE, RawToken: "}", Line: 1}, {Type: token.EOF}},
			wantErrors: []error{},
		},
		{
			input:      ",",
			wantTokens: []token.Token{{Type: token.COMMA, RawToken: ",", Line: 1}, {Type: token.EOF}},
			wantErrors: []error{},
		},
		{
			input:      ".",
			wantTokens: []token.Token{{Type: token.DOT, RawToken: ".", Line: 1}, {Type: token.EOF}},
			wantErrors: []error{},
		},
		{
			input:      "-",
			wantTokens: []token.Token{{Type: token.MINUS, RawToken: "-", Line: 1}, {Type: token.EOF}},
			wantErrors: []error{},
		},
		{
			input:      "+",
			wantTokens: []token.Token{{Type: token.PLUS, RawToken: "+", Line: 1}, {Type: token.EOF}},
			wantErrors: []error{},
		},
		{
			input:      ";",
			wantTokens: []token.Token{{Type: token.SEMICOLON, RawToken: ";", Line: 1}, {Type: token.EOF}},
			wantErrors: []error{},
		},
		{
			input:      "*",
			wantTokens: []token.Token{{Type: token.STAR, RawToken: "*", Line: 1}, {Type: token.EOF}},
			wantErrors: []error{},
		},
		{
			input:      "!=",
			wantTokens: []token.Token{{Type: token.NOT_EQUAL, RawToken: "!=", Line: 1}, {Type: token.EOF}},
			wantErrors: []error{},
		},
		{
			input:      "=",
			wantTokens: []token.Token{{Type: token.EQUAL, RawToken: "=", Line: 1}, {Type: token.EOF}},
			wantErrors: []error{},
		},
		{
			input:      "<=",
			wantTokens: []token.Token{{Type: token.LESS_EQUAL, RawToken: "<=", Line: 1}, {Type: token.EOF}},
			wantErrors: []error{},
		},
		{
			input:      ">=",
			wantTokens: []token.Token{{Type: token.GREATER_EQUAL, RawToken: ">=", Line: 1}, {Type: token.EOF}},
			wantErrors: []error{},
		},
		{
			input:      "/",
			wantTokens: []token.Token{{Type: token.SLASH, RawToken: "/", Line: 1}, {Type: token.EOF}},
			wantErrors: []error{},
		},
		{
			input:      "\n",
			wantTokens: []token.Token{{Type: token.EOF}},
			wantErrors: []error{},
		},
		{
			input:      "123",
			wantTokens: []token.Token{{Type: token.INTEGER, RawToken: "123", Literal: 123, Line: 1}, {Type: token.EOF}},
			wantErrors: []error{},
		},
		{
			input:      "123.456",
			wantTokens: []token.Token{{Type: token.FLOAT, RawToken: "123.456", Literal: 123.456, Line: 1}, {Type: token.EOF}},
			wantErrors: []error{},
		},
		{
			input:      "abc",
			wantTokens: []token.Token{{Type: token.IDENTIFIER, RawToken: "abc", Literal: "abc", Line: 1}, {Type: token.EOF}},
			wantErrors: []error{},
		},
		{
			input:      "@",
			wantTokens: []token.Token{{Type: token.EOF}},
			wantErrors: []error{fmt.Errorf("line 1: Unexpected character: @")},
		},
		{
			input: "var a = 1;",
			wantTokens: []token.Token{
				{Type: token.VAR, RawToken: "var", Line: 1},
				{Type: token.IDENTIFIER, RawToken: "a", Literal: "a", Line: 1},
				{Type: token.EQUAL, RawToken: "=", Line: 1},
				{Type: token.INTEGER, RawToken: "1", Literal: 1, Line: 1},
				{Type: token.SEMICOLON, RawToken: ";", Line: 1},
				{Type: token.EOF},
			},
			wantErrors: []error{},
		},
		{
			input: inputWithNewLines,
			wantTokens: []token.Token{
				{Type: token.VAR, RawToken: "var", Literal: nil, Line: 1},
				{Type: token.IDENTIFIER, RawToken: "a", Literal: "a", Line: 1},
				{Type: token.EQUAL, RawToken: "=", Literal: nil, Line: 1},
				{Type: token.INTEGER, RawToken: "1", Literal: 1, Line: 1},
				{Type: token.SEMICOLON, RawToken: ";", Literal: nil, Line: 1},
				{Type: token.IF, RawToken: "if", Literal: nil, Line: 2},
				{Type: token.IDENTIFIER, RawToken: "hoge", Literal: "hoge", Line: 2},
				{Type: token.LEFT_BRACE, RawToken: "{", Literal: nil, Line: 2},
				{Type: token.VAR, RawToken: "var", Literal: nil, Line: 3},
				{Type: token.IDENTIFIER, RawToken: "c", Literal: "c", Line: 3},
				{Type: token.EQUAL, RawToken: "=", Literal: nil, Line: 3},
				{Type: token.FLOAT, RawToken: "10.21", Literal: 10.21, Line: 3},
				{Type: token.SEMICOLON, RawToken: ";", Literal: nil, Line: 3},
				{Type: token.RIGHT_BRACE, RawToken: "}", Literal: nil, Line: 4},
				{Type: token.IDENTIFIER, RawToken: "a", Literal: "a", Line: 5},
				{Type: token.PLUS, RawToken: "+", Literal: nil, Line: 5},
				{Type: token.PLUS, RawToken: "+", Literal: nil, Line: 5},
				{Type: token.SEMICOLON, RawToken: ";", Literal: nil, Line: 5},
				{Type: token.VAR, RawToken: "var", Literal: nil, Line: 7},
				{Type: token.IDENTIFIER, RawToken: "b", Literal: "b", Line: 7},
				{Type: token.EQUAL, RawToken: "=", Literal: nil, Line: 7},
				{Type: token.TRUE, RawToken: "false", Literal: false, Line: 7},
				{Type: token.SEMICOLON, RawToken: ";", Literal: nil, Line: 7},
				{Type: token.EOF, RawToken: "", Literal: nil, Line: 0},
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
