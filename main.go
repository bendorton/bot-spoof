package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// TODO verify that http requests are denied
const (
	endpointURL             = "https://cloudflare.liveaddress.us"
	randomRequestVariations = 10
	requestIterations       = 10
	workerCount             = 5
)

func main() {
	fmt.Println("Starting bot tester...")

	var botRequests []*BotRequest
	requestChan := make(chan *BotRequest, requestIterations*randomRequestVariations)
	resultChan := make(chan Response)
	var wg sync.WaitGroup

	for i := range randomRequestVariations {
		botRequests = append(botRequests, NewRandomizedBotRequest(i, "POST", endpointURL))
	}
	botRequests = append(botRequests, NewCurlBotRequest(len(botRequests), "POST", endpointURL))

	// Start worker goroutines
	for i := range workerCount {
		wg.Add(1)
		go worker(i, requestChan, resultChan, &wg)
	}

	go func() {
		for _, request := range botRequests {
			requestChan <- request
		}
		close(requestChan)
	}()

	// Collect Results
	responses := make(map[int][]Response)
	var totalResponseCount, failedRequestCount int
	var mu sync.Mutex
	go func() {
		for result := range resultChan {
			mu.Lock()
			totalResponseCount++
			responses[result.RequestID] = append(responses[result.RequestID], result)
			if result.Error != nil {
				failedRequestCount++
			}
			mu.Unlock()
		}
	}()

	// Wait for all workers to finish
	wg.Wait()
	close(resultChan)

	// Output results
	fmt.Println("\nBot Spoof Report")
	fmt.Printf("Requests:\n")
	for _, request := range botRequests {
		fmt.Printf("Request %d:\n", request.ID)
		fmt.Printf("\t%+v\n", request)
		fmt.Printf("Requests Sent: %d\n", len(responses[request.ID]))
		for _, response := range responses[request.ID] {
			if response.Error != nil {
				fmt.Printf("\t%d: %s\n", response.StatusCode, response.Error)
			}
		}
		fmt.Println()
	}
	fmt.Printf("Total Requests Sent: %d\n", totalResponseCount)
	fmt.Printf("Total Failed Requests: %d\n", failedRequestCount)
}

func worker(id int, requests <-chan *BotRequest, results chan<- Response, wg *sync.WaitGroup) {
	defer wg.Done()
	for request := range requests {
		fmt.Printf("Worker %d processing request ID %d\n", id, request.ID)
		for i := range requestIterations {
			resp, err := request.Send()
			if err != nil {
				results <- Response{RequestID: request.ID, ID: i, Error: err}
			} else {
				results <- parseResponse(request.ID, i, resp)
			}
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)))
		}
		fmt.Printf("Worker %d finished processing request ID %d\n", id, request.ID)
	}
}

func parseResponse(requestID, requestIteration int, response *http.Response) Response {
	defer response.Body.Close()

	res := Response{
		RequestID:  requestID,
		ID:         requestIteration,
		StatusCode: response.StatusCode,
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		res.Error = fmt.Errorf("failed to read response: %w", err)
		return res
	}
	res.Body = string(body)

	if response.StatusCode != 200 {
		res.Error = fmt.Errorf("unexpected status code")
	}

	return res
}
