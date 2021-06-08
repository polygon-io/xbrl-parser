package xbrl

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalUnit(t *testing.T) {
	t.Run("simple unit", func(t *testing.T) {
		// language=xml
		unitXml := `<unit>
			<measure>shares</measure>
		</unit>`

		var unit Unit
		require.NoError(t, xml.Unmarshal([]byte(unitXml), &unit))

		require.Len(t, unit.Measures, 1)
		assert.Equal(t, "shares", unit.Measures[0].Value)
		assert.Equal(t, "shares", unit.Measures[0].String())
		assert.Nil(t, unit.Divide)

		assert.Equal(t, unit.String(), "shares")
	})

	t.Run("product of measures", func(t *testing.T) {
		// language=xml
		unitXml := `<unit>
			<measure>myns:feet</measure>
			<measure>myns:feet</measure>
		</unit>`

		var unit Unit
		require.NoError(t, xml.Unmarshal([]byte(unitXml), &unit))

		require.Len(t, unit.Measures, 2)
		assert.Equal(t, "myns:feet", unit.Measures[0].Value)
		assert.Equal(t, "feet", unit.Measures[0].String())

		assert.Equal(t, "myns:feet", unit.Measures[1].Value)
		assert.Equal(t, "feet", unit.Measures[1].String())

		assert.Nil(t, unit.Divide)

		assert.Equal(t, "feet * feet", unit.String())
	})

	t.Run("ratio of simple measures", func(t *testing.T) {
		// language=xml
		unitXml := `<unit>
			<divide>
				<unitNumerator>
					<measure>iso4127:USD</measure>
				</unitNumerator>
				<unitDenominator>
					<measure>shares</measure>
				</unitDenominator>
			</divide>
		</unit>`

		var unit Unit
		require.NoError(t, xml.Unmarshal([]byte(unitXml), &unit))

		require.Len(t, unit.Measures, 0)

		assert.NotNil(t, unit.Divide)
		assert.Len(t, unit.Divide.Numerator, 1)
		assert.Equal(t, "iso4127:USD", unit.Divide.Numerator[0].Value)

		assert.Len(t, unit.Divide.Denominator, 1)
		assert.Equal(t, "shares", unit.Divide.Denominator[0].Value)

		assert.Equal(t, "USD / shares", unit.String())
	})

	t.Run("ratio of products of measures", func(t *testing.T) {
		// language=xml
		unitXml := `<unit>
			<divide>
				<unitNumerator>
					<measure>iso4127:USD</measure>
				</unitNumerator>
				<unitDenominator>
					<measure>myns:feet</measure>
					<measure>myns:feet</measure>
				</unitDenominator>
			</divide>
		</unit>`

		var unit Unit
		require.NoError(t, xml.Unmarshal([]byte(unitXml), &unit))

		require.Len(t, unit.Measures, 0)

		assert.NotNil(t, unit.Divide)
		assert.Len(t, unit.Divide.Numerator, 1)
		assert.Equal(t, "iso4127:USD", unit.Divide.Numerator[0].Value)

		assert.Len(t, unit.Divide.Denominator, 2)
		assert.Equal(t, "myns:feet", unit.Divide.Denominator[0].Value)
		assert.Equal(t, "myns:feet", unit.Divide.Denominator[1].Value)

		assert.Equal(t, "USD / feet * feet", unit.String())
	})
}
