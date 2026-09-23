---
title: Quickstart
description: Run Human with Docker Compose, create the first account and build a small app.
---

# Quickstart

Run Human on your machine with Docker, sign in as the administrator, and build
a small app. It takes about ten minutes.

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
One container runs the Human server and serves the web app. The other runs
PostgreSQL. Files people upload are kept in the `data` volume. Every setting is
an [environment variable](/reference/environment).
:::

## 2. Create your account

Choose **Create a new account**. **The first account on a new
instance becomes the administrator.** It gets every administrative role, so use
an address you control.

There is no email server in this setup, so your address is confirmed straight
away. To send email, set the [`SMTP_*`](/reference/environment#email-sending)
variables.

## 3. Build an app

1. Open **Namespaces** from the menu, or go to `/compose/namespaces`. A
   namespace holds one app.
2. Click **New Namespace**. Name it `Contacts`, give it the short name
   `contacts`, and click **Save**.
3. Open **Modules** and click **New Module**. Call it `Contact`, and use
   **Add new field** to add `Name` (Text), `Email` (Email) and `Company`
   (Text). Click **Save**.
4. Still in the module, click **Create record page**, then **Create record list
   page**. Human creates a form for one contact and a list of all of them.
5. Open the namespace. Your list page is there. Click **New Record** and
   add a contact.

That is a working app: a data model, a list and a form, with permissions
already applied.

## 4. Add an agent

Agents need a model to talk to. Add one first:

1. In **Admin Area**, open **System → LLM Providers** and click **New LLM Provider**.
   Pick **Anthropic** or **Mistral** and paste an API key.
2. Open **Agentic** from the menu and click **Create agent**. Give it a name,
   pick the provider and a model, and write a system prompt such as
   _"You help the team look up contacts."_
3. Under **Access**, set **Works in** to `Contacts` and, in **Tools**, allow
   the **Records** tools. Save, and chat with it right there in the editor.

The [Agents guide](/guides/agents) covers tools, knowledge and invocation in
depth.

## Stop or reset

```sh
docker compose down      # stop; your data is kept
docker compose down -v   # stop and delete all data
```

## Next steps

- [Core concepts](./concepts): the vocabulary the rest of the docs use.
- [Build an app](/guides/build-an-app): charts, dashboards and linked modules.
- [Environment variables](/reference/environment): everything the server can be
  configured with.
