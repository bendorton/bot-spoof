package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/bendorton/bot-spoof/bot"
)

const (
	endpointURL             = "https://cloudflare.liveaddress.us"
	randomRequestVariations = 10
	requestIterations       = 10
)

type Result struct {
	BotID      int
	StatusCode int
	Body       string
	Error      error
}

func main() {
	fmt.Println("Starting bot tester...")

	// Set up bots to test
	bots := []bot.Bot{
		bot.NewCurlBot("POST", endpointURL),
		bot.NewHeadlessBrowser(endpointURL),
	}
	for range randomRequestVariations {
		bots = append(bots, bot.NewRandomizedBot("POST", endpointURL))
	}

	LogBotConfigs(bots)

	var wg sync.WaitGroup
	resultChan := make(chan Result, len(bots)*requestIterations)

	// Run bots in parallel
	for i, b := range bots {
		wg.Add(1)
		go func(b bot.Bot) {
			defer wg.Done()
			for range requestIterations {
				resp := b.SendRequest()
				resultChan <- Result{BotID: i, StatusCode: resp.StatusCode, Body: resp.Body, Error: resp.Error}
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)))
			}
		}(b)
	}

	// Collect results
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Log results
	results := CollectResults(resultChan)
	LogResults(results)
	fmt.Println("Bot tester finished.")
}

func LogBotConfigs(bots []bot.Bot) {
	fmt.Printf("Bot Configs:\n")
	for i, b := range bots {
		fmt.Printf("[%d]: %+v\n", i, b.Config())
	}
	fmt.Printf("\n")
}

func CollectResults(results chan Result) map[int][]Result {
	var mu sync.Mutex
	botResponses := make(map[int][]Result)

	for res := range results {
		mu.Lock()
		botResponses[res.BotID] = append(botResponses[res.BotID], res)
		mu.Unlock()
	}

	return botResponses
}

func LogResults(results map[int][]Result) {
	var totalRequests, totalFailedRequests int

	for botID, result := range results {
		fmt.Printf("Bot %d:\n", botID)

		var requestCount, errorCount int
		failedResponses := make(map[string]int)
		for _, res := range result {
			totalRequests++
			requestCount++
			if res.Error != nil {
				totalFailedRequests++
				errorCount++

				err := fmt.Sprintf("FAILED Status Code: %d (Error: %v)\n", res.StatusCode, res.Error)
				failedResponses[err]++
			}
		}

		for err, count := range failedResponses {
			fmt.Printf("\t(%d) %s", count, err)
		}

		fmt.Printf("Failed %d/%d requests\n\n", errorCount, requestCount)
	}

	fmt.Printf("\nTest Summary:\n")
	fmt.Printf("Total Requests: %d\n", totalRequests)
	fmt.Printf("Failed Requests: %d\n", totalFailedRequests)
}
