package models

type SaleDocument struct {
	ID                         int          `json:"id"`
	Number                     string       `json:"number"`
	Date                       string       `json:"date"`
	Time                       string       `json:"time"`
	PaymentDays                string       `json:"paymentDays"`
	Notes                      string       `json:"notes"`
	CurrencyCode               string       `json:"currencyCode"`
	ShipToName                 string       `json:"shipToName"`
	Type                       string       `json:"type"`
	InvoiceState               string       `json:"invoiceState"`
	PaymentType                string       `json:"paymentType"`
	ShipToAddress              string       `json:"shipToAddress"`
	NetTotal                   float64      `json:"netTotal"`
	VatTotal                   float64      `json:"vatTotal"`
	Rounding                   float64      `json:"rounding"`
	Total                      float64      `json:"total"`
	Paid                       string       `json:"paid"`
	TaxExemptCertificateNumber string       `json:"taxExemptCertificateNumber"`
	DeliveryDate               string       `json:"deliveryDate"`
	ShippingDate               string       `json:"shippingDate"`
	ReferenceNumber            string       `json:"referenceNumber"`
	CustomReferenceNumber      string       `json:"customReferenceNumber"`
	PaymentStatus              string       `json:"paymentStatus"`
	CompanyInfo                *CompanyInfo `json:"companyInfo"`
	Customer                   *Customer    `json:"companyInfo"`
	InvoiceRows                InvoiceRows  `json:"rows"`
	VatRate                    *VatRate     `json:"rows"`
}

type InvoiceRows []InvoiceRow

type InvoiceRow struct {
	RowID    string `json:"rowID"`
	ItemName string `json:"itemName"`
	Barcode  string `json:"barcode"`
	Amount   string `json:"amount"`
	Price    string `json:"price"`
	Discount string `json:"discount"`

	DeliveryDate      string  `json:"deliveryDate"`
	BillingStartDate  string  `json:"billingStartDate"`
	BillingEndDate    string  `json:"billingEndDate"`
	Code              string  `json:"code"`
	FinalNetPrice     float64 `json:"finalNetPrice"`
	FinalPriceWithVAT float64 `json:"finalPriceWithVAT"`
	RowNetTotal       float64 `json:"rowNetTotal"`
	RowVAT            float64 `json:"rowVAT"`
	RowTotal          float64 `json:"rowTotal"`
}

type CompanyInfo struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	Code               string `json:"code"`
	VAT                string `json:"VAT"`
	Phone              string `json:"phone"`
	Mobile             string `json:"mobile"`
	Fax                string `json:"fax"`
	Email              string `json:"email"`
	Web                string `json:"web"`
	BankAccountNumber  string `json:"bankAccountNumber"`
	BankName           string `json:"bankName"`
	BankSWIFT          string `json:"bankSWIFT"`
	BankIBAN           string `json:"bankIBAN"`
	BankAccountNumber2 string `json:"bankAccountNumber2"`
	BankName2          string `json:"bankName2"`
	BankSWIFT2         string `json:"bankSWIFT2"`
	BankIBAN2          string `json:"bankIBAN2"`
	Address            string `json:"address"`
	Country            string `json:"country"`
	DefaultCurrency    string `json:"defaultCurrency"`
}

type VatRate struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Rate         string `json:"rate"`
	Code         string `json:"code"`
	Active       string `json:"active"`
	Added        string `json:"added"`
	LastModified string `json:"lastModified"`
	IsReverseVat int    `json:"isReverseVat"`
	ReverseRate  int    `json:"reverseRate"`
}

type VatRates []VatRate

func VatRateById(Vt []VatRate, id string) VatRate {
	for _, vatrate := range Vt {
		if vatrate.ID == id {
			return vatrate
		}
	}
	return VatRate{}
}

type Customer struct {
	ID                 int    `json:"id"`
	FullName           string `json:"fullName"`
	CompanyName        string `json:"companyName"`
	FirstName          string `json:"firstName"`
	LastName           string `json:"lastName"`
	Phone              string `json:"phone"`
	Mobile             string `json:"mobile"`
	Email              string `json:"email"`
	Fax                string `json:"fax"`
	Code               string `json:"code"`
	ReferenceNumber    string `json:"referenceNumber"`
	CustomerCardNumber string `json:"customerCardNumber"`
	OutsideEU          int    `json:"outsideEU"`
	VatNumber          string `json:"vatNumber"`
	Skype              string `json:"skype"`
	Website            string `json:"website"`
	BankName           string `json:"bankName"`
	BankAccountNumber  string `json:"bankAccountNumber"`
	BankIBAN           string `json:"bankIBAN"`
	BankSWIFT          string `json:"bankSWIFT"`
	Notes              string `json:"notes"`
	CustomerType       string `json:"customerType"`
	Address            string `json:"address"`
	Street             string `json:"street"`
	Address2           string `json:"address2"`
	City               string `json:"city"`
	PostalCode         string `json:"postalCode"`
	Country            string `json:"country"`
	State              string `json:"state"`
	AddressTypeName    string `json:"addressTypeName"`
	TaxExempt          int    `json:"taxExempt"`
}
