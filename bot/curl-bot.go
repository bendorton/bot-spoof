package bot

func NewCurlBot(requestMethod, endpointURL string) *BasicHTTPBot {
	return &BasicHTTPBot{
		EndpointURL:    endpointURL,
		RequestMethod:  requestMethod,
		UserAgent:      "curl/7.64.1",
		Accept:         "*/*",
		Connection:     "close",
		ContentType:    "*/*",
		AcceptLanguage: "en-US",
		AcceptEncoding: "*",
		Referer:        "",
		Cookies:        "",
		ProxyURL:       "",
		Payload:        "",
	}
}
