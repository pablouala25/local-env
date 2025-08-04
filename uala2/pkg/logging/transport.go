package logging

import (
	"log/slog"
	"net/http"

	"github.com/Bancar/uala-auth-team-go-dependencies/libs/http/endpoint"
	"github.com/Bancar/uala-auth-team-go-dependencies/libs/logger"
	"github.com/Bancar/uala-auth-team-go-dependencies/libs/metrics"
	"github.com/Bancar/uala-auth-team-go-dependencies/libs/tracing"
)

func Transport(
	name string,
	metricsEnabled, traceEnabled bool,
	endpointPaths []string,
) http.RoundTripper {
	var roundTripper http.RoundTripper

	base := http.DefaultTransport

	roundTripper = logger.NewLoggingTransport(
		logger.WithName(name),
		logger.WithLogger(slog.Default()),
		logger.WithRoundTripper(base),
	)

	if metricsEnabled {
		roundTripper = metrics.NewTransport(
			roundTripper,
			metrics.Default(),
			name,
			endpointPaths,
			metrics.WithTags(endpoint.MetricTags),
		)
	}

	if traceEnabled {
		roundTripper = tracing.NewTransport(roundTripper)
	}

	return roundTripper
}
