package earveldaja

const (
	xsi     = "http://www.w3.org/2001/XMLSchema-instance"
	xsd     = "http://www.w3.org/2001/XMLSchema"
	soapenv = "http://schemas.xmlsoap.org/soap/envelope/"
	sec     = "https://secure.maventa.com/"
	soapenc = "http://schemas.xmlsoap.org/soap/encoding/"
)

type Envelope struct {
	XMLName string `xml:"soap:Envelope"`
	Text    string `xml:",chardata"`
	Xsi     string `xml:"xmlns:xsi,attr"`
	Xsd     string `xml:"xmlns:xsd,attr"`
	Soap    string `xml:"xmlns:soap,attr"`
	Body    Body   `xml:"soap:Body"`
}

type ResponseEnvelope struct {
	XMLName string `xml:"Envelope"`
	Text    string `xml:",chardata"`
	Xsi     string `xml:"xmlns:xsi,attr"`
	Xsd     string `xml:"xmlns:xsd,attr"`
	Soap    string `xml:"xmlns:soap,attr"`
	Body    Body   `xml:"Body"`
}

type Body struct {
	CreateForwardingContractService         *CreateForwardingContractService         `xml:"CreateForwardingContractService,omitempty"`
	CreateForwardingContractServiceResponse *CreateForwardingContractServiceResponse `xml:"CreateForwardingContractServiceResponse,omitempty"`
	ReceiveInvoiceService                   *ReceiveInvoiceService                   `xml:"ReceiveInvoiceService,omitempty"`
	ReceiveInvoiceServiceResponse           *ReceiveInvoiceServiceResponse           `xml:"ReceiveInvoiceServiceResponse,omitempty"`
	GetInvoiceListService                   *GetInvoiceListService                   `xml:"GetInvoiceListService,omitempty"`
	GetInvoiceListServiceResponse           *GetInvoiceListServiceResponse           `xml:"GetInvoiceListServiceResponse,omitempty"`
	Fault                                   *Fault                                   `xml:"Fault,omitempty"`
}

type ReceiveInvoiceQuery struct {
	DocId            int
	TransferDatetime string `xml:"transfer_datetime"`
	SenderRegcode    string `xml:"sender_regcode"`
	SenderName       string `xml:"sender_name"`
	ReceiverName     string `xml:"receiver_name"`
	SenderInvoiceID  string `xml:"sender_invoice_id"`
	InvoiceXMLFile   string `xml:"invoice_xml_file"`
	InvoicePdfFile   string `xml:"invoice_pdf_file"`
	ProviderCode     string `xml:"provider_code"`
	InvoiceSignature string `xml:"invoice_signature"`
	ReceiverRegcode  string `xml:"receiver_regcode"`
}

type CreateForwardingContractQuery struct {
	TransferDatetime string `xml:"transfer_datetime"`
	ProviderCode     string `xml:"provider_code"`
	ForwarderRegcode string `xml:"forwarder_regcode"`
	Signature        string `xml:"signature"`
}

type CreateForwardingContractServiceResponse struct {
	Response CreateForwardingContractQueryResponse `xml:"response"`
}

type CreateForwardingContractQueryResponse struct {
	Status  string `xml:"status"`
	Message string `xml:"message"`
}

type CreateForwardingContractService struct {
	Text                          string                        `xml:",chardata"`
	Xmlns                         string                        `xml:"xmlns,attr"`
	CreateForwardingContractQuery CreateForwardingContractQuery `xml:"request"`
}

type ReceiveInvoiceService struct {
	Text                string              `xml:",chardata"`
	Xmlns               string              `xml:"xmlns,attr"`
	ReceiveInvoiceQuery ReceiveInvoiceQuery `xml:"request"`
}

type ReceiveInvoiceServiceResponse struct {
	Response ReceiveInvoiceQueryResponse `xml:"response"`
}

type ReceiveInvoiceQueryResponse struct {
	Status          string `xml:"status"`
	Message         string `xml:"message"`
	SenderInvoiceID string `xml:"sender_invoice_id"`
}

type GetInvoiceListService struct {
	Text                string              `xml:",chardata"`
	Xmlns               string              `xml:"xmlns,attr"`
	GetInvoiceListQuery GetInvoiceListQuery `xml:"request"`
}

type GetInvoiceListQuery struct {
	ProviderCode     string `xml:"provider_code"`
	FromDatetime     string `xml:"from_datetime"`
	ToDatetime       string `xml:"to_datetime"`
	TransferDatetime string `xml:"transfer_datetime"`
	Signature        string `xml:"signature"`
}

type GetInvoiceListServiceResponse struct {
	GetInvoiceListQueryResponse GetInvoiceListQueryResponse `xml:"response"`
}

type GetInvoiceListQueryResponse struct {
	Text           string `xml:",chardata"`
	Status         string `xml:"status"`
	Message        string `xml:"message"`
	InvoiceListing InvoiceListing
}

type InvoiceListing struct {
	Invoices Invoices `xml:"invoice"`
}
type Invoices []Invoice

type Invoice struct {
	InvoiceId        int    `xml:"invoice_id"`
	ReceivedDatetime string `xml:"received_datetime"`
}

type Fault struct {
	Code        string `xml:"faultcode"`
	Description string `xml:"faultstring"`
	Detail      string `xml:"detail"`
}
