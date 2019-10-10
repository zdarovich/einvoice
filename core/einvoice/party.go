package einvoice

import (
	"github.com/zdarovich/einvoice/core/models"
)

//Party ...
type party struct {
	Text         string            `xml:",chardata"`
	Name         string            `xml:"Name"`
	RegNumber    string            `xml:"RegNumber,omitempty"`
	VATRegNumber string            `xml:"VATRegNumber,omitempty"`
	ContactData  *contactData      `xml:"ContactData,omitempty"`
	AccountInfo  []bankAccountInfo `xml:"AccountInfo,omitempty"`
}

//ContactData ...
type contactData struct {
	Text              string   `xml:",chardata"`
	ContactName       string   `xml:"ContactName,omitempty"`
	ContactPersonCode string   `xml:"ContactPersonCode,omitempty"`
	PhoneNumber       string   `xml:"PhoneNumber,omitempty"`
	FaxNumber         string   `xml:"FaxNumber,omitempty"`
	EMailAddress      string   `xml:"E-mailAddress,omitempty"`
	LegalAddress      *address `xml:"LegalAddress,omitempty"`
	MailAddress       *address `xml:"MailAddress,omitempty"`
}

//BankAccountInfo ...
type bankAccountInfo struct {
	Text          string `xml:",chardata"`
	AccountNumber string `xml:"AccountNumber"`
	IBAN          string `xml:"IBAN,omitempty"`
	BIC           string `xml:"BIC,omitempty"`
	BankName      string `xml:"BankName,omitempty"`
}

//Address ...
type address struct {
	Text           string `xml:",chardata"`
	PostalAddress1 string `xml:"PostalAddress1"`           //Tänav, maja, korter.
	PostalAddress2 string `xml:"PostalAddress2,omitempty"` //Küla, alev.
	City           string `xml:"City"`
	PostalCode     string `xml:"PostalCode,omitempty"` //Linn või maakond.
	Country        string `xml:"Country,omitempty"`
}

//NewBuyer ...
func newBuyer(in *models.Customer) (*party, error) {

	b := &party{
		Name:         in.FullName, //may use company name, but it is for companies only
		RegNumber:    in.Code,     //alert - can be personal ID code - is that okay? can check via CustomerType "COMPANY" or "PERSON".
		VATRegNumber: in.VatNumber,
		ContactData: &contactData{
			ContactName: in.FirstName,
			PhoneNumber: in.Phone,
			FaxNumber:   in.Fax,
			LegalAddress: &address{
				PostalAddress1: in.Street,
				PostalAddress2: in.Address2,
				City:           in.City,
				PostalCode:     in.PostalCode,
				Country:        in.Country,
			},
		},
		AccountInfo: []bankAccountInfo{
			bankAccountInfo{
				AccountNumber: in.BankAccountNumber,
				IBAN:          in.BankIBAN,
				BIC:           in.BankSWIFT,
				BankName:      in.BankName,
			},
		},
	}
	return b, nil
}

//NewSeller ...
func newSeller(in *models.CompanyInfo) (*party, error) {

	s := &party{
		Name:         in.Name, //may use company name, but it is for companies only
		RegNumber:    in.Code, //alert - can be personal ID code - is that okay? can check via CustomerType "COMPANY" or "PERSON".
		VATRegNumber: in.VAT,
		ContactData: &contactData{
			PhoneNumber:  in.Phone,
			FaxNumber:    in.Fax,
			EMailAddress: in.Email,
			LegalAddress: &address{
				PostalAddress1: in.Address,
				Country:        in.Country,
			},
		},
		AccountInfo: []bankAccountInfo{
			bankAccountInfo{
				AccountNumber: in.BankAccountNumber,
				IBAN:          in.BankIBAN,
				BIC:           in.BankSWIFT,
				BankName:      in.BankName,
			},
			bankAccountInfo{
				AccountNumber: in.BankAccountNumber2,
				IBAN:          in.BankIBAN2,
				BIC:           in.BankSWIFT2,
				BankName:      in.BankName2,
			},
		},
	}
	return s, nil
}
