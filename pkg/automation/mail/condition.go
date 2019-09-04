package mail

import (
	"encoding/json"
	"errors"
	"net/mail"
	"net/textproto"
	"regexp"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/automation"
	"github.com/cortezaproject/corteza-server/system/types"
)

// Trigger condition:
// Matcher for mail headers

type (
	Condition struct {
		MatchAll bool            `json:"matchAll"`
		Headers  []HeaderMatcher `json:"headers"`
	}

	HeaderMatcher struct {
		Name  HMName `json:"name"`
		Op    HMOp   `json:"op"`
		Match string `json:"match"`

		// Compiled regexp
		re *regexp.Regexp
	}

	HMName string
	HMOp   string

	userExistanceVerifier func(string) bool
)

const (
	HeaderMatchNameFrom    HMName = "from"
	HeaderMatchNameTo             = "to"
	HeaderMatchNameCC             = "cc"
	HeaderMatchNameBCC            = "bcc"
	HeaderMatchNameReplyTo        = "reply-to"
	HeaderMatchNameSubject        = "subject"

	// Keeping -ci suffix in case we get
	// feature request to separete ci & cs operators
	HMOpEqualCi  HMOp = "equal-ci"
	HMOpSuffixCi      = "suffix-ci"
	HMOpPrefixCi      = "prefix-ci"
	HMOpRegex         = "regex"
	HMOpUser          = "user"
)

var (
	ErrUnknownHeaderMatcherName     = errors.New("unknown header matcher field")
	ErrUnknownHeaderMatcherOperator = errors.New("unknown header matcher operator")
	ErrInvalidHeaderMatcherValue    = errors.New("invalid header matcher value")
)

func (c *Condition) Prepare() (err error) {
	for i := range c.Headers {
		err = c.Headers[i].prepare()
		if err != nil {
			return
		}
	}

	return
}

// IsValid verifies if header matcher is valid
func (m *HeaderMatcher) prepare() (err error) {
	switch m.Name {
	case HeaderMatchNameFrom,
		HeaderMatchNameTo,
		HeaderMatchNameCC,
		HeaderMatchNameBCC,
		HeaderMatchNameReplyTo,
		HeaderMatchNameSubject:
	// ok fields
	default:
		return ErrUnknownHeaderMatcherName
	}

	switch m.Op {
	case HMOpRegex:
		// Try to compile given regex
		m.re, err = regexp.Compile(m.Match)
		if err != nil {
			return ErrInvalidHeaderMatcherValue
		}
	case HMOpUser:
		// When matching against existing user,
		// there should be no value set
		if m.Match != "" {
			return ErrInvalidHeaderMatcherValue
		}
	case HMOpEqualCi, HMOpSuffixCi, HMOpPrefixCi:
		// no special validation here
		m.Match = strings.ToLower(m.Match)
	default:
		return ErrUnknownHeaderMatcherOperator
	}

	return
}

func (n HMName) match(name string) bool {
	return string(n) == strings.ToLower(name)
}

// IsMatch checks if header matcher matches against given headers
func (m *HeaderMatcher) isMatch(header mail.Header, exists userExistanceVerifier, matchAll bool) (match bool) {
	var lcHeader string

	for name, vv := range header {
		if !m.Name.match(name) {
			continue
		}

		for _, v := range vv {
			lcHeader = strings.ToLower(v)

			switch m.Op {
			case HMOpEqualCi:
				match = m.Match == lcHeader
			case HMOpSuffixCi:
				match = strings.HasSuffix(lcHeader, m.Match)
			case HMOpPrefixCi:
				match = strings.HasPrefix(lcHeader, m.Match)
			case HMOpRegex:
				match = m.re.MatchString(v)
			case HMOpUser:
				match = exists(v)
			default:
				return false
			}

			if !match && matchAll {
				// fail in first non-match
				return false
			} else if match && !matchAll {
				// match in first
				return true
			}
		}
	}

	return match
}

func (c *Condition) CheckHeader(header mail.Header, uev userExistanceVerifier) (match bool) {
	_ = c.Prepare()

	// Pre-process & simplify header values: parse all addresses,
	// extract emails and toss away names, we do not need them
	for name := range header {
		switch textproto.CanonicalMIMEHeaderKey(name) {
		case "From", "To", "Cc", "Bcc", "Reply-To":
			for i, v := range header[name] {
				addr, _ := mail.ParseAddress(v)
				header[name][i] = strings.Trim(addr.Address, "><")
			}
		}
	}

	for _, h := range c.Headers {
		match = h.isMatch(header, uev, c.MatchAll)

		if !match && c.MatchAll {
			// fail in first non-match
			return false
		} else if match && !c.MatchAll {
			// match in first
			return true
		}
	}

	return match
}

func MakeChecker(headers types.MailMessageHeader, uev userExistanceVerifier) automation.TriggerConditionChecker {
	return func(c string) bool {
		var (
			tc = Condition{}
		)

		if err := json.Unmarshal([]byte(c), &tc); err == nil {
			return tc.CheckHeader(headers.Raw, uev)
		}

		return false
	}
}
