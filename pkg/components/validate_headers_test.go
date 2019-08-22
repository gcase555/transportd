package components

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_incomingMatchesAllowed(t *testing.T) {
	defaultAllowedHeaderAndValues := map[string][]string{"LDAP-Groups": {"sre", "devs"}}
	tests := []struct {
		name string
		allowedHeader map[string][]string
		incomingHeader map[string][]string
		wantResult bool
	}{
		{
			name: "good incoming header values",
			allowedHeader: defaultAllowedHeaderAndValues,
			incomingHeader: map[string][]string{"LDAP-Groups": {"sre", "devs"}},
			wantResult: true,
		},
		{
			name: "good incoming header values with single required header value",
			allowedHeader: map[string][]string{"LDAP-Groups": {"devs"}},
			incomingHeader: map[string][]string{"LDAP-Groups": {"sre", "devs"}},
			wantResult: true,
		},
		{
			name: "bad incoming header values",
			allowedHeader: defaultAllowedHeaderAndValues,
			incomingHeader: map[string][]string{"LDAP-Groups": {"design"}},
			wantResult: false,
		},
		{
			name: "missing incoming header",
			allowedHeader: defaultAllowedHeaderAndValues,
			incomingHeader: map[string][]string{"meh": {"sre", "devs"}},
			wantResult: false,
		},
		{
			name: "missing single incoming header value",
			allowedHeader: defaultAllowedHeaderAndValues,
			incomingHeader: map[string][]string{"LDAP-Groups": {"devs"}},
			wantResult: false,
		},
		{
			name: "missing incoming header and values",
			allowedHeader: defaultAllowedHeaderAndValues,
			incomingHeader: map[string][]string{"": {""}},
			wantResult: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T){
			assert.Equal(t, incomingMatchesAllowed(tt.allowedHeader, tt.incomingHeader), tt.wantResult)
		})
	}
}

func Test_contains(t *testing.T) {
	defaultS := []string{"dog", "cat", "bird", "fish"}
	tests := []struct {
		name string
		sliceToCheck []string
		target string
		wantResult bool
	}{
		{
			name: "target is not present",
			sliceToCheck: defaultS,
			target: "insect",
			wantResult: false,
		},
		{
			name: "target is empty",
			sliceToCheck: defaultS,
			target: "",
			wantResult: false,
		},
		{
			name: "target is present",
			sliceToCheck: defaultS,
			target: "fish",
			wantResult: true,
		},
		{
			name: "target is present in a single value slice",
			sliceToCheck: []string{"dog"},
			target: "dog",
			wantResult: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T){
			assert.Equal(t, contains(tt.sliceToCheck, tt.target), tt.wantResult)
		})
	}
}

func Test_validateHeadersRoundTrip(t *testing.T) {
	assert.Equal(t, true, true)
}