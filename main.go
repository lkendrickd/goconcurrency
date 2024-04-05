package main

import (
	"net/http"
	"sync"

	"github.com/lkendrickd/concurrency/urlchecker"
)

// ##### NOTE: Drop in any of the example code from the respective directories
// markdown file example to see the output of the code in the main function. There are also table
// driven tests in the respective directories that can be run with `go test`.

func main() {
	// List of URLs to check
	urls := []string{
		"http://api.open-notify.org/astros.json", // This URL returns the number of astronauts currently in space.
		"https://catfact.ninja/fact",             // This URL returns a random cat fact.
	}

	// Create a new Checker instance
	client := &http.Client{}                    // A client is needed to call the URLs
	wg := &sync.WaitGroup{}                     // A waitgroup is needed to wait for all goroutines to finish
	checker := urlchecker.New(client, urls, wg) // Create a new checker instance

	// Perform the checks
	checker.Check()

	// Retrieve and print the results
	for result := range checker.Results() {
		println(result)
	}
}
