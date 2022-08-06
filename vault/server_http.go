package vault

import (
	"context"
	httptransport "github.com/go-kit/kit/transport/http"
	"net/http"
)

func NewHTTPServer(ctx context.Context, endpoints Endpoints) http.Handler {
	m := http.NewServeMux()
	m.Handle("/hash", httptransport.NewServer(
		endpoints.HashEndpoint, decodeHashReq, encodeRes,
	))
	m.Handle("/validate", httptransport.NewServer(
		endpoints.ValidateEndpoint, decodeValidateReq, encodeRes,
	))
	return m
}
