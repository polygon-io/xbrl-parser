package xbrl

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalContext(t *testing.T) {
	t.Run("instant period | no segments", func(t *testing.T) {
		// language=xml
		contextXML := `<context id="i132ffb9faa364951a8fec32b601d621b_I20210416">
    <entity>
        <identifier scheme="http://www.sec.gov/CIK">0000320193</identifier>
    </entity>
    <period>
        <instant>2021-04-16</instant>
    </period>
</context>`

		var context Context
		require.NoError(t, xml.Unmarshal([]byte(contextXML), &context))

		assert.Equal(t, "http://www.sec.gov/CIK", context.Entity.Identifier.Scheme)
		assert.Equal(t, "0000320193", context.Entity.Identifier.Value)
		assert.Equal(t, PeriodTypeInstant, context.Period.Type())
		assert.Equal(t, "2021-04-16", *context.Period.Instant)
	})

	t.Run("forever period | no segments", func(t *testing.T) {
		// language=xml
		contextXML := `<context id="i132ffb9faa364951a8fec32b601d621b_I20210416">
    <entity>
        <identifier scheme="http://www.sec.gov/CIK">0000320193</identifier>
    </entity>
    <period>
        <forever/>
    </period>
</context>`

		var context Context
		require.NoError(t, xml.Unmarshal([]byte(contextXML), &context))

		assert.Equal(t, "http://www.sec.gov/CIK", context.Entity.Identifier.Scheme)
		assert.Equal(t, "0000320193", context.Entity.Identifier.Value)
		assert.Equal(t, PeriodTypeForever, context.Period.Type())
	})

	t.Run("duration period | has segments", func(t *testing.T) {
		// language=xml
		contextXML := `<context id="iff44040cd61344d085f7a2b7a1076cb1_D20200927-20210327">
    <entity>
        <identifier scheme="http://www.sec.gov/CIK">0000320193</identifier>
        <segment>
            <xbrldi:explicitMember dimension="us-gaap:StatementClassOfStockAxis">us-gaap:CommonStockMember</xbrldi:explicitMember>
            <myns:cool_segment>I follow my own rules</myns:cool_segment>
        </segment>
    </entity>
    <period>
        <startDate>2020-09-27</startDate>
        <endDate>2021-03-27</endDate>
    </period>
</context>`

		var context Context
		require.NoError(t, xml.Unmarshal([]byte(contextXML), &context))

		assert.Equal(t, "http://www.sec.gov/CIK", context.Entity.Identifier.Scheme)
		assert.Equal(t, "0000320193", context.Entity.Identifier.Value)

		assert.Len(t, context.Entity.Segments, 2)
		assert.Equal(t, xml.Name{Space: "xbrldi", Local: "explicitMember"}, context.Entity.Segments[0].XMLName)
		assert.Equal(t, []xml.Attr{{Name: xml.Name{Local: "dimension"}, Value: "us-gaap:StatementClassOfStockAxis"}}, context.Entity.Segments[0].Attributes)
		assert.Equal(t, "us-gaap:CommonStockMember", context.Entity.Segments[0].Value)

		assert.Equal(t, xml.Name{Space: "myns", Local: "cool_segment"}, context.Entity.Segments[1].XMLName)
		assert.Empty(t, context.Entity.Segments[1].Attributes)
		assert.Equal(t, "I follow my own rules", context.Entity.Segments[1].Value)

		assert.Equal(t, PeriodTypeDuration, context.Period.Type())
		assert.Equal(t, "2020-09-27", *context.Period.StartDate)
		assert.Equal(t, "2021-03-27", *context.Period.EndDate)
	})
}
