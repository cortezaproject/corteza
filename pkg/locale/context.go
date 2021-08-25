package locale

import (
	"context"

	"golang.org/x/text/language"
)

func SetLanguageToContext(ctx context.Context, code language.Tag) context.Context {
	return context.WithValue(ctx, language.Tag{}, code)
}

// GetLanguageFromContext always returns language, either valid or default
func GetLanguageFromContext(ctx context.Context) language.Tag {
	if tag, ok := ctx.Value(language.Tag{}).(language.Tag); ok {
		return tag
	} else {
		return language.Tag{}
	}
}
