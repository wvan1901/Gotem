package lexer

import (
	"testing"

	"github.com/wvan1901/Gotem/internal/gohttp/token"
)

func TestNextToken(t *testing.T) {
	input := `
# Some comment!
@Name=request2
@Description=submit
POST http://localhost:42069/submit
Host: localhost:42069
Content-Length: 13

hello world!

# Comment 2
@Name=health
@Description=health-check
# Comment before request
GET http://localhost:42069/health
Host: localhost:42069
`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedLine    int
	}{
		{token.COMMENT, "# Some comment!", 1},
		{token.AT_SIGN, "@", 2},
		{token.LABEL_NAME, "Name", 2},
		{token.EQUAL, "=", 2},
		{token.LABEL_VALUE, "request2", 2},
		{token.AT_SIGN, "@", 3},
		{token.LABEL_NAME, "Description", 3},
		{token.EQUAL, "=", 3},
		{token.LABEL_VALUE, "submit", 3},
		{token.HTTP_TEMPLATE, "POST http://localhost:42069/submit\nHost: localhost:42069\nContent-Length: 13\n\nhello world!\n", 9},
		{token.COMMENT, "# Comment 2", 10},
		{token.AT_SIGN, "@", 11},
		{token.LABEL_NAME, "Name", 11},
		{token.EQUAL, "=", 11},
		{token.LABEL_VALUE, "health", 11},
		{token.AT_SIGN, "@", 12},
		{token.LABEL_NAME, "Description", 12},
		{token.EQUAL, "=", 12},
		{token.LABEL_VALUE, "health-check", 12},
		{token.COMMENT, "# Comment before request", 13},
		{token.HTTP_TEMPLATE, "GET http://localhost:42069/health\nHost: localhost:42069\n", 16},
		{token.EOF, "", 16},
	}
	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q line=%d",
				i, tt.expectedType, tok.Type, tok.Line)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q line=%d",
				i, tt.expectedLiteral, tok.Literal, tok.Line)
		}
		if tok.Line != tt.expectedLine {
			t.Fatalf("tests[%d] - line wrong. expected=%d, got=%d",
				i, tt.expectedLine, tok.Line)
		}
	}
}
