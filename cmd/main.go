package cmd

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		proxyHttp := httputil.NewSingleHostReverseProxy(&url.URL{})
		proxyHttp.ServeHTTP(writer, request)
	})

	http.ListenAndServe(":8080", mux)
}
