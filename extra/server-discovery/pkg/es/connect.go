package es

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"

	"github.com/cortezaproject/corteza/extra/server-discovery/pkg/options"
	"github.com/elastic/go-elasticsearch/v7"
)

func Connect(opt options.EsOpt) (*elasticsearch.Client, error) {
	config := elasticsearch.Config{
		Addresses:            opt.Addresses,
		EnableRetryOnTimeout: opt.EnableRetryOnTimeout,
		MaxRetries:           opt.MaxRetries,
	}
	if len(opt.Username) > 0 {
		config.Username = opt.Username
		config.Password = opt.Password
	}

	// if the user provided a CA certificate file, we need to load it
	// and set it in the TLS config
	if opt.CertFile != "" {
		rootCACert, err := os.ReadFile(opt.CertFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA certificate file: %w", err)
		}

		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(rootCACert) {
			return nil, fmt.Errorf("failed parse root certificate: %w", err)
		}

		config.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		}
	} else if !opt.Secure {
		config.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	return elasticsearch.NewClient(config)
}
