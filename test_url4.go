package main

import (
"fmt"
"net/url"
)

func buildURL(baseURL, path string) (string, error) {
	if baseURL == "" {
		return path, nil
	}

	rel, err := url.Parse(path)
	if err != nil {
		return "", fmt.Errorf("invalid path: %w", err)
	}

	if rel.IsAbs() {
		return rel.String(), nil
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}

	base.Path, err = url.JoinPath(base.Path, rel.Path)
	if err != nil {
		return "", err
	}

	rel.Path = ""
	return base.ResolveReference(rel).String(), nil
}

func main() {
	res, _ := buildURL("https://api.github.com:443/repos/tjerman/dev-test", "/issues?foo=bar")
	fmt.Println(res)
    res2, _ := buildURL("https://api.github.com", "/issues")
	fmt.Println(res2)
    res3, _ := buildURL("https://api.github.com", "https://api.github.com/issues")
	fmt.Println(res3)
}
