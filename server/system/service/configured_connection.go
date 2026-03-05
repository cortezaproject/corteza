package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	automationService "github.com/cortezaproject/corteza/server/automation/service"
	atypes "github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/actionlog"
	a "github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/label"
	labelTypes "github.com/cortezaproject/corteza/server/pkg/label/types"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	configuredConnection struct {
		actionlog     actionlog.Recorder
		dalConnection dalConMngmntSvc
		store         store.Storer
		ac            configuredConnectionAccessController
		connectionSvc *connection
	}

	configuredConnectionAccessController interface {
		CanSearchConfiguredConnections(ctx context.Context) bool

		CanCreateConfiguredConnection(context.Context) bool
		CanReadConfiguredConnection(context.Context, *types.ConfiguredConnection) bool
		CanDeleteConfiguredConnection(context.Context, *types.ConfiguredConnection) bool
	}

	dalConMngmntSvc interface {
		Create(ctx context.Context, new *types.DalConnection) (*types.DalConnection, error)
	}
)

func ConfiguredConnectionSvc() *configuredConnection {
	return &configuredConnection{
		actionlog:     DefaultActionlog,
		store:         DefaultStore,
		ac:            DefaultAccessControl,
		connectionSvc: DefaultConnection,
	}
}

func (svc *configuredConnection) WithDalConnection(s dalConMngmntSvc) *configuredConnection {
	svc.dalConnection = s
	return svc
}

func (svc *configuredConnection) FindByID(ctx context.Context, ID uint64) (res *types.ConfiguredConnection, err error) {
	var (
		aProps = &configuredConnectionActionProps{connection: &types.ConfiguredConnection{ID: ID}}
	)

	err = func() error {
		if res, err = loadConfiguredConnection(ctx, svc.store, ID); err != nil {
			return err
		}

		aProps.setConnection(res)

		if !svc.ac.CanReadConfiguredConnection(ctx, res) {
			return ConfiguredConnectionErrNotAllowedToRead()
		}

		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ConfiguredConnectionActionLookup, err)
}

func (svc *configuredConnection) Create(ctx context.Context, new *types.ConfiguredConnection) (res *types.ConfiguredConnection, err error) {
	var (
		aProps = &configuredConnectionActionProps{new: new}
	)

	err = func() (err error) {
		if !svc.ac.CanCreateConfiguredConnection(ctx) {
			return ConfiguredConnectionErrNotAllowedToCreate()
		}

		// Fetch and snapshot the connection definition
		var conn *types.Connection
		if conn, err = svc.connectionSvc.FindByID(ctx, new.ConnectionID); err != nil {
			return err
		}

		new.ID = nextID()
		new.Connection = *conn
		new.CreatedAt = *now()
		new.CreatedBy = a.GetIdentityFromContext(ctx).Identity()

		if new.Status == "" {
			new.Status = "draft"
		}

		if new.Labels == nil {
			new.Labels = make(map[string]labelTypes.LabelValue)
		}
		new.Labels["corteza/connection-id"] = labelTypes.LabelValue{Val: strconv.FormatUint(conn.ID, 10)}
		new.Labels["corteza/connection-revision"] = labelTypes.LabelValue{Val: strconv.Itoa(conn.Revision)}

		if err = store.CreateConfiguredConnection(ctx, svc.store, new); err != nil {
			return err
		}

		if err = label.Create(ctx, svc.store, new); err != nil {
			return err
		}

		return nil
	}()

	return new, svc.recordAction(ctx, aProps, ConfiguredConnectionActionCreate, err)
}

