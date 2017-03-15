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
		XMLNS: SitemapXMLNamespace,
		URLSet: []*SitemapURL{
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

func TestTargetSite404ResultsInError(t *testing.T) {
	addr := ":3000"
	srv := setupStaticFileServer("test/basic", addr)
	defer srv.Close()

	_, err := Site("http://localhost:3000/doesnotexist.html")

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

	if len(expected.URLSet) != len(got.URLSet) {
		t.Errorf("Expected %d urls, got %d urls", len(expected.URLSet), len(got.URLSet))
	}

	// TODO these comparrissons could be done a lot faster
	// lets just KISS for now

	for _, url := range got.URLSet {
		found := false
		for _, test := range expected.URLSet {
			if url.Loc == test.Loc {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Got unexpected url '%s'", url.Loc)
		}
	}

	for _, test := range expected.URLSet {
		found := false
		for _, url := range got.URLSet {
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
