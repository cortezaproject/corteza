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

		paths = map[string]string{
			"corredor CA public key":      opt.TlsCertCA,
			"corredor client public key":  opt.TlsCertPublic,
			"corredor client private key": opt.TlsCertPrivate,
		}

		expl = "\n\n" +
			"Check and change configured path to Corredor TLS certificates " +
			"in env var CORREDOR_CLIENT_CERTIFICATES_PATH or disable " +
			"this check with CORREDOR_CLIENT_CERTIFICATES_ENABLED=false" +
			"\n\n"
	)

	// opt.PublicKey = ""
	if opt.TlsCertEnabled {
		// Check paths

		for label, path := range paths {
			if _, err = os.Stat(path); os.IsNotExist(err) {
				return nil, fmt.Errorf("%s (%s) not found"+expl, label, path)
			} else if err != nil {
				return nil, fmt.Errorf("%s (%s) could not be loaded: %s"+expl, label, path, err)
			}
		}

		// Load the certificates from disk
		certificate, err := tls.LoadX509KeyPair(opt.TlsCertPublic, opt.TlsCertPrivate)
		if err != nil {
			return nil, fmt.Errorf("could not load client key pair: %s"+expl, err)
		}

		// Create a certificate pool from the certificate authority
		certPool := x509.NewCertPool()
		ca, err := ioutil.ReadFile(opt.TlsCertCA)
		if err != nil {
			return nil, fmt.Errorf("could not read ca certificate: %s"+expl, err)
		}

		// Append the client certificates from the CA
		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			return nil, errors.New("failed to append ca certs" + expl)
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

	if opt.MaxReceiveMessageSize > 0 {
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(opt.MaxReceiveMessageSize))
	}

	if opt.MaxBackoffDelay > 0 {
		dialOpts = append(dialOpts, grpc.WithBackoffMaxDelay(opt.MaxBackoffDelay))
	}

	return grpc.DialContext(ctx, opt.Addr, dialOpts...)
}
