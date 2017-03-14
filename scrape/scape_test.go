package scrape_test

import (
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
