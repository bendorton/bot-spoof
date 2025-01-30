package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

const (
	endpointURL = "http://cloudflare.liveaddress.us"
	iterations  = 20
	workerCount = 3
)

type Response struct {
	StatusCode int
	Body       string
	Error      error
}

func main() {
	fmt.Println("Starting bot tester...")

	requestChan := make(chan int, iterations)
	resultChan := make(chan Response)
	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(i, requestChan, resultChan, &wg)
	}

	go func() {
		for i := 0; i < iterations; i++ {
			requestChan <- i
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)))
		}
		close(requestChan)
	}()

	// Collect Results
	var totalRequests, failedRequests []Response
	var mu sync.Mutex
	go func() {
		for result := range resultChan {
			mu.Lock()
			totalRequests = append(totalRequests, result)
			if result.Error != nil {
				failedRequests = append(failedRequests, result)
			}
			mu.Unlock()
		}
	}()

	// Wait for all workers to finish
	wg.Wait()
	close(resultChan)

	// Output results
	fmt.Println("Bot Spoof Report.")
	fmt.Printf("Total Requests: %d\n", len(totalRequests))
	fmt.Printf("Failed requests: %v\n", len(failedRequests))
	for _, failure := range failedRequests {
		fmt.Printf("Failed Request - Status: %d, Error: %v, Body: %s\n", failure.StatusCode, failure.Error, failure.Body)
	}
}

func worker(id int, requests <-chan int, results chan<- Response, wg *sync.WaitGroup) {
	defer wg.Done()
	for requestID := range requests {
		fmt.Printf("Worker %d processing request %d\n", id, requestID)
		botRequest := NewRandomBotRequest("POST", endpointURL)
		resp, err := botRequest.Send()
		if err != nil {
			results <- Response{Error: err}
		} else {
			results <- parseResponse(resp)
		}
	}
}

func parseResponse(response *http.Response) Response {
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return Response{StatusCode: response.StatusCode, Body: "", Error: fmt.Errorf("failed to read response: %w", err)}
	}

	if response.StatusCode != 200 {
		return Response{StatusCode: response.StatusCode, Body: string(body), Error: fmt.Errorf("unexpected status code")}
	}

	return Response{StatusCode: response.StatusCode, Body: string(body), Error: nil}
}
