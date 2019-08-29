package decoder

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/internal/test"
)

func TestMapify(t *testing.T) {
	t.Run("Fail if lengths missmatch", func(t *testing.T) {
		var h []string
		h = append(h, "h1", "h2")

		var v []string
		v = append(v, "v1")

		test.Assert(t,
			mapify(h, v) == nil,
			"Value should be nil",
		)

		test.Assert(t,
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
		test.Assert(t,
			len(mpd) == 2,
			fmt.Sprintf("Invalid length %d; should be %d", len(mpd), 2),
		)

		test.Assert(t,
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
		test.Assert(t,
			err == nil,
			"Returned with error",
		)

		test.Assert(t,
			is,
			"Couldn't determine it's a system field",
		)

		test.Assert(t,
			r.ID == 123,
			fmt.Sprintf("Determined value (%d) not valid; should be %s", r.ID, value),
		)
	})

	t.Run("Correctly determine that it's not", func(t *testing.T) {
		r := &types.Record{}
		name := "customField"
		value := "123"
		is, err := setSystemField(r, name, value)
		test.Assert(t,
			err == nil,
			"Returned with error",
		)

		test.Assert(t,
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
			test.Assert(t,
				len(mod.Values) == 2,
				"Not enough values",
			)

			test.Assert(t,
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
			test.Assert(t,
				len(mod.Values) == 2,
				"Not enough values",
			)

			test.Assert(t,
				mod.ID == uint64(row),
				"Not enough values",
			)

			return nil
		})
	})
}
