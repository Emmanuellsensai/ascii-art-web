package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// helper: sends a POST with form data and returns the status code
func postForm(text, banner string) int {
	form := url.Values{"text": {text}, "banner": {banner}}

	// "text" : emmanuel, "banner" : shadow (text=emmanuel&banner=shadow)
	req := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	asciiArtHandler(rec, req)
	return rec.Code
}

// GET / → 200
func TestHome_GET(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	homeHandler(rec, req)
	if rec.Code != 200 {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

// POST / → 400
func TestHome_POST(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	homeHandler(rec, req)
	if rec.Code != 400 {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

// GET /wrong → 404
func TestHome_WrongPath(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/wrong", nil)
	rec := httptest.NewRecorder()
	homeHandler(rec, req)
	if rec.Code != 404 {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

// GET /ascii-art → 400
func TestAscii_GET(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ascii-art", nil)
	rec := httptest.NewRecorder()
	asciiArtHandler(rec, req)
	if rec.Code != 400 {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

// POST valid input → 200
func TestAscii_Valid(t *testing.T) {
	if code := postForm("Hello", "standard"); code != 200 {
		t.Errorf("expected 200, got %d", code)
	}
}

// POST empty text → 400
func TestAscii_EmptyText(t *testing.T) {
	if code := postForm("", "standard"); code != 400 {
		t.Errorf("expected 400, got %d", code)
	}
}

// POST fake banner → 400
func TestAscii_BadBanner(t *testing.T) {
	if code := postForm("Hello", "fakefont"); code != 400 {
		t.Errorf("expected 400, got %d", code)
	}
}
