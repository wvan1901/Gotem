package internal_test

import (
	"fmt"
	"github.com/wvan1901/Gotem/internal"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMethodRequests(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.WriteHeader(200)
			_, err := w.Write([]byte("get"))
			if err != nil {
				t.Errorf("unexpected handler err: %s", err.Error())
			}
		case http.MethodPost:
			w.WriteHeader(200)
			_, err := w.Write([]byte("post"))
			if err != nil {
				t.Errorf("unexpected handler err: %s", err.Error())
			}
		case http.MethodDelete:
			w.WriteHeader(200)
			_, err := w.Write([]byte("delete"))
			if err != nil {
				t.Errorf("unexpected handler err: %s", err.Error())
			}
		case http.MethodPut:
			w.WriteHeader(200)
			_, err := w.Write([]byte("put"))
			if err != nil {
				t.Errorf("unexpected handler err: %s", err.Error())
			}
		case http.MethodOptions:
			w.WriteHeader(200)
			_, err := w.Write([]byte("options"))
			if err != nil {
				t.Errorf("unexpected handler err: %s", err.Error())
			}
		default:
			w.WriteHeader(500)
		}
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	tests := []struct {
		reqMethod    string
		expectedBody string
	}{
		{reqMethod: http.MethodGet, expectedBody: "get"},
		{reqMethod: http.MethodPost, expectedBody: "post"},
		{reqMethod: http.MethodDelete, expectedBody: "delete"},
		{reqMethod: http.MethodPut, expectedBody: "put"},
		{reqMethod: http.MethodOptions, expectedBody: "options"},
	}

	expectedStatus := 200

	for _, tt := range tests {
		req := internal.Request{
			Method: tt.reqMethod,
			Url:    server.URL,
		}
		resp, err := req.Execute()
		if err != nil {
			t.Errorf("unexpected err, err: %s", err.Error())
		}
		if resp.StatusCode != expectedStatus {
			t.Errorf("Status: got %q, want %q", resp.StatusCode, expectedStatus)
		}
		if string(resp.Body) != tt.expectedBody {
			t.Errorf("body: got %s, want %s", string(resp.Body), tt.expectedBody)
		}
	}

}

func TestHeaderRequests(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerStr := fmt.Sprint(r.Header.Values("key"))
		w.WriteHeader(200)
		_, err := w.Write([]byte(headerStr))
		if err != nil {
			t.Errorf("unexpected handler err: %s", err.Error())
		}
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	keyStr := "key"

	tests := []struct {
		reqHeaders map[string][]string
	}{
		{reqHeaders: map[string][]string{keyStr: {"val1"}}},
		{reqHeaders: map[string][]string{keyStr: {"val1", "val2"}}},
		{reqHeaders: map[string][]string{keyStr: {"val1", "val2", "val3"}}},
	}

	for _, tt := range tests {
		req := internal.Request{
			Method:  http.MethodGet,
			Url:     server.URL,
			Headers: tt.reqHeaders,
		}
		resp, err := req.Execute()
		if err != nil {
			t.Errorf("unexpected err, err: %s", err.Error())
		}
		if resp.StatusCode != 200 {
			t.Errorf("Status: got %q, want %q", resp.StatusCode, 200)
		}
		headerVal := tt.reqHeaders[keyStr]
		if string(resp.Body) != fmt.Sprint(headerVal) {
			t.Errorf("body: got %s, want %v", string(resp.Body), headerVal)
		}
	}

}

func TestParamRequests(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paramStr := fmt.Sprint(r.URL.Query())
		w.WriteHeader(200)
		_, err := w.Write([]byte(paramStr))
		if err != nil {
			t.Errorf("unexpected handler err: %s", err.Error())
		}
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	keyStr := "key"

	tests := []struct {
		reqParams map[string][]string
	}{
		{reqParams: map[string][]string{keyStr: {"val1"}}},
		{reqParams: map[string][]string{keyStr: {"val1", "val2"}}},
		{reqParams: map[string][]string{keyStr: {"val1", "val2", "val3"}}},
	}

	for _, tt := range tests {
		req := internal.Request{
			Method:  http.MethodGet,
			Url:     server.URL,
			Headers: tt.reqParams,
			Params: internal.Params{
				Values: tt.reqParams,
			},
		}
		resp, err := req.Execute()
		if err != nil {
			t.Errorf("unexpected err, err: %s", err.Error())
		}
		if resp.StatusCode != 200 {
			t.Errorf("Status: got %q, want %q", resp.StatusCode, 200)
		}
		if string(resp.Body) != fmt.Sprint(tt.reqParams) {
			t.Errorf("body: got %s, want %v", string(resp.Body), tt.reqParams)
		}
	}

}
