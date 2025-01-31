package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
)

type BotRequest struct {
	ID             int
	EndpointURL    string
	RequestMethod  string
	UserAgent      string
	ContentType    string
	AcceptLanguage string
	AcceptEncoding string
	Referer        string
	Cookies        string
	ProxyURL       string
	Payload        string
}
type Response struct {
	RequestID  int
	ID         int
	StatusCode int
	Body       string
	Error      error
}

func NewCurlBotRequest(id int, requestMethod, endpointURL string) *BotRequest {
	return &BotRequest{
		ID:             id,
		EndpointURL:    endpointURL,
		RequestMethod:  requestMethod,
		UserAgent:      "curl/7.64.1",
		ContentType:    "*/*",
		AcceptLanguage: "en-US",
		AcceptEncoding: "*",
		Referer:        "",
		Cookies:        "",
		ProxyURL:       "",
		Payload:        "",
	}
}

func NewRandomizedBotRequest(id int, requestMethod, endpointURL string) *BotRequest {
	userAgent := userAgents[rand.Intn(len(userAgents))]
	contentType := contentTypes[rand.Intn(len(contentTypes))]
	acceptLanguage := acceptLanguages[rand.Intn(len(acceptLanguages))]
	acceptEncoding := acceptEncodings[rand.Intn(len(acceptEncodings))]
	referer := referers[rand.Intn(len(referers))]
	cookie := cookies[rand.Intn(len(cookies))]
	proxyURL := proxies[rand.Intn(len(proxies))]
	payload := payloads[rand.Intn(len(payloads))]

	return &BotRequest{
		ID:             id,
		EndpointURL:    endpointURL,
		RequestMethod:  requestMethod,
		UserAgent:      userAgent,
		ContentType:    contentType,
		AcceptLanguage: acceptLanguage,
		AcceptEncoding: acceptEncoding,
		Referer:        referer,
		Cookies:        cookie,
		ProxyURL:       proxyURL,
		Payload:        payload,
	}
}

func (this *BotRequest) Send() (*http.Response, error) {
	request, err := http.NewRequest(this.RequestMethod, this.EndpointURL, bytes.NewBuffer([]byte(this.Payload)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	request.Header.Set("Accept", "*/*")
	request.Header.Set("Connection", "close")
	request.Header.Set("User-Agent", this.UserAgent)
	request.Header.Set("Content-Type", this.ContentType)
	request.Header.Set("Accept-Language", this.AcceptLanguage)
	request.Header.Set("Accept-Encoding", this.AcceptEncoding)
	request.Header.Set("Referer", this.Referer)
	request.Header.Set("Cookie", this.Cookies)

	client := &http.Client{}
	//if this.ProxyURL != "" {
	//	proxyURL, err := url.Parse(this.ProxyURL)
	//	if err == nil {
	//		client = &http.Client{
	//			Transport: &http.Transport{
	//				Proxy: http.ProxyURL(proxyURL),
	//			},
	//		}
	//	}
	//}

	res, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	return res, nil
}

var (
	userAgents = []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; rv:49.0) Gecko/20100101 Firefox/49.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.121 Safari/537.36 Edge/85.0.564.51",
		"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
		"Chrome/51.0.2704.103 Safari/537.36",
		"BotSpoof/1.0",
		"curl/7.64.1",
	}
	contentTypes = []string{
		"*/*",
		"application/json",
		"plain/text",
		"application/x-www-form-urlencoded",
	}
	acceptLanguages = []string{
		"en-US,en;q=0.9",
		"en-GB,en;q=0.8",
		"es-ES,es;q=0.9",
	}
	acceptEncodings = []string{
		"*",
		"identity",
		"br;q=1.0, gzip;q=0.6, *;q=0.1",
		"gzip, deflate, zstd",
		"gzip, deflate, br",
	}
	referers = []string{
		"https://www.google.com",
		"https://www.twitter.com",
		"https://www.facebook.com",
		"https://example.com",
		"https://www.youtube.com",
		"https://developer.mozilla.org/en-US/docs/Web/JavaScript",
	}
	cookies = []string{
		"_cfduid=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		"name=value; test=2025-01-01",
		"PHPSESSID=298zf09hf012fh2; csrftoken=u32t4o3tb3gg43; _gat=1",
	}

	proxies = []string{
		"https://example.com",
		"https://test.example.com",
		"https://www.example.com",
		"https://www.test.com",
	}
	payloads = []string{
		`{"test": "test"}`,
		"payload",
	}
)
