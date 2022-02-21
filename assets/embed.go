package assets

import (
	"embed"
)

var (
	//go:embed src/*
	Embedded embed.FS
)
