# Human Corredor

Corredor is Human's automation **script runner and bundler**. It loads server and
client automation scripts from extension directories, parses their trigger and
security definitions, executes server scripts on request, and packs client scripts
into browser bundles with webpack.

The Human server talks to Corredor over gRPC, using two services defined in
[`def/protobuf/service-corredor.proto`](../def/protobuf/service-corredor.proto):

- `ServerScripts` — `List`, `Exec`
- `ClientScripts` — `List`, `Bundle`

The script API is Corteza-compatible: scripts written for Corteza Corredor run
here unchanged, and every `CORREDOR_*` environment variable keeps its Corteza name
and default.

## Running

```sh
pnpm --filter @planetcrust/human-corredor serve        # run
pnpm --filter @planetcrust/human-corredor serve:dev    # run, restart on change
pnpm --filter @planetcrust/human-corredor test:unit    # unit tests
pnpm --filter @planetcrust/human-corredor lint:check   # lint
```

Corredor needs to reach the Human API to start, so point it at a running server:

```sh
CORREDOR_EXEC_CSERVERS_API_HOST=localhost:1343 \
CORREDOR_EXEC_CSERVERS_API_BASEURL_TEMPLATE='http://{host}/api/{service}' \
pnpm --filter @planetcrust/human-corredor serve
```

An `.env` file in this directory is read on start; see `.env.example`.

[`SCRIPTING.md`](SCRIPTING.md) is the reference for writing one: every resource
a trigger can bind to, the arguments each event hands the script, the constraint
names, and what the execution context exposes. It is generated, so it tracks the
server rather than drifting from it.

## Extension layout

Extensions live under the paths in `CORREDOR_EXT_SEARCH_PATHS`:

```
usr/
    <extension>/
        package.json              # optional, dependencies installed on start
        server-scripts/           # *.js, any nesting
            hello.js
        client-scripts/
            <bundle>/             # one webpack bundle per directory
                index.js
```

A script exports its definition as the default export:

```js
export default {
  label: 'Hello',
  description: 'What this script does',
  triggers({ on }) {
    return on('manual').for('compose:record')
  },
  security: { runAs: 'someone@example.tld', allow: ['everyone'] },
  async exec(args, ctx) {
    return { hello: 'world' }
  },
}
```

`triggers` may return one trigger, an array, or be a generator yielding several.
`iterator` replaces `triggers` for resource iterators. Files named `*.test.js` and
anything under `node_modules` are skipped.

## Environment

Two variables are specific to Human:

| Variable                              | Default               | Notes                                                                                                 |
| ------------------------------------- | --------------------- | ----------------------------------------------------------------------------------------------------- |
| `CORREDOR_PROTOBUF_PATH`              | `<repo>/def/protobuf` | Directory holding `service-corredor.proto`. `CORREDOR_CORTEZA_PROTOBUF_PATH` is accepted as an alias. |
| `CORREDOR_EXT_DEPENDENCIES_INSTALLER` | yarn, else npm        | Package manager for extension dependencies: `yarn`, `npm` or `pnpm`. Auto-detected when unset.        |

Everything else keeps its Corteza name and default:

| Variable                                                | Default                        | Notes                                                                                                          |
| ------------------------------------------------------- | ------------------------------ | -------------------------------------------------------------------------------------------------------------- |
| `CORREDOR_ENVIRONMENT`                                  | `prod`                         | `dev*` turns on pretty logs and `trace` level. `CORREDOR_ENV` and `NODE_ENV` are fallbacks.                    |
| `CORREDOR_ADDR`                                         | `localhost:50051`              | gRPC listen address.                                                                                           |
| `CORREDOR_SERVER_CERTIFICATES_ENABLED`                  | on in prod                     | Serve gRPC over TLS.                                                                                           |
| `CORREDOR_SERVER_CERTIFICATES_PATH`                     | `/certs`                       | Holds `ca.crt`, `private.key`, `public.crt`. Individual paths: `..._CA`, `..._PRIVATE`, `..._PUBLIC`.          |
| `CORREDOR_LOG_ENABLED`                                  | `true`                         |                                                                                                                |
| `CORREDOR_LOG_PRETTY`                                   | on in dev                      | Pretty multiline logs instead of JSON.                                                                         |
| `CORREDOR_LOG_LEVEL`                                    | `info` (`trace` in dev)        | `fatal`, `error`, `warn`, `info`, `debug`, `trace`, `silent`.                                                  |
| `CORREDOR_EXEC_CSERVERS_API_HOST`                       | —                              | API host; `DOMAIN`, `HOSTNAME`, `HOST` are fallbacks.                                                          |
| `CORREDOR_EXEC_CSERVERS_API_BASEURL_TEMPLATE`           | `https://api.{host}/{service}` | `{host}` and `{service}` are substituted.                                                                      |
| `CORREDOR_EXEC_CTX_CORTEZA_SERVERS_SYSTEM_API_BASEURL`  | derived                        | Overrides the system API base URL.                                                                             |
| `CORREDOR_EXEC_CTX_CORTEZA_SERVERS_COMPOSE_API_BASEURL` | derived                        | Overrides the compose API base URL.                                                                            |
| `CORREDOR_EXEC_CTX_FRONTEND_BASEURL`                    | derived                        | Frontend base URL handed to scripts; `PROVISION_SETTINGS_AUTH_FRONTEND_URL_BASE` is a fallback.                |
| `CORREDOR_EXT_SEARCH_PATHS`                             | `usr:usr/*`                    | Colon separated; relative paths resolve against this directory.                                                |
| `CORREDOR_EXT_DEPENDENCIES_AUTO_UPDATE`                 | `true`                         | Install extension dependencies on start and on `package.json` change.                                          |
| `CORREDOR_EXT_SERVER_SCRIPTS_ENABLED`                   | `true`                         |                                                                                                                |
| `CORREDOR_EXT_SERVER_SCRIPTS_WATCH`                     | `true`                         | Reload server scripts on change.                                                                               |
| `CORREDOR_EXT_CLIENT_SCRIPTS_ENABLED`                   | `true`                         |                                                                                                                |
| `CORREDOR_EXT_CLIENT_SCRIPTS_WATCH`                     | `true`                         | Rebundle client scripts on change.                                                                             |
| `CORREDOR_BUNDLER_ENABLED`                              | `true`                         | Disabled automatically when the output path is not writable.                                                   |
| `CORREDOR_BUNDLER_OUTPUT_PATH`                          | `/tmp/corredor/bundler-dist`   | Where client-script bundles are written.                                                                       |
| `SENTRY_DSN`                                            | —                              | Enables Sentry. `SENTRY_SERVERNAME`, `SENTRY_RELEASE`, `SENTRY_DIST`, `SENTRY_ENVIRONMENT` are passed through. |

## TLS certificates

`certs/Makefile` generates a CA plus server and client pairs for local testing:

```sh
cd certs && make all_certs
```

The generated directories are git-ignored; deployments mount their own into
`CORREDOR_SERVER_CERTIFICATES_PATH`.

## Container image

The build context is the repository root, because the image needs the pnpm
workspace manifests, `lib/js` and `def/protobuf` alongside this package:

```sh
docker build -f corredor/Dockerfile -t human-corredor .
```

Inside the image the workspace is `/app`, with the package at `/app/corredor`
and the protobuf definitions at `/app/def/protobuf`. `/corredor/usr` (scripts)
and `/corredor/certs` (TLS certificates) are volumes.
