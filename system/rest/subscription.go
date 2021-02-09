package rest

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/system/rest/request"
)

type Subscription struct{}

func (Subscription) New() *Subscription {
	return &Subscription{}
}

func (ctrl *Subscription) Current(ctx context.Context, r *request.SubscriptionCurrent) (interface{}, error) {
	return api.OK(), nil
	//if service.CurrentSubscription == nil {
	//	// Nothing to do here
	//}
	//
	//// Returning function that gets called with writter & request
	////
	//// This is the only way to get to the request URL we need to do a domain check
	//// for the permit
	//return func(w http.ResponseWriter, r *http.Request) {
	//	var (
	//		domain = r.Host
	//		pos    = strings.IndexByte(domain, ':')
	//
	//		// Anyone that has access permissions is considered admin
	//		isAdmin = service.DefaultAccessControl.CanAccess(ctx)
	//	)
	//
	//	if pos > -1 {
	//		// Strip port
	//		domain = domain[:pos]
	//	}
	//
	//	if err := service.CurrentSubscription.Validate(domain, isAdmin); err != nil {
	//		api.Send(w, r, err)
	//	} else {
	//		api.Send(w, r, api.OK())
	//	}
	//}, nil
}
