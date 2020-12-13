package main

import (
	"BEAT_THA/etaMicroservice/httpUtil"
	"BEAT_THA/protocol"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ExtEtaMicroService interface {
	url() string
	method() protocol.METHOD
	toResponse([]byte) (error, *protocol.MicroserviceResponse)
	performRequest(calculate protocol.Calculate) (error, *http.Response)
	ToString() string
}

// Factory method for creating the two available external microservices.
func Factory(calculate protocol.Calculate) []ExtEtaMicroService {

	switch calculate.Provider {

	case protocol.SERVICE_A:

		return []ExtEtaMicroService{
			ExtEtaMicroserviceA{
				serviceUrl:    protocol.SERVICE_A_URL,
				serviceMethod: protocol.POST,
			}}

	case protocol.SERVICE_B:

		url := fmt.Sprintf("%s?from=%f|%f&to=%f|%f",
			protocol.SERVICE_B_URL,
			calculate.Origin.Lat, calculate.Origin.Lng,
			calculate.Destination.Lat, calculate.Destination.Lng)
		return []ExtEtaMicroService{
			ExtEtaMicroserviceB{
				serviceUrl:    url,
				serviceMethod: protocol.GET,
			}}

	case protocol.UNSPECIFIED:

		url := fmt.Sprintf("%s?from=%f|%f&to=%f|%f",
			protocol.SERVICE_B_URL,
			calculate.Origin.Lat, calculate.Origin.Lng,
			calculate.Destination.Lat, calculate.Destination.Lng)
		return []ExtEtaMicroService{
			ExtEtaMicroserviceA{
				serviceUrl:    protocol.SERVICE_A_URL,
				serviceMethod: protocol.POST,
			},
			ExtEtaMicroserviceB{
				serviceUrl:    url,
				serviceMethod: protocol.GET,
			}}

	default:
		return nil

	}

}

type ExtMicroServiceResponse struct {
	err  error
	resp *protocol.MicroserviceResponse
}

// Exported function for using abstractly the external microservices.
func Call(service ExtEtaMicroService, message protocol.Calculate) <-chan ExtMicroServiceResponse {

	output := make(chan ExtMicroServiceResponse)

	go func() {

		err, response := service.performRequest(message)

		// Fail if something went wrong.
		if err != nil || response == nil {

			output <- ExtMicroServiceResponse{err: err, resp: nil}
			return

		}

		// Read the response
		read, readError := ioutil.ReadAll(response.Body)
		if readError != nil {

			output <- ExtMicroServiceResponse{err: readError, resp: nil}
			return

		}

		// Convert according to protocol.
		convertError, converted := service.toResponse(read)
		if convertError != nil {

			output <- ExtMicroServiceResponse{err: convertError, resp: nil}
			return

		}

		output <- ExtMicroServiceResponse{err: nil, resp: converted}

	}()

	return output

}

// ======= Service implementations follow

// ======= SERVICE A
type ExtEtaMicroserviceA struct {
	serviceUrl    string
	serviceMethod protocol.METHOD
}

func (s ExtEtaMicroserviceA) url() string { return s.serviceUrl }

func (s ExtEtaMicroserviceA) method() protocol.METHOD { return s.serviceMethod }

func (s ExtEtaMicroserviceA) toResponse(body []byte) (error, *protocol.MicroserviceResponse) {

	var AResponse protocol.ServiceAResponse
	unmarshalError, resp := AResponse.UnmarshalJSON(body)
	if unmarshalError != nil {
		return unmarshalError, nil
	}

	result := protocol.MicroserviceResponse{
		Eta:      resp.Duration,
		Provider: protocol.SERVICE_A,
	}

	return nil, &result

}

func (s ExtEtaMicroserviceA) performRequest(message protocol.Calculate) (error, *http.Response) {

	err, bytes := s.body(message)
	if err != nil {
		return err, nil
	}

	return httpUtil.Post(s.url(), bytes)

}

func (s ExtEtaMicroserviceA) ToString() string {
	return fmt.Sprintf("Service A: [{%s}: %s]", s.serviceMethod, s.serviceUrl)
}

func (s ExtEtaMicroserviceA) body(calculate protocol.Calculate) (error, []byte) {

	request := protocol.ServiceARequest{
		Origin:      protocol.Spot{Lat: calculate.Origin.Lat, Lng: calculate.Origin.Lng},
		Destination: protocol.Spot{Lat: calculate.Destination.Lat, Lng: calculate.Destination.Lng},
	}

	marshal, marshalError := json.Marshal(request)
	if marshalError != nil {
		return marshalError, nil
	}

	return nil, marshal

}

// ======= SERVICE B
type ExtEtaMicroserviceB struct {
	serviceUrl    string
	serviceMethod protocol.METHOD
}

func (s ExtEtaMicroserviceB) url() string { return s.serviceUrl }

func (s ExtEtaMicroserviceB) method() protocol.METHOD { return s.serviceMethod }

func (s ExtEtaMicroserviceB) toResponse(body []byte) (error, *protocol.MicroserviceResponse) {

	var BResponse protocol.ServiceBResponse
	unmarshalError, resp := BResponse.UnmarshalJSON(body)
	if unmarshalError != nil {
		return unmarshalError, nil
	}

	result := protocol.MicroserviceResponse{
		Eta:      resp.Duration,
		Provider: protocol.SERVICE_B,
	}

	return nil, &result

}

func (s ExtEtaMicroserviceB) performRequest(message protocol.Calculate) (error, *http.Response) {
	return httpUtil.Get(s.url())
}

func (s ExtEtaMicroserviceB) ToString() string {
	return fmt.Sprintf("Service B: [{%s}:%s]", s.serviceMethod, s.serviceUrl)
}
