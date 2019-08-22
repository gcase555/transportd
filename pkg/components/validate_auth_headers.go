package components

import (
	"context"
	"fmt"
	"net/http"
)

type validateAuthHeaderTransport struct {
	Wrapped http.RoundTripper
	Allowed map[string][]string
}

func contains(s []string, target string) bool {
	for _, c := range s {
		if target == c {
			return true
		}
	}
	return false
}

func incomingMatchesAllowed(allowed map[string][]string, incoming map[string][]string) bool{
	allowedValuesFound := false
	for allowedKey, allowedValues := range allowed {
		matchedIncomingHeaderValues := incoming[allowedKey]
		for _, item := range allowedValues {
			allowedValuesFound = contains(matchedIncomingHeaderValues, item)
			if allowedValuesFound == false {
				return allowedValuesFound
			}
		}
	}
	return allowedValuesFound
}

func (r *validateAuthHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	incomingLdapGroups := req.Header["X-Slauth-User-Groups"]
	allowedList := r.Allowed["x-slauth-user-groups"]
	fmt.Println("##### ALLOWED GROUPS #####: ", allowedList)
	fmt.Println("##### ALLOWED #####: ", r.Allowed)
	fmt.Println("##### INCOMING LDAP GROUPS ####: ", incomingLdapGroups)
	if incomingMatchesAllowed(r.Allowed, req.Header) {
		return r.Wrapped.RoundTrip(req)
	}
	return newError(http.StatusUnauthorized, "missing required LDAP group"), nil
}

// AuthConfig is used to configure authorization based on ldap group membership sent in a header
type AuthConfig struct {
	Allowed map[string][]string `description:"List of allowed headers and "`
}

// Name of the config root
func (*AuthConfig) Name() string {
	return "authheaders"
}

// AuthConfigComponent is a plugin
type AuthConfigComponent struct{}

// ValidateAuthHeaders satisfies the NewComponent signature
func ValidateAuthHeaders(_ context.Context, _ string, _ string, _ string) (interface{}, error) {
	return &AuthConfigComponent{}, nil
}

// Settings generates a config populated with defaults.
func (*AuthConfigComponent) Settings() *AuthConfig {
	return &AuthConfig{}
}

// New generates the middleware.
func (*AuthConfigComponent) New(_ context.Context, conf *AuthConfig) (func(tripper http.RoundTripper) http.RoundTripper, error) {
	return func(wrapped http.RoundTripper) http.RoundTripper {
		return &validateAuthHeaderTransport{
			Wrapped: wrapped,
			Allowed: conf.Allowed,
		}
	}, nil
}
