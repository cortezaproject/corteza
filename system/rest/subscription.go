package rest

import (
	"context"
	"net/http"
	"strings"

	"github.com/titpetric/factory/resputil"

	"github.com/cortezaproject/corteza-server/system/rest/request"
	"github.com/cortezaproject/corteza-server/system/service"
)

type Subscription struct{}

func (Subscription) New() *Subscription {
	return &Subscription{}
}

func (ctrl *Subscription) Current(ctx context.Context, r *request.SubscriptionCurrent) (interface{}, error) {
	if service.CurrentSubscription == nil {
		// Nothing to do here
		return resputil.OK(), nil
	}

	// Returning function that gets called with writter & request
	//
	// This is the only way to get to the request URL we need to do a domain check
	// for the permit
	return func(w http.ResponseWriter, req *http.Request) {
		var (
			domain = req.Host
			pos    = strings.IndexByte(domain, ':')

			// Anyone that has access permissions is considered admin
			isAdmin = service.DefaultAccessControl.CanAccess(ctx)
		)

		if pos > -1 {
			// Strip port
			domain = domain[:pos]
		}

		if err := service.CurrentSubscription.Validate(domain, isAdmin); err != nil {
			resputil.JSON(w, err)
		} else {
			resputil.JSON(w, resputil.OK())
		}
	}, nil
}
