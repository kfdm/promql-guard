package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// https://blog.questionable.services/article/testing-http-handlers-go/

func TestQuery(t *testing.T) {
	targetHandler := Query()

	req, err := http.NewRequest("GET", "/api/v1/query", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	targetHandler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestSeries(t *testing.T) {
	targetHandler := Series()
	req, err := http.NewRequest("GET", "/api/v1/series", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	targetHandler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
