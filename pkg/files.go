package files

import (
	"errors"
	"path"
	"regexp"
	"strings"
)

// ExtractExtFromURL extracts file's ext from given url
func ExtractExtFromURL(url string) (string, error) {
	// regex to remove any trailing params/hashes
	r, _ := regexp.Compile("\\.?(?P<ext>\\w+)")
	url = path.Ext(strings.Trim(url, "."))
	rz := r.FindStringSubmatch(url)
	if rz == nil {
		return "", errors.New("Can not extract ext")
	}
	return rz[len(rz)-1], nil
}

// ExtractNameFromURL extracts file's name from given URL
func ExtractNameFromURL(url string) (string, error) {
	ext, err := ExtractExtFromURL(url)
	if err != nil {
		return "", err
	}

	// manually add ext to base name to remove any extra query params/hashes
	url = path.Base(strings.Trim(url, "."))
	url = strings.Split(url, ".")[0]

	return url + "." + ext, nil
}
