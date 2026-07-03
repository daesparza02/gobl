package fi

import (
	"fmt"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// Section 209e of the Finnish Value Added Tax Act and the Business
// Information Act (244/2001) require the seller to be identified on
// invoices, normally with the Business ID or the VAT number derived
// from it.
// Source: https://www.vero.fi/en/detailed-guidance/guidance/48090/vat-invoice-requirements/
func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(
			is.InContext(tax.RegimeIn(l10n.FI.Tax())),
			rules.Field("supplier",
				rules.Assert("01", fmt.Sprintf("invoice FI supplier must have either tax ID code or identity with '%s' type", IdentityTypeYTunnus),
					is.Func(
						fmt.Sprintf("has tax ID code or identity with '%s' type", IdentityTypeYTunnus),
						hasSupplierTaxIDOrIdentity,
					),
				),
			),
		),
	)
}

func hasSupplierTaxIDOrIdentity(value any) bool {
	party, _ := value.(*org.Party)
	return hasTaxIDCode(party) || hasIdentityYTunnus(party)
}

func hasTaxIDCode(party *org.Party) bool {
	return party != nil && party.TaxID != nil && party.TaxID.Code != ""
}

func hasIdentityYTunnus(party *org.Party) bool {
	if party == nil || len(party.Identities) == 0 {
		return false
	}
	return org.IdentityForType(party.Identities, IdentityTypeYTunnus) != nil
}
