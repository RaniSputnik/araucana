package scrape_test

import (
	"net/http"
	"testing"

	. "github.com/RaniSputnik/araucana/scrape"
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
	srv := setupStaticFileServer("test/basic", addr)
	defer srv.Close()

	expected := &Sitemap{
		Pages: []*SitemapURL{
			&SitemapURL{
				Loc: "http://localhost:3000/index.html",
			},
		},
	}
	sitemap, err := Site("http://localhost:3000/index.html")

	if err != nil {
		t.Errorf("Expected no error but got '%v'", err)
	}
	ensureSitemapsMatch(t, sitemap, expected)
}

func TestLinksAreAlsoScraped(t *testing.T) {
	addr := ":3000"
	srv := setupStaticFileServer("test/basic2", addr)
	defer srv.Close()

	expected := &Sitemap{
		Pages: []*SitemapURL{
			&SitemapURL{
				Loc: "http://localhost:3000/index.html",
			},
			&SitemapURL{
				Loc: "http://localhost:3000/contact.html",
			},
		},
	}
	sitemap, err := Site("http://localhost:3000/index.html")

	if err != nil {
		t.Errorf("Expected no error but got '%v'", err)
	}
	ensureSitemapsMatch(t, sitemap, expected)
}

func TestExternalLinksAreNotScraped(t *testing.T) {
	addr := ":3000"
	srv := setupStaticFileServer("test/external", addr)
	defer srv.Close()

	expected := &Sitemap{
		Pages: []*SitemapURL{
			&SitemapURL{
				Loc: "http://localhost:3000/index.html",
			},
		},
	}
	sitemap, err := Site("http://localhost:3000/index.html")

	if err != nil {
		t.Errorf("Expected no error but got '%v'", err)
	}
	ensureSitemapsMatch(t, sitemap, expected)
}

func TestHashAndQueryStringAreIgnored(t *testing.T) {
	addr := ":3000"
	srv := setupStaticFileServer("test/basic3", addr)
	defer srv.Close()

	expected := &Sitemap{
		Pages: []*SitemapURL{
			&SitemapURL{
				Loc: "http://localhost:3000/index.html",
			},
			&SitemapURL{
				Loc: "http://localhost:3000/contact.html",
			},
		},
	}
	sitemap, err := Site("http://localhost:3000/index.html")

	if err != nil {
		t.Errorf("Expected no error but got '%v'", err)
	}
	ensureSitemapsMatch(t, sitemap, expected)
}

func TestTargetSite404ResultsInError(t *testing.T) {
	addr := ":3000"
	srv := setupStaticFileServer("test/basic", addr)
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

func setupStaticFileServer(dir string, addr string) *http.Server {
	srv := &http.Server{
		Handler: http.FileServer(http.Dir(dir)),
		Addr:    addr,
	}

	go func(srv *http.Server) {
		srv.ListenAndServe()
	}(srv)

	return srv
}

func ensureSitemapsMatch(t *testing.T, got *Sitemap, expected *Sitemap) {
	// Handle nil cases
	if got == nil {
		if expected != nil {
			t.Fatalf("Expected valid sitemap but got 'nil'")
		} else {
			// Both got & expected are nil
			return
		}
	}

	if len(expected.Pages) != len(got.Pages) {
		t.Errorf("Expected %d url(s), got %d url(s)", len(expected.Pages), len(got.Pages))
	}

	// TODO these comparrissons could be done a lot faster
	// lets just KISS for now

	for _, url := range got.Pages {
		found := false
		for _, test := range expected.Pages {
			if url.Loc == test.Loc {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Got unexpected url '%s'", url.Loc)
		}
	}

	for _, test := range expected.Pages {
		found := false
		for _, url := range got.Pages {
			if url.Loc == test.Loc {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected url '%s' was not found", test.Loc)
		}
	}
}
