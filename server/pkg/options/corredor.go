package options

import (
	"path"
)

func (o *CorredorOpt) Cleanup() {
	o.TlsCertCA = path.Join(o.TlsCertPath, o.TlsCertCA)
	o.TlsCertPrivate = path.Join(o.TlsCertPath, o.TlsCertPrivate)
	o.TlsCertPublic = path.Join(o.TlsCertPath, o.TlsCertPublic)
}
