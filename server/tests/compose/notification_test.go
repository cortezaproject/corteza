package compose

import (
	"encoding/json"
	"github.com/crusttech/human/server/compose/rest/request"
	"github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/tests/helpers"
	sqlxTypes "github.com/jmoiron/sqlx/types"
	"github.com/steinfletcher/apitest"
	"net/http"
	"testing"
)

func (h helper) apiSendEmailNotification(req request.NotificationEmailSend) *apitest.Response {
	payload, err := json.Marshal(req)
	h.noError(err)

	return h.apiInit().
		Post("/notification/email").
		Header("Accept", "application/json").
		JSON(string(payload)).
		Expect(h.t).
		Status(http.StatusOK)
}

func TestEmailNotification(t *testing.T) {
	t.Skip("we need smtp server mock")
	h := newHelper(t)

	h.apiSendEmailNotification(request.NotificationEmailSend{
		To:                []string{"foo+to@test.tld"},
		Cc:                []string{"foo+cc1@test.tld", "foo+cc2@test.tld"},
		Subject:           "Subject!",
		Content:           sqlxTypes.JSONText(`{}`),
		RemoteAttachments: []string{"file1", "file2"},
	}).End()
}

func TestEmailNotificationForbidden(t *testing.T) {
	h := newHelper(t)
	helpers.DenyMe(h, types.ComponentRbacResource(), "email-notifications.send")

	h.apiSendEmailNotification(request.NotificationEmailSend{
		To:      []string{"foo+to@test.tld"},
		Subject: "Subject!",
		Content: sqlxTypes.JSONText(`{}`),
	}).
		Assert(helpers.AssertError("notification.errors.notAllowedToSend")).
		End()
}

// Without SMTP server mock we can only verify that request gets past the access control
func TestEmailNotificationAllowed(t *testing.T) {
	h := newHelper(t)
	helpers.AllowMe(h, types.ComponentRbacResource(), "email-notifications.send")

	h.apiSendEmailNotification(request.NotificationEmailSend{
		Subject: "Subject!",
		Content: sqlxTypes.JSONText(`{}`),
	}).
		Assert(helpers.AssertError("notification.errors.noRecipients")).
		End()
}
