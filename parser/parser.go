package parser

import (
	"fmt"
	"lox-by-go/ast"
	"lox-by-go/scanner"
	"lox-by-go/token"
	"strings"
)

type Parser struct {
	scanner      *scanner.Scanner
	tokens       []token.Token
	errors       []error
	currentAt    int
	currentToken token.Token

	prefixParseFns map[token.TokenType]func() ast.Expression
	infixParseFns  map[token.TokenType]func(ast.Expression) ast.Expression
}

func NewParser(s *scanner.Scanner) *Parser {
	p := &Parser{
		scanner: s,
	}
	p.prefixParseFns = make(map[token.TokenType]func() ast.Expression)
	p.infixParseFns = make(map[token.TokenType]func(ast.Expression) ast.Expression)

	p.registerPrefix(token.IDENTIFIER, p.parseIdentifier)
	p.registerPrefix(token.INTEGER, p.parseIntegerLiteral)
	p.registerPrefix(token.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)

	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.STAR, p.parseInfixExpression)
	p.registerInfix(token.EQUAL_EQUAL, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQUAL, p.parseInfixExpression)
	p.registerInfix(token.LESS, p.parseInfixExpression)
	p.registerInfix(token.LESS_EQUAL, p.parseInfixExpression)
	p.registerInfix(token.GREATER, p.parseInfixExpression)
	p.registerInfix(token.GREATER_EQUAL, p.parseInfixExpression)

	return p
}

func (p *Parser) GetErrors() []error {
	return p.errors
}

func (p *Parser) addError(message string) {
	p.errors = append(p.errors, fmt.Errorf(message))
}

func (p *Parser) advance() {
	p.currentAt++
	p.refreshCurrentToken()
}

func (p *Parser) isAtEnd() bool {
	return p.currentToken.Type == token.EOF || p.currentAt >= len(p.tokens)-1
}

func (p *Parser) refreshCurrentToken() token.Token {
	if len(p.tokens) == 1 {
		return token.Token{Type: token.EOF}
	}
	p.currentToken = p.tokens[p.currentAt]
	return p.currentToken
}

func (p *Parser) nextToken() token.Token {
	if p.isAtEnd() {
		return token.Token{Type: token.EOF}
	}
	return p.tokens[p.currentAt+1]
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn func() ast.Expression) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn func(ast.Expression) ast.Expression) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) peekNextPriority() int {
	if priority, ok := priorityMap[p.nextToken().Type]; ok {
		return priority
	}
	return LOWEST
}

func (p *Parser) peekCurrentPriority() int {
	if priority, ok := priorityMap[p.currentToken.Type]; ok {
		return priority
	}
	return LOWEST
}

func (p *Parser) Parse() (*ast.Program, error) {
	// 字句解析
	if err := p.runScanner(); err != nil {
		return nil, err
	}
	p.refreshCurrentToken()
	program := p.parseProgram()

	if errs := p.errors; len(errs) > 0 {
		var msg strings.Builder
		for _, err := range errs {
			msg.WriteString(err.Error())
		}
		return nil, fmt.Errorf(msg.String())
	}
	return program, nil
}

func (p *Parser) runScanner() error {
	p.scanner.Reset()
	p.scanner.ScanTokens()
	if errs := p.scanner.GetErrors(); len(errs) > 0 {
		var msg strings.Builder
		for _, err := range errs {
			msg.WriteString(err.Error())
		}
		return fmt.Errorf(msg.String())
	}
	p.tokens = p.scanner.Tokens()
	return nil
}

func (p *Parser) parseProgram() *ast.Program {
	program := &ast.Program{}
	for !p.isAtEnd() {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, *stmt)
		}
	}
	return program
}

func (p *Parser) parseStatement() *ast.Statement {
	var stmt ast.Statement
	switch p.tokens[p.currentAt].Type {
	case token.VAR:
		stmt = p.parseVarStatement()
	case token.RETURN:
		stmt = p.parseReturnStatement()
	default:
		stmt = p.parseExpressionStatement()
	}
	return &stmt
}

func (p *Parser) parseEnd() {
	isValidEnd := false
	for p.currentToken.Type == token.SEMICOLON || p.currentToken.Type == token.LINE_BREAK {
		p.advance()
		isValidEnd = true
	}
	if !isValidEnd {
		p.addError(fmt.Sprintf("line %v; expected ';', but got %s", p.currentToken.Line, p.currentToken.RawToken))
	}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	expressionStatement := &ast.ExpressionStatement{}
	expressionStatement.Token = p.currentToken
	expressionStatement.Expression = p.parseExpression(LOWEST)
	p.advance()
	p.parseEnd()
	return expressionStatement
}

// var identifier = expression;
func (p *Parser) parseVarStatement() *ast.VarStatement {
	varStatement := &ast.VarStatement{}
	varStatement.Token = p.currentToken
	p.advance()
	if p.currentToken.Type != token.IDENTIFIER {
		p.addError(fmt.Sprintf("line %v; expected identifier, but got %s", p.currentToken.Line, p.currentToken.RawToken))
		return nil
	}

	varStatement.Name = p.parseIdentifier().(ast.Identifier)

	p.advance()
	if p.currentToken.Type != token.EQUAL {
		p.addError(fmt.Sprintf("line %v; expected '=', but got %s", p.currentToken.Line, p.currentToken.RawToken))
		return nil
	}
	p.advance()
	varStatement.Value = p.parseExpression(LOWEST)
	p.advance()
	p.parseEnd()
	return varStatement
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	returnStatement := &ast.ReturnStatement{}
	returnStatement.Token = p.currentToken
	p.advance()
	returnStatement.ReturnValue = p.parseExpression(LOWEST)
	return returnStatement
}

func (p *Parser) parseExpression(priority int) ast.Expression {
	prefix := p.prefixParseFns[p.currentToken.Type]
	if prefix == nil {
		p.addError(fmt.Sprintf("line %v; no prefix parse function for %s", p.currentToken.Line, p.currentToken.RawToken))
		return nil
	}
	leftExp := prefix()

	for priority < p.peekNextPriority() {
		infix := p.infixParseFns[p.nextToken().Type]
		if infix == nil {
			return leftExp
		}
		p.advance()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.RawToken,
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	return &ast.IntegerLiteral{
		Token: p.currentToken,
		Value: p.currentToken.Literal.(int),
	}
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	return &ast.FloatLiteral{
		Token: p.currentToken,
		Value: p.currentToken.Literal.(float64),
	}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.RawToken,
	}
	p.advance()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.RawToken,
		Left:     left,
	}
	priority := p.peekCurrentPriority()
	p.advance()
	expression.Right = p.parseExpression(priority)
	return expression
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{
		Token: p.currentToken,
		Value: p.currentToken.Type == token.TRUE,
	}
}
