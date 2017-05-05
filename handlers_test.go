package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type fakeclock struct {
	Seconds int `json:"uptime"`
}

func (f *fakeclock) elapsed() time.Duration {
	duration, _ := time.ParseDuration(fmt.Sprintf("%ds", f.Seconds))
	return duration
}

func TestIndex(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatalf("request: %v", err)
	}

	res := httptest.NewRecorder()
	handler := http.HandlerFunc(indexHandler)
	handler.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusOK {
		t.Errorf("index returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestHTML(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/html", nil)
	if err != nil {
		t.Fatalf("request: %v", err)
	}

	res := httptest.NewRecorder()
	handler := upHandler{&fakeclock{100}, htmlHandler}
	handler.ServeHTTP(res, req)

	expected := "<p>Up 1m40s</p>"
	result := res.Body.String()
	if result != expected {
		t.Errorf("returned unexpected html: got %v want %v", result, expected)
	}
}

func TestJSON(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/json", nil)
	if err != nil {
		t.Fatalf("request: %v", err)
	}

	res := httptest.NewRecorder()
	handler := upHandler{&fakeclock{100}, jsonHandler}
	handler.ServeHTTP(res, req)

	expected := `{"uptime":100}`
	result := res.Body.String()
	if result != expected {
		t.Errorf("returned unexpected json: got %v want %v", result, expected)
	}
}
