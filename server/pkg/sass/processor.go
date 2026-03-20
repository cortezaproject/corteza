package sass

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bep/godartsass/v2"
	"go.uber.org/zap"
)

const (
	GeneralTheme = "general"
	LightTheme   = "light"
	DarkTheme    = "dark"
)

var (
	StylesheetCache = newStylesheetCache()
)

// Transpile compiles theme SASS variables + custom CSS into final CSS for a given theme
func Transpile(transpiler *godartsass.Transpiler, log *zap.Logger, themeID, themeSASS, customCSS string) (err error) {
	var stringsBuilder strings.Builder

	// add theme-mode variable
	stringsBuilder.WriteString(fmt.Sprintf("$theme-mode: %s;\n", themeID))

	// convert theme JSON variables to SASS variable assignments
	if themeSASS != "" {
		if err = jsonToSass(themeSASS, &stringsBuilder); err != nil {
			log.Error("failed to unmarshal theme sass variables", zap.Error(err))
			return err
		}
	}

	// append custom CSS
	if customCSS != "" {
		stringsBuilder.WriteString(customCSS)
	}

	source := stringsBuilder.String()
	if strings.TrimSpace(source) == "" {
		return nil
	}

	transpiledCSS, err := transpileSass(transpiler, source)
	if err != nil {
		log.Error("sass compilation failure", zap.Error(err))
		// If transpilation fails, store the raw custom CSS as-is
		StylesheetCache.Set(themeID, customCSS)
		return nil
	}

	StylesheetCache.Set(themeID, transpiledCSS)
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
