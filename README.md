# What is Corteza?

[![Build Status](https://drone.crust.tech/api/badges/cortezaproject/corteza/status.svg)](https://drone.crust.tech/cortezaproject/corteza)
[![Go Report Card](https://goreportcard.com/badge/github.com/cortezaproject/corteza-server)](https://goreportcard.com/report/github.com/cortezaproject/corteza-server)

Corteza brings your user ecosystem and essential applications together on one platform, unifying them via CRM, Team Messaging and Advanced Identity and Access Management.

**Corteza Messaging** is a secure, high performance, open source Slack alternative that allows your teams to collaborate more efficiently, as well as communicate safely with other organisations or customers.

**Corteza Compose** is an open source Rapid Application Development (RAD) platform for custom web based business applications. Deliver the application you need more easily and faster then ever before with the drag and drop page builder, protect users with integrated Identity, Access and Privacy Management, and automate tasks with Compose’s advanced automation functionality. Corteza Compose is easy, fast and secure – your perfect ally to digitize your organisation’s business processes and customer engagement.

**Corteza Unify** manages user experience for Corteza applications, such as Compose and Messaging, as well as providing an integrated interface for third party or other bespoke applications. 100% responsive and with an intuitive design, Corteza Unify increases productivity and ease of access to all IT resources.

## Contributing

### Setup

Copy `.env.example` to `.env` and make proper modifications for your local environment.

An access to a (local) instance of MySQL must be available.
Configure access to your database with `SYSTEM_DB_DSN`, `MESSAGING_DB_DSN` and `COMPOSE_DB_DSN`.

The database will be populated with migrations at the start of each service. You don't need to pre-populate the database, just make sure that your permissions include CREATE and ALTER capabilities.

### Running in local environment for development

Everything should be set and ready to run with `make realize`. This utilizes realize tool that monitors codebase for changes and restarts api http server for every file change. It is not 100% so it needs help (manual restart) in certain cases (new files added, changes in non .go files etc..)

### Making changes

Please refer to each project's style guidelines and guidelines for submitting patches and additions. In general, we follow the "fork-and-pull" Git workflow.

 1. **Fork** the repo on GitHub
 2. **Clone** the project to your own machine
 3. **Commit** changes to your own branch
 4. **Push** your work back up to your fork
 5. Submit a **Pull request** so that we can review your changes

NOTE: Be sure to merge the latest master from "upstream" before making a pull request!
