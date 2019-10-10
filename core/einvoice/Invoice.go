package einvoice

import (
	"github.com/zdarovich/einvoice/core/models"
	"strconv"
)

//Invoice - Invoice section ...
type Invoice struct {
	Text               string          `xml:",chardata"`
	InvoiceID          string          `xml:"invoiceId,attr"`
	RegNumber          string          `xml:"regNumber,attr"`
	SellerRegnumber    string          `xml:"sellerRegnumber,attr"`
	InvoiceParties     *parties        `xml:"InvoiceParties"`
	InvoiceInformation *information    `xml:"InvoiceInformation"`
	InvoiceSumGroup    *sumGroup       `xml:"InvoiceSumGroup"`
	InvoiceItem        *item           `xml:"InvoiceItem"`
	AttachmentFile     *attachmentFile `xml:"AttachmentFile,omitempty"`
	PaymentInfo        *paymentInfo    `xml:"PaymentInfo"`
}

//TODO: implement all of the parties

//Parties ...
type parties struct {
	Text           string `xml:",chardata"`
	SellerParty    *party `xml:"SellerParty"`
	BuyerParty     *party `xml:"BuyerParty"`
	RecipientParty *party `xml:"RecipientParty,omitempty"` //Kasutatakse juhul, kui arve saajaks ei ole ostja, vaid näiteks mõni muu ettevõtte või raamatupidamisfirma.
	DeliveryParty  *party `xml:"DeliveryParty,omitempty"`  //Kasutatakse juhul, kui kauba või teenuse sihtpunkt on erinev ostja aadressist.
	PayerParty     *party `xml:"PayerParty,omitempty"`     //Kajastab arve eest tasuva osapoole andmeid juhul, kui see ei ühti ostja andmetega.
	FactorParty    *party `xml:"FactorParty,omitempty"`
}

//AttachmentFile ...
type attachmentFile struct {
	Text       string `xml:",chardata"`
	FileName   string `xml:"FileName,omitempty"`
	FileBase64 string `xml:"FileBase64,omitempty"`
	FileSize   string `xml:"FileSize,omitempty"`
}

//NewInvoice ...
func NewInvoice(salesDoc *models.SaleDocument, customer *models.Customer, companyInfo *models.CompanyInfo, vatRate *models.VatRate) (*Invoice, error) {
	//check RegNumber & SellerRegnumber

	buyerParty, err := newBuyer(customer)
	if err != nil {
		return nil, err
	}

	sellerParty, err := newSeller(companyInfo)
	if err != nil {
		return nil, err
	}

	invoiceItem, err := newInvoiceItem(salesDoc.InvoiceRows, salesDoc.CurrencyCode, vatRate)
	if err != nil {
		return nil, err
	}

	invoiceInfo, err := newInvoiceInformation(salesDoc)
	if err != nil {
		return nil, err
	}

	paymentInfo, err := newPaymentInfo(sellerParty, buyerParty, salesDoc, invoiceInfo.Type.Type)
	if err != nil {
		return nil, err
	}

	sumGroup, err := newSumGroup(salesDoc, invoiceInfo.Type.Type)
	if err != nil {
		return nil, err
	}

	i := &Invoice{
		InvoiceID:       strconv.Itoa(salesDoc.ID),
		RegNumber:       buyerParty.RegNumber,  //buyer party RegNumber
		SellerRegnumber: sellerParty.RegNumber, //seller party RegNumber
		InvoiceParties: &parties{
			SellerParty: sellerParty,
			BuyerParty:  buyerParty,
		},
		InvoiceInformation: invoiceInfo,
		InvoiceSumGroup:    sumGroup,
		InvoiceItem:        invoiceItem,
		PaymentInfo:        paymentInfo,
	}
	return i, nil
}
