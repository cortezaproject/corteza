package factory

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"database/sql"

	"github.com/go-sql-driver/mysql"
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
	case *string:
		r.credentials[name] = &DatabaseCredential{DSN: *val}
	case string:
		r.credentials[name] = &DatabaseCredential{DSN: val}
	case DatabaseCredential:
		r.credentials[name] = &val
	default:
		panic("factory.Database.Add can take config as string|*string|factory.DatabaseCredential")
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
			handle, err := sqlx.Connect("mysql", dsn)
			if err != nil {
				return nil, err
			}
			r.instances[name] = &DB{
				handle,
				context.Background(),
				0,
				nil,
				&sql.TxOptions{
					ReadOnly: false,
				},
				nil,
			}
			return r.instances[name], nil
		}
	}
	return nil, fmt.Errorf("No configuration found for database: %v", names)
}

// MustGet is the same as Get, except it will panic on connection error
func (r *DatabaseFactory) MustGet(dbName ...string) *DB {
	db, err := r.Get(dbName...)
	if err != nil {
		panic(err)
	}
	return db
}

// DB struct encapsulates sqlx.DB to add new functions
type DB struct {
	*sqlx.DB

	ctx context.Context

	inTx   int32
	Tx     *sqlx.Tx
	TxOpts *sql.TxOptions

	Profiler DatabaseProfiler
}

// Quiet will return a DB handle without a profiler (throw-away)
func (r *DB) Quiet() *DB {
	return &DB{
		r.DB,
		r.ctx,
		r.inTx,
		r.Tx,
		r.TxOpts,
		nil,
	}
}

// With will return a DB handle with a bound context (throw-away)
func (r *DB) With(ctx context.Context) *DB {
	return &DB{
		r.DB,
		ctx,
		r.inTx,
		r.Tx,
		r.TxOpts,
		r.Profiler,
	}
}

// Begin will create a transaction in the DB with a context
func (r *DB) Begin() (err error) {
	if r.inTx > 0 {
		_, err = r.Exec(fmt.Sprintf("SAVEPOINT sp_%d", r.inTx))
	}
	if r.inTx == 0 {
		if r.ctx == nil {
			r.Log(func() {
				r.Tx, err = r.DB.Beginx()
			}, "BEGIN;")
		} else {
			r.Log(func() {
				r.Tx, err = r.DB.BeginTxx(r.ctx, r.TxOpts)
			}, "BEGIN; -- with context")
		}
	}
	if err != nil {
		return errors.WithStack(err)
	}

	r.inTx++
	return nil
}

// Transaction will create a transaction and invoke a callback
func (r *DB) Transaction(callback func() error) (err error) {
	var tries int

	// Perform transaction statements
	tries = 0
	for {
		// Start transaction
		if err = r.Begin(); err != nil {
			return
		}

		if err = callback(); err == nil {
			break
		}

		tries++
		if tries > 3 {
			log.Printf("Retried transaction %d times, aborting", tries-1)
			break
		}

		// Break out if the causer is not a MySQL error
		cause, ok := (errors.Cause(err)).(*mysql.MySQLError)
		if !ok {
			log.Printf("Returned error cause is not a MySQLError, %#v", err)
			break
		}

		// restart transaction:
		//   - 1205: lock within transaction (unit tested),
		//   - 1213: deadlock found
		if cause.Number == 1205 || cause.Number == 1213 {
			log.Printf("Restarting transaction (try=%d): %s", tries, cause)
			r.Rollback()
			time.Sleep(50 * time.Millisecond)
			continue
		}
		log.Println("Can't handle transaction error")
		break
	}

	if err != nil {
		r.Rollback()
		return err
	}
	return r.Commit()
}

