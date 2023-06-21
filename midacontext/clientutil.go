package midacontext

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func newHttpClient(token string) *http.Client {
	client := http.DefaultClient
	rt := withHeader(client.Transport)
	rt.Set("X-Mida-Id", token)
	client.Transport = otelhttp.NewTransport(rt)
	return client
}

type HeaderParam struct {
	http.Header
	rt http.RoundTripper
}

func withHeader(rt http.RoundTripper) HeaderParam {
	if rt == nil {
		rt = http.DefaultTransport
	}
	return HeaderParam{Header: make(http.Header), rt: rt}
}

func (w HeaderParam) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range w.Header {
		req.Header[k] = v
	}
	return w.rt.RoundTrip(req)
}
