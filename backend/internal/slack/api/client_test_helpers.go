package api

import (
	"net/http"
	"net/http/httptest"
)

const (
	testAccessToken = "xoxb-test-token"
	testChannelID   = "C12345678"
	testMessageText = "JPEG analysis completed"
)

type roundTripFunc func(
	req *http.Request,
) (*http.Response, error)

func (f roundTripFunc) RoundTrip(
	req *http.Request,
) (*http.Response, error) {
	return f(req)
}

func newTestClient(server *httptest.Server) *Client {
	client := NewClient(server.Client())
	client.postMessageURL = server.URL

	return client
}
