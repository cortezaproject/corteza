# Testing RDBMS drivers 

Package provides helper functions for writing tests for Corteza's RDBMS implementations (`store/adapters/rdbms/drivers` package). 

Please note that this is currently set-up to support RDBMS only and should not be used for other kinds of store implementations.

# Configuration

By default, in-memory SQLite is used if no other databases are configured.

Tests scan environmental variables and use all variables with `TEST_STORE_ADAPTERS_RDBMS_` prefix.

See `.env.example` file for more details.