func (svc *configuredConnection) Update(ctx context.Context, upd *types.ConfiguredConnection) (res *types.ConfiguredConnection, err error) {
	var (
		uaProps = &configuredConnectionActionProps{update: upd}
	)

	err = func() (err error) {
		if res, err = loadConfiguredConnection(ctx, svc.store, upd.ID); err != nil {
			return
		}

		if res.Status != "draft" {
			return ConfiguredConnectionErrCannotUpdateInstalled()
		}

		uaProps.setUpdate(upd)

		// Test if stale (update has an older version of data)
		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return UserErrStaleData()
		}

		// Assign changed values
		res.Name = upd.Name
		res.UpdatedAt = now()
		res.Config = upd.Config

		if err = store.UpdateConfiguredConnection(ctx, svc.store, res); err != nil {
			return
		}

		if label.Changed(res.Labels, upd.Labels) {
			if err = label.Update(ctx, svc.store, upd); err != nil {
				return
			}

			res.Labels = upd.Labels
		}

		return
	}()

	return res, svc.recordAction(ctx, uaProps, ConfiguredConnectionActionUpdate, err)
}

func (svc *configuredConnection) DeleteByID(ctx context.Context, ID uint64) (err error) {
	return ConfiguredConnectionErrDeletionNotSupported()
}

func (svc *configuredConnection) Enable(ctx context.Context, ID uint64) (res *types.ConfiguredConnection, err error) {
	var (
		aProps = &configuredConnectionActionProps{connection: &types.ConfiguredConnection{ID: ID}}
	)

	err = func() (err error) {
		if res, err = loadConfiguredConnection(ctx, svc.store, ID); err != nil {
			return err
		}

		aProps.setConnection(res)

		if res.Status != "draft" {
			return errors.InvalidData("only draft connections can be enabled")
		}

		// Resolve templates and provision sub-systems
		resolved := svc.resolveTemplates(&res.Connection, res.Config.Params)
		rsp, err := svc.dispatch(ctx, resolved, res)
		if err != nil {
			return err
		}

		res.Config.DalConnectionID = rsp.dalConnection.ID
		res.Status = "active"

		if res.Labels == nil {
			res.Labels = make(map[string]labelTypes.LabelValue)
		}

		res.Labels["corteza/configured-connection-id"] = labelTypes.LabelValue{Val: strconv.FormatUint(res.ID, 10)}

		n := now()
		res.UpdatedAt = n
		res.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()

		if err = store.UpdateConfiguredConnection(ctx, svc.store, res); err != nil {
			return err
		}

		svc.registerOperations([]types.ConfiguredConnection{*res})
		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ConfiguredConnectionActionEnable, err)
}

func (svc *configuredConnection) Search(ctx context.Context, filter types.ConfiguredConnectionFilter) (set types.ConfiguredConnectionSet, f types.ConfiguredConnectionFilter, err error) {
	var (
		aProps = &configuredConnectionActionProps{filter: &filter}
	)

	err = func() error {
		if !svc.ac.CanSearchConfiguredConnections(ctx) {
			return ConfiguredConnectionErrNotAllowedToSearch()
		}

		set, f, err = store.SearchConfiguredConnections(ctx, svc.store, filter)
		return err
	}()

	return set, f, svc.recordAction(ctx, aProps, ConfiguredConnectionActionSearch, err)
}

func loadConfiguredConnection(ctx context.Context, s store.ConfiguredConnections, ID uint64) (res *types.ConfiguredConnection, err error) {
	if ID == 0 {
		return nil, errors.NotFound("connection not found")
	}

	if res, err = store.LookupConfiguredConnectionByID(ctx, s, ID); errors.IsNotFound(err) {
		err = errors.NotFound("connection not found")
	}

	return
}

