package main

import (
	"net/http"
)

func main() {
	http.Handle("/hello", http.HandlerFunc(handle))
}

func handle(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("world"))
}
