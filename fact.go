package xbrl

import (
	"encoding/xml"
	"errors"
	"strconv"
)

// ErrInvalidNumericFact is returned when a numeric fact does not conform to the XBRL spec and a numeric value cannot be extracted from it.
// For example, if it's a fraction type fact and the denominator is 0, or the fact is malformed and missing attributes.
var ErrInvalidNumericFact = errors.New("numeric fact contains invalid data")

// NilFact is a fact in an XBRL document that has an xsi:nil set to a truthy value.
// These facts still must have a contextRef, but nothing else is guaranteed.
//
// Detailed info on Facts in general can be found here:
// https://www.xbrl.org/Specification/XBRL-2.1/REC-2003-12-31/XBRL-2.1-REC-2003-12-31+corrected-errata-2013-02-20.html#_4.6
type NilFact struct {
	XMLName xml.Name

	// ID uniquely identifies a fact within an XBRL document.
	// The spec does not require an ID attribute, but many items have an ID attribute, which is why it's included in this model.
	ID         string `xml:"id,attr"`
	ContextRef string `xml:"ContextRef,attr"`
}

// NonNumericFact is a non-nil fact that does not describe a numeric value (ie text, dates, encoded binary data, etc).
// For example:
// <ci:concentrationsNote contextRef="c1">Some cool text block about concentrations.</ci:concentrationsNote>
//
// Detailed info on Facts in general can be found here:
// https://www.xbrl.org/Specification/XBRL-2.1/REC-2003-12-31/XBRL-2.1-REC-2003-12-31+corrected-errata-2013-02-20.html#_4.6
type NonNumericFact struct {
	XMLName xml.Name

	// ID uniquely identifies a fact within an XBRL document.
	// The spec does not require an ID attribute, but many items have an ID attribute, which is why it's included in this model.
	ID         string `xml:"id,attr"`
	ContextRef string `xml:"contextRef,attr"`
	Value      string `xml:",chardata"`
}

// NumericFact is a non-nil fact that describes a numeric value.
//
// For example, a simple numeric fact could look like:
// <ci:capitalLeases contextRef="c1" unitRef="u1" precision="3">727432</ci:capitalLeases>
//
// And a fraction type numeric fact could look like:
// <myTaxonomy:oneThird id="oneThird" unitRef="u1" contextRef="numC1">
//     <numerator>1</numerator>
//     <denominator>3</denominator>
// </myTaxonomy:oneThird>
//
// Detailed info on Facts in general can be found here:
// https://www.xbrl.org/Specification/XBRL-2.1/REC-2003-12-31/XBRL-2.1-REC-2003-12-31+corrected-errata-2013-02-20.html#_4.6
type NumericFact struct {
	XMLName xml.Name

	// ID uniquely identifies a fact within an XBRL document.
	// The spec does not require an ID attribute, but many items have an ID attribute, which is why it's included in this model.
	ID         string `xml:"id,attr"`
	ContextRef string `xml:"contextRef,attr"`
	UnitRef    string `xml:"unitRef,attr"`

	// Precision conveys the arithmetic precision of a measurement.
	// It can be either a non-negative integer or the special value "INF", which represents infinite precision.
	// Precision must be non-nil if Decimals is nil, unless this fact is a fraction type (in which case both Precision and Decimals may be nil)
	//
	// Examples and more info here:
	// https://www.xbrl.org/Specification/XBRL-2.1/REC-2003-12-31/XBRL-2.1-REC-2003-12-31+corrected-errata-2013-02-20.html#_4.6.4
	Precision *string `xml:"precision,attr"`

	// Decimals specifies the number of decimal places to which the value of the fact represented may be considered accurate.
	// It can be either an integer (positive or negative) or the special value "INF", which represents accuracy to infinite decimal places.
	// Decimals must be non-nil if Precision is nil, unless this fact is a fraction type (in which case both Precision and Decimals may be nil)
	//
	// Examples and more info here:
	// https://www.xbrl.org/Specification/XBRL-2.1/REC-2003-12-31/XBRL-2.1-REC-2003-12-31+corrected-errata-2013-02-20.html#_4.6.5
	Decimals *string `xml:"decimals,attr"`

	// ValueStr will be non-nil if this is not a fraction type Fact.
	// Use NumericValue() to easily get the numeric value that this fact represents, regardless of whether or not it's a fraction type.
	ValueStr *string `xml:",chardata"`

	// Numerator and Denominator will be non-nil values if this is a fraction type Fact.
	// Use NumericValue() to easily get the numeric value that this fact represents, regardless of whether or not it's a fraction type.
	Numerator   *float64 `xml:"numerator"`
	Denominator *float64 `xml:"denominator"`
}

// IsFractionType returns true if this fact is described in the XBRL as a fraction with a numerator and non-zero denominator.
// NumericValue() will return the result of numerator / denominator if some potential rounding error is acceptable for your use-case.
func (f NumericFact) IsFractionType() bool {
	if f.Numerator != nil && f.Denominator != nil && *f.Denominator != 0 {
		return true
	}

	return false
}

// NumericValue attempts to return the numeric value this fact represents.
// If this fact is a fraction type, this function returns the value of numerator / denominator.
// Note that fraction type facts generally cannot be precisely represented as a float64 and may have some rounding error.
func (f NumericFact) NumericValue() (float64, error) {
	if f.IsFractionType() {
		return *f.Numerator / *f.Denominator, nil
	}

	if f.ValueStr != nil {
		return strconv.ParseFloat(*f.ValueStr, 64)
	}

	return 0, ErrInvalidNumericFact
}
