package network

import "net/http"

func GetTransport(skipVerifySSL bool) http.RoundTripper {
	return _transports[skipVerifySSL]
}
