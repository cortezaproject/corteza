package sass

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/bep/godartsass/v2"
	"github.com/cortezaproject/corteza/server/assets"
	"github.com/cortezaproject/corteza/server/pkg/logger"
	"go.uber.org/zap"
)

func GenerateCSS(brandingSass, customCSS string) (string, error) {
	var (
		stringsBuilder strings.Builder
	)

	if brandingSass != "" {
		brandingSassString, err := jsonToSass(brandingSass)
		if err != nil {
			return "", err
		}
		stringsBuilder.WriteString(brandingSassString)
	}

	// Get SASS variables from the custom CSS editor and give them precedence over branding variables
	variablesPattern := regexp.MustCompile(`(\$[a-zA-Z_-]+):\s*([^;]+);`)
	customVariables := variablesPattern.FindAllString(customCSS, -1)

	for _, customVariable := range customVariables {
		stringsBuilder.WriteString(fmt.Sprintf("%s \n", customVariable))
	}

	// get Boostrap, bootstrap-vue and custom variables sass content
	scssFromAssets, err := readSassFiles()
	if err != nil {
		return "", err
	}
	stringsBuilder.WriteString(scssFromAssets)

	//Custom CSS editor selector block
	selectorBlock := variablesPattern.ReplaceAllString(customCSS, "")
	stringsBuilder.WriteString(selectorBlock)

	// Process Sass
	args := godartsass.Args{
		Source: stringsBuilder.String(),
	}

	t, err := godartsass.Start(godartsass.Options{
		DartSassEmbeddedFilename: "sass",
	})

	if err != nil {
		panic(err.Error())
	}

	execute, err := t.Execute(args)
	if err != nil {
		logger.Default().Error("Sass syntax error", zap.Error(err))
		return "", err
	}

	return execute.CSS, nil
}

// readSassFiles reads SASS files from assets and converts them to a string
func readSassFiles() (string, error) {
	var stringsBuilder strings.Builder

	assetFiles := assets.Files(logger.Default(), "")
	fileNames, err := assets.DirFileNames("scss")
	if err != nil {
		logger.Default().Error("failed to read the file names", zap.Error(err))
		return "", err
	}

	for _, fileName := range fileNames {
		open, err := assetFiles.Open("scss/" + fileName)
		if err != nil {
			logger.Default().Error(fmt.Sprintf("failed to open asset %s file", fileName), zap.Error(err))
			return "", err
		}

		reader := bufio.NewReader(open)
		_, err = io.Copy(&stringsBuilder, reader)
		if err != nil {
			logger.Default().Error("failed to write scss content", zap.Error(err))
			return "", err
		}
	}

	return stringsBuilder.String(), nil
}

// jsonToSass converts JSON string to SASS variable assignment string
func jsonToSass(jsonStr string) (string, error) {
	var (
		colorMap       map[string]string
		stringsBuilder strings.Builder
	)

	err := json.Unmarshal([]byte(jsonStr), &colorMap)
	if err != nil {
		return "", err
	}

	for key, value := range colorMap {
		stringsBuilder.WriteString(fmt.Sprintf("$%s: %s;\n", key, value))
	}

	return stringsBuilder.String(), nil
}
