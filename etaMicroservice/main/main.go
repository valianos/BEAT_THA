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
		logErrorAndRespond(w, err, protocol.METHOD_NOT_ALLOWED)
		return

	}

	if r.Body != nil {
		defer r.Body.Close()
	}

	read, readError := ioutil.ReadAll(r.Body)

	if readError != nil {

		logErrorAndRespond(w, readError, protocol.BAD_REQUEST)
		return

	}

	var calculate protocol.Calculate
	// TODO: pass pointer here, so that l58 is not necessary.
	unmarshalError, calculated := protocol.Calculate.UnmarshalJSON(calculate, read)

	if unmarshalError != nil {

		logErrorAndRespond(w, unmarshalError, protocol.BAD_REQUEST)
		return

	}

	calculate = *calculated

	l.LogInfo("Received a valid calculated object:\n" + calculate.ToString())

	// We have received the command. Now it is time to use the
	// external microservices.
	service := Factory(calculate)
	if service == nil {

		err := fmt.Sprintf("Unexpected service [%s]", calculate.Provider)
		logErrorAndRespond(w, errors.New(err), protocol.SERVER_ERROR)
		return

	}

	l.LogDebug(fmt.Sprintf("Will use [%s] endpoint.", service.ToString()))

	extError, response := Call(service, calculate)
	if extError != nil {

		logErrorAndRespond(w, extError, protocol.SERVER_ERROR)
		return

	}

	marshal, marshalError := json.Marshal(response)
	if marshalError != nil {

		logErrorAndRespond(w, marshalError, protocol.SERVER_ERROR)
		return

	}

	w.Write(marshal)

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
