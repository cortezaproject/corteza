package main

import (
"fmt"
"net/url"
)

func main() {
	res, _ := url.JoinPath("https://api.github.com:443/repos/tjerman/dev-test", "/issues")
	fmt.Println(res)
	res2, _ := url.JoinPath("https://api.github.com:443/repos/tjerman/dev-test", "issues")
	fmt.Println(res2)
}
