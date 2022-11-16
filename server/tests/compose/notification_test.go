package compose

import (
	"encoding/json"
	"github.com/cortezaproject/corteza/server/compose/rest/request"
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
