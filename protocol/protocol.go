package protocol

import (
	"encoding/json"
	"errors"
	"fmt"
)

// ========== Inner protocol
type Calculate struct {
	Origin      Spot     `json:"origin"`
	Destination Spot     `json:"destination"`
	Provider    PROVIDER `json:"provider"`
}

func (calculate Calculate) Validate() error {

	err := calculate.Origin.validate()
	if err == nil {
		err = calculate.Destination.validate()
	}
	if err == nil {
		err = calculate.Provider.validate()
	}

	return err

}

// Friendly representation of the Calculate
func (calculate Calculate) ToString() string {

	return fmt.Sprintf("Origin: %s\nDestination: %s\nProviderService: %s",
		calculate.Origin.toString(), calculate.Destination.toString(), calculate.Provider)

}

// Helper method to JSON unmarshal a Calculate.
func UnmarshalCalculateToJSON(calculate *Calculate, input <-chan []byte) <-chan UnmarshalCalculateResult {

	// This is the outbound channel.
	output := make(chan UnmarshalCalculateResult)

	go func() {

		// Unmarshal channel data.
		err := json.Unmarshal(<-input, &calculate)

		if err != nil {
			output <- UnmarshalCalculateResult{nil, err}
		}

		if err = calculate.Provider.validate(); err != nil {
			output <- UnmarshalCalculateResult{nil, err}
		}

		output <- UnmarshalCalculateResult{calculate, err}

	}()

	return output

}

type Spot struct {
	Lat float32
	Lng float32
}

func (s Spot) validate() error {

	if s.Lat == 0 || s.Lng == 0 {
		return errors.New(fmt.Sprintf("invalid spot: [%s]", s.toString()))
	}
	return nil

}

func (s Spot) toString() string {
	return fmt.Sprintf("lat: %f \t lng: %f", s.Lat, s.Lng)
}

type PROVIDER string

type METHOD string

type MicroserviceResponse struct {
	Eta      int32
	Provider PROVIDER
}

// ======= Outer microservices protocol

type ExternalMicroserviceResponse interface {
	Eta() int32
}

type ServiceBResponse struct {
	From     string
	To       string
	Distance int32
	Duration int32
}

func (s ServiceBResponse) ToString() string {
	return fmt.Sprintf("From: %s \tTo: %s\tDistance: %d\tDuration: %d",
		s.From, s.To, s.Distance, s.Duration)
}

func (s ServiceBResponse) Eta() int32 { return s.Duration }

type ServiceARequest struct {
	Origin      Spot
	Destination Spot
}

func (s ServiceARequest) toString() string {
	return fmt.Sprintf("Origin: %s \tDestination: %s",
		s.Origin.toString(), s.Destination.toString())
}

type ServiceAResponse struct {
	Origin      Spot
	Destination Spot
	Distance    int32
	Duration    int32
}

func (s ServiceAResponse) ToString() string {
	return fmt.Sprintf("From: %s \tTo: %s\tDistance: %d\tDuration: %d",
		s.Origin.toString(), s.Origin.toString(), s.Distance, s.Duration)
}

func (s ServiceAResponse) Eta() int32 { return s.Duration }

type UnmarshalCalculateResult struct {
	Unmarshal *Calculate
	Err       error
}

func (s ServiceBResponse) UnmarshalJSON(b []byte) (error, *ServiceBResponse) {

	// Define a secondary type to avoid ending up with a recursive call to json.Unmarshal
	type resp ServiceBResponse
	var r resp = (resp)(s)

	err := json.Unmarshal(b, &r)

	if err != nil {
		return err, nil
	}

	return nil, (*ServiceBResponse)(&r)

}

func (s ServiceAResponse) UnmarshalJSON(b []byte) (error, *ServiceAResponse) {

	// Define a secondary type to avoid ending up with a recursive call to json.Unmarshal
	type resp ServiceAResponse
	var r resp = (resp)(s)

	err := json.Unmarshal(b, &r)

	if err != nil {
		return err, nil
	}

	return nil, (*ServiceAResponse)(&r)

}

func (provider PROVIDER) validate() error {

	switch provider {
	case SERVICE_A, SERVICE_B, UNSPECIFIED:
		return nil
	}

	return errors.New(fmt.Sprintf("invalid provider service type: %s.", provider))

}
