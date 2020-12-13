package httpUtil

import (
	"BEAT_THA/etaMicroservice/logger"
	"bytes"
	"net/http"
)

var l = new(logger.Logger)

func Get(url string) (error, *http.Response) {

	l.LogInfo("1. Performing Http Get...")
	resp, err := http.Get(url)
	if err != nil {
		l.LogError(err.Error())
	}

	return err, resp

}

func Post(url string, body []byte) (error, *http.Response) {

	l.LogInfo("1. Performing Http Post...")
	resp, err := http.Post(url, "application/json; charset=utf-8", bytes.NewBuffer(body))
	if err != nil {
		l.LogError(err.Error())
	}

	return err, resp

}
