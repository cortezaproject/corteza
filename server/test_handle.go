package main

import (
	"fmt"
	"regexp"
)

var isValidHandle = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_\-\.]*[a-zA-Z0-9_]$`).MatchString

func main() {
	fmt.Println("agent-:", isValidHandle("agent-"))
	fmt.Println("agent:", isValidHandle("agent"))
}
