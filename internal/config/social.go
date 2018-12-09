package config

import (
	"errors"
	"strings"

	"github.com/namsral/flag"
)

type (
	Social struct {
		Enabled bool

		FacebookKey    string
		FacebookSecret string
		GPlusKey       string
		GPlusSecret    string
		GitHubKey      string
		GitHubSecret   string
		LinkedInKey    string
		LinkedInSecret string

		Url string

		SessionStoreSecret string
		SessionStoreExpiry int // seconds!
	}
)

var social *Social

func (c *Social) Validate() error {
	if c == nil {
		return nil
	}

	if c.Enabled == false {
		return nil
	}

	if c.SessionStoreSecret == "" {
		return errors.New("Session store secret not set for SOCIAL")
	}

	return nil
}

func (*Social) Init(prefix ...string) *Social {
	if social != nil {
		return social
	}

	b := func(name string, k, s *string) {
		flag.StringVar(k, "auth-social-"+strings.ToLower(name)+"-key", "", name+" key")
		flag.StringVar(s, "auth-social-"+strings.ToLower(name)+"-secret", "", name+" secret")

	}

	social = new(Social)
	flag.BoolVar(&social.Enabled, "auth-social-enabled", true, "SocialAuth enabled")

	b("Facebook", &social.FacebookKey, &social.FacebookSecret)
	b("GPlus", &social.GPlusKey, &social.GPlusSecret)
	b("GitHub", &social.GitHubKey, &social.GitHubSecret)
	b("LinkedIn", &social.LinkedInKey, &social.LinkedInSecret)

	flag.StringVar(&social.Url, "auth-social-url", "", "Base URL")
	flag.StringVar(&social.SessionStoreSecret, "auth-social-session-store-secret", "", "Session store secret")
	flag.IntVar(&social.SessionStoreExpiry, "auth-social-state-cookie-expiry", 60*15, "SocialAuth State cookie expiry in seconds")
	return social
}
