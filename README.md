# What is CRUST?

**CRUST Messaging** is a high performance, self-hosted, open source Slack alternative. It has an API centric design and all data exchange is in JSON format. CRUST Messaging is tightly coupled with CRUST IAM.

**CRUST CRM** is a scalable, self-hosted, open source Salesforce alternative. It provides a suite of tools to build API centric microservice modules. CRUST CRM can be loosely coupled with CRUST Messaging for the purposes of customer engagement and addition of rich external data sources. All data exchange is in JSON format. Business Logic, Workflow, Search and AI (later!) can be applied across both CRM and Messaging stores. CRUST CRM is tightly coupled with CRUST IAM.

**CRUST IAM** is an advanced Identity and Access Management infrastructure which includes:
 - Social Logins
 - Single Sign On
 - Multi-factor Authentication
 - User Identity Self-Service
 - User and Organisational Privacy Controls
 - Standardised ANSI Role Based Access Control

**CRUST Client** is a multi-functional client which unifies the user experience of CRUST IAM, Messaging and CRM, allowing organisations to extend access to third party applications internal and external to their firewall (e.g. Video, Docs, Dev Tools) via a common user interface. The design is inspired by popular browser UX, though not identical. CRUST Client will be available on web, desktop and mobile platforms.

## Contributing

### Setup

Copy `.env.example` to `.env` and make proper modifications for your 
local environment.

An access to a (local) instance of MySQL must be available.
Configure access to your database with `SAM_DB_DSN` and `CRM_DB_DSN`.

Please check the options available with `./app -h`.

The database will be populated with migrations at the start of each
service. You don't need to pre-populate the database, just make sure
that your permissions include CREATE and ALTER capabilities.

### Running in local environment for development

Everything should be set and ready to run with `make realize`. This
utilizes realize tool that monitors codebase for changes and restarts
api http server for every file change. It is not 100% so it needs help 
(manual restart) in certain cases (new files added, changes in non .go files etc..)

### Making changes

Please refer to each project's style guidelines and guidelines for submitting patches and additions.
In general, we follow the "fork-and-pull" Git workflow.

 1. **Fork** the repo on GitHub
 2. **Clone** the project to your own machine
 3. **Commit** changes to your own branch
 4. **Push** your work back up to your fork
 5. Submit a **Pull request** so that we can review your changes

NOTE: Be sure to merge the latest master from "upstream" before making a pull request!