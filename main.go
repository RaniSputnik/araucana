package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/RaniSputnik/araucana/scrape"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// TODO set address from configuration
	addr := ":8080"
	// TODO set timeouts from configuration
	writeTimeout := 15 * time.Second
	readTimeout := 15 * time.Second

	router := mux.NewRouter()
	router.HandleFunc("/sitemap", sitemapHandler).Methods(http.MethodGet)
	router.HandleFunc("/ping", pingHandler).Methods(http.MethodGet)
	router.HandleFunc("/", indexHandler).Methods(http.MethodGet)

	srv := &http.Server{
		Handler:      wrapGlobalMiddleware(router),
		Addr:         addr,
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
	}

	// TODO configure logger for file output
	log.Printf("Server now listening at '%s'", addr)
	log.Fatal(srv.ListenAndServe())
}

var sitemapHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	site := q.Get("site")

	sitemap, err := scrape.Site(r.Context(), site)
	if err != nil {
		if err == scrape.ErrURLInvalid {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Required parameter 'site' is invalid"))
		} else if err == scrape.ErrHTTPError {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Could not reach '%s'", site)))
		} else {
			log.Printf("Failed to scrape '%s': %v'", site, err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 internal server error"))
		}
		return
	}

	resBytes, err := json.Marshal(sitemap)
	if err != nil {
		// TODO helper method for 500's? OR do we panic and recover in middleware
		log.Printf("Failed to marshall JSON: %v'", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 internal server error"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resBytes)
})

var pingHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
})

var defaultSiteURL = "/sitemap?site=http://ryanloader.me"
var indexHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, defaultSiteURL, http.StatusFound)
})

func wrapGlobalMiddleware(handler http.Handler) http.Handler {
	r := handlers.RecoveryHandler()(handler)
	r = handlers.LoggingHandler(os.Stdout, r)
	return r
}
