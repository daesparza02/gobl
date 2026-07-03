package fi_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeTaxIdentity(t *testing.T) {
	tests := []struct {
		Code     cbc.Code
		Expected cbc.Code
	}{
		{
			Code:     "01120389",
			Expected: "01120389",
		},
		{
			Code:     "FI01120389",
			Expected: "01120389",
		},
		{
			Code:     "fi 0112038-9",
			Expected: "01120389",
		},
		{
			Code:     " 1234567-1 ",
			Expected: "12345671",
		},
	}
	for _, ts := range tests {
		tID := &tax.Identity{Country: "FI", Code: ts.Code}
		norm.Normalize(tID)
		assert.Equal(t, ts.Expected, tID.Code)
	}
}

func TestTaxIdentityRules(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{
			name: "empty",
			code: "",
		},
		{
			// Nokia Oyj's publicly listed Business ID (0112038-9).
			name: "valid",
			code: "01120389",
		},
		{
			name: "valid 2",
			code: "12345671",
		},
		{
			// Weighted sum is a multiple of 11, so the check digit is 0.
			name: "valid with zero check digit",
			code: "00000190",
		},
		{
			name: "too short",
			code: "1234567",
			err:  "IDENTITY-01",
		},
		{
			name: "too long",
			code: "123456789",
			err:  "IDENTITY-01",
		},
		{
			name: "not normalized",
			code: "1234567-1",
			err:  "IDENTITY-01",
		},
		{
			name: "non numeric",
			code: "123456A1",
			err:  "IDENTITY-01",
		},
		{
			name: "non numeric check digit",
			code: "1234567A",
			err:  "IDENTITY-01",
		},
		{
			name: "checksum mismatch",
			code: "12345678",
			err:  "IDENTITY-01",
		},
		{
			// The weighted sum leaves a remainder of 1, which is never
			// assigned: no valid Business ID produces it.
			name: "remainder one",
			code: "00000060",
			err:  "IDENTITY-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "FI", Code: tt.code}
			err := rules.Validate(tID)
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
