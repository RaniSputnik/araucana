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
}

type scrapeResult struct {
	NextURLs []string
	Err      error
}

func (s *crawler) Crawl(ctx context.Context, startAddr string) (map[string]*Page, error) {
	results := make(map[string]*Page)
	cResults := make(chan *scrapeResult)
	inflight := 0

	startPageScrape := func(addr string) {
		thisPage := &Page{addr, []*Asset{}, []string{}}
		results[addr] = thisPage
		inflight++

		go s.scrape(ctx, thisPage, cResults)
	}

	startPageScrape(startAddr)

	for inflight > 0 {
		inflight--
		select {
		case res := <-cResults:
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
