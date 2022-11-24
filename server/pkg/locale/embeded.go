package locale

///////////////////////////////////////////////////////////////////////////////
// This helps us import translations into corteza-server as a module
// dependency

import (
	"embed"
	"io/fs"
)

//go:embed src/*
var languages embed.FS

// LoadEmbedded returns embedded translation files
// as a virtual filesystem
func LoadEmbedded() fs.FS {
	sub, _ := fs.Sub(languages, "src")
	return sub
}
