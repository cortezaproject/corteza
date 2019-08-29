package count

import (
	"bytes"
	"io"
)

// Lines provides a line count
//
// https://stackoverflow.com/a/24563853
func Lines(r io.ReadSeeker) (count uint64, err error) {
	defer r.Seek(0, 0)
	buf := make([]byte, 32*1024)
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += uint64(bytes.Count(buf[:c], lineSep))

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
