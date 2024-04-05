# Go Concurrency Steps to Success: URL Checker

In the file checker.go here are some key operations that are taking place to make concurrent requests to a list of URLs and check their status codes.

### Step 1: Define Your Data Structure
Create a Checker struct with an HTTP client, a slice of URLs, a sync.WaitGroup, and a results channel.

### Step 2: Initialize Your Checker
Implement a New function to instantiate your Checker with a given HTTP client, URLs, and an initialized results channel sized to the number of URLs.

### Step 3: Set Up Concurrency
For each URL, increment the WaitGroup and launch a goroutine that calls a checkStatus function.

### Step 4: Implement the checkStatus Function
In checkStatus, make an HTTP GET request to the URL.
Use defer to decrement the WaitGroup and to close the response body.
Send the status code (or error) back through the results channel.

### Step 5: Wait for Completion
After launching all goroutines, start another goroutine that waits for all tasks to finish (using WaitGroup.Wait()) and then closes the results channel.

### Step 6: Collect Results
Provide a Results method to return the results channel for reading the status of each URL check.

### Step 7: Clean Up
Ensure proper error handling and resource management, especially with closing the response bodies and the results channel.

#### Key Points to Remember:
- Goroutines for concurrent execution.
- Channels to communicate between goroutines.
- WaitGroup to synchronize completion.
- Resource Management to avoid leaks by closing response bodies and channels.

### Usage

To use the Concurrency Checker, you need to specify the URLs you want to check.
This can be done by modifying the main.go file.
Below is an example of how to use the Concurrency Checker:

**main.go**

```go
package main

import (
	"net/http"
	"sync"

	"github.com/lkendrickd/concurrency/urlchecker"
)

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
```
