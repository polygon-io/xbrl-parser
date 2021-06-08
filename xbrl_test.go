package xbrl

import (
	"encoding/xml"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalXBRL(t *testing.T) {
	t.Run("real-world xbrl from 2021", func(t *testing.T) {
		f, err := os.Open("test_data/aapl-20210327_htm.xml")
		require.NoError(t, err)
		defer f.Close()

		var content XBRL
		decoder := xml.NewDecoder(f)

		require.NoError(t, decoder.Decode(&content))
		require.NoError(t, content.Validate())

		// Not using assert.Len here because with the output is painful for large maps/slices
		assert.Equal(t, 283, len(content.ContextsByID))
		assert.Equal(t, 9, len(content.UnitsByID))
		assert.Equal(t, 1070, len(content.Facts))
	})

	t.Run("real-world xbrl from 2004", func(t *testing.T) {
		f, err := os.Open("test_data/edgr-2004_10k.xml") // The very first XBRL submission to the SEC!
		require.NoError(t, err)
		defer f.Close()

		var content XBRL
		decoder := xml.NewDecoder(f)
		decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
			return input, nil
		}

		require.NoError(t, decoder.Decode(&content))
		require.NoError(t, content.Validate())

		// Not using assert.Len here because with the output is painful for large maps/slices
		assert.Equal(t, 4, len(content.ContextsByID))
		assert.Equal(t, 2, len(content.UnitsByID))
		assert.Equal(t, 154, len(content.Facts))
	})

	t.Run("simple xbrl happy path", func(t *testing.T) {
		xbrlBytes, err := os.ReadFile("test_data/simple_xbrl.xml")
		require.NoError(t, err)

		var content XBRL

		require.NoError(t, xml.Unmarshal(xbrlBytes, &content))
		require.NoError(t, content.Validate())

		require.Len(t, content.ContextsByID, 1)
		expectedContext := Context{
			ID: "c1",
			Period: Period{
				Instant: stringPtr("2021-04-16"),
			},
			Entity: Entity{
				Identifier: Identifier{
					Scheme: "http://www.sec.gov/CIK",
					Value:  "0000320193",
				},
			},
		}

		assert.Equal(t, expectedContext, content.ContextsByID["c1"])

		require.Len(t, content.UnitsByID, 1)
		expectedUnit := Unit{
			ID:       "u1",
			Measures: Measures{{Value: "shares"}},
		}

		assert.Equal(t, expectedUnit, content.UnitsByID["u1"])

		require.Len(t, content.Facts, 2)
		expectedFacts := []Fact{
			{
				XMLName:    xml.Name{Space: "http://www.xbrl.org/us/gaap/ci/2003/usfr-ci-2003", Local: "assets"},
				ContextRef: "c1",
				UnitRef:    stringPtr("u1"),
				Precision:  stringPtr("3"),
				ValueStr:   stringPtr("727"),
			},
			{
				XMLName:    xml.Name{Space: "fakens", Local: "textItem"},
				ContextRef: "c1",
				ValueStr:   stringPtr("this is a text item"),
			},
		}

		assert.Equal(t, expectedFacts, content.Facts)
	})

	t.Run("invalid xbrl", func(t *testing.T) {
		xbrlBytes, err := os.ReadFile("test_data/invalid_xbrl.xml")
		require.NoError(t, err)

		var content XBRL

		require.NoError(t, xml.Unmarshal(xbrlBytes, &content))
		assert.Error(t, content.Validate())
	})
}

func stringPtr(str string) *string {
	return &str
}
