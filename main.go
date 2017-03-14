package main

import (
	"log"
	"net/http"
	"time"

	"encoding/xml"

	"github.com/RaniSputnik/araucana/scrape"
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

	srv := &http.Server{
		Handler:      router,
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

	sitemap, err := scrape.Site(site)
	if err != nil {
		// TODO handle error codes
		// - report when couldn't reach the input site
		if err == scrape.ErrURLInvalid {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Required parameter 'site' is invalid"))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 internal server error"))
		}
		return
	}

	resBytes, err := xml.Marshal(sitemap)
	if err != nil {
		// TODO helper method for 500's? OR do we panic and recover in middleware
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 internal server error"))
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	w.Write(resBytes)
})
