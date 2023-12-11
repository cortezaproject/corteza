package provision

import (
	"context"
	"encoding/json"
	"github.com/cortezaproject/corteza/server/pkg/sass"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
	"go.uber.org/zap"
	"strconv"
)

// updateWebappTheme is a function that provisions webapp themes.
// It migrates the old custom css and branding sass settings to the new webapp themes setting.
func updateWebappTheme(ctx context.Context, log *zap.Logger, s store.Storer) (err error) {
	vv, _, err := store.SearchSettingValues(ctx, s, types.SettingsFilter{})
	if err != nil {
		return err
	}

	oldCustomCSS := vv.FindByName("ui.custom-css")
	oldBranding := vv.FindByName("ui.studio.branding-sass")

	provisionTheme := func(name string, oldValue *types.SettingValue, themeIDs ...string) (err error) {
		oldValueStr, err := strconv.Unquote(oldValue.Value.String())
		if err != nil {
			return err
		}

		var themes []types.Theme
		for _, themeID := range themeIDs {
			if len(themeIDs) == 2 {
				themes = append(themes, types.Theme{
					ID:     themeID,
					Values: oldValueStr,
				})
				continue
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

		// delete old custom css and branding sass settings from the database
		err = store.DeleteSettingValue(ctx, s, oldValue)
		if err != nil {
			return err
		}

		return nil
	}

	// provision custom CSS
	if !oldCustomCSS.IsNull() {
		err = provisionTheme("ui.studio.custom-css", oldCustomCSS, sass.GeneralTheme, sass.LightTheme, sass.DarkTheme)
		if err != nil {
			return err
		}
	}

	// provision branding sass
	if !oldBranding.IsNull() {
		err = provisionTheme("ui.studio.themes", oldBranding, sass.LightTheme, sass.DarkTheme)
	}

	return nil
}
