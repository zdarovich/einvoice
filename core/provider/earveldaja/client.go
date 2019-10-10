package earveldaja

import (
	"github.com/zdarovich/einvoice/logging"
	"io"
	"net"
	"net/http"
	"time"
)

const (
	timeout               = 10 * time.Second
	keepAlive             = 10 * time.Second
	tlsHandshakeTimeout   = 10 * time.Second
	expectContinueTimeout = 10 * time.Second
	responseHeaderTimeout = 10 * time.Second
	httpClientTimeout     = 10 * time.Second
	maxIdleConns          = 1
	maxConnsPerHost       = 1
)

type HTTPClient struct {
	//HttpClient ...
	Cli *http.Client
	URL string
}

func NewHTTPClient(url string) *HTTPClient {
	tr := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   timeout,
			KeepAlive: keepAlive,
		}).DialContext,
		TLSHandshakeTimeout: tlsHandshakeTimeout,

		ExpectContinueTimeout: expectContinueTimeout,
		ResponseHeaderTimeout: responseHeaderTimeout,

		MaxIdleConns:    maxIdleConns,
		MaxConnsPerHost: maxConnsPerHost,
	}

	cli := HTTPClient{
		URL: url,
		Cli: &http.Client{
			Transport: tr,
			Timeout:   httpClientTimeout,
		},
	}
	return &cli
}

func (cli *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	start := time.Now()
	resp, err := cli.Cli.Do(req)
	elapsed := time.Since(start)
	logging.HTTP(req, resp, err, elapsed, req.Method)
	return resp, err
}

func (cli *HTTPClient) NewRequest(method string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, cli.URL, body)
}
