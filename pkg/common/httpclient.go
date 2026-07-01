package common

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"
)

// UserAgent is the User-Agent header sent with database download requests.
const UserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36"

// MaxResponseSize caps the number of bytes read from a single response.
const MaxResponseSize = 500 * 1024 * 1024

// HttpClient is an http.Client wrapper used to fetch database files.
type HttpClient struct {
	*http.Client
}

var httpClient *HttpClient

func init() {
	httpClient = &HttpClient{
		&http.Client{
			Timeout: time.Second * 300,
			Transport: &http.Transport{
				TLSHandshakeTimeout:   time.Second * 5,
				IdleConnTimeout:       time.Second * 10,
				ResponseHeaderTimeout: time.Second * 10,
				ExpectContinueTimeout: time.Second * 20,
				Proxy:                 http.ProxyFromEnvironment,
			},
		},
	}
}

// GetHttpClient returns the shared HTTP client.
func GetHttpClient() *HttpClient {
	return httpClient
}

// Get fetches the first URL that responds with 200 OK and returns its body.
func (c *HttpClient) Get(ctx context.Context, urls ...string) (body []byte, err error) {
	for _, url := range urls {
		req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if reqErr != nil {
			log.Println(reqErr)
			err = reqErr
			continue
		}
		req.Header.Set("User-Agent", UserAgent)

		resp, doErr := c.Do(req)
		if doErr != nil {
			log.Println(doErr)
			err = doErr
			continue
		}

		body, err = func() ([]byte, error) {
			defer func() { _ = resp.Body.Close() }()
			if resp.StatusCode != http.StatusOK {
				return nil, nil
			}
			return io.ReadAll(io.LimitReader(resp.Body, MaxResponseSize))
		}()
		if err != nil {
			continue
		}
		if body != nil {
			return body, nil
		}
	}

	return nil, err
}
