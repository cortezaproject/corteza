package saml

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"

	"github.com/cortezaproject/corteza-server/system/types"
)

const settingsPrefix = "auth.external.providers.saml"

type (
	Cert struct {
		Cert []byte
		Key  []byte
	}

	certManager struct {
		loader CertLoader
	}

	CertLoader interface {
		Load(ctx context.Context) (*Cert, error)
		Key(ctx context.Context) ([]byte, error)
		Cert(ctx context.Context) ([]byte, error)
	}

	CertStoreLoader struct {
		Storer storer
	}

	CertFsLoader struct {
		KeyFile  string
		CertFile string
	}

	storer interface {
		SearchSettings(ctx context.Context, f types.SettingsFilter) (types.SettingValueSet, types.SettingsFilter, error)
	}
)

func (cm *certManager) Load(ctx context.Context) (cc *Cert, err error) {
	return cm.loader.Load(ctx)
}

func (cm *certManager) Parse(cert []byte, key []byte) (tls.Certificate, error) {
	return tls.X509KeyPair(cert, key)
}

func (l *CertStoreLoader) Cert(ctx context.Context) (c []byte, err error) {
	return l.get(ctx, "cert")
}

func (l *CertStoreLoader) Key(ctx context.Context) ([]byte, error) {
	return l.get(ctx, "key")
}

func (l *CertStoreLoader) get(ctx context.Context, key string) (c []byte, err error) {
	var s types.SettingValueSet

	if s, _, err = l.Storer.SearchSettings(ctx, types.SettingsFilter{Prefix: settingsPrefix}); err != nil {
		return
	}

	if cc := s.First(fmt.Sprintf("%s.%s", settingsPrefix, key)); cc != nil {
		return []byte(cc.String()), nil
	}

	return
}

func (l *CertStoreLoader) Load(ctx context.Context) (cc *Cert, err error) {
	var (
		c, k []byte
	)

	if c, err = l.Cert(ctx); err != nil && c != nil {
		return
	}

	if k, err = l.Key(ctx); err != nil && k != nil {
		return
	}

	cc = &Cert{Cert: c, Key: k}

	return
}

func (l *CertFsLoader) Cert(ctx context.Context) (c []byte, err error) {
	return ioutil.ReadFile(l.CertFile)
}

func (l *CertFsLoader) Key(ctx context.Context) ([]byte, error) {
	return ioutil.ReadFile(l.KeyFile)
}

func (l *CertFsLoader) Load(ctx context.Context) (cc *Cert, err error) {
	var (
		c, k []byte
	)

	if c, err = l.Cert(ctx); err != nil && c != nil {
		return
	}

	if k, err = l.Key(ctx); err != nil && k != nil {
		return
	}

	cc = &Cert{Cert: c, Key: k}

	return
}

func NewCertManager(l CertLoader) *certManager {
	return &certManager{
		loader: l,
	}
}
