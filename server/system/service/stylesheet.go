package service

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/logger"
	"github.com/cortezaproject/corteza/server/pkg/sass"
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
func GenerateCSS(ctx context.Context, brandingSass, customCSS, sassDirPath string) (string, error) {
	var (
		log        = logger.Default()
		transpiler = sass.DartSass(log)
	)

	//check if cache has already compiled css value
	if !sass.StylesheetCache.Empty() {
		return sass.StylesheetCache.Get("css"), nil
	}

	// if dart sass is not installed, or when the sass transpiler creation and startup process fails.
	// return contents from default css
	if transpiler == nil {
		updateSassInstalledSetting(ctx, log, false)
		return sass.DefaultCSS(log, customCSS), nil
	}

	updateSassInstalledSetting(ctx, log, true)

	// transpile sass to css
	err := sass.Transpiler(log, brandingSass, customCSS, sassDirPath, transpiler)

	if err != nil {
		return sass.DefaultCSS(log, customCSS), err
	}

	return sass.StylesheetCache.Get("css"), nil
}

func updateSassInstalledSetting(ctx context.Context, log *zap.Logger, installed bool) {
	sv := &types.SettingValue{
		Name: "ui.studio.sass-installed",
	}

	err := sv.SetSetting(installed)
	if err != nil {
		log.Warn("failed to set ui.studio.sass-installed setting", zap.Error(err))
	}

	err = DefaultSettings.Set(ctx, sv)
	if err != nil {
		log.Warn("failed to set ui.studio.sass-installed setting", zap.Error(err))
	}
}
