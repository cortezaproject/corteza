# Corteza → Human Feature Migration TODO

Comparison of all features between **Corteza (Vue2)** and **Human (Vue3)**.
Corteza is the source of truth — all its features should eventually be in Human.

Legend: ✅ Implemented · ⚠️ Partial · ❌ Not implemented · 🆕 Human-only

---

## Admin

| Feature                                        | Corteza | Human | Notes                    |
| ---------------------------------------------- | ------- | ----- | ------------------------ |
| User list, create, edit                        | ✅      | ✅    | Added admin users screen |
| User avatar/password/MFA management            | ✅      | ❌    |                          |
| User import/export                             | ✅      | ❌    |                          |
| Role list, create, edit                        | ✅      | ❌    |                          |
| Role member management                         | ✅      | ❌    |                          |
| User group list, create, edit                  | ✅      | ❌    |                          |
| System settings (auth, email, branding)        | ✅      | ❌    |                          |
| SAML / OIDC external auth config               | ✅      | ❌    |                          |
| Applications list, create, edit                | ✅      | ❌    |                          |
| API Gateway route management                   | ✅      | ❌    |                          |
| API Gateway profiler                           | ✅      | ❌    |                          |
| Auth clients (OAuth2/OIDC) management          | ✅      | ❌    |                          |
| Data connections list, create, edit            | ✅      | ✅    | Human admin covers this  |
| Data sources list, create, edit                | ✅      | ✅    | Human admin covers this  |
| Email templates editor                         | ✅      | ❌    |                          |
| Code snippets editor                           | ✅      | ❌    |                          |
| Sensitivity levels management                  | ✅      | ❌    |                          |
| Message queues management                      | ✅      | ❌    |                          |
| Automation/workflow management (admin)         | ✅      | ❌    | TAQ is separate app      |
| Automation sessions management                 | ✅      | ❌    |                          |
| Compose module settings                        | ✅      | ❌    |                          |
| Federation management                          | ✅      | ❌    |                          |
| Theming/navigation/location settings (UI mgmt) | ✅      | ❌    |                          |
| System audit/action log                        | ✅      | ❌    |                          |
| Admin dashboard with metrics                   | ✅      | ✅    | Basic dashboard exists   |
| System-wide permissions editor                 | ✅      | ❌    |                          |
| RBAC permission cloning                        | ✅      | ❌    |                          |

---

## Compose

