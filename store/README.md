# Corteza storage layer

Provides unified storage for Corteza with composable and interchangeable backends.
Backends can be any RDBMS, Key-Value, NoSQL, FS, Memory.

## FAQ

* Why are we changing create/update function signature (input struct is no longer returned)?
** Because store functions are no longer manipulating the input.

* Why naming inconsistency between search/lookup and create/update/...?
** To ensure function names sound more natural

* Why changing find prefix to search/lookup?
** To be consistent with actions

* Why do we use custom mapping (and not db:... tag on struct)?
   Separation of concerns
** consistency with store backends that do not support db tags
** de-cluttering types* namespace

## Testing

Running store tests:
```shell script
make test.store
```

See [Makefile](Makefile) for details

## Known issues

SQLite, transactions & locking::
Transactions in SQLite are explicitly disabled due to issues
with locking (see `sqlite/sqlite.go`, line with `cfg.TxDisabled # true`)

Until this is resolved, SQLite storage backend should not be used
in production environment.

We're still keeping and maintaining it to provide the simplest and most performant backend for (service) unit testing to avoid overly complex mocking of storage.
This will likely change when we have support for built-in in-memory storage (MEMORY_DSN)
