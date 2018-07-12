# Crust

@todo What is Crust

## Contributing

### Setup

Copy `.env.example` to `.env` and make proper modifications for your 
local environment.

An access to a (local) instance of MySQL must be available.
Configure access to your database with `SAM_DB_DSN`.

@todo how to setup crust database

### Running in local environment for development

Everything should be set and ready to run with `make realize`. This
utilizes realize tool that monitors codebase for changes and restarts
api http server for every file change. It is not 100% so it needs help 
(manual restart) in certain cases (new files added, changes in non .go files etc..)
