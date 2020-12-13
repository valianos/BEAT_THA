package main

import (
	"BEAT_THA/protocol"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testCase struct {
	requestBody protocol.Calculate
	expected    string
}

func Test_basic_functionality(t *testing.T) {

	cases := []testCase{
		{
			requestBody: protocol.Calculate{
				Origin: protocol.Spot{
					Lat: 37.0816818,
					Lng: 23.5035676,
				},
				Destination: protocol.Spot{
					Lat: 37.5142881,
					Lng: 22.5093049,
				},
				Provider: "ETAServiceB"},
			expected: `{"Eta":16707,"Provider":"ETAServiceB"}`,
		}, {
			requestBody: protocol.Calculate{
				Origin: protocol.Spot{
					Lat: 37.0816818,
					Lng: 23.5035676,
				},
				Destination: protocol.Spot{
					Lat: 37.5142881,
					Lng: 22.5093049,
				},
				Provider: "ETAServiceA"},
			expected: `{"Eta":20048,"Provider":"ETAServiceA"}`,
		}, {
			requestBody: protocol.Calculate{
				Origin: protocol.Spot{
					Lat: 37.0816818,
					Lng: 23.5035676,
				},
				Destination: protocol.Spot{
					Lat: 37.5142881,
					Lng: 22.5093049,
				},
				Provider: ""},
			expected: `{"Eta":20048,"Provider":"ETAServiceA"}`, // Service A responds quicker.
		},
	}

	for _, c := range cases {

		marshal, marshalError := json.Marshal(c.requestBody)
		if marshalError != nil {
			t.Fatal(marshalError)
		}

		req, err := http.NewRequest("POST", "/calculate", bytes.NewBuffer(marshal))
		if err != nil {
			t.Fatal(err)
		}

		recorder := httptest.NewRecorder()
		handlerFunc := http.HandlerFunc(handle)
		handlerFunc.ServeHTTP(recorder, req)

		if status := recorder.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		// Check the response body is what we expect.
		if recorder.Body.String() != c.expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				recorder.Body.String(), c.expected)
		}
	}

}
