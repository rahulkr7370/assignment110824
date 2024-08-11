package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestScrape_Success(t *testing.T) {
	// Mock server to return a valid HTML response
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<html><head><title>Test Page</title></head><body>Hello, this is a test page with some content.</body></html>`))
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	ch := make(chan PageData, 1)
	timeout := 2 * time.Second
	go Scrape(server.URL, ch, timeout)

	result := <-ch
	if result.Title != "Test Page" {
		t.Errorf("Expected title 'Test Page', got '%s'", result.Title)
	}

	expectedContent := "Hello, this is a test page with some content."
	if result.Content != expectedContent {
		t.Errorf("Expected content '%s', got '%s'", expectedContent, result.Content)
	}
}

func TestScrape_Timeout(t *testing.T) {
	// Mock server to simulate a long response time
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
		w.WriteHeader(http.StatusOK)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	ch := make(chan PageData, 1)
	timeout := 1 * time.Second
	go Scrape(server.URL, ch, timeout)

	result := <-ch
	if result.Title != "Timeout" {
		t.Errorf("Expected title 'Timeout', got '%s'", result.Title)
	}

	expectedContent := "The request timed out."
	if result.Content != expectedContent {
		t.Errorf("Expected content '%s', got '%s'", expectedContent, result.Content)
	}
}

func TestScrape_Error(t *testing.T) {
	// Unreachable server
	ch := make(chan PageData, 1)
	timeout := 2 * time.Second
	go Scrape("http://unreachable.url", ch, timeout)

	result := <-ch
	if result.Title != "Error" {
		t.Errorf("Expected title 'Error', got '%s'", result.Title)
	}

	if len(result.Content) == 0 {
		t.Errorf("Expected an error message, but got an empty content")
	}
}

func TestScrape_NoContent(t *testing.T) {
	// Mock server with minimal content
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<html><head><title>Empty Page</title></head><body></body></html>`))
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	ch := make(chan PageData, 1)
	timeout := 2 * time.Second
	go Scrape(server.URL, ch, timeout)

	result := <-ch
	if result.Title != "Empty Page" {
		t.Errorf("Expected title 'Empty Page', got '%s'", result.Title)
	}

	expectedContent := "No content found"
	if result.Content != expectedContent {
		t.Errorf("Expected content '%s', got '%s'", expectedContent, result.Content)
	}
}

func TestScrape_PartialContent(t *testing.T) {
	// Mock server to return more than 100 words
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<html><head><title>Long Content Page</title></head><body>` + strings.Repeat("word ", 150) + `</body></html>`))
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	ch := make(chan PageData, 1)
	timeout := 2 * time.Second
	go Scrape(server.URL, ch, timeout)

	result := <-ch
	if result.Title != "Long Content Page" {
		t.Errorf("Expected title 'Long Content Page', got '%s'", result.Title)
	}

	words := strings.Fields(result.Content)
	if len(words) != 100 {
		t.Errorf("Expected 100 words, but got %d", len(words))
	}
}
