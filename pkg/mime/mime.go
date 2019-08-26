package mime

import (
	"bufio"
	"io"

	"github.com/gabriel-vasile/mimetype"
)

func Type(file io.ReadSeeker) (mt string, ext string, err error) {
	if _, err = file.Seek(0, 0); err != nil {
		return
	}

	// Make sure we rewind when we're done
	defer file.Seek(0, 0)
	return mimetype.DetectReader(file)
}

func JsonL(file io.ReadSeeker) (bool, error) {
	// ExtractMimetype fails to detect json if jsonl is used
	// For now check if first rune is {
	r := bufio.NewReader(file)
	rn, _, err := r.ReadRune()
	defer file.Seek(0, 0)
	if err != nil {
		return false, err
	}

	if string(rn) == "{" {
		return true, nil
	}
	return false, nil
}
