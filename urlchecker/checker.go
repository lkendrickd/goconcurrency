package urlchecker

import (
	"fmt"
	"net/http"
	"sync"
)

/*
Note the url checking is just a trivial example of some unit of work.
The actual work could be anything that can be done concurrently.
Operations such as database queries, file I/O, sending notifications, etc.

Key Concepts:
	- Waitgroups - WaitGroups are used to wait for a collection of goroutines to finish.
	- Channels - Channels are used to communicate between goroutines or gather results.
	- HTTP Client - The HTTP client is used to send requests to URLs.
	- Goroutines - Goroutines are used to check the status of each URL concurrently.
*/

// Checker is a struct that contains a list of URLs to check in a concurrent manner.
// client is an HTTP client used to send GET requests.
// urls is a list of URLs to fetch.
// wg is a WaitGroup used to wait for all goroutines to finish.
// results is a channel used to store the results of the fetch operation.
type Checker struct {
	client  *http.Client
	urls    []string
	wg      *sync.WaitGroup
	results chan string
}

// Check checks status for all stored URLs concurrently.
func (c *Checker) Check() {
	// the channel is the same size as the number of URLs to check to avoid blocking.
	c.results = make(chan string, len(c.urls))

	// Start a goroutine for each URL to check.
	for _, url := range c.urls {
		// Increment the WaitGroup counter this is done before starting the goroutine.
		c.wg.Add(1)
		// Start a goroutine to check the status of the URL.
		go c.checkStatus(url)
	}

	// Wait for all checks to complete and then close the results channel.
	go func() {
		c.wg.Wait()
		close(c.results)
	}()
}

// checkStatus Helper method to check the status of a single URL.
// this is called in a goroutine for each URL to check.
func (c *Checker) checkStatus(url string) {
	// defer the WaitGroup Done method to decrement the counter when the goroutine completes.
	defer c.wg.Done()
	// Send a GET request to the URL.
	resp, err := c.client.Get(url)
	if err != nil {
		c.results <- fmt.Sprintf("%s: ERROR (%v)", url, err)
		return
	}

	// defer closing the response body to avoid resource leaks.
	defer resp.Body.Close()

	// Send the result to the results channel.
	c.results <- fmt.Sprintf(
		"%s: %d %s",
		url,
		resp.StatusCode,
		http.StatusText(resp.StatusCode),
	)
}

// Results returns the results from the check operation by reading from the results channel.
func (c *Checker) Results() <-chan string {
	return c.results
}

// New creates a new Checker instance that will check the status of the provided URLs.
func New(
	client *http.Client,
	urls []string,
	wg *sync.WaitGroup,
) *Checker {
	return &Checker{
		client: client,
		urls:   urls,
		wg:     wg,
		// initialize the results channel with a buffer size equal to
		// the number of URLs to check to avoid blocking.
		results: make(chan string, len(urls)),
	}
}
