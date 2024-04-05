package urlchecker_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/lkendrickd/concurrency/urlchecker"
)

// TestChecker_Check uses table-driven tests to verify the Check method.
func TestChecker(t *testing.T) {
	// Define the test cases
	tests := []struct {
		name       string
		urls       []string
		wantResult []string
	}{
		{
			name: "single URL success",
			urls: []string{"/success"},
			wantResult: []string{
				"/success: 200 OK",
			},
		},
		{
			name: "multiple URLs with mixed results",
			urls: []string{"/success", "/not-found"},
			wantResult: []string{
				"/success: 200 OK",
				"/not-found: 404 Not Found",
			},
		},
	}

	// create the mock http server as for the unit tests we don't need to make real requests
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/success":
			w.WriteHeader(http.StatusOK)
		case "/not-found":
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	// Run the tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			urls := make([]string, len(tc.urls))
			for i, path := range tc.urls {
				urls[i] = fmt.Sprintf("%s%s", server.URL, path)
			}

			client := server.Client()
			client.Timeout = 2 * time.Second
			wg := &sync.WaitGroup{}
			checker := urlchecker.New(client, urls, wg)

			checker.Check()

			var gotResults []string
			for result := range checker.Results() {
				gotResults = append(gotResults, result)
			}

			// Check the results
			if len(gotResults) != len(tc.wantResult) {
				t.Errorf("Expected %d results, got %d", len(tc.wantResult), len(gotResults))
			}

			for _, want := range tc.wantResult {
				found := false
				for _, got := range gotResults {
					if got == fmt.Sprintf("%s%s", server.URL, want) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected result %s not found in results", want)
				}
			}
		})
	}
}

func TestChecker_Check_ErrorHandling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Immediately close the connection to simulate an error condition
		if r.URL.Path == "/error" {
			conn, _, _ := w.(http.Hijacker).Hijack()
			conn.Close()
			return
		}
	}))
	defer server.Close()

	client := server.Client()
	client.Timeout = 1 * time.Second
	wg := &sync.WaitGroup{}

	// Test with a URL that will trigger the error handling path
	urls := []string{server.URL + "/error"}
	checker := urlchecker.New(client, urls, wg)

	checker.Check()

	// Attempt to read results
	var gotResults []string
	for result := range checker.Results() {
		gotResults = append(gotResults, result)
	}

	if len(gotResults) != 1 || !strings.Contains(gotResults[0], "ERROR") {
		t.Errorf("Expected an error result, got: %v", gotResults)
	}
}
