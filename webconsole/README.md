# Corteza Server Web Console

When enabled (`HTTP_WEB_CONSOLE_ENABLED=true`), it allows insight and management of corteza internals. 

## Running server and web console both from source

See [DEVELOPMENT.md](DEVELOPMENT.md) for details.

## Running server from source and using bundled web console

If you want to use web console while do changes on the server you need to bundle it.
Run `yarn build` and make sure you enable dev mode with `ENVIRONMENT=dev`, restart the server and navigate your browser to `<server>/console/`.  