func (r *DB) Commit() (err error) {
	if r.Tx != nil {
		if r.inTx <= 1 {
			r.Log(func() {
				err = r.Tx.Commit()
			}, "COMMIT;")
			if err != nil {
				return errors.WithStack(err)
			}
			r.Tx = nil
			r.inTx = 0
			return nil
		}
		if _, err = r.Exec(fmt.Sprintf("RELEASE SAVEPOINT sp_%d", r.inTx-1)); err != nil {
			return errors.WithStack(err)
		}
		r.inTx--
		return nil
	}
	return errors.WithStack(sql.ErrTxDone)
}

func (r *DB) Rollback() (err error) {
	if r.Tx != nil {
		if r.inTx <= 1 {
			r.Log(func() {
				err = r.Tx.Rollback()
			}, "ROLLBACK;")
			if err != nil {
				return errors.WithStack(err)
			}
			r.Tx = nil
			r.inTx = 0
			return nil
		}
		if _, err = r.Exec(fmt.Sprintf("ROLLBACK SAVEPOINT sp_%d", r.inTx-1)); err != nil {
			return errors.WithStack(err)
		}
		r.inTx--
		return nil
	}
	return errors.WithStack(sql.ErrTxDone)
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

// NamedExec adds profiling on top of the parent DB.NameExec
func (r *DB) NamedExec(query string, arg interface{}) (res sql.Result, err error) {
	exec := func() (sql.Result, error) {
		if r.Tx != nil {
			return r.Tx.NamedExecContext(r.ctx, query, arg)
		}
		return r.DB.NamedExecContext(r.ctx, query, arg)
	}

	r.Log(func() {
		res, err = exec()
	}, query, arg)

	return res, errors.Wrap(err, "exec query failed")
}

// Exec adds profiling on top of the parent DB.Exec
func (r *DB) Exec(query string, args ...interface{}) (res sql.Result, err error) {
	exec := func() (sql.Result, error) {
		if r.Tx != nil {
			return r.Tx.ExecContext(r.ctx, query, args...)
		}
		return r.DB.ExecContext(r.ctx, query, args...)
	}

	r.Log(func() {
		res, err = exec()
	}, query, args...)

	return res, errors.Wrap(err, "exec query failed")
}

// Select is a helper function that will ignore sql.ErrNoRows
func (r *DB) Select(dest interface{}, query string, args ...interface{}) error {
	var err error
	exec := func() error {
		if r.Tx != nil {
			return r.Tx.SelectContext(r.ctx, dest, query, args...)
		}
		return r.DB.SelectContext(r.ctx, dest, query, args...)
	}

	r.Log(func() {
		err = exec()
	}, query, args...)

	// clear no rows returned error
	if err == sql.ErrNoRows {
		return nil
	}

	return errors.Wrap(err, "select query failed")
}

// Get is a helper function that will ignore sql.ErrNoRows
func (r *DB) Get(dest interface{}, query string, args ...interface{}) error {
	var err error
	exec := func() error {
		if r.Tx != nil {
			return r.Tx.GetContext(r.ctx, dest, query, args...)
		}
		return r.DB.GetContext(r.ctx, dest, query, args...)
	}

	r.Log(func() {
		err = exec()
	}, query, args...)

	// clear no rows returned error
	if err == sql.ErrNoRows {
		return nil
	}

	return errors.Wrap(err, "get query failed")
}

// set uses reflection to iterate over struct fields tags, producing bindings for struct values
func (r *DB) set(data interface{}, allowed ...string) string {
	set := r.setMap(data, allowed...)
	return r.setImplode(", ", set)
}

func (r *DB) tag(tag string) string {
	if tag != "" && tag != "-" {
		return strings.Split(tag, ",")[0]
	}
	return ""
}

func (r *DB) setMap(data interface{}, allowed ...string) map[string]string {
	message_value := reflect.ValueOf(data)
	if message_value.Kind() == reflect.Ptr {
		message_value = message_value.Elem()
	}

	set := make(map[string]string)
	length := message_value.NumField()
	for i := 0; i < length; i++ {
		fieldType := message_value.Type().Field(i)
		if tag := r.tag(fieldType.Tag.Get("db")); tag != "" {
			set[tag] = ":" + tag
		}
	}

	// limit only to allowed fields
	if len(allowed) > 0 {
		for tag, _ := range set {
			canDelete := true
			for _, key := range allowed {
				if tag == key {
					canDelete = false
					break
				}
			}
			if canDelete {
				delete(set, tag)
			}
		}
	}

	return set
}

func (r *DB) setImplode(delimiter string, set map[string]string) string {
	result := ""
	count := 0
	for key, value := range set {
		if count > 0 {
			result = result + delimiter
		}
		result = result + key + "=" + value
		count++
	}
	return result
}

// Update is a helper function which will issue an `update` statement to the db
func (r *DB) Update(table string, args interface{}, keys ...string) error {
	var err error
	if len(keys) == 0 {
		return errors.New("Full-table update not supported")
	}
	set := r.setMap(args)
	setWhere := make(map[string]string)
	for _, key := range keys {
		value, ok := set[key]
		if !ok {
			return errors.New("Can't update table " + table + " by key " + key + " (no such field in struct)")
		}
		delete(set, key)
		setWhere[key] = value
	}
	if len(set) == 0 {
		return errors.New("Encountered update struct with no fields")
	}
	query := "update " + table + " set " + r.setImplode(", ", set) + " where " + r.setImplode(" AND ", setWhere)
	_, err = r.NamedExec(query, args)
	return err
}

// UpdatePartial is a helper function which will issue an `update` statement to the db
func (r *DB) UpdatePartial(table string, args interface{}, allowed []string, keys ...string) error {
	var err error
	if len(keys) == 0 {
		return errors.New("Full-table update not supported")
	}
	set := r.setMap(args, allowed...)
	setWhere := make(map[string]string)
	for _, key := range keys {
		value, ok := set[key]
		if !ok {
			return errors.New("Can't update table " + table + " by key " + key + " (no such field in struct)")
		}
		delete(set, key)
		setWhere[key] = value
	}
	if len(set) == 0 {
		return errors.New("Encountered update struct with no fields")
	}
	query := "update " + table + " set " + r.setImplode(", ", set) + " where " + r.setImplode(" AND ", setWhere)
	_, err = r.NamedExec(query, args)
	return err
}

// Delete is a helper function which will issue an `delete` statement to the db
func (r *DB) Delete(table string, args interface{}, keys ...string) error {
	var err error
	if len(keys) == 0 {
		return errors.New("Full-table delete not supported")
	}
	set := r.setMap(args)
	setWhere := make(map[string]string)
	for _, key := range keys {
		value, ok := set[key]
		if !ok {
			return errors.New("Can't update table " + table + " by key " + key + " (no such field in struct)")
		}
		delete(set, key)
		setWhere[key] = value
	}
	query := "delete from " + table + " where " + r.setImplode(" AND ", setWhere)
	_, err = r.NamedExec(query, args)
	return err
}

// Replace is a helper function which will issue an `replace` statement to the database
func (r *DB) Replace(table string, args interface{}) error {
	var err error
	query := "replace into " + table + " set " + r.set(args)
	_, err = r.NamedExec(query, args)
	return err
}

// Insert is a helper function which will issue an `insert` statement to the database
func (r *DB) Insert(table string, args interface{}) error {
	var err error
	query := "insert into " + table + " set " + r.set(args)
	_, err = r.NamedExec(query, args)
	return err
}

// InsertIgnore is a helper function which will issue an `insert ignore` statement to the database
func (r *DB) InsertIgnore(table string, args interface{}) error {
	var err error
	query := "insert ignore into " + table + " set " + r.set(args)
	_, err = r.NamedExec(query, args)
	return err
}

func (r *DB) Log(callback func(), query string, args ...interface{}) {
	if r.Profiler != nil {
		ctx := DatabaseProfilerContext{}.new(query, args...)
		callback()
		r.Profiler.Post(ctx)
		return
	}
	callback()
}
