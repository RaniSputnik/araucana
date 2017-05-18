package scrape

import (
	"context"
	"log"
	"net/http"
	"net/url"
)

type crawler struct {
	rootURL *url.URL
	client  *http.Client
	logger  *log.Logger
	scraper Scraper
}

type scrapeResult struct {
	NextURLs []string
	Err      error
}

func (s *crawler) Crawl(ctx context.Context, startAddr string) (map[string]*Page, error) {
	results := make(map[string]*Page)
	resultsChan := make(chan *scrapeResult)
	inflight := 0

	startPageScrape := func(addr string) {
		thisPage := &Page{addr, []*Asset{}, []string{}}
		results[addr] = thisPage
		inflight++

		go s.downloadAndScrapePage(ctx, thisPage, resultsChan)
	}

	startPageScrape(startAddr)

	for inflight > 0 {
		inflight--
		select {
		case res := <-resultsChan:
			if res.Err != nil {
				return nil, res.Err
			}

			for _, nextURL := range res.NextURLs {
				if _, alreadyScraped := results[nextURL]; !alreadyScraped {
					startPageScrape(nextURL)
				} else {
					s.logger.Printf("We've already scraped '%s'", nextURL)
				}
			}

		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return results, nil
}

func (s *crawler) downloadAndScrapePage(ctx context.Context, page *Page, resultsChan chan<- *scrapeResult) {
	s.logger.Printf("Scraping %s", page.URL)

	response, err := s.client.Get(page.URL)
	if err != nil {
		resultsChan <- &scrapeResult{Err: ErrHTTPError}
		return
	}
	defer response.Body.Close()

	if httpStatusIsError(response.StatusCode) {
		select {
		case resultsChan <- &scrapeResult{Err: ErrHTTPError}:
		case <-ctx.Done():
		}
	}

	if err = s.scraper.Scrape(response.Body, page); err != nil {
		select {
		case resultsChan <- &scrapeResult{Err: err}:
		case <-ctx.Done():
		}
	}

	select {
	case resultsChan <- &scrapeResult{NextURLs: page.Pages}:
	case <-ctx.Done():
	}
}
