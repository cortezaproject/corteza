package provision

import (
	"context"
	"fmt"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	"go.uber.org/zap"
)

// Pre 2021.3 versions had email templates stored in settings
// from 2021.3 onwards we have dedicated subsystem for managing templates
//
// This migration moves email templates from settings templates.
func migrateEmailTemplates(ctx context.Context, log *zap.Logger, s store.Storer) error {
	var (
		// setting name => template handle
		m = map[string]*types.Template{
			//"general.mail.logo",
			"general.mail.header.en": {
				Type:    "text/html",
				Handle:  "email_general_header",
				Partial: true,
				Meta: types.TemplateMeta{
					Short:       "General template header",
					Description: "General template header to use with system email notifications",
				},
			},
			"general.mail.footer.en": {
				Type:    "text/html",
				Handle:  "email_general_footer",
				Partial: true,
				Meta: types.TemplateMeta{
					Short:       "General template footer",
					Description: "General template footer to use with system email notifications",
				},
			},
			"auth.mail.email-confirmation.subject.en": {
				Type:   "text/plain",
				Handle: "auth_email_confirmation_subject",
				Meta:   types.TemplateMeta{Short: "Password reset subject"},
			},
			"auth.mail.email-confirmation.body.en": {
				Type:   "text/html",
				Handle: "auth_email_confirmation_body",
				Meta:   types.TemplateMeta{Short: "Password reset content"},
			},
			"auth.mail.password-reset.subject.en": {
				Type:   "text/plain",
				Handle: "auth_email_password_reset_subject",
				Meta:   types.TemplateMeta{Short: "Email confirmation subject"},
			},
			"auth.mail.password-reset.body.en": {
				Type:   "text/html",
				Handle: "auth_email_password_reset_body",
				Meta:   types.TemplateMeta{Short: "Email confirmation content"},
			},
		}
	)

	return store.Tx(ctx, s, func(ctx context.Context, s store.Storer) error {
		for name, tmpl := range m {
			sval, err := store.LookupSettingByNameOwnedBy(ctx, s, name, 0)
			if errors.IsNotFound(err) {
				// setting not found, that's ok, see the next one
				continue
			} else if err != nil {
				return fmt.Errorf("failed to lookup for setting by bame: %w", err)
			}

			_, err = store.LookupTemplateByHandle(ctx, s, tmpl.Handle)
			if err != nil && !errors.IsNotFound(err) {
				// any error but not-found is fatal.
				return fmt.Errorf("failed to lookup for template by handle: %w", err)
			} else if err == nil {
				// template exists
				continue
			}

			sval.String()
			tmpl.ID = id.Next()
			tmpl.CreatedAt = time.Now()
			tmpl.Template = sval.String()

			err = store.CreateTemplate(ctx, s, tmpl)
			if err != nil {
				return fmt.Errorf("failed to store migrated template: %w", err)
			}

			log.Debug("migrated template from settings",
				zap.String("setting", name),
				zap.String("handle", tmpl.Handle),
			)
		}

		// Go over all settings again and remove them
		for name := range m {
			if err := store.DeleteSettingByNameOwnedBy(ctx, s, name, 0); err != nil {
				return err
			}
		}

		return nil
	})
}
