package rules

import (
	"fmt"
	"testing"

	"encoding/json"

	"github.com/crusttech/crust/internal/test"
)

func TestResource(t *testing.T) {
	var (
		assert = test.Assert
	)
	r := Resource{123, "Test name", "channel", "messaging"}
	assert(t, r.String() == "messaging:channel:123", "Resource ID doesn't match, messaging:channel:123 != '%s'", r.String())

	b, _ := json.Marshal(r)
	fmt.Println(string(b))

	{
		r := ResourceJSON{}
		json.Unmarshal(b, &r)
		assert(t, r.ResourceID == "messaging:channel:123", "Decoded full-json resource ID doesn't match, messaging:channel:123 != '%s'", r.ResourceID)
	}

	{
		r := Resource{}
		json.Unmarshal(b, &r)
		assert(t, r.String() == "messaging:channel:123", "Decoded full-json resource ID doesn't match, messaging:channel:123 != '%s'", r.String())
	}

	{
		r.ID = 0
		assert(t, r.String() == "", "Empty resource should return empty string, got '%s'", r.String())
	}
}
