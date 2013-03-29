package up

import (
	"net/http"
)

// -----------------------------------------------------------------------------

type Transport struct {
	upToken string
	
	transport http.RoundTripper
}

func NewTransport(upToken string, transport http.RoundTripper) *Transport {
	if transport == nil {
		transport = http.DefaultTransport
	}
	return &Transport{upToken, transport}
}

func (t *Transport) RoundTrip(req *http.Request) (resp *http.Response, err error){
	req.Header.Set("Authorization", "UpToken " + t.upToken)
	return t.transport.RoundTrip(req)
}

// -----------------------------------------------------------------------------

func NewClient(upToken string, transport http.RoundTripper) *http.Client {
	t := NewTransport(upToken, transport)
	client := &http.Client{Transport: t}
	return client
}
