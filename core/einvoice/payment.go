package einvoice

import (
	"github.com/zdarovich/einvoice/core/models"

	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

//PaymentInfo ...
type paymentInfo struct {
	Text               string `xml:",chardata"`
	Currency           string `xml:"Currency"`
	PaymentRefID       string `xml:"PaymentRefId,omitempty"`
	PaymentDescription string `xml:"PaymentDescription"`
	Payable            string `xml:"Payable"`
	PayDueDate         string `xml:"PayDueDate,omitempty"`
	PaymentTotalSum    string `xml:"PaymentTotalSum"`
	PayerName          string `xml:"PayerName"`
	PaymentID          string `xml:"PaymentId"`
	PayToAccount       string `xml:"PayToAccount,omitempty"`
	PayToName          string `xml:"PayToName"`
}

func newPaymentInfo(seller *party, buyer *party, s *models.SaleDocument, invoiceType string) (*paymentInfo, error) {

	if len(seller.AccountInfo) == 0 {
		return nil, errors.New("Seller's account info was empty")
	}
	//payment may have bankDescription field not empty.
	p := &paymentInfo{
		Currency:     s.CurrencyCode,
		PayerName:    buyer.Name,
		PaymentID:    s.Number,
		PayToAccount: seller.AccountInfo[0].AccountNumber,
		PayToName:    seller.Name,
	}

	if s.ReferenceNumber != "" {
		p.PaymentRefID = s.ReferenceNumber
	} else {
		p.PaymentDescription = fmt.Sprintf("%s number %s", s.Type, s.Number)
	}

	//len(payments) can be 0, or 1, or many.

	if s.Paid == "" {
		p.PaymentTotalSum = twoDigit(s.Total)
	} else {
		paid, err := strconv.ParseFloat(s.Paid, 32)
		if err != nil {
			return nil, errors.Wrap(err, "Could not convert Sales document Paid field")
		}
		p.PaymentTotalSum = twoDigit(s.Total - paid)
	}

	//NO – arve ei kuulu tasumisele (PayDueDate täitmine ei olekohustuslik). See on kui salesDoc.Paid = salesDoc.TotalSum*/
	p.Payable = "NO"

	//check if payable
	if invoiceType == "DEB" {
		//PAYABLE: YES – arve kuulub tasumisele. (PayDueDate täitmine on kohustuslik).
		p.Payable = "YES"
		dueDate, err := newDueDate(s.Date, s.PaymentDays)
		if err != nil {
			return nil, err
		}
		p.PayDueDate = dueDate
	}

	return p, nil

}
