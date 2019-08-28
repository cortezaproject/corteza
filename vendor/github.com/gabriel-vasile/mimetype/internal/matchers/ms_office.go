package matchers

import (
	"bytes"
	"fmt"
	"strings"
)

// Xlsx matches a Microsoft Excel 2007 file.
func Xlsx(in []byte) bool {
	return bytes.Contains(in, []byte("xl/"))
}

// Docx matches a Microsoft Office 2007 file.
func Docx(in []byte) bool {
	return bytes.Contains(in, []byte("word/"))
}

// Pptx matches a Microsoft PowerPoint 2007 file.
func Pptx(in []byte) bool {
	return bytes.Contains(in, []byte("ppt/"))
}

// Doc matches a Microsoft Office 97-2003 file.
func Doc(in []byte) bool {
	return bytes.HasPrefix(in, []byte{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1})
}

// Ppt  matches a Microsoft PowerPoint 97-2003 file.
func Ppt(in []byte) bool {
	if len(in) < 520 {
		return false
	}

	if fmt.Sprintf("%X", in[:8]) == "D0CF11E0A1B11AE1" {
		offset512 := fmt.Sprintf("%X", in[512:516])
		if offset512 == "A0461DF0" || offset512 == "006E1EF0" || offset512 == "0F00E803" {
			return true
		}
		if offset512 == "FDFFFFFF" && fmt.Sprintf("%x", in[518:520]) == "0000" {
			return true
		}
	}

	return false
}

// Xls  matches a Microsoft Excel 97-2003 file.
func Xls(in []byte) bool {
	if len(in) < 520 {
		return false
	}

	if fmt.Sprintf("%X", in[:8]) == "D0CF11E0A1B11AE1" {
		offset512 := fmt.Sprintf("%X", in[512:520])
		subheaders := []string{
			"0908100000060500",
			"FDFFFFFF10",
			"FDFFFFFF1F",
			"FDFFFFFF22",
			"FDFFFFFF23",
			"FDFFFFFF28",
			"FDFFFFFF29",
		}
		for _, h := range subheaders {
			if strings.HasPrefix(offset512, h) {
				return true
			}
		}
	}

	return false
}
