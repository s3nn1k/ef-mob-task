package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestCase struct {
	Name       string
	Url        string
	Method     string
	Body       string
	WantStatus int
	WantRes    string
}

func TestEndpoint(t *testing.T, handler http.Handler, test TestCase) {
	t.Run(test.Name, func(t *testing.T) {
		req, err := http.NewRequest(test.Method, test.Url, bytes.NewBufferString(test.Body))
		if err != nil {
			t.Fatalf("error not expected while creating request: %s", err)
		}

		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		if res.Code != test.WantStatus {
			t.Fatalf("error: status missmatch. Want %v, but got %v", test.WantStatus, res.Code)
		}

		if test.WantRes != "" {
			if test.WantRes != res.Body.String() {
				t.Fatalf("error: body missmatch. Want %s, but got %s", test.WantRes, res.Body.String())
			}
		}
	})
}
