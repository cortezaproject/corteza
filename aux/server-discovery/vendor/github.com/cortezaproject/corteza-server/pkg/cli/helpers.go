package cli

import (
	"errors"
	"fmt"
	"os"
)

func HandleError(err error) {
	if err = unwrapGeneric(err); err == nil {
		return
	}

	_, _ = fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}

func unwrapGeneric(err error) error {
	for {
		g, ok := err.(interface{ IsGeneric() bool })
		if ok && g != nil && g.IsGeneric() {
			err = errors.Unwrap(err)
			continue
		}

		return err
	}
}
