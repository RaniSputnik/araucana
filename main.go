package main

import (
	"log"
	"net/http"
	"time"

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
	// TODO the sitemap tool should validate the site parameter
	if site == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Required parameter 'site' is invalid"))
		return
	}

	// TODO generate sitemap
	w.WriteHeader(http.StatusNotImplemented)
})
