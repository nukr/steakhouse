package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	server     *httptest.Server
	reader     io.Reader
	healthzURL string
)

func init() {
	server = httptest.NewServer(CreateRouter())
	healthzURL = fmt.Sprintf("%s/healthz", server.URL)
}

func TestIndex(t *testing.T) {
	res, err := http.Get(server.URL)
	if err != nil {
		t.Error("fail")
	}
	if res.StatusCode != 200 {
		t.Errorf("fail %d", res.StatusCode)
	}
}

func TestHealthCheck(t *testing.T) {
	res, _ := http.Get(healthzURL)
	if res.StatusCode != 200 {
		t.Errorf("fail %d", res.StatusCode)
	}
	expected := "no-cache, no-store, must-revalidate"
	header := res.Header.Get("Cache-Control")
	if header != expected {
		t.Errorf("we expected %s but got %s", expected, header)
	}

	expected = "no-cache"
	header = res.Header.Get("Pragma")
	if header != expected {
		t.Errorf("we expected %s but got %s", expected, header)
	}

	expected = "0"
	header = res.Header.Get("Expires")
	if header != expected {
		t.Errorf("we expected %s but got %s", expected, header)
	}
}
