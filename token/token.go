package token

type TokenType string

const (
	LEFT_PAREN  TokenType = "LEFT_PAREN"  // (
	RIGHT_PAREN TokenType = "RIGHT_PAREN" // )
	LEFT_BRACE  TokenType = "LEFT_BRACE"  // {
	RIGHT_BRACE TokenType = "RIGHT_BRACE" // }
	COMMA       TokenType = "COMMA"       // ,
	DOT         TokenType = "DOT"         // .
	MINUS       TokenType = "MINUS"       // -
	PLUS        TokenType = "PLUS"        // +
	SEMICOLON   TokenType = "SEMICOLON"   // ;
	SLASH       TokenType = "SLASH"       // /
	STAR        TokenType = "STAR"        // *
	LINE_BREAK  TokenType = "LINE_BREAK"  // \n

	BANG          TokenType = "BANG"          // !
	NOT_EQUAL     TokenType = "NOT_EQUAL"     // !=
	EQUAL         TokenType = "EQUAL"         // =
	EQUAL_EQUAL   TokenType = "EQUAL_EQUAL"   // ==
	GREATER       TokenType = "GREATER"       // >
	GREATER_EQUAL TokenType = "GREATER_EQUAL" // >=
	LESS          TokenType = "LESS"          // <
	LESS_EQUAL    TokenType = "LESS_EQUAL"    // <=

	IDENTIFIER TokenType = "IDENTIFIER" // variable name
	STRING     TokenType = "STRING"     // "string"
	INTEGER    TokenType = "INTEGER"    // 123
	FLOAT      TokenType = "FLOAT"      // 123.45

	AND    TokenType = "AND"    // and
	CLASS  TokenType = "CLASS"  // class
	ELSE   TokenType = "ELSE"   // else
	FALSE  TokenType = "FALSE"  // false
	FUN    TokenType = "FUN"    // fun
	FOR    TokenType = "FOR"    // for
	IF     TokenType = "IF"     // if
	NIL    TokenType = "NIL"    // nil
	OR     TokenType = "OR"     // or
	PRINT  TokenType = "PRINT"  // print
	RETURN TokenType = "RETURN" // return
	SUPER  TokenType = "SUPER"  // super
	THIS   TokenType = "THIS"   // this
	TRUE   TokenType = "TRUE"   // true
	VAR    TokenType = "VAR"    // var
	WHILE  TokenType = "WHILE"  // while

	EOF TokenType = "EOF" // end of file
)

type Token struct {
	Type     TokenType
	RawToken string // ソース・ファイルから取得した生の状態、
	Literal  any
	Line     int // 行番号
}

func (t Token) String() string {
	return ""
}
