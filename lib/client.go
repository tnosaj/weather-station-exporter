package lib

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type HttpClient struct {
	BaseUri string
	client  *http.Client
}

func (c HttpClient) Request(method string, uri string, body io.Reader) (respData []byte, err error) {
	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		return nil, fmt.Errorf("Failed to build %s request to url '%s' with error: %s", method, uri, err)
	}
	if body != nil {
		req.Header = http.Header{
			"Content-Type": []string{"application/json"},
		}
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to do %s request to url '%s' with error: %s", method, uri, err)
	}

	respData, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		if err != nil {
			respData = []byte(err.Error())
		}
		return nil, fmt.Errorf("Request failed with status %s (%d): %s", resp.Status, resp.StatusCode, respData)
	}

	return respData, nil
}

func NewHttpClient(uri string, timeout int) *HttpClient {
	log.Debugf("Starting call to '%s' with timeout: %d", uri, timeout)
	httpClient := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
		},
	}

	return &HttpClient{
		BaseUri: uri,
		client:  httpClient,
	}
}
