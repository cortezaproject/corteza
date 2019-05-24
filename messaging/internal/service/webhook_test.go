// +build integration

package service

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cortezaproject/corteza-server/internal/auth"
	"github.com/cortezaproject/corteza-server/internal/http"
	"github.com/cortezaproject/corteza-server/internal/test"
	"github.com/cortezaproject/corteza-server/messaging/internal/repository"
	"github.com/cortezaproject/corteza-server/messaging/types"
)

func TestOutgoingWebhook(t *testing.T) {
	handler := &Fortune{}
	server := httptest.NewServer(handler)
	defer server.Close()

	var channel = &types.Channel{ID: 1}

	ctx := context.WithValue(context.Background(), "testing", true)
	ctx = auth.SetIdentityToContext(ctx, auth.NewIdentity(1337))

	// set up channel
	{
		rpo := repository.Channel(ctx, repository.DB(ctx))
		_, err := rpo.Create(channel)
		test.Assert(t, err == nil, "Error when creating channel: %+v", err)
	}

	client, err := http.New(&http.Config{
		Timeout: 10,
	})
	test.Assert(t, err == nil, "Error creating HTTP client: %+v", err)

	svc := Webhook(ctx, client)

	{
		/* create outgoing webhook */
		webhook, err := svc.Create(types.OutgoingWebhook, channel.ID, types.WebhookRequest{
			Username:        "test-webhook",
			UserID:          1337,
			OutgoingTrigger: "fortune",
			OutgoingURL:     server.URL,
		})
		test.Assert(t, err == nil, "Error when creating webhook: %+v", err)

		/* find outgoing webhook */
		webhooks, err := svc.Find(&types.WebhookFilter{
			OutgoingTrigger: webhook.OutgoingTrigger,
		})
		test.Assert(t, err == nil, "Error when finding webhook: %+v", err)
		test.Assert(t, len(webhooks) == 1, "Expected to find 1 webhook, got %d", len(webhooks))

		/* trigger outgoing webhook */
		{
			message, err := svc.Do(webhooks[0], "")
			test.Assert(t, err == nil, "Error when triggering webhook: %+v", err)
			test.Assert(t, strings.Contains(message.Message, "Louis Pasteur"), "Unexpected webhook output: %s", message.Message)
		}

		// update webhook
		wh, err := svc.Update(webhook.ID, types.OutgoingWebhook, channel.ID, types.WebhookRequest{
			Username:        "test-webhook-json",
			UserID:          1337,
			OutgoingTrigger: "fortune-json",
			OutgoingURL:     server.URL + "?username=test",
		})
		test.Assert(t, err == nil, "Error when updating webhook: %+v", err)

		/* trigger outgoing webhook */
		{
			message, err := svc.Do(wh, "")
			test.Assert(t, err == nil, "Error when triggering webhook: %+v", err)
			test.Assert(t, message.Meta.Username == "test", "Expected message.meta.username = 'test', got: '%s'", message.Meta.Username)
			test.Assert(t, strings.Contains(message.Message, "Louis Pasteur"), "Unexpected webhook output: %s", message.Message)
		}
	}
}
