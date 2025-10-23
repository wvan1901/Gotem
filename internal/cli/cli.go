package cli

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/textproto"
	"strings"
	"text/template"

	"github.com/wvan1901/Gotem/internal/gohttp/ast"
	"github.com/wvan1901/Gotem/internal/gohttp/lexer"
	"github.com/wvan1901/Gotem/internal/gohttp/parse"
)

type response struct {
	StatusCode int
	Body       string
	// Cookies
	// Headers
	// Time
	// Size
}

func ParseInputIntoProgram(input string) (*ast.Progarm, error) {
	l := lexer.New(input)
	p := parse.New(l)

	program := p.ParseRequests()
	errs := p.Errors()

	if program == nil {
		return nil, errors.New("error while parsing request, parser returned nil")
	}
	if len(errs) != 0 {
		errMsg := ""
		for i, msg := range errs {
			errMsg += fmt.Sprintf("[%d] parser error: %q", i, msg)
		}
		return nil, errors.New("errors while parsing request: " + errMsg)
	}
	return program, nil
}

func ListAstRequests(rs []ast.UserRequest) string {
	nameHeader := "Name"
	descHeader := "Description"
	nameMaxWidth := len(nameHeader)
	descMaxWidth := len(descHeader)
	rows := [][2]string{}
	rows = append(rows, [2]string{nameHeader, descHeader})
	for _, r := range rs {
		rows = append(rows, [2]string{r.Name, r.Description})
		if len(r.Name) > nameMaxWidth {
			nameMaxWidth = len(r.Name)
		}
		if len(r.Description) > descMaxWidth {
			descMaxWidth = len(r.Description)
		}
	}

	listStr := ""
	for _, row := range rows {
		rowStr := row[0] + strings.Repeat(" ", nameMaxWidth+2-len(row[0]))
		rowStr += row[1] + strings.Repeat(" ", descMaxWidth+2-len(row[1]))
		rowStr += "\n"
		listStr += rowStr
	}
	return listStr
}

func MakeRequest(ur *ast.UserRequest, extraHeaders map[string][]string) (*response, error) {
	url, bodyTempl, err := createUrlAndBodyFromTemplate(ur.Url, ur.HttpRequestBody, ur.ExtraLabels)

	headers, bodyStr, err := getHeadersAndBody(bodyTempl)
	if err != nil {
		return nil, fmt.Errorf("MakeRequest: reading headers & body: %w", err)
	}

	httpHeaders := make(map[string][]string)
	if len(headers) > 0 {
		httpHeaders = headers
	}

	// Create request
	req, err := http.NewRequest(ur.HttpMethod, url, bytes.NewBufferString(bodyStr))
	req.Header = httpHeaders
	for h, hv := range extraHeaders {
		for _, v := range hv {
			req.Header.Add(h, v)
		}
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Execute: res: %w", err)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Execute: body: %w", err)
	}

	resp := &response{
		StatusCode: res.StatusCode,
		Body:       string(resBody),
	}

	return resp, nil
}

func getHeadersAndBody(httpReqBody string) (http.Header, string, error) {
	// NOTE: We are parsing as http message so new lines matter!
	// "header: 1\n{"some": "value"}" Will not return a body since due to rfc we need 2 '\n' to return a body we need this string:
	// "header: 1\n\n{"some": "value"}"
	// So we should follow the following format:
	// Headers...
	//
	// Body string...
	reader := bufio.NewReader(strings.NewReader(httpReqBody))

	// Create a Textproto Reader
	tp := textproto.NewReader(reader)

	// Read the MIME-style headers
	headers, err := tp.ReadMIMEHeader()
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, "", fmt.Errorf("getHeadersAndBody: headers: %w", err)
	}

	// Cast the result to http.Header for easier use.
	httpHeaders := http.Header(headers)

	// After the headers, the remaining data is the body.
	bodyStr, err := reader.ReadString(0) // Read until EOF
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, "", fmt.Errorf("MakeRequest: body: %w", err)
	}

	return httpHeaders, bodyStr, nil
}

func GetAstRequest(name string, reqs []ast.UserRequest) (*ast.UserRequest, error) {
	if name == "" {
		if len(reqs) == 1 {
			return &reqs[0], nil
		}
		return nil, errors.New("GetAstRequest: name is empty")
	}
	for _, r := range reqs {
		if r.Name == name {
			return &r, nil
		}
	}

	return nil, errors.New("no request found")
}

func createUrlAndBodyFromTemplate(urlTempl, bodyTempl string, labels map[string]string) (string, string, error) {
	t1, err := template.New("url").Parse(urlTempl)
	if err != nil {
		return "", "", fmt.Errorf("createUrlAndBodyFromTemplate: url: parse: %w", err)
	}
	newUrl := &strings.Builder{}
	err = t1.Execute(newUrl, labels)
	if err != nil {
		return "", "", fmt.Errorf("createUrlAndBodyFromTemplate: url: execute: %w", err)
	}

	t2, err := template.New("body").Parse(bodyTempl)
	if err != nil {
		return "", "", fmt.Errorf("createUrlAndBodyFromTemplate: body: %w", err)
	}
	newBody := &strings.Builder{}
	err = t2.Execute(newBody, labels)
	if err != nil {
		return "", "", fmt.Errorf("createUrlAndBodyFromTemplate: body: execute: %w", err)
	}

	return newUrl.String(), newBody.String(), nil
}
