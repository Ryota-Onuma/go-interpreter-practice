package parser

import (
	"go-interpreter-practice/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

var priorityMap = map[token.TokenType]int{
	token.EQUAL_EQUAL:   EQUALS,
	token.NOT_EQUAL:     EQUALS,
	token.LESS:          LESSGREATER,
	token.LESS_EQUAL:    LESSGREATER,
	token.GREATER:       LESSGREATER,
	token.GREATER_EQUAL: LESSGREATER,
	token.PLUS:          SUM,
	token.MINUS:         SUM,
	token.SLASH:         PRODUCT,
	token.STAR:          PRODUCT,
	token.LEFT_PAREN:    CALL,
}
