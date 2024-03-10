package parser

import (
	"fmt"
	"go-interpreter-practice/ast"
	"go-interpreter-practice/scanner"
	"go-interpreter-practice/token"
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
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LEFT_PAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUN, p.parseFuncExpression)

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
	p.registerInfix(token.LEFT_PAREN, p.parseCallExpression)

	return p
}

func (p *Parser) GetErrors() []error {
	return p.errors
}

func (p *Parser) addError(message string) {
	p.errors = append(p.errors, fmt.Errorf(message))
}

func (p *Parser) advance() {
	if p.currentAt >= len(p.tokens)-1 {
		return
	}
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
	p.skipMeaningless()
	return &stmt
}

func (p *Parser) skipMeaningless() {
	if p.nextToken().Type != token.LINE_BREAK && p.nextToken().Type != token.SEMICOLON && p.currentToken.Type != token.LINE_BREAK {
		return
	}
	for p.nextToken().Type == token.LINE_BREAK || p.nextToken().Type == token.SEMICOLON {
		p.advance()
	}
	p.advance()
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	expressionStatement := &ast.ExpressionStatement{}
	expressionStatement.Token = p.currentToken
	expressionStatement.Expression = p.parseExpression(LOWEST)
	p.advance()
	// p.parseEnd()
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

	name := p.parseIdentifier().(*ast.Identifier)
	varStatement.Name = name

	p.advance()
	if p.currentToken.Type != token.EQUAL {
		p.addError(fmt.Sprintf("line %v; expected '=', but got %s", p.currentToken.Line, p.currentToken.RawToken))
		return nil
	}
	p.advance()
	varStatement.Value = p.parseExpression(LOWEST)
	p.advance()
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
	return &ast.Identifier{
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

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{
		Token: p.currentToken,
		Value: p.currentToken.Literal.(string),
	}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.advance()
	expression := p.parseExpression(LOWEST)
	if p.nextToken().Type != token.RIGHT_PAREN {
		p.addError(fmt.Sprintf("line %v; expected ')', but got %s", p.currentToken.Line, p.currentToken.RawToken))
		return nil
	}
	p.advance()
	return expression
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{
		Token: p.currentToken,
	}
	p.advance() // ( を消費
	if p.currentToken.Type != token.LEFT_PAREN {
		p.addError(fmt.Sprintf("line %v; expected '(', but got %s", p.currentToken.Line, p.currentToken.RawToken))
		return nil
	}
	p.advance() // 条件式の最初のやつを消費
	expression.Condition = p.parseExpression(LOWEST)
	p.advance() // ) を消費
	if p.currentToken.Type != token.RIGHT_PAREN {
		p.addError(fmt.Sprintf("line %v; expected ')', but got %s", p.currentToken.Line, p.currentToken.RawToken))
		return nil
	}
	p.advance() // { を消費
	if p.currentToken.Type != token.LEFT_BRACE {
		p.addError(fmt.Sprintf("line %v; expected '{', but got %s", p.currentToken.Line, p.currentToken.RawToken))
		return nil
	}
	expression.Consequence = p.parseBlockStatement()
	p.advance() // } の次のやつを消費　else句かもしれないし、そうじゃないかもしれない

	if p.currentToken.Type == token.ELSE {
		p.advance() // { を消費
		if p.currentToken.Type != token.LEFT_BRACE {
			p.addError(fmt.Sprintf("line %v; expected '{', but got %s", p.currentToken.Line, p.currentToken.RawToken))
			return nil
		}
		expression.Alternative = p.parseBlockStatement()
		p.skipMeaningless()
		// p.advance() // } の次のやつを消費
	}
	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	blockStatement := &ast.BlockStatement{
		Token: p.currentToken,
	}
	blockStatement.Statements = []ast.Statement{}
	p.advance()
	for p.currentToken.Type != token.RIGHT_BRACE && p.currentToken.Type != token.EOF {
		p.skipMeaningless() // 改行じゃなくなるまで改行を消費
		// ボディが空の場合
		if p.currentToken.Type == token.RIGHT_BRACE || p.currentToken.Type == token.EOF {
			p.errors = append(p.errors, fmt.Errorf("line %v: empty if body", p.currentToken.Line))
			break
		}
		stmt := p.parseStatement()
		if stmt != nil {
			blockStatement.Statements = append(blockStatement.Statements, *stmt)
		}
	}
	return blockStatement
}

func (p *Parser) parseFuncExpression() ast.Expression {
	expression := &ast.FunctionExpression{
		Token: p.currentToken,
	}
	p.advance()
	if p.currentToken.Type == token.IDENTIFIER {
		expression.Name = p.parseIdentifier().(*ast.Identifier)
		p.advance()
	}

	if p.currentToken.Type != token.LEFT_PAREN {
		p.addError(fmt.Sprintf("line %v; expected '(', but got %s", p.currentToken.Line, p.currentToken.RawToken))
		return nil
	}
	expression.Parameters = p.parseFuncParameters()
	p.advance()
	if p.currentToken.Type != token.LEFT_BRACE {
		p.addError(fmt.Sprintf("line %v; expected '{', but got %s", p.currentToken.Line, p.currentToken.RawToken))
		return nil
	}
	expression.Body = p.parseBlockStatement()
	return expression
}

func (p *Parser) parseFuncParameters() []*ast.Identifier {
	parameters := []*ast.Identifier{}
	if p.nextToken().Type == token.RIGHT_PAREN {
		p.advance()
		return parameters
	}
	p.advance()
	for p.currentToken.Type != token.RIGHT_PAREN && p.currentToken.Type != token.EOF {
		p.skipMeaningless()
		if p.currentToken.Type != token.IDENTIFIER {
			p.addError(fmt.Sprintf("line %v; expected identifier, but got %s", p.currentToken.Line, p.currentToken.RawToken))
			return nil
		}
		identifier := p.parseIdentifier().(*ast.Identifier)
		parameters = append(parameters, identifier)
		p.advance()
		if p.currentToken.Type != token.COMMA && p.currentToken.Type != token.RIGHT_PAREN {
			p.addError(fmt.Sprintf("line %v; expected ',' or ')', but got %s", p.currentToken.Line, p.currentToken.RawToken))
			return nil
		}
		if p.currentToken.Type == token.COMMA {
			p.advance()
		}
	}
	return parameters
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	expression := &ast.CallExpression{
		Token:    p.currentToken,
		Function: &function,
	}
	expression.Arguments = p.parseCallArguments()
	if p.currentToken.Type != token.RIGHT_PAREN {
		p.addError(fmt.Sprintf("line %v; expected ')', but got %s", p.currentToken.Line, p.currentToken.RawToken))
		return nil
	}
	return expression
}

func (p *Parser) parseCallArguments() []*ast.Expression {
	args := []*ast.Expression{}
	if p.nextToken().Type == token.RIGHT_PAREN {
		p.advance()
		return args
	}
	p.advance()
	for p.currentToken.Type != token.RIGHT_PAREN && p.currentToken.Type != token.EOF {
		p.skipMeaningless()
		arg := p.parseExpression(LOWEST)
		args = append(args, &arg)
		p.advance()
		if p.currentToken.Type != token.COMMA && p.currentToken.Type != token.RIGHT_PAREN {
			p.addError(fmt.Sprintf("line %v; expected ',' or ')', but got %s", p.currentToken.Line, p.currentToken.RawToken))
			return nil
		}
		if p.currentToken.Type == token.COMMA {
			p.advance()
		}
	}
	return args
}
