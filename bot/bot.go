package bot

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type Bot interface {
	Config() any
	SendRequest() Response
}

type BasicHTTPBot struct {
	EndpointURL    string
	RequestMethod  string
	Accept         string
	Connection     string
	UserAgent      string
	ContentType    string
	AcceptLanguage string
	AcceptEncoding string
	Referer        string
	Cookies        string
	ProxyURL       string
	Payload        string
}

func (this *BasicHTTPBot) Config() any {
	return this
}

func (this *BasicHTTPBot) SendRequest() Response {
	request, err := http.NewRequest(this.RequestMethod, this.EndpointURL, bytes.NewBuffer([]byte(this.Payload)))
	if err != nil {
		return Response{Error: fmt.Errorf("failed to create request: %w", err)}
	}

	request.Header.Set("Accept", this.Accept)
	request.Header.Set("Connection", this.Connection)
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
		return Response{Error: fmt.Errorf("failed to send request: %w", err)}
	}

	return ParseResponse(res)
}

type Response struct {
	StatusCode int
	Body       string
	Error      error
}

func ParseResponse(response *http.Response) Response {
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
