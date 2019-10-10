package einvoice

import (
	"crypto/rand"
	"fmt"
	"io"
	"strconv"
	"time"
)

const (
	TimeFormat      = "2006-01-02"
	twoDigitFormat  = "%.2f"
	fourDigitFormat = "%.4f"
)

//TwoDigit Decimal2FractionDigitsType
func twoDigit(float float64) string {
	return fmt.Sprintf(twoDigitFormat, float)
}

//FourDigit Decimal4FractionDigitsType
func fourDigit(float float64) string {
	return fmt.Sprintf(fourDigitFormat, float)
}

// NewUUID does not generate a random UUID according to RFC 4122. It is of length 10, not 16. and without "-" in between the slices
func newUUID() (string, error) {
	uuid := make([]byte, 10)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x", uuid[:]), nil
}

//NewDueDate ...
func newDueDate(date string, paymentDays string) (string, error) {

	//There is always a due date (even when empty all the configuration for default payment days it puts 14 or 0 days for it),
	if paymentDays == "" {
		//when still no duedate then use invoice date as a due date
		return date, nil
	}
	numberOfDays, err := strconv.ParseUint(paymentDays, 10, 32)
	if err != nil {
		return "", err
	}
	dueDate, err := getDueDate(date, numberOfDays)
	if err != nil {
		return "", err
	}

	return dueDate, nil
}

//add number of days to invoice Date (YYYY-MM-DD format)
func getDueDate(inDate string, numberOfDays uint64) (string, error) {
	myTime, err := time.Parse(TimeFormat, inDate)

	if err != nil {
		return "", nil
	}

	day := 24 * time.Hour

	dueDate := myTime.Add(time.Duration(numberOfDays) * day).Local().Format(TimeFormat)

	return dueDate, nil
}
