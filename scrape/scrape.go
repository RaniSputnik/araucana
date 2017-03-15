package scrape

import (
	"encoding/xml"
	"errors"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

const SitemapXMLNamespace = "http://www.sitemaps.org/schemas/sitemap/0.9"

var (
	// ErrURLInvalid is given when the URL provided to the 'Site'
	// method is empty or invalid
	ErrURLInvalid = errors.New("The given URL is invalid")

	// ErrHTTPError is given when the URL provided results in a
	// HTTP error code or could not be reached.
	ErrHTTPError = errors.New("The given URL gave a http error code")
)

// Sitemap represents a heirachy of pages within a webiste
type Sitemap struct {
	XMLName xml.Name      `xml:"urlset"`
	XMLNS   string        `xml:"xmlns,attr"`
	URLSet  []*SitemapURL `xml:"url"`
}

// SitemapURL represents a location within a sitemap.
// Should be indicative of a page within the website.
type SitemapURL struct {
	Loc string `xml:"loc"`
	// TODO sitemap image
	// TODO sitemap js
	// TODO sitemap css
}

// Site will generate a sitemap for the given URL.
// The sitemap will be constrained to a given domain,
// external links will not be followed.
//
// An error will be thrown if the url is invalid or
// the site can not be reached for any reason. Partial
// sitemaps will not be returned.
func Site(site string) (*Sitemap, error) {
	// Validation
	if site == "" {
		return nil, ErrURLInvalid
	}
	siteURL, err := url.Parse(site)
	if err != nil {
		return nil, ErrURLInvalid
	}

	// Run the scraping of the site
	s := &scraper{
		rootURL: siteURL,
		results: map[string]*SitemapURL{},
	}
	if err = s.Scrape(siteURL.String()); err != nil {
		return nil, err
	}

	// Copy map of scraped sites into URL set
	i := 0
	urlset := make([]*SitemapURL, len(s.results))
	for _, val := range s.results {
		urlset[i] = val
		i++
	}

	// TODO dynamically generate sitemap
	return &Sitemap{
		XMLNS:  SitemapXMLNamespace,
		URLSet: urlset,
	}, nil
}

type scraper struct {
	rootURL *url.URL
	results map[string]*SitemapURL
}

func (s *scraper) Scrape(addr string) error {
	s.results[addr] = &SitemapURL{
		Loc: addr,
	}
	log.Printf("Scraping %s", addr)

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
		if n.Type == html.ElementNode && n.Data == "a" {
			if ok, href := getHref(n); ok {
				if href, err := s.GetFullURL(href); err == nil {
					if _, ok := s.results[href]; !ok {
						s.Scrape(href)
					} else {
						// TODO use logger on scraper
						log.Printf("We've already scraped '%s'", href)
					}
				} else {
					// TODO use logger on scraper
					log.Printf("<a> tag has a href attribute (%s) we can't parse: '%v'", href, err)
				}
			} else {
				// TODO use logger on scraper
				log.Printf("<a> tag appears to have no 'href' attribute")
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return nil
}

func (s *scraper) GetFullURL(val string) (string, error) {
	parsedVal, err := url.Parse(val)
	if err != nil {
		return val, err
	}

	if parsedVal.Scheme == "" {
		parsedVal.Scheme = s.rootURL.Scheme
	}
	if parsedVal.Host == "" {
		parsedVal.Host = s.rootURL.Host
	}

	return parsedVal.String(), nil
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
