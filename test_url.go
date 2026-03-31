package main

import (
"fmt"
"net/url"
)

func main() {
	base, _ := url.Parse("https://api.github.com:443/repos/tjerman/dev-test")
	rel, _ := url.Parse("/issues")
	fmt.Println(base.ResolveReference(rel).String())
}
