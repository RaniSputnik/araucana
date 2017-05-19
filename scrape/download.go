package scrape

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type Downloader interface {
	Download(url string) (io.ReadCloser, error)
}

type DownloaderFunc func(url string) (io.ReadCloser, error)

func (f DownloaderFunc) Download(url string) (io.ReadCloser, error) {
	return f(url)
}

var defaultHTTPClient = &http.Client{
	Timeout: time.Second * 30,
}

var DefaultDownloaderFunc = DownloaderFunc(func(url string) (io.ReadCloser, error) {
	var body io.ReadCloser
	response, err := defaultHTTPClient.Get(url)
	if response != nil {
		body = response.Body
	}
	if err != nil {
		return body, err
	}

	if httpStatusIsError(response.StatusCode) {
		return body, fmt.Errorf("Got error http status, '%s'", response.Status)
	}

	return body, nil
})
