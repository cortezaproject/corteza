package assets

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path"

	"go.uber.org/zap"
)

var (
	//go:embed src/*
	ff embed.FS
)

// Files return files from assets path or falls back to embedded files
//
// @todo try to find a way to merge this with auth/assets
func Files(log *zap.Logger, aPath string) (files fs.FS) {
	var err error
	if len(aPath) > 0 {
		if files, err = fromPath(aPath); err != nil {
			// log warning but fallback to embedded assets
			log.Warn(
				fmt.Sprintf("failed to use custom assets path (HTTP_SERVER_ASSETS_PATH=%s)", aPath),
				zap.Error(err),
			)
		}
	}

	if files == nil {
		aPath = "embedded"
		files, err = fs.Sub(ff, "src")
		if err != nil {
			// something is seriously wrong, we might as well panic
			panic(err)
		}
	}

	return
}

func fromPath(path string) (assets fs.FS, err error) {
	// at least favicon file should exist in the custom asset path
	// otherwise we default to embedded files
	const check = "favicon32x32.png"

	var (
		fi os.FileInfo
	)

	if fi, err = os.Stat(path); err != nil {
		return

	}

	if !fi.IsDir() {
		return nil, fmt.Errorf("expecting directory")
	}

	assets = os.DirFS(path)
	if _, err = assets.Open(check); err != nil {
		return nil, err
	}

	return
}

func DirEntries(dir string) (fileNames, subDirs []string, err error) {
	dirEntries, err := fs.ReadDir(ff, path.Join("src", dir))
	if err != nil {
		return nil, nil, err
	}

	for _, dirEntry := range dirEntries {
		fileInfo, err := dirEntry.Info()
		if err != nil {
			return nil, nil, err
		}

		// if the entry is a directory skip it
		if fileInfo.IsDir() {
			subDirs = append(subDirs, dirEntry.Name())
			continue
		}

		fileNames = append(fileNames, dirEntry.Name())
	}

	return fileNames, subDirs, err
}
