package fi

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

const (
	// IdentityTypeYTunnus represents the Finnish Business ID (Y-tunnus)
	// issued by the Business Information System (YTJ) to every registered
	// company. It is the same number used for the VAT identity, but
	// presented in its domestic format with a hyphen before the check
	// digit (e.g. "0112038-9").
	IdentityTypeYTunnus cbc.Code = "Y-TUNNUS"
)

// businessIDLen is the length of a Business ID in its domestic format:
// seven digits, a hyphen and a check digit.
const businessIDLen = 9

var identityDefinitions = []*cbc.Definition{
	{
		Code: IdentityTypeYTunnus,
		Name: i18n.String{
			i18n.EN: "Business ID",
			i18n.FI: "Y-tunnus",
			i18n.SV: "FO-nummer",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("YTJ - Business ID"),
				URL:   "https://www.ytj.fi/en/index/businessid.html",
			},
		},
	},
}

func orgIdentityRules() *rules.Set {
	return rules.For(new(org.Identity),
		rules.When(
			is.InContext(tax.RegimeIn(CountryCode)),
			rules.When(
				org.IdentityTypeIn(IdentityTypeYTunnus),
				rules.Field("code",
					rules.Assert("01", "identity code for type Y-TUNNUS must be valid",
						is.Func("valid", isValidBusinessID),
					),
				),
			),
		),
	)
}

// isValidBusinessID expects a Business ID in its domestic format and reuses
// the VAT code validation for the digits and the check digit.
func isValidBusinessID(value any) bool {
	code, ok := value.(cbc.Code)
	if !ok || len(code) != businessIDLen || code[7] != '-' {
		return false
	}
	return validateTaxCode(code[:7]+code[8:]) == nil
}
