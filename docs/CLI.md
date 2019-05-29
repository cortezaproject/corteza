# Corteza Command Line Interface

## General

You can choose from one of the four (4) app/service CLI entrypoints:

 - system (`corteza-server-system`)
 - messaging (`corteza-server-messaging`)
 - compose (`corteza-server-compose`)
 - all-in-one monolith (`corteza-server`)

When system is build into a all-in-one binary or image,
app-specific commands can be accessed under a sub-command with the 
same name as the app.

## Config
Command and application behaviour can be adjusted using environmental
variables. Defaults that are pre-set internally are optimized for production mode.

Make sure you read `/.env.example` for all the available options.

Note that defaults shown under `--help` are overridden by value from the accompanying 
environmental variable.

## How to run:

### Compile source files and run with `go run`

Source entrypoints can be found under `cmd/`:

```sh
go run cmd/system/*.go [commands and flags]
```

### Docker or Docker Compose

Docker images are built with handy presets for 
`ENTRYPOINT` and `CMD`. By default, `CMD` is set to `serve-api`.

See docker manual for more info.  

```sh
docker exec -it -rm corteza-server-system [commands and flags]
```

## Base commands:

These commands are accessible from the base level (never 
wrapped under an app-sub-command) and always affect all apps.

### `provision`

Command `provision` Wraps 2 sub-commands: 
 - `access-control-rules`
   Reset access control rules for Everyone and Administrator roles.
 - `migrate-database`
   Run database migration scripts

Both these commands are executed when you run `serve-api` command,
after database is connected and before initialization. This behaviour can
be changed with env variables.

## System specific commands & sub-commands

| Command | Subcommand | |
| --- | --- | --- |
| auth      | auto-discovery     | Auto discovers new OIDC client
| auth      | jwt                | Generates new JWT for a user
| auth      | test-notifications | Sends samples of all authentication notification to receipient
| roles     | useradd            | Add user to role
| settings  | auto-configure     | Run autoconfiguration
| settings  | delete             | Set value (raw JSON) for a specific key
| settings  | get                | Get value (raw JSON) for a specific key
| settings  | import             | Import settings as JSON from stdin or file
| settings  | list               | List all
| settings  | set                | Set value (raw JSON) for a specific key
| users     | add                | Add new user
| users     | list               | List users
| users     | password           | Change user's password

See help (`--help`) under each (sub) command for details about additional flags and arguments.
