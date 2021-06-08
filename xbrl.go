package xbrl

import (
	"encoding/xml"
	"fmt"
)

// NotImplemented represents an expected element in the XBRL that isn't handled yet, but should not be considered a Fact.
type NotImplemented []*struct{}

// RawXBRL represents the XML structure of an XBRL document.
// This is not a feature complete XBRL parser!
// See the fields of type NotImplemented for an idea of what's missing.
// Also note that this struct doesn't support Tuple facts (https://www.xbrl.org/Specification/XBRL-2.1/REC-2003-12-31/XBRL-2.1-REC-2003-12-31+corrected-errata-2013-02-20.html#_4.9)
//
// You can use this struct directly, but XBRL is structured in a more convenient way.
// See the comment on XBRL for more info.
type RawXBRL struct {
	Contexts []Context `xml:"context"`
	Units    []Unit    `xml:"unit"`

	Facts []Fact `xml:",any"`

	// The fields below are not properly implemented, but need to be here so they aren't lumped into the `Facts` slice.

	SchemaRef    NotImplemented `xml:"schemaRef"`
	LinkbaseRef  NotImplemented `xml:"linkbaseRef"`
	RoleRef      NotImplemented `xml:"roleRef"`
	ArcRoleRef   NotImplemented `xml:"arcroleRef"`
	FootnoteLink NotImplemented `xml:"footnoteLink"`
}

// XBRL contains maps for contexts and units so they can be accessed easier when looping through facts.
// You can either unmarshal XML directly into this struct (it has a custom unmarshaller),
// or you can unmarshal XML into a RawXBRL struct and call NewProcessedXBRL(RawXBRL) to process the raw XBRL into this format.
type XBRL struct {
	ContextsByID map[string]Context
	UnitsByID    map[string]Unit

	Facts []Fact
}

// NewProcessedXBRL constructs a XBRL struct from a RawXBRL struct.
func NewProcessedXBRL(raw RawXBRL) XBRL {
	contextsByID := make(map[string]Context, len(raw.Contexts))
	unitsByID := make(map[string]Unit, len(raw.Units))

	for _, context := range raw.Contexts {
		contextsByID[context.ID] = context
	}

	for _, unit := range raw.Units {
		unitsByID[unit.ID] = unit
	}

	return XBRL{
		ContextsByID: contextsByID,
		UnitsByID:    unitsByID,
		Facts:        raw.Facts,
	}
}

// UnmarshalXML implements xml.Unmarshaler and unmarshals the contents as a RawXBRL,
// then processes it and populates this struct's fields.
func (x *XBRL) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var raw RawXBRL
	if err := d.DecodeElement(&raw, &start); err != nil {
		return nil
	}

	*x = NewProcessedXBRL(raw)
	return nil
}

// Validate checks that all Facts are valid and reference contexts and units that also exist.
// Note that since this parser does not properly handle Tuple elements, it's possible that some malformed Facts were unmarshalled.
func (x XBRL) Validate() error {
	for _, fact := range x.Facts {
		if !fact.IsValid() {
			return fmt.Errorf("invalid fact: %s:%s", fact.XMLName.Space, fact.XMLName.Local)
		}

		if _, exists := x.ContextsByID[fact.ContextRef]; !exists {
			return fmt.Errorf("fact (%s:%s) references non-existent context: %s", fact.XMLName.Space, fact.XMLName.Local, fact.ContextRef)
		}

		if fact.UnitRef != nil {
			if _, exists := x.UnitsByID[*fact.UnitRef]; !exists {
				return fmt.Errorf("fact (%s:%s) references non-existent unit: %s", fact.XMLName.Space, fact.XMLName.Local, *fact.UnitRef)
			}
		}
	}

	return nil
}

// IsValid validates the Facts in this struct and returns true if no error was found.
func (x XBRL) IsValid() bool {
	return x.Validate() == nil
}
