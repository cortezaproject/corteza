package sam

import (
	"testing"
)

func TestUser_SetPassword(t *testing.T) {
	const dummyPwd = "Lorem ipsum dolor sit amet, consectetur adipiscing elit..."
	dummy := &User{}
	initChLen := len(dummy.changed)
	dummy.SetPassword(dummyPwd)

	if dummyPwd == string(dummy.Password) {
		t.Error("Internal password value should be encrypted")
	}

	if len(dummy.changed) <= initChLen {
		t.Error("Password change should be recorded")
	}

	if !dummy.ValidatePassword(dummyPwd) {
		t.Error("Revalidating password should succeed")
	}
}
