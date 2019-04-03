// +build unit

package service

import (
	"testing"

	"github.com/crusttech/crust/internal/test"
)

func TestPermissionsValidation(t *testing.T) {
	test.Error(t, validatePermission("bogus", "bogus"), "expected error")
	test.Error(t, validatePermission("bogus", "bogus"), "expected error")
	test.Error(t, validatePermission("messaging:channel", "bogus"), "expected error")
	test.Error(t, validatePermission("messaging:channel:", "message.send"), "expected error")
	test.NoError(t, validatePermission("messaging:channel:1", "message.send"), "expected valid response")
}
