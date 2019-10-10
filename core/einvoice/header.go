package einvoice

//Header ...
type header struct {
	Text    string `xml:",chardata"`
	Date    string `xml:"Date"`
	FileID  string `xml:"FileId"`
	Version string `xml:"Version"`
}
