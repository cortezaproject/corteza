package decoder

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/stretchr/testify/require"
)

func TestMapify(t *testing.T) {
	t.Run("Fail if lengths missmatch", func(t *testing.T) {
		var h []string
		h = append(h, "h1", "h2")

		var v []string
		v = append(v, "v1")

		require.True(t,
			mapify(h, v) == nil,
			"Value should be nil",
		)

		require.True(t,
			mapify(v, h) == nil,
			"Value should be nil",
		)
	})

	t.Run("Successfully mapped", func(t *testing.T) {
		var h []string
		h = append(h, "h1", "h2")

		var v []string
		v = append(v, "v1", "v2")

		mpd := mapify(h, v)
		require.True(t,
			len(mpd) == 2,
			fmt.Sprintf("Invalid length %d; should be %d", len(mpd), 2),
		)

		require.True(t,
			mpd["h1"] == "v1" && mpd["h2"] == "v2",
			"Invalid values",
		)
	})
}

func TestSetSystemField(t *testing.T) {
	t.Run("Correctly determine & set", func(t *testing.T) {
		r := &types.Record{}
		name := "recordID"
		value := "123"
		is, err := setSystemField(r, name, value)
		require.True(t,
			err == nil,
			"Returned with error",
		)

		require.True(t,
			is,
			"Couldn't determine it's a system field",
		)

		require.True(t,
			r.ID == 123,
			fmt.Sprintf("Determined value (%d) not valid; should be %s", r.ID, value),
		)
	})

	t.Run("Correctly determine that it's not", func(t *testing.T) {
		r := &types.Record{}
		name := "customField"
		value := "123"
		is, err := setSystemField(r, name, value)
		require.True(t,
			err == nil,
			"Returned with error",
		)

		require.True(t,
			!is,
			"Couldn't determine it's not a system field",
		)
	})
}

func TestRecords(t *testing.T) {
	testFields := make(map[string]string)
	testFields["f1"] = "f1"
	testFields["f2"] = "f2"
	testFields["ID"] = "ID"

	t.Run("Flat reader", func(t *testing.T) {
		fr := NewFlatReader(csv.NewReader(strings.NewReader(testCSV)), nil)
		row := 0
		fr.Records(testFields, func(mod *types.Record) error {
			row++
			require.True(t,
				len(mod.Values) == 2,
				"Not enough values",
			)

			require.True(t,
				mod.ID == uint64(row),
				"Not enough values",
			)

			return nil
		})
	})

	t.Run("Structured decoder", func(t *testing.T) {
		sd := NewStructuredDecoder(json.NewDecoder(strings.NewReader(testJSONL)), nil)
		row := 0
		sd.Records(testFields, func(mod *types.Record) error {
			row++
			require.True(t,
				len(mod.Values) == 2,
				"Not enough values",
			)

			require.True(t,
				mod.ID == uint64(row),
				"Not enough values",
			)

			return nil
		})
	})
}
