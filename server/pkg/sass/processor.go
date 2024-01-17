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

const (
	GeneralTheme = "general"
	LightTheme   = "light"
	DarkTheme    = "dark"
	SectionRoot  = "root"
	SectionMain  = "main"
	SectionTheme = "theme"
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

	StylesheetCache.Set("default-theme", processedCSS)

	return processedCSS
}

func Transpile(transpiler *godartsass.Transpiler, log *zap.Logger, themeID, themeSASS, customCSS, sassDirPath string) (err error) {
	// process root section
	err = processSass(transpiler, log, SectionRoot, themeID, themeSASS, "", sassDirPath)
	if err != nil {
		return err
	}

	// process main section
	if themeID == LightTheme {
		err = processSass(transpiler, log, SectionMain, themeID, themeSASS, customCSS, sassDirPath)
		if err != nil {
			return err
		}
	}

	//process theme section
	err = processSass(transpiler, log, SectionTheme, themeID, themeSASS, customCSS, sassDirPath)
	if err != nil {
		return err
	}

	return nil
}

func processSass(transpiler *godartsass.Transpiler, log *zap.Logger, section, themeID, themeSASS, customCSS, sassDirPath string) (err error) {
	var (
		stringsBuilder         strings.Builder
		isCustomCssSyntaxValid bool
	)

	// add theme-mode variable to the top of the section
	stringsBuilder.WriteString(fmt.Sprintf("$theme-mode: %s;\n", themeID))

	if themeSASS != "" {
		err = jsonToSass(themeSASS, &stringsBuilder)
		if err != nil {
			log.Error("failed to unmarshal branding sass variables", zap.Error(err))
			return err
		}
	}

	if customCSS != "" {
		_, err = transpileSass(transpiler, customCSS)
		if err != nil {
			log.Error("sass compilation for custom css failed", zap.Error(err))
		} else {
			isCustomCssSyntaxValid = true
			// Get SASS variables from the custom CSS editor and give them precedence over branding variables
			customVariables := sassVariablesPattern.FindAllString(customCSS, -1)
			for _, customVariable := range customVariables {
				stringsBuilder.WriteString(fmt.Sprintf("%s \n", customVariable))
			}
		}
	}

	//Save branding sass variables to cache
	sassVariables, err := readSassFiles(log, "scss/variables")
	if err != nil {
		return err
	}
	stringsBuilder.WriteString(sassVariables)
	StylesheetCache.Set("sass", stringsBuilder.String())

	// start processing a section
	sassSection, err := readSassFiles(log, path.Join("scss", section))
	if err != nil {
		return err
	}
	stringsBuilder.WriteString(StylesheetCache.Get("sass"))
	stringsBuilder.WriteString(sassSection)

	// when a user provides sets WEBAPP_SCSS_DIR_PATH environment variable
	if sassDirPath != "" {
		customSass, err := readSassFiles(log, sassDirPath)
		if err != nil {
			return err
		}
		stringsBuilder.WriteString(customSass)
	}

	if customCSS != "" && isCustomCssSyntaxValid {
		//Custom CSS editor selector block
		selectorBlock := sassVariablesPattern.ReplaceAllString(customCSS, "")
		stringsBuilder.WriteString(selectorBlock)
	}

	transpiledCss, err := transpileSass(transpiler, stringsBuilder.String())
	if err != nil {
		log.Error("sass compilation failure", zap.Error(err))
	}

	// in case of sass error in custom css sass compilation,
	// use compiled css from branding sass and append custom css value to it
	if !isCustomCssSyntaxValid {
		stringsBuilder.Reset()
		stringsBuilder.WriteString(transpiledCss)
		stringsBuilder.WriteString(customCSS)

		// append un-compiled custom css content to transpiled css
		transpiledCss = stringsBuilder.String()
	}

	//save the transpiled css to stylesheet cache
	sectionKey := fmt.Sprintf("%s-%s", section, themeID)
	StylesheetCache.Set(sectionKey, transpiledCss)

	return nil
}

// transpileSass computes sass to css by the transpiler
func transpileSass(transpiler *godartsass.Transpiler, sass string) (string, error) {
	args := godartsass.Args{
		Source: sass,
	}
	execute, err := transpiler.Execute(args)
	if err != nil {
		return "", err
	}

	return execute.CSS, nil
}

func DartSassTranspiler(log *zap.Logger) *godartsass.Transpiler {
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
