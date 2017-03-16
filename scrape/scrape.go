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

func (s *scraper) Scrape(addr string) error {
	// TODO provide context to method so timeout can be provided
	// TODO limit the recursion to a fixed max
	thisPage := &Page{addr, []*Asset{}}
	s.results[addr] = thisPage
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
		s.tryFollowLinkIfExistsInNode(n)
		s.tryAddAsset(n, thisPage)
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	// TODO do we need to remove duplicate assets from a page?
	// That could be a neat feature - find when you have duplicate
	// links to a single resource

	return nil
}

func (s *scraper) tryFollowLinkIfExistsInNode(n *html.Node) {
	if n.Type != html.ElementNode || n.Data != "a" {
		// Skip this node, it's not an <a> tag
		return
	}
	ok, href := attr(n, "href")
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

// tryAddAsset will check if the node is an asset reference
// (link|img|script). If the node does represent an asset then
// the asset will be created and added to the specified page
func (s *scraper) tryAddAsset(n *html.Node, page *Page) {
	if n.Type != html.ElementNode {
		return
	}

	switch n.Data {
	case "link":
		if ok, src := attr(n, "href"); ok {
			// TODO check rel=stylesheet || rel="" && ext=css
			if fullURL, err := s.GetFullURL(src); err == nil {
				page.Assets = append(page.Assets, &Asset{Type: AssetTypeLink, URL: fullURL.String()})
			}
		}

	case "img":
		if ok, src := attr(n, "src"); ok {
			if fullURL, err := s.GetFullURL(src); err == nil {
				page.Assets = append(page.Assets, &Asset{Type: AssetTypeImage, URL: fullURL.String()})
			}
		}

	case "script":
		if ok, src := attr(n, "src"); ok {
			if fullURL, err := s.GetFullURL(src); err == nil {
				page.Assets = append(page.Assets, &Asset{Type: AssetTypeScript, URL: fullURL.String()})
			}
		}
	}
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