// resolveTemplates substitutes all {{placeholder}} variables in a connection's
// Template fields with values from the connection's scoped params.
// Returns a deep copy of the connection with all templates resolved.
func (svc *configuredConnection) resolveTemplates(conn *types.Connection, params []types.ConfiguredConnectionParam) *types.Connection {
	// Build a quick lookup: "scope|name" → value
	lookup := make(map[string]string, len(params))
	for _, p := range params {
		key := joinScope(p.Scope) + "|" + p.Name
		lookup[key] = p.Value
	}

	resolve := func(tpl *types.ConnectionTemplate, scope []string) {
		scopeKey := joinScope(scope)
		tpl.Value = placeholderRe.ReplaceAllStringFunc(tpl.Value, func(match string) string {
			sub := placeholderRe.FindStringSubmatch(match)
			if len(sub) < 2 {
				return match
			}
			if v, ok := lookup[scopeKey+"|"+sub[1]]; ok {
				return v
			}
			return match
		})
	}

	resolveHTTPAction := func(op *types.ConnectionHTTPAction, scope []string) {
		if op == nil {
			return
		}
		resolve(&op.Path, scope)
		resolve(&op.BodyTemplate, scope)
		for k, h := range op.Headers {
			resolve(&h, scope)
			op.Headers[k] = h
		}
		for k, q := range op.QueryParams {
			resolve(&q, scope)
			op.QueryParams[k] = q
		}
	}

	// Deep-copy connection so we don't mutate the original
	out := *conn

	// Service-level
	resolve(&out.Service.BaseURL, []string{"service", "baseURL"})

	headers := make(map[string]types.ConnectionTemplate, len(out.Service.Headers))
	for k, h := range out.Service.Headers {
		resolve(&h, []string{"service", "headers"})
		headers[k] = h
	}
	out.Service.Headers = headers

	authParams := make(map[string]types.ConnectionTemplate, len(out.Service.Auth.Params))
	for k, p := range out.Service.Auth.Params {
		resolve(&p, []string{"service", "auth"})
		authParams[k] = p
	}
	out.Service.Auth.Params = authParams

	// Resources
	resources := make(types.ConnectionResources, len(out.Resources))
	copy(resources, out.Resources)
	for i := range resources {
		resolve(&resources[i].Endpoint, []string{"resources", resources[i].Handle, "endpoint"})

		resolveHTTPAction(resources[i].Operations.List, []string{"resources", resources[i].Handle, "standardOperations", "list"})
		resolveHTTPAction(resources[i].Operations.Read, []string{"resources", resources[i].Handle, "standardOperations", "read"})
		resolveHTTPAction(resources[i].Operations.Create, []string{"resources", resources[i].Handle, "standardOperations", "create"})
		resolveHTTPAction(resources[i].Operations.Update, []string{"resources", resources[i].Handle, "standardOperations", "update"})
		resolveHTTPAction(resources[i].Operations.Delete, []string{"resources", resources[i].Handle, "standardOperations", "delete"})
	}
	out.Resources = resources

	// Standard operations

	// Custom operations
	ops := make(types.ConnectionOperations, len(out.Operations))
	copy(ops, out.Operations)
	for i := range ops {
		for j := range ops[i].Steps {
			if ops[i].Steps[j].HTTP != nil {
				resolveHTTPAction(ops[i].Steps[j].HTTP, []string{"operations", ops[i].Handle})
			}
		}
	}
	out.Operations = ops

	return &out
}

// dispatch orchestrates sub-system provisioning from a resolved connection definition.
// TODO: implement Compose, Automation, and IG provisioning.
func (svc *configuredConnection) dispatch(ctx context.Context, resolved *types.Connection, conn *types.ConfiguredConnection) (rsp dispatchRsp, err error) {
	rsp.dalConnection, err = svc.provisionDAL(ctx, resolved, conn)
	if err != nil {
		return rsp, err
	}

	// modules, err := svc.provisionModules(ctx, resolved.Resources, conn.Config.NamespaceID, dalConn.ID)
	// err = svc.provisionAutomation(ctx, resolved.Operations, modules, dalConn.ID)
	// err = svc.provisionWebhooks(ctx, resolved.Resources, modules)

	return
}

