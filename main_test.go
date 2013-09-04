package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestHelloWorld(t *testing.T) {
	record := httptest.NewRecorder()
	req := &http.Request{
		Method: "GET",
		URL:    &url.URL{Path: "/world"},
	}

	hello(record, req)
	if got, want := record.Code, 400; got != want {
		t.Errorf("%s: response code = %d, want %d", "Hello world test", got, want)
	}
	fmt.Print("HI")
}
