package dal

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/ql"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

type (
	PKValues map[string]any

	ConnectionParams struct {
		Type   string         `json:"type"`
		Params map[string]any `json:"params"`
	}

	Connection interface {
		// Meta

		// Models returns all the models the underlying connection already supports
		//
		// This is useful when adding support for new models since we can find out what
		// can work out of the box.
		Models(context.Context) (ModelSet, error)

		// Operations returns all of the operations the given store supports
		Operations() OperationSet

		// Can returns true if this store can handle the given operations
		Can(operations ...Operation) bool

		// DML stuff

		// Create stores the given data into the underlying database
		Create(ctx context.Context, m *Model, rr ...ValueGetter) error

		// Update updates the given value in the underlying connection
		Update(ctx context.Context, m *Model, r ValueGetter) error

		// Lookup returns one bit of data
		Lookup(context.Context, *Model, ValueGetter, ValueSetter) error

		// Search returns an iterator which can be used to access all if the bits
		Search(context.Context, *Model, filter.Filter) (Iterator, error)

		// Count returns the total number of rows matching the filter
		Count(context.Context, *Model, filter.Filter) (uint, error)

		// Analyze returns the operation analysis the connection can perform for the model
		Analyze(ctx context.Context, m *Model) (map[string]OpAnalysis, error)

		// Aggregate returns the iterator with aggregated data from the base model
		Aggregate(ctx context.Context, m *Model, f filter.Filter, groupBy []AggregateAttr, aggrExpr []AggregateAttr, having *ql.ASTNode) (i Iterator, _ error)

		// Delete deletes the given value
		Delete(ctx context.Context, m *Model, pkv ValueGetter) error

		// Truncate deletes all the data for the given model
		Truncate(ctx context.Context, m *Model) error

		// DDL stuff

		// // returns all attribute types that driver supports
		// AttributeTypes() []data.AttributeType

		// CreateModel adds support for the given models to the underlying database
		//
		// The operation returns an error if any of the models already exists.
		CreateModel(context.Context, ...*Model) error

		// DeleteModel removes support for the given model from the underlying database
		DeleteModel(context.Context, ...*Model) error

		// UpdateModel requests for metadata changes to the existing model
		//
		// Only metadata (such as idents) are affected; attributes can not be changed here
		UpdateModel(ctx context.Context, old *Model, new *Model) error

		// AssertSchemaAlterations returns a new set of Alterations based on what the underlying
		// schema already provides -- it discards alterations for column additions that already exist, etc.
		AssertSchemaAlterations(ctx context.Context, sch *Model, aa ...*Alteration) ([]*Alteration, error)

		// ApplyAlteration applies the given alterations to the underlying schema
		ApplyAlteration(ctx context.Context, sch *Model, aa ...*Alteration) []error
	}

	ConnectionCloser interface {
		// Close closes the store connection allowing the driver to perform potential
		// cleanup operations
		Close(ctx context.Context) error
	}

	// Store provides an interface which CRS uses to interact with the underlying database

	ValueGetter interface {
		CountValues() map[string]uint
		GetValue(string, uint) (any, error)
	}

	ValueSetter interface {
		SetValue(string, uint, any) error
	}

	ConnectorFn func(ctx context.Context, dsn string) (Connection, error)

	DriverConnectionParam struct {
		Key        string `json:"key"`
		ValueType  string `json:"valueType"`
		MultiValue bool   `json:"multiValue"`
	}

	DriverConnectionConfig struct {
		Type   string                  `json:"type"`
		Params []DriverConnectionParam `json:"params"`
	}

	Driver struct {
		Type       string                 `json:"type"`
		Connection DriverConnectionConfig `json:"connection"`
		Operations OperationSet           `json:"operations"`
	}
)

var (
	registeredConnectors = make(map[string]ConnectorFn)
	registeredDrivers    = make(map[string]Driver)
)

func (pkv PKValues) CountValues() map[string]uint {
	c := make(map[string]uint)
	for k := range pkv {
		c[k] = 1
	}

	return c
}

func (pkv PKValues) GetValue(key string, _ uint) (any, error) {
	if val, has := pkv[key]; has {
		return val, nil
	} else {
		return nil, fmt.Errorf("missing")
	}
}

// RegisterConnector registers a new connector for the given DSN schema
//
// In case of a duplicate schema the latter overwrites the prior
func RegisterConnector(fn ConnectorFn, tt ...string) {
	for _, t := range tt {
		registeredConnectors[t] = fn
	}
}

func RegisterDriver(d Driver) {
	registeredDrivers[d.Type] = d
}

// connect opens a new StoreConnection for the given CRS
func connect(ctx context.Context, log *zap.Logger, isDevelopment bool, cp ConnectionParams) (Connection, error) {
	if cp.Params == nil {
		return nil, fmt.Errorf("cannot open connection: connection parameters not defined")
	}

	switch cp.Type {
	case "corteza::dal:connection:dsn":
		return connectRDBMS(ctx, log, isDevelopment, cp)

	case "corteza::dal:connection:rest":
		return connectREST(ctx, log, isDevelopment, cp)

	default:
		return nil, fmt.Errorf("cannot open connection: unsupported connection (got: %q)", cp.Type)
	}
}

