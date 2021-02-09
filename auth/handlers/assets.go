package handlers

import (
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/cortezaproject/corteza-server/pkg/version"
	"github.com/goware/statik/fs"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type (
	TemplateExecutor interface {
		ExecuteTemplate(io.Writer, string, interface{}) error
	}
)

const (
	tmplRoot                     = "/templates"
	TmplAuthorizedClients        = "authorized-clients.html.tpl"
	TmplChangePassword           = "change-password.html.tpl"
	TmplLogin                    = "login.html.tpl"
	TmplLogout                   = "logout.html.tpl"
	TmplOAuth2AuthorizeClient    = "oauth2-authorize-client.html.tpl"
	TmplRequestPasswordReset     = "request-password-reset.html.tpl"
	TmplPasswordResetRequested   = "password-reset-requested.html.tpl"
	TmplResetPassword            = "reset-password.html.tpl"
	TmplProfile                  = "profile.html.tpl"
	TmplSessions                 = "sessions.html.tpl"
	TmplSignup                   = "signup.html.tpl"
	TmplPendingEmailConfirmation = "pending-email-confirmation.html.tpl"
	TmplInternalError            = "error-internal.html.tpl"
)

//var (
//	Templates = []string{
//		TmplAuthorizedClients,
//		TmplChangePassword,
//		TmplLogin,
//		TmplLogout,
//		TmplOAuth2AuthorizeClient,
//		TmplRequestPasswordReset,
//		TmplPasswordResetRequested,
//		TmplResetPassword,
//		TmplProfile,
//		TmplSessions,
//		TmplSignup,
//		TmplPendingEmailConfirmation,
//		TmplInternalError,
//	}
//)

func TemplateBase() *template.Template {
	return template.New("").
		Funcs(sprig.FuncMap()).
		Funcs(template.FuncMap{
			"version":   func() string { return version.Version },
			"buildtime": func() string { return version.BuildTime },
			"links":     GetLinks,
		})
}

// EmbeddedTemplates returns embedded templates.
//
// @todo migrate to go:embed as soon as we can
func EmbeddedTemplates(t *template.Template) (tpl *template.Template, err error) {
	var (
		f       http.File
		tplBody []byte
		sfs     http.FileSystem
	)

	if sfs, err = fs.New(embedded_assets); err != nil {
		return
	}

	if t, err = t.Clone(); err != nil {
		return
	}

	return t, fs.Walk(sfs, tmplRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// templates are always referenced with path, relative to template root
		var tplName = path[len(tmplRoot)+1:]

		if f, err = sfs.Open(path); err != nil {
			return fmt.Errorf("could not open %s: %w", path, err)
		} else if tplBody, err = ioutil.ReadAll(f); err != nil {
			return fmt.Errorf("could not read %s: %w", path, err)
		} else if t, err = t.New(tplName).Parse(string(tplBody)); err != nil {
			return fmt.Errorf("could not parse %s: %w", path, err)
		}

		return nil
	})
}

func EmbeddedPublicAssets() http.HandlerFunc {
	sfs, err := fs.New(embedded_assets)
	if err != nil {
		return func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, fmt.Sprintf("could not read embeded filesystem: %v", err.Error()), http.StatusInternalServerError)
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/auth/assets", http.FileServer(sfs)).ServeHTTP(w, r)
	}
}

func DirectPublicAssets(prefix, root string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix(prefix, http.FileServer(http.Dir(root))).ServeHTTP(w, r)
	}
}
