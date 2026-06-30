package system

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	composeTypes "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/system/types"
	"github.com/crusttech/human/server/tests/helpers"
	"github.com/davecgh/go-spew/spew"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

func Test_dml_full_flow(t *testing.T) {
	// @todo this should not be hard coded
	dsn := "postgres://corteza:corteza@127.0.0.1:3402/action_log_testing?sslmode=disable"

	h := newHelper(t)
	helpers.AllowMe(h, types.ComponentRbacResource(), "dal-connections.search")
	helpers.AllowMe(h, types.ComponentRbacResource(), "dal-connection.create")
	helpers.AllowMe(h, composeTypes.ComponentRbacResource(), "namespace.create")
	helpers.AllowMe(h, composeTypes.NamespaceRbacResource(0), "module.create", "modules.search", "read")
	helpers.AllowMe(h, composeTypes.ModuleRbacResource(0, 0), "read", "records.search", "record.create")

	// 1. create connection
	connBody, _ := json.Marshal(types.DmlConnectionInput{
		Handle: "dml_test_conn",
		Params: &types.DmlConnectionParams{
			Type:   "corteza::dal:connection:dsn",
			Params: map[string]any{"dsn": dsn},
		},
	})
	connResult := h.apiInit().
		Post("/dml/connections").
		Body(string(connBody)).
		Header("Content-Type", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.Present("$.response.connectionID")).
		End()

	var connResp struct {
		Response types.DmlConnection `json:"response"`
	}
	h.a.NoError(json.NewDecoder(connResult.Response.Body).Decode(&connResp))
	conn := connResp.Response
	h.a.NotZero(conn.ID)

	// 2. get models
	modelsResult := h.apiInit().
		Getf("/dml/connections/%d/models", conn.ID).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.Present("$.response.set")).
		End()

	var modelsResp struct {
		Response struct {
			Set []*types.DmlModel `json:"set"`
		} `json:"response"`
	}
	h.a.NoError(json.NewDecoder(modelsResult.Response.Body).Decode(&modelsResp))
	for _, m := range modelsResp.Response.Set {
		t.Logf("model: %s  attrs: %d", m.Ident, len(m.Attributes))
		spew.Dump(m)
	}
	h.a.NotEmpty(modelsResp.Response.Set)

	// 3. create one mapping per discovered model
	var mappings []types.DmlMapping
	for _, model := range modelsResp.Response.Set {
		var columns types.DmlColumnMapSet
		for _, attr := range model.Attributes {
			columns = append(columns, &types.DmlColumnMap{
				SourceIdent: attr.Ident,
				FieldName:   attr.Ident,
				FieldKind:   "String",
			})
		}
		mpBody, _ := json.Marshal(struct {
			SourceIdent     string                `json:"sourceIdent"`
			ModuleHandle    string                `json:"moduleHandle"`
			ModuleName      string                `json:"moduleName"`
			NamespaceHandle string                `json:"namespaceHandle"`
			Columns         types.DmlColumnMapSet `json:"columns"`
		}{
			SourceIdent:     model.Ident,
			ModuleHandle:    model.Ident,
			ModuleName:      model.Label,
			NamespaceHandle: "dml_test",
			Columns:         columns,
		})
		mpResult := h.apiInit().
			Postf("/dml/connections/%d/mapping", conn.ID).
			Header("Accept", "application/json").
			Header("Content-Type", "application/json").
			Body(string(mpBody)).
			Expect(t).
			Status(http.StatusOK).
			Assert(jsonpath.Present("$.response.mappingID")).
			End()

		var mpResp struct {
			Response types.DmlMapping `json:"response"`
		}
		h.a.NoError(json.NewDecoder(mpResult.Response.Body).Decode(&mpResp))
		mp := mpResp.Response
		h.a.NotZero(mp.ID)
		h.a.Equal(model.Ident, mp.SourceIdent)
		mappings = append(mappings, mp)
	}

	// 4. import first mapping — applies schema + pulls data
	mp := mappings[0]
	importBody, _ := json.Marshal(struct {
		Method types.DmlImportMethod `json:"method"`
	}{Method: types.DmlImportMethodBackground})
	importResult := h.apiInit().
		Postf("/dml/connections/%d/mapping/%d/import", conn.ID, mp.ID).
		Header("Accept", "application/json").
		Header("Content-Type", "application/json").
		Body(string(importBody)).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.Present("$.response.runID")).
		End()

	var runResp struct {
		Response types.DmlImportRun `json:"response"`
	}
	h.a.NoError(json.NewDecoder(importResult.Response.Body).Decode(&runResp))
	run := runResp.Response
	h.a.NotZero(run.ID)

	for run.Status == "pending" || run.Status == "running" {
		time.Sleep(2 * time.Second)

		pollResult := h.apiInit().
			Getf("/dml/connections/%d/mapping/%d/import/%d", conn.ID, mp.ID, run.ID).
			Header("Accept", "application/json").
			Expect(t).
			Status(http.StatusOK).
			Assert(jsonpath.Present("$.response.runID")).
			End()

		var pollResp struct {
			Response types.DmlImportRun `json:"response"`
		}
		h.a.NoError(json.NewDecoder(pollResult.Response.Body).Decode(&pollResp))
		run = pollResp.Response
	}

	h.apiInit().
		Getf("/dml/connections/%d/mapping/%d/import/%d", conn.ID, mp.ID, run.ID).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.Equal("$.response.status", run.Status)).
		End()

	spew.Dump(run)
	h.a.Equal("completed", run.Status)
}
