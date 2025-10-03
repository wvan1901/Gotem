package token

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	// Literal
	// NEW_LINE       = "NEW_LINE"
	STRING_LITERAL = "STRING_LITERAL"
	// Delimiters
	AT_SIGN = "@"
	EQUAL   = "="
	// Non functional
	COMMENT = "COMMENT"
	// Labels
	LABEL_NAME  = "LABEL_NAME"
	LABEL_VALUE = "LABEL_VALUE"
	// Other
	HTTP_TEMPLATE = "HTTP_TEMPLATE"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Line    int
}
