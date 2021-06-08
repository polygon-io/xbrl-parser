package xbrl

import "strings"

// Unit specifies the unit in which a numeric fact has been measured.
// A Unit can be either a simple measure, product of measures, or a ratio of products of measures with a numerator and a denominator.
// https://www.xbrl.org/Specification/XBRL-2.1/REC-2003-12-31/XBRL-2.1-REC-2003-12-31+corrected-errata-2013-02-20.html#_4.8
type Unit struct {
	ID       string   `xml:"id,attr"`
	Measures Measures `xml:"measure"`
	Divide   *Divide  `xml:"divide"`
}

// Divide is represents a complex Unit that has a numerator and a denominator.
// For example, you can represent a complex unit like earnings per share (EPS) as dollars per share (USD / share).
// https://www.xbrl.org/Specification/XBRL-2.1/REC-2003-12-31/XBRL-2.1-REC-2003-12-31+corrected-errata-2013-02-20.html#_4.8.2
type Divide struct {
	Numerator   Measures `xml:"unitNumerator>measure"`
	Denominator Measures `xml:"unitDenominator>measure"`
}

// Measure represents a unit of measure.
// https://www.xbrl.org/Specification/XBRL-2.1/REC-2003-12-31/XBRL-2.1-REC-2003-12-31+corrected-errata-2013-02-20.html#_4.8.2
type Measure struct {
	Value string `xml:",chardata"`
}

type Measures []Measure

// String returns a human readable representation of the Unit.
func (u Unit) String() string {
	// If the Divide element is not nil, there can be no top-level Meaures.
	if u.Divide != nil {
		return u.Divide.Numerator.String() + " / " + u.Divide.Denominator.String()
	}

	// If the divider element is nil, there must be 1+ top-level Measures.
	return u.Measures.String()
}

// String returns the local name of the measure if the value is formatted as 'xsd:Qname', otherwise the value itself is returned.
// Ex: `<measure>iso4127:USD</measure>` -> "USD"
//     `<measure>shares</measure>`      -> "shares"
func (m Measure) String() string {
	if index := strings.IndexRune(m.Value, ':'); index != -1 && index < len(m.Value) {
		return m.Value[index+1 : len(m.Value)]
	}

	return m.Value
}

func (m Measures) String() string {
	// More than one Measure implies multiplication.
	var builder strings.Builder
	for index, measure := range m {
		if index > 0 {
			builder.WriteString(" * ")
		}

		builder.WriteString(measure.String())
	}

	return builder.String()
}
