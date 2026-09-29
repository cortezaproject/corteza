---
title: Upgrading from Corteza
description: What changes when you move a Corteza installation to Human, starting with Corredor and its automation scripts.
---

# Upgrading from Corteza

Human is built on Corteza. This page lists what you need to change when you
move a Corteza installation to Human. It covers Corredor so far.

## Corredor

Corredor runs automation scripts. Your scripts need no changes: the script API,
the events scripts bind to and every `CORREDOR_*` environment variable keep
their Corteza names and defaults. What changes is the container around them.

### 1. Switch the image

Replace `cortezaproject/corteza-server-corredor` with
`planetcrust/human-corredor`. Keep the service named `corredor`: the Human
server looks for it at `corredor:80`, as Corteza did.

### 2. Point Corredor at the Human server

The Corredor image expects the Human server to be the `human` service in the
same Compose file. Corteza named it `server`. Either rename the service to
`human`, or keep `server` and set this on the `corredor` service:

```yaml
CORREDOR_EXEC_CSERVERS_API_BASEURL_TEMPLATE: http://server/api/{service}
```

Without one of the two, scripts cannot call the API.

### 3. Mount your scripts

Mount the folder that holds your extensions at `/corredor/usr`. Corredor loads
every extension in it, each with its scripts under `server-scripts/` and
`client-scripts/`:

```txt
scripts/
  my-extension/
    server-scripts/
    client-scripts/
```

You no longer need to list the extension in `CORREDOR_EXT_SEARCH_PATHS`, so
you can remove that line from your `.env`. If you keep it, it still works: it
replaces the default search path.

### 4. Move client scripts for other webapps

Corteza had a separate webapp for each area. Human has one, and it loads
client scripts only from these folders under `client-scripts/`:

| Folder     | Holds scripts for                              |
| ---------- | ---------------------------------------------- |
| `compose/` | Namespaces, pages and records                  |
| `admin/`   | The Admin Area                                 |
| `unify/`   | Scripts that belong to no single part of Human |

Scripts in any other folder, such as one written for Corteza's workflow
webapp, are not loaded. Move them into one of the three.

### 5. Enable Corredor on the server

Set `CORREDOR_ENABLED=true` on the Human server, as in Corteza.

### Example

The `corredor` service and the setting it needs on the server, with your
scripts in a `scripts` folder next to `docker-compose.yml`:

```yaml
services:
  human:
    image: planetcrust/human:latest
    environment:
      CORREDOR_ENABLED: 'true'
      # … your other settings

  corredor:
    image: planetcrust/human-corredor:latest
    restart: unless-stopped
    environment:
      CORREDOR_EXEC_CTX_FRONTEND_BASEURL: https://your-human-domain.example
    volumes:
      - ./scripts:/corredor/usr
```

### Check that it worked

Run `docker compose restart corredor`, then `docker compose logs corredor`.
Corredor logs a `processed` line for server scripts and one for client
scripts. Each gives `valid` and `total` counts. When a script fails to parse,
the two counts differ.