| Feature                                           | Corteza | Human | Notes                                                                  |
| ------------------------------------------------- | ------- | ----- | ---------------------------------------------------------------------- |
| Namespace list view                               | ✅      | ✅    |                                                                        |
| Namespace create/edit/delete                      | ✅      | ✅    |                                                                        |
| Namespace clone                                   | ✅      | ⚠️    | API client exists but no UI to trigger it                              |
| Namespace import/export                           | ✅      | ❌    |                                                                        |
| Namespace translation support                     | ✅      | ❌    |                                                                        |
| Reminders within namespace                        | ✅      | ❌    |                                                                        |
| Module list, create, edit, delete                 | ✅      | ✅    |                                                                        |
| Field type: String                                | ✅      | ✅    |                                                                        |
| Field type: Number                                | ✅      | ✅    |                                                                        |
| Field type: Boolean                               | ✅      | ✅    |                                                                        |
| Field type: DateTime                              | ✅      | ✅    |                                                                        |
| Field type: Select                                | ✅      | ✅    |                                                                        |
| Field type: Email                                 | ✅      | ✅    |                                                                        |
| Field type: URL                                   | ✅      | ✅    |                                                                        |
| Field type: Record (relationship)                 | ✅      | ❌    | Not registered in lib/vue/src/components/field/registry.ts             |
| Field type: User                                  | ✅      | ❌    | Not registered in lib/vue/src/components/field/registry.ts             |
| Field type: File                                  | ✅      | ❌    |                                                                        |
| Field type: Geometry                              | ✅      | ❌    |                                                                        |
| Field validation configuration                    | ✅      | ⚠️    | Only DateTime has Past/Future constraints; no general validator system |
| Field translation/localization                    | ✅      | ❌    |                                                                        |
| DAL field encoding / schema alterations           | ✅      | ❌    |                                                                        |
| Data privacy per field (sensitivity)              | ✅      | ⚠️    | Privacy config exists, not field-level                                 |
| Record revision tracking                          | ✅      | ❌    |                                                                        |
| Federation per module                             | ✅      | ❌    |                                                                        |
| Module discovery settings                         | ✅      | ❌    |                                                                        |
| Unique value constraints                          | ✅      | ❌    |                                                                        |
| Record list with pagination                       | ✅      | ✅    |                                                                        |
| Record create/edit/delete                         | ✅      | ✅    |                                                                        |
| Record bulk edit                                  | ✅      | ❌    |                                                                        |
| Record import from file                           | ✅      | ❌    |                                                                        |
| Record export to file                             | ✅      | ❌    |                                                                        |
| Record organizer (visual hierarchy)               | ✅      | ❌    |                                                                        |
| Record tree view                                  | ✅      | ❌    |                                                                        |
| Record revision history view                      | ✅      | ❌    |                                                                        |
| Advanced filtering / filter presets               | ✅      | ⚠️    | Basic filtering exists; presets/advanced missing                       |
| Draft management (auto-save)                      | ✅      | ❌    |                                                                        |
| Bulk operations toolbar                           | ✅      | ❌    |                                                                        |
| Record translation (multilingual content)         | ✅      | ❌    |                                                                        |
| Page list, create, edit, delete                   | ✅      | ✅    |                                                                        |
| Page hierarchy (parent-child)                     | ✅      | ✅    |                                                                        |
| Page drag-drop reorder/reparent                   | ✅      | ✅    |                                                                        |
| Page translation                                  | ✅      | ❌    |                                                                        |
| Page builder (visual, drag-drop)                  | ✅      | ⚠️    | Route exists, listed as WIP                                            |
| Page attachments management                       | ✅      | ❌    |                                                                        |
| Page layout templates                             | ✅      | ✅    | `usePageLayoutStore` exists                                            |
| Block: Record (single record editor)              | ✅      | ✅    |                                                                        |
| Block: Record List                                | ✅      | ✅    |                                                                        |
| Block: Content (rich text/HTML)                   | ✅      | ✅    |                                                                        |
| Block: Page (container block)                     | ✅      | ✅    |                                                                        |
| Block: Record Organizer                           | ✅      | ❌    |                                                                        |
| Block: Record Revisions                           | ✅      | ❌    |                                                                        |
| Block: Chart/Report                               | ✅      | ❌    | Chart types implemented in lib but not wired into page block registry  |
| Block: Calendar                                   | ✅      | ❌    |                                                                        |
| Block: Geometry/Map                               | ✅      | ❌    |                                                                        |
| Block: Metric                                     | ✅      | ❌    |                                                                        |
| Block: Progress                                   | ✅      | ❌    |                                                                        |
| Block: File                                       | ✅      | ❌    |                                                                        |
| Block: IFrame                                     | ✅      | ❌    |                                                                        |
| Block: Automation (trigger buttons)               | ✅      | ❌    |                                                                        |
| Block: Tabs                                       | ✅      | ❌    |                                                                        |
| Block: Navigation                                 | ✅      | ❌    |                                                                        |
| Block: Social Feed (comments)                     | ✅      | ❌    |                                                                        |
| Block: Report embed                               | ✅      | ❌    |                                                                        |
| Grid layout for blocks                            | ✅      | ✅    |                                                                        |
| Chart list, create, edit                          | ✅      | ✅    |                                                                        |
| Chart translation                                 | ✅      | ❌    |                                                                        |
| Chart types: bar, line, pie, funnel, gauge, radar | ✅      | ✅    | All types confirmed in lib/js/src/compose/types/chart/                 |
| Public page view                                  | ✅      | ✅    |                                                                        |
| Public record create/view via page                | ✅      | ✅    |                                                                        |
| Resource/content translator UI                    | ✅      | ❌    |                                                                        |
| Compose permissions editor                        | ✅      | ❌    |                                                                        |

