package factory

import (
	"fmt"
	"reflect"
	"strings"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// DatabaseCredential is a configuration struct for a database connection
type DatabaseCredential struct {
	DSN string
}

// DatabaseFactory contains all database credentials and instances
type DatabaseFactory struct {
	credentials map[string]*DatabaseCredential
	instances   map[string]*DB

	ProfilerStdout DatabaseProfilerStdout
	ProfilerMemory DatabaseProfilerMemory
}

// The default database factory
var Database *DatabaseFactory

func init() {
	Database = &DatabaseFactory{}
	Database.credentials = make(map[string]*DatabaseCredential)
	Database.instances = make(map[string]*DB)
}

// Add a new named database credential to the database factory
//
// The function will store a named database credential for use with
// the Get function. Generally you'd need to call Add at least once.
// In order to set the default database credential, use "default" as
// the first parameter.
//
// Example:
//
// ```
// factory.Database.Add("default", "sqlapi:sqlapi@tcp(db1:3306)/sqlapi?collation=utf8mb4_general_ci")
// ```
//
// By default, additional options will be added to the credentials:
//
// - collation will be set to `utf8_general_ci`
// - parseTime will be set to `true`
// - loc will be set to `Local`
//
// If your passed DSN will include any of these options, the default values will not
// be applied, and your custom settings will be honored.
func (r *DatabaseFactory) Add(name string, config interface{}) {
	switch val := config.(type) {
	case string:
		r.credentials[name] = &DatabaseCredential{DSN: val}
	case DatabaseCredential:
		r.credentials[name] = &val
	default:
		panic("factory.Database.Add can take config as string|factory.DatabaseCredential")
	}
}

// GetDSN returns the augmented DSN from the named database
func (r *DatabaseFactory) GetDSN(name string) (string, error) {
	addOption := func(s, match, option string) string {
		if !strings.Contains(s, match) {
			s += option
		}
		return s
	}

	if value, ok := r.credentials[name]; ok {
		value.DSN = addOption(value.DSN, "?", "?")
		value.DSN = addOption(value.DSN, "collation=", "&collation=utf8_general_ci")
		value.DSN = addOption(value.DSN, "parseTime=", "&parseTime=true")
		value.DSN = addOption(value.DSN, "loc=", "&loc=Local")
		value.DSN = strings.Replace(value.DSN, "?&", "?", 1)
		return value.DSN, nil
	}
	return "", errors.New("No configuration found for database: " + name)
}

// Get returns a database connection
//
// If you don't request a database connection by name, `factory.Database.Get()` will
// return the connection with the name "default". If you supply one or more names as
// the parameter, the first successful connection will be returned.
//
// It is fine to call this function per request, as a singleton instance is returned
// for each call with the same parameters.
//
// This behavior enables sharding workloads between hosts. One could randomize the
// parameters using [math/rand#Shuffle](https://tip.golang.org/pkg/math/rand/#Shuffle), or
// provide a consistent hashing method based on server hostname/IP, or even go so far
// to retrieve the hosts to connect to from some sort of inventory like etcd.
//
// The most general use case is that you will only call `factory.Database.Get()` once,
// and then pass the resulting `*DB` forward in your application. Requesting custom named
// connections also provides a way to access different parts of the database, or different
// database altogether, depending on your microservice data distribution.
func (r *DatabaseFactory) Get(dbName ...string) (*DB, error) {
	names := dbName
	if len(names) == 0 {
		names = []string{"default"}
	}
	for _, name := range names {
		if value, ok := r.instances[name]; ok {
			return value, nil
		}
		dsn, _ := r.GetDSN(name)
		if dsn != "" {
			handle, err := sqlx.Open("mysql", dsn)
			if err != nil {
				return nil, err
			}
			r.instances[name] = &DB{handle, nil}
			return r.instances[name], nil
		}
	}
	return nil, fmt.Errorf("No configuration found for database: %v", names)
}

// DB struct encapsulates sqlx.DB to add new functions
type DB struct {
	*sqlx.DB

	Profiler DatabaseProfiler
}

// Quiet will return a DB handle without a profiler (throw-away)
func (r *DB) Quiet() *DB {
	return &DB{
		DB: r.DB,
	}
}

// SetFields will provide a string with SQL named bindings from a string slice
func (r *DB) SetFields(fields []string) string {
	idx := 0
	sql := ""
	for _, field := range fields {
		if idx > 0 {
			sql = sql + ", "
		}
		idx++
		sql = sql + field + "=:" + field
	}
	return sql
}

// Select is a helper function that will ignore sql.ErrNoRows
func (r *DB) Select(dest interface{}, query string, args ...interface{}) error {
	var err error
	if r.Profiler != nil {
		ctx := DatabaseProfilerContext{}.new(query, args...)
		err = r.DB.Select(dest, query, args...)
		r.Profiler.Post(ctx)
	} else {
		err = r.DB.Select(dest, query, args...)
	}
	// clear no rows returned error
	if err == sql.ErrNoRows {
		return nil
	}
	return errors.Wrap(err, "select query failed")
}

// Get is a helper function that will ignore sql.ErrNoRows
func (r *DB) Get(dest interface{}, query string, args ...interface{}) error {
	var err error
	if r.Profiler != nil {
		ctx := DatabaseProfilerContext{}.new(query, args...)
		err = r.DB.Get(dest, query, args...)
		r.Profiler.Post(ctx)
	} else {
		err = r.DB.Get(dest, query, args...)
	}
	// clear no rows returned error
	if err == sql.ErrNoRows {
		return nil
	}
	return errors.Wrap(err, "get query failed")
}

// set uses reflection to iterate over struct fields tags, producing bindings for struct values
func (r *DB) set(data interface{}) string {
	message_value := reflect.ValueOf(data)
	if message_value.Kind() == reflect.Ptr {
		message_value = message_value.Elem()
	}

	message_fields := make([]string, message_value.NumField())

	for i := 0; i < len(message_fields); i++ {
		fieldType := message_value.Type().Field(i)
		message_fields[i] = fieldType.Tag.Get("db")
	}

	sql := ""
	for _, tagFull := range message_fields {
		if tagFull != "" && tagFull != "-" {
			tag := strings.Split(tagFull, ",")
			sql = sql + " " + tag[0] + "=:" + tag[0] + ","
		}
	}
	return sql[1 : len(sql)-1]
}

// Replace is a helper function which will issue an `replace` statement to the database
func (r *DB) Replace(table string, data interface{}) error {
	sql := "replace into " + table + " set " + r.set(data)
	_, err := r.NamedExec(sql, data)
	return err
}

// Insert is a helper function which will issue an `insert` statement to the database
func (r *DB) Insert(table string, data interface{}) error {
	sql := "insert into " + table + " set " + r.set(data)
	_, err := r.NamedExec(sql, data)
	return err
}

// InsertIgnore is a helper function which will issue an `insert ignore` statement to the database
func (r *DB) InsertIgnore(table string, data interface{}) error {
	sql := "insert ignore into " + table + " set " + r.set(data)
	_, err := r.NamedExec(sql, data)
	return err
}
