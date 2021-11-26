package es

import "net/http"

type mockTransport struct {
	roundTripFunc func(*http.Request) (*http.Response, error)
}

func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.roundTripFunc(req)
}
