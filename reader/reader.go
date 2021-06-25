package reader

import (
	"errors"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"
)

type IWebClient interface {
	Get(string, *zap.Logger) (string, error)
}

type WebClient struct {
}

//NewWebClient returns a pointer to a new WebClient instance
func NewWebClient() IWebClient {
	return &WebClient{}
}

//Get makes a http Get request and returns the contents or an error
func (w *WebClient) Get(uri string, l *zap.Logger) (string, error) {
	resp, err := http.Get(uri)
	if err != nil {
		l.Warn("Read uri failed: ", zap.String("uri", uri))
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		l.Warn("Read uri Status Code error ", zap.String("StatusCode", resp.Status))
		return "", errors.New("Invalid Status code")
	}

	//close the response
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		l.Fatal("Could not read body ", zap.Error(err))
	}
	bodyString := string(bodyBytes)

	return bodyString, nil
}
