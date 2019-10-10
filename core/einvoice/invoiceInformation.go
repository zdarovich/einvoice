package einvoice

import "github.com/zdarovich/einvoice/core/models"

//Information ...
type information struct {
	Text                   string           `xml:",chardata"`
	Type                   *informationType `xml:"Type"`
	FactorContractNumber   string           `xml:"FactorContractNumber,omitempty"` //will not implement
	ContractNumber         string           `xml:"ContractNumber,omitempty"`
	DocumentName           string           `xml:"DocumentName"`
	InvoiceNumber          string           `xml:"InvoiceNumber"`
	InvoiceContentCode     string           `xml:"InvoiceContentCode,omitempty"`
	InvoiceContentText     string           `xml:"InvoiceContentText,omitempty"`
	PaymentReferenceNumber string           `xml:"PaymentReferenceNumber,omitempty"`
	PaymentMethod          string           `xml:"PaymentMethod,omitempty"` //TODO
	InvoiceDate            string           `xml:"InvoiceDate"`
	DueDate                string           `xml:"DueDate,omitempty"`
	PaymentTerm            string           `xml:"PaymentTerm,omitempty"` //TODO
	FineRatePerDay         string           `xml:"FineRatePerDay,omitempty"`
	Period                 *period          `xml:"Period,omitempty"`           //TODO
	InvoiceDeliverer       string           `xml:"InvoiceDeliverer,omitempty"` //TODO
	Extension              []extension      `xml:"Extension,omitempty"`
}

type extension struct {
	Text               string `xml:",chardata"`
	ExtensionID        string `xml:"extensionId,attr"`
	InformationContent string `xml:"InformationContent,omitempty"`
}

//period ...
type period struct {
	PeriodName string `xml:"PeriodName,omitempty"` //TODO
	StartDate  string `xml:"StartDate,omitempty"`  //TODO
	EndDate    string `xml:"EndDate,omitempty"`    //TODO
}

//InformationType ...
type informationType struct {
	Text          string `xml:",chardata"`
	Type          string `xml:"type,attr"`
	SourceInvoice string `xml:"SourceInvoice,omitempty"`
}

type sourceInvoice struct {
	string
}

const (
	creditInvoiceType = "CREDITINVOICE"
)

func newInvoiceInformation(s *models.SaleDocument) (*information, error) {

	dueDate, err := newDueDate(s.Date, s.PaymentDays)
	if err != nil {
		return nil, err
	}

	i := &information{
		Extension:              getExtensions(),
		DocumentName:           getDocumentName(s.Type),
		InvoiceNumber:          s.Number,
		PaymentReferenceNumber: s.ReferenceNumber,
		PaymentMethod:          s.PaymentType,
		InvoiceDate:            s.Date,
		DueDate:                dueDate, // can be calculated from `salesDoc.PaymentDays 	Integer 	In how many days the invoice is due.`
		PaymentTerm:            s.PaymentDays,
		FineRatePerDay:         "", // in Back Office penalty % per day. Can be retrieved from getCustomers
	}

	if s.PaymentType != "" {
		i.PaymentMethod = s.PaymentType
	}

	i.Type = &informationType{Type: "DEB"}
	if s.Type == creditInvoiceType {
		i.Type = &informationType{Type: "CRE"}
	}

	return i, nil

}

func getExtensions() []extension {
	const eakChannel = "eakChannel"
	return []extension{
		extension{
			ExtensionID:        eakChannel,
			InformationContent: "INTERNET_BANK",
		},
		extension{
			ExtensionID:        eakChannel,
			InformationContent: "POST",
		},
		extension{
			ExtensionID:        eakChannel,
			InformationContent: "EMAIL",
		},
		extension{
			ExtensionID:        eakChannel,
			InformationContent: "PORTAL",
		},
		extension{
			ExtensionID:        "SoftwareId",
			InformationContent: "e-invoices 1.0.0",
		},
	}
}

func getDocumentName(documentType string) string {
	switch documentType {
	case "INVWAYBILL":
		return "Arve-saateleht"
	case "CASHINVOICE":
		return "Kassaarve"
	case "WAYBILL":
		return "Saateleht"
	case "PREPAYMENT":
		return "Ettemaksuarve"
	case "OFFER":
		return "Pakkumine"
	case "EXPORTINVOICE":
		return "Eksportarve"
	case "RESERVATION":
		return "Broneering"
	case "CREDITINVOICE":
		return "Kreeditarve"
	case "ORDER":
		return "Tellimus"
	case "INVOICE":
		return "Arve"
	}
	return ""
}
