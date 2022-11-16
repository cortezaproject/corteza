package corredor

import (
	"github.com/cortezaproject/corteza/server/pkg/slice"
	"strings"
)

type (
	Filter struct {
		ResourceTypePrefixes []string `json:"resourceTypePrefixes"`
		ResourceTypes        []string `json:"resourceTypes"`
		EventTypes           []string `json:"eventTypes"`
		ExcludeServerScripts bool     `json:"excludeServerScripts"`
		ExcludeClientScripts bool     `json:"excludeClientScripts"`
		ExcludeInvalid       bool     `json:"excludeInvalid"`

		//Page    uint `json:"page"`
		//PerPage uint `json:"perPage"`
		Count uint `json:"count"`
	}
)

func (f Filter) makeFilterFn() func(s *Script) (b bool, err error) {
	collectTypes := func(tt ...*Trigger) (rr []string, ee []string) {
		for _, t := range tt {
			rr = append(rr, t.ResourceTypes...)
			ee = append(ee, t.EventTypes...)
		}

		return
	}

	return func(s *Script) (b bool, err error) {
		if f.ExcludeInvalid {
			if len(s.Errors) > 0 {
				// Skip scripts with errors
				return
			}
		}

		var (
			resourceTypes, eventTypes = collectTypes(s.Triggers...)
		)

		if !f.checkPrefixes(resourceTypes...) {
			return
		}

		if !f.checkEventTypes(eventTypes...) {
			return
		}

		if !f.checkResourceTypes(resourceTypes...) {
			return
		}

		return true, nil
	}
}

// parses resource type prefixes, adds ui: and service prefix if not present
func (f *Filter) procRTPrefixes(service string) {
	var (
		hasService, hasUi bool
	)

	if service == "" {
		return
	}

	for i, rtp := range f.ResourceTypePrefixes {
		if rtp == service {
			hasService = true
		}

		if strings.HasPrefix(rtp, "ui") {
			hasUi = true
		}

		// Check all requested resource-type prefixed and make sure
		// all start with either "ui" or resourcePrefix + ":"
		if !strings.HasPrefix(rtp, "ui:") && !strings.HasPrefix(rtp, service) {
			f.ResourceTypePrefixes[i] = service + ":" + rtp
		}
	}

	if !hasService {
		f.ResourceTypePrefixes = append(f.ResourceTypePrefixes, service)
	}

	if !hasUi {
		f.ResourceTypePrefixes = append(f.ResourceTypePrefixes, "ui:")
	}
}

func (f Filter) checkPrefixes(rr ...string) bool {
	if len(f.ResourceTypePrefixes) == 0 {
		return true
	}

	if len(rr) == 0 {
		return false
	}

	for _, p := range f.ResourceTypePrefixes {
		for _, r := range rr {
			if strings.HasPrefix(r, p) {
				return true
			}
		}
	}

	return false
}

func (f Filter) checkResourceTypes(rr ...string) bool {
	return len(f.ResourceTypes) == 0 || len(slice.IntersectStrings(f.ResourceTypes, rr)) > 0
}

func (f Filter) checkEventTypes(ee ...string) bool {
	return len(f.EventTypes) == 0 || len(slice.IntersectStrings(f.EventTypes, ee)) > 0
}
