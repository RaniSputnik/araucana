package scrape

import (
	"log"
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

type scraper struct {
	rootURL *url.URL
	results map[string]*Page
	logger  *log.Logger
}

type scrapeResult struct {
	NextURLs []string
	Err      error
}

func (s *scraper) Scrape(addr string) error {
	// TODO provide context to method so timeout can be provided
	// TODO limit the recursion to a fixed max
	thisPage := &Page{addr, []*Asset{}, []string{}}
	s.results[addr] = thisPage
	s.logger.Printf("Scraping %s", addr)

	cResults := make(chan *scrapeResult)
	cStop := make(chan interface{})
	defer close(cStop)
	go s.doScrape(thisPage, cResults, cStop)

	inflight := 1
	for inflight > 0 {
		inflight--
		res := <-cResults
		if res.Err != nil {
			return res.Err
		}

		for _, nextURL := range res.NextURLs {
			if _, done := s.results[nextURL]; !done {
				nextPage := &Page{nextURL, []*Asset{}, []string{}}
				s.results[nextURL] = nextPage
				inflight++
				go s.doScrape(nextPage, cResults, cStop)
			} else {
				s.logger.Printf("We've already scraped '%s'", nextURL)
			}
		}
	}

	return nil
}

func (s *scraper) doScrape(page *Page, cResults chan<- *scrapeResult, done <-chan interface{}) {
	// TODO never ever use the default client in production
	response, err := http.Get(page.URL)
	if err != nil {
		cResults <- &scrapeResult{Err: ErrHTTPError}
	}

	if httpStatusIsError(response.StatusCode) {
		cResults <- &scrapeResult{Err: ErrHTTPError}
	}

	doc, err := html.Parse(response.Body)
	if err != nil {
		// TODO return defined error
		cResults <- &scrapeResult{Err: err}
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if nextURL := s.getLinkIfExistsInNode(n); nextURL != "" {
			page.Pages = appendPageIfNotPresent(page.Pages, nextURL)
		} else if asset := s.getAssetIfExistsInNode(n); asset != nil {
			page.Assets = append(page.Assets, asset)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	// TODO do we need to remove duplicate assets from a page?
	// That could be a neat feature - find when you have duplicate
	// links to a single resource

	select {
	case cResults <- &scrapeResult{NextURLs: page.Pages}:
	case <-done:
	}
}

func (s *scraper) getLinkIfExistsInNode(n *html.Node) string {
	if n.Type != html.ElementNode || n.Data != "a" {
		// Skip this node, it's not an <a> tag
		return ""
	}

	ok, href := attr(n, "href")
	if !ok {
		s.logger.Printf("<a> tag appears to have no 'href' attribute")
		return ""
	}

	parsedHref, err := s.GetFullURLWithoutHashAndQuery(href)
	if err != nil {
		s.logger.Printf("<a> tag has a href attribute (%s) we can't parse: '%v'", href, err)
		return ""
	}

	if parsedHref.Host != s.rootURL.Host {
		s.logger.Printf("External link will not be followed '%s'", href)
		return ""
	}

	return parsedHref.String()
}

func (s *scraper) getAssetIfExistsInNode(n *html.Node) *Asset {
	if n.Type != html.ElementNode {
		return nil
	}

	var attrName, assetType string
	switch n.Data {
	case "link":
		attrName, assetType = "href", AssetTypeLink
	case "img":
		attrName, assetType = "src", AssetTypeImage
	case "script":
		attrName, assetType = "src", AssetTypeScript
	default:
		return nil
	}

	if ok, src := attr(n, attrName); ok {
		if fullURL, err := s.GetFullURL(src); err == nil {
			return &Asset{Type: assetType, URL: fullURL.String()}
		}
	}
	return nil
}

func (s *scraper) GetFullURL(val string) (*url.URL, error) {
	parsedVal, err := url.Parse(val)
	if err != nil {
		return nil, err
	}

	parsedVal = s.rootURL.ResolveReference(parsedVal)
	return parsedVal, nil
}

func (s *scraper) GetFullURLWithoutHashAndQuery(val string) (*url.URL, error) {
	parsedVal, err := s.GetFullURL(val)
	parsedVal.RawQuery = ""
	parsedVal.Fragment = ""
	return parsedVal, err
}

func attr(t *html.Node, name string) (bool, string) {
	for _, a := range t.Attr {
		if a.Key == name {
			return true, a.Val
		}
	}
	return false, ""
}

func httpStatusIsError(status int) bool {
	return status == 0 || status >= 400
}

func appendPageIfNotPresent(pages []string, pageURL string) []string {
	for _, existingPageURL := range pages {
		if existingPageURL == pageURL {
			return pages
		}
	}
	return append(pages, pageURL)
}
