package dal

import (
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/locale"
)

func errModelHigherSensitivity(model, connection string) error {
	return errors.New(
		errors.KindSensitiveData,

		"model sensitivity surpasses connection sensitivity",

		errors.Meta("type", "invalid sensitivity"),

		// Translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "internal"),
		errors.Meta(locale.ErrorMetaKey{}, "dal.sensitivity.model-exceeds-connection"),
		errors.Meta("model", model),
		errors.Meta("connection", connection),

		errors.StackSkip(1),
		errors.StackTrimAtFn("http.HandlerFunc.ServeHTTP"),
	)
}

func errAttributeHigherSensitivity(model, attribute string) error {
	return errors.New(
		errors.KindSensitiveData,

		"attribute sensitivity surpasses model sensitivity",

		errors.Meta("type", "invalid sensitivity"),

		// Translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "internal"),
		errors.Meta(locale.ErrorMetaKey{}, "dal.sensitivity.attribute-exceeds-model"),
		errors.Meta("model", model),
		errors.Meta("attribute", attribute),

		errors.StackSkip(1),
		errors.StackTrimAtFn("http.HandlerFunc.ServeHTTP"),
	)
}
