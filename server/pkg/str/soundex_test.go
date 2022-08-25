package str

import (
	"testing"
)

func Test_soundex(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			"Robert",
			"R163",
		},
		{
			"Rupert",
			"R163",
		},
		{
			"Rubin",
			"R150",
		},
		{
			"Ashcraft",
			"A261",
		},
		{
			"Ashcroft",
			"A261",
		},
		{
			"Tymczak",
			"T522",
		},
		{
			"Pfister",
			"P123",
		},
		{
			"AH KEY",
			"A000",
		},
		{
			"The quick brown fox",
			"T221",
		},
		{
			"h3110 w021d",
			"3000",
		},
		{
			"1337",
			"1000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToSoundex(tt.name); got != tt.want {
				t.Errorf("soundex() = %v, want %v", got, tt.want)
			}
		})
	}
}
