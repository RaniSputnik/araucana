package scrape

import (
	"errors"
	"log"
	"net/url"
	"os"
)

// Sitemap represents a heirachy of pages within a webiste
type Sitemap struct {
	Pages []*Page `json:"pages"`
}

// Page represents a location within a sitemap.
// Should be indicative of a page within the website.
type Page struct {
	URL    string   `json:"url"`
	Assets []*Asset `json:"assets"`
	// TODO crosslink pages
}

// Asset represents a reference to a piece of static content.
// Assets include stylesheets, images and scripts.
// Assets can be external because they will not be followed by
// the scraper. They do not represent the content that was served
// but the reference from the page.
type Asset struct {
	Type string
	URL  string
}

const (
	// AssetTypeLink is used for <link> assets
	AssetTypeLink = "link"
	// AssetTypeImage is used for <image> assets
	AssetTypeImage = "image"
	// AssetTypeScript is used for <script> assets
	AssetTypeScript = "script"
)

var (
	// ErrURLInvalid is given when the URL provided to the 'Site'
	// method is empty or invalid
	ErrURLInvalid = errors.New("The given URL is invalid")

	// ErrHTTPError is given when the URL provided results in a
	// HTTP error code or could not be reached.
	ErrHTTPError = errors.New("The given URL gave a http error code")
)

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
		results: map[string]*Page{},
		logger:  log.New(os.Stdout, "", log.LstdFlags),
	}
	if err = s.Scrape(siteURL.String()); err != nil {
		return nil, err
	}

	// Copy map of scraped sites into URL set
	i := 0
	urlset := make([]*Page, len(s.results))
	for _, val := range s.results {
		urlset[i] = val
		i++
	}

	// Return the sitemap response
	return &Sitemap{
		Pages: urlset,
	}, nil
}
