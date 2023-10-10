package handler_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hararudoka/shrt/handler"
	"github.com/hararudoka/shrt/service"
)

var (
	// common inputs for all test cases
	googleURL   = "{\"url\":\"google.com\"}"
	XXXXXXShort = "{\"short\":\"XXXXXX\"}\n"
	error404    = "404 page not found\n"
)

type MockStore struct {
	store map[string]string
}

func (s MockStore) Insert(URL, short string) error {
	s.store[short] = URL
	return nil
}

func (s MockStore) ByShort(short string) (string, error) {
	for url, s := range s.store {
		if short == s {
			return url, nil
		}
	}
	return "", fmt.Errorf("404 page not found")
}

func (s MockStore) ByURL(url string) (string, error) {
	if short, ok := s.store[url]; ok {
		return short, nil
	}
	return "", fmt.Errorf("404 page not found")
}

func TestHandler_ServeHTTP(t *testing.T) {
	// create common mock for all test cases
	dbMock := MockStore{store: map[string]string{"google.com": "XXXXXX"}}
	s := service.New(dbMock)
	h := handler.New(*s)

	tests := []struct {
		name     string
		input    string // is main string field
		expected string
		path     string
		method   string
		status   int
	}{
		// URL cases
		{
			name:     "POST google.com -> /api/short = XXXXXX 200",
			input:    googleURL,
			expected: XXXXXXShort,
			path:     "/api/short",
			method:   http.MethodPost,
			status:   http.StatusOK,
		},
		{
			name:     "POST new url -> /api/short = 404",
			input:    "{\"url\":\"absent.url\"}",
			expected: error404,
			path:     "/api/short",
			method:   http.MethodPost,
			status:   http.StatusNotFound,
		},

		// Short cases
		{
			name:     "POST XXXXXX -> /api/url = 200",
			input:    XXXXXXShort,
			expected: "{\"url\":\"google.com\"}\n",
			path:     "/api/url",
			method:   http.MethodPost,
			status:   http.StatusOK,
		},
		{
			name:     "POST new short -> /api/url = 404",
			input:    "{\"short\":\"BBBBBB\"}\n",
			expected: error404,
			path:     "/api/url",
			method:   http.MethodPost,
			status:   http.StatusNotFound,
		},

		// redirect cases
		{
			name:     "GET /XXXXXX = 200",
			expected: "<a href=\"https://google.com\">Found</a>.\n\n",
			path:     "/XXXXXX",
			method:   http.MethodGet,
			status:   http.StatusFound,
		},
		{
			name:     "GET /BBBBBB = 404",
			expected: error404,
			path:     "/BBBBBB",
			method:   http.MethodGet,
			status:   http.StatusNotFound,
		},

		// unusual cases
		{
			name:     "GET BBBBBB -> /api/short = 400",
			input:    "BBBBBB",
			expected: error404,
			path:     "/api/url",
			method:   http.MethodGet,
			status:   http.StatusNotFound,
		},
		{
			name:     "GET BBBBBB -> /api/url = 400",
			input:    "BBBBBB",
			expected: error404,
			path:     "/api/short",
			method:   http.MethodGet,
			status:   http.StatusNotFound,
		},
		{
			name:     "POST BBBBBB -> /api/short = 403",
			input:    "{\"url\":\"BBBBBB\"}\n",
			expected: error404,
			path:     "/api/short",
			method:   http.MethodPost,
			status:   http.StatusNotFound,
		},
		{
			name:     "POST BBBBBB -> /api/url = 404",
			input:    "{\"short\":\"BBBBBB\"}\n",
			expected: error404,
			path:     "/api/url",
			method:   http.MethodPost,
			status:   http.StatusNotFound,
		},

		// // wrong method cases
		// {
		// 	name:     "PUT google.com -> /api/short = 405",
		// 	input:    googleURL,
		// 	expected: "\n",
		// 	path:     "/api/short",
		// 	method:   http.MethodPut,
		// 	status:   http.StatusMethodNotAllowed,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// tt.input -> reader
			reader := bytes.NewReader([]byte(tt.input))

			// reader -> r
			r, err := http.NewRequest(tt.method, tt.path, reader)
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			handler := http.HandlerFunc(h.ServeHTTP)

			// r+w -> request
			handler.ServeHTTP(w, r)

			// check status code
			if status := w.Code; status != tt.status {
				t.Errorf("handler returned wrong status code: got '%v' want '%v'",
					status, tt.status)
			}

			// check response body
			if w.Body.String() != tt.expected {
				t.Errorf("handler returned unexpected body: got '%v' want '%v'",
					w.Body.String(), tt.expected)
			}
		})
	}
}
