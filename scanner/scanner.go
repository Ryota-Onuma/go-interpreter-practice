package scanner

import (
	"fmt"
	"lox-by-go/token"
	"strconv"
	"strings"
)

type Scanner struct {
	source    string
	tokens    []token.Token
	start     int
	currentAt int
	line      int
	errors    []error
}

func NewScanner(source string) *Scanner {
	return &Scanner{source: source, line: 1}
}

func (s *Scanner) Reset() string {
	s.tokens = []token.Token{}
	s.start = 0
	s.currentAt = 0
	s.line = 1
	s.errors = []error{}
	return s.source
}

func (s *Scanner) SetSource(source string) {
	s.source = source
}

func (s *Scanner) GetErrors() []error {
	return s.errors
}

func (s *Scanner) addError(message string) {
	s.errors = append(s.errors, fmt.Errorf("line %d: %s", s.line, message))
}

func (s *Scanner) Tokens() []token.Token {
	return s.tokens
}

func (s *Scanner) ScanTokens() {
	for {
		s.start = s.currentAt
		s.scanToken()
		if s.isAtEnd() {
			break
		}
		s.advance()
	}
	s.tokens = append(s.tokens, token.Token{Type: token.LINE_BREAK}, token.Token{Type: token.EOF})
}

func (s *Scanner) isAtEnd() bool {
	return s.currentAt >= len([]rune(s.source))-1
}

func (s *Scanner) scanToken() {
	c := s.current()
	if s.shouldSkip(c) {
		return
	}

	switch c {
	case '(':
		t := s.createToken(token.LEFT_PAREN)
		s.addToken(t)
	case ')':
		t := s.createToken(token.RIGHT_PAREN)
		s.addToken(t)
	case '{':
		t := s.createToken(token.LEFT_BRACE)
		s.addToken(t)
	case '}':
		t := s.createToken(token.RIGHT_BRACE)
		s.addToken(t)
	case ',':
		t := s.createToken(token.COMMA)
		s.addToken(t)
	case '.':
		t := s.createToken(token.DOT)
		s.addToken(t)
	case '-':
		t := s.createToken(token.MINUS)
		s.addToken(t)
	case '+':
		t := s.createToken(token.PLUS)
		s.addToken(t)
	case ';':
		t := s.createToken(token.SEMICOLON)
		s.addToken(t)
	case '*':
		t := s.createToken(token.STAR)
		s.addToken(t)
	case '!':
		if s.peekNext() == '=' {
			s.advance()
			t := s.createToken(token.NOT_EQUAL)
			s.addToken(t)
		} else {
			t := s.createToken(token.BANG)
			s.addToken(t)
		}
	case '=':
		if s.peekNext() == '=' {
			s.advance()
			t := s.createToken(token.EQUAL_EQUAL)
			s.addToken(t)
		} else {
			t := s.createToken(token.EQUAL)
			s.addToken(t)
		}
	case '<':
		if s.peekNext() == '=' {
			s.advance()
			t := s.createToken(token.LESS_EQUAL)
			s.addToken(t)
		} else {
			t := s.createToken(token.LESS)
			s.addToken(t)
		}
	case '>':
		if s.peekNext() == '=' {
			s.advance()
			t := s.createToken(token.GREATER_EQUAL)
			s.addToken(t)
		} else {
			t := s.createToken(token.GREATER)
			s.addToken(t)
		}
	case '/':
		if s.peekNext() == '/' {
			s.advance()
			for s.peekNext() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			t := s.createToken(token.SLASH)
			s.addToken(t)
		}
	case '"':
		s.createString()
	case '\n':
		t := s.createToken(token.LINE_BREAK)
		s.addToken(t)
		s.line++
	default:
		if isDigit(c) {
			s.createNumber()
		} else if isAlphabet(c) {
			s.identifier()
		} else {
			s.addError(fmt.Sprintf("Unexpected character: %s", string(c)))
		}
	}
	return
}

func (s *Scanner) current() rune {
	return []rune(s.source)[s.currentAt]
}

func (s *Scanner) advance() rune {
	s.currentAt++
	return []rune(s.source)[s.currentAt]
}

func (s *Scanner) peekNext() rune {
	if s.isAtEnd() {
		return '\x00'
	}
	return []rune(s.source)[s.currentAt+1]
}

func (s *Scanner) peekNextNext() rune {
	if s.currentAt+1 >= len([]rune(s.source)) {
		return '\x00'
	}
	return []rune(s.source)[s.currentAt+2]
}

func (s *Scanner) createToken(tokenType token.TokenType) token.Token {
	return token.Token{
		Type:     tokenType,
		RawToken: "",
		Literal:  nil,
		Line:     s.line,
	}
}

func (s *Scanner) addToken(token token.Token) {
	if token.RawToken == "" {
		text := []rune(s.source)[s.start : s.currentAt+1]
		token.RawToken = string(text)
	}

	s.tokens = append(s.tokens, token)
}

func (s *Scanner) shouldSkip(c rune) bool {
	return c == ' ' || c == '\r' || c == '\t'
}

// 数字かどうかを判定する
func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

// アルファベット or アンダースコアかどうかを判定する
func isAlphabet(c rune) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '_'
}

func (s *Scanner) createNumber() {
	for isDigit(s.peekNext()) {
		s.advance()
	}

	// floatの場合
	if s.peekNext() == '.' && isDigit(s.peekNextNext()) {
		s.advance()
		for isDigit(s.peekNext()) {
			s.advance()
		}
		floatLiteral, err := strconv.ParseFloat(string([]rune(s.source)[s.start:s.currentAt+1]), 64)
		if err != nil {
			s.addError(err.Error())
			return
		}
		t := s.createToken(token.FLOAT)
		t.Literal = floatLiteral
		s.addToken(t)
		return
	}

	// intの場合
	intLiteral, err := strconv.Atoi(string([]rune(s.source)[s.start : s.currentAt+1]))
	if err != nil {
		s.addError(err.Error())
		return
	}
	t := s.createToken(token.INTEGER)
	t.Literal = intLiteral
	s.addToken(t)
}

func (s *Scanner) createString() {
	for s.peekNext() != '"' && !s.isAtEnd() {
		if s.peekNext() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.addError("Unterminated string.")
		return
	}

	s.advance()
	t := s.createToken(token.STRING)

	textLiteral := string([]rune(s.source)[s.start+1 : s.currentAt+1])
	// \"を消す
	textLiteral = strings.Replace(textLiteral, `"`, "", -1)
	t.Literal = textLiteral
	s.addToken(t)
}

func (s *Scanner) identifier() {
	for isAlphabet(s.peekNext()) || isDigit(s.peekNext()) {
		s.advance()
	}
	textLiteral := string([]rune(s.source)[s.start : s.currentAt+1])

	// 予約語かどうかを判定する
	tokenType, ok := keywords[textLiteral]
	if !ok {
		t := s.createToken(token.IDENTIFIER)
		t.Literal = textLiteral
		s.addToken(t)
		return
	}

	if tokenType == token.TRUE || tokenType == token.FALSE {
		t := s.createToken(tokenType)
		t.Literal = textLiteral == "true"
		s.addToken(t)
		return
	}

	t := s.createToken(tokenType)
	s.addToken(t)
}
