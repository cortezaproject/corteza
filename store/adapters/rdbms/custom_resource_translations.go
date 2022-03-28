package rdbms

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/locale"
	"golang.org/x/text/language"
)

func (s Store) TransformResource(ctx context.Context, lang language.Tag) (map[string]map[string]*locale.ResourceTranslation, error) {
	//TODO implement me
	panic("implement me")
}
