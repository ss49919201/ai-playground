package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIEndpoints(t *testing.T) {
	// TODO: implement test
	t.Skip()

	server := NewServer()

	ts := httptest.NewServer(server)
	defer ts.Close()

	cases := []struct {
		name   string
		method string
		path   string
		want   int
	}{
		{"signup", "POST", "/api/v1/signup", 200},
		{"signin", "POST", "/api/v1/signin", 200},
		{"logout", "POST", "/api/v1/logout", 200},
		{"threads GET", "GET", "/api/v1/threads/123", 200},
		{"threads GET", "GET", "/api/v1/threads", 200},
		{"threads POST", "POST", "/api/v1/threads", 200},
	}
	client := &http.Client{}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, err := http.NewRequest(c.method, ts.URL+c.path, nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != c.want {
				t.Errorf("unexpected status: got %d, want %d", resp.StatusCode, c.want)
			}
		})
	}
}
