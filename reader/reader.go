package reader

import (
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
	ret := ""

	resp, err := http.Get(uri)
	if err != nil {
		l.Warn("Read uri failed: ", zap.String("error", err.Error()))
	}

	//close the response
	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}

	return ret, nil
}
