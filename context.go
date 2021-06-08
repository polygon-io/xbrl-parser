package xbrl

import "encoding/xml"

// Context contains information about the Entity being described, the reporting Period, and the reporting Scenario (scenario is NOT implemented).
// All of which are necessary for understanding a business Fact captured as an XBRL item.
// https://www.xbrl.org/Specification/XBRL-2.1/REC-2003-12-31/XBRL-2.1-REC-2003-12-31+corrected-errata-2013-02-20.html#_4.7
type Context struct {
	ID string `xml:"id,attr"`

	Period Period `xml:"period"`
	Entity Entity `xml:"entity"`
}

// Entity documents the business entity for a Context (business, government department, individual, etc.).
// https://www.xbrl.org/Specification/XBRL-2.1/REC-2003-12-31/XBRL-2.1-REC-2003-12-31+corrected-errata-2013-02-20.html#_4.7.3
type Entity struct {
	Identifier Identifier `xml:"identifier"`
	Segments   Segments   `xml:"segment"`
}

// Identifier specifies a scheme for identifying business entities and an identifier that follows the scheme.
// https://www.xbrl.org/Specification/XBRL-2.1/REC-2003-12-31/XBRL-2.1-REC-2003-12-31+corrected-errata-2013-02-20.html#_4.7.3.1
// For Example:
// <identifier scheme="http://www.sec.gov/CIK">0000320193</identifier>
//
// The above `identifier` element specifies that the scheme for identifying the entity is through SEC CIK numbers,
// and that the identifier itself is 0000320193 (the CIK for Apple Inc.).
type Identifier struct {
	Scheme string `xml:"scheme,attr"`
	Value  string `xml:",chardata"`
}

// Segments is a type alias for a slice of Segment structs.
// It implements xml.Unmarshaller and puts and unmarshals any sub-elements as a Segment and puts it into the slice.
type Segments []Segment

// Segment is an optional container for additional information used to identify a business segment more completely
// for cases where the Identifier is insufficient.
// There are no Segments defined in the base XBRL spec, they must be defined in other XML schemas.
// https://www.xbrl.org/Specification/XBRL-2.1/REC-2003-12-31/XBRL-2.1-REC-2003-12-31+corrected-errata-2013-02-20.html#_4.7.3.2
type Segment struct {
	XMLName    xml.Name
	Attributes []xml.Attr `xml:",any,attr"`
	Value      string     `xml:",chardata"`
}

// UnmarshalXML implements xml.Unmarshaller for Segments.
// It unmarshals any sub-elements as Segments and puts them into this slice.
func (s *Segments) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var segmentsAnon struct {
		Segments []Segment `xml:",any"`
	}

	if err := d.DecodeElement(&segmentsAnon, &start); err != nil {
		return err
	}

	*s = segmentsAnon.Segments
	return nil
}

type PeriodType string

// All the supported PeriodType values. See Period.Type() for more information.
const (
	PeriodTypeDuration PeriodType = "duration"
	PeriodTypeInstant  PeriodType = "instant"
	PeriodTypeForever  PeriodType = "forever"
	PeriodTypeInvalid  PeriodType = "invalid"
)

// Period contains an instant or interval of time for a Context.
// https://www.xbrl.org/Specification/XBRL-2.1/REC-2003-12-31/XBRL-2.1-REC-2003-12-31+corrected-errata-2013-02-20.html#_4.7.2
type Period struct {
	// StartDate is non-nil and guaranteed to be before EndDate if Period.Type() returns Duration.
	StartDate *string `xml:"startDate"`
	// EndDate is non-nil and guaranteed to be after StartDate if Period.Type() returns Duration.
	EndDate *string `xml:"endDate"`

	// Instant is non-nil if Period.Type() returns Instant
	Instant *string `xml:"instant"`

	// Forever is non-nil if Period.Type() returns Forever.
	// Note ideally this would be a bool, but the XML Unmarshaller doesn't support
	// setting a boolean flag based on the existence of an empty tag (in this case `<forever/>`).
	Forever *struct{} `xml:"forever"`
}

// Type returns the type of this period to help clarify what fields in the Period struct are non-nil and valid to use.
// The comments on the attributes inside the Period struct explain when they can be used depending on what this function returns.
func (p Period) Type() PeriodType {
	if p.Forever != nil {
		return PeriodTypeForever
	}

	if p.Instant != nil {
		return PeriodTypeInstant
	}

	if p.StartDate != nil && p.EndDate != nil {
		return PeriodTypeDuration
	}

	return PeriodTypeInvalid
}
