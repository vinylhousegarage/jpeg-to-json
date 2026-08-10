package api

import (
	"net/http"
	"net/http/httptest"
)

const (
	testAccessToken = "xoxb-test-token"
	testChannelID   = "D123456"
	testUserID      = "U123456"
	testMessage     = "test message"
)

func newTestClient(
	server *httptest.Server,
) *Client {
	client := NewClient(
		server.Client(),
	)

	client.postMessageURL = server.URL
	client.openConversationURL = server.URL

	return client
}

func newTestServer(
	handler http.HandlerFunc,
) *httptest.Server {
	return httptest.NewServer(handler)
}
