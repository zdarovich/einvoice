package einvoice

import (
	"github.com/zdarovich/einvoice/core/models"
)

const tax = "TAX"

//VAT for ItemEnty & SumGroup
type vat struct {
	Text         string `xml:",chardata"`
	VatID        string `xml:"vatId,attr,omitempty"`   //Käibemaksu tüübi määramiseks.  NOTTAX – mitte-maksustatav käive. TAXEX – maksuvabastus. TAX – maksustatav käive.
	SumBeforeVAT string `xml:"SumBeforeVAT,omitempty"` //Summa, millelt käibemaksu arvutatakse. Decimal4Fraction
	VATRate      string `xml:"VATRate"`
	VATSum       string `xml:"VATSum"`
	Currency     string `xml:"Currency,omitempty"`    //ISO 4217 valuuta
	SumAfterVAT  string `xml:"SumAfterVAT,omitempty"` //Summa käibemaksuga. Decimal4Fraction
}

//note - FinalNetPrice & FinalPriceWith VAT are already DISCOUNTED prices (per item)
func newRowVat(in *models.InvoiceRow, currency string, vatRate string) (*vat, error) {

	//convert in.VatrateID to VatID
	v := &vat{
		VatID:        tax,
		SumBeforeVAT: fourDigit(in.RowNetTotal),
		VATRate:      vatRate,
		VATSum:       twoDigit(in.RowVAT),
		Currency:     currency,
		SumAfterVAT:  fourDigit(in.RowTotal),
	}

	return v, nil
}
