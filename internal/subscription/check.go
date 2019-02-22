package subscription

import (
	"context"

	"github.com/pkg/errors"

	"github.com/crusttech/permit/pkg/permit"
)

// Check for subscription
func Check(ctx context.Context) error {
	p := permit.Permit{
		Key:    flags.subscription.Key,
		Domain: flags.subscription.Domain,
	}

	// Do not collect stats on local domains.
	// if p.Domain != "local.crust.tech" {
	// @todo collect & pass attributes (no of users....) to be validated by permit.crust.tech subscription server.
	p.Attributes = map[string]int{
		"messaging.enabled": 1,
		// "messaging.max-public-channels": 1,
		// "messaging.max-messages": 1,
		// "messaging.max-users": 1,
		// "messaging.max-private-channels": 1,

		"system.enabled": 1,
		// "system.max-organisations": 1,
		// "system.max-users": 1,
		// "system.max-teams": 1,

		"compose.enabled": 1,
		// "compose.max-modules": 1,
		// "compose.max-pages": 1,
		// "compose.max-triggers": 1,
		// "compose.max-users": 1,
		// "compose.max-namespaces": 1,
		// "compose.max-charts": 1,
	}
	// }

	if p, err := permit.Check(ctx, p); err != nil {
		return errors.Wrap(err, "unable to check for licence")
	} else if !p.IsValid() {
		return err
	} else if p.Domain != flags.subscription.Domain {
		return errors.Errorf("subscription domains do not match (%s <> %s)", p.Domain, flags.subscription.Domain)
	}

	return nil
}
