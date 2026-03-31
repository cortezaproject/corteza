package main

import (
"fmt"
"net/url"
)

func main() {
	base, _ := url.Parse("https://api.github.com/repos/tjerman/dev-test")
	rel, _ := url.Parse("https://api.github.com/repos/tjerman/dev-test/issues")
	fmt.Println(base.ResolveReference(rel).String())
}
