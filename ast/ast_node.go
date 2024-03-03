package ast

type Node interface {
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p Program) String() string {
	return ""
}

func (p *Program) ParseStatement() {
	for _, s := range p.Statements {
		s.statementNode()
	}
}
