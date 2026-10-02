# End-to-end tests

Playwright specs that drive the web applications in a real browser against a
running Corteza stack. They attach to whatever is running; they never start or
seed anything themselves.

## Stack

1. A server with a throwaway database. Start it with these on top of your usual
   environment:

   ```
   DB_DSN=postgres://user:pass@localhost:5432/corteza_e2e
   HTTP_ADDR=:1043
   DOMAIN=localhost:1043
   AUTH_REQUEST_RATE_LIMIT=0
   AUTH_DEFAULT_REDIRECT_URIS=http://127.0.0.1:8081 http://localhost:8081
   PROVISION_ALWAYS=true
   ```

   Enable internal login and create the test user once:

   ```
   corteza settings set auth.internal.enabled true   # then restart the server
   corteza users add e2e-admin@example.tld --password '<pass>' --role super-admin
   corteza import e2e/fixtures/
   ```

2. The compose webapp dev server pointed at that API
   (`public/config.js`: `window.CortezaAPI = 'http://localhost:1043/api'`):

   ```
   cd client/web/compose && yarn serve --port 8081
   ```

3. `cp .env.e2e.example .env.e2e` and fill in the user.

## Running

```
cd e2e
npm install
npx playwright test            # all specs
npx playwright test --ui       # interactive
```

The `setup` project logs in once and stores the session in `.auth/`; every
other project reuses it. Traces and screenshots of failures land in `.results/`.

## Fixtures

`fixtures/` is an envoy YAML bundle (namespace, modules, records, pages) that
the specs rely on. Import it with `corteza import e2e/fixtures/`.
