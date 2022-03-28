# RDBMS Schema package

Package contains 

1. DDL 
2. Current schema definition instructions
3. Upgrade instructions using state detection and modification instructions

## DDL or Data Definition Language

This package provides building blocks and generic dialect for specifying RDBMS schema creation and alternation instructions.

## Flow

### General steps
1. Application (or tests) start the boot procedure 
2. `store.Upgrade()` (`store/upgrade.go`) is called with current primary store connection

### RDBMS specific steps

RDBMS implementation of schema upgrade is unified across all drivers (store implementations)
3. `schema.Upgrade()` initializes common (across implementations) upgrade procedure and use `schemaUpgrader` interface that each store implementation provides.
4. 
