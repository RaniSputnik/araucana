package scrape

import "encoding/xml"

const SitemapXMLNamespace = "http://www.sitemaps.org/schemas/sitemap/0.9"

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
func Site(url string) (*Sitemap, error) {
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
