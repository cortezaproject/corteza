Goal
 - multi-repository (different DB) support
 - refactor old raw-sql migrations to logical and dynamic schema Upgrading
 - (tbd) consolidation of system & compose (& messaging) subsystems

File system:
 /repository            Holds all repository logic for all subsystems, for all core repository implementation
   /internal            Internal repository tools (pkg/ql, pgk/rh should be moved here)
 /<implementation>      Individual core repository implementation
                        [mysql|postgresql|redis|memory|sqlite|elasticsearch|mongo]
   /schema              Schema Upgrading for individual repository implementation





/corteza/store/rdbms
    Basic RDBMS logic that should be reused across all
    RDBMS-like implementations

/corteza/store/<rdbms-store-type>/sql_*.go
        Database table definitions (used by Upgrade_*.go)

/corteza/store/<rdbms-store-type>/Upgrade_*.go
    Upgrade logic, cascading (eg: Upgrade => UpgradeSystem => UpgradeUsers)



FAQ:
Why are we changing create/update function signature (input struct is no longer returned)?
Because store functions are no longer manipulating the input.

Why naming inconsistency between search/lookup and create/update/...?
To ensure function names sound more natural

Why changing find prefix to search/lookup?
To be consistent with actions

Why do we use custom mapping (and not db:... tag on struct)?
Separation of concerns
 + consistency with store backends that do not support db tags
 + de-cluttering types* namespace

