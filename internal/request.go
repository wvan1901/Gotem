package internal

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
)

type Request struct {
	Name    string
	Method  string
	Url     string
	Headers map[string][]string
	Body    Body
	Params  Params
}

type Body struct {
	Type     string // Opts: none, form-data, x-www-form-urlencoded, raw, binary
	FormData map[string][]string
	Raw      string
	Binary   []byte
}

type Params struct {
	Values map[string][]string
}

type Response struct {
	StatusCode int
	Body       []byte
	// Cookies
	// Headers
	// Time
	// Size
}

func (r *Request) IsValid() error {
	if r.Method == "" {
		return errors.New("missing http method")
	}
	validMethods := []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodOptions}
	if !slices.Contains(validMethods, r.Method) {
		return errors.New("invalid http method")
	}

	if r.Url == "" {
		return errors.New("missing url")
	}

	return nil
}

func (r *Request) Execute() (*Response, error) {
	if err := r.IsValid(); err != nil {
		return nil, fmt.Errorf("Execute: invalid: %w", err)
	}

	params := r.Params.GetParams()
	fullUrl := generateUrl(r.Url, params)

	req, err := http.NewRequest(r.Method, fullUrl, r.Body.GetBody())
	if err != nil {
		return nil, fmt.Errorf("Execute: req: %w", err)
	}

	for k, values := range r.Headers {
		for _, v := range values {
			req.Header.Add(k, v)
		}
	}

	for k, v := range r.Body.GetHeaders() {
		req.Header.Add(k, v)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Execute: res: %w", err)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Execute: body: %w", err)
	}

	resp := Response{
		StatusCode: res.StatusCode,
		Body:       resBody,
	}

	return &resp, nil
}

func (b *Body) GetBody() io.Reader {
	switch b.Type {
	case "", "none":
		return nil
	case "form-data", "x-www-form-urlencoded":
		form := url.Values{}
		for key, values := range b.FormData {
			for _, value := range values {
				form.Add(key, value)
			}
		}
		return strings.NewReader(form.Encode())
	case "raw":
		return strings.NewReader(b.Raw)
	case "binary":
		//return bytes.NewReader(b.Binary)
		return nil
	default:
		return nil
	}
}

func (b *Body) GetHeaders() map[string]string {
	switch b.Type {
	case "", "none", "raw", "binary":
		return nil
	case "form-data":
		return map[string]string{"Content-Type": "multipart/form-data"}
	case "x-www-form-urlencoded":
		return map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	default:
		return nil
	}
}

func (p *Params) GetParams() *url.Values {
	if len(p.Values) == 0 {
		return nil
	}
	form := url.Values{}
	for key, values := range p.Values {
		for _, value := range values {
			form.Add(key, value)
		}
	}

	return &form
}

func generateUrl(url string, p *url.Values) string {
	if p == nil {
		return url
	}
	return url + "?" + p.Encode()
}

func GetRequest(name string, reqs []Request) (*Request, error) {
	if name == "" {
		return nil, errors.New("GetRequest: name is empty")
	}
	for _, r := range reqs {
		if r.Name == name {
			return &r, nil
		}
	}

	return nil, errors.New("no request found")
}
