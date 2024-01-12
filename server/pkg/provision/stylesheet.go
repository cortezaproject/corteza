package provision

import (
	"context"
	"encoding/json"
	"strconv"

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

	//provision new studio themes setting
	newThemes := processNewTheme()
	if oldBranding.IsNull() && studioThemes.IsNull() {
		err = provisionTheme(ctx, s, "ui.studio.themes", newThemes, log)
		if err != nil {
			return err
		}
	}

	// provision old branding sass setting to studio themes setting
	if !oldBranding.IsNull() {
		themes, err := processOldTheme(oldBranding, sass.LightTheme, sass.DarkTheme)
		if err != nil {
			return err
		}

		//append dark mode from new themes
		themes = append(themes, newThemes[1])

		err = provisionTheme(ctx, s, "ui.studio.themes", themes, log)
		if err != nil {
			return err
		}

		// delete old custom css and branding sass settings from the database
		err = store.DeleteSettingValue(ctx, s, oldBranding)
		if err != nil {
			return err
		}
	}

	// provision custom CSS
	if !oldCustomCSS.IsNull() {
		themes, err := processOldTheme(oldCustomCSS, sass.GeneralTheme, sass.LightTheme, sass.DarkTheme)
		if err != nil {
			return err
		}

		err = provisionTheme(ctx, s, "ui.studio.custom-css", themes, log)
		if err != nil {
			return err
		}

		// delete old custom css and branding sass settings from the database
		err = store.DeleteSettingValue(ctx, s, oldCustomCSS)
		if err != nil {
			return err
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
	}

	err = store.CreateSettingValue(ctx, s, newThemeSetting)
	if err != nil {
		log.Error("failed to provision webapp themes", zap.Error(err))
		return err
	}

	return nil
}

func processNewTheme() (themes []types.Theme) {
	lightModeValues := `
    {
        "black":"#162425",
        "white":"#FFFFFF",
        "primary":"#0B344E",
        "secondary":"#758D9B",
        "success":"#43AA8B",
        "warning":"#E2A046",
        "danger":"#E54122",
        "light":"#F3F5F7",
        "extra-light":"#E4E9EF",
        "body-bg":"#F3F5F7",
        "sidebar-bg": "#F3F5F7",
        "topbar-bg": "#F3F5F7"
    }`

	darkModeValues := `
    {
        "black":"#FBF7F4",
        "white":"#0B344E",
        "primary":"#FF9661",
        "secondary":"#758D9B",
        "success":"#43AA8B",
        "warning":"#E2A046",
        "danger":"#E54122",
        "light":"#768D9A",
        "extra-light":"#23495F",
        "body-bg":"#092B40",
        "sidebar-bg": "#768D9A",
        "topbar-bg": "#768D9A"
    }`

	themes = []types.Theme{
		{
			ID:     "light",
			Values: lightModeValues,
		},
		{
			ID:     "dark",
			Values: darkModeValues,
		},
	}

	return themes
}

func processOldTheme(oldValue *types.SettingValue, themeIDs ...string) (themes []types.Theme, err error) {
	oldValueStr, err := strconv.Unquote(oldValue.Value.String())
	if err != nil {
		return
	}

	for _, themeID := range themeIDs {
		//append only light mode on studio themes
		if len(themeIDs) == 2 {
			themes = append(themes, types.Theme{
				ID:     themeID,
				Values: oldValueStr,
			})
			break
		}

		if themeID == sass.GeneralTheme {
			themes = append(themes, types.Theme{
				ID:     themeID,
				Values: oldValueStr,
			})
			continue
		}

		themes = append(themes, types.Theme{
			ID:     themeID,
			Values: "",
		})
	}

	return
}
