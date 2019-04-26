// +build integration,external

package service

import (
	"context"
	"strings"
	"testing"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/internal/config"
	"github.com/crusttech/crust/internal/http"
	"github.com/crusttech/crust/internal/test"
	"github.com/crusttech/crust/messaging/internal/repository"
	"github.com/crusttech/crust/messaging/types"
	systemService "github.com/crusttech/crust/system/service"
	systemTypes "github.com/crusttech/crust/system/types"
)

func TestOutgoingWebhook(t *testing.T) {
	var user = &systemTypes.User{ID: 1}
	var channel = &types.Channel{ID: 1}

	ctx := context.WithValue(context.Background(), "testing", true)
	ctx = auth.SetIdentityToContext(ctx, user)

	// set up user
	{
		svc := systemService.User(ctx)
		_, err := svc.Create(user)
		test.Assert(t, err == nil, "Error when creating user: %+v", err)
	}

	// set up channel
	{
		rpo := repository.Channel(ctx, repository.DB(ctx))
		_, err := rpo.Create(channel)
		test.Assert(t, err == nil, "Error when creating channel: %+v", err)
	}

	client, err := http.New(&config.HTTPClient{
		Timeout: 10,
	})
	test.Assert(t, err == nil, "Error creating HTTP client: %+v", err)

	svc := Webhook(ctx, client)

	{
		/* create outgoing webhook */
		webhook, err := svc.Create(types.OutgoingWebhook, channel.ID, types.WebhookRequest{
			Username:        "test-webhook",
			OutgoingTrigger: "fortune",
			OutgoingURL:     "https://api.scene-si.org/fortune.php",
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
			test.Assert(t, strings.Contains(message.Message, "BOFH"), "Unexpected webhook output: %s", message.Message)
		}

		// update webhook
		wh, err := svc.Update(webhook.ID, types.OutgoingWebhook, channel.ID, types.WebhookRequest{
			Username:        "test-webhook-json",
			OutgoingTrigger: "fortune-json",
			OutgoingURL:     "https://api.scene-si.org/fortune.php?username=test",
		})
		test.Assert(t, err == nil, "Error when updating webhook: %+v", err)

		/* trigger outgoing webhook */
		{
			message, err := svc.Do(wh, "")
			test.Assert(t, err == nil, "Error when triggering webhook: %+v", err)
			test.Assert(t, message.Meta.Username == "test", "Expected message.meta.username = 'test', got: '%s'", message.Meta.Username)
			test.Assert(t, strings.Contains(message.Message, "BOFH"), "Unexpected webhook output: %s", message.Message)
		}
	}
}
