package http

import (
	"io"
	"net/http"
	"net/url"

	"github.com/cookieo9/resources-go"
)

type HttpBundle struct {
	BaseURL *url.URL
}

func NewBundle(baseurl string) (resources.Bundle, error) {
	u, err := url.Parse(baseurl)
	if err != nil {
		return nil, err
	}
	return &HttpBundle{BaseURL: u}, nil
}

func (hb *HttpBundle) Open(path string) (io.ReadCloser, error) {
	dest, err := hb.BaseURL.Parse(path)
	if err != nil {
		return nil, err
	}
	response, err := http.Get(dest.String())
	if err != nil {
		return nil, err
	}
	return response.Body, nil
}

func (hb *HttpBundle) Close() error {
	return nil
}
