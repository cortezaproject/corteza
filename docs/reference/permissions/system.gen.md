---
title: System permissions
outline: [2, 2]
---

<!-- This file is auto-generated from server/app and the component definitions. -->

# System permissions

Every operation is denied unless a role is granted it.

## Component

| Operation | Description |
| --- | --- |
| `grant` | Manage system permissions |
| `action-log.read` | Access to action log |
| `settings.read` | Read system settings |
| `settings.manage` | Manage system settings |
| `auth-client.create` | Create auth clients |
| `auth-clients.search` | List, search or filter auth clients |
| `role.create` | Create roles |
| `roles.search` | List, search or filter roles |
| `user-group.create` | Create user groups |
| `user-groups.search` | List, search or filter user groups |
| `user.create` | Create users |
| `users.search` | List, search or filter users |
| `dal-connection.create` | Create DAL connections |
| `dal-connections.search` | List, search or filter DAL connections |
| `dal-sensitivity-level.manage` | Can manage DAL sensitivity levels |
| `labels.search` | List, search or filter labels |
| `corredor-scripts.search` | List, search or filter Corredor scripts |
| `application.create` | Create applications |
| `applications.search` | List, search or filter auth clients |
| `application.flag.self` | Manage private flags for applications |
| `application.flag.global` | Manage global flags for applications |
| `template.create` | Create template |
| `templates.search` | List, search or filter templates |
| `report.create` | Create report |
| `reports.search` | List, search or filter reports |
| `reminder.assign` |  Assign reminders |
| `queue.create` | Create messagebus queues |
| `queues.search` | List, search or filter messagebus queues |
| `apigw-route.create` | Create API gateway route |
| `apigw-routes.search` | List search or filter API gateway routes |
| `resource-translations.manage` | List, search, create, or update resource translations |
| `dal-schema-alterations.manage` | List, search, apply, or dismiss DAL alterations |
| `connection.create` | Create connections |
| `connections.search` | List, search or filter connections |
| `configured-connection.create` | Install connection connections |
| `configured-connections.search` | List, search or filter connection connections |
| `data-privacy-request.create` | Create data privacy requests |
| `data-privacy-requests.search` | List, search or filter data privacy requests |
| `notification.assign` | Assign notifications to other users |
| `llm-provider.create` | Create LLM providers |
| `llm-providers.search` | List, search or filter LLM providers |
| `agent.create` | Create agents |
| `agents.search` | List, search or filter agents |
| `project-incident.create` | Create project incidents |
| `project-incidents.search` | List, search or filter project incidents |
| `project-feature.create` | Create project features |
| `project-features.search` | List, search or filter project features |
| `project-privacy.create` | Create project privacy items |
| `project-privacys.search` | List, search or filter project privacy items |
| `project-task.create` | Create project tasks |
| `project-tasks.search` | List, search or filter project tasks |
| `project-review.create` | Create project reviews |
| `project-reviews.search` | List, search or filter project reviews |
| `project-backlog-item.create` | Create project backlog items |
| `project-backlog-items.search` | List, search or filter project backlog items |
| `ai-conversation.create` | Create AI conversations |
| `ai-conversations.search` | List, search or filter AI conversations |
| `knowledge-base.create` | Create knowledge bases |
| `knowledge-bases.search` | List, search or filter knowledge bases |
| `chatbot.create` | Create chatbots |
| `chatbots.search` | List, search or filter chatbots |
| `chatbot-session.create` | Create chatbot sessions |
| `chatbot-sessions.search` | List, search or filter chatbot sessions |
| `chatbot-session-step.create` | Create chatbot session steps |
| `chatbot-session-steps.search` | List, search or filter chatbot session steps |
| `chatbot-session-handoff.create` | Create chatbot session handoffs |
| `chatbot-session-handoffs.search` | List, search or filter chatbot session handoffs |
| `tenant.create` | Create tenants |
| `tenants.search` | List, search or filter tenants |
| `project.create` | Create projects |
| `projects.search` | List, search or filter projects |
| `project-ai-system.create` | Create AI systems |
| `project-ai-systems.search` | List, search or filter AI systems |
| `project-fria-scenario.create` | Create FRIA risk scenarios |
| `project-fria-scenarios.search` | List, search or filter FRIA risk scenarios |

## Application

| Operation | Description |
| --- | --- |
| `read` | Read application |
| `access` | Access application |
| `update` | Update application |
| `delete` | Delete application |
| `source.manage` | Replace the HTML source of a custom application |

## Apigw Route

| Operation | Description |
| --- | --- |
| `read` | Read API Gateway route |
| `update` | Update API Gateway route |
| `delete` | Delete API Gateway route |

## Auth Client

| Operation | Description |
| --- | --- |
| `read` | Read authorization client |
| `update` | Update authorization client |
| `delete` | Delete authorization client |
| `authorize` | Authorize authorization client |

## Data Privacy Request

| Operation | Description |
| --- | --- |
| `read` | Read data privacy request |
| `approve` | Approve/Reject data privacy request |

## Queue

