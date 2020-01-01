package corredor

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"

	"github.com/cortezaproject/corteza-server/pkg/app/options"
)

// Corredor standard connector to Corredor service via gRPC
func NewConnection(ctx context.Context, opt options.CorredorOpt, logger *zap.Logger) (c *grpc.ClientConn, err error) {
	if !opt.Enabled {
		// Do not connect when script runner is not enabled
		return
	}

	if opt.Log {
		// Send logs to zap
		//
		// waiting for https://github.com/uber-go/zap/pull/538
		grpclog.SetLogger(zapgrpc.NewLogger(logger.Named("grpc")))
	}

	var (
		dialOpts = make([]grpc.DialOption, 0)
	)

	// opt.PublicKey = ""
	if opt.TlsCertEnabled {
		// Check paths
		paths := map[string]string{
			"CA public key ":     opt.TlsCertCA,
			"client public key":  opt.TlsCertPublic,
			"client private key": opt.TlsCertPrivate,
		}
		for label, path := range paths {
			if _, err = os.Stat(path); os.IsNotExist(err) {
				return nil, fmt.Errorf("%s (%s) not found", label, path)
			} else if err != nil {
				return nil, fmt.Errorf("%s (%s) could not be loaded: %s", label, path, err)
			}
		}

		// Load the certificates from disk
		certificate, err := tls.LoadX509KeyPair(opt.TlsCertPublic, opt.TlsCertPrivate)
		if err != nil {
			return nil, fmt.Errorf("could not load client key pair: %s", err)
		}

		// Create a certificate pool from the certificate authority
		certPool := x509.NewCertPool()
		ca, err := ioutil.ReadFile(opt.TlsCertCA)
		if err != nil {
			return nil, fmt.Errorf("could not read ca certificate: %s", err)
		}

		// Append the client certificates from the CA
		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			return nil, errors.New("failed to append ca certs")
		}

		crds := credentials.NewTLS(&tls.Config{
			ServerName:   opt.TlsServerName,
			Certificates: []tls.Certificate{certificate},
			RootCAs:      certPool,
		})

		dialOpts = append(dialOpts, grpc.WithTransportCredentials(crds))
	} else {
		dialOpts = append(dialOpts, grpc.WithInsecure())
	}

	if opt.MaxBackoffDelay > 0 {
		dialOpts = append(dialOpts, grpc.WithBackoffMaxDelay(opt.MaxBackoffDelay))
	}

	return grpc.DialContext(ctx, opt.Addr, dialOpts...)
}
