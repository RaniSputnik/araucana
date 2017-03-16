package scrape_test

import (
	"fmt"
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
	ensureSitemapsMatch(t, sitemap, expected)
}

func TestLinksAreAlsoScraped(t *testing.T) {
	addr := ":3000"
	srv := setupStaticFileServer("test/basic2", addr)
	defer srv.Close()

	expected := &Sitemap{
		Pages: []*Page{
			&Page{
				URL: "http://localhost:3000/index.html",
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
	ensureSitemapsMatch(t, sitemap, expected)
}

func TestExternalLinksAreNotScraped(t *testing.T) {
	addr := ":3000"
	srv := setupStaticFileServer("test/external", addr)
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
	ensureSitemapsMatch(t, sitemap, expected)
}

func TestHashAndQueryStringAreIgnored(t *testing.T) {
	addr := ":3000"
	srv := setupStaticFileServer("test/basic3", addr)
	defer srv.Close()

	expected := &Sitemap{
		Pages: []*Page{
			&Page{
				URL: "http://localhost:3000/index.html",
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
	ensureSitemapsMatch(t, sitemap, expected)
}

func TestAssetReferencesAreIncluded(t *testing.T) {
	addr := ":3000"
	srv := setupStaticFileServer("test/assets", addr)
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
						Type: AssetTypeStylesheet,
						URL:  fmt.Sprintf("%s/index.css", site),
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

	// TODO these comparrissons could be done a lot faster
	// lets just KISS for now

	// Check the number of pages match
	if len(expected.Pages) != len(got.Pages) {
		t.Errorf("Expected %d pages(s), got %d pages(s)", len(expected.Pages), len(got.Pages))
	}

	// Check that no unexpected pages were found
	for _, url := range got.Pages {
		foundPage := findPageInList(url.URL, expected.Pages)
		if foundPage == nil {
			t.Errorf("Got unexpected page '%s'", url.URL)
		}
	}

	// Ensure that the expected pages were all found in the results
	for _, expectedPage := range expected.Pages {
		gotPage := findPageInList(expectedPage.URL, got.Pages)
		if gotPage == nil {
			t.Errorf("Expected page '%s' was not found", expectedPage.URL)
		} else {

			// Check the number of assets match
			if len(gotPage.Assets) != len(expectedPage.Assets) {
				t.Errorf("Expected %d assets(s), got %d assets(s)", len(expectedPage.Assets), len(gotPage.Assets))
			}

			// Check that no unexpected assets were found
			for _, gotAsset := range gotPage.Assets {
				found := findAssetInList(gotAsset.URL, expectedPage.Assets)
				if found == nil {
					t.Errorf("Got unexpected asset '%s'", gotAsset.URL)
				}
			}

			// Check that all the expected assets were found
			for _, expectedAsset := range expectedPage.Assets {
				found := findAssetInList(expectedAsset.URL, gotPage.Assets)
				if found == nil {
					t.Errorf("Expected asset '%s' was not found", expectedAsset.URL)
				}
			}
		}
	}
}

func findPageInList(findURL string, pages []*Page) *Page {
	for _, page := range pages {
		if page.URL == findURL {
			return page
		}
	}
	return nil
}

func findAssetInList(findURL string, assets []*Asset) *Asset {
	for _, asset := range assets {
		if asset.URL == findURL {
			return asset
		}
	}
	return nil
}
