package xbrl

type XBRL struct {
	Contexts []Context `xml:"context"`
	Units    []Unit    `xml:"unit"`

	Facts []Fact `xml:",any"`
}