func (svc *configuredConnection) provisionDAL(ctx context.Context, resolved *types.Connection, conn *types.ConfiguredConnection) (*types.DalConnection, error) {
	// Build DAL connection params
	params := map[string]any{
		"url": resolved.Service.BaseURL.Value,
	}

	for k, h := range resolved.Service.Headers {
		if params["headers"] == nil {
			params["headers"] = map[string]string{}
		}
		params["headers"].(map[string]string)[k] = h.Value
	}

	// For Google service accounts or OAuth2, we use the credential ID.
	// We'll pass it in params so the driver can retrieve it via cred_registry.
	if conn.Config.CredentialID > 0 {
		params["credentialID"] = conn.Config.CredentialID
	}

	if resolved.Service.Auth.Method != "" {
		authParams := make(map[string]any)
		for k, p := range resolved.Service.Auth.Params {
			authParams[k] = p.Value
		}

		params["auth"] = map[string]any{
			"method": resolved.Service.Auth.Method,
			"params": authParams,
		}
	}

	dalConn := &types.DalConnection{
		Handle: fmt.Sprintf("%s_%d", resolved.Handle, conn.ID),
		Type:   "corteza::system:dal-connection",
		Meta: types.DalConnectionMeta{
			Name: conn.Name,
		},
		Config: types.DalConnectionConfig{
			DAL: &types.DalConnectionConfigDAL{
				Type:   "corteza::dal:connection:rest",
				Params: params,
			},
		},
		Labels: map[string]string{
			"corteza/connector-id":          strconv.FormatUint(resolved.ID, 10),
			"corteza/connector-revision":    strconv.Itoa(resolved.Revision),
			"corteza/configured-connection": strconv.FormatUint(conn.ID, 10),
		},
	}

	// Create via dal connection service
	if svc.dalConnection == nil {
		return nil, fmt.Errorf("dalConnection service not injected")
	}

	return svc.dalConnection.Create(ctx, dalConn)
}

// RegisterAllOperations loads all active configured connections and registers
// their operations into the automation construct library.
// Called on boot from service.go.
func (svc *configuredConnection) RegisterAllOperations(ctx context.Context) {
	set, _, err := store.SearchConfiguredConnections(ctx, svc.store, types.ConfiguredConnectionFilter{
		Status: []string{"active"},
	})
	if err != nil {
		return
	}

	// Group configured connections by their source connection (connector)
	byConn := make(map[uint64][]types.ConfiguredConnection)
	for _, cc := range set {
		byConn[cc.ConnectionID] = append(byConn[cc.ConnectionID], *cc)
	}

	for _, ccs := range byConn {
		svc.registerOperations(ccs)
	}
}

// registerOperations converts each ConnectionOperation into a ConstructFunction
// and adds it to the automation construct library.
func (svc *configuredConnection) registerOperations(ccs []types.ConfiguredConnection) {
	if len(ccs) == 0 || len(ccs[0].Connection.Operations) == 0 {
		return
	}

	fns := make([]atypes.ConstructFunction, 0, len(ccs[0].Connection.Operations))
	for _, op := range ccs[0].Connection.Operations {
		fn := operationToFunction(ccs, op)
		fns = append(fns, fn)
	}

	automationService.ConstructLibrary().AddFunctions(fns...)
}

