package transportd

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/asecurityteam/settings"
	"github.com/asecurityteam/transport"
)

const (
	backendsSetting = "backends"
	hostSetting     = "host"
	countSetting    = "count"
	ttlSetting      = "ttl"
	poolSetting     = "pool"
)

// NewBaseTransports generates a mapping of backend names to http.RoundTripper instances.
// This method is used to handle the top-level x-transportd block and configure a set of
// base http.RoundTripper instances with some core connection pooling settings applied.
func NewBaseTransports(ctx context.Context, s settings.Source) (BackendRegistry, error) {
	// Determine what backends are available for endpoints.
	backendsInstalled := settings.NewStringSliceSetting(backendsSetting, "", []string{})
	backendsG := &settings.SettingGroup{
		NameValue:     ExtensionKey,
		SettingValues: []settings.Setting{backendsInstalled},
	}
	err := settings.LoadGroups(ctx, s, []settings.Group{backendsG})
	if err != nil {
		return nil, fmt.Errorf("failed to load backend list: %s", err.Error())
	}
	backends := *backendsInstalled.StringSliceValue

	// Load the base transport for each backend found.
	result := NewStaticBackendRegistry()
	s = &settings.PrefixSource{
		Source: s,
		Prefix: []string{ExtensionKey},
	}
	for _, backend := range backends {
		host := settings.NewStringSetting(hostSetting, "", "")
		poolCount := settings.NewIntSetting(countSetting, "", 1)
		poolTTL := settings.NewDurationSetting(ttlSetting, "", time.Hour)
		pool := &settings.SettingGroup{
			NameValue:     poolSetting,
			SettingValues: []settings.Setting{poolCount, poolTTL},
		}
		g := &settings.SettingGroup{
			NameValue:     backend,
			SettingValues: []settings.Setting{host},
			GroupValues:   []settings.Group{pool},
		}

		err := settings.LoadGroups(ctx, s, []settings.Group{g})
		if err != nil {
			return nil, fmt.Errorf("failed to load backend %s: %s", backend, err.Error())
		}
		hostVal, err := url.Parse(*host.StringValue)
		if err == nil {
			// Only try to validate the content if parsing passed.
			err = validateHost(hostVal)
		}
		if err != nil {
			return nil, fmt.Errorf(
				"failed to parse host %s for backend %s: %s",
				*host.StringValue, backend, err.Error(),
			)
		}
		f := transport.NewFactory(transport.OptionDefaultTransport)
		f = transport.NewRecyclerFactory(
			f,
			transport.RecycleOptionTTL(*poolTTL.DurationValue),
			transport.RecycleOptionTTLJitter(*poolTTL.DurationValue/time.Duration(5)),
		)
		f = transport.NewRotatorFactory(f, transport.RotatorOptionInstances(*poolCount.IntValue))
		result.Store(ctx, backend, &hostRewrite{
			Scheme:  hostVal.Scheme,
			Host:    hostVal.Host,
			Wrapped: f(),
		})
	}
	return result, nil
}

func validateHost(u *url.URL) error {
	if u.Scheme == "" {
		return fmt.Errorf("missing url scheme for %s", u.String())
	}
	if u.Host == "" {
		return fmt.Errorf("missing url host for %s", u.String())
	}
	return nil
}

// hostRewrite is an http.RoundTripper decorator that implements the same logic
// as we'd normally implement as a Directory method for the http.ReverseProxy.
// Structuring this as a decorator comes with the benefits. For one, it prevents
// us from needing to perform a request-time lookup of the matched route, extract
// the relevant extensions, and load the relevant backend data in order to rewrite
// the requests. Instead, we can rely on a static binding that is attached to the
// transport itself. Additionally, this decouples our rewrite logic from the
// ReverseProxy implementation should we ever need to diverge from it.
type hostRewrite struct {
	Scheme  string
	Host    string
	Wrapped http.RoundTripper
}

func (r *hostRewrite) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Host = r.Host
	req.URL.Host = r.Host
	req.URL.Scheme = r.Scheme
	return r.Wrapped.RoundTrip(req)
}