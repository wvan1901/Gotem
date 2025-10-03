package lexer

import (
	"net/http"

	"github.com/wvan1901/Gotem/internal/gohttp/token"
)

type Lexer struct {
	input        string
	curPosition  int  // Current position in input. Points to current char
	readPosition int  // Current read position in input. After current char
	curChar      byte //Current char under examination
	curLine      int  // Current Line curChar is on
	prevToken    token.Token
}

func New(input string) *Lexer {
	l := &Lexer{input: input, prevToken: token.Token{}}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.curChar = 0
	} else {
		l.curChar = l.input[l.readPosition]
	}
	l.curPosition = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() token.Token {
	tok := token.Token{}

	l.skipNewLine()
	l.skipWhitespace()

	switch l.curChar {
	case '#':
		tok.Literal = l.readComment()
		tok.Type = token.COMMENT
		tok.Line = l.curLine
		// NOTE: Omit comment from prev token so we can set the semi colons properly
		// We could solve this by ommiting comments from the lexer
		return tok

	case '@':
		tok = newToken(token.AT_SIGN, l.curChar, l.curLine)
	case '=':
		tok = newToken(token.EQUAL, l.curChar, l.curLine)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
		tok.Line = l.curLine
	default:
		if l.peekCurUpperCaseWordIsHttpVerb() {
			// Read char until we reach EOF || ('\n' && ('@' || '#' ))
			tok.Literal += l.readHttpRequest()
			tok.Type = token.HTTP_TEMPLATE
			tok.Line = l.curLine
			l.prevToken = tok
			return tok
		} else if isAlphaNumeric(l.curChar) {
			switch l.prevToken.Type {
			case token.AT_SIGN: // This means string is a label name
				tok.Literal += l.readAlphaNumeric()
				tok.Type = token.LABEL_NAME
				tok.Line = l.curLine
				l.prevToken = tok
				return tok
			case token.EQUAL: // String should be a label value
				tok.Literal += l.readAlphaNumeric()
				tok.Type = token.LABEL_VALUE
				tok.Line = l.curLine
				l.prevToken = tok
				return tok
			default:
				tok = newToken(token.ILLEGAL, l.curChar, l.curLine)
			}
		} else {
			tok = newToken(token.ILLEGAL, l.curChar, l.curLine)
		}
	}

	l.readChar()
	l.prevToken = tok
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.curChar == ' ' || l.curChar == '\t' || l.curChar == '\r' {
		l.readChar()
	}
}

func (l *Lexer) skipNewLine() {
	for l.curChar == '\n' {
		l.curLine += 1
		l.readChar()
	}
}

func (l *Lexer) readComment() string {
	position := l.curPosition
	for !isNewLineOrEof(l.curChar) {
		l.readChar()
	}
	return l.input[position:l.curPosition]
}

func (l *Lexer) readAlphaNumeric() string {
	position := l.curPosition
	for isAlphaNumeric(l.curChar) {
		l.readChar()
	}
	return l.input[position:l.curPosition]
}

func (l *Lexer) readHttpRequest() string {
	position := l.curPosition
	curChar, nextChar := l.curChar, l.peekChar()
	for !(curChar == 0) && !(curChar == '\n' && (nextChar == '@' || nextChar == '#')) {
		if curChar == '\n' {
			l.curLine += 1
		}
		l.readChar()
		curChar, nextChar = l.curChar, l.peekChar()
	}
	return l.input[position:l.curPosition]
}

func (l *Lexer) peekCurUpperCaseWordIsHttpVerb() bool {
	endPos := l.curPosition
	for endPos < len(l.input) && isUpperCaseLetter(l.input[endPos]) {
		endPos++
	}

	wordLen := endPos - l.curPosition
	if wordLen < 3 || wordLen > 6 {
		return false
	}

	curUpperWord := l.input[l.curPosition:endPos]

	return curUpperWord == http.MethodGet || curUpperWord == http.MethodPost || curUpperWord == http.MethodDelete || curUpperWord == http.MethodPut || curUpperWord == http.MethodPatch
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch == '-'
}

func isUpperCaseLetter(ch byte) bool {
	return 'A' <= ch && ch <= 'Z'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isAlphaNumeric(ch byte) bool {
	return isLetter(ch) || isDigit(ch)
}

func newToken(tokenType token.TokenType, ch byte, line int) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch), Line: line}
}

func isNewLineOrEof(ch byte) bool {
	return ch == '\n' || ch == 0
}
