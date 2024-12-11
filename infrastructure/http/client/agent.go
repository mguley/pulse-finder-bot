package client

import "net/http"

// RoundTripWithUserAgent allows to set a User-Agent header.
type RoundTripWithUserAgent struct {
	rt        http.RoundTripper
	userAgent string
}

// RoundTrip executes a single HTTP transaction and sets the User-Agent header.
func (r *RoundTripWithUserAgent) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", r.userAgent)
	return r.rt.RoundTrip(req)
}