| Operation | Description |
| --- | --- |
| `read` | Read queue |
| `update` | Update queue |
| `delete` | Delete queue |
| `queue.read` | Read from queue |
| `queue.write` | Write to queue |

## Report

| Operation | Description |
| --- | --- |
| `read` | Read report |
| `update` | Update report |
| `delete` | Delete report |
| `run` | Run report |

## Role

| Operation | Description |
| --- | --- |
| `read` | Read role |
| `update` | Update role |
| `delete` | Delete role |
| `members.manage` | Manage members |

## User Group

| Operation | Description |
| --- | --- |
| `read` | Read user group |
| `update` | Update user group |
| `delete` | Delete user group |
| `members.manage` | Manage members |

## Template

| Operation | Description |
| --- | --- |
| `read` | Read template |
| `update` | Update template |
| `delete` | Delete template |
| `render` | Render template |

## User

| Operation | Description |
| --- | --- |
| `read` | Read user |
| `update` | Update user |
| `delete` | Delete user |
| `suspend` | Suspend user |
| `unsuspend` | Unsuspend user |
| `email.unmask` | Unmask email |
| `name.unmask` | Unmask name |
| `impersonate` | Impersonate user |
| `credentials.manage` | Manage user's credentials |

## Dal Connection

| Operation | Description |
| --- | --- |
| `read` | Read connection |
| `update` | Update connection |
| `delete` | Delete connection |
| `dal-config.manage` | Manage DAL configuration |

## Connection

| Operation | Description |
| --- | --- |
| `read` | Read connection |
| `update` | Update connection |
| `delete` | Delete connection |
| `install` | Install connection |

## Configured Connection

| Operation | Description |
| --- | --- |
| `read` | Read connection |
| `delete` | Delete connection |
| `update` | Update connection |

## Llm Provider

| Operation | Description |
| --- | --- |
| `read` | Read LLM provider |
| `update` | Update LLM provider |
| `delete` | Delete LLM provider |

## Agent

| Operation | Description |
| --- | --- |
| `read` | Read agent |
| `update` | Update agent |
| `delete` | Delete agent |

## Ai Conversation

| Operation | Description |
| --- | --- |
| `read` | Read AI conversation |
| `update` | Update AI conversation |
| `delete` | Delete AI conversation |

## Knowledge Base

| Operation | Description |
| --- | --- |
| `read` | Read knowledge base |
| `update` | Update knowledge base |
| `delete` | Delete knowledge base |

## Chatbot

| Operation | Description |
| --- | --- |
| `read` | Read chatbot |
| `update` | Update chatbot |
| `delete` | Delete chatbot |
| `sessions.view` | View chatbot sessions |
| `sessions.manage` | Manage chatbot sessions (advance, close) |
| `sessions.handoff.manage` | Manage chatbot session handoffs |

## Chatbot Session

| Operation | Description |
| --- | --- |
| `read` | Read chatbot session |
| `update` | Update chatbot session |
| `delete` | Delete chatbot session |

## Chatbot Session Step

| Operation | Description |
| --- | --- |
| `read` | Read chatbot session step |
| `update` | Update chatbot session step |
| `delete` | Delete chatbot session step |

## Chatbot Session Handoff

| Operation | Description |
| --- | --- |
| `read` | Read chatbot session handoff |
| `update` | Update chatbot session handoff |
| `delete` | Delete chatbot session handoff |

## Tenant

| Operation | Description |
| --- | --- |
| `read` | Read tenant |
| `update` | Update tenant |
| `delete` | Delete tenant |
| `suspend` | Suspend tenant |
| `members.manage` | Manage tenant members |

## Project

| Operation | Description |
| --- | --- |
| `read` | Read project |
| `update` | Update project |
| `delete` | Delete project |
| `members.manage` | Manage project members |
| `revise` | Create a revision of a project |
| `publish` | Publish a project revision |

## Project Ai System

| Operation | Description |
| --- | --- |
| `read` | Read AI system |
| `update` | Update AI system |
| `delete` | Delete AI system |
| `resources.manage` | Manage AI system resources |

## Project Fria Scenario

| Operation | Description |
| --- | --- |
| `read` | Read FRIA risk scenario |
| `update` | Update FRIA risk scenario |
| `delete` | Delete FRIA risk scenario |

## Project Incident

| Operation | Description |
| --- | --- |
| `read` | Read project incident |
| `update` | Update project incident |
| `delete` | Delete project incident |

## Project Feature

| Operation | Description |
| --- | --- |
| `read` | Read project feature |
| `update` | Update project feature |
| `delete` | Delete project feature |

## Project Privacy

| Operation | Description |
| --- | --- |
| `read` | Read project privacy |
| `update` | Update project privacy |
| `delete` | Delete project privacy |

## Project Task

| Operation | Description |
| --- | --- |
| `read` | Read project task |
| `update` | Update project task |
| `delete` | Delete project task |

## Project Review

| Operation | Description |
| --- | --- |
| `read` | Read project review |
| `update` | Update project review |
| `delete` | Delete project review |

## Project Backlog Item

| Operation | Description |
| --- | --- |
| `read` | Read project backlog item |
| `update` | Update project backlog item |
| `delete` | Delete project backlog item |
