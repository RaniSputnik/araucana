package scrape_test

import (
	"fmt"
	"testing"

	. "github.com/RaniSputnik/araucana/scrape"
	"github.com/RaniSputnik/araucana/scrape/test"
)

func TestScrapeThrowsErrorWhenURLInvalid(t *testing.T) {
	// For now we only want to ensure the site isn't empty
	// TODO determine other invalid cases
	invalidSites := []string{
		"",
	}

	expected := ErrURLInvalid
	for _, url := range invalidSites {
		_, err := Site(url)
		if err != expected {
			t.Errorf("Expected '%s' but got '%v' for url '%s'", expected, err, url)
		}
	}
}

func TestTargetSiteIsSinglePage(t *testing.T) {
	addr := ":3000"
	srv := test.SetupStaticFileServer("test/basic", addr)
	defer srv.Close()

	expected := &Sitemap{
		Pages: []*Page{
			&Page{
				URL: "http://localhost:3000/index.html",
			},
		},
	}
	sitemap, err := Site("http://localhost:3000/index.html")

	if err != nil {
		t.Errorf("Expected no error but got '%v'", err)
	}
	test.EnsureSitemapsMatch(t, sitemap, expected)
}

func TestLinksAreAlsoScraped(t *testing.T) {
	addr := ":3000"
	srv := test.SetupStaticFileServer("test/basic2", addr)
	defer srv.Close()

	expected := &Sitemap{
		Pages: []*Page{
			&Page{
				URL:   "http://localhost:3000/index.html",
				Pages: []string{"http://localhost:3000/contact.html"},
			},
			&Page{
				URL:   "http://localhost:3000/contact.html",
				Pages: []string{"http://localhost:3000/index.html"},
			},
		},
	}
	sitemap, err := Site("http://localhost:3000/index.html")

	if err != nil {
		t.Errorf("Expected no error but got '%v'", err)
	}
	test.EnsureSitemapsMatch(t, sitemap, expected)
}

func TestExternalLinksAreNotScraped(t *testing.T) {
	addr := ":3000"
	srv := test.SetupStaticFileServer("test/external", addr)
	defer srv.Close()

	expected := &Sitemap{
		Pages: []*Page{
			&Page{
				URL: "http://localhost:3000/index.html",
			},
		},
	}
	sitemap, err := Site("http://localhost:3000/index.html")

	if err != nil {
		t.Errorf("Expected no error but got '%v'", err)
	}
	test.EnsureSitemapsMatch(t, sitemap, expected)
}

func TestHashAndQueryStringAreIgnored(t *testing.T) {
	addr := ":3000"
	srv := test.SetupStaticFileServer("test/basic3", addr)
	defer srv.Close()

	expected := &Sitemap{
		Pages: []*Page{
			&Page{
				URL:   "http://localhost:3000/index.html",
				Pages: []string{"http://localhost:3000/contact.html"},
			},
			&Page{
				URL: "http://localhost:3000/contact.html",
			},
		},
	}
	sitemap, err := Site("http://localhost:3000/index.html")

	if err != nil {
		t.Errorf("Expected no error but got '%v'", err)
	}
	test.EnsureSitemapsMatch(t, sitemap, expected)
}

func TestAssetReferencesAreIncluded(t *testing.T) {
	addr := ":3000"
	srv := test.SetupStaticFileServer("test/assets", addr)
	defer srv.Close()

	site := "http://localhost:3000"
	expected := &Sitemap{
		Pages: []*Page{
			&Page{
				URL: fmt.Sprintf("%s/index.html", site),
				Assets: []*Asset{
					&Asset{
						Type: AssetTypeImage,
						URL:  fmt.Sprintf("%s/hello-world.jpg", site),
					},
					&Asset{
						Type: AssetTypeLink,
						URL:  fmt.Sprintf("%s/index.css", site),
					},
					&Asset{
						Type: AssetTypeLink,
						URL:  fmt.Sprintf("%s/favicon.ico", site),
					},
					&Asset{
						Type: AssetTypeScript,
						URL:  fmt.Sprintf("%s/index.js", site),
					},
				},
			},
		},
	}
	sitemap, err := Site(fmt.Sprintf("%s/index.html", site))

	if err != nil {
		t.Errorf("Expected no error but got '%v'", err)
	}
	test.EnsureSitemapsMatch(t, sitemap, expected)
}

func TestTargetSite404ResultsInError(t *testing.T) {
	addr := ":3000"
	srv := test.SetupStaticFileServer("test/basic", addr)
	defer srv.Close()

	_, err := Site("http://localhost:3000/doesnotexist.html")

	if err != ErrHTTPError {
		t.Errorf("Expected '%v' but got '%v'", ErrHTTPError, err)
	}
}

func TestConnectionTroubleResultsInError(t *testing.T) {
	// Don't setup file server this time

	_, err := Site("http://localhost:9999")
	if err != ErrHTTPError {
		t.Errorf("Expected '%v' but got '%v'", ErrHTTPError, err)
	}
}
