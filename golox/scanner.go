package main

import (
	"strconv"
)

type Scanner struct {
	source  string
	tokens  []Token
	start   int
	current int
	line    int
}

var keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

func (sc *Scanner) ScanTokens() []Token {
	for !sc.isAtEnd() {
		// We are at the beginning of the next lexeme.
		sc.start = sc.current
		sc.scanToken()
	}

	sc.tokens = append(sc.tokens, Token{tokenType: EOF, line: sc.line})
	return sc.tokens
}

func (sc *Scanner) scanToken() {
	b := sc.advance()

	switch b {
	case '(':
		sc.addToken(LEFT_PAREN)
	case ')':
		sc.addToken(RIGHT_PAREN)
	case '{':
		sc.addToken(LEFT_BRACE)
	case '}':
		sc.addToken(RIGHT_BRACE)
	case ',':
		sc.addToken(COMMA)
	case '.':
		sc.addToken(DOT)
	case '-':
		sc.addToken(MINUS)
	case '+':
		sc.addToken(PLUS)
	case ';':
		sc.addToken(SEMICOLON)
	case '*':
		sc.addToken(STAR)

	case '!':
		if sc.match('=') {
			sc.addToken(BANG_EQUAL)
		} else {
			sc.addToken(BANG)
		}
	case '=':
		if sc.match('=') {
			sc.addToken(EQUAL_EQUAL)
		} else {
			sc.addToken(EQUAL)
		}
	case '<':
		if sc.match('=') {
			sc.addToken(LESS_EQUAL)
		} else {
			sc.addToken(LESS)
		}
	case '>':
		if sc.match('=') {
			sc.addToken(GREATER_EQUAL)
		} else {
			sc.addToken(GREATER)
		}
	case '/':
		if sc.match('/') {
			// A comment goes until the end of the line.
			for sc.peek() != '\n' && !sc.isAtEnd() {
				sc.advance()
			}
		} else {
			sc.addToken(SLASH)
		}

	case ' ':
	case '\r':
	case '\t':

	case '\n':
		sc.line++

	case '"':
		sc.handleString()

	default:
		if isDigit(b) {
			sc.handleNumber()
		} else if isAlpha(b) {
			sc.handleIdentifier()
		} else {
			reportError(sc.line, "Unexpected character.")
		}
	}
}

func (sc *Scanner) handleIdentifier() {
	for isAlphaNumeric(sc.peek()) {
		sc.advance()
	}

	text := sc.source[sc.start:sc.current]
	tokenType, ok := keywords[text]
	if !ok {
		tokenType = IDENTIFIER
	}
	sc.addToken(tokenType)
}

func (sc *Scanner) handleNumber() {
	for isDigit(sc.peek()) {
		sc.advance()
	}

	// Look for a fractional part.
	if sc.peek() == '.' && isDigit(sc.peekNext()) {
		// Consume the "."
		sc.advance()

		for isDigit(sc.peek()) {
			sc.advance()
		}
	}

	double, _ := strconv.ParseFloat(sc.source[sc.start:sc.current], 64)
	sc.addTokenLit(NUMBER, double)
}

func (sc *Scanner) handleString() {
	for sc.peek() != '"' && !sc.isAtEnd() {
		if sc.peek() == '\n' {
			sc.line++
		}
		sc.advance()
	}

	if sc.isAtEnd() {
		reportError(sc.line, "Unterminated string.")
		return
	}

	// The closing ".
	sc.advance()

	// Trim the surrounding quotes.
	value := sc.source[(sc.start + 1):(sc.current - 1)]
	sc.addTokenLit(STRING, value)
}

func (sc *Scanner) match(expected byte) bool {
	if sc.isAtEnd() {
		return false
	}
	if sc.source[sc.current] != expected {
		return false
	}

	sc.current++
	return true
}

func (sc *Scanner) peek() byte {
	if sc.isAtEnd() {
		return '\000'
	}
	return sc.source[sc.current]
}

func (sc *Scanner) peekNext() byte {
	if sc.current+1 >= len(sc.source) {
		return '\000'
	}
	return sc.source[sc.current+1]
}

func (sc *Scanner) isAtEnd() bool {
	return sc.current >= len(sc.source)
}

func (sc *Scanner) advance() byte {
	c := sc.source[sc.current]
	sc.current++
	return c
}

func (sc *Scanner) addToken(tokenType TokenType) {
	text := sc.source[sc.start:sc.current]
	sc.tokens = append(sc.tokens, Token{tokenType: tokenType, lexeme: text, line: sc.line})
}

func (sc *Scanner) addTokenLit(tokenType TokenType, literal any) {
	text := sc.source[sc.start:sc.current]
	sc.tokens = append(sc.tokens, Token{tokenType: tokenType, lexeme: text, literal: literal, line: sc.line})
}

func isAlpha(b byte) bool {
	return (b >= 'a' && b <= 'z') ||
		(b >= 'A' && b <= 'Z') ||
		b == '_'
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func isAlphaNumeric(b byte) bool {
	return isAlpha(b) || isDigit(b)
}
