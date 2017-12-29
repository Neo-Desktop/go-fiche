package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/spf13/viper"
)

func serveHttp() {
	files := http.FileServer(http.Dir(viper.GetString("output")))
	http.Handle("/", disableDirectoryListing(files))

	log.Printf("Starting embedded http server on port %d\n", viper.GetInt("httpport"))
	http.ListenAndServe(fmt.Sprintf(":%d", viper.GetInt("httpport")), handlers.LoggingHandler(os.Stdout, http.DefaultServeMux))
}

func disableDirectoryListing(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		h.ServeHTTP(w, r)
	})
}
