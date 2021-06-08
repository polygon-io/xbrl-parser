package xbrl

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalFact(t *testing.T) {
	t.Run("nil fact", func(t *testing.T) {
		// language=xml
		factXML := `<myns:sillyFact contextRef="c1" xsi:nil="true"/>`

		var fact Fact
		require.NoError(t, xml.Unmarshal([]byte(factXML), &fact))

		assert.Equal(t, xml.Name{Space: "myns", Local: "sillyFact"}, fact.XMLName)
		assert.Equal(t, FactTypeNil, fact.Type())
		assert.True(t, fact.IsValid())
		assert.Equal(t, "c1", fact.ContextRef)
	})

	t.Run("non-numeric fact", func(t *testing.T) {
		// language=xml
		factXML := `<ci:concentrationsNote contextRef="c1">Some cool text block about concentrations.</ci:concentrationsNote>`

		var fact Fact
		require.NoError(t, xml.Unmarshal([]byte(factXML), &fact))

		assert.Equal(t, xml.Name{Space: "ci", Local: "concentrationsNote"}, fact.XMLName)
		assert.Equal(t, FactTypeNonNumeric, fact.Type())
		assert.True(t, fact.IsValid())
		assert.Equal(t, "c1", fact.ContextRef)
		assert.Nil(t, fact.UnitRef)
		assert.Equal(t, "Some cool text block about concentrations.", fact.Value())
	})

	t.Run("simple numeric fact", func(t *testing.T) {
		// language=xml
		factXML := `<ci:capitalLeases id="id123" contextRef="c1" unitRef="u1" precision="3">727432</ci:capitalLeases>`

		var fact Fact
		require.NoError(t, xml.Unmarshal([]byte(factXML), &fact))

		assert.Equal(t, xml.Name{Space: "ci", Local: "capitalLeases"}, fact.XMLName)
		assert.Equal(t, FactTypeNonFraction, fact.Type())
		assert.True(t, fact.IsValid())
		assert.Equal(t, "id123", fact.ID)
		assert.Equal(t, "c1", fact.ContextRef)
		require.NotNil(t, fact.UnitRef)
		assert.Equal(t, "u1", *fact.UnitRef)
		require.NotNil(t, fact.Precision)
		assert.Equal(t, "3", *fact.Precision)
		assert.Nil(t, fact.Decimals)

		val, err := fact.NumericValue()
		require.NoError(t, err)
		assert.EqualValues(t, 727432, val)
	})

	t.Run("simple decimal numeric fact", func(t *testing.T) {
		// language=xml
		factXML := `<us-gaap:EarningsPerShareBasic contextRef="i0ad" decimals="2" id="id3Vyb" unitRef="usdPerShare">0.64</us-gaap:EarningsPerShareBasic>`

		var fact Fact
		require.NoError(t, xml.Unmarshal([]byte(factXML), &fact))

		assert.Equal(t, xml.Name{Space: "us-gaap", Local: "EarningsPerShareBasic"}, fact.XMLName)
		assert.Equal(t, FactTypeNonFraction, fact.Type())
		assert.True(t, fact.IsValid())
		assert.Equal(t, "id3Vyb", fact.ID)
		assert.Equal(t, "i0ad", fact.ContextRef)
		require.NotNil(t, fact.UnitRef)
		assert.Equal(t, "usdPerShare", *fact.UnitRef)
		assert.Nil(t, fact.Precision)
		require.NotNil(t, fact.Decimals)
		assert.Equal(t, "2", *fact.Decimals)

		val, err := fact.NumericValue()
		require.NoError(t, err)
		assert.EqualValues(t, 0.64, val)
	})

	t.Run("fraction type numeric fact", func(t *testing.T) {
		// language=xml
		factXML := `<myTaxonomy:oneThird id="oneThird" unitRef="u1" contextRef="numC1">
	<numerator>1</numerator>
	<denominator>3</denominator>
</myTaxonomy:oneThird>`

		var fact Fact
		require.NoError(t, xml.Unmarshal([]byte(factXML), &fact))

		assert.Equal(t, xml.Name{Space: "myTaxonomy", Local: "oneThird"}, fact.XMLName)
		assert.Equal(t, FactTypeFraction, fact.Type())
		assert.True(t, fact.IsValid())
		assert.Equal(t, "oneThird", fact.ID)
		assert.Equal(t, "numC1", fact.ContextRef)
		require.NotNil(t, fact.UnitRef)
		assert.Equal(t, "u1", *fact.UnitRef)
		assert.Nil(t, fact.Precision)
		assert.Nil(t, fact.Decimals)

		val, err := fact.NumericValue()
		require.NoError(t, err)
		assert.EqualValues(t, 1.0/3.0, val)
	})
}
