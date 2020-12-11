package httpUtil

import (
	"bytes"
	"myFirstProject/logger"
	"net/http"
)

var l = new(logger.Logger)

func Get(url string) (*http.Response, error) {

	l.LogInfo("1. Performing Http Get...")
	resp, err := http.Get(url)
	if err != nil {
		l.LogError(err.Error())
	}

	return resp, err

}

func Post(url string, body []byte) (*http.Response, error) {

	l.LogInfo("1. Performing Http Post...")
	resp, err := http.Post(url, "application/json; charset=utf-8", bytes.NewBuffer(body))
	if err != nil {
		l.LogError(err.Error())
	}

	return resp, err

}