func operationToFunction(ccs []types.ConfiguredConnection, op types.ConnectionOperation) atypes.ConstructFunction {
	conn := ccs[0].Connection
	ref := fmt.Sprintf("conn_%d_%s", ccs[0].ConnectionID, op.Handle)

	// Build a lookup map: configurationID → dalConnectionID
	dalByConfig := make(map[uint64]uint64, len(ccs))
	for _, cc := range ccs {
		dalByConfig[cc.ID] = cc.Config.DalConnectionID
	}

	configParam := &atypes.Param{
		ArgumentName: "configurationID",
		Types:        []string{"ID"},
		Required:     true,
		Meta: &atypes.ParamMeta{
			Label:       "Configuration",
			Description: "Which configured connection to use",
		},
	}

	params := generateFunctionArguments(conn, op)
	segments := generateFunctionSegments(conn, op, params)
	results := generateFunctionResults(conn, op)

	// Add the parameter and segment for config ID
	params = append(atypes.ParamSet{configParam}, params...)

	input := generateSegmentInput(atypes.ParamSet{configParam})
	configOptions := make([]atypes.SelectItem, len(ccs))
	for i, cc := range ccs {
		configOptions[i] = atypes.SelectItem{
			Label: cc.Name,
			Value: strconv.FormatUint(cc.ID, 10),
		}
	}
	input[0].Input.Options = configOptions

	segments[0].Sections[0].Elements = append(input, segments[0].Sections[0].Elements...)

	var icon *atypes.NgAutomationIcon
	if conn.Meta.Icon != "" {
		icon = &atypes.NgAutomationIcon{Type: "name", Value: conn.Meta.Icon}
	}

	return atypes.ConstructFunction{
		Ref:    ref,
		Kind:   "function",
		Groups: []string{conn.Meta.Short},
		Meta: &atypes.ConstructFunctionMeta{
			Short:       op.Meta.Short,
			Description: op.Meta.Description,
			Icon:        icon,
		},
		Parameters: params,
		Results:    results,
		Segments:   segments,
		Labels: map[string]string{
			"connection": "step,workflow",
			op.Handle:    "step",
		},
		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			out = &expr.Vars{}

			// Resolve which configured connection to use
			var configID uint64
			if v, ok := in.Dict()["configurationID"]; ok {
				switch id := v.(type) {
				case uint64:
					configID = id
				case string:
					configID, _ = strconv.ParseUint(id, 10, 64)
				}
			}
			dalConnectionID, ok := dalByConfig[configID]
			if !ok {
				return nil, fmt.Errorf("unknown configurationID: %d", configID)
			}

			cw := dal.Service().GetConnectionByID(dalConnectionID)
			if cw == nil {
				return nil, fmt.Errorf("service DAL connection not found: %d", dalConnectionID)
			}

			resolveTemplate := makeTemplateResolver(in.Dict())

			var respBody []byte
			for _, step := range op.Steps {
				if step.Type != "http" || step.HTTP == nil {
					continue
				}

				path, headers, payload, err := buildHTTPRequest(*step.HTTP, in, resolveTemplate)
				if err != nil {
					return nil, err
				}

				statusCode, outHeaders, body, err := cw.Execute(ctx, step.HTTP.Method, path, headers, payload)
				if err != nil {
					return nil, fmt.Errorf("operation execution failed: %w", err)
				}

				if err = checkHTTPResponse(statusCode, outHeaders, body); err != nil {
					return nil, err
				}

				respBody = body

				// @todo fix up when we support multi-step operations
				break
			}

			if len(respBody) == 0 || len(op.Output) == 0 {
				return out, nil
			}

			var respData any
			if err := json.Unmarshal(respBody, &respData); err != nil {
				return out, fmt.Errorf("failed to parse response JSON: %w", err)
			}

			for _, outField := range op.Output {
				if len(outField.Selector) == 0 {
					_ = out.Set(outField.Name, extractByPath(respData, []string{outField.Name}))
					continue
				}
				_ = out.Set(outField.Name, extractByPath(respData, outField.Selector))
			}

			return
		},
	}
}

// buildHTTPRequest resolves templates in path, query params, headers, and body
// for a single HTTP operation, returning the final path, headers, and payload.
func buildHTTPRequest(http types.ConnectionHTTPAction, in *expr.Vars, resolve func(string) string) (path string, headers map[string][]string, payload []byte, err error) {
	path = resolve(http.Path.Value)

	u, err := url.Parse(path)
	if err != nil {
		return "", nil, nil, fmt.Errorf("invalid path %q: %w", path, err)
	}
	q := u.Query()
	for k, tpl := range http.QueryParams {
		q.Set(k, resolve(tpl.Value))
	}
	u.RawQuery = q.Encode()
	path = u.String()

	headers = make(map[string][]string, len(http.Headers))
	for k, tpl := range http.Headers {
		headers[k] = []string{resolve(tpl.Value)}
	}

	if http.BodyTemplate.Value != "" {
		payload = []byte(resolve(http.BodyTemplate.Value))
	} else if in.Len() > 0 && (http.Method == "POST" || http.Method == "PUT" || http.Method == "PATCH") {
		payload, _ = json.Marshal(in.Dict())
	}

	return
}

