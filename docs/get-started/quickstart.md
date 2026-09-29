---
title: Quickstart
description: Run Human with Docker Compose and sign in as the administrator.
---

# Quickstart

Run Human on your machine with Docker and sign in as the administrator. It
takes about five minutes.

You need [Docker](https://docs.docker.com/get-docker/) with Compose.

## 1. Start Human

Save this as `docker-compose.yml` in an empty folder, or
[download it](/docker-compose.yml).

<<< @/public/docker-compose.yml{yaml}

Then start it:

```sh
docker compose up -d
```

The first start creates the database and takes a minute. Once
`docker compose ps` reports the `human` service as healthy, open
[http://localhost:8080](http://localhost:8080).

::: tip What is in the file
One container runs the Human server and serves the web app, another runs
PostgreSQL, and a third, Corredor, runs automation scripts. Files people upload
are kept in the `data` volume. Every setting is an
[environment variable](/reference/environment).
:::

::: tip Automation scripts
Corredor loads scripts from the `scripts` folder next to
`docker-compose.yml`. Put each extension in its own folder, with its scripts
under `server-scripts/` and `client-scripts/`:

```txt
scripts/
  my-extension/
    server-scripts/
    client-scripts/
```

Copy the folder there with `git clone`, `rsync` or anything else, then run
`docker compose restart corredor` to load the changes.
:::

## 2. Create your account

Choose **Create a new account**. **The first account on a new
instance becomes the administrator.** It gets every administrative role, so use
an address you control.

There is no email server in this setup, so your address is confirmed straight
away. To send email, set the [`SMTP_*`](/reference/environment#email-sending)
variables.

## 3. Look around

<Screenshot name="home" alt="The home screen after signing in" />

After you sign in, the home screen shows three columns: the **Menu** of apps,
the **Assistant** for chatting with agents, and your **Notifications**. A new
instance lists these apps:

| App                  | What it is for                                                   |
| -------------------- | ---------------------------------------------------------------- |
| **Namespaces**       | Build apps: data models, pages and charts                        |
| **Automation (TAQ)** | Automations that react to events, run on schedules or on demand  |
| **Agentic**          | AI agents, their models, tools and access                        |
| **Chatbot**          | Chatbots for your website, and the inbox for their conversations |
| **Projects**         | Group related resources into one governed unit                   |
| **Admin Area**       | Users, roles, permissions and system settings                    |

[The platform](/platform/) describes each of them.

## Stop or reset

```sh
docker compose down      # stop; your data is kept
docker compose down -v   # stop and delete all data
```

## Next steps

- [Core concepts](./concepts): the vocabulary the rest of the docs use.
- [How building works](./building): how data, pages, automation and agents fit
  together.
- [Environment variables](/reference/environment): everything the server can be
  configured with.
