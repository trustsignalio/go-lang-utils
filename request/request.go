package request

import (
	"bytes"
	"io/ioutil"
	netHttp "net/http"
	"strings"
	"time"
)

// ClientOptions
type ClientOptions struct {
	Timeout int
}

// Client
type Client struct {
	backendClient *netHttp.Client
}

// RequestOptions struct
type RequestOptions struct {
	Method  string
	URL     string
	Body    string
	Retries int
	Query   map[string]string
	Headers map[string]string

	RetryInterval time.Duration // retry after x milliseconds
}

// Response struct
type Response struct {
	Body, Status string
	StatusCode   int
	Header       map[string][]string
	Cookies      []*netHttp.Cookie
}

// NewClient
func NewClient(opts *ClientOptions) *Client {
	var c = &netHttp.Client{
		Timeout: time.Duration(opts.Timeout) * time.Second,
	}
	var client = &Client{backendClient: c}
	return client
}

// Request method will make a request
func (c *Client) Request(opts *RequestOptions) (*Response, error) {
	var url = opts.URL
	if !strings.Contains(opts.URL, "?") {
		url += "?"
	}
	for k, v := range opts.Query {
		url += "&" + k + "=" + v
	}
	byteBody := []byte(opts.Body)
	var body = bytes.NewBuffer(byteBody)
	var req, err = netHttp.NewRequest(opts.Method, url, body)
	if err != nil {
		return nil, err
	}
	var lastError error
	if opts.RetryInterval == 0 {
		opts.RetryInterval = 1 * time.Second
	}
	var respData = &Response{}
	for index := 0; index < opts.Retries; index++ {
		var resp, err = c.backendClient.Do(req)
		if err != nil {
			lastError = err
			// If request is failed then retry after sleeping for some time
			time.Sleep(opts.RetryInterval)
			continue
		}
		defer resp.Body.Close()
		body, readErr := ioutil.ReadAll(resp.Body)
		respData.Body = string(body)
		respData.Status = resp.Status
		respData.StatusCode = resp.StatusCode
		respData.Cookies = resp.Cookies()

		for k, v := range resp.Header {
			respData.Header[k] = v
		}

		if readErr != nil {
			return respData, readErr
		}
		break // since it is success break out of the for loop
	}
	return respData, lastError
}
