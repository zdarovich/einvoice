package einvoice

import (
	"crypto/md5"
	"encoding/xml"
	"fmt"
)

//Hash ...
func Hash(in *EInvoice) (string, error) {
	xmlInvoice, err := xml.Marshal(in)
	if err != nil {
		return "", nil
	}
	hash := md5.Sum(xmlInvoice)
	md5string := fmt.Sprintf("%x", hash)
	return md5string, nil
}
