package servicebus

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	azadmin "github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/admin"
	"go.uber.org/zap"
)

type (
	sbClient struct {
		logger *zap.Logger

		azClient *azadmin.Client
		client   *azservicebus.Client
		// store    azservicebus.Client
	}
)

func newClient(logger *zap.Logger, connStr string) (cc *sbClient, err error) {
	cc = &sbClient{
		logger: logger,
	}

	// to create the queues
	cc.azClient, err = azadmin.NewClientFromConnectionString(connStr, nil)
	if err != nil {
		cc.logger.Error("could not create service bus admin client", zap.Error(err))
		return
	}

	// to send and receive the messages
	cc.client, err = azservicebus.NewClientFromConnectionString(connStr, nil)
	if err != nil {
		cc.logger.Error("could not create service bus client", zap.Error(err))
		return
	}

	return
}

func (c *sbClient) CreateQueue(ctx context.Context, q string) (err error) {
	_, err = c.azClient.CreateQueue(ctx, q, nil)
	if err != nil {
		// todo err
		return
	}
	return
}

func (c *sbClient) DeleteQueue(ctx context.Context, q string) (err error) {
	_, err = c.azClient.DeleteQueue(ctx, q, nil)
	if err != nil {
		// todo err
		return
	}
	return
}

func (c *sbClient) SendMessage(ctx context.Context, q string, payload []byte) (err error) {
	// fixme: get queue name from event
	fmt.Println("SendMessage q: ", q)
	fmt.Println("SendMessage payload: ", string(payload))
	fmt.Println("SendMessage client 000: ", c)
	fmt.Println("SendMessage client: ", c.client)
	sender, err := c.client.NewSender(q, nil)
	if err != nil {
		// todo err
		return
	}
	defer func(sender *azservicebus.Sender, ctx context.Context) {
		err = sender.Close(ctx)
		fmt.Println("SENDMESSAGE err: ", err)
		if err != nil {
			// fixme err
			return
		}
	}(sender, ctx)

	sbMessage := &azservicebus.Message{
		Body: payload,
	}
	err = sender.SendMessage(ctx, sbMessage, nil)
	fmt.Println("SENDMESSAGE err 3333: ", err)
	if err != nil {
		// todo err
		return
	}

	fmt.Println("SENDMESSAGE ENDDDD: ", string(payload))
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
// func (c *sbClient) GetMessage(ctx context.Context, count int, qm types.QueueMessage) {
// 	receiver, err := c.client.NewReceiverForQueue(qm.Queue, nil) // Change myqueue to env var
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer func(receiver *azservicebus.Receiver, ctx context.Context) {
// 		err := receiver.Close(ctx)
// 		if err != nil {
// 			// todo err
// 			return
// 		}
// 	}(receiver, ctx)
//
// 	messages, err := receiver.ReceiveMessages(ctx, count, nil)
// 	if err != nil {
// 		// todo err
// 		return
// 	}
//
// 	for _, message := range messages {
// 		body := message.Body
// 		// fixme: handle message
// 		fmt.Printf("%s\n", string(body))
//
// 		err = receiver.CompleteMessage(ctx, message, nil)
// 		if err != nil {
// 			// todo err
// 			return
// 		}
// 	}
// }
