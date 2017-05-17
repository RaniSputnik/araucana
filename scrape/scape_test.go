package scrape_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/RaniSputnik/araucana/scrape"
	"github.com/RaniSputnik/araucana/scrape/test"
)

func TestScrapeThrowsErrorWhenURLInvalid(t *testing.T) {
	// For now we only want to ensure the site isn't empty
	invalidSites := []string{
		"",
	}

	expected := ErrURLInvalid
	for _, url := range invalidSites {
		_, err := Site(context.Background(), url)
		if err != expected {
			t.Errorf("Expected '%s' but got '%v' for url '%s'", expected, err, url)
		}
	}
}

func TestTargetSiteIsSinglePage(t *testing.T) {
	srv := setupStaticFileServer("test/basic")
	defer srv.Close()

	expected := &Sitemap{
		Pages: []*Page{
			&Page{
				URL: srv.URL + "/index.html",
			},
		},
	}
	sitemap, err := Site(context.Background(), srv.URL+"/index.html")

	if err != nil {
		t.Errorf("Expected no error but got '%v'", err)
	}
	test.EnsureSitemapsMatch(t, sitemap, expected)
}

func TestLinksAreAlsoScraped(t *testing.T) {
	srv := setupStaticFileServer("test/basic2")
	defer srv.Close()

	expected := &Sitemap{
		Pages: []*Page{
			&Page{
				URL:   srv.URL + "/index.html",
				Pages: []string{srv.URL + "/contact.html"},
			},
			&Page{
				URL:   srv.URL + "/contact.html",
				Pages: []string{srv.URL + "/index.html"},
			},
		},
	}
	sitemap, err := Site(context.Background(), srv.URL+"/index.html")

	if err != nil {
		t.Errorf("Expected no error but got '%v'", err)
	}
	test.EnsureSitemapsMatch(t, sitemap, expected)
}

func TestExternalLinksAreNotScraped(t *testing.T) {
	srv := setupStaticFileServer("test/external")
	defer srv.Close()

	expected := &Sitemap{
		Pages: []*Page{
			&Page{
				URL: srv.URL + "/index.html",
			},
		},
	}
	sitemap, err := Site(context.Background(), srv.URL+"/index.html")

	if err != nil {
		t.Errorf("Expected no error but got '%v'", err)
	}
	test.EnsureSitemapsMatch(t, sitemap, expected)
}

func TestHashAndQueryStringAreIgnored(t *testing.T) {
	srv := setupStaticFileServer("test/basic3")
	defer srv.Close()

	expected := &Sitemap{
		Pages: []*Page{
			&Page{
				URL:   srv.URL + "/index.html",
				Pages: []string{srv.URL + "/contact.html"},
			},
			&Page{
				URL: srv.URL + "/contact.html",
			},
		},
	}
	sitemap, err := Site(context.Background(), srv.URL+"/index.html")

	if err != nil {
		t.Errorf("Expected no error but got '%v'", err)
	}
	test.EnsureSitemapsMatch(t, sitemap, expected)
}

func TestAssetReferencesAreIncluded(t *testing.T) {
	srv := setupStaticFileServer("test/assets")
	defer srv.Close()

	expected := &Sitemap{
		Pages: []*Page{
			&Page{
				URL: fmt.Sprintf("%s/index.html", srv.URL),
				Assets: []*Asset{
					&Asset{
						Type: AssetTypeImage,
						URL:  fmt.Sprintf("%s/hello-world.jpg", srv.URL),
					},
					&Asset{
						Type: AssetTypeLink,
						URL:  fmt.Sprintf("%s/index.css", srv.URL),
					},
					&Asset{
						Type: AssetTypeLink,
						URL:  fmt.Sprintf("%s/favicon.ico", srv.URL),
					},
					&Asset{
						Type: AssetTypeScript,
						URL:  fmt.Sprintf("%s/index.js", srv.URL),
					},
				},
			},
		},
	}
	sitemap, err := Site(context.Background(), fmt.Sprintf("%s/index.html", srv.URL))

	if err != nil {
		t.Errorf("Expected no error but got '%v'", err)
	}
	test.EnsureSitemapsMatch(t, sitemap, expected)
}

func TestTargetSite404ResultsInError(t *testing.T) {
	srv := setupStaticFileServer("test/basic")
	defer srv.Close()

	_, err := Site(context.Background(), srv.URL+"/doesnotexist.html")

	if err != ErrHTTPError {
		t.Errorf("Expected '%v' but got '%v'", ErrHTTPError, err)
	}
}

func TestConnectionTroubleResultsInError(t *testing.T) {
	// Don't setup file server this time

	_, err := Site(context.Background(), "http://localhost:9999")
	if err != ErrHTTPError {
		t.Errorf("Expected '%v' but got '%v'", ErrHTTPError, err)
	}
}

func setupStaticFileServer(dir string) *httptest.Server {
	return httptest.NewServer(http.FileServer(http.Dir(dir)))
}
