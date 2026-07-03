// Package fi provides the tax region definition for Finland.
package fi

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

// CountryCode is the tax country code for Finland.
const CountryCode = "FI"

// init registers the Finnish regime, its validation rules and its
// normalization with GOBL so that they become available as soon as this
// package is imported.
func init() {
	tax.RegisterRegimeDef(New())
	rules.Register("fi", rules.GOBL.Add(CountryCode),
		billInvoiceRules(),
		orgIdentityRules(),
		taxIdentityRules(),
	)
	norm.Register(
		norm.When(tax.IdentityIn(CountryCode), norm.For(func(id *tax.Identity) { tax.NormalizeIdentity(id) })),
	)
}

// New provides the tax region definition for Finland.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   CountryCode,
		Currency:  currency.EUR,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Finland",
			i18n.FI: "Suomi",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Finland's tax system is administered by the Finnish Tax Administration
				(Verohallinto / Vero). As an EU member state, Finland follows the EU VAT
				Directive with locally adapted rates.

				VAT (arvonlisävero, ALV) applies at a standard rate with reduced rates for
				specific goods and services such as food, books, medicines, passenger
				transport, and accommodation. The standard rate rose from 24% to 25.5% on
				1 September 2024.

				Businesses are identified by their Business ID (Y-tunnus) in the format
				"1234567-8" (seven digits, a hyphen and a check digit). The VAT number is
				the same identifier without the hyphen, prefixed with the FI country code
				(e.g. FI12345678).

				Finland supports credit notes for invoice corrections. E-invoicing follows
				the European EN 16931 standard and is commonly exchanged over the Peppol
				network.
			`),
		},
		TimeZone:   "Europe/Helsinki",
		Identities: identityDefinitions,
		Scenarios: []*tax.ScenarioSet{
			bill.InvoiceScenarios(),
		},
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
				},
			},
		},
		Categories: taxCategories,
	}
}