---

## One (App Launcher)

| Feature                              | Corteza | Human | Notes |
| ------------------------------------ | ------- | ----- | ----- |
| Application list/launcher            | ✅      | ✅    |       |
| Unify-filtered app list              | ✅      | ✅    |       |
| Jitsi video conferencing integration | ✅      | ❌    |       |

---

## Workflow Builder

| Feature                                  | Corteza | Human | Notes                                                |
| ---------------------------------------- | ------- | ----- | ---------------------------------------------------- |
| Workflow list view                       | ✅      | 🆕    | TAQ handles this differently                         |
| Visual workflow canvas                   | ✅      | 🆕    | TAQ uses Vue Flow                                    |
| Trigger nodes                            | ✅      | 🆕    | TAQ has Trigger nodes                                |
| Function/step nodes                      | ✅      | 🆕    | TAQ has Step nodes                                   |
| Conditional/gateway nodes                | ✅      | 🆕    | TAQ has Branch nodes                                 |
| Iterator/loop nodes                      | ✅      | ❌    | Not in TAQ                                           |
| Delay nodes                              | ✅      | ❌    | Not in TAQ                                           |
| Error handler nodes                      | ✅      | ❌    | Not in TAQ                                           |
| Prompt nodes (user interaction mid-flow) | ✅      | ❌    | Not in TAQ                                           |
| Execute-other-workflow nodes             | ✅      | ❌    | Not in TAQ                                           |
| Expression editor                        | ✅      | ⚠️    | TAQ has DynamicInput; no dedicated expression editor |
| Workflow export / import                 | ✅      | ❌    |                                                      |
| Workflow validation                      | ✅      | ❌    |                                                      |
| Workflow labels                          | ✅      | ❌    |                                                      |
| Undo/redo in flow editor                 | ❌      | 🆕    | TAQ-only feature                                     |
| Manual test run from builder             | ❌      | 🆕    | TAQ-only feature                                     |
| Execution sessions tracking              | ✅      | ⚠️    | Backend has it; no frontend UI in TAQ                |
| Workflow prompts via WebSocket           | ✅      | ❌    |                                                      |

---

## Reporter

| Feature                                  | Corteza | Human | Notes                       |
| ---------------------------------------- | ------- | ----- | --------------------------- |
| Report list, create, edit, view          | ✅      | ❌    | Entire Reporter app missing |
| Compose record loader datasource         | ✅      | ❌    |                             |
| Data join/link between sources           | ✅      | ❌    |                             |
| Data aggregation (sum, avg, count, etc.) | ✅      | ❌    |                             |
| Report block: Table                      | ✅      | ❌    |                             |
| Report block: Chart                      | ✅      | ❌    |                             |
| Report block: Metric                     | ✅      | ❌    |                             |
| Report block: Text                       | ✅      | ❌    |                             |
| Report scenarios (saved filter combos)   | ✅      | ❌    |                             |
| Report sidebar / editor toolbar          | ✅      | ❌    |                             |

---

## Privacy

| Feature                            | Corteza | Human | Notes                          |
| ---------------------------------- | ------- | ----- | ------------------------------ |
| Privacy dashboard                  | ✅      | ❌    | Entire Privacy app missing     |
| Sensitive data overview            | ✅      | ❌    |                                |
| Privacy request: Export            | ✅      | ❌    |                                |
| Privacy request: Delete            | ✅      | ❌    |                                |
| Privacy request: Correct           | ✅      | ❌    |                                |
| Privacy request list/filter/status | ✅      | ❌    |                                |
| Request comments/threading         | ✅      | ❌    |                                |
| Data connection map visualization  | ✅      | ❌    |                                |
| Sensitivity level configuration    | ✅      | ❌    | Backend exists; no UI in Human |

