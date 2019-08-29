package decoder

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/cortezaproject/corteza-server/internal/test"
)

func makeReadSeeker(c string) io.ReadSeeker {
	bb := []byte(c)
	return bytes.NewReader(bb)
}

const (
	testCSV   string = "f1,f2,ID\nr1v1,r1v2,1\n"
	testJSONL string = "{ \"f1\": \"nr1v1\", \"f2\": \"r1v2\", \"ID\": \"1\" }\n"
)

func TestEntryCount(t *testing.T) {
	t.Run("Flat reader", func(t *testing.T) {
		rs := makeReadSeeker(testCSV)
		fr := NewFlatReader(csv.NewReader(rs), rs)

		c, err := fr.EntryCount()
		test.Assert(t,
			c == 1,
			fmt.Sprintf("Invalid number of entries determines; found %d, expected %d", c, 1),
		)

		test.Assert(t,
			err == nil,
			"Returned with error",
		)
	})

	t.Run("Structured decoder", func(t *testing.T) {
		rs := makeReadSeeker(testJSONL)
		sd := NewStructuredDecoder(json.NewDecoder(rs), rs)

		c, err := sd.EntryCount()
		test.Assert(t,
			c == 1,
			fmt.Sprintf("Invalid number of entries determines; found %d, expected %d", c, 1),
		)

		test.Assert(t,
			err == nil,
			"Returned with error",
		)
	})
}

func TestGet(t *testing.T) {
	t.Run("Flat reader", func(t *testing.T) {
		rs := makeReadSeeker(testCSV)
		fr := NewFlatReader(csv.NewReader(rs), rs)

		// dump first line
		fr.get(func(f []string) error { return nil })
		err := fr.get(func(f []string) error {
			return errors.New("called")
		})
		test.Assert(t,
			err != nil,
			"Error should be returned to indicate that get did read",
		)

		err = fr.get(func(f []string) error {
			return errors.New("called")
		})
		test.Assert(t,
			err == nil,
			"Error should NOT be returned to indicate that get didn't read",
		)
	})

	t.Run("Structured decoder", func(t *testing.T) {
		rs := makeReadSeeker(testJSONL)
		sd := NewStructuredDecoder(json.NewDecoder(rs), rs)

		err := sd.get(func(f map[string]interface{}) error {
			return errors.New("called")
		})
		test.Assert(t,
			err != nil,
			"Error should be returned to indicate that get did read",
		)

		err = sd.get(func(f map[string]interface{}) error {
			return errors.New("called")
		})
		test.Assert(t,
			err == nil,
			"Error should NOT be returned to indicate that get didn't read",
		)
	})
}

func TestWalk(t *testing.T) {
	t.Run("Flat reader", func(t *testing.T) {
		rs := makeReadSeeker(testCSV)
		fr := NewFlatReader(csv.NewReader(rs), rs)

		i := 0
		err := fr.walk(func(f []string) error {
			i++
			return nil
		})

		test.Assert(t,
			i == 2,
			"Invalid number of reads",
		)

		test.Assert(t,
			err == nil,
			"Returned with error",
		)
	})

	t.Run("Structured decoder", func(t *testing.T) {
		rs := makeReadSeeker(testJSONL)
		sd := NewStructuredDecoder(json.NewDecoder(rs), rs)

		i := 0
		err := sd.walk(func(f map[string]interface{}) error {
			i++
			return nil
		})

		test.Assert(t,
			i == 1,
			"Invalid number of reads",
		)

		test.Assert(t,
			err == nil,
			"Returned with error",
		)
	})
}

func TestHeader(t *testing.T) {
	t.Run("Flat reader", func(t *testing.T) {
		rs := makeReadSeeker(testCSV)
		fr := NewFlatReader(csv.NewReader(rs), rs)

		h := fr.Header()
		test.Assert(t,
			len(h) == 3,
			"Invalid number of header fields",
		)

		expect := [...]string{"f1", "f2", "ID"}
		for i, h := range h {
			test.Assert(t,
				h == expect[i],
				"Invalid header value",
			)
		}
	})

	t.Run("Structured decoder", func(t *testing.T) {
		rs := makeReadSeeker(testJSONL)
		sd := NewStructuredDecoder(json.NewDecoder(rs), rs)

		h := sd.Header()
		test.Assert(t,
			len(h) == 3,
			"Invalid number of header fields",
		)
	})
}
