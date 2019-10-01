package mail

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestDialerInvalidSetup(t *testing.T) {
	defaultDialer = nil
	defaultDialerError = nil

	SetupDialer("", 0, "", "", "")
	require.True(t, defaultDialerError != nil, "'Missing SMTP configuration' error should be set, got: %v", defaultDialerError)
	require.True(t, defaultDialer == nil, "defaultDialer should n be set, got: %v", defaultDialer)
}

func TestDialerValidSetup(t *testing.T) {
	defaultDialer = nil
	defaultDialerError = nil

	SetupDialer("localhost:321", 0, "", "", "some@email.tld")
	require.True(t, defaultDialerError == nil, "defaultDialerError should be nil, got %v", defaultDialerError)
	require.True(t, defaultDialer != nil, "defaultDialer should be set, got %v", defaultDialer)

}

func TestMailSendWithoutDialer(t *testing.T) {
	msg := New()
	defaultDialer = nil
	defaultDialerError = nil
	{
		err := Send(msg)
		require.True(t, err != nil, "Send() should return an error, got %v", err)
	}

	defaultDialer = nil
	defaultDialerError = errors.New("Default dialer init error")
	{
		err := Send(msg)
		require.True(t, err != nil, "Send() should return an error, got %v", err)
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

	require.NoError(t, Send(msg))
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
	require.NoError(t, Send(msg, sDailer))
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
	require.True(t, err != nil, "Send() should return an error, got: %v", err)
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
		require.True(t, IsValidAddress(tc.addr) == tc.ok, "Validation of %s should return %v", tc.addr, tc.ok)
	}
}
