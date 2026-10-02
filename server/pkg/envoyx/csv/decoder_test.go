package csv

import (
	"io"
	"strings"
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/envoyx/datasource"
	"github.com/stretchr/testify/require"
)

func TestDecoder(t *testing.T) {
	req := require.New(t)

	t.Run("init & meta", func(t *testing.T) {
		dc, err := Decoder(testReader(), "test.csv", nil)
		req.NoError(err)

		req.Equal(uint64(3), dc.Count())
	})

	t.Run("fields", func(t *testing.T) {
		dc, err := Decoder(testReader(), "test.csv", nil)
		req.NoError(err)

		hh := dc.Fields()
		req.Contains(hh, "f1")
		req.Contains(hh, "f2")
		req.Contains(hh, "f3")
	})

	t.Run("iterate", func(t *testing.T) {
		dc, err := Decoder(testReader(), "test.csv", nil)
		req.NoError(err)

		aux := make(datasource.RawRecord)
		var more bool

		more, err = dc.Next(nil, aux)
		req.NoError(err)
		req.True(more)
		req.Equal("r1f1", aux["f1"].Values[0])
		req.Equal("r1f2", aux["f2"].Values[0])
		req.Equal("r1f3", aux["f3"].Values[0])

		more, err = dc.Next(nil, aux)
		req.NoError(err)
		req.True(more)
		req.Equal("r2f1", aux["f1"].Values[0])
		req.Equal("r2f2", aux["f2"].Values[0])
		req.Equal("r2f3", aux["f3"].Values[0])

		more, err = dc.Next(nil, aux)
		req.NoError(err)
		req.True(more)
		req.Equal("r3f1", aux["f1"].Values[0])
		req.Equal("r3f2", aux["f2"].Values[0])
		req.Equal("r3f3", aux["f3"].Values[0])

		more, err = dc.Next(nil, aux)
		req.NoError(err)
		req.False(more)
	})
}

func TestParseComplexCSVCell(t *testing.T) {
	t.Run("plain value", func(t *testing.T) {
		d := &decoder{}

		out := d.parseComplexCSVCell("cell")
		require.Equal(t, []string{"cell"}, out)
	})

	t.Run("multi value without brackets", func(t *testing.T) {
		d := &decoder{
			multiValueDelimiter: ";",
		}

		out := d.parseComplexCSVCell("v1;v2")
		require.Equal(t, []string{"v1", "v2"}, out)
	})

	t.Run("multi value with brackets", func(t *testing.T) {
		d := &decoder{
			multiValueDelimiter: ";",
			multiValueBrackets:  true,
		}

		out := d.parseComplexCSVCell("[v1;v2]")
		require.Equal(t, []string{"v1", "v2"}, out)
	})

	t.Run("multi value wrong delimiter", func(t *testing.T) {
		d := &decoder{
			multiValueDelimiter: ",",
			multiValueBrackets:  true,
		}

		out := d.parseComplexCSVCell("[v1;v2]")
		require.Equal(t, []string{"v1;v2"}, out)
	})

	t.Run("json", func(t *testing.T) {
		d := &decoder{}

		out := d.parseComplexCSVCell("{}")
		require.Equal(t, []string{"{}"}, out)
	})

	t.Run("multi value json", func(t *testing.T) {
		d := &decoder{
			multiValueDelimiter: ";",
			multiValueBrackets:  true,
		}

		out := d.parseComplexCSVCell("[{};{}]")
		require.Equal(t, []string{"{}", "{}"}, out)
	})

	t.Run("single value geometry json", func(t *testing.T) {
		d := &decoder{
			multiValueDelimiter: ";",
		}

		out := d.parseComplexCSVCell(`{"coordinates":[47.754098,1.9335938]}`)
		require.Equal(t, []string{`{"coordinates":[47.754098,1.9335938]}`}, out)
	})

	t.Run("multi value geometry json without brackets", func(t *testing.T) {
		d := &decoder{
			multiValueDelimiter: ";",
		}

		out := d.parseComplexCSVCell(`{"coordinates":[47.754098,1.9335938]};{"coordinates":[9.7956776,23.3789063]}`)
		require.Equal(t, []string{
			`{"coordinates":[47.754098,1.9335938]}`,
			`{"coordinates":[9.7956776,23.3789063]}`,
		}, out)
	})

	t.Run("multi value geometry json with brackets", func(t *testing.T) {
		d := &decoder{
			multiValueDelimiter: ";",
			multiValueBrackets:  true,
		}

		out := d.parseComplexCSVCell(`[{"coordinates":[47.754098,1.9335938]};{"coordinates":[9.7956776,23.3789063]}]`)
		require.Equal(t, []string{
			`{"coordinates":[47.754098,1.9335938]}`,
			`{"coordinates":[9.7956776,23.3789063]}`,
		}, out)
	})

	t.Run("single value geometry json with comma delimiter", func(t *testing.T) {
		d := &decoder{
			multiValueDelimiter: ",",
		}

		// commas inside the JSON must not split the value
		out := d.parseComplexCSVCell(`{"coordinates":[47.754098,1.9335938]}`)
		require.Equal(t, []string{`{"coordinates":[47.754098,1.9335938]}`}, out)
	})

	t.Run("single value json with nested objects", func(t *testing.T) {
		d := &decoder{
			multiValueDelimiter: ",",
		}

		// "},{" inside the JSON must not split the value
		out := d.parseComplexCSVCell(`{"points":[{"x":1},{"y":2}]}`)
		require.Equal(t, []string{`{"points":[{"x":1},{"y":2}]}`}, out)
	})

	t.Run("multi value json with nested objects and braces in strings", func(t *testing.T) {
		d := &decoder{
			multiValueDelimiter: ";",
		}

		out := d.parseComplexCSVCell(`{"a":[{"x":1},{"y":2}],"s":"};{"};{"b":2}`)
		require.Equal(t, []string{
			`{"a":[{"x":1},{"y":2}],"s":"};{"}`,
			`{"b":2}`,
		}, out)
	})

	t.Run("multi value json split is delimiter agnostic", func(t *testing.T) {
		// exports default to ";" while an unconfigured decoder falls back to
		// ","; complete-object splitting covers both
		d := &decoder{}

		out := d.parseComplexCSVCell(`{"a":1};{"b":2}`)
		require.Equal(t, []string{`{"a":1}`, `{"b":2}`}, out)

		out = d.parseComplexCSVCell(`{"a":1},{"b":2}`)
		require.Equal(t, []string{`{"a":1}`, `{"b":2}`}, out)
	})

	t.Run("malformed json cell returned whole", func(t *testing.T) {
		d := &decoder{
			multiValueDelimiter: ";",
		}

		out := d.parseComplexCSVCell(`{"a":{1};{"b"}`)
		require.Equal(t, []string{`{"a":{1};{"b"}`}, out)
	})
}

func testReader() io.Reader {
	src := `f1,f2,f3
r1f1,r1f2,r1f3
r2f1,r2f2,r2f3
r3f1,r3f2,r3f3`

	return strings.NewReader(src)
}