func connectRDBMS(ctx context.Context, log *zap.Logger, isDevelopment bool, cp ConnectionParams) (Connection, error) {
	if _, ok := cp.Params["dsn"]; !ok {
		return nil, fmt.Errorf("cannot open connection: DSN not provided")
	}

	dsn := cp.Params["dsn"].(string)

	if isDevelopment {
		if strings.Contains(dsn, "{version}") {
			log.Warn("You're using DB_DSN with {version}, It is still in EXPERIMENTAL phase")
			log.Warn("Should be used only for development mode")
			log.Warn("You may experience instability")
		}
		expr := regexp.MustCompile(`[.\-]+`)
		version := expr.ReplaceAllString(os.Getenv("BUILD_VERSION"), "_")
		dsn = strings.Replace(dsn, "{version}", version, 1)
	}

	var storeType = strings.SplitN(dsn, "://", 2)[0]
	if storeType == "" {
		// Backward compatibility
		storeType = "mysql"
	}

	if conn, ok := registeredConnectors[storeType]; ok {
		return conn(ctx, dsn)
	} else {
		return nil, fmt.Errorf("unknown store type used: %q (check your database configuration)", storeType)
	}
}

func connectREST(ctx context.Context, log *zap.Logger, isDevelopment bool, cp ConnectionParams) (c Connection, err error) {
	if _, ok := cp.Params["url"]; !ok {
		return nil, fmt.Errorf("cannot connect to the REST API: missing parameter: url")
	}

	// For now, the only store type is REST API so this is fine
	storeType := "restapi"
	url := cast.ToString(cp.Params["url"])
	if len(url) == 0 {
		return nil, fmt.Errorf("cannot connect to the REST API: invalid parameter: url")
	}

	d, err := ParseDSN(url)
	if err != nil {
		return
	}

	d, err = expandDSN(d, cp)
	if err != nil {
		return
	}

	if conn, ok := registeredConnectors[storeType]; ok {
		return conn(ctx, d.ToDSN())
	} else {
		return nil, fmt.Errorf("unknown store type used: %q (check your database configuration)", storeType)
	}
}

func expandDSN(base DSN, cp ConnectionParams) (out DSN, err error) {
	out = base

	if auth, ok := cp.Params["auth"]; ok {

		// @todo improve this
		bb, err := json.Marshal(auth)
		if err != nil {
			return out, err
		}

		aux := struct {
			Method string
			Params map[string]any
		}{}

		err = json.Unmarshal(bb, &aux)
		if err != nil {
			return out, err
		}

		// @todo validation and all that :)
		out.AuthType = aux.Method
		out.Token = aux.Params["token"].(string)
		out.APIKey = aux.Params["APIKey"].(string)
		out.APIKeyHeader = aux.Params["APIKeyHeader"].(string)
		out.ClientID = aux.Params["clientID"].(string)
		out.ClientSecret = aux.Params["clientSecret"].(string)
		out.TokenURL = aux.Params["tokenURL"].(string)
	}

	if dops, ok := cp.Params["defaultOps"]; ok {
		aux := map[string]any{}

		bb, err := json.Marshal(dops)
		if err != nil {
			return out, err
		}

		err = json.Unmarshal(bb, &aux)
		if err != nil {
			return out, err
		}

		out.Arbitrary = aux
	}

	return
}

func NewDSNDriverConnectionConfig() DriverConnectionConfig {
	return DriverConnectionConfig{
		Type: "corteza::dal:connection:dsn",
		Params: []DriverConnectionParam{{
			Key:       "dsn",
			ValueType: "string",
		}},
	}
}

func NewRESTriverConnectionConfig() DriverConnectionConfig {
	return DriverConnectionConfig{
		Type: "corteza::dal:connection:rest",
		Params: []DriverConnectionParam{{
			Key:       "dsn",
			ValueType: "string",
		}},
	}
}

//func NewHTTPDriverConnectionConfig() DriverConnectionConfig {
//	panic("not implemented NewHTTPDriverConnectionConfig")
//	return DriverConnectionConfig{
//		Type:   "corteza::dal:connection:http",
//		Params: []DriverConnectionParam{{}},
//	}
//}
//func NewFederatedNodeDriverConnectionConfig() DriverConnectionConfig {
//	panic("not implemented NewFederatedNodeDriverConnectionConfig")
//	return DriverConnectionConfig{
//		Type:   "corteza::dal:connection:federated-node",
//		Params: []DriverConnectionParam{{}},
//	}
//}

//func NewDSNConnection(dsn string) ConnectionParams {
//	return ConnectionParams{
//		Type: "corteza::dal:connection:dsn",
//		Params: map[string]any{
//			"dsn": dsn,
//		},
//	}
//}

//func NewHTTPConnection(url string, headers, query map[string][]string) ConnectionParams {
//	return ConnectionParams{
//		Type: "corteza::dal:connection:http",
//		Params: map[string]any{
//			"url":     url,
//			"headers": headers,
//			"query":   query,
//		},
//	}
//}
//
//func NewFederatedNodeConnection(url string, pairToken, authToken string) ConnectionParams {
//	return ConnectionParams{
//		Type: "corteza::dal:connection:federation-node",
//		Params: map[string]any{
//			"baseURL":   url,
//			"pairToken": pairToken,
//			"authToken": authToken,
//		},
//	}
//}
