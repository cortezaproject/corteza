package corredor

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"

	"github.com/cortezaproject/corteza-server/pkg/app/options"
	"github.com/cortezaproject/corteza-server/pkg/logger"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestNewConnectionWithDisabled(t *testing.T) {
	c, err := NewConnection(nil, options.CorredorOpt{Enabled: false}, nil)
	assert.Nil(t, c)
	assert.Nil(t, err)
}

func TestNewConnection(t *testing.T) {
	var (
		ctx = context.Background()

		dbgLog = logger.MakeDebugLogger()

		a  = assert.New(t)
		wg = &sync.WaitGroup{}

		lstnr      = openListener(t)
		grpcServer = grpc.NewServer()

		opt = options.CorredorOpt{
			Enabled:         true,
			Log:             true,
			MaxBackoffDelay: 1,
			Addr:            lstnr.Addr().String(),
		}
	)

	a.NotNil(lstnr)
	defer lstnr.Close()
	wg.Add(1)
	go func() {
		defer wg.Done()
		a.NoError(grpcServer.Serve(lstnr))
	}()

	grpcClientConn, err := NewConnection(ctx, opt, dbgLog)
	a.NoError(err)

	// Go and
	NewService(grpcClientConn, nil, dbgLog, opt)

	grpcClientConn.WaitForStateChange(ctx, connectivity.Ready)
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
