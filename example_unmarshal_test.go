package xbrl_test

import (
	"encoding/xml"
	"fmt"

	"github.com/polygon-io/xbrl-parser"
)

const doc = `<xbrl>
    <link:schemaRef xlink:type="simple" xlink:href="http://www.xbrl.org/us/fr/ci/2000-07-31/usfr-ci-2003.xsd"/>

    <context id="c1">
        <entity>
            <identifier scheme="http://www.sec.gov/CIK">0000320193</identifier>
        </entity>
        <period>
            <instant>2021-04-16</instant>
        </period>
    </context>

    <ci:assets precision="3" unitRef="u1" contextRef="c1">727</ci:assets>

    <unit id="u1">
        <measure>shares</measure>
    </unit>
</xbrl>`

func Example() {
	var processed xbrl.XBRL

	if err := xml.Unmarshal([]byte(doc), &processed); err != nil {
		panic(err)
	}

	fact := processed.Facts[0]
	if !fact.IsValid() {
		panic("fact invalid!")
	}

	factType := fact.Type()
	numericValue, err := fact.NumericValue()

	factContext := processed.ContextsByID[fact.ContextRef]
	factUnit := processed.UnitsByID[*fact.UnitRef]

	if err != nil {
		panic(err)
	}

	fmt.Printf("Fact: %s:%s (type: %s)\n", fact.XMLName.Space, fact.XMLName.Local, factType)
	fmt.Printf("      %.0f %s on %s\n", numericValue, factUnit.String(), *factContext.Period.Instant)

	// Output: Fact: ci:assets (type: non_fraction)
	//       727 shares on 2021-04-16
}
