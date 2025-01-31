package bot

import (
	"context"
	"time"

	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
)

type Scraper struct {
	endpointURL string
}

func NewScraper(endpoint string) *Scraper {
	return &Scraper{endpointURL: endpoint}
}

func (this *Scraper) Config() any {
	return this
}

func (this *Scraper) SendRequest() Response {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var html string

	err := chromedp.Run(ctx,
		chromedp.Navigate(this.endpointURL),
		chromedp.Sleep(2*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			rootNode, err := dom.GetDocument().Do(ctx)
			if err != nil {
				return err
			}

			html, err = dom.GetOuterHTML().WithNodeID(rootNode.NodeID).Do(ctx)
			return err
		}),
	)
	if err != nil {
		return Response{
			Body:  html,
			Error: err,
		}
	}

	return Response{Body: html}
}
