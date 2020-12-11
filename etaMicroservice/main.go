package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"myFirstProject/httpUtil"
	"myFirstProject/logger"
	"net/http"
)

// TODO: protocol module maybe

// ========== Inner protocol
type Calculate struct {
	Origin      spot
	Destination spot
	Provider    PROVIDER
}

func (calculate Calculate) toString() string {

	return fmt.Sprintf("Origin: %s\nDestination: %s\nProviderService: %s",
		calculate.Origin.toString(), calculate.Destination.toString(), calculate.Provider)

}

type spot struct {
	Lat float32
	Lng float32
}

func (s spot) toString() string {
	return fmt.Sprintf("lat: %f \t lng: %f", s.Lat, s.Lng)
}

// Create enum-like providerService.
type PROVIDER string

const (
	SERVICE_A   PROVIDER = "ETAServiceA"
	SERVICE_B   PROVIDER = "ETAServiceB"
	UNSPECIFIED PROVIDER = ""
)

type microserviceResponse struct {
	Eta      int32
	Provider PROVIDER
}

// ======= Outer microservices protocol

type externalMicroserviceResponse interface {
	eta() int32
}

type serviceBResponse struct {
	From     string
	To       string
	Distance int32
	Duration int32
}

func (s serviceBResponse) toString() string {
	return fmt.Sprintf("From: %s \tTo: %s\tDistance: %d\tDuration: %d",
		s.From, s.To, s.Distance, s.Duration)
}

func (s serviceBResponse) eta() int32 { return s.Duration }

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

	var calculate Calculate
	unmarshalError, calculated := Calculate.UnmarshalJSON(calculate, read)

	if unmarshalError != nil {

		logErrorAndRespond(w, unmarshalError)
		return

	}

	calculate = *calculated

	l.LogInfo("Received a valid calculated object:\n" + calculate.toString())

	var url []string
	switch calculate.Provider {
	case SERVICE_A:
		url = []string{"http://localhost:8001/eta/calculate"}
	case SERVICE_B:
		url = []string{"http://localhost:8002/calculateETA"}
	default:
		url = []string{"http://localhost:8001/eta/calculate", "http://localhost:8002/calculateETA"}
	}

	l.LogDebug(fmt.Sprintf("Will use [%s] endpoint.", url))

	// TODO: use protocol
	if calculate.Provider == SERVICE_B {

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

		var BResponse serviceBResponse
		unmarshalError, resp := serviceBResponse.UnmarshalJSON(BResponse, read)

		if unmarshalError != nil {

			logErrorAndRespond(w, unmarshalError)
			return

		}

		BResponse = *resp
		l.LogInfo("Received a valid service B response:\n" + resp.toString())

		// Time to respond.
		result := microserviceResponse{
			Eta:      BResponse.Duration,
			Provider: calculate.Provider,
		}

		marshal, marshalError := json.Marshal(result)
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

func (calculate Calculate) UnmarshalJSON(b []byte) (error, *Calculate) {

	// Define a secondary type to avoid ending up with a recursive call to json.Unmarshal
	type calc Calculate
	var c calc = (calc)(calculate)

	err := json.Unmarshal(b, &c)

	if err != nil {
		panic(err)
	}

	if err := c.Provider.IsValid(); err != nil {
		return err, nil
	}

	return nil, (*Calculate)(&c)

}

func (s serviceBResponse) UnmarshalJSON(b []byte) (error, *serviceBResponse) {

	// Define a secondary type to avoid ending up with a recursive call to json.Unmarshal
	type resp serviceBResponse
	var r resp = (resp)(s)

	err := json.Unmarshal(b, &r)

	if err != nil {
		panic(err)
	}

	return nil, (*serviceBResponse)(&r)

}

func (provider PROVIDER) IsValid() error {

	switch provider {
	case SERVICE_A, SERVICE_B, UNSPECIFIED:
		return nil
	}

	return errors.New(fmt.Sprintf("invalid provider service type: %s.", provider))

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
