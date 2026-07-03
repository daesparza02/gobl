package fi_test

import (
	"testing"

	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/fi"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestValidateOrgIdentity(t *testing.T) {
	tests := []struct {
		name     string
		identity *org.Identity
		err      string
	}{
		{
			// Nokia Oyj's publicly listed Business ID.
			name: "valid",
			identity: &org.Identity{
				Type: fi.IdentityTypeYTunnus,
				Code: "0112038-9",
			},
		},
		{
			name: "valid 2",
			identity: &org.Identity{
				Type: fi.IdentityTypeYTunnus,
				Code: "1234567-1",
			},
		},
		{
			// Weighted sum is a multiple of 11, so the check digit is 0.
			name: "valid with zero check digit",
			identity: &org.Identity{
				Type: fi.IdentityTypeYTunnus,
				Code: "0000019-0",
			},
		},
		{
			name: "missing hyphen",
			identity: &org.Identity{
				Type: fi.IdentityTypeYTunnus,
				Code: "01120389",
			},
			err: "ORG-IDENTITY-01",
		},
		{
			name: "hyphen in wrong position",
			identity: &org.Identity{
				Type: fi.IdentityTypeYTunnus,
				Code: "011203-89",
			},
			err: "ORG-IDENTITY-01",
		},
		{
			name: "too short",
			identity: &org.Identity{
				Type: fi.IdentityTypeYTunnus,
				Code: "112038-9",
			},
			err: "ORG-IDENTITY-01",
		},
		{
			name: "checksum mismatch",
			identity: &org.Identity{
				Type: fi.IdentityTypeYTunnus,
				Code: "1234567-8",
			},
			err: "ORG-IDENTITY-01",
		},
		{
			// The weighted sum leaves a remainder of 1, which is never
			// assigned: no valid Business ID produces it.
			name: "remainder one",
			identity: &org.Identity{
				Type: fi.IdentityTypeYTunnus,
				Code: "0000006-0",
			},
			err: "ORG-IDENTITY-01",
		},
		{
			name: "empty code",
			identity: &org.Identity{
				Type: fi.IdentityTypeYTunnus,
				Code: "",
			},
			err: "ORG-IDENTITY-01",
		},
		{
			name: "other identity type",
			identity: &org.Identity{
				Type: "OTHER",
				Code: "not checked",
			},
		},
	}

	opts := []rules.WithContext{
		tax.RegimeContext(fi.CountryCode),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rules.Validate(tt.identity, opts...)
			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), tt.err)
				}
			}
		})
	}
}
