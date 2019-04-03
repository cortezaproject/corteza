// +build unit

package mail

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/crusttech/crust/internal/config"
	"github.com/crusttech/crust/internal/test"
)

func TestDialerInvalidSetup(t *testing.T) {
	defaultDialer = nil
	defaultDialerError = nil

	SetupDialer(nil)
	test.Assert(t, defaultDialerError != nil, "'Missing SMTP configuration' error should be set, got: %v", defaultDialerError)
	test.Assert(t, defaultDialer == nil, "defaultDialer should n be set, got: %v", defaultDialer)
}

func TestDialerValidSetup(t *testing.T) {
	defaultDialer = nil
	defaultDialerError = nil

	cfg := &config.SMTP{
		Host: "localhost:321",
		From: "some@email.tld",
	}
	cfg.Validate()

	SetupDialer(cfg)
	test.Assert(t, defaultDialerError == nil, "defaultDialerError should be nil, got %v", defaultDialerError)
	test.Assert(t, defaultDialer != nil, "defaultDialer should be set, got %v", defaultDialer)

}

func TestMailSendWithoutDialer(t *testing.T) {
	msg := New()
	defaultDialer = nil
	defaultDialerError = nil
	{
		err := Send(msg)
		test.Assert(t, err != nil, "Send() should return an error, got %v", err)
	}

	defaultDialer = nil
	defaultDialerError = errors.New("Default dialer init error")
	{
		err := Send(msg)
		test.Assert(t, err != nil, "Send() should return an error, got %v", err)
	}
}

func TestMailSendWithDefaultDialer(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	msg := New()

	dDialer := NewMockDialer(mockCtrl)
	dDialer.EXPECT().DialAndSend(msg).Times(1).Return(nil)

	defaultDialerError = nil
	defaultDialer = dDialer

	test.NoError(t, Send(msg), "Send() returned an error: %v")
	defaultDialer = nil
}

func TestMailSendWithSpecificDialer(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	msg := New()

	sDailer := NewMockDialer(mockCtrl)
	sDailer.EXPECT().DialAndSend(msg).Times(1).Return(nil)

	defaultDialerError = nil
	dDialer := NewMockDialer(mockCtrl)
	dDialer.EXPECT().DialAndSend(msg).Times(0)

	defaultDialer = dDialer
	test.NoError(t, Send(msg, sDailer), "Send() returned an error: %v")
	defaultDialer = nil
}

func TestMailSendErrors(t *testing.T) {
	defaultDialerError = nil
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	msg := New()

	sDailer := NewMockDialer(mockCtrl)
	sDailer.EXPECT().DialAndSend(msg).Times(1).Return(errors.New("some-error"))

	err := Send(msg, sDailer)
	test.Assert(t, err != nil, "Send() should return an error, got: %v", err)
}

func TestMailValidator(t *testing.T) {
	ttc := []struct {
		addr string
		ok   bool
	}{
		{"ç$€§/az@gmail.com", false},
		{"abcd@gmail_yahoo.com", false},
		{"abcd@gmail-yahoo.com", true},
		{"abcd@gmailyahoo", true},
		{"abcd@gmail.yahoo", true},
		{"info @ crust tech", false},
		{"info@crust.tech", true},
	}

	for _, tc := range ttc {
		test.Assert(t, IsValidAddress(tc.addr) == tc.ok, "Validation of %s should return %v", tc.addr, tc.ok)
	}
}
