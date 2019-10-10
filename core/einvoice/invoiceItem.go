package einvoice

import (
	"github.com/zdarovich/einvoice/core/models"
	"strconv"
	"strings"
)

/*
type GroupEntry struct {
	GroupDescription string          `xml:"GroupDescription,omitempty"`
	GroupAmount      string          `xml:"GroupAmount,omitempty"` //ItemAmount elementide summa Decimal4FractionDigitsType
	GroupSum         string          `xml:"GroupSum,omitempty"`    //Decimal4FractionDigitsType
	Addition         *AdditionRecord `xml:"Addition,omitempty"`
	VAT              *VAT            `xml:"VAT,omitempty"`
	GroupTotal       string          `xml:"GroupTotal,omitempty"` //ItemTotal elementide summa Decimal4FractionDigitsType
}
*/

//Group - ItemGroup, TotalGroup ...
type group struct {
	Text           string      `xml:",chardata"`
	GroupID        string      `xml:"groupId,attr,omitempty"`
	ItemEntryArray []itemEntry `xml:"ItemEntry"`
	//GroupEntry     *GroupEntry `xml:"GroupEntry,omitempty"` not implementing
}

//Item ...
type item struct {
	Text              string `xml:",chardata"`
	InvoiceTotalGroup *group `xml:"InvoiceTotalGroup,omitempty"` //not implementing
	InvoiceItemGroup  *group `xml:"InvoiceItemGroup"`
	//InvoiceItemTotalGroup - not implementing
}

//ItemEntry ...
type itemEntry struct {
	Text            string          `xml:",chardata"`
	RowNo           string          `xml:"RowNo,omitempty"`
	SerialNumber    string          `xml:"SerialNumber,omitempty"`
	SellerProductID string          `xml:"SellerProductId,omitempty"`
	BuyerProductID  string          `xml:"BuyerProductId,omitempty"` //we are not implementing that
	TaricCode       string          `xml:"TaricCode,omitempty"`      //we are not implementing that
	CustomerRef     string          `xml:"CustomerRef,omitempty"`
	Description     string          `xml:"Description"`
	EAN             string          `xml:"EAN,omitempty"`            //barcode
	InitialReading  string          `xml:"InitialReading,omitempty"` //Perioodilise arvelduse puhul perioodi algnäit
	FinalReading    string          `xml:"FinalReading,omitempty"`   //Perioodilise arvelduse puhul perioodi lõppnäit.
	ItemPrice       string          `xml:"ItemPrice,omitempty"`      //field only for parsing invoice from Omniva
	ItemAmount      string          `xml:"ItemAmount,omitempty"`     //field only for parsing invoice from Omniva
	DetailInfo      *itemDetailInfo `xml:"ItemDetailInfo,omitempty"`
	ItemSum         string          `xml:"ItemSum,omitempty"`
	Addition        *additionRecord `xml:"Addition,omitempty"`
	VAT             *vat            `xml:"VAT,omitempty"`
	ItemTotal       string          `xml:"ItemTotal,omitempty"`
}

//ItemDetailInfo ...
type itemDetailInfo struct {
	Text   string `xml:",chardata"`
	Amount string `xml:"ItemAmount,omitempty"`
	Price  string `xml:"ItemPrice,omitempty"`
}

const (
	emptyDate = "0000-00-00"
)

//NewInvoiceItem ...
func newInvoiceItem(invoiceRows models.InvoiceRows, currency string, vatRate *models.VatRate) (*item, error) {

	var rows []itemEntry

	for i, row := range invoiceRows {

		twoDigitPrice := twoDigitDecimalFromString(row.Price)
		priceFloat, err := strconv.ParseFloat(twoDigitPrice, 64)
		if err != nil {
			return nil, err
		}
		//amount can be 2.34 KG
		amountFloat, err := strconv.ParseFloat(row.Amount, 64)
		if err != nil {
			return nil, err
		}
		//price it should be multiplied by amount
		itemSum := twoDigit(amountFloat * priceFloat)
		addition, err := newRowAdditionRecord(&row)
		if err != nil {
			return nil, err
		}

		item := itemEntry{
			RowNo:           strconv.Itoa(i + 1),
			SerialNumber:    "",       //Leave empty. E-invoice standard says this field is for the sold item's  (unique) serial number. We do not have such information. Before it was row.Code,
			SellerProductID: row.Code, //Use product code instead. At least one customer has already clearly indicated that they want product code. Product code would allow interoperability between systems. Internal product ID is not very useful information.
			Description:     row.ItemName,
			EAN:             row.Barcode, //International Article Number (barcode)
			DetailInfo: &itemDetailInfo{
				Amount: row.Amount,
				Price:  twoDigitPrice,
			},
			ItemSum:   itemSum, //before discount & VAT.
			Addition:  addition,
			ItemTotal: fourDigit(row.RowTotal),
		}

		if vatRate.Rate != "0" {
			vat, err := newRowVat(&row, currency, vatRate.Rate)
			if err != nil {
				return nil, err
			}
			item.VAT = vat
		}
		/*
			Does not quite match e-invoice spec (these fields should rather show meter readings, not dates), so perhaps should be omitted.


			if row.BillingStartDate != emptyDate {
				item.InitialReading = row.BillingStartDate
			}

			if row.BillingEndDate != emptyDate {
				item.FinalReading = row.BillingEndDate
			}
		*/

		rows = append(rows, item)
	}

	return &item{
		InvoiceItemGroup: &group{
			ItemEntryArray: rows,
		},
	}, nil
}

func twoDigitDecimalFromString(s string) string {
	sep := "."
	if strings.Contains(s, sep) {
		xs := strings.Split(s, sep)
		remainder := xs[1]
		if len(remainder) > 1 {
			remainder = remainder[:2]
		} else {
			remainder = remainder[:1]
		}
		xs[1] = remainder
		twoDigit := strings.Join(xs, sep)
		return twoDigit
	}

	return s
}
