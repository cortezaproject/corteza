package request

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/cortezaproject/corteza-server/pkg/rand"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

// cortezaSessionStore implements the session store and bridge to corteza store
type (
	cortezaSessionStore struct {
		// session store does not accept context on functions
		// so we'll store the general one here on struct so we
		// can capture termination signals..
		//ctx         context.Context
		store store.AuthSessions

		//cookies sessions.CookieStore
		//stopCleanup chan bool

		Codecs []securecookie.Codec

		Options *sessions.Options
	}
)

var (
	// wrapper around time.Now() that will aid service testing
	now = func() *time.Time {
		c := time.Now().Round(time.Second)
		return &c
	}
)

func init() {
	gob.Register(map[string]string{})
}

func CortezaSessionStore(store store.AuthSessions, opt options.AuthOpt) *cortezaSessionStore {
	var domain = opt.SessionCookieDomain
	if strings.Contains(domain, ":") {
		// do not set domain on the cookie if it contains port
		domain = ""
	}

	return &cortezaSessionStore{
		store:  store,
		Codecs: securecookie.CodecsFromPairs([]byte(opt.Secret)),
		Options: &sessions.Options{
			Path:     opt.SessionCookiePath,
			Domain:   domain,
			MaxAge:   int(opt.SessionLifetime / time.Second),
			Secure:   opt.SessionCookieSecure,
			HttpOnly: true,
		},
	}
}

func (s cortezaSessionStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(s, name)
}

func (s cortezaSessionStore) New(r *http.Request, name string) (*sessions.Session, error) {
	session := sessions.NewSession(s, name)
	session.Options = &sessions.Options{
		Path:     s.Options.Path,
		Domain:   s.Options.Domain,
		MaxAge:   s.Options.MaxAge,
		Secure:   s.Options.Secure,
		HttpOnly: s.Options.HttpOnly,
	}
	session.IsNew = true
	var err error
	if cook, errCookie := r.Cookie(name); errCookie == nil {
		err = securecookie.DecodeMulti(name, cook.Value, &session.ID, s.Codecs...)
		if err == nil {
			err = s.load(r.Context(), session)
			if err == nil {
				session.IsNew = false
			} else {
				err = nil
			}
		}
	}
	return session, err
}

func (s cortezaSessionStore) load(ctx context.Context, ses *sessions.Session) error {
	cortezaSession, err := s.store.LookupAuthSessionByID(ctx, ses.ID)
	if err != nil {
		return err
	}

	var maxAge = int(cortezaSession.ExpiresAt.Sub(*now()) / time.Second)
	if maxAge < 0 {
		return fmt.Errorf("session expired")
	}

	ses.Options.MaxAge = maxAge

	if err = gob.NewDecoder(bytes.NewReader(cortezaSession.Data)).Decode(&ses.Values); err != nil {
		return fmt.Errorf("failed to decode session: %w", err)
	}

	// Store original session back for ref
	ses.Values[keyOriginalSession] = cortezaSession

	return nil
}

func (s cortezaSessionStore) Save(r *http.Request, w http.ResponseWriter, ses *sessions.Session) (err error) {
	// Set delete if max-age is < 0
	if ses.Options.MaxAge < 0 {
		err = s.Delete(r.Context(), ses)
		if err != nil {
			return err
		}

		http.SetCookie(w, sessions.NewCookie(ses.Name(), "", ses.Options))

		newSes, err := s.New(r, ses.Name())
		if err != nil {
			return err
		}

		ses.Options = newSes.Options
		ses.IsNew = true
		ses.ID = ""
		ses.Values = nil
		return nil
	}

	if len(ses.Values) == 0 && ses.IsNew {
		// no values set, nothing to save
		return nil
	}

	if ses.ID == "" || ses.IsNew {
		ses.ID = string(rand.Bytes(64))
	}

	if err = s.save(r.Context(), ses); err != nil {
		return err
	}

	// Keep the session ID key in a cookie so it can be looked up in DB later.
	encoded, err := securecookie.EncodeMulti(ses.Name(), ses.ID, s.Codecs...)
	if err != nil {
		return err
	}

	http.SetCookie(w, sessions.NewCookie(ses.Name(), encoded, ses.Options))
	return nil
}

// save writes encoded session.Values to a database record.
// writes to http_sessions table by default.
func (s cortezaSessionStore) save(ctx context.Context, ses *sessions.Session) (err error) {
	var (
		buf            = &bytes.Buffer{}
		cortezaSession *types.AuthSession
	)

	// retrieve original session value
	if aux, is := ses.Values[keyOriginalSession]; is {
		cortezaSession = aux.(*types.AuthSession)
	} else {
		cortezaSession = &types.AuthSession{
			ID:        ses.ID,
			CreatedAt: *now(),
		}

		// new session does not belong to anyone yet.
		// retrieve user id from ses. values
		if au := GetAuthUser(ses); au != nil {
			cortezaSession.UserID = au.User.ID
		}

		extra := auth.GetExtraReqInfoFromContext(ctx)
		cortezaSession.UserAgent = extra.UserAgent
		cortezaSession.RemoteAddr = extra.RemoteAddr

		// calculate expiration date from max-age
		cortezaSession.ExpiresAt = now().Add(time.Second * time.Duration(ses.Options.MaxAge))
	}

	delete(ses.Values, keyOriginalSession)
	if err = gob.NewEncoder(buf).Encode(ses.Values); err != nil {
		return fmt.Errorf("failed to encode session: %w", err)
	}

	cortezaSession.Data = buf.Bytes()

	return s.store.UpsertAuthSession(ctx, cortezaSession)
}

func (s cortezaSessionStore) Delete(ctx context.Context, ses *sessions.Session) error {
	if len(ses.ID) > 0 {
		return s.store.DeleteAuthSessionByID(ctx, ses.ID)
	}

	return nil
}

var _ sessions.Store = &cortezaSessionStore{}
