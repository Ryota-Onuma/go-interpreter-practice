package ast

import (
	"lox-by-go/token"
)

type Identifier struct {
	Token token.Token
	Value string
}

func (i Identifier) expressionNode() {}
func (i Identifier) String() string {
	return ""
}

type IntegerLiteral struct {
	Token token.Token
	Value int
}

func (il IntegerLiteral) expressionNode() {}
func (il IntegerLiteral) String() string {
	return ""
}

type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode() {}
func (fl *FloatLiteral) String() string {
	return ""
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe PrefixExpression) expressionNode() {}
func (pe PrefixExpression) String() string {
	return ""
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie InfixExpression) expressionNode() {}
func (ie InfixExpression) String() string {
	return ""
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b Boolean) expressionNode() {}
func (b Boolean) String() string {
	return ""
}
