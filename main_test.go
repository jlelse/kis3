package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHelloResponseHandler(t *testing.T) {
	t.Run("Hello response", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(HelloResponseHandler)
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
		expected := "Hello from KISSS"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	})
}
