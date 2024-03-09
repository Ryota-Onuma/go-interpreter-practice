package ast

import (
	"lox-by-go/token"
)

type VarStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (vs VarStatement) statementNode() {}
func (vs VarStatement) String() string {
	return ""
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (es ReturnStatement) statementNode() {}
func (rs ReturnStatement) String() string {
	return ""
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es ExpressionStatement) statementNode() {}
func (es ExpressionStatement) String() string {
	return ""
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs BlockStatement) statementNode() {}
func (bs BlockStatement) String() string {
	return ""
}
