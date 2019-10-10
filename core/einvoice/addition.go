package einvoice

import (
	"github.com/zdarovich/einvoice/core/models"
	"strconv"
)

//AdditionRecord - Element soodustuste ja juurdehindluse kajastamiseks (to reflect - отражать).
//Used in InvoiceSumGroup, InvoiceItem, ItemEntry, GroupEntry

const (
	discount     = "DSC"
	extraCharge  = "CHR"
	floatBitSize = 32
)

//AdditionRecord ...
type additionRecord struct {
	Text       string `xml:",chardata"`
	AddCode    string `xml:"addCode,attr"`      //DSC – soodustus. CHR – juurdehindlus.
	AddContent string `xml:"AddContent"`        //Juurdehindluse/soodustuse nimetus (näiteks: „Püsikliendi soodustus“).
	AddRate    string `xml:"AddRate,omitempty"` //Decimal2Type protsent.
	AddSum     string `xml:"AddSum,omitempty"`  //Decimal4Type summa.
}

//currently Discount by default
func newRowAdditionRecord(row *models.InvoiceRow) (*additionRecord, error) {
	price, err := strconv.ParseFloat(row.Price, floatBitSize)
	if err != nil {
		return nil, err
	}
	disc, err := strconv.ParseFloat(row.Discount, floatBitSize)
	if err != nil {
		return nil, err
	}
	amount, err := strconv.ParseFloat(row.Amount, floatBitSize)
	if err != nil {
		return nil, err
	}
	//if there was no discount do not create a dummy object
	if disc == 0 {
		return nil, nil
	}

	a := &additionRecord{
		AddCode:    discount,
		AddContent: "selgitus",
		AddRate:    twoDigit(disc),
		AddSum:     fourDigit(amount * price * disc / 100), //PDF generator does not have this, so we can skip it
	}

	return a, nil
}
