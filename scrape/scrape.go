package scrape

import (
	"log"
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

type scraper struct {
	rootURL *url.URL
	results map[string]*SitemapURL
	logger  *log.Logger
}

func (s *scraper) Scrape(addr string) error {
	// TODO provide context to method so timeout can be provided
	// TODO limit the recursion to a fixed max
	s.results[addr] = &SitemapURL{
		Loc: addr,
	}
	s.logger.Printf("Scraping %s", addr)

	// TODO never ever use the default client in production
	response, err := http.Get(addr)
	if err != nil {
		return ErrHTTPError
	}

	if httpStatusIsError(response.StatusCode) {
		return ErrHTTPError
	}

	doc, err := html.Parse(response.Body)
	if err != nil {
		// TODO return defined error
		return err
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		s.followLinkIfExistsInNode(n)
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return nil
}

func (s *scraper) followLinkIfExistsInNode(n *html.Node) {
	if n.Type != html.ElementNode || n.Data != "a" {
		// Skip this node, it's not an <a> tag
		return
	}
	ok, href := getHref(n)
	if !ok {
		s.logger.Printf("<a> tag appears to have no 'href' attribute")
		return
	}

	parsedHref, err := s.GetFullURLWithoutHashAndQuery(href)
	if err != nil {
		s.logger.Printf("<a> tag has a href attribute (%s) we can't parse: '%v'", href, err)
		return
	}

	if parsedHref.Host != s.rootURL.Host {
		s.logger.Printf("External link will not be followed '%s'", href)
		return
	}

	href = parsedHref.String()
	if _, ok := s.results[href]; ok {
		s.logger.Printf("We've already scraped '%s'", href)
		return
	}

	// All the checks have run, we can safely scrape now
	s.Scrape(href)
}

func (s *scraper) GetFullURLWithoutHashAndQuery(val string) (*url.URL, error) {
	parsedVal, err := url.Parse(val)
	if err != nil {
		return nil, err
	}

	parsedVal = s.rootURL.ResolveReference(parsedVal)
	parsedVal.RawQuery = ""
	parsedVal.Fragment = ""

	return parsedVal, nil
}

func getHref(t *html.Node) (ok bool, href string) {
	for _, a := range t.Attr {
		if a.Key == "href" {
			ok = true
			href = a.Val
			break
		}
	}
	return
}

func httpStatusIsError(status int) bool {
	return status == 0 || status >= 400
}
