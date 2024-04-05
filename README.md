# Concurrency in Go

This is a multipart series on Go Concurrency. Each directory will have it's own "runbook" to help you understand the code and get you up and running with concurrent tasks quickly in Go.

## Go Concurrency Scenarios
- **URL Checker**: Demonstrates how to use goroutines, channels, and waitgroups to concurrently check the status of multiple URLs.
- **Semaphore Pattern**: Shows how to use a buffered channel as a semaphore to limit the number of concurrent goroutines. This is useful when limitations on the number of concurrent database reads, or calls to an API are required. 

### NOTE ON DOCUMENTATION:
Each directory I have a markdown file that will walk you through the code and explain how it works. They contain code you can copy and paste into the main.go file to get you up and running quickly.

For instance urlchecker/checker.md will walk you through the code in the urlchecker directory.

## Documentation
#### URL Checker:
[URL Checker Documentation](./urlchecker/checker.md)
#### Semaphore Pattern:
[Semaphore Pattern Documentation](./semaphore/semaphore.md)


## Getting Started

### Prerequisites
Go (version 1.22 or higher)

### Execution
Clone the repository:

```bash
git clone https://github.com/lkendrickd/goconcurrency.git
cd goconcurrency
```

Consult one of the Mardown files in the directories to get started. Place the code of your chosen example in the main.go file and run the project.

Run the Project
```bash
go run main.go
```

### Run the tests
Alternatively, you can run the tests for the project using the following command:

```bash
go test -v ./...
```
