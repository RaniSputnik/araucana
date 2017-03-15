# Araucana

Tool for creating sitemaps from a given URL.

### Brief

Requirements;

* Write a simple web-crawler in a language of your choice
* Accepts an input URL and generates a sitemap for that domain
* Crawler does not follow external links
* Sitemap should include links between pages
* Sitemap should include static assets a page depends on

Bonus points for tests and making it as fast as possible!

### Assumptions

* Timebox to 4hrs
* Web API to make the results easily consumable
* Limited use of libraries (preferably standard lib only!)
* Fully qualified URL's
* Stick as close as possible to Googles sitemap format
* Respect robots.txt

### Known Issues

* Sites with an `<a>` tag link to a non-HTML page will result in the url still being added to the list of pages eg. `<a href="index.js">View Source</a>`
