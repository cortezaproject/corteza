package cli

import (
	"fmt"
	"os"
)

func HandleError(err error) {
	if err == nil {
		return
	}

	_, _ = fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}
