package scrape

import (
	"context"

	"golang.org/x/net/html"
)

func (s *crawler) scrape(ctx context.Context, page *Page, cResults chan<- *scrapeResult) {
	s.logger.Printf("Scraping %s", page.URL)

	response, err := s.client.Get(page.URL)
	if err != nil {
		cResults <- &scrapeResult{Err: ErrHTTPError}
		return
	}
	defer response.Body.Close()

	if httpStatusIsError(response.StatusCode) {
		cResults <- &scrapeResult{Err: ErrHTTPError}
		return
	}

	doc, err := html.Parse(response.Body)
	if err != nil {
		cResults <- &scrapeResult{Err: ErrParseError}
		return
	}

	var parseNextToken func(*html.Node)
	parseNextToken = func(n *html.Node) {
		if nextURL := s.getLinkIfExistsInNode(n); nextURL != "" {
			page.Pages = appendPageIfNotPresent(page.Pages, nextURL)
		} else if asset := s.getAssetIfExistsInNode(n); asset != nil {
			page.Assets = append(page.Assets, asset)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			parseNextToken(c)
		}
	}
	parseNextToken(doc)

	// TODO do we need to remove duplicate assets from a page?
	// That could be a neat feature - find when you have duplicate
	// links to a single resource

	select {
	case cResults <- &scrapeResult{NextURLs: page.Pages}:
	case <-ctx.Done():
	}
}

func (s *crawler) getLinkIfExistsInNode(n *html.Node) string {
	if n.Type != html.ElementNode || n.Data != "a" {
		// Skip this node, it's not an <a> tag
		return ""
	}

	ok, href := attr(n, "href")
	if !ok {
		s.logger.Printf("<a> tag appears to have no 'href' attribute")
		return ""
	}

	parsedHref, err := resolveURL(href, s.rootURL)
	if err != nil {
		s.logger.Printf("<a> tag has a href attribute (%s) we can't parse: '%v'", href, err)
		return ""
	}

	if parsedHref.Host != s.rootURL.Host {
		s.logger.Printf("External link will not be followed '%s'", href)
		return ""
	}

	// Ignore query & fragment
	parsedHref.RawQuery = ""
	parsedHref.Fragment = ""

	return parsedHref.String()
}

func (s *crawler) getAssetIfExistsInNode(n *html.Node) *Asset {
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
		if fullURL, err := resolveURL(src, s.rootURL); err == nil {
			return &Asset{Type: assetType, URL: fullURL.String()}
		}
	}
	return nil
}