---

## Discovery

| Feature                                                  | Corteza | Human | Notes                        |
| -------------------------------------------------------- | ------- | ----- | ---------------------------- |
| Discovery map (visual resource nav)                      | ✅      | ❌    | Entire Discovery app missing |
| Global cross-resource search                             | ✅      | ❌    |                              |
| Search results by type (namespace, module, record, user) | ✅      | ❌    |                              |
| Search filters                                           | ✅      | ❌    |                              |

---

## Shared / Cross-cutting

| Feature                                       | Corteza | Human | Notes                                               |
| --------------------------------------------- | ------- | ----- | --------------------------------------------------- |
| OAuth2 / SAML / OIDC authentication           | ✅      | ✅    | `$Auth` plugin present                              |
| Multi-factor authentication (MFA)             | ✅      | ❌    | No MFA UI                                           |
| Role-based access control (RBAC)              | ✅      | ⚠️    | `useRBAC` composable exists; no full permission UI  |
| Multi-language / i18n                         | ✅      | ✅    | Vue I18n integrated                                 |
| RTL text direction support                    | ✅      | ❌    |                                                     |
| Translation management UI                     | ✅      | ❌    |                                                     |
| Dark / light theme toggle                     | ✅      | ✅    |                                                     |
| Custom logo/favicon/branding                  | ✅      | ⚠️    | `$Settings` supports logo; no admin UI to configure |
| WebSocket / real-time notifications           | ✅      | ❌    | No WS found in Human frontend                       |
| In-app notification list (read/unread/delete) | ✅      | ❌    | Toast only                                          |
| Toast notifications (error/success/warn)      | ✅      | ✅    |                                                     |
| File upload (single + batch)                  | ✅      | ❌    | No file upload components found                     |
| File preview / download                       | ✅      | ❌    |                                                     |
| Federation (multi-instance sync)              | ✅      | ❌    |                                                     |
| Server-side automation scripts (Corredor)     | ✅      | ❌    |                                                     |
| Client-side automation scripts                | ✅      | ❌    |                                                     |
| Lazy loading / code splitting                 | ✅      | ✅    |                                                     |
| Server-side pagination                        | ✅      | ✅    | Cursor-based pagination in record store             |
| Breadcrumb navigation                         | ✅      | ⚠️    | Topbar exists; breadcrumbs unclear                  |
| Collapsible sidebar                           | ✅      | ✅    | CSidebar component                                  |
| Namespace switcher in sidebar                 | ✅      | ✅    | CSidebarNamespaceSwitcher                           |

---

## Human-only (TAQ)

| Feature                                              | Corteza | Human |
| ---------------------------------------------------- | ------- | ----- |
| TAQ automation builder dashboard                     | ❌      | 🆕    |
| Vue Flow visual node editor                          | ❌      | 🆕    |
| Undo/redo in flow editor                             | ❌      | 🆕    |
| Manual run/test from builder                         | ❌      | 🆕    |
| Dynamic segment forms (auto-generated from metadata) | ❌      | 🆕    |
| Node reference panel with docs                       | ❌      | 🆕    |
| Enable/disable automation toggle                     | ❌      | 🆕    |

---

## Summary

| Status                             | Count (approx.) |
| ---------------------------------- | --------------- |
| ✅ Implemented in both             | ~33             |
| ⚠️ Partial in Human                | ~14             |
| ❌ Corteza only (not yet in Human) | ~85             |
| 🆕 Human only                      | ~7              |

**Biggest gaps:** Full Admin panel (users/roles/settings/audit), Reporter app, Privacy app, Discovery app, 13/16 page block types, Record/User field types, advanced record features (bulk edit, import/export, revisions), and full workflow node types (iterator, delay, error handler, prompt, exec-workflow).
