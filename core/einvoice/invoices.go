package einvoice

//Invoices ...
type Invoices []Invoice

//ElementExists ...
func (invoices Invoices) ElementExists(invoiceNumber string) bool {
	for _, invoice := range invoices {
		if invoice.InvoiceInformation.InvoiceNumber == invoiceNumber {
			return true
		}
	}
	return false
}

//ElementsExist ...
func (invoices Invoices) ElementsExist(invoiceNumbers []string) bool {
	elementMap := make(map[string]bool)

	for _, invoiceNum := range invoiceNumbers {
		elementMap[invoiceNum] = true
	}
	for _, invoice := range invoices {
		if elementMap[invoice.InvoiceInformation.InvoiceNumber] == false {
			return false
		}
	}
	return true
}

//NumberList returns every invoice's number in the given slice
func (invoices Invoices) NumberList() []string {
	var list []string
	for _, invoice := range invoices {
		list = append(list, invoice.InvoiceInformation.InvoiceNumber)
	}
	return list
}
