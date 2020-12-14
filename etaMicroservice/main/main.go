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

	// Configure the server.
	http.Handle("/calculate", logAndHandle())
	http.ListenAndServe(":8080", nil)

}

func handle(w http.ResponseWriter, r *http.Request) {

	// Only accept HTTP POST.
	if r.Method != "POST" {

		err := errors.New(fmt.Sprintf("invalid request method [%s]", r.Method))
		logErrorAndRespond(w, err, http.StatusMethodNotAllowed)
		return

	}

	// Require a valid body.
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

	// Try to unmarshal the 'calculate' entity from input.
	var calculate protocol.Calculate
	result := <-protocol.UnmarshalCalculateToJSON(&calculate, input)

	if result.Err != nil {

		logErrorAndRespond(w, result.Err, http.StatusBadRequest)
		return

	}

	if err := calculate.Validate(); err != nil {

		logErrorAndRespond(w, err, http.StatusBadRequest)
		return

	}

	// We have received the command. Now it is time to use the
	// external microservices.
	l.LogInfo("Received a valid calculated object:\n" + calculate.ToString())

	// Fetch the appropriate service(s).
	service := Factory(calculate)
	if service == nil {

		err := fmt.Sprintf("Unexpected service [%s]", calculate.Provider)
		logErrorAndRespond(w, errors.New(err), http.StatusInternalServerError)
		return

	}

	// This is the concurrent call to all required external microservices.
	// Note: we will only keep the first response, in case of more services.
	output := make(chan []byte)
	for _, serv := range service {

		l.LogDebug(fmt.Sprintf("Will use [%s] endpoint.", serv.ToString(calculate)))

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

	// Use the first available response, if more than one is due.
	w.Write(<-output)

}

// Helper method that will log the occurred error and fail the request appropriately.
func logErrorAndRespond(w http.ResponseWriter, getErr error, statusCode int) {

	l.LogError(getErr)
	w.WriteHeader(statusCode)

}

// logAndHandle can wrap any handler function and will print out some
// information about the request.
func logAndHandle() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		l.LogInfo(fmt.Sprintf("Received [%s] request from [%s] to: [%s]",
			r.Method, r.RemoteAddr, r.RequestURI))
		handle(w, r)

	}

}
