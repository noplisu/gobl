package it

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// invoiceValidator adds validation checks to invoices which are relevant
// for the region.
type invoiceValidator struct {
	inv *bill.Invoice
}

func validateInvoice(inv *bill.Invoice) error {
	v := &invoiceValidator{inv: inv}
	return v.validate()
}

func (v *invoiceValidator) validate() error {
	inv := v.inv
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Supplier, validation.By(v.supplier)),
		validation.Field(&inv.Customer, validation.By(v.customer)),
	)
}

func (v *invoiceValidator) supplier(value interface{}) error {
	obj, _ := value.(*org.Party)
	if obj == nil {
		return nil
	}

	return validation.ValidateStruct(obj,
		validation.Field(&obj.TaxID,
			validation.Required,
			tax.RequireIdentityCode,
			tax.IdentityTypeIn(TaxIdentityTypeBusiness, TaxIdentityTypeGovernment),
		),
	)
}

func (v *invoiceValidator) customer(value interface{}) error {
	p, _ := value.(*org.Party)
	if p == nil {
		return nil
	}

	// Customers must have a tax ID (PartitaIVA) if they are legal entities like
	// government offices and companies.
	return validation.ValidateStruct(p,
		validation.Field(&p.TaxID,
			validation.When(
				p.TaxID.Country.In(l10n.IT), // if destination is Italian
				validation.Required,
				tax.RequireIdentityCode,
				tax.IdentityTypeIn(TaxIdentityTypeBusiness, TaxIdentityTypeGovernment, TaxIdentityTypeIndividual),
			),
		),
	)
}