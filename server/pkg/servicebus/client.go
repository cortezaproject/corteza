package servicebus

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/cortezaproject/corteza/server/pkg/messagebus/types"
	"os"
)

type (
	Client interface {
		SendMessage(context.Context, types.QueueMessage) error
		GetMessage(context.Context, int, types.QueueMessage)
	}

	sbClient struct {
		client *azservicebus.Client
		store  azservicebus.Client
	}
)

func NewClient() *sbClient {
	// ex: myservicebus.servicebus.windows.net
	// @fixme: put up this to proper channels
	namespace, ok := os.LookupEnv("AZURE_SERVICEBUS_HOSTNAME")
	if !ok {
		// todo err
		fmt.Println("AZURE_SERVICEBUS_HOSTNAME environment variable not found")
		return nil
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// todo err
		fmt.Println(err)
		return nil
	}

	client, err := azservicebus.NewClient(namespace, cred, nil)
	if err != nil {
		// todo err
		return nil
	}

	return &sbClient{client: client}
}

func (c *sbClient) SendMessage(ctx context.Context, qm types.QueueMessage) (err error) {
	sender, err := c.client.NewSender(qm.Queue, nil)
	if err != nil {
		// todo err
		return
	}
	defer func(sender *azservicebus.Sender, ctx context.Context) {
		err = sender.Close(ctx)
		if err != nil {
			// fixme err
			return
		}
	}(sender, ctx)

	sbMessage := &azservicebus.Message{
		Body: qm.Payload,
	}
	err = sender.SendMessage(ctx, sbMessage, nil)
	if err != nil {
		// todo err
		return
	}

	c.GetMessage(ctx, 1, qm)

	return
}

func (c *sbClient) SendMessageBatch(ctx context.Context, q string, payload [][]byte) (err error) {
	sender, err := c.client.NewSender("myqueue", nil)
	if err != nil {
		// todo err
		return
	}
	defer func(sender *azservicebus.Sender, ctx context.Context) {
		err = sender.Close(ctx)
		if err != nil {
			// fixme err
			return
		}
	}(sender, ctx)

	batch, err := sender.NewMessageBatch(ctx, nil)
	if err != nil {
		// todo err
		return
	}

	for _, message := range payload {
		if err = batch.AddMessage(&azservicebus.Message{Body: message}, nil); err != nil {
			// todo err
			return
		}
	}
	if err = sender.SendMessageBatch(ctx, batch, nil); err != nil {
		// todo err
		return
	}

	return
}

// func (c *sbClient) GetMessage(ctx context.Context, q string, count int) {
func (c *sbClient) GetMessage(ctx context.Context, count int, qm types.QueueMessage) {
	receiver, err := c.client.NewReceiverForQueue(qm.Queue, nil) // Change myqueue to env var
	if err != nil {
		panic(err)
	}
	defer func(receiver *azservicebus.Receiver, ctx context.Context) {
		err := receiver.Close(ctx)
		if err != nil {
			// todo err
			return
		}
	}(receiver, ctx)

	messages, err := receiver.ReceiveMessages(ctx, count, nil)
	if err != nil {
		// todo err
		return
	}

	for _, message := range messages {
		body := message.Body
		// fixme: handle message
		fmt.Printf("%s\n", string(body))

		err = receiver.CompleteMessage(ctx, message, nil)
		if err != nil {
			// todo err
			return
		}
	}
}
