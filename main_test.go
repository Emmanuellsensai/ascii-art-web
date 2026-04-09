package main

import (
	"ascii-art-web/ascii"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// Test: GET / returns 200
func TestHomeHandler_GET(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	homeHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

// Test: POST / returns 400 (only GET allowed)
func TestHomeHandler_POST(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	homeHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

// Test: GET /wrong returns 404
func TestHomeHandler_WrongPath(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/wrong", nil)
	rec := httptest.NewRecorder()
	homeHandler(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

// Test: GET /ascii-art returns 400 (only POST allowed)
func TestAsciiHandler_GET(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ascii-art", nil)
	rec := httptest.NewRecorder()
	asciiArtHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

// Test: POST /ascii-art with empty text returns 400
func TestAsciiHandler_EmptyText(t *testing.T) {
	form := url.Values{"text": {""}, "banner": {"standard"}}
	req := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	asciiArtHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

// Test: POST /ascii-art with invalid banner returns 400
func TestAsciiHandler_InvalidBanner(t *testing.T) {
	form := url.Values{"text": {"hello"}, "banner": {"fakefont"}}
	req := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	asciiArtHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

// Test: POST /ascii-art with valid input returns 200
func TestAsciiHandler_ValidInput(t *testing.T) {
	form := url.Values{"text": {"Hello"}, "banner": {"standard"}}
	req := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	asciiArtHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

// ========== ASCII LOGIC TESTS ==========

// Test: ReadBanner on missing file returns error
func TestReadBanner_MissingFile(t *testing.T) {
	_, err := ascii.ReadBanner("banners/doesnotexist.txt")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

// Test: ReadBanner on valid file returns lines
func TestReadBanner_ValidFile(t *testing.T) {
	lines, err := ascii.ReadBanner("banners/standard.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 95 characters × 9 lines each + 1 initial line = 856 lines minimum
	if len(lines) < 855 {
		t.Errorf("expected at least 855 lines, got %d", len(lines))
	}
}

// Test: BuildAsciiMap produces 95 entries (space through tilde)
func TestBuildAsciiMap(t *testing.T) {
	lines, _ := ascii.ReadBanner("banners/standard.txt")
	m := ascii.BuildAsciiMap(lines)

	if len(m) != 95 {
		t.Errorf("expected 95 characters in map, got %d", len(m))
	}
}

// Test: PrintAscii on empty string returns empty
func TestPrintAscii_Empty(t *testing.T) {
	lines, _ := ascii.ReadBanner("banners/standard.txt")
	m := ascii.BuildAsciiMap(lines)
	result := ascii.PrintAscii("", m)

	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

// Test: PrintAscii on newline input returns just a newline
func TestPrintAscii_Newline(t *testing.T) {
	lines, _ := ascii.ReadBanner("banners/standard.txt")
	m := ascii.BuildAsciiMap(lines)
	result := ascii.PrintAscii("\r\n", m)

	if result != "\n" {
		t.Errorf("expected single newline, got %q", result)
	}
}

// Test: Each banner produces output for simple input
func TestAllBanners(t *testing.T) {
	banners := []string{"standard", "shadow", "thinkertoy"}
	for _, b := range banners {
		lines, err := ascii.ReadBanner("banners/" + b + ".txt")
		if err != nil {
			t.Fatalf("banner %s: %v", b, err)
		}
		m := ascii.BuildAsciiMap(lines)
		result := ascii.PrintAscii("A", m)
		if result == "" {
			t.Errorf("banner %s: expected output for 'A', got empty", b)
		}
	}
}
