package auth

import (
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/locale"
)

func ErrUnauthorized() error {
	return errors.New(
		errors.KindUnauthorized,

		"unauthorized",

		errors.Meta("type", "unauthorized"),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "internal"),
		errors.Meta(locale.ErrorMetaKey{}, "auth.errors.unauthorized"),

		errors.StackSkip(1),
		errors.StackTrimAtFn("http.HandlerFunc.ServeHTTP"),
	)
}

func ErrUnauthorizedScope() error {
	return errors.New(
		errors.KindUnauthorized,

		"unauthorized scope",

		errors.Meta("type", "unauthorizedScope"),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "internal"),
		errors.Meta(locale.ErrorMetaKey{}, "auth.errors.unauthorizedScope"),

		errors.StackSkip(1),
		errors.StackTrimAtFn("http.HandlerFunc.ServeHTTP"),
	)
}

func ErrMalformedToken(details string) error {
	return errors.New(
		errors.KindUnauthorized,

		"malformed token: "+details,

		errors.Meta("type", "malformedToken"),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "internal"),
		errors.Meta(locale.ErrorMetaKey{}, "auth.errors.malformedToken"),

		errors.StackSkip(1),
		errors.StackTrimAtFn("http.HandlerFunc.ServeHTTP"),
	)
}
