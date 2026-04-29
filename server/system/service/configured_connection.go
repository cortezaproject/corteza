package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	automationService "github.com/crusttech/human/server/automation/service"
	atypes "github.com/crusttech/human/server/automation/types"
	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/pkg/apigw"
	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/dal"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/expr"
	"github.com/crusttech/human/server/pkg/label"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/store/adapters/api/cred_registry"
	"github.com/crusttech/human/server/store/adapters/api/drivers/google"
	restDriver "github.com/crusttech/human/server/store/adapters/api/drivers/rest"
	"github.com/crusttech/human/server/system/types"
	"go.uber.org/zap"
)

type (
	configuredConnection struct {
		actionlog     actionlog.Recorder
		dalConnection dalConMngmntSvc
		store         store.Storer
		ac            configuredConnectionAccessController
		connectionSvc *connection
		logger        *zap.Logger
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
		logger:        DefaultLogger.Named("configured-connection"),
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

		// temp: always use the live connection definition instead of the snapshot
		if err = svc.liveConnection(ctx, res); err != nil {
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
		new.Labels["human/connection-id"] = labelTypes.LabelValue{Val: strconv.FormatUint(conn.ID, 10)}
		new.Labels["human/connection-revision"] = labelTypes.LabelValue{Val: strconv.Itoa(conn.Revision)}

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

		// temp: always use the live connection definition instead of the snapshot
		if err = svc.liveConnection(ctx, res); err != nil {
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

		res.Labels["human/configured-connection-id"] = labelTypes.LabelValue{Val: strconv.FormatUint(res.ID, 10)}

		n := now()
		res.UpdatedAt = n
		res.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()

		if err = store.UpdateConfiguredConnection(ctx, svc.store, res); err != nil {
			return err
		}

		// Fetch and cache Google resource discovery (spreadsheets + tabs).
		// Non-fatal: failures are logged and silently skipped.
		_ = svc.syncGoogleDiscovery(ctx, res)

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

		conn, connErr := loadConnection(ctx, svc.store, res.ConnectionID)
		if connErr == nil {
			svc.registerOperations(conn, allCCs)
		}
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
		if err != nil {
			return err
		}

		// temp: always use the live connection definition instead of the snapshot
		for _, cc := range set {
			_ = svc.liveConnection(ctx, cc)
		}
		return nil
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

		var sa struct {
			ClientEmail string `json:"client_email"`
			PrivateKey  string `json:"private_key"`
		}

		if err := json.Unmarshal([]byte(saJSON), &sa); err != nil {
			// Private key PEM blocks may contain literal newlines when stored;
			// replacing them with JSON escape strings helps Unmarshal succeed
			// for invalid single-line strings.
			fallbackJSON := strings.ReplaceAll(saJSON, "\n", `\n`)
			if err2 := json.Unmarshal([]byte(fallbackJSON), &sa); err2 != nil {
				return
			}
		}

		if sa.ClientEmail == "" {
			return
		}

		var scopes []string
		baseURL := conn.Service.BaseURL.Value
		switch {
		case strings.Contains(baseURL, "googleapis.com/calendar"):
			scopes = []string{"https://www.googleapis.com/auth/calendar"}
		case strings.Contains(baseURL, "sheets.googleapis.com"):
			scopes = []string{"https://www.googleapis.com/auth/spreadsheets", "https://www.googleapis.com/auth/drive.readonly"}
		case strings.Contains(baseURL, "googleapis.com/drive"):
			scopes = []string{"https://www.googleapis.com/auth/drive"}
		case strings.Contains(baseURL, "tasks.googleapis.com"):
			scopes = []string{"https://www.googleapis.com/auth/tasks"}
		case strings.Contains(baseURL, "gmail.googleapis.com"):
			scopes = []string{"https://mail.google.com/"}
		}

		var subject string
		for _, p := range cc.Config.Params {
			if p.Name == "dwdSubject" {
				subject = p.Value
				break
			}
		}

		cred, err := cred_registry.NewCredential(cred_registry.CredentialConfig{
			ConnectionID:        cc.ID,
			AuthType:            "google_service_account",
			ServiceAccountEmail: sa.ClientEmail,
			PrivateKey:          sa.PrivateKey,
			Scopes:              scopes,
			Subject:             subject,
		})
		if err == nil && cred != nil {
			_ = cred_registry.Default().Store(cred)
		}
	}
}

// resolveExecutor returns the appropriate HTTP executor for the given configured
// connection. Connections targeting googleapis.com use the Google wrapper
// (auth + transport). All other connectors are routed through the DAL service.
func resolveExecutor(
	ctx context.Context,
	cc *types.ConfiguredConnection,
	conn *types.Connection,
	baseURL string,
	dalConnectionID uint64,
) (func(ctx context.Context, method, path string, headers map[string][]string, payload []byte) (int, map[string][]string, []byte, error), error) {
	if strings.Contains(baseURL, "googleapis.com") {
		// Google APIs require OAuth2 token injection via the Google wrapper.
		// This check is URL-based and works regardless of the user-defined
		// connection handle or whether a gsheets DAL connection was provisioned.
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

	// temp: always use the live connection definition instead of the snapshot
	if err = svc.liveConnection(ctx, cc); err != nil {
		return nil, err
	}

	// Resolve templates from stored params — works for both draft and active CCs
	resolved := svc.resolveTemplates(&cc.Connection, cc.Config.Params)

	var runner connectionRunner
	if strings.Contains(resolved.Service.BaseURL.Value, "googleapis.com") {
		ensureGoogleCredential(ctx, svc.store, cc, &cc.Connection)
		runner = google.NewWrapper(resolved.Service.BaseURL.Value, cc.ID)
	} else {
		runner, err = restDriver.RunnerFromConnection(resolved, cc.ID)
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

	if err != nil && statusCode == 0 {
		return types.ConfiguredConnectionCheckStatus{OK: false, Message: "network error: " + err.Error()}
	}

	if err != nil {
		if statusCode == 404 || statusCode == 405 {
			// If the server returns Not Found or Method Not Allowed for the root path,
			// it means the request successfully passed the authentication layer.
			return types.ConfiguredConnectionCheckStatus{OK: true}
		}

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

// liveConnection replaces cc.Connection with the current live connection definition
// fetched from the store. This is a temporary measure to avoid stale snapshots
// until a proper sync/versioning flow is implemented.
func (svc *configuredConnection) liveConnection(ctx context.Context, cc *types.ConfiguredConnection) error {
	conn, err := loadConnection(ctx, svc.store, cc.ConnectionID)
	if err != nil {
		return err
	}
	cc.Connection = *conn
	return nil
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

	// ConfigParams are declared in the spec under service.params and stored
	// with scope ["service"]. Ensure they are also reachable by unscoped templates.
	svcScope := joinScope([]string{"service"})
	for _, cp := range conn.Service.Params {
		if v, ok := lookup[svcScope+"|"+cp.Name]; ok {
			lookup["|"+cp.Name] = v
		}
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
			"human/connection-id":         strconv.FormatUint(resolved.ID, 10),
			"human/connection-revision":   strconv.Itoa(resolved.Revision),
			"human/configured-connection": strconv.FormatUint(conn.ID, 10),
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

	for connID, ccs := range byConn {
		conn, err := loadConnection(ctx, svc.store, connID)
		if err != nil {
			continue
		}
		svc.registerOperations(conn, ccs)
		for _, cc := range ccs {
			svc.registerWebhookTriggers(cc)
		}
	}
}

// registerOperations converts each ConnectionOperation into a ConstructFunction
// and adds it to the automation construct library.
func (svc *configuredConnection) registerOperations(conn *types.Connection, ccs []types.ConfiguredConnection) {
	if len(ccs) == 0 || len(conn.Operations) == 0 {
		return
	}

	fns := make([]atypes.ConstructFunction, 0, len(conn.Operations))
	for _, op := range conn.Operations {
		fn := operationToFunction(*conn, ccs, op)
		fns = append(fns, fn)
	}

	automationService.ConstructLibrary().AddFunctions(fns...)
}

func operationToFunction(conn types.Connection, ccs []types.ConfiguredConnection, op types.ConnectionOperation) atypes.ConstructFunction {
	ref := fmt.Sprintf("conn_%d_%s", conn.ID, op.Handle)

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
	input[0].Input.Type = "Select"

	segments[0].Sections[0].Elements = append(input, segments[0].Sections[0].Elements...)

	// Inject discovered resource options (spreadsheetId → Select, sheetName → flat tab list)
	if len(ccs) > 0 && len(ccs[0].Config.Discovery) > 0 {
		injectDiscoveredOptions(segments, ccs[0].Config.Discovery)
	}

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
				case int64:
					configID = uint64(id)
				case float64:
					configID = uint64(id)
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

			// Build shared vars map once; mime_build steps may add to it.
			vars := make(map[string]any, len(cc.Config.Params))
			for _, p := range cc.Config.Params {
				if len(p.Scope) > 0 && p.Scope[0] == "service" {
					vars[p.Name] = p.Value
				}
			}
			for k, v := range in.Dict() {
				vars[k] = v
			}

			if dbg, _ := json.Marshal(vars); dbg != nil {
				fmt.Printf("[DEBUG] connector vars before template resolution: %s\n", string(dbg))
			}

			var respBody []byte
			for _, step := range op.Steps {
				switch step.Type {
				case "mime_build":
					if step.MimeBuild == nil {
						continue
					}
					mb := step.MimeBuild
					to, _ := resolveTemplate(types.ConnectionTemplate{Value: mb.To}, vars)
					subject, _ := resolveTemplate(types.ConnectionTemplate{Value: mb.Subject}, vars)
					body, _ := resolveTemplate(types.ConnectionTemplate{Value: mb.Body}, vars)
					from, _ := resolveTemplate(types.ConnectionTemplate{Value: mb.From}, vars)
					outKey := mb.Output
					if outKey == "" {
						outKey = "raw"
					}
					raw, err := buildMIMEEmail(from, to, subject, body)
					if err != nil {
						return nil, fmt.Errorf("mime_build: %w", err)
					}
					vars[outKey] = raw

				case "http":
					if step.HTTP == nil {
						continue
					}
					path, headers, payload, err := buildHTTPRequest(*step.HTTP, in, vars)
					if err != nil {
						return nil, err
					}

					statusCode, outHeaders, body, err := execute(ctx, step.HTTP.Method, path, headers, payload)
					if err != nil {
						return nil, fmt.Errorf("operation execution failed: %w", err)
					}

					if err = checkHTTPResponse(statusCode, outHeaders, body); err != nil {
						return nil, err
					}

					respBody = body
				}
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

	// Google AIP-style custom verbs (e.g. ":batchUpdate") start with ":" and
	// cannot be parsed by url.Parse (it treats ":" as a scheme separator).
	// For these we skip the parse-and-merge step; query params are appended
	// manually and the executor's buildURL handles the base-URL concatenation.
	if strings.HasPrefix(path, ":") {
		if len(http.QueryParams) > 0 {
			q := url.Values{}
			for k, tpl := range http.QueryParams {
				v, e := resolveTemplate(tpl, vars)
				if e != nil {
					return "", nil, nil, fmt.Errorf("queryParam %q: %w", k, e)
				}
				q.Set(k, v)
			}
			path = path + "?" + q.Encode()
		}
	} else {
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
	}

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
			raw := v
			if tv, is := v.(expr.TypedValue); is {
				type dictI interface{ Dict() map[string]any }
				type sliceI interface{ Slice() []any }

				if d, ok := tv.(dictI); ok {
					raw = d.Dict()
				} else if s, ok := tv.(sliceI); ok {
					raw = s.Slice()
				} else {
					raw = tv.Get()
				}
			}

			switch s := raw.(type) {
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

// normalizeParamType maps common type aliases from connection operation specs
// to the canonical type names registered in the automation expr registry.
func normalizeParamType(t string) string {
	switch strings.ToLower(t) {
	case "number":
		return "Integer"
	}
	return t
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

	// Build a set of variables produced by intermediate steps
	// so we don't expose them to the user as derived params.
	for _, step := range op.Steps {
		if step.MimeBuild != nil {
			outKey := step.MimeBuild.Output
			if outKey == "" {
				outKey = "raw"
			}
			inLookup[outKey] = true
		}
	}

	// Explicit arguments — skip any that are already covered by service-level params
	for _, in := range op.Input {
		if serviceParams[in.Name] {
			continue
		}
		inLookup[in.Name] = true
		params = append(params, &atypes.Param{
			ArgumentName: in.Name,
			Types:        []string{normalizeParamType(in.Type)},
			Required:     in.Required,
			Aggregate:    in.Aggregate,
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
						"human.connectionID": {Val: strconv.FormatUint(cc.ID, 10)},
						"human.webhookEvent": {Val: wh.Event},
						"human.resource":     {Val: res.Handle},
						"human.connection":   {Val: resolved.Handle},
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
				Groups:       []string{cc.Connection.Meta.Short},
				Properties:   props,
				Meta: &atypes.ConstructTriggerMeta{
					Short: fmt.Sprintf("%s: %s", labelFromName(res.Handle), labelFromName(wh.Event)),
				},
			})
		}
	}

	if len(tt) > 0 {
		automationService.ConstructLibrary().AddTriggers(tt...)
	}
}

// buildMIMEEmail constructs a minimal RFC 2822 email message and returns it
// base64url-encoded, suitable for the Gmail API "raw" field.
func buildMIMEEmail(from, to, subject, body string) (string, error) {
	var buf bytes.Buffer
	if from != "" {
		buf.WriteString("From: " + from + "\r\n")
	}
	buf.WriteString("To: " + to + "\r\n")
	buf.WriteString("Subject: " + subject + "\r\n")
	buf.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	buf.WriteString("\r\n")
	buf.WriteString(body)
	return base64.URLEncoding.EncodeToString(buf.Bytes()), nil
}

// syncGoogleDiscovery fetches accessible Google resources for the given
// configured connection and stores them in cc.Config.Discovery.
// Only runs for Google connectors (baseURL contains "googleapis.com").
// On any error, logs a warning and returns nil (non-fatal).
func (svc *configuredConnection) syncGoogleDiscovery(ctx context.Context, cc *types.ConfiguredConnection) error {
	resolved := svc.resolveTemplates(&cc.Connection, cc.Config.Params)
	baseURL := resolved.Service.BaseURL.Value

	if !strings.Contains(baseURL, "googleapis.com") {
		return nil // not a Google connector
	}

	ensureGoogleCredential(ctx, svc.store, cc, &cc.Connection)

	if cc.Config.Discovery == nil {
		cc.Config.Discovery = make(map[string]json.RawMessage)
	}

	// Google Calendar: fetch the impersonated user's calendar list.
	if strings.Contains(baseURL, "googleapis.com/calendar") {
		calWrapper := google.NewWrapper("https://www.googleapis.com/calendar/v3", cc.ID)
		_, _, calBody, err := calWrapper.Run(ctx, "GET",
			"/users/me/calendarList?fields=items(id,summary)&maxResults=250",
			nil, nil)
		if err != nil {
			svc.logger.Warn("google discovery: calendarList unavailable",
				zap.Uint64("ccID", cc.ID), zap.Error(err))
			return nil
		}

		var calResp struct {
			Items []struct {
				ID      string `json:"id"`
				Summary string `json:"summary"`
			} `json:"items"`
		}
		if err := json.Unmarshal(calBody, &calResp); err != nil {
			return nil
		}

		calendars := make([]atypes.SelectItem, 0, len(calResp.Items))
		for _, c := range calResp.Items {
			calendars = append(calendars, atypes.SelectItem{Value: c.ID, Label: c.Summary})
		}
		calsJSON, _ := json.Marshal(calendars)
		cc.Config.Discovery["calendars"] = calsJSON
		return store.UpdateConfiguredConnection(ctx, svc.store, cc)
	}

	// Google Sheets: list spreadsheets via Drive, then fetch tabs per spreadsheet.
	driveWrapper := google.NewWrapper("https://www.googleapis.com/drive/v3", cc.ID)

	q := url.QueryEscape("mimeType='application/vnd.google-apps.spreadsheet'")
	// 1. List all spreadsheets the service account can access
	_, _, body, err := driveWrapper.Run(ctx, "GET",
		"/files?q="+q+"&fields=files(id,name)&pageSize=200",
		nil, nil)
	if err != nil {
		svc.logger.Warn("google discovery: drive API unavailable",
			zap.Uint64("ccID", cc.ID), zap.Error(err))
		return nil
	}

	var driveResp struct {
		Files []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"files"`
	}
	if err := json.Unmarshal(body, &driveResp); err != nil {
		return nil
	}

	spreadsheets := make([]atypes.SelectItem, 0, len(driveResp.Files))
	for _, f := range driveResp.Files {
		spreadsheets = append(spreadsheets, atypes.SelectItem{Value: f.ID, Label: f.Name})
	}

	// 2. Fetch tabs for each spreadsheet
	sheetsWrapper := google.NewWrapper("https://sheets.googleapis.com/v4/spreadsheets", cc.ID)
	tabs := make(map[string][]atypes.SelectItem, len(driveResp.Files))

	for _, f := range driveResp.Files {
		_, _, shBody, err := sheetsWrapper.Run(ctx, "GET",
			"/"+f.ID+"?fields=sheets.properties%28title%29",
			nil, nil)
		if err != nil {
			continue // partial data is fine
		}

		var shResp struct {
			Sheets []struct {
				Properties struct {
					Title string `json:"title"`
				} `json:"properties"`
			} `json:"sheets"`
		}
		if err := json.Unmarshal(shBody, &shResp); err != nil {
			continue
		}

		items := make([]atypes.SelectItem, 0, len(shResp.Sheets))
		for _, sh := range shResp.Sheets {
			items = append(items, atypes.SelectItem{Value: sh.Properties.Title, Label: sh.Properties.Title})
		}
		tabs[f.ID] = items
	}

	// 3. Persist into cc.Config.Discovery
	ssJSON, _ := json.Marshal(spreadsheets)
	tabsJSON, _ := json.Marshal(tabs)
	cc.Config.Discovery["spreadsheets"] = ssJSON
	cc.Config.Discovery["tabs"] = tabsJSON

	return store.UpdateConfiguredConnection(ctx, svc.store, cc)
}

// injectDiscoveredOptions walks segment elements and sets Options + Type="Select"
// for known discoverable param names.
// - spreadsheetId → list of spreadsheets
// - sheetName     → deduplicated flat list of all tab names across all spreadsheets
func injectDiscoveredOptions(segments []atypes.ConstructSegment, discovery map[string]json.RawMessage) {
	var spreadsheets []atypes.SelectItem
	if raw, ok := discovery["spreadsheets"]; ok {
		_ = json.Unmarshal(raw, &spreadsheets)
	}

	// Build flat, deduplicated tab list
	var tabs []atypes.SelectItem
	if raw, ok := discovery["tabs"]; ok {
		var tabMap map[string][]atypes.SelectItem
		if err := json.Unmarshal(raw, &tabMap); err == nil {
			seen := make(map[string]bool)
			for _, items := range tabMap {
				for _, item := range items {
					if !seen[item.Value] {
						seen[item.Value] = true
						tabs = append(tabs, item)
					}
				}
			}
		}
	}

	for si := range segments {
		for seci := range segments[si].Sections {
			for ei := range segments[si].Sections[seci].Elements {
				el := &segments[si].Sections[seci].Elements[ei]
				switch el.Input.Argument {
				case "spreadsheetId":
					if len(spreadsheets) > 0 {
						el.Input.Type = "Select"
						el.Input.Options = spreadsheets
					}
				case "sheetName":
					if len(tabs) > 0 {
						el.Input.Type = "Select"
						el.Input.Options = tabs
					}
				case "calendarId":
					var calendars []atypes.SelectItem
					if raw, ok := discovery["calendars"]; ok {
						_ = json.Unmarshal(raw, &calendars)
					}
					if len(calendars) > 0 {
						el.Input.Type = "Select"
						el.Input.Options = calendars
					}
				}
			}
		}
	}
}

// RefreshDiscovery re-fetches Google resource discovery for the given CC,
// re-registers operations so the construct library reflects the new options,
// and returns a summary.
func (svc *configuredConnection) RefreshDiscovery(ctx context.Context, ID uint64) (map[string]any, error) {
	cc, err := loadConfiguredConnection(ctx, svc.store, ID)
	if err != nil {
		return nil, err
	}

	// Evict the cached credential so ensureGoogleCredential re-reads all params
	// (including dwdSubject) from the stored configuration rather than reusing
	// a stale token that may have been minted without impersonation.
	_ = cred_registry.Default().Delete(cc.ID)

	if err = svc.syncGoogleDiscovery(ctx, cc); err != nil {
		return nil, err
	}

	// Re-register operations so the construct library picks up new options
	siblings, _, _ := store.SearchConfiguredConnections(ctx, svc.store, types.ConfiguredConnectionFilter{
		ConnectionID: cc.ConnectionID,
		Status:       []string{"active"},
	})
	allCCs := make([]types.ConfiguredConnection, 0, len(siblings))
	for _, s := range siblings {
		allCCs = append(allCCs, *s)
	}
	if conn, connErr := loadConnection(ctx, svc.store, cc.ConnectionID); connErr == nil {
		svc.registerOperations(conn, allCCs)
	}

	var spreadsheetCount int
	if raw, ok := cc.Config.Discovery["spreadsheets"]; ok {
		var ss []atypes.SelectItem
		if json.Unmarshal(raw, &ss) == nil {
			spreadsheetCount = len(ss)
		}
	}

	var calendarCount int
	if raw, ok := cc.Config.Discovery["calendars"]; ok {
		var cals []atypes.SelectItem
		if json.Unmarshal(raw, &cals) == nil {
			calendarCount = len(cals)
		}
	}

	return map[string]any{
		"ok":               true,
		"spreadsheetCount": spreadsheetCount,
		"calendarCount":    calendarCount,
	}, nil
}

// StartDiscoveryRefreshLoop starts a background goroutine that periodically
// re-fetches Google resource discovery for all active Google configured
// connections and re-registers their operations.
func (svc *configuredConnection) StartDiscoveryRefreshLoop(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				svc.refreshAllGoogleConnections(ctx)
			}
		}
	}()
}

func (svc *configuredConnection) refreshAllGoogleConnections(ctx context.Context) {
	set, _, err := store.SearchConfiguredConnections(ctx, svc.store, types.ConfiguredConnectionFilter{
		Status: []string{"active"},
	})
	if err != nil {
		svc.logger.Warn("google discovery refresh: failed to load connections", zap.Error(err))
		return
	}

	// Group by connection so registerOperations sees all siblings
	byConn := make(map[uint64][]types.ConfiguredConnection)
	for _, cc := range set {
		resolved := svc.resolveTemplates(&cc.Connection, cc.Config.Params)
		if !strings.Contains(resolved.Service.BaseURL.Value, "googleapis.com") {
			continue
		}
		if err := svc.syncGoogleDiscovery(ctx, cc); err != nil {
			svc.logger.Warn("google discovery refresh: sync failed",
				zap.Uint64("ccID", cc.ID), zap.Error(err))
		}
		byConn[cc.ConnectionID] = append(byConn[cc.ConnectionID], *cc)
	}
	for connID, ccs := range byConn {
		conn, err := loadConnection(ctx, svc.store, connID)
		if err != nil {
			continue
		}
		svc.registerOperations(conn, ccs)
	}
}
