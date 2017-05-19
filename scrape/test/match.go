package test

import (
	"testing"

	"github.com/RaniSputnik/araucana/scrape"
)

// EnsureSitemapsMatch will check that two given sitemaps are equivelent.
// Uses testing to report any differences between the two sitemaps and their pages / assets.
func EnsureSitemapsMatch(t *testing.T, got *scrape.Sitemap, expected *scrape.Sitemap) {
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
			continue
		}

		// Check the number of assets match
		if len(gotPage.Assets) != len(expectedPage.Assets) {
			t.Errorf("Expected %d asset(s) on page '%s', got %d asset(s)",
				len(expectedPage.Assets), expectedPage.URL, len(gotPage.Assets))
		}
		// Check that no unexpected assets were found
		for _, gotAsset := range gotPage.Assets {
			found := findAssetInList(gotAsset, expectedPage.Assets)
			if found == nil {
				t.Errorf("Got unexpected asset '%v'", gotAsset)
			}
		}
		// Check that all the expected assets were found
		for _, expectedAsset := range expectedPage.Assets {
			found := findAssetInList(expectedAsset, gotPage.Assets)
			if found == nil {
				t.Errorf("Expected asset '%v' was not found", expectedAsset)
			}
		}

		// Check that the number of page links match
		if len(gotPage.Pages) != len(expectedPage.Pages) {
			t.Errorf("Expected %d page link(s) on page '%s', got %d page link(s)",
				len(expectedPage.Pages), expectedPage.URL, len(gotPage.Pages))
		}
		// Check that no unexpected page links were found
		for _, gotPageLink := range gotPage.Pages {
			if !listContainsString(gotPageLink, expectedPage.Pages) {
				t.Errorf("Got unexpected page link '%s'", gotPageLink)
			}
		}
		// Check that all the expected page links were found
		for _, expectedPageLink := range expectedPage.Pages {
			if !listContainsString(expectedPageLink, gotPage.Pages) {
				t.Errorf("Expected page link '%s' was not found", expectedPageLink)
			}
		}
	}
}

func listContainsString(find string, list []string) bool {
	for _, str := range list {
		if str == find {
			return true
		}
	}
	return false
}

func findPageInList(findURL string, pages []*scrape.Page) *scrape.Page {
	for _, page := range pages {
		if page.URL == findURL {
			return page
		}
	}
	return nil
}

func findAssetInList(find *scrape.Asset, assets []*scrape.Asset) *scrape.Asset {
	for _, asset := range assets {
		if asset.URL == find.URL && asset.Type == find.Type {
			return asset
		}
	}
	return nil
}
