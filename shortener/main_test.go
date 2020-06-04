package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var sampleRedirects = map[string]string{
	"/test": "http://whatever/test",
}

func TestBasicRedirectFailure(t *testing.T) {
	h := LookupHandler(LookupMap(sampleRedirects), notFoundHandler)
	req := httptest.NewRequest(http.MethodGet, "http://test.com/foo", nil)
	w := httptest.NewRecorder()
	h(w, req)

	if w.Result().StatusCode != 404 {
		t.Errorf("expected not found, got %d", w.Result().StatusCode)
	}
}

func TestBasicRedirectSuccess(t *testing.T) {
	h := LookupHandler(LookupMap(sampleRedirects), notFoundHandler)
	req := httptest.NewRequest(http.MethodGet, "http://test.com/test", nil)
	w := httptest.NewRecorder()
	h(w, req)

	if w.Result().StatusCode != http.StatusTemporaryRedirect {
		t.Errorf("expected success, got %d", w.Result().StatusCode)
	}

	if loc := w.Header().Get("Location"); loc != sampleRedirects["/test"] {
		t.Errorf("expected redirect, got %s", loc)
	}
}

func TestFallback(t *testing.T) {
	called := 0
	fb := func(w http.ResponseWriter, r *http.Request) {
		called++
	}
	h := LookupHandler(LookupMap(map[string]string{}), fb)
	req := httptest.NewRequest(http.MethodGet, "http://test.com/test", nil)
	w := httptest.NewRecorder()
	h(w, req)

	if called == 0 {
		t.Errorf("expected fallback handler to be called but it wasn't")
	}
}
