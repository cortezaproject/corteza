package types

import (
	"fmt"
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
			expect: fmt.Sprintf("/relative/path?limit=0&lastSync=%d", now.Unix()),
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

	tests := []struct {
		name   string
		url    string
		expect *SyncerURI
		err    string
	}{
		{
			name:   "parse limit",
			url:    "https://example.url?limit=11",
			expect: &SyncerURI{BaseURL: "https://example.url", Limit: 11},
			err:    "",
		},
		{
			name:   "parse path",
			url:    "/path/to/endpoint/",
			expect: &SyncerURI{Path: "/path/to/endpoint/"},
			err:    "",
		},
		{
			name:   "parse invalid url",
			url:    "ht tps:/ /in valid",
			expect: &SyncerURI{},
			err:    `parse "ht tps:/ /in valid": first path segment in URL cannot contain colon`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SyncerURI{}
			err := s.Parse(tt.url)

			errString := ""
			if err != nil {
				errString = err.Error()
			}

			req.Equal(tt.expect, s)
			req.Equal(tt.err, errString)
		})
	}
}

func TestSyncerURIParseLastSync(t *testing.T) {
	var (
		req = require.New(t)
	)

	tests := []struct {
		name   string
		url    string
		expect string
	}{
		{
			name:   "parse timestamp",
			url:    "https://example.url?limit=11&lastSync=1603451471",
			expect: "2020-10-23T13:11:11+02:00",
		},
		{
			name:   "parse RFC3339",
			url:    "https://example.url?limit=11&lastSync=2020-10-23T13:11:11Z",
			expect: "2020-10-23T13:11:11Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SyncerURI{}
			err := s.Parse(tt.url)

			req.NoError(err)
			req.Equal(tt.expect, s.LastSync.Format(time.RFC3339))
		})
	}

}
