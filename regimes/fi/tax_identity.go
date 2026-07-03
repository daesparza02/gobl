package fi

import (
	"errors"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// taxCodeLen is the expected length of a normalized Finnish VAT code: the
// seven digits of the Business ID (Y-tunnus) followed by the check digit.
const taxCodeLen = 8

// taxCodeMultipliers are the weights applied to the first seven digits of
// the Business ID to calculate the check digit.
//
// Format source: https://www.ytj.fi/en/index/businessid.html
// Algorithm reference: https://hetut.fi/en/business-id-guide/
var taxCodeMultipliers = []int{7, 9, 10, 5, 8, 4, 2}

func taxIdentityRules() *rules.Set {
	return rules.For(new(tax.Identity),
		rules.When(tax.IdentityIn(CountryCode),
			rules.Field("code",
				rules.AssertIfPresent("01", "invalid Finnish VAT identity code",
					is.Func("valid", isValidTaxIdentityCode),
				),
			),
		),
	)
}

func isValidTaxIdentityCode(value any) bool {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return false
	}
	return validateTaxCode(code) == nil
}

func validateTaxCode(code cbc.Code) error {
	if code == "" {
		return nil
	}
	if len(code) != taxCodeLen {
		return errors.New("invalid length")
	}
	sum := 0
	for i := 0; i < 7; i++ {
		c := code[i]
		if c < '0' || c > '9' {
			return errors.New("invalid characters")
		}
		sum += int(c-'0') * taxCodeMultipliers[i]
	}
	check := code[7]
	if check < '0' || check > '9' {
		return errors.New("invalid characters")
	}

	// The check digit is derived from the remainder of the weighted sum
	// divided by 11: a remainder of 0 means a check digit of 0, a remainder
	// of 1 is never assigned (no valid ID produces it), and any other
	// remainder means a check digit of 11 minus the remainder.
	switch r := sum % 11; r {
	case 1:
		return errors.New("invalid checksum")
	case 0:
		if check != '0' {
			return errors.New("checksum mismatch")
		}
	default:
		if int(check-'0') != 11-r {
			return errors.New("checksum mismatch")
		}
	}
	return nil
}
