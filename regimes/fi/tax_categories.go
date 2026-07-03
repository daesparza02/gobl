package fi

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

var taxCategories = []*tax.CategoryDef{
	{
		Code: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "VAT",
			i18n.FI: "ALV",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
			i18n.FI: "Arvonlisävero",
		},
		Retained: false,
		Keys:     tax.GlobalVATKeys(),
		Rates: []*tax.RateDef{
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateGeneral,
				Name: i18n.String{
					i18n.EN: "Standard Rate",
					i18n.FI: "Yleinen verokanta",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2024, 9, 1),
						Percent: num.MakePercentage(255, 3),
					},
					{
						Since:   cal.NewDate(2013, 1, 1),
						Percent: num.MakePercentage(240, 3),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate",
					i18n.FI: "Alennettu verokanta",
				},
				// Food, restaurant and catering services, books, pharmaceuticals,
				// passenger transport, accommodation, and cultural/sporting events.
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(135, 3),
					},
				},
			},
			{
				Keys: []cbc.Key{tax.KeyStandard},
				Rate: tax.RateSuperReduced,
				Name: i18n.String{
					i18n.EN: "Reduced Rate for Publications",
					i18n.FI: "Lehtien verokanta",
				},
				// Newspapers and magazines (print and electronic).
				Values: []*tax.RateValueDef{
					{
						Percent: num.MakePercentage(100, 3),
					},
				},
			},
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Finnish Tax Administration (Vero) - Rates of VAT",
				},
				URL: "https://www.vero.fi/en/businesses-and-corporations/taxes-and-charges/vat/rates-of-vat/",
				At:  cal.NewDateTime(2026, 7, 1, 0, 0, 0),
			},
		},
	},
}
