package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/bendorton/bot-spoof/bot"
)

// TODO verify that http requests are denied
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
		bot.NewScraper(endpointURL),
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
	LogResults(resultChan)
	fmt.Println("Bot tester finished.")
}

func LogBotConfigs(bots []bot.Bot) {
	fmt.Printf("Bot Configs:\n")
	for i, b := range bots {
		fmt.Printf("[%d]: %+v\n", i, b.Config())
	}
	fmt.Printf("\n")
}

func LogResults(results chan Result) {
	var totalRequests, failedRequests int

	for res := range results {
		totalRequests++
		if res.Error != nil {
			failedRequests++
			fmt.Printf("Bot %d: FAILED Status Code: %d (Error: %v)\n", res.BotID, res.StatusCode, res.Error)
		}
	}

	fmt.Printf("\nTest Summary:\n")
	fmt.Printf("Total Requests: %d\n", totalRequests)
	fmt.Printf("Failed Requests: %d\n", failedRequests)
}
