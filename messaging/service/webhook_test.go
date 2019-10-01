// +build integration

package service

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cortezaproject/corteza-server/internal/http"
	"github.com/cortezaproject/corteza-server/messaging/repository"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/auth"
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
		require.True(t, err == nil, "Error when creating channel: %+v", err)
	}

	client, err := http.New(&http.Config{
		Timeout: 10,
	})
	require.True(t, err == nil, "Error creating HTTP client: %+v", err)

	svc := Webhook(ctx, client)

	{
		/* create outgoing webhook */
		webhook, err := svc.Create(types.OutgoingWebhook, channel.ID, types.WebhookRequest{
			Username:        "test-webhook",
			UserID:          1337,
			OutgoingTrigger: "fortune",
			OutgoingURL:     server.URL,
		})
		require.True(t, err == nil, "Error when creating webhook: %+v", err)

		/* find outgoing webhook */
		webhooks, err := svc.Find(&types.WebhookFilter{
			OutgoingTrigger: webhook.OutgoingTrigger,
		})
		require.True(t, err == nil, "Error when finding webhook: %+v", err)
		require.True(t, len(webhooks) == 1, "Expected to find 1 webhook, got %d", len(webhooks))

		/* trigger outgoing webhook */
		{
			message, err := svc.Do(webhooks[0], "")
			require.True(t, err == nil, "Error when triggering webhook: %+v", err)
			require.True(t, strings.Contains(message.Message, "Louis Pasteur"), "Unexpected webhook output: %s", message.Message)
		}

		// update webhook
		wh, err := svc.Update(webhook.ID, types.OutgoingWebhook, channel.ID, types.WebhookRequest{
			Username:        "test-webhook-json",
			UserID:          1337,
			OutgoingTrigger: "fortune-json",
			OutgoingURL:     server.URL + "?username=test",
		})
		require.True(t, err == nil, "Error when updating webhook: %+v", err)

		/* trigger outgoing webhook */
		{
			message, err := svc.Do(wh, "")
			require.True(t, err == nil, "Error when triggering webhook: %+v", err)
			require.True(t, message.Meta.Username == "test", "Expected message.meta.username = 'test', got: '%s'", message.Meta.Username)
			require.True(t, strings.Contains(message.Message, "Louis Pasteur"), "Unexpected webhook output: %s", message.Message)
		}
	}
}
