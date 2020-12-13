package main

import (
	"BEAT_THA/etaMicroservice/logger"
	"BEAT_THA/protocol"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var l = new(logger.Logger)

func main() {

	// wrap our hello handler function
	http.Handle("/calculate", loggerware(handle))
	http.ListenAndServe(":5500", nil)

}

func handle(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {

		err := errors.New(fmt.Sprintf("invalid request method [%s]", r.Method))
		logErrorAndRespond(w, err, http.StatusMethodNotAllowed)
		return

	}

	if r.Body == nil {

		err := errors.New("missing request body ")
		logErrorAndRespond(w, err, http.StatusBadRequest)
		return

	}

	defer r.Body.Close()

	input := make(chan []byte)
	go func() {

		read, readError := ioutil.ReadAll(r.Body)

		if readError != nil {

			logErrorAndRespond(w, readError, http.StatusBadRequest)
			return

		}

		input <- read

	}()

	var calculate protocol.Calculate
	result := <-protocol.UnmarshalCalculateToJSON(&calculate, input)

	if result.Err != nil {

		logErrorAndRespond(w, result.Err, http.StatusInternalServerError)
		return

	}

	l.LogInfo("Received a valid calculated object:\n" + calculate.ToString())

	// We have received the command. Now it is time to use the
	// external microservices.
	service := Factory(calculate)
	if service == nil {

		err := fmt.Sprintf("Unexpected service [%s]", calculate.Provider)
		logErrorAndRespond(w, errors.New(err), http.StatusInternalServerError)
		return

	}

	output := make(chan []byte)
	for _, serv := range service {

		l.LogDebug(fmt.Sprintf("Will use [%s] endpoint.", serv.ToString()))

		extResponse := <-Call(serv, calculate)
		if extResponse.err != nil {

			logErrorAndRespond(w, extResponse.err, http.StatusInternalServerError)
			return

		}

		go func() {

			marshal, marshalError := json.Marshal(extResponse.resp)
			if marshalError != nil {

				logErrorAndRespond(w, marshalError, http.StatusInternalServerError)
				return

			}

			output <- marshal

		}()

	}

	w.Write(<-output)

}

func logErrorAndRespond(w http.ResponseWriter, getErr error, statusCode int) {

	l.LogError(getErr)
	w.WriteHeader(statusCode)

}

// loggerware can wrap any handler function and will print out the datetime of the request
// as well as the path that the request was made to.
func loggerware(handler http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		requestPath := r.URL
		l.LogInfo("Request Made To: ", requestPath)
		handle(w, r)

	})

}
