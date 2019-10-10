package einvoice

//Footer ...
type footer struct {
	Text                string `xml:",chardata"`
	TotalNumberInvoices string `xml:"TotalNumberInvoices"`
	TotalAmount         string `xml:"TotalAmount"`
}
