package main

import (
	"myFirstProject/logger"
	"net/http"
)

func main() {

	l := new(logger.Logger)

	// wrap our hello handler function
	http.Handle("/hello", loggerware(l, http.HandlerFunc(handle)))
	http.ListenAndServe(":5500", nil)

}

func handle(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("world"))
}

// loggerware can wrap any handler function and will print out the datetime of the request
// as well as the path that the request was made to.
func loggerware(l *logger.Logger, next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestPath := r.URL
		l.LogInfo("Request Made To: ", requestPath)
	})

}
