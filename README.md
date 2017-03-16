# Araucana

Tool for creating sitemaps from a given URL. Output format is as follows;

```
pages: [{
	url: http://example.com/
	assets: [{
		type: link|image|script
		url: http://example.com/path-to-asset.png
	}]
	pages: [
		http://example.com/link
	]
}]
```

To run, you will need [go 1.8](https://golang.org/doc/go1.8) or later installed;

```
cd $GOPATH/src/github.com/RaniSputnik/araucana
go build
./araucana
```

Open [localhost:8080/sitemap?site=http://ryanloader.me](http://localhost:8080/sitemap?site=http://ryanloader.me) in your browser to set the results.

Use `go test ./...` to run the tests.

### Brief

Requirements;

* Write a simple web-crawler in a language of your choice
* Accepts an input URL and generates a sitemap for that domain
* Crawler does not follow external links
* Sitemap should include links between pages
* Sitemap should include static assets a page depends on

Bonus points for tests and making it as fast as possible!

### Assumptions

* ~~Timebox to 4hrs~~ Ha!
* Web API to make the results easily consumable
* Limited use of libraries (preferably standard lib only!)
* Fully qualified URL's
* ~~Stick as close as possible to Googles sitemap format~~
* ~~Respect robots.txt~~

### Known Issues

* Sites with an `<a>` tag link to a non-HTML page will result in the url still being added to the list of pages eg. `<a href="index.js">View Source</a>`
* Links that have a trailing slash are counted separate from those without eg. `//localhost:3000/blog` vs `//localhost:3000/blog/`
