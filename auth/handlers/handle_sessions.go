package handlers

import (
	"crypto/md5"
	"fmt"
	"sort"
	"time"

	"github.com/cortezaproject/corteza-server/auth/request"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	userSessions []userSession
	userSession  struct {
		id string
		// this will store hashed ID value of the session ID + auth secret
		ID             string
		Current        bool
		ExpiresAt      time.Time
		CreatedAt      time.Time
		Expired        bool
		ExpiresIn      int
		RemoteAddr     string
		SameRemoteAddr bool
		UserAgent      string
		SameUserAgent  bool
	}
)

// Will sort sessions - current first, then the rest in order of creation
var _ sort.Interface = userSessions{}

func (set userSessions) Len() int      { return len(set) }
func (set userSessions) Swap(i, j int) { set[i], set[j] = set[j], set[i] }
func (set userSessions) Less(i, j int) bool {
	switch {
	case set[i].Current:
		return true
	default:
		return set[i].CreatedAt.After(set[j].CreatedAt)
	}

}

func (h *AuthHandlers) sessionsView(req *request.AuthReq) error {
	req.Template = TmplSessions

	ss, err := h.getSessions(req)
	if err != nil {
		return err
	}

	sort.Sort(ss)
	req.Data["sessions"] = ss

	return nil
}

func (h *AuthHandlers) sessionsProc(req *request.AuthReq) error {
	ss, err := h.getSessions(req)
	if err != nil {
		return err
	}

	switch {
	case len(req.Request.PostFormValue("delete-all-but-current")) > 0:
		for _, s := range ss {
			if s.Current {
				continue
			}

			if err = h.SessionManager.DeleteByID(req.Context(), s.id); err != nil {
				return err
			}
		}

		t := translator(req, "auth")
		req.NewAlerts = append(req.NewAlerts, request.Alert{
			Type: "primary",
			Text: t("sessions.alerts.delete-sessions-but-current"),
		})
	case len(req.Request.PostFormValue("delete")) > 0:
		for _, s := range ss {
			if s.ID != req.Request.PostFormValue("delete") {
				continue
			}

			if err = h.SessionManager.DeleteByID(req.Context(), s.id); err != nil {
				return err
			}
		}

		t := translator(req, "auth")
		req.NewAlerts = append(req.NewAlerts, request.Alert{
			Type: "primary",
			Text: t("sessions.alerts.session-deleted"),
		})
	}

	req.RedirectTo = GetLinks().Sessions
	return nil
}

func (h *AuthHandlers) getSessions(req *request.AuthReq) (ss userSessions, err error) {
	var set []*types.AuthSession
	if set, err = h.SessionManager.Search(req.Context(), req.AuthUser.User.ID); err != nil {
		return
	} else {
		ss = make(userSessions, len(set))
		for i := range set {
			isCurrent := req.Session.ID == set[i].ID
			ss[i] = userSession{
				id:             set[i].ID,
				ID:             fmt.Sprintf("%x", md5.Sum([]byte(set[i].ID+h.Opt.Secret))),
				Current:        isCurrent,
				CreatedAt:      set[i].CreatedAt,
				ExpiresAt:      set[i].ExpiresAt,
				Expired:        set[i].ExpiresAt.Before(time.Now()),
				ExpiresIn:      int(set[i].ExpiresAt.Sub(time.Now()).Hours() / 24),
				RemoteAddr:     set[i].RemoteAddr,
				SameRemoteAddr: set[i].RemoteAddr == req.Request.RemoteAddr,
				UserAgent:      set[i].UserAgent,
				SameUserAgent:  set[i].UserAgent == req.Request.UserAgent(),
			}
		}
	}

	return
}
