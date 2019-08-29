package components

import (
	"context"
	"fmt"
	"net/http"
)

type validateHeaderTransport struct {
	Wrapped  http.RoundTripper
	Required map[string][]string
}

func contains(s []string, target string) bool {
	for _, c := range s {
		if target == c {
			return true
		}
	}
	return false
}

func incomingMatchesRequired(allowed map[string][]string, incoming map[string][]string) error {
	checkedHeaderAndValues := make(map[string][]string)
	for requiredKey, requiredValues := range allowed {
		if matchedIncomingHeaderValues, present := incoming[http.CanonicalHeaderKey(requiredKey)]; present {
			for _, item := range requiredValues {
				checkedHeaderAndValues[requiredKey] = append(checkedHeaderAndValues[requiredKey], item)
				if contains(matchedIncomingHeaderValues, item) {
					return nil
				}
			}
		}
	}
	return fmt.Errorf("no matching values for headers: %s. given matching headers and values: %s", allowed, checkedHeaderAndValues)
}

func (r *validateHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	err := incomingMatchesRequired(r.Required, req.Header)
	if err != nil {
		return newError(http.StatusBadRequest, fmt.Sprintf("header validation failed due to: %s", err)), nil
	}
	return r.Wrapped.RoundTrip(req)
}

// ValidateHeaderConfig is used to validate a map of headers and their required values against an incoming requests headers
type ValidateHeaderConfig struct {
	Required map[string][]string `description:"Map of headers to validate and the required values for those headers."`
}

// Name of the config root
func (*ValidateHeaderConfig) Name() string {
	return "validateheaders"
}

// ValidateHeaderConfigComponent is a plugin
type ValidateHeaderConfigComponent struct{}

// ValidateHeaders satisfies the NewComponent signature
func ValidateHeaders(_ context.Context, _ string, _ string, _ string) (interface{}, error) {
	return &ValidateHeaderConfigComponent{}, nil
}

// Settings generates a config populated with defaults.
func (*ValidateHeaderConfigComponent) Settings() *ValidateHeaderConfig {
	return &ValidateHeaderConfig{}
}

// New generates the middleware.
func (*ValidateHeaderConfigComponent) New(_ context.Context, conf *ValidateHeaderConfig) (func(tripper http.RoundTripper) http.RoundTripper, error) {
	return func(wrapped http.RoundTripper) http.RoundTripper {
		return &validateHeaderTransport{
			Wrapped:  wrapped,
			Required: conf.Required,
		}
	}, nil
}
