package scrape

import (
	"encoding/xml"
	"errors"
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

const SitemapXMLNamespace = "http://www.sitemaps.org/schemas/sitemap/0.9"

var (
	// ErrURLInvalid is given when the URL provided to the 'Site'
	// method is empty or invalid
	ErrURLInvalid = errors.New("The given URL is invalid")
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
	results map[string]*SitemapURL
}

func (s *scraper) Scrape(addr string) error {
	s.results[addr] = &SitemapURL{
		Loc: addr,
	}

	response, err := http.Get(addr)
	if err != nil {
		// TODO return defined error
		return err
	}

	_, err = html.Parse(response.Body)
	if err != nil {
		// TODO return defined error
		return err
	}

	// TODO read the body and scrape again

	return nil
}