// checkHTTPResponse returns an error if the response indicates a failure,
// combining the HTTP status code, any standard error headers, and the body.
func checkHTTPResponse(statusCode int, headers map[string][]string, body []byte) error {
	if statusCode < 400 {
		return nil
	}

	var headerErr string
	for _, h := range []string{"X-Error", "X-Error-Message", "X-Api-Error"} {
		if vals := headers[h]; len(vals) > 0 && vals[0] != "" {
			headerErr = vals[0]
			break
		}
	}

	switch {
	case headerErr != "" && len(body) > 0:
		return fmt.Errorf("request failed (status %d, %s): %s", statusCode, headerErr, body)
	case headerErr != "":
		return fmt.Errorf("request failed (status %d): %s", statusCode, headerErr)
	case len(body) > 0:
		return fmt.Errorf("request failed (status %d): %s", statusCode, body)
	default:
		return fmt.Errorf("request failed (status %d)", statusCode)
	}
}

// makeTemplateResolver builds a single-pass replacer from a map of variables.
// Keys are wrapped in {{...}} delimiters.
func makeTemplateResolver(vars map[string]any) func(string) string {
	pairs := make([]string, 0, len(vars)*2)
	for k, v := range vars {
		pairs = append(pairs, "{{"+k+"}}", fmt.Sprintf("%v", v))
	}
	r := strings.NewReplacer(pairs...)
	return r.Replace
}

// extractByPath walks a parsed JSON tree (`any`) using a slice of string keys.
func extractByPath(data any, path []string) any {
	current := data
	for _, p := range path {
		switch v := current.(type) {
		case map[string]any:
			current = v[p]
		default:
			return nil // Cannot traverse further
		}
	}
	return current
}

func generateFunctionArguments(conn types.Connection, op types.ConnectionOperation) (params atypes.ParamSet) {
	inLookup := make(map[string]bool)

	// Explicit arguments
	for _, in := range op.Input {
		inLookup[in.Name] = true
		params = append(params, &atypes.Param{
			ArgumentName: in.Name,
			// @todo improve type mapping/determination; we might need to enforce this when defining the connection
			Types:    []string{in.Type},
			Required: in.Required,
			Meta: &atypes.ParamMeta{
				Label: in.Name,
			},
		})
	}

	// Implicit from derived parameters
	for _, dp := range conn.DerivedParams {
		if len(dp.Scope) == 2 && dp.Scope[0] == "operations" && dp.Scope[1] == op.Handle {
			if !inLookup[dp.Name] {
				inLookup[dp.Name] = true
				params = append(params, &atypes.Param{
					ArgumentName: dp.Name,
					Types:        []string{dp.Type},
					Required:     dp.Required,
					Meta: &atypes.ParamMeta{
						Label:       dp.Name,
						Description: dp.Description,
					},
				})
			}
		}
	}

	return
}

func generateFunctionResults(conn types.Connection, op types.ConnectionOperation) (results atypes.ParamSet) {
	for _, out := range op.Output {
		results = append(results, &atypes.Param{
			ArgumentName: out.Name,
			// @todo improve type mapping/determination; we might need to enforce this when defining the connection
			Types: []string{out.Type},
			Meta: &atypes.ParamMeta{
				Label: out.Name,
			},
		})
	}

	return
}

func generateFunctionSegments(conn types.Connection, op types.ConnectionOperation, parameters atypes.ParamSet) (out []atypes.ConstructSegment) {
	var elements []atypes.SectionElement

	elements = generateSegmentInput(parameters)

	return []atypes.ConstructSegment{
		{
			Sections: []atypes.ConstructSection{
				{
					Elements: elements,
				},
			},
		},
	}
}

func generateSegmentInput(parameters atypes.ParamSet) (elements []atypes.SectionElement) {
	for _, p := range parameters {
		inputType := "string"
		if len(p.Types) > 0 {
			inputType = p.Types[0]
		}

		elements = append(elements, atypes.SectionElement{
			Input: atypes.SectionElementInput{
				Type:     inputType,
				Label:    p.ArgumentName,
				Argument: p.ArgumentName,
			},
		})
	}

	return
}
