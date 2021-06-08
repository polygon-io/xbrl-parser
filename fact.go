package xbrl

import (
	"encoding/xml"
	"errors"
	"strconv"
)

type FactType string

const (
	// FactTypeNil is a fact which has an `xsi:nil` attribute set to a truthy value.
	// A nil fact is only guaranteed to have an XMLName and ContextRef.
	FactTypeNil FactType = "nil"

	// FactTypeNonNumeric is a non-nil fact that does not describe a numeric value (ie text, dates, encoded binary data, etc).
	// A non-numeric fact is guaranteed to have an XMLName, ContextRef, and ValueStr.
	FactTypeNonNumeric FactType = "non_numeric"

	// FactTypeNonFraction is a non-nil fact describing a numeric value that can precisely expressed as a simple value.
	// A non-fraction fact is guaranteed to have an XMLName, ContextRef, UnitRef, ValueStr, and exactly one of Precision or Decimals.
	//
	// For example: <ci:capitalLeases contextRef="c1" unitRef="u1" precision="3">727432</ci:capitalLeases>
	//
	// Use Fact.NumericValue() for easy access to the numeric value as a float64.
	FactTypeNonFraction FactType = "non_fraction"

	// FactTypeFraction is a non-nil fact describing a numeric value that is the result of a numerator / denominator.
	// Usually the numeric value that these facts describe cannot be precisely expressed by a float64 (ie 1/3 = 0.3333...)
	// A fraction fact is guaranteed to have an XMLName, ContextRef, UnitRef, Numerator, and Denominator.
	//
	// For example:
	// <myTaxonomy:oneThird id="oneThird" unitRef="u1" contextRef="numC1">
	//     <numerator>1</numerator>
	//     <denominator>3</denominator>
	// </myTaxonomy:oneThird>
	//
	// Use Fact.NumericValue() for easy access to the numeric value as a float64,
	// but be aware that the float64 representation may not be able to precisely represent the Facts actual value.
	FactTypeFraction FactType = "fraction"
)

// ErrNonNumericFactType is returned when a fact is expected to be numeric, but is not.
var ErrNonNumericFactType = errors.New("fact is not of type FactTypeFraction or FactTypeNonFraction")

// Fact represents an item in an XBRL document.
// A Fact is a simple value which is tied to a context that gives the fact more meaning.
//
// This struct contains fields that may or may not be nil depending on what type of Fact you're dealing with.
// See Fact.Type() to determine what type of fact you're dealing with.
// Then the various FactTypes to understand what fields in this struct are expected to exist for each FactType.
//
// For general information and details on XBRL Facts, see here:
// https://www.xbrl.org/Specification/XBRL-2.1/REC-2003-12-31/XBRL-2.1-REC-2003-12-31+corrected-errata-2013-02-20.html#_4.6
type Fact struct {
	XMLName xml.Name

	// ID uniquely identifies a fact within an XBRL document.
	// The spec does not require an ID attribute, but many items have an ID attribute, which is why it's included in this model.
	ID string `xml:"id,attr"`

	// Nil is an attribute denoting whether or not this fact is expressed as nil.
	Nil *bool `xml:"nil,attr"`

	// ContextRef is the ID of the context in the XBRL document that gives more meaning to this fact.
	ContextRef string `xml:"contextRef,attr"`

	// UnitRef is the ID of the unit in the XBRL document that this fact is expressed in.
	// It is non-nil for numeric facts only.
	UnitRef *string `xml:"unitRef,attr"`

	// Precision conveys the arithmetic precision of a measurement.
	// It can be either a non-negative integer or the special value "INF", which represents infinite precision.
	// If this is a numeric fact but NOT a fraction type, Precision will be non-nil if Decimals is nil,
	//
	// Examples and more info here:
	// https://www.xbrl.org/Specification/XBRL-2.1/REC-2003-12-31/XBRL-2.1-REC-2003-12-31+corrected-errata-2013-02-20.html#_4.6.4
	Precision *string `xml:"precision,attr"`

	// Decimals specifies the number of decimal places to which the value of the fact represented may be considered accurate.
	// It can be either an integer (positive or negative) or the special value "INF", which represents accuracy to infinite decimal places.
	// If this is a numeric fact but NOT a fraction type, Decimals will be non-nil if Precision is nil,
	//
	// Examples and more info here:
	// https://www.xbrl.org/Specification/XBRL-2.1/REC-2003-12-31/XBRL-2.1-REC-2003-12-31+corrected-errata-2013-02-20.html#_4.6.5
	Decimals *string `xml:"decimals,attr"`

	// ValueStr will be non-nil unless this is a numeric fraction type Fact.
	// Use NumericValue() to easily get the numeric value that this fact represents, regardless of whether or not it's a fraction type.
	ValueStr *string `xml:",chardata"`

	// Numerator and Denominator will be non-nil values if this is a fraction type Fact.
	// Use NumericValue() to easily get the numeric value that this fact represents, regardless of whether or not it's a fraction type.
	Numerator   *float64 `xml:"numerator"`
	Denominator *float64 `xml:"denominator"`
}

// Type returns the type of this Fact. See the comments on the various FactTypes for more information.
// Note that this function returning a particular type does not necessarily mean that the fact is semantically correct.
// See IsValid() to be certain that the fact is valid.
func (f Fact) Type() FactType {
	// If the nil attribute exists and is true, this is simply a nil fact
	if f.Nil != nil && *f.Nil {
		return FactTypeNil
	}

	// If the unitRef attribute exists, this is some kind of numeric attribute
	if f.UnitRef != nil {
		// If we have a numerator and denominator, it's a fraction fact
		if f.Numerator != nil && f.Denominator != nil {
			return FactTypeFraction
		}

		// Otherwise it's a simple non fraction numeric type.
		return FactTypeNonFraction
	}

	// All that's left is a plain non-numeric fact type
	return FactTypeNonNumeric
}

// IsValid confirms that f has at least the required fields that the FactType requires.
// Note that this function is not strict about extra fields existing.
func (f Fact) IsValid() bool {
	// All facts must have a context ref
	if f.ContextRef == "" {
		return false
	}

	// Some types have particular rules beyond what Type() checks for that must be true to be considered valid.
	switch f.Type() {
	case FactTypeFraction:
		// Fraction must have a non-zero Denominator
		return *f.Denominator != 0
	case FactTypeNonFraction:
		// NonFractions must have either a non-nil Precision or non-nil Decimals field
		return (f.Precision == nil) != (f.Decimals == nil)
	case FactTypeNonNumeric:
		return f.ValueStr != nil
	default:
		return true
	}
}

// NumericValue attempts to return the numeric value this fact represents.
// This function returns
// If this fact is a fraction type, this function returns the value of numerator / denominator.
// Note that fraction type facts generally cannot be precisely represented as a float64 and may have some rounding error.
func (f Fact) NumericValue() (float64, error) {
	switch f.Type() {
	case FactTypeFraction:
		return *f.Numerator / *f.Denominator, nil
	case FactTypeNonFraction:
		return strconv.ParseFloat(*f.ValueStr, 64)
	default:
		return 0, ErrNonNumericFactType
	}
}

// Value returns the ValueStr of this Fact, or empty string if f.ValueStr is nil.
func (f Fact) Value() string {
	if f.ValueStr != nil {
		return *f.ValueStr
	}

	return ""
}
