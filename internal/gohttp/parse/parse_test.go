package parse

import (
	"strings"
	"testing"

	"github.com/wvan1901/Gotem/internal/gohttp/ast"
	"github.com/wvan1901/Gotem/internal/gohttp/lexer"
)

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestParseRequest(t *testing.T) {
	input := `
# Some comment!
@Name="request2"
@Description="submit"
POST http://localhost:42069/submit
Host: localhost:42069
Content-Length: 13

hello world!

# Comment 2
@Name="health"
@Description="health-check"
@testLabel="random"
# Comment before request
GET http://localhost:42069/health
Host: localhost:42069

@Name="templ"
@Description="example"
@url="http://localhost:8090"
@header_one = "1"
@header_two = "2"
@type_value = "none"
GET {{.url}}/health
header: {{.header_one}}
header: {{.header_two}}

{"type": "{{.type_value}}"}
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseRequests()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseRequests() returned nil")
	}

	tests := []struct {
		expectedRequest ast.UserRequest
	}{
		{expectedRequest: ast.UserRequest{
			Name:            "request2",
			Description:     "submit",
			HttpMethod:      "POST",
			Url:             "http://localhost:42069/submit",
			HttpRequestBody: "Host: localhost:42069\nContent-Length: 13\n\nhello world!\n",
			ExtraLabels:     nil,
		}},
		{expectedRequest: ast.UserRequest{
			Name:            "health",
			Description:     "health-check",
			HttpMethod:      "GET",
			Url:             "http://localhost:42069/health",
			HttpRequestBody: "Host: localhost:42069\n",
			ExtraLabels:     map[string]string{"testlabel": "random"},
		}},
		{expectedRequest: ast.UserRequest{
			Name:            "templ",
			Description:     "example",
			HttpMethod:      "GET",
			Url:             "{{.url}}/health",
			HttpRequestBody: "header: {{.header_one}}\nheader: {{.header_two}}\n\n{\"type\": \"{{.type_value}}\"}\n",
			ExtraLabels:     map[string]string{"url": "http://localhost:8090", "header_one": "1", "header_two": "2", "type_value": "none"},
		}},
	}

	if len(program.Requests) != len(tests) {
		t.Fatalf("program.Requests does not contain 1 request. got=%d", len(program.Requests))
	}

	for i, tt := range tests {
		req := program.Requests[i]
		if !testRequest(t, req, tt.expectedRequest) {
			return
		}
	}
}

func TestParseMissingRequestErrorRequest(t *testing.T) {
	input := `
# ...
@Name="request2"
@Description="submit"
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseRequests()
	if program == nil {
		t.Fatalf("ParseRequests() returned nil")
	}
	if len(program.Requests) != 0 {
		t.Fatalf("ParseRequests() returned have returned no request")
	}
	errs := p.Errors()
	firstErr := errs[0]

	if !strings.Contains(firstErr, "request is invalid") {
		t.Fatalf("expected invalid request error")
	}
}

func TestParseMissingRequestNameError(t *testing.T) {
	input := `
# ...
@Description="submit"
GET http://localhost:42069/health
Host: localhost:42069
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseRequests()
	if program == nil {
		t.Fatalf("ParseRequests() returned nil")
	}
	if len(program.Requests) != 0 {
		t.Fatalf("ParseRequests() returned have returned no request")
	}
	errs := p.Errors()
	firstErr := errs[0]

	if !strings.Contains(firstErr, "request is invalid") {
		t.Fatalf("expected invalid request error")
	}
}

func testRequest(t *testing.T, s ast.UserRequest, r ast.UserRequest) bool {
	if isValid, problems := s.IsValid(); !isValid {
		t.Errorf("request not valid. problems=%q", problems)
		return false
	}

	isEqual, issue := isRequestEqual(s, r)
	if !isEqual {
		t.Errorf("request doesn't match expected. issue=%q", issue)
		return false
	}
	return true
}

func isRequestEqual(r1, r2 ast.UserRequest) (bool, string) {
	if r1.Name != r2.Name {
		return false, "name"
	}
	if r1.Description != r2.Description {
		return false, "description"
	}
	if r1.HttpMethod != r2.HttpMethod {
		return false, "http Method"
	}
	if r1.Url != r2.Url {
		return false, "http Method"
	}
	if r1.HttpRequestBody != r2.HttpRequestBody {
		return false, "http body"
	}
	if len(r1.ExtraLabels) != len(r2.ExtraLabels) {
		return false, "labels"
	}
	for k, v := range r1.ExtraLabels {
		v2, ok := r2.ExtraLabels[k]
		if !ok {
			return false, "label missing"
		}
		if v2 != v {
			return false, "label value different"
		}
	}
	return true, ""
}
