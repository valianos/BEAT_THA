package httpUtil

import (
	"BEAT_THA/etaMicroservice/logger"
	"bytes"
	"fmt"
	"net/http"
)

var l = new(logger.Logger)

func Get(url string) (error, *http.Response) {

	l.LogInfo(fmt.Sprintf("Performing Http [GET] to [%s]", url))
	resp, err := http.Get(url)
	if err != nil {
		l.LogError(err.Error())
	}

	return err, resp

}

func Post(url string, body []byte) (error, *http.Response) {

	l.LogInfo(fmt.Sprintf("Performing Http [POST] to [%s]", url))
	resp, err := http.Post(url, "application/json; charset=utf-8", bytes.NewBuffer(body))
	if err != nil {
		l.LogError(err.Error())
	}

	return err, resp

}
