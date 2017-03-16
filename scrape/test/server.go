package test

import "net/http"

// SetupStaticFileServer will create a new server that will serve
// the given directory on the given address. The server runs on addr
// separate goroutine, so do not forget to defer Server.Close()
func SetupStaticFileServer(dir string, addr string) *http.Server {
	srv := &http.Server{
		Handler: http.FileServer(http.Dir(dir)),
		Addr:    addr,
	}

	go func(srv *http.Server) {
		srv.ListenAndServe()
	}(srv)

	return srv
}
