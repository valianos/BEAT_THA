package main

import (
	"BEAT_THA/etaMicroservice/httpUtil"
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
		logErrorAndRespond(w, err)
		return

	}

	if r.Body != nil {
		defer r.Body.Close()
	}

	read, readError := ioutil.ReadAll(r.Body)

	if readError != nil {

		logErrorAndRespond(w, readError)
		return

	}

	var calculate protocol.Calculate
	// TODO: pass pointer here, so that l58 is not necessary.
	unmarshalError, calculated := protocol.Calculate.UnmarshalJSON(calculate, read)

	if unmarshalError != nil {

		logErrorAndRespond(w, unmarshalError)
		return

	}

	calculate = *calculated

	l.LogInfo("Received a valid calculated object:\n" + calculate.ToString())

	var url []string
	switch calculate.Provider {
	case protocol.SERVICE_A:
		url = []string{protocol.SERVICE_A_URL}
	case protocol.SERVICE_B:
		url = []string{protocol.SERVICE_B_URL}
	case protocol.UNSPECIFIED:
		url = []string{protocol.SERVICE_A_URL, protocol.SERVICE_B_URL}
	}

	l.LogDebug(fmt.Sprintf("Will use [%s] endpoint.", url))

	if calculate.Provider == protocol.SERVICE_B {

		// We should make a get request to the appropriate url
		api := fmt.Sprintf("%s?from=%f|%f&to=%f|%f",
			url[0],
			calculate.Origin.Lat, calculate.Origin.Lng,
			calculate.Destination.Lat, calculate.Destination.Lng)

		response, getErr := httpUtil.Get(api) // TODO:consider buffers ,this could ease which response comes first

		if getErr != nil {

			logErrorAndRespond(w, getErr)
			return

		}

		read, readError := ioutil.ReadAll(response.Body)

		if readError != nil {

			logErrorAndRespond(w, readError)
			return

		}

		var BResponse protocol.ServiceBResponse
		unmarshalError, resp := protocol.ServiceBResponse.UnmarshalJSON(BResponse, read)

		if unmarshalError != nil {

			logErrorAndRespond(w, unmarshalError)
			return

		}

		BResponse = *resp
		l.LogInfo("Received a valid service B response:\n" + BResponse.ToString())

		// Time to respond.
		result := protocol.MicroserviceResponse{
			Eta:      BResponse.Duration,
			Provider: calculate.Provider,
		}

		marshal, marshalError := json.Marshal(result)
		if marshalError != nil {

			logErrorAndRespond(w, marshalError)
			return

		}

		w.Write(marshal)

	} else if calculate.Provider == protocol.SERVICE_A {

		api := url[0]
		request := protocol.ServiceARequest{
			Origin:      protocol.Spot{Lat: calculate.Origin.Lat, Lng: calculate.Origin.Lng},
			Destination: protocol.Spot{Lat: calculate.Destination.Lat, Lng: calculate.Destination.Lng},
		}

		marshal, marshalError := json.Marshal(request)
		if marshalError != nil {

			logErrorAndRespond(w, marshalError)
			return

		}

		response, postError := httpUtil.Post(api, marshal)
		if postError != nil {

			logErrorAndRespond(w, postError)
			return

		}

		read, readError := ioutil.ReadAll(response.Body)

		if readError != nil {

			logErrorAndRespond(w, readError)
			return

		}

		var AResponse protocol.ServiceAResponse
		unmarshalError, resp := protocol.ServiceAResponse.UnmarshalJSON(AResponse, read)
		if unmarshalError != nil {

			logErrorAndRespond(w, unmarshalError)
			return

		}
		AResponse = *resp
		l.LogInfo("Received a valid service A response:\n" + AResponse.ToString())

		// Time to respond.
		result := protocol.MicroserviceResponse{
			Eta:      AResponse.Duration,
			Provider: calculate.Provider,
		}

		marshal, marshalError = json.Marshal(result)
		if marshalError != nil {

			logErrorAndRespond(w, marshalError)
			return

		}

		w.Write(marshal)

	} else {
		w.Write([]byte("world"))
	}
}

func logErrorAndRespond(w http.ResponseWriter, getErr error) {

	l.LogError(getErr)
	w.WriteHeader(500)

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
