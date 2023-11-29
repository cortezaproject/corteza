package service

import (
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/logger"
	"github.com/cortezaproject/corteza/server/pkg/sass"
	"github.com/cortezaproject/corteza/server/pkg/str"
	"github.com/cortezaproject/corteza/server/system/types"
	"go.uber.org/zap"
)

// GenerateCSS takes care of creating CSS for webapps by reading SASS content from embedded assets,
// combining it with brandingSass and customCSS, and then transpiling it using the dart-sass compiler.
//
// If dart sass isn't installed on the host machine, it will default to css content from the minified-custom.css which is
// generated from [Boostrap, bootstrap-vue and custom variables sass content].
// If dart isn't installed on the host machine, customCustom css will continue to function, but without sass support.
//
// In case of an error, it will return default css and log out the error
func GenerateCSS(settings *types.AppSettings, sassDirPath string) (err error) {
	var (
		log          = logger.Default()
		studio       = settings.UI.Studio
		customCSSMap = make(map[string]string)
	)

	for _, customCSS := range studio.CustomCSS {
		customCSSMap[customCSS.ID] = customCSS.Values
	}

	sass.DefaultCSS(log, customCSSMap["general"])

	if studio.Themes == nil && studio.CustomCSS == nil {
		return
	}

	// if dart sass is not installed, or when the sass transpiler creation and startup process fails.
	if !studio.SassInstalled {
		return
	}

	transpiler := sass.DartSassTranspiler(log)

	// transpile sass to css for each theme
	if studio.Themes != nil {
		for _, theme := range studio.Themes {
			if studio.CustomCSS == nil {
				err := sass.Transpile(transpiler, log, theme.ID, theme.Values, "", sassDirPath)
				if err != nil {
					continue
				}
			}

			customCSS := processCustomCSS(theme.ID, customCSSMap)
			// transpile sass to css
			err := sass.Transpile(transpiler, log, theme.ID, theme.Values, customCSS, sassDirPath)
			if err != nil {
				continue
			}
		}
	} else {
		if studio.CustomCSS != nil {
			for key := range customCSSMap {
				if key == "general" {
					continue
				}

				customCSS := processCustomCSS(key, customCSSMap)
				// transpile sass to css
				err := sass.Transpile(transpiler, log, key, "", customCSS, sassDirPath)
				if err != nil {
					continue
				}
			}
		}
	}

	return
}

// processCustomCSS, processes CustomCSS input and gives priority to theme specific customCSS
func processCustomCSS(themeID string, customCSSMap map[string]string) (customCSS string) {
	var stringsBuilder strings.Builder

	stringsBuilder.WriteString(customCSSMap["general"])
	stringsBuilder.WriteString("\n")
	stringsBuilder.WriteString(customCSSMap[themeID])

	return stringsBuilder.String()
}

// updateCSS, updates theme css when ui.studio.themes or ui.studio.custom-css settings are updated
func updateCSS(current, old, compStyles *types.SettingValue, name, sassDirPath string, log *zap.Logger) {
	transpiler := sass.DartSassTranspiler(log)

	complimentaryStylesMap := themeMap(compStyles)
	oldThemesMap := themeMap(old)
	currentThemesMap := themeMap(current)

	transpileSASS := func(themeID, themeSASS string, themeCustomCSS map[string]string) {
		customCSS := processCustomCSS(themeID, themeCustomCSS)
		err := sass.Transpile(transpiler, log, themeID, themeSASS, customCSS, sassDirPath)
		if err != nil {
			log.Error("failed to transpile sass to css", zap.Error(err))
		}
	}

	for key := range currentThemesMap {
		if str.HashStringSHA256(oldThemesMap[key]) == str.HashStringSHA256(currentThemesMap[key]) {
			continue
		}

		if name == "ui.studio.themes" {
			transpileSASS(key, currentThemesMap[key], complimentaryStylesMap)
			continue
		}

		if key == "general" {
			if complimentaryStylesMap == nil {
				themeID := "light"
				transpileSASS(themeID, complimentaryStylesMap[themeID], currentThemesMap)
				continue
			}

			for themeID := range complimentaryStylesMap {
				transpileSASS(themeID, complimentaryStylesMap[themeID], currentThemesMap)
			}
			continue
		}

		transpileSASS(key, complimentaryStylesMap[key], currentThemesMap)
	}

}

func themeMap(settingsValue *types.SettingValue) (themeMap map[string]string) {
	var themes []types.Theme

	if settingsValue == nil {
		return themeMap
	}

	_ = settingsValue.Value.Unmarshal(&themes)

	themeMap = make(map[string]string)
	for _, theme := range themes {
		themeMap[theme.ID] = theme.Values
	}

	return themeMap
}

func FetchCSS() string {
	var stringsBuilder strings.Builder

	if sass.StylesheetCache.Get("root-light") == "" {
		return sass.StylesheetCache.Get("default-theme")
	}

	// root css section
	stringsBuilder.WriteString(sass.StylesheetCache.Get("root-light"))
	stringsBuilder.WriteString("\n")
	stringsBuilder.WriteString(sass.StylesheetCache.Get("root-dark"))

	stringsBuilder.WriteString("\n")

	//theme css section
	stringsBuilder.WriteString(sass.StylesheetCache.Get("theme-dark"))
	stringsBuilder.WriteString("\n")

	// body css section
	stringsBuilder.WriteString(sass.StylesheetCache.Get("main-light"))
	stringsBuilder.WriteString("\n")

	return stringsBuilder.String()
}
