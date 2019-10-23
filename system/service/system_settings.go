package service

import (
	intset "github.com/cortezaproject/corteza-server/pkg/settings"
)

type (
	SystemSettings struct {
		DefaultLogo string
		MailHeader  string
		MailFooter  string
	}
)

// ParseAuthSettings maps from plain values to AuthSettings struct
//
// see settings.Initialize() func
func ParseSystemSettings(kv intset.KV) (ss *SystemSettings, err error) {
	ss = &SystemSettings{}
	ss.ReadKV(kv)
	return
}

func (ss *SystemSettings) ReadKV(kv intset.KV) (err error) {
	ss.DefaultLogo = kv.String("system.defaultLogo")
	ss.MailHeader = kv.String("system.mail.header.en")
	ss.MailFooter = kv.String("system.mail.footer.en")

	return
}
