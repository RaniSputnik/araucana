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
	response, err := defaultHTTPClient.Get(url)
	if err != nil {
		return response.Body, err
	}

	if httpStatusIsError(response.StatusCode) {
		return response.Body, fmt.Errorf("Got error http status, '%s'", response.Status)
	}

	return response.Body, nil
})
