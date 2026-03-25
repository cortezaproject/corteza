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
	"github.com/cortezaproject/corteza/server/pkg/apigw"
	a "github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/label"
	labelTypes "github.com/cortezaproject/corteza/server/pkg/label/types"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/store/adapters/api/cred_registry"
	"github.com/cortezaproject/corteza/server/store/adapters/api/drivers/google"
	restDriver "github.com/cortezaproject/corteza/server/store/adapters/api/drivers/rest"
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

	// connectionRunner is the minimal interface needed by health-check helpers.
	// Both *restAPIWrapper (via Run) and any future runner satisfy it.
	connectionRunner interface {
		Run(ctx context.Context, method, path string, payload []byte, headers map[string][]string) (statusCode int, outHeaders map[string][]string, rsp []byte, err error)
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

		if rsp.dalConnection != nil {
			res.Config.DalConnectionID = rsp.dalConnection.ID
		}
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

		// Load all active configured connections for the same source connection
		// so that registerOperations can merge them into a single function entry
		// with all configurationID options. Without this, each Enable call would
		// produce a duplicate entry containing only the newly-enabled CC.
		siblings, _, _ := store.SearchConfiguredConnections(ctx, svc.store, types.ConfiguredConnectionFilter{
			ConnectionID: res.ConnectionID,
			Status:       []string{"active"},
		})
		allCCs := make([]types.ConfiguredConnection, 0, len(siblings))
		for _, s := range siblings {
			allCCs = append(allCCs, *s)
		}
		svc.registerOperations(allCCs)
		svc.registerWebhookTriggers(*res)
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

func ensureGoogleCredential(ctx context.Context, s store.Storer, cc *types.ConfiguredConnection, conn *types.Connection) {
	if _, err := cred_registry.Default().Get(cc.ID); err != nil {
		var saJSON string

		// Prefer inline param
		for _, p := range cc.Config.Params {
			if p.Name == "serviceAccountJSON" {
				saJSON = p.Value
				break
			}
		}

		// Fallback: load from Credential store via CredentialID
		if saJSON == "" && cc.Config.CredentialID > 0 {
			if storedCred, loadErr := store.LookupCredentialByID(ctx, s, cc.Config.CredentialID); loadErr == nil {
				saJSON = storedCred.Credentials
			}
		}

		if saJSON == "" {
			return
		}

		// Private key PEM blocks may contain literal newlines when stored;
		// replace them with JSON escape sequences so Unmarshal succeeds.
		saJSON = strings.ReplaceAll(saJSON, "\n", `\n`)

		var sa struct {
			ClientEmail string `json:"client_email"`
			PrivateKey  string `json:"private_key"`
		}
		if err := json.Unmarshal([]byte(saJSON), &sa); err != nil || sa.ClientEmail == "" {
			return
		}

		var scopes []string
		switch {
		case strings.HasPrefix(conn.Handle, "google-calendar"):
			scopes = []string{"https://www.googleapis.com/auth/calendar"}
		case strings.HasPrefix(conn.Handle, "google-sheets"):
			scopes = []string{"https://www.googleapis.com/auth/spreadsheets"}
		case strings.HasPrefix(conn.Handle, "google-drive"):
			scopes = []string{"https://www.googleapis.com/auth/drive"}
		case strings.HasPrefix(conn.Handle, "google-tasks"):
			scopes = []string{"https://www.googleapis.com/auth/tasks"}
		}

		cred, err := cred_registry.NewCredential(cred_registry.CredentialConfig{
			ConnectionID:        cc.ID,
			AuthType:            "google_service_account",
			ServiceAccountEmail: sa.ClientEmail,
			PrivateKey:          sa.PrivateKey,
			Scopes:              scopes,
		})
		if err == nil && cred != nil {
			_ = cred_registry.Default().Store(cred)
		}
	}
}

// resolveExecutor returns the appropriate HTTP executor for the given configured
// connection. Google connectors use the google wrapper directly (bypassing DAL)
// because they have not yet been integrated as proper DAL connection types.
// All other connectors are looked up via the DAL service.
func resolveExecutor(
	ctx context.Context,
	cc *types.ConfiguredConnection,
	conn *types.Connection,
	baseURL string,
	dalConnectionID uint64,
) (func(ctx context.Context, method, path string, headers map[string][]string, payload []byte) (int, map[string][]string, []byte, error), error) {
	if strings.HasPrefix(conn.Handle, "google-") {
		// Ensure the service-account credential is cached in the registry
		// before the wrapper tries to use it.
		ensureGoogleCredential(ctx, DefaultStore, cc, conn)
		gw := google.NewWrapper(baseURL, cc.ID)
		return func(ctx context.Context, method, path string, headers map[string][]string, payload []byte) (int, map[string][]string, []byte, error) {
			return gw.Run(ctx, method, path, payload, headers)
		}, nil
	}

	cw := dal.Service().GetConnectionByID(dalConnectionID)
	if cw == nil {
		return nil, fmt.Errorf("DAL connection not found: %d", dalConnectionID)
	}
	return cw.Execute, nil
}

func (svc *configuredConnection) Check(ctx context.Context, ID uint64) (*types.ConfiguredConnectionCheckResult, error) {
	cc, err := loadConfiguredConnection(ctx, svc.store, ID)
	if err != nil {
		return nil, err
	}

	// Resolve templates from stored params — works for both draft and active CCs
	resolved := svc.resolveTemplates(&cc.Connection, cc.Config.Params)

	var runner connectionRunner
	if strings.HasPrefix(resolved.Handle, "google-") {
		ensureGoogleCredential(ctx, svc.store, cc, &cc.Connection)
		runner = google.NewWrapper(resolved.Service.BaseURL.Value, cc.ID)
	} else {
		runner, err = restDriver.RunnerFromConnection(resolved)
		if err != nil {
			return nil, fmt.Errorf("could not build connection runner: %w", err)
		}
	}

	result := &types.ConfiguredConnectionCheckResult{}
	result.Connectivity = svc.checkConnectivity(ctx, runner)
	result.Auth = svc.checkAuth(ctx, runner)

	if probe := resolved.Service.Probe; probe != nil {
		ps := svc.checkProbe(ctx, runner, probe)
		result.Probe = &ps
	}

	return result, nil
}

func (svc *configuredConnection) checkConnectivity(ctx context.Context, r connectionRunner) types.ConfiguredConnectionCheckStatus {
	statusCode, _, _, err := r.Run(ctx, "HEAD", "/", nil, nil)
	// If err != nil but statusCode > 0, we reached the server but got an HTTP error. Connectivity is OK!
	if err != nil && statusCode == 0 {
		return types.ConfiguredConnectionCheckStatus{OK: false, Message: err.Error()}
	}
	return types.ConfiguredConnectionCheckStatus{OK: true}
}

func (svc *configuredConnection) checkAuth(ctx context.Context, r connectionRunner) types.ConfiguredConnectionCheckStatus {
	statusCode, _, _, err := r.Run(ctx, "GET", "/", nil, nil)

	if statusCode == 401 || statusCode == 403 {
		return types.ConfiguredConnectionCheckStatus{OK: false, Message: fmt.Sprintf("authentication failed (HTTP %d)", statusCode)}
	}

	if statusCode == 404 {
		return types.ConfiguredConnectionCheckStatus{OK: false, Message: "calendar or resource not found (HTTP 404), or you lack permissions to view it"}
	}

	if err != nil && statusCode == 0 {
		return types.ConfiguredConnectionCheckStatus{OK: false, Message: "network error: " + err.Error()}
	}

	if err != nil {
		// some other HTTP error >= 400
		return types.ConfiguredConnectionCheckStatus{OK: false, Message: fmt.Sprintf("API returned HTTP %d", statusCode)}
	}

	return types.ConfiguredConnectionCheckStatus{OK: true}
}

func (svc *configuredConnection) checkProbe(ctx context.Context, r connectionRunner, probe *types.ConnectionProbe) types.ConfiguredConnectionCheckStatus {
	expected := probe.ExpectedStatus
	if expected == 0 {
		expected = 200
	}

	statusCode, _, _, err := r.Run(ctx, "GET", probe.Path.Value, nil, nil)
	if err != nil {
		return types.ConfiguredConnectionCheckStatus{OK: false, Message: err.Error()}
	}
	if statusCode != expected {
		return types.ConfiguredConnectionCheckStatus{OK: false, Message: fmt.Sprintf("probe returned HTTP %d, expected %d", statusCode, expected)}
	}
	return types.ConfiguredConnectionCheckStatus{OK: true}
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
	if !strings.HasPrefix(resolved.Handle, "google-") {
		rsp.dalConnection, err = svc.provisionDAL(ctx, resolved, conn)
		if err != nil {
			return rsp, err
		}
	} else {
		rsp.dalConnection = &types.DalConnection{ID: 0}
	}

	if err = svc.provisionWebhooks(ctx, resolved, conn); err != nil {
		return rsp, err
	}

	// modules, err := svc.provisionModules(ctx, resolved.Resources, conn.Config.NamespaceID, dalConn.ID)
	// err = svc.provisionAutomation(ctx, resolved.Operations, modules, dalConn.ID)

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
			"corteza/connection-id":         strconv.FormatUint(resolved.ID, 10),
			"corteza/connection-revision":   strconv.Itoa(resolved.Revision),
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
		for _, cc := range ccs {
			svc.registerWebhookTriggers(cc)
		}
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

			var cc *types.ConfiguredConnection
			for _, c := range ccs {
				if c.ID == configID {
					cc = &c
					break
				}
			}
			if cc == nil {
				return nil, fmt.Errorf("configured connection not found: %d", configID)
			}

			// Re-resolve to get correct BaseURL for this cc
			resolved := ConfiguredConnectionSvc().resolveTemplates(&cc.Connection, cc.Config.Params)
			baseURL := resolved.Service.BaseURL.Value

			// @todo driver selection does not belong in the service layer.
			// Once Google connectors are provisioned as proper DAL connections
			// (with the google wrapper as the underlying transport), this
			// branching can be removed and all connectors can go through
			// dal.Service().GetConnectionByID uniformly.
			execute, err := resolveExecutor(ctx, cc, &conn, baseURL, dalConnectionID)
			if err != nil {
				return nil, err
			}

			var respBody []byte
			for _, step := range op.Steps {
				if step.Type != "http" || step.HTTP == nil {
					continue
				}

				// Merge service-level param values into vars so that path/body
				// templates referencing them (e.g. {{calendarId}}) resolve
				// correctly. Runtime input takes precedence.
				vars := make(map[string]any, len(cc.Config.Params))
				for _, p := range cc.Config.Params {
					if len(p.Scope) > 0 && p.Scope[0] == "service" {
						vars[p.Name] = p.Value
					}
				}
				for k, v := range in.Dict() {
					vars[k] = v
				}

				path, headers, payload, err := buildHTTPRequest(*step.HTTP, in, vars)
				if err != nil {
					return nil, err
				}

				if baseURL != "" && strings.HasPrefix(path, baseURL) {
					path = strings.TrimPrefix(path, baseURL)
					if path != "" && !strings.HasPrefix(path, "/") && !strings.HasPrefix(path, "?") {
						path = "/" + path
					}
				}
				statusCode, outHeaders, body, err := execute(ctx, step.HTTP.Method, path, headers, payload)
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
func buildHTTPRequest(http types.ConnectionHTTPAction, in *expr.Vars, vars map[string]any) (path string, headers map[string][]string, payload []byte, err error) {
	if path, err = resolveTemplate(http.Path, vars); err != nil {
		return "", nil, nil, fmt.Errorf("path: %w", err)
	}

	u, err := url.Parse(path)
	if err != nil {
		return "", nil, nil, fmt.Errorf("invalid path %q: %w", path, err)
	}
	q := u.Query()
	for k, tpl := range http.QueryParams {
		v, e := resolveTemplate(tpl, vars)
		if e != nil {
			return "", nil, nil, fmt.Errorf("queryParam %q: %w", k, e)
		}
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	path = u.String()

	headers = make(map[string][]string, len(http.Headers))
	for k, tpl := range http.Headers {
		v, e := resolveTemplate(tpl, vars)
		if e != nil {
			return "", nil, nil, fmt.Errorf("header %q: %w", k, e)
		}
		headers[k] = []string{v}
	}

	if http.BodyTemplate.Value != "" {
		body, e := resolveTemplate(http.BodyTemplate, vars)
		if e != nil {
			return "", nil, nil, fmt.Errorf("bodyTemplate: %w", e)
		}
		payload = []byte(body)
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

// resolveTemplate interpolates values and omits missing optional placeholders
func resolveTemplate(tpl types.ConnectionTemplate, vars map[string]any) (string, error) {
	if tpl.Value == "" {
		return "", nil
	}

	// Build metadata lookup
	meta := make(map[string]types.ConnectionPlaceholder, len(tpl.Placeholders))
	for _, p := range tpl.Placeholders {
		meta[p.Name] = p
	}

	var resolveErr error
	result := placeholderRe.ReplaceAllStringFunc(tpl.Value, func(match string) string {
		if resolveErr != nil {
			return match
		}
		sub := placeholderRe.FindStringSubmatch(match)
		if len(sub) < 2 {
			return match
		}
		name := sub[1]

		if v, ok := vars[name]; ok {
				switch s := v.(type) {
				case string:
					return s
				default:
					b, _ := json.Marshal(s)
					return string(b)
				}
			}

		// Not provided — consult placeholder metadata
		if p, ok := meta[name]; ok {
			if !p.Required {
				return p.Default
			}
		}

		resolveErr = fmt.Errorf("missing required placeholder %q", name)
		return match
	})
	return result, resolveErr
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
	// Build a set of service-level param names so we can silently skip
	// operation-level inputs that duplicate them — those are resolved from
	// the configured connection's stored values, not from runtime user input.
	serviceParams := make(map[string]bool)
	for _, dp := range conn.DerivedParams {
		if len(dp.Scope) > 0 && dp.Scope[0] == "service" {
			serviceParams[dp.Name] = true
		}
	}

	inLookup := make(map[string]bool)

	// Explicit arguments — skip any that are already covered by service-level params
	for _, in := range op.Input {
		if serviceParams[in.Name] {
			continue
		}
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
						Label:       dp.Label,
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

func connectionWebhookResourceType(connectionID uint64) string {
	return fmt.Sprintf("corteza::system:connection-webhook/%d", connectionID)
}

// provisionWebhooks creates an routes and processors for webhook defs
func (svc *configuredConnection) provisionWebhooks(ctx context.Context, resolved *types.Connection, cc *types.ConfiguredConnection) error {
	var (
		routes  []*types.ApigwRoute
		filters []*types.ApigwFilter
		invoker = a.GetIdentityFromContext(ctx).Identity()
	)

	for _, res := range resolved.Resources {
		for _, wh := range res.Webhooks {
			route := &types.ApigwRoute{
				ID:       nextID(),
				Endpoint: fmt.Sprintf("/%s/%d/%d/%s", resolved.Handle, resolved.Revision, cc.ID, wh.Event),
				Method:   "POST",
				Enabled:  true,
				Meta: types.ApigwRouteMeta{
					Desc: fmt.Sprintf("Webhook: %s / %s", res.Handle, wh.Event),
					Labels: map[string]labelTypes.LabelValue{
						"corteza.connectionID": {Val: strconv.FormatUint(cc.ID, 10)},
						"corteza.webhookEvent": {Val: wh.Event},
						"corteza.resource":     {Val: res.Handle},
						"corteza.connection":   {Val: resolved.Handle},
					},
				},
				CreatedAt: *now(),
				CreatedBy: invoker,
			}
			routes = append(routes, route)

			mapping := make(map[string][]string, len(wh.Payload))
			for _, f := range wh.Payload {
				mapping[f.Name] = f.Selector
			}

			paramsMap := map[string]any{
				"connectionID":           strconv.FormatUint(cc.ConnectionID, 10),
				"configuredConnectionID": strconv.FormatUint(cc.ID, 10),
				"eventType":              wh.Event,
				"mapping":                mapping,
			}
			paramsJSON, _ := json.Marshal(paramsMap)

			f := &types.ApigwFilter{
				ID:        nextID(),
				Route:     route.ID,
				Ref:       "eventDispatch",
				Kind:      "processer",
				Enabled:   true,
				Weight:    100, // processer weight
				Params:    types.ApigwFilterParams{},
				CreatedAt: *now(),
				CreatedBy: invoker,
			}
			if err := json.Unmarshal(paramsJSON, &f.Params); err != nil {
				return fmt.Errorf("provisionWebhooks: could not build filter params: %w", err)
			}
			filters = append(filters, f)
		}
	}

	if len(routes) == 0 {
		return nil
	}

	if err := store.CreateApigwRoute(ctx, svc.store, routes...); err != nil {
		return fmt.Errorf("provisionWebhooks: could not create apigw routes: %w", err)
	}

	if err := store.CreateApigwFilter(ctx, svc.store, filters...); err != nil {
		return fmt.Errorf("provisionWebhooks: could not create apigw filters: %w", err)
	}

	for _, route := range routes {
		if err := apigw.Service().ReloadEndpoint(ctx, route.Method, route.Endpoint); err != nil {
			return fmt.Errorf("provisionWebhooks: could not reload apigw endpoint %s: %w", route.Endpoint, err)
		}
	}

	return nil
}

// registerWebhookTriggers registers Webhook-based triggers to the registry
func (svc *configuredConnection) registerWebhookTriggers(cc types.ConfiguredConnection) {
	existing := automationService.ConstructLibrary().Triggers()
	seen := make(map[string]bool, len(existing))
	for _, t := range existing {
		seen[t.ResourceType+"|"+t.EventType] = true
	}

	var tt []atypes.ConstructTrigger

	for _, res := range cc.Connection.Resources {
		for _, wh := range res.Webhooks {
			rt := connectionWebhookResourceType(cc.ConnectionID)
			key := rt + "|" + wh.Event
			if seen[key] {
				continue
			}
			seen[key] = true

			props := make([]atypes.ConstructTriggerProperty, 0, len(wh.Payload))
			for _, f := range wh.Payload {
				props = append(props, atypes.ConstructTriggerProperty{
					Name: f.Name,
					Type: f.Type,
				})
			}

			tt = append(tt, atypes.ConstructTrigger{
				ResourceType: rt,
				EventType:    wh.Event,
				Properties:   props,
				Meta: &atypes.ConstructTriggerMeta{
					Short: fmt.Sprintf("%s: %s", res.Handle, wh.Event),
				},
			})
		}
	}

	if len(tt) > 0 {
		automationService.ConstructLibrary().AddTriggers(tt...)
	}
}
