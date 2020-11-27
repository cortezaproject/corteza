package mime

import (
	"bufio"
	"io"
)

func JsonL(f io.Reader) (bool, error) {
	// ExtractMimetype fails to detect json if jsonl is used
	// For now check if first rune is {
	r := bufio.NewReader(f)
	rn, _, err := r.ReadRune()
	if err != nil {
		return false, err
	}

	if string(rn) == "{" {
		return true, nil
	}
	return false, nil
}
