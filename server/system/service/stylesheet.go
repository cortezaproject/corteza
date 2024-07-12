package service

import (
    "fmt"
    "github.com/bep/godartsass/v2"
    "github.com/cespare/xxhash/v2"
    "github.com/cortezaproject/corteza/server/pkg/sass"
    "github.com/cortezaproject/corteza/server/system/types"
    "go.uber.org/zap"
    "strings"
)

type (
    stylesheet struct {
        transpiler *godartsass.Transpiler
        logger     *zap.Logger
    }
)

func Stylesheet(transpiler *godartsass.Transpiler, logger *zap.Logger) *stylesheet {
    return &stylesheet{
        transpiler: transpiler,
        logger:     logger,
    }
}

// GenerateCSS takes care of creating CSS for webapps by reading SASS content from embedded assets,
// combining it with different themeSASS and customCSS themes, and then transpiling it using the dart-sass compiler.
//
// If dart sass isn't installed on the host machine, it will default to css content from the minified-custom.css which is
// generated from [Boostrap, bootstrap-vue and custom variables sass content].
// If dart isn't installed on the host machine, customCustom css will continue to function, but without sass support.
//
// In case of an error, it will return default css and log out the error
func (svc *stylesheet) GenerateCSS(settings *types.AppSettings, sassDirPath string, log *zap.Logger) (err error) {
	var (
		studio       = settings.UI.Studio
		customCSSMap = make(map[string]string)
	)

	for _, customCSS := range studio.CustomCSS {
		customCSSMap[customCSS.ID] = customCSS.Values
	}

	sass.DefaultCSS(log, customCSSMap[sass.GeneralTheme])

	if studio.Themes == nil && studio.CustomCSS == nil {
		return
	}

	// if dart sass is not installed, or when the sass transpiler creation and startup process fails.
	if !studio.SassInstalled {
		return
	}

	// transpile sass to css for each theme
	for _, theme := range studio.Themes {
		if studio.CustomCSS == nil {
            err := sass.Transpile(svc.transpiler, log, theme.ID, theme.Values, "", sassDirPath)
			if err != nil {
				continue
			}
		}

		customCSS := processCustomCSS(theme.ID, customCSSMap)
		// transpile sass to css
        err := sass.Transpile(svc.transpiler, log, theme.ID, theme.Values, customCSS, sassDirPath)
		if err != nil {
			continue
		}
	}

	return
}

func (svc *stylesheet) SassInstalled() bool {
    return svc.transpiler != nil
}
// processCustomCSS, processes CustomCSS input and gives priority to theme specific customCSS
func processCustomCSS(themeID string, customCSSMap map[string]string) (customCSS string) {
	var stringsBuilder strings.Builder

	// add theme mode on customCSS
	if themeID == sass.DarkTheme {
		stringsBuilder.WriteString(fmt.Sprintf("\n[data-color-mode=\"%s\"] {\n", themeID))
	}

	stringsBuilder.WriteString(customCSSMap[sass.GeneralTheme])
	stringsBuilder.WriteString("\n")
	stringsBuilder.WriteString(customCSSMap[themeID])

	if themeID == sass.DarkTheme {
		stringsBuilder.WriteString("}\n")
	}

	return stringsBuilder.String()
}

// updateCSS, updates theme css when ui.studio.themes or ui.studio.custom-css settings are updated
func (svc *stylesheet) updateCSS(current, old, compStyles *types.SettingValue, name, sassDirPath string, log *zap.Logger) {
	complimentaryStylesMap := themeMap(compStyles)
	oldThemesMap := themeMap(old)
	currentThemesMap := themeMap(current)

	transpileSASS := func(themeID, themeSASS string, themeCustomCSS map[string]string) {
		customCSS := processCustomCSS(themeID, themeCustomCSS)
        err := sass.Transpile(svc.transpiler, log, themeID, themeSASS, customCSS, sassDirPath)
		if err != nil {
			log.Error("failed to transpile sass to css", zap.Error(err))
		}
	}

	for key := range currentThemesMap {
		if xxhash.Sum64String(oldThemesMap[key]) == xxhash.Sum64String(currentThemesMap[key]) {
			continue
		}

		if name == "ui.studio.themes" {
			transpileSASS(key, currentThemesMap[key], complimentaryStylesMap)
			continue
		}

		if key == sass.GeneralTheme {
			if complimentaryStylesMap == nil {
				transpileSASS(sass.LightTheme, complimentaryStylesMap[sass.LightTheme], currentThemesMap)
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
	var (
		stringsBuilder strings.Builder
		rootLight      = sass.StylesheetCache.Get(fmt.Sprintf("%s-%s", sass.SectionRoot, sass.LightTheme))
	)

	if rootLight == "" {
		return sass.StylesheetCache.Get("default-theme")
	}

	// root css section
	stringsBuilder.WriteString(rootLight)
	stringsBuilder.WriteString("\n")
	stringsBuilder.WriteString(sass.StylesheetCache.Get(fmt.Sprintf("%s-%s", sass.SectionRoot, sass.DarkTheme)))
	stringsBuilder.WriteString("\n")

	//theme css section
	stringsBuilder.WriteString(sass.StylesheetCache.Get(fmt.Sprintf("%s-%s", sass.SectionTheme, sass.DarkTheme)))
	stringsBuilder.WriteString("\n")

	// body css section
	stringsBuilder.WriteString(sass.StylesheetCache.Get(fmt.Sprintf("%s-%s", sass.SectionMain, sass.LightTheme)))
	stringsBuilder.WriteString("\n")

	return stringsBuilder.String()
}
