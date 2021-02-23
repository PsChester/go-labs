package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/", router)
	http.ListenAndServe(":80", nil)
}

func router(w http.ResponseWriter, r *http.Request) {
	shortUrls := map[string]string{
		"/go-http":    "https://golang.org/pkg/net/http/",
		"/go-gophers": "https://github.com/shalakhin/gophericons/blob/master/preview.jpg",
	}
	if r.URL.Path == "/go-http" {
		http.Redirect(w, r, shortUrls["/go-http"], http.StatusFound)
		return
	}
	if r.URL.Path == "/go-gophers" {
		http.Redirect(w, r, shortUrls["/go-gophers"], http.StatusFound)
		return
	}
	http.NotFound(w, r)
}
