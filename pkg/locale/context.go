package locale

import (
	"context"

	"golang.org/x/text/language"
)

func SetAcceptLanguageToContext(ctx context.Context, code language.Tag) context.Context {
	return context.WithValue(ctx, language.Tag{}, code)
}

// GetAcceptLanguageFromContext always returns language, either valid or default
func GetAcceptLanguageFromContext(ctx context.Context) language.Tag {
	if tag, ok := ctx.Value(language.Tag{}).(language.Tag); ok {
		return tag
	} else {
		return defaultLanguage
	}
}

type contentLanguageCtx struct{}

func SetContentLanguageToContext(ctx context.Context, code language.Tag) context.Context {
	return context.WithValue(ctx, contentLanguageCtx{}, code)
}

// GetContentLanguageFromContext always returns language, either valid or default
func GetContentLanguageFromContext(ctx context.Context) language.Tag {
	if tag, ok := ctx.Value(contentLanguageCtx{}).(language.Tag); ok {
		return tag
	} else {
		return language.Und
	}
}
