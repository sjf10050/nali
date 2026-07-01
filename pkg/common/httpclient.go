package common

import (
	"io"
	"log"
	"net/http"
	"time"
)

const UserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36"

const MaxResponseSize = 500 * 1024 * 1024

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

func GetHttpClient() *HttpClient {
	return httpClient
}

func (c *HttpClient) Get(urls ...string) (body []byte, err error) {
	for _, url := range urls {
		req, reqErr := http.NewRequest(http.MethodGet, url, nil)
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
			defer resp.Body.Close()
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
