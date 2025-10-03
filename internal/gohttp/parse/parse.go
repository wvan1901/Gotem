package parse

import (
	"fmt"
	"strings"

	"github.com/wvan1901/Gotem/internal/gohttp/ast"
	"github.com/wvan1901/Gotem/internal/gohttp/lexer"
	"github.com/wvan1901/Gotem/internal/gohttp/token"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Read two tokens so curToken & peekToken are set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseRequests() *ast.Progarm {
	prog := &ast.Progarm{
		Requests: []ast.UserRequest{},
	}

	for !p.curTokenIs(token.EOF) {
		req := p.parseRequest()
		if req != nil {
			prog.Requests = append(prog.Requests, *req)
		}
		p.nextToken()
	}

	return prog
}

func (p *Parser) parseRequest() *ast.UserRequest {
	req := &ast.UserRequest{ExtraLabels: map[string]string{}}
	for !p.curTokenIs(token.EOF) && !p.curTokenIs(token.HTTP_TEMPLATE) {
		switch p.curToken.Type {
		case token.COMMENT:
			p.nextToken()
		case token.AT_SIGN:
			lName, lValue, isValidLabel := p.parseLabel()
			if isValidLabel {
				switch strings.ToLower(lName) {
				case "name":
					req.Name = lValue
				case "description":
					req.Description = lValue
				default:
					req.ExtraLabels[lName] = lValue
				}
			}
		default:
			p.invalidToken(p.curToken)
			return nil
		}
	}
	if p.curTokenIs(token.HTTP_TEMPLATE) {
		requestLine, body := p.parseHttpTemplate()
		httpMethod, url := p.parseRequestLine(requestLine)
		req.HttpMethod = httpMethod
		req.Url = url
		req.HttpRequestBody = body
	}
	valid, problems := req.IsValid()
	if !valid {
		p.addProblems(problems, "request is invalid")
		return nil
	}

	return req
}

func (p *Parser) parseLabel() (string, string, bool) {
	if !p.expectPeek(token.LABEL_NAME) {
		return "", "", false
	}
	name := p.curToken.Literal
	if !p.expectPeek(token.EQUAL) {
		return "", "", false
	}
	if !p.expectPeek(token.LABEL_VALUE) {
		return "", "", false
	}
	val := p.curToken.Literal
	p.nextToken()
	return name, val, true
}

func (p *Parser) parseHttpTemplate() (string, string) {
	splitStr := strings.SplitN(p.curToken.Literal, "\n", 2)
	customRequestLine := splitStr[0]
	httpBody := splitStr[1]

	return customRequestLine, httpBody
}

func (p *Parser) parseRequestLine(reqLine string) (string, string) {
	httpMethod := ""
	url := ""
	splitReq := strings.Split(strings.TrimSpace(reqLine), " ")
	if len(splitReq) != 2 {
		return httpMethod, url
	}
	httpMethod = splitReq[0]
	url = splitReq[1]
	return httpMethod, url
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) invalidToken(t token.Token) {
	msg := fmt.Sprintf("token should not be here, %v found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) addProblems(prob map[string]string, title string) {
	msg := title + ": "
	for k, v := range prob {
		msg += fmt.Sprintf("issue: %s, err: %s,", k, v)
	}
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead, line %d", t, p.peekToken.Type, p.peekToken.Line)
	p.errors = append(p.errors, msg)
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}
