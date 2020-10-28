package types

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
)

// SyncerURI serves as a default struct for handling url queues
// for fetcher
type (
	SyncerURI struct {
		Limit    int
		LastSync *time.Time
		// todo
		// LastSync *LastSyncTime
		Path     string
		BaseURL  string
		NextPage string
		LastPage string
	}
)

// Does the transformation to string, does not
// url encode page cursor, since it needs to be base64
// encoded
func (s *SyncerURI) String() (string, error) {
	var query []string

	// using query instead of url.Values
	// because pageCursor should not be url encoded
	query = append(query, fmt.Sprintf("limit=%d", s.Limit))

	if s.NextPage != "" {
		query = append(query, fmt.Sprintf("pageCursor=%s", s.NextPage))
	} else if s.LastPage != "" {
		query = append(query, fmt.Sprintf("pageCursor=%s", s.LastPage))
	}

	if s.LastSync != nil {
		query = append(query, fmt.Sprintf("lastSync=%d", s.LastSync.Unix()))
	}

	return fmt.Sprintf("%s%s?%s", s.BaseURL, s.Path, strings.Join(query[:], "&")), nil
}

// Parse the uri to the struct
func (s *SyncerURI) Parse(uri string) error {
	u, err := url.Parse(uri)

	if err != nil {
		return err
	}

	var limit int
	l := u.Query().Get("limit")

	if l != "" {
		if limit, err = strconv.Atoi(l); err != nil {
			return err
		}
	}

	s.Path = u.Path
	s.Limit = limit

	if u.Host != "" {
		s.BaseURL = fmt.Sprintf("%s://%s", u.Scheme, u.Host)
	}

	ls := u.Query().Get("lastSync")

	if ls != "" {
		parsed, err := parseLastSync(ls)

		if err != nil {
			return err
		}

		s.LastSync = parsed
	}

	return nil
}

func parseLastSync(lastSync string) (*time.Time, error) {
	spew.Dump("parse last sync")
	if i, err := strconv.ParseInt(lastSync, 10, 64); err == nil {
		spew.Dump("returning")
		t := time.Unix(i, 0)
		return &t, nil
	}

	// try different format if above fails
	spew.Dump("TRY HERE")
	if t, err := time.Parse(time.RFC3339, lastSync); err == nil {
		return &t, nil
	}

	t, err := time.Parse("2006-01-02", lastSync)

	if err != nil {
		return nil, err
	}

	return &t, nil
}
