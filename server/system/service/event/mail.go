package event

import (
	"github.com/cortezaproject/corteza/server/pkg/eventbus"
	"github.com/cortezaproject/corteza/server/system/types"
	"net/mail"
)

// Match returns false if given conditions do not match event & resource internals
func (res mailBase) Match(c eventbus.ConstraintMatcher) bool {
	// By default we match no mather what kind of constraints we receive
	//
	// Function will be called multiple times - once for every trigger constraint
	// All should match (return true):
	//   constraint#1 AND constraint#2 AND constraint#3 ...
	//
	// When there are multiple values, Match() can decide how to treat them (OR, AND...)
	return mailMatch(res.message, c)
}

// Handles role matchers
func mailMatch(r *types.MailMessage, c eventbus.ConstraintMatcher) bool {
	switch c.Name() {
	case "message.header.subject", "mail.header.subject":
		return c.Match(r.Subject)
	case "message.header.from", "mail.header.from":
		return mailMatchAnyAddress(c, r.Header.From...)
	case "message.header.to", "mail.header.to":
		return mailMatchAnyAddress(c, r.Header.To...)
	case "message.header.reply-to", "mail.header.replyTo":
		return mailMatchAnyAddress(c, r.Header.ReplyTo...)
	case "message.header.cc", "mail.header.cc":
		return mailMatchAnyAddress(c, r.Header.CC...)
	case "message.header.bcc", "mail.header.bcc":
		return mailMatchAnyAddress(c, r.Header.BCC...)
	}

	return false
}

func mailMatchAnyAddress(c eventbus.ConstraintMatcher, aa ...*mail.Address) bool {
	for _, a := range aa {
		if c.Match(a.Address) {
			return true
		}
	}
	return false
}
