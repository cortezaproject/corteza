package main

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"

	city311 "github.com/cortezaproject/corteza/server/compose/types/city311"
)

func main() {
	out := flag.String("out", "compose/types/city311", "artifact output directory")
	flag.Parse()
	writeJSON(filepath.Join(*out, "contract.json"), city311.NewContractDocument())
	writeJSON(filepath.Join(*out, "openapi.json"), city311.NewOpenAPIDocument())
}

func writeJSON(path string, value interface{}) {
	raw, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		panic(err)
	}
	raw = append(raw, '\n')
	if err = os.WriteFile(path, raw, 0o644); err != nil {
		panic(err)
	}
}
