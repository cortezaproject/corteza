package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSyncerURIString(t *testing.T) {
	var (
		req = require.New(t)
	)

	now, _ := time.Parse("2006-01-02 15:04:05", "2020-10-23 11:11:11")

	tests := []struct {
		name   string
		url    *SyncerURI
		expect string
	}{
		{
			name:   "next page cursor",
			url:    &SyncerURI{BaseURL: "https://example.url", NextPage: "NEXT_PAGE_CURSOR=="},
			expect: "https://example.url?limit=0&pageCursor=NEXT_PAGE_CURSOR==",
		},
		{
			name:   "last page cursor",
			url:    &SyncerURI{BaseURL: "https://example.url", LastPage: "LAST_PAGE_CURSOR=="},
			expect: "https://example.url?limit=0&pageCursor=LAST_PAGE_CURSOR==",
		},
		{
			name:   "limit results",
			url:    &SyncerURI{BaseURL: "https://example.url", Limit: 666},
			expect: "https://example.url?limit=666",
		},
		{
			name:   "base path",
			url:    &SyncerURI{Path: "/relative/path"},
			expect: "/relative/path?limit=0",
		},
		{
			name:   "last sync",
			url:    &SyncerURI{Path: "/relative/path", LastSync: &now},
			expect: "/relative/path?limit=0&lastSync=2020-10-23T11:11:11Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := tt.url.String()

			req.Equal(tt.expect, r)
		})
	}
}

func TestSyncerURIParse(t *testing.T) {
	var (
		req = require.New(t)
	)

	// TODO - add a parse last sync date
	// now, _ := time.Parse("2006-01-02 15:04:05", "2020-10-23 11:11:11")

	tests := []struct {
		name   string
		url    string
		expect *SyncerURI
	}{
		{
			name:   "parse limit",
			url:    "https://example.url?limit=11",
			expect: &SyncerURI{BaseURL: "https://example.url", Limit: 11},
		},
		{
			name:   "parse path",
			url:    "/path/to/endpoint/",
			expect: &SyncerURI{Path: "/path/to/endpoint/"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SyncerURI{}
			err := s.Parse(tt.url)

			req.Equal(tt.expect, s)
			req.NoError(err)
		})
	}
}
