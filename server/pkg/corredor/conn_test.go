package corredor

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"math/rand"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"

	"github.com/crusttech/human/server/pkg/options"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestNewConnectionWithDisabled(t *testing.T) {
	c, err := NewConnection(nil, options.CorredorOpt{Enabled: false}, zap.NewNop())
	assert.Nil(t, c)
	assert.Nil(t, err)
}

func TestNewConnection(t *testing.T) {
	var (
		ctx = context.Background()

		a  = assert.New(t)
		wg = &sync.WaitGroup{}

		lstnr      = openListener(t)
		grpcServer = grpc.NewServer()

		opt = options.CorredorOpt{
			Enabled:         true,
			MaxBackoffDelay: 1,
			Addr:            lstnr.Addr().String(),
		}
	)

	a.NotNil(lstnr)
	defer lstnr.Close()
	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := grpcServer.Serve(lstnr); err != nil && err != grpc.ErrServerStopped {
			a.NoError(err)
		}
	}()

	grpcClientConn, err := NewConnection(ctx, opt, zap.NewNop())
	a.NoError(err)

	// wait until connected, the dial may already be done
	waitCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	for state := grpcClientConn.GetState(); state != connectivity.Ready; state = grpcClientConn.GetState() {
		if !grpcClientConn.WaitForStateChange(waitCtx, state) {
			a.FailNow("connection did not become ready")
		}
	}

	grpcServer.GracefulStop()
	lstnr.Close()

	t.Log("waiting for connection to close")
	wg.Wait()
}

func openListener(t *testing.T) (ln net.Listener) {
	var (
		tries   = 1000
		minPort = 50000
		maxPort = 63000
		port    int
		err     error
	)

	for tries > 0 {
		port = minPort + rand.Intn(maxPort-minPort)
		t.Logf("trying to find available port for gRPC connection: %d", port)
		ln, err = net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
		if err == nil {
			return ln
		}

		t.Errorf("error: %s", err)
		tries--
	}

	return nil
}
