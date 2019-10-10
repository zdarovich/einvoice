package einvoice

import (
	"github.com/zdarovich/einvoice/core/models"
	"time"

	"github.com/pkg/errors"
)

//EInvoice ...
type EInvoice struct {
	Text                      string   `xml:",chardata"`
	Xsi                       string   `xml:"xmlns:xsi,attr"`
	NoNamespaceSchemaLocation string   `xml:"xsi:noNamespaceSchemaLocation,attr"`
	Header                    *header  `xml:"Header"`
	Invoices                  Invoices `xml:"Invoice"`
	Footer                    *footer  `xml:"Footer"`
}

const (
	xsi                       = "http://www.w3.org/2001/XMLSchema-instance"
	noNamespaceSchemaLocation = "e-invoice_ver1.2.xsd"
	version                   = "1.2"
	timeFormat                = "2006-01-02"
)

func NewEinvoice(salesDoc *models.SaleDocument, customer *models.Customer, companyInfo *models.CompanyInfo, vatRate *models.VatRate) (*EInvoice, error) {

	ice, err := NewInvoice(salesDoc, customer, companyInfo, vatRate)
	if err != nil {
		return nil, err
	}
	invs := append(Invoices{}, *ice)

	fileID, err := newUUID()
	if err != nil {
		return nil, errors.Wrap(err, "Generation of UUID failed")
	}
	//create e-invoice
	eInvoice := &EInvoice{
		Xsi:                       xsi,
		NoNamespaceSchemaLocation: noNamespaceSchemaLocation,
		Header: &header{
			Date:    time.Now().Local().Format(timeFormat),
			FileID:  fileID, //fmt.Sprint(invoiceID),
			Version: version,
		},
		Invoices: invs,
		Footer: &footer{
			TotalNumberInvoices: "1",
			TotalAmount:         twoDigit(salesDoc.Total),
		},
	}

	return eInvoice, nil
}
