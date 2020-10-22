package yaml

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

func parseDocument(name string) (*Document, error) {
	doc := &Document{}
	f, err := os.Open(fmt.Sprintf("testdata/%s.yaml", name))
	if err != nil {
		return nil, err
	}

	return doc, yaml.NewDecoder(f).Decode(doc)
}
