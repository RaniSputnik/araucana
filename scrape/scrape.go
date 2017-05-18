package scrape

import (
	"io"
	"net/url"

	"golang.org/x/net/html"
)

type Scraper interface {
	Scrape(body io.Reader, page *Page) error
}

type ScraperFunc func(body io.Reader, page *Page) error

func (f ScraperFunc) Scrape(body io.Reader, page *Page) error {
	return f(body, page)
}

var DefaultScraperFunc = ScraperFunc(func(body io.Reader, page *Page) error {
	doc, err := html.Parse(body)
	if err != nil {
		return err
	}

	rootURL, err := url.Parse(page.URL)
	if err != nil {
		return err
	}

	var parseNextToken func(*html.Node)
	parseNextToken = func(n *html.Node) {
		if nextURL := getLinkIfExistsInNode(n, rootURL); nextURL != "" {
			page.Pages = appendPageIfNotPresent(page.Pages, nextURL)
		} else if asset := getAssetIfExistsInNode(n, rootURL); asset != nil {
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

	return nil
})

func getLinkIfExistsInNode(n *html.Node, rootURL *url.URL) string {
	if n.Type != html.ElementNode || n.Data != "a" {
		// Skip this node, it's not an <a> tag
		return ""
	}

	ok, href := attr(n, "href")
	if !ok {
		//s.logger.Printf("<a> tag appears to have no 'href' attribute")
		return ""
	}

	parsedHref, err := resolveURL(href, rootURL)
	if err != nil {
		//s.logger.Printf("<a> tag has a href attribute (%s) we can't parse: '%v'", href, err)
		return ""
	}

	if parsedHref.Host != rootURL.Host {
		//s.logger.Printf("External link will not be followed '%s'", href)
		return ""
	}

	// Ignore query & fragment
	parsedHref.RawQuery = ""
	parsedHref.Fragment = ""

	return parsedHref.String()
}

func getAssetIfExistsInNode(n *html.Node, rootURL *url.URL) *Asset {
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
		if fullURL, err := resolveURL(src, rootURL); err == nil {
			return &Asset{Type: assetType, URL: fullURL.String()}
		}
	}
	return nil
}
