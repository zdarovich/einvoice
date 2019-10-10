package einvoice

import (
	"github.com/zdarovich/einvoice/core/models"
)

const (
	creditTotal = 0.00
)

//SumGroup - InvoiceRow Summary
type sumGroup struct {
	Text        string          `xml:",chardata"`
	InvoiceSum  string          `xml:"InvoiceSum,omitempty"`
	PenaltySum  string          `xml:"PenaltySum,omitempty"` // Viivise summa.
	Addition    *additionRecord `xml:"Addition,omitempty"`   // Element soodustuste ja juurdehindluse kajastamiseks.
	Rounding    string          `xml:"Rounding,omitempty"`
	VATSummary  []vat           `xml:"VAT,omitempty"`
	TotalVATSum string          `xml:"TotalVATSum,omitempty"`
	TotalSum    string          `xml:"TotalSum"`
	TotalToPay  string          `xml:"TotalToPay,omitempty"`
	Currency    string          `xml:"Currency,omitempty"`
}

//InvoiceSumGroup block - information about invoice amounts (invoice amount, balance of consumption, etc.)

func newSumGroup(s *models.SaleDocument, invoiceType string) (*sumGroup, error) {

	//calculate addition for SumGroup

	//calculate VAT
	// vat, err := newInvoiceVat(s)
	// if err != nil {
	// 	return nil, err
	// }

	g := &sumGroup{
		/*
			Arve summa ilma maksudeta. Väärtus
			on võrdne arvel kajastatud elementide
			InvoiceItemGroup/ItemEntry/ItemSum
			summaga. - WRITE TEST TO CHECK THIS
		*/
		InvoiceSum: fourDigit(s.NetTotal),
		//Addition:    addition,
		//VATSummary:  vatArray,
		TotalSum:   twoDigit(s.Total),
		TotalToPay: twoDigit(s.Total),
		Currency:   s.CurrencyCode,
	}

	if s.VatTotal != 0 {
		g.TotalVATSum = twoDigit(s.VatTotal)
	}

	if s.Rounding != 0 {
		g.Rounding = fourDigit(s.Rounding)
	}

	if invoiceType == "CRE" {
		g.TotalToPay = twoDigit(creditTotal) //get invoice type. if CREDIT -> 0.00
	}

	return g, nil
}
