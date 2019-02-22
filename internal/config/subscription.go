package config

import (
	"github.com/namsral/flag"
)

type (
	Subscription struct {
		Key    string
		Domain string
	}
)

var subscription *Subscription

func (c *Subscription) Validate() error {
	if c == nil {
		return nil
	}

	return nil
}

func (*Subscription) Init(prefix ...string) *Subscription {
	if subscription != nil {
		return subscription
	}

	subscription = new(Subscription)
	flag.StringVar(&subscription.Key, "subscription-key", "", "Subscription key")
	flag.StringVar(&subscription.Domain, "subscription-domain", "", "Domain")
	return subscription
}
