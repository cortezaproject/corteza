package encoder

import (
	"io"

	"github.com/360EntSecGroup-Skylar/excelize"
)

type (
	excelizeEncoder struct {
		row int
		f   *excelize.File
		w   io.Writer
		ff  []field
	}
)

func NewExcelizeEncoder(w io.Writer, header bool, ff ...field) *excelizeEncoder {
	enc := &excelizeEncoder{
		f:  excelize.NewFile(),
		w:  w,
		ff: ff,
	}

	if header {
		enc.writeHeader()
	}

	return enc
}

func (enc *excelizeEncoder) Flush() {
	_ = enc.f.Write(enc.w)
}

// Returns current row + column to alphanumeric cell name
func (enc excelizeEncoder) pos(col int) string {
	cn, _ := excelize.CoordinatesToCellName(col, enc.row)
	return cn
}

func (enc excelizeEncoder) sheet() string {
	return "Sheet1"
}

func (enc *excelizeEncoder) writeHeader() {
	enc.row++
	for p := range enc.ff {
		_ = enc.f.SetCellStr(enc.sheet(), enc.pos(p+1), enc.ff[p].name)
	}
}
