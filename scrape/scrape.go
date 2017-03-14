package scrape

import (
	"encoding/xml"
	"errors"
	"net/url"
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
	if site == "" {
		return nil, ErrURLInvalid
	}
	_, err := url.Parse(site)
	if err != nil {
		return nil, ErrURLInvalid
	}

	// TODO dynamically generate sitemap
	return &Sitemap{
		XMLNS: SitemapXMLNamespace,
		URLSet: []*SitemapURL{
			&SitemapURL{
				Loc: "http://tomblomfield.com/about",
			},
			&SitemapURL{
				Loc: "http://tomblomfield.com/random",
			},
		},
	}, nil
}
