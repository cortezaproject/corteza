package sass

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"path"
	"regexp"
	"strings"

	"github.com/bep/godartsass/v2"
	"github.com/cortezaproject/corteza/server/assets"
	"github.com/cortezaproject/corteza/server/pkg/logger"
	"go.uber.org/zap"
)

var (
	StylesheetCache      = newStylesheetCache()
	sassVariablesPattern = regexp.MustCompile(`(\$[a-zA-Z_-]+):\s*([^;]+);`)
)

// DefaultCSS contains css contents from minified-custom.css and custom-css editor content
func DefaultCSS(log *zap.Logger, customCSS string) string {
	var strBuilder strings.Builder

	//read css contents from minified-custom.css
	minifiedCSSFile, err := assets.Files(logger.Default(), "").Open(path.Join("css", "minified-custom.css"))
	if err != nil {
		log.Error("failed to open minified-custom.css file", zap.Error(err))
	}

	reader := bufio.NewReader(minifiedCSSFile)
	_, err = io.Copy(&strBuilder, reader)
	if err != nil {
		log.Error("failed to read css content from minified-custom.css", zap.Error(err))
	}

	if customCSS != "" {
		strBuilder.WriteString(customCSS)
	}

	processedCSS := strBuilder.String()

	StylesheetCache.Set(
		map[string]string{
			"css": processedCSS,
		},
	)

	return processedCSS
}

func Transpiler(log *zap.Logger, brandingSass, customCSS, sassDirPath string, transpiler *godartsass.Transpiler) error {
	var stringsBuilder strings.Builder

	if brandingSass != "" {
		err := jsonToSass(brandingSass, &stringsBuilder)
		if err != nil {
			log.Error("failed to unmarshal branding sass variables", zap.Error(err))
			return err
		}
	}

	if customCSS != "" {
		// Get SASS variables from the custom CSS editor and give them precedence over branding variables
		customVariables := sassVariablesPattern.FindAllString(customCSS, -1)

		for _, customVariable := range customVariables {
			stringsBuilder.WriteString(fmt.Sprintf("%s \n", customVariable))
		}
	}

	// get Boostrap, bootstrap-vue and custom variables sass content
	mainSass, err := readSassFiles(log, "scss")
	if err != nil {
		return err
	}
	stringsBuilder.WriteString(mainSass)

	// when a user provides sets WEBAPP_SCSS_DIR_PATH environment variable
	if sassDirPath != "" {
		customSass, err := readSassFiles(log, sassDirPath)
		if err != nil {
			return err
		}
		stringsBuilder.WriteString(customSass)
	}

	if customCSS != "" {
		//Custom CSS editor selector block
		selectorBlock := sassVariablesPattern.ReplaceAllString(customCSS, "")
		stringsBuilder.WriteString(selectorBlock)
	}

	// compute sass content to CSS
	args := godartsass.Args{
		Source: stringsBuilder.String(),
	}
	execute, err := transpiler.Execute(args)
	if err != nil {
		log.Error("sass compilation failure", zap.Error(err))
		return err
	}

	// save computed css to Stylesheet in-memory cache
	StylesheetCache.Set(
		map[string]string{
			"css": execute.CSS,
		},
	)

	return nil
}

func DartSass(log *zap.Logger) *godartsass.Transpiler {
	transpiler, err := godartsass.Start(godartsass.Options{
		DartSassEmbeddedFilename: "sass",
	})

	if err != nil {
		log.Warn("dart sass is not installed in your system", zap.Error(err))
		return nil
	}

	return transpiler
}

// readSassFiles reads SASS contents from provided embedded directory and subdirectories then converts them to a string
func readSassFiles(log *zap.Logger, dirPath string) (string, error) {
	var stringsBuilder strings.Builder

	filenames, subDirs, err := assets.DirEntries(dirPath)
	if err != nil {
		log.Error(fmt.Sprintf("failed to read assets/src/%s entries", dirPath), zap.Error(err))
		return "", err
	}

	if len(filenames) > 0 {
		err := readSassContents(log, dirPath, filenames, &stringsBuilder)
		if err != nil {
			return "", err
		}
	}

	if len(subDirs) > 0 {
		for _, subDir := range subDirs {
			sassContents, err := readSassFiles(log, path.Join(dirPath, subDir))
			if err != nil {
				return "", err
			}
			stringsBuilder.WriteString(sassContents)
		}
	}

	return stringsBuilder.String(), nil
}

func readSassContents(log *zap.Logger, dirPath string, filenames []string, stringsBuilder *strings.Builder) error {
	assetFiles := assets.Files(logger.Default(), "")

	for _, fileName := range filenames {
		open, err := assetFiles.Open(path.Join(dirPath, fileName))
		if err != nil {
			log.Error(fmt.Sprintf("failed to open asset %s file", fileName), zap.Error(err))
			return err
		}

		reader := bufio.NewReader(open)
		_, err = io.Copy(stringsBuilder, reader)
		if err != nil {
			log.Error(fmt.Sprintf("failed to copy sass content from %s", fileName), zap.Error(err))
			return err
		}
	}

	return nil
}

// jsonToSass converts JSON string to SASS variable assignment string
func jsonToSass(jsonStr string, strBuilder *strings.Builder) error {
	var colorMap map[string]string

	err := json.Unmarshal([]byte(jsonStr), &colorMap)
	if err != nil {
		return err
	}

	for key, value := range colorMap {
		strBuilder.WriteString(fmt.Sprintf("$%s: %s;\n", key, value))
	}

	return nil
}
