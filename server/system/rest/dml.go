package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
)

// DmlController exposes the DML feature over HTTP.
// Routes are mounted directly (no codegen required).
type DmlController struct{}

func NewDmlController() *DmlController { return &DmlController{} }

func (c *DmlController) MountRoutes(r chi.Router) {
	r.Route("/dml", func(r chi.Router) {
		// Reject anonymous callers — the surrounding HttpTokenValidator only
		// enforces token validity, it lets requests with no token through as
		// anonymous. requireDmlRead then gates every route on the DAL-connection
		// search permission so external connections/schemas are not enumerable
		// by any authenticated user. Mutating routes add manage rights below.
		r.Use(requireAuthenticated)
		r.Use(requireDmlRead)

		r.Get("/connections", c.listConnections)
		r.Get("/connections/{connectionID}", c.getConnection)
		r.Get("/connections/{connectionID}/models", c.getConnectionModels)

		r.Post("/connections/{connectionID}/mapping", c.generateMapping)
		r.Get("/mappings/{mappingID}", c.getMapping)
		r.Put("/mappings/{mappingID}", c.updateMapping)

		r.Post("/mappings/{mappingID}/apply", c.applyMapping)
		r.Post("/mappings/{mappingID}/import", c.startImport)
		r.Get("/runs/{runID}", c.getRun)
	})
}

// requireAuthenticated blocks anonymous requests.
func requireAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !auth.GetIdentityFromContext(r.Context()).Valid() {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// requireDmlRead gates read access to DML routes on the DAL-connection search
// permission — DML exposes external connection metadata and introspected
// schemas, which should not be enumerable by every authenticated user.
func requireDmlRead(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !service.DefaultAccessControl.CanSearchDalConnections(r.Context()) {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// canManage gates mutating DML operations (generate mapping, update, apply,
// import). Reuse the DAL-connection create permission: orchestrating external
// connections + provisioning modules from them is at least as privileged.
//
// @todo replace with a first-class corteza::system:dml-* RBAC resource
//
//	(read/manage) once it's added to the codegen rbac blocks.
func canManage(r *http.Request) bool {
	return service.DefaultAccessControl.CanCreateDalConnection(r.Context())
}

func (c *DmlController) listConnections(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	f := types.DmlConnectionFilter{
		Handle: q.Get("handle"),
		Type:   q.Get("type"),
	}
	if ids, ok := q["connectionID"]; ok {
		f.ConnectionID = ids
	}

	set, err := service.DefaultDmlConnection.Find(r.Context(), f)
	if err != nil {
		dmlError(w, err)
		return
	}
	dmlJSON(w, http.StatusOK, map[string]interface{}{"set": set})
}

func (c *DmlController) getConnection(w http.ResponseWriter, r *http.Request) {
	id, err := parseUint64(chi.URLParam(r, "connectionID"))
	if err != nil {
		http.Error(w, "bad connectionID", http.StatusBadRequest)
		return
	}
	conn, err := service.DefaultDmlConnection.FindByID(r.Context(), id)
	if err != nil {
		dmlError(w, err)
		return
	}
	dmlJSON(w, http.StatusOK, conn)
}

func (c *DmlController) getConnectionModels(w http.ResponseWriter, r *http.Request) {
	id, err := parseUint64(chi.URLParam(r, "connectionID"))
	if err != nil {
		http.Error(w, "bad connectionID", http.StatusBadRequest)
		return
	}
	models, err := service.DefaultDmlConnection.FindModels(r.Context(), id)
	if err != nil {
		dmlError(w, err)
		return
	}
	dmlJSON(w, http.StatusOK, map[string]interface{}{"set": models})
}

func (c *DmlController) generateMapping(w http.ResponseWriter, r *http.Request) {
	if !canManage(r) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	id, err := parseUint64(chi.URLParam(r, "connectionID"))
	if err != nil {
		http.Error(w, "bad connectionID", http.StatusBadRequest)
		return
	}
	mp, err := service.DefaultDmlMapping.GenerateMapping(r.Context(), id)
	if err != nil {
		dmlError(w, err)
		return
	}
	dmlJSON(w, http.StatusCreated, mp)
}

func (c *DmlController) getMapping(w http.ResponseWriter, r *http.Request) {
	id, err := parseUint64(chi.URLParam(r, "mappingID"))
	if err != nil {
		http.Error(w, "bad mappingID", http.StatusBadRequest)
		return
	}
	mp, err := service.DefaultDmlMapping.FindByID(r.Context(), id)
	if err != nil {
		dmlError(w, err)
		return
	}
	dmlJSON(w, http.StatusOK, mp)
}

func (c *DmlController) updateMapping(w http.ResponseWriter, r *http.Request) {
	if !canManage(r) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	id, err := parseUint64(chi.URLParam(r, "mappingID"))
	if err != nil {
		http.Error(w, "bad mappingID", http.StatusBadRequest)
		return
	}
	var mp types.DmlMapping
	if err := json.NewDecoder(r.Body).Decode(&mp); err != nil {
		http.Error(w, "bad body: "+err.Error(), http.StatusBadRequest)
		return
	}
	mp.ID = id
	updated, err := service.DefaultDmlMapping.Update(r.Context(), &mp)
	if err != nil {
		dmlError(w, err)
		return
	}
	dmlJSON(w, http.StatusOK, updated)
}

func (c *DmlController) applyMapping(w http.ResponseWriter, r *http.Request) {
	if !canManage(r) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	id, err := parseUint64(chi.URLParam(r, "mappingID"))
	if err != nil {
		http.Error(w, "bad mappingID", http.StatusBadRequest)
		return
	}
	if err := service.DefaultDmlApplier.Apply(r.Context(), id); err != nil {
		dmlError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (c *DmlController) startImport(w http.ResponseWriter, r *http.Request) {
	if !canManage(r) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	id, err := parseUint64(chi.URLParam(r, "mappingID"))
	if err != nil {
		http.Error(w, "bad mappingID", http.StatusBadRequest)
		return
	}
	run, err := service.DefaultDmlImporter.RunImport(r.Context(), id)
	if err != nil {
		dmlError(w, err)
		return
	}
	dmlJSON(w, http.StatusAccepted, run)
}

func (c *DmlController) getRun(w http.ResponseWriter, r *http.Request) {
	id, err := parseUint64(chi.URLParam(r, "runID"))
	if err != nil {
		http.Error(w, "bad runID", http.StatusBadRequest)
		return
	}
	run, err := service.DefaultDmlImporter.GetRun(r.Context(), id)
	if err != nil {
		dmlError(w, err)
		return
	}
	dmlJSON(w, http.StatusOK, run)
}

func dmlJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func dmlError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func parseUint64(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}
