# Notes

Rationale;

* Contain logic in library, easy to plug into web-server / command-line tool etc.
* Keep the scrape library simple to use
* Test business logic, avoid internals
* (Generally) Readability over performance

Codebase improvements;

* Split page scraper from site crawler - two interfaces
* Tests would not need to be end-to-end, could test parsing and crawling separately
* Cap the number of concurrent scrape requests
* Many more security measures for server, idle timeout, max header size, etc.
* Better error format from server + JSON encoding

Preparation for Deployment;

* Log level + more liberal use of logger
* Dynamic Configuration
* Metrics
* Docker Image
* Ping & Healthcheck + diagnostic endpoints
