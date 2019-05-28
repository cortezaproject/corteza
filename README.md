[![Build Status](https://drone.crust.tech/api/badges/cortezaproject/corteza/status.svg)](https://drone.crust.tech/cortezaproject/corteza)
[![Go Report Card](https://goreportcard.com/badge/github.com/cortezaproject/corteza-server)](https://goreportcard.com/report/github.com/cortezaproject/corteza-server)

# What is Corteza?

Corteza brings your user ecosystem and essential applications together on one 
platform, unifying them via CRM, Team Messaging and Advanced Identity and 
Access Management.

**Crust CRM**
Crust CRM is the highly flexible, scalable and open source Salesforce alternative, that enables your team to sell faster. It provides a 360 degree view of your customers, enabling you to service your prospects better and detect new opportunities.

**Crust Messaging**
Crust Messaging is the secure, high performance, open source Slack alternative that allows your teams to collaborate more efficiently, as well as communicate safely with other organisations or customers.

**Crust Compose**
Crust Compose is the flexible and easy to use open source Low Code Development platform for custom web based business applications, with drag and drop builder features, integrated Identity, Access and Privacy Management, and powerful automation options. Crust CRM is based on Compose.

**Crust Unify**
Crust Unify manages the user experience for Crust applications, such as CRM, Messaging and Compose, as well as providing an integrated interface for third party or other bespoke applications. 100% responsive and with an intuitive design, Crust Unify increases productivity and ease of access to all IT resources.

# Setup

Copy `.env.example` to `.env` and make proper modifications for your local environment.

An access to a (local) instance of MySQL must be available.
Configure access to your database with `SYSTEM_DB_DSN`, `MESSAGING_DB_DSN` and `COMPOSE_DB_DSN`.

The database will be populated with migrations at the start of each service. You don't need to pre-populate the database, just make sure that your permissions include CREATE and ALTER capabilities.

# Running in local environment for development

Everything should be set and ready to run with `make realize`. This utilizes realize tool that monitors codebase for changes and restarts api http server for every file change. It is not 100% so it needs help (manual restart) in certain cases (new files added, changes in non .go files etc..)

# Making changes

Please refer to each project's style guidelines and guidelines for submitting patches and additions. In general, we follow the "fork-and-pull" Git workflow.

 1. **Fork** the repo on GitHub
 2. **Clone** the project to your own machine
 3. **Commit** changes to your own branch
 4. **Push** your work back up to your fork
 5. Submit a **Pull request** so that we can review your changes

NOTE: Be sure to merge the latest master from "upstream" before making a pull request!
