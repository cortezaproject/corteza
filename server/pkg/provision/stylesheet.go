package provision

import (
    "context"
    "encoding/json"
    "strconv"
    "time"

    "github.com/cortezaproject/corteza/server/pkg/sass"
    "github.com/cortezaproject/corteza/server/store"
    "github.com/cortezaproject/corteza/server/system/types"
    "go.uber.org/zap"
)

// updateWebappTheme is a function that provisions new webapp themes,
// and migrates the old custom css and branding sass settings to new webapp themes setting.
func updateWebappTheme(ctx context.Context, log *zap.Logger, s store.Storer) (err error) {
	vv, _, err := store.SearchSettingValues(ctx, s, types.SettingsFilter{})
	if err != nil {
		return err
	}

	oldCustomCSS := vv.FindByName("ui.custom-css")
	oldBranding := vv.FindByName("ui.studio.branding-sass")
	studioThemes := vv.FindByName("ui.studio.themes")
	customCSSThemes := vv.FindByName("ui.studio.custom-css")

	//get branding themes setting value
	brandingTheme, err := processBrandingTheme(oldBranding)
	if err != nil {
		return err
	}

	//get custom css themes setting value
	customCSSTheme, err := processCustomCSSTheme(oldCustomCSS)
	if err != nil {
		return err
	}

	// provision new themes
	if studioThemes.IsNull() {
		// provision branding themes setting
		err = provisionTheme(ctx, s, "ui.studio.themes", brandingTheme, log)
		if err != nil {
			return err
		}

		if !oldBranding.IsNull() {
			// delete old branding sass settings from the database
			err = store.DeleteSettingValue(ctx, s, oldBranding)
			if err != nil {
				return err
			}
		}
	}

	if customCSSThemes.IsNull() {
		// provision custom CSS themes setting
		err = provisionTheme(ctx, s, "ui.studio.custom-css", customCSSTheme, log)
		if err != nil {
			return err
		}

		if !oldCustomCSS.IsNull() {
			// delete old custom css settings from the database
			err = store.DeleteSettingValue(ctx, s, oldCustomCSS)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func provisionTheme(ctx context.Context, s store.Storer, name string, themes []types.Theme, log *zap.Logger) (err error) {
	value, err := json.Marshal(themes)
	if err != nil {
		return err
	}

	newThemeSetting := &types.SettingValue{
		Name:  name,
		Value: value,
        UpdatedAt: time.Now(),
	}

	err = store.CreateSettingValue(ctx, s, newThemeSetting)
	if err != nil {
		log.Error("failed to provision webapp themes", zap.Error(err))
		return err
	}

	return nil
}

func processBrandingTheme(oldBranding *types.SettingValue) (themes []types.Theme, err error) {
	var brandingMap map[string]string

	lightModeMap := map[string]string{
		"black":       "#162425",
		"white":       "#FFFFFF",
		"primary":     "#0B344E",
		"secondary":   "#758D9B",
		"success":     "#43AA8B",
		"warning":     "#E2A046",
		"danger":      "#E54122",
		"light":       "#F3F5F7",
		"extra-light": "#E4E9EF",
		"body-bg":     "#F3F5F7",
		"sidebar-bg":  "#FFFFFF",
		"topbar-bg":   "#F3F5F7",
	}

	// process old branding sass settings and match them with the new branding themes setting
	if !oldBranding.IsNull() {
		oldBrandingString, err := strconv.Unquote(oldBranding.Value.String())
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal([]byte(oldBrandingString), &brandingMap); err != nil {
			return nil, err
		}

		for key, bmValue := range brandingMap {
			if key == "light" || key == "extra-light" {
				continue
			}

			if _, ok := lightModeMap[key]; ok {
				lightModeMap[key] = bmValue
			}
		}
	}

	darkModeValues := `
    {
        "black":"#FBF7F4",
        "white":"#0B344E",
        "primary":"#FF9661",
        "secondary":"#758D9B",
        "success":"#43AA8B",
        "warning":"#E2A046",
        "danger":"#E54122",
        "light":"#23495F",
        "extra-light":"#3E5A6F",
        "body-bg":"#092B40",
        "sidebar-bg": "#0B344E",
        "topbar-bg": "#092B40"
    }`

	lightModeValues, _ := json.Marshal(lightModeMap)

	themes = []types.Theme{
		{
			ID:     sass.LightTheme,
			Values: string(lightModeValues),
		},
		{
			ID:     sass.DarkTheme,
			Values: darkModeValues,
		},
	}

	return themes, nil
}

func processCustomCSSTheme(oldValue *types.SettingValue) (themes []types.Theme, err error) {
	var generalCSS string

	if oldValue.IsNull() {
		generalCSS = ""
	} else {
		generalCSS, err = strconv.Unquote(oldValue.Value.String())
		if err != nil {
			return nil, err
		}
	}

	themes = []types.Theme{
		{
			ID:     sass.GeneralTheme,
			Values: generalCSS,
		},
		{
			ID:     sass.LightTheme,
			Values: "",
		},
		{
			ID:     sass.DarkTheme,
			Values: "",
		},
	}

	return themes, nil
}
