package main

import (
	"context"
	"encoding/json"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/ricardo-ch/go-tracing"
	"os"
)

// this name is use to identify traces inside zipkin
const (
	appName = "my-http-server-middleware"
)

func main() {
	os.Setenv("JAEGER_SERVICE_NAME", appName)
	os.Setenv("JAEGER_AGENT_HOST", "localhost")
	os.Setenv("JAEGER_AGENT_PORT", "6831")
	os.Setenv("JAEGER_SAMPLER_TYPE", "const")
	os.Setenv("JAEGER_SAMPLER_PARAM", "1")

	tracing.SetGlobalTracer()
	defer tracing.FlushCollector()

	svc := stringService{}
	uppercaseHandler := httptransport.NewServer(
		makeUppercaseEndpoint(svc),
		decodeUppercaseRequest,
		encodeResponse,
	)

	http.Handle("/uppercase", tracing.HTTPMiddleware("uppercase-handler", uppercaseHandler))
	http.ListenAndServe(":8000", nil)
}

func decodeUppercaseRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request uppercaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
