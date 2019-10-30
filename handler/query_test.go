package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// https://blog.questionable.services/article/testing-http-handlers-go/

func TestQuery(t *testing.T) {
	// Build Reqeust
	req, err := http.NewRequest("GET", "/api/v1/query", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add Test Query
	q := req.URL.Query()
	q.Add("query", "foo")
	req.URL.RawQuery = q.Encode()
	t.Logf("%s", req.URL.String())

	// Test Request
	rr := httptest.NewRecorder()
	targetHandler := Query()
	targetHandler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestSeries(t *testing.T) {
	// Build Reqeust
	req, err := http.NewRequest("GET", "/api/v1/series", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add Test Query
	q := req.URL.Query()
	q.Add("query", "foo")
	req.URL.RawQuery = q.Encode()
	t.Logf("%s", req.URL.String())

	// Test Request
	rr := httptest.NewRecorder()
	targetHandler := Series()
	targetHandler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
