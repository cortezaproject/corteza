// +build unit

package repository

import (
	"context"
	"testing"
	"time"

	"github.com/crusttech/crust/internal/config"
)

func TestPubSubMemory(t *testing.T) {
	p := PubSubMemory{}.New(&config.PubSub{
		PollingInterval: time.Second,
	})

	calledOnConnect := false
	calledOnMessage := 0

	onConnect := func() error {
		calledOnConnect = true
		return nil
	}

	onMessage := func(channel string, message []byte) error {
		calledOnMessage++
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)

	go func() {
		done <- p.Subscribe(ctx, "events", onConnect, onMessage)
	}()

	p.Publish(ctx, "events", "new message event")
	p.Publish(ctx, "events", "new message event")

	time.Sleep(2 * time.Millisecond)

	if !calledOnConnect {
		t.Fatalf("Expected initial call to 'onConnect'")
	}
	if calledOnMessage != 2 {
		t.Fatalf("Expected calledOnMessage to be 2, got %d", calledOnMessage)
	}

	cancel()

	select {
	case <-done:
	case <-time.After(10 * time.Millisecond):
		t.Fatalf("Expected PubSub channel exit after context cancellation")
	}

	if err := p.Publish(ctx, "events", "new message event"); err == nil {
		t.Fatalf("Expected error from sending message on closed channel")
	}
}
