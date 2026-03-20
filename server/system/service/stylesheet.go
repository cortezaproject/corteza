package service

import (
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

// GenerateCSS takes care of creating CSS for webapps by reading theme variables
// and custom CSS from settings, then transpiling with dart-sass.
//
// If dart sass isn't installed, custom CSS is stored as-is without SCSS support.
func (svc *stylesheet) GenerateCSS(settings *types.AppSettings, log *zap.Logger) (err error) {
	var (
		studio       = settings.UI.Studio
		customCSSMap = make(map[string]string)
	)

	for _, customCSS := range studio.CustomCSS {
		customCSSMap[customCSS.ID] = customCSS.Values
	}

	if studio.Themes == nil && studio.CustomCSS == nil {
		return
	}

	// if dart sass is not installed, store custom CSS as-is
	if svc.transpiler == nil {
		generalCSS := customCSSMap[sass.GeneralTheme]
		if generalCSS != "" {
			sass.StylesheetCache.Set("custom", generalCSS)
		}
		return
	}

	// transpile sass to css for each theme
	for _, theme := range studio.Themes {
		customCSS := processCustomCSS(theme.ID, customCSSMap)
		err := sass.Transpile(svc.transpiler, log, theme.ID, theme.Values, customCSS)
		if err != nil {
			continue
		}
	}

	return
}




// processCustomCSS processes CustomCSS input and gives priority to theme specific customCSS
// Light theme: general CSS + light-specific CSS (no wrapper)
// Dark theme: only dark-specific CSS wrapped in .dark { }
func processCustomCSS(themeID string, customCSSMap map[string]string) (customCSS string) {
	var stringsBuilder strings.Builder

	if themeID == sass.DarkTheme {
		// Dark: only dark-specific CSS, nested in .dark { }
		darkCSS := customCSSMap[themeID]
		if darkCSS != "" {
			stringsBuilder.WriteString(".dark {\n")
			stringsBuilder.WriteString(darkCSS)
			stringsBuilder.WriteString("\n}\n")
		}
	} else {
		// Light: general CSS first, then light-specific CSS
		stringsBuilder.WriteString(customCSSMap[sass.GeneralTheme])
		stringsBuilder.WriteString("\n")
		stringsBuilder.WriteString(customCSSMap[themeID])
	}

	return stringsBuilder.String()
}

// updateCSS updates theme css when ui.studio.themes or ui.studio.custom-css settings are updated
func (svc *stylesheet) updateCSS(current, old, compStyles *types.SettingValue, name string, log *zap.Logger) {
	complimentaryStylesMap := themeMap(compStyles)
	oldThemesMap := themeMap(old)
	currentThemesMap := themeMap(current)

	transpileSASS := func(themeID, themeSASS string, themeCustomCSS map[string]string) {
		customCSS := processCustomCSS(themeID, themeCustomCSS)
		err := sass.Transpile(svc.transpiler, log, themeID, themeSASS, customCSS)
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

// FetchCSS returns the compiled CSS for all themes
func FetchCSS() string {
	var stringsBuilder strings.Builder

	// Check if we have transpiled theme CSS
	lightCSS := sass.StylesheetCache.Get(sass.LightTheme)
	darkCSS := sass.StylesheetCache.Get(sass.DarkTheme)

	if lightCSS == "" && darkCSS == "" {
		// No transpiled CSS, return raw custom CSS if available
		return sass.StylesheetCache.Get("custom")
	}

	if lightCSS != "" {
		stringsBuilder.WriteString(lightCSS)
		stringsBuilder.WriteString("\n")
	}

	if darkCSS != "" {
		stringsBuilder.WriteString(darkCSS)
		stringsBuilder.WriteString("\n")
	}

	return stringsBuilder.String()
}
