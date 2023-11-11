package main

import "net/http"

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("User-Agent", "PostmanRuntime/7.34.0")
		w.Header().Add("Accept", "*/*")
		w.Header().Add("Accept-Encoding", "gzip, deflate, br")
		w.Header().Add("Connection", "keep-alive")
		w.Header().Add("Host", "localhost:2358")
		next.ServeHTTP(w, r)
	})
}
