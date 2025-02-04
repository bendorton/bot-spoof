package bot

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

type HeadlessBrowser struct {
	endpointURL string
}

func NewHeadlessBrowser(endpoint string) *HeadlessBrowser {
	return &HeadlessBrowser{endpointURL: endpoint}
}

func (this *HeadlessBrowser) Config() any {
	return this
}

func (this *HeadlessBrowser) SendRequest() Response {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var responseBody string

	tasks := chromedp.Tasks{
		chromedp.Navigate(this.endpointURL),
		chromedp.Sleep(3 * time.Second),
		chromedp.OuterHTML("html", &responseBody),
	}

	err := chromedp.Run(ctx, tasks)
	if err != nil {
		return Response{Error: err}
	}

	if containsCloudflareChallenge(responseBody) {
		return Response{StatusCode: 403, Body: "Blocked by Cloudflare", Error: fmt.Errorf("cloudflare challenge detected")}
	}

	return Response{StatusCode: 200, Body: responseBody, Error: nil}
}

func containsCloudflareChallenge(body string) bool {
	return strings.Contains(body, "Checking your browser") ||
		strings.Contains(body, "Please turn JavaScript on") ||
		strings.Contains(body, "challenge-platform")
}
