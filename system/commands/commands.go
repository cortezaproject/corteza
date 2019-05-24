package commands

import (
	"fmt"
	"os"
)

func exit(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}
