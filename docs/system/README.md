# Applications

| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
| `GET` | `/application/` | List applications |
| `POST` | `/application/` | Create application |
| `PUT` | `/application/{applicationID}` | Update user details |
| `GET` | `/application/{applicationID}` | Read application details |
| `DELETE` | `/application/{applicationID}` | Remove application |
| `POST` | `/application/{applicationID}/undelete` | Undelete application |

## List applications

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/application/` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| name | string | GET | Application name | N/A | NO |
| query | string | GET | Filter applications | N/A | NO |
| deleted | uint | GET | Exclude (0, default), include (1) or return only (2) deleted roles | N/A | NO |
| page | uint | GET | Page number | N/A | NO |
| perPage | uint | GET | Returned items per page (default 50) | N/A | NO |
| sort | string | GET | Sort | N/A | NO |

## Create application

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/application/` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| name | string | POST | Application name | N/A | YES |
| enabled | bool | POST | Enabled | N/A | NO |
| unify | sqlxTypes.JSONText | POST | Unify properties | N/A | NO |
| config | sqlxTypes.JSONText | POST | Arbitrary JSON holding application configuration | N/A | NO |

## Update user details

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/application/{applicationID}` | HTTP/S | PUT |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| applicationID | uint64 | PATH | Application ID | N/A | YES |
| name | string | POST | Email | N/A | YES |
| enabled | bool | POST | Enabled | N/A | NO |
| unify | sqlxTypes.JSONText | POST | Unify properties | N/A | NO |
| config | sqlxTypes.JSONText | POST | Arbitrary JSON holding application configuration | N/A | NO |

## Read application details

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/application/{applicationID}` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| applicationID | uint64 | PATH | Application ID | N/A | YES |

## Remove application

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/application/{applicationID}` | HTTP/S | DELETE |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| applicationID | uint64 | PATH | Application ID | N/A | YES |

## Undelete application

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/application/{applicationID}/undelete` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| applicationID | uint64 | PATH | Application ID | N/A | YES |

---




# Authentication

| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
| `GET` | `/auth/` | Returns auth settings |
| `GET` | `/auth/check` | Check JWT token |
| `POST` | `/auth/exchange` | Exchange auth token for JWT |
| `GET` | `/auth/logout` | Logout |

## Returns auth settings

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/auth/` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |

## Check JWT token

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/auth/check` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |

## Exchange auth token for JWT

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/auth/exchange` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| token | string | POST | Token to be exchanged for JWT | N/A | YES |

## Logout

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/auth/logout` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |

---




# Internal authentication

| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
| `POST` | `/auth/internal/login` | Login user |
| `POST` | `/auth/internal/signup` | User signup/registration |
| `POST` | `/auth/internal/request-password-reset` | Request password reset token (via email) |
| `POST` | `/auth/internal/exchange-password-reset-token` | Exchange password reset token for new token and user info |
| `POST` | `/auth/internal/reset-password` | Reset password with exchanged password reset token |
| `POST` | `/auth/internal/confirm-email` | Confirm email with token |
| `POST` | `/auth/internal/change-password` | Changes password for current user, requires current password |

## Login user

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/auth/internal/login` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| email | string | POST | Email | N/A | YES |
| password | string | POST | Password | N/A | YES |

## User signup/registration

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/auth/internal/signup` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| email | string | POST | Email | N/A | YES |
| username | string | POST | Username | N/A | NO |
| password | string | POST | Password | N/A | YES |
| handle | string | POST | User handle | N/A | NO |
| name | string | POST | Display name | N/A | NO |

## Request password reset token (via email)

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/auth/internal/request-password-reset` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| email | string | POST | Email | N/A | YES |

## Exchange password reset token for new token and user info

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/auth/internal/exchange-password-reset-token` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| token | string | POST | Token | N/A | YES |

## Reset password with exchanged password reset token

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/auth/internal/reset-password` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| token | string | POST | Token | N/A | YES |
| password | string | POST | Password | N/A | YES |

## Confirm email with token

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/auth/internal/confirm-email` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| token | string | POST | Token | N/A | YES |

## Changes password for current user, requires current password

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/auth/internal/change-password` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| oldPassword | string | POST | Old password | N/A | YES |
| newPassword | string | POST | New password | N/A | YES |

---




# Organisations

Organisations represent a top-level grouping entity. There may be many organisations defined in a single deployment.

| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
| `GET` | `/organisations/` | List organisations |
| `POST` | `/organisations/` | Create organisation |
| `PUT` | `/organisations/{id}` | Update organisation details |
| `DELETE` | `/organisations/{id}` | Remove organisation |
| `GET` | `/organisations/{id}` | Read organisation details |
| `POST` | `/organisations/{id}/archive` | Archive organisation |

## List organisations

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/organisations/` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| query | string | GET | Search query | N/A | NO |

## Create organisation

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/organisations/` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| name | string | POST | Organisation Name | N/A | YES |

## Update organisation details

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/organisations/{id}` | HTTP/S | PUT | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | PATH | Organisation ID | N/A | NO |
| name | string | POST | Organisation Name | N/A | YES |

## Remove organisation

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/organisations/{id}` | HTTP/S | DELETE | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | PATH | Organisation ID | N/A | YES |

## Read organisation details

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/organisations/{id}` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | GET | Organisation ID | N/A | YES |

## Archive organisation

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/organisations/{id}/archive` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | PATH | Organisation ID | N/A | YES |

---




# Permissions

| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
| `GET` | `/permissions/` | Retrieve defined permissions |
| `GET` | `/permissions/effective` | Effective rules for current user |
| `GET` | `/permissions/{roleID}/rules` | Retrieve role permissions |
| `DELETE` | `/permissions/{roleID}/rules` | Remove all defined role permissions |
| `PATCH` | `/permissions/{roleID}/rules` | Update permission settings |

## Retrieve defined permissions

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/permissions/` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |

## Effective rules for current user

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/permissions/effective` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| resource | string | GET | Show only rules for a specific resource | N/A | NO |

## Retrieve role permissions

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/permissions/{roleID}/rules` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| roleID | uint64 | PATH | Role ID | N/A | YES |

## Remove all defined role permissions

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/permissions/{roleID}/rules` | HTTP/S | DELETE | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| roleID | uint64 | PATH | Role ID | N/A | YES |

## Update permission settings

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/permissions/{roleID}/rules` | HTTP/S | PATCH | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| roleID | uint64 | PATH | Role ID | N/A | YES |
| rules | permissions.RuleSet | POST | List of permission rules to set | N/A | YES |

---




# Reminders

| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
| `GET` | `/reminder/` | List/read reminders |
| `POST` | `/reminder/` | Add new reminder |
| `PUT` | `/reminder/{reminderID}` | Update reminder |
| `GET` | `/reminder/{reminderID}` | Read reminder by ID |
| `DELETE` | `/reminder/{reminderID}` | Delete reminder |
| `PATCH` | `/reminder/{reminderID}/dismiss` | Dismiss reminder |
| `PATCH` | `/reminder/{reminderID}/snooze` | Snooze reminder |

## List/read reminders

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/reminder/` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| reminderID | []string | GET | Filter by reminder ID | N/A | NO |
| resource | string | GET | Only reminders of a specific resource | N/A | NO |
| assignedTo | uint64 | GET | Only reminders for a given user | N/A | NO |
| scheduledFrom | *time.Time | GET | Only reminders from this time (included) | N/A | NO |
| scheduledUntil | *time.Time | GET | Only reminders up to this time (included) | N/A | NO |
| scheduledOnly | bool | GET | Only scheduled reminders | N/A | NO |
| excludeDismissed | bool | GET | Filter out dismissed reminders | N/A | NO |
| page | uint | GET | Page number | N/A | NO |
| perPage | uint | GET | Returned items per page (default 50) | N/A | NO |
| sort | string | GET | Sort | N/A | NO |

## Add new reminder

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/reminder/` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| resource | string | POST | Resource | N/A | YES |
| assignedTo | uint64 | POST | Assigned To | N/A | YES |
| payload | sqlxTypes.JSONText | POST | Payload | N/A | YES |
| remindAt | *time.Time | POST | Remind At | N/A | NO |

## Update reminder

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/reminder/{reminderID}` | HTTP/S | PUT | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| reminderID | uint64 | PATH | Reminder ID | N/A | YES |
| resource | string | POST | Resource | N/A | YES |
| assignedTo | uint64 | POST | Assigned To | N/A | YES |
| payload | sqlxTypes.JSONText | POST | Payload | N/A | YES |
| remindAt | *time.Time | POST | Remind At | N/A | NO |

## Read reminder by ID

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/reminder/{reminderID}` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| reminderID | uint64 | PATH | Reminder ID | N/A | YES |

## Delete reminder

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/reminder/{reminderID}` | HTTP/S | DELETE | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| reminderID | uint64 | PATH | Reminder ID | N/A | YES |

## Dismiss reminder

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/reminder/{reminderID}/dismiss` | HTTP/S | PATCH | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| reminderID | uint64 | PATH | reminder ID | N/A | YES |

## Snooze reminder

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/reminder/{reminderID}/snooze` | HTTP/S | PATCH | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| reminderID | uint64 | PATH | reminder ID | N/A | YES |
| remindAt | *time.Time | POST | New Remind At Time | N/A | YES |

---




# Roles

An organisation may have many roles. Roles may have many channels available. Access to channels may be shared between roles.

| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
| `GET` | `/roles/` | List roles |
| `POST` | `/roles/` | Update role details |
| `PUT` | `/roles/{roleID}` | Update role details |
| `GET` | `/roles/{roleID}` | Read role details and memberships |
| `DELETE` | `/roles/{roleID}` | Remove role |
| `POST` | `/roles/{roleID}/archive` | Archive role |
| `POST` | `/roles/{roleID}/unarchive` | Unarchive role |
| `POST` | `/roles/{roleID}/undelete` | Undelete role |
| `POST` | `/roles/{roleID}/move` | Move role to different organisation |
| `POST` | `/roles/{roleID}/merge` | Merge one role into another |
| `GET` | `/roles/{roleID}/members` | Returns all role members |
| `POST` | `/roles/{roleID}/member/{userID}` | Add member to a role |
| `DELETE` | `/roles/{roleID}/member/{userID}` | Remove member from a role |

## List roles

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/roles/` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| query | string | GET | Search query | N/A | NO |
| deleted | uint | GET | Exclude (0, default), include (1) or return only (2) deleted roles | N/A | NO |
| archived | uint | GET | Exclude (0, default), include (1) or return only (2) achived roles | N/A | NO |
| page | uint | GET | Page number | N/A | NO |
| perPage | uint | GET | Returned items per page (default 50) | N/A | NO |
| sort | string | GET | Sort | N/A | NO |

## Update role details

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/roles/` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| name | string | POST | Name of Role | N/A | YES |
| handle | string | POST | Handle for Role | N/A | YES |
| members | []string | POST | Role member IDs | N/A | NO |

## Update role details

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/roles/{roleID}` | HTTP/S | PUT | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| roleID | uint64 | PATH | Role ID | N/A | YES |
| name | string | POST | Name of Role | N/A | NO |
| handle | string | POST | Handle for Role | N/A | NO |
| members | []string | POST | Role member IDs | N/A | NO |

## Read role details and memberships

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/roles/{roleID}` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| roleID | uint64 | PATH | Role ID | N/A | YES |

## Remove role

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/roles/{roleID}` | HTTP/S | DELETE | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| roleID | uint64 | PATH | Role ID | N/A | YES |

## Archive role

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/roles/{roleID}/archive` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| roleID | uint64 | PATH | Role ID | N/A | YES |

## Unarchive role

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/roles/{roleID}/unarchive` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| roleID | uint64 | PATH | Role ID | N/A | YES |

## Undelete role

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/roles/{roleID}/undelete` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| roleID | uint64 | PATH | Role ID | N/A | YES |

## Move role to different organisation

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/roles/{roleID}/move` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| roleID | uint64 | PATH | Role ID | N/A | YES |
| organisationID | uint64 | POST | Role ID | N/A | YES |

## Merge one role into another

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/roles/{roleID}/merge` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| roleID | uint64 | PATH | Source Role ID | N/A | YES |
| destination | uint64 | POST | Destination Role ID | N/A | YES |

## Returns all role members

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/roles/{roleID}/members` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| roleID | uint64 | PATH | Source Role ID | N/A | YES |

## Add member to a role

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/roles/{roleID}/member/{userID}` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| roleID | uint64 | PATH | Source Role ID | N/A | YES |
| userID | uint64 | PATH | User ID | N/A | YES |

## Remove member from a role

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/roles/{roleID}/member/{userID}` | HTTP/S | DELETE | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| roleID | uint64 | PATH | Source Role ID | N/A | YES |
| userID | uint64 | PATH | User ID | N/A | YES |

---




# Settings

| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
| `GET` | `/settings/` | List settings |
| `PATCH` | `/settings/` | Update settings |
| `GET` | `/settings/{key}` | Get a value for a key |
| `GET` | `/settings/current` | Current compose settings |

## List settings

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/settings/` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| prefix | string | GET | Key prefix | N/A | NO |

## Update settings

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/settings/` | HTTP/S | PATCH |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| values | settings.ValueSet | POST | Array of new settings: `[{ name: ..., value: ... }]`. Omit value to remove setting | N/A | YES |

## Get a value for a key

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/settings/{key}` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| ownerID | uint64 | GET | Owner ID | N/A | NO |
| key | string | PATH | Setting key | N/A | YES |

## Current compose settings

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/settings/current` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |

---




# Statistics

| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
| `GET` | `/stats/` | List system statistics |

## List system statistics

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/stats/` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |

---




# Subscription

| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
| `GET` | `/subscription/` | Returns current subscription status |

## Returns current subscription status

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/subscription/` | HTTP/S | GET |
Warning: implode(): Invalid arguments passed in /private/tmp/Users/darh/Work.crust/corteza-server/codegen/templates/README.tpl on line 32
 |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |

---




# Users

| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
| `GET` | `/users/` | Search users (Directory) |
| `POST` | `/users/` | Create user |
| `PUT` | `/users/{userID}` | Update user details |
| `GET` | `/users/{userID}` | Read user details |
| `DELETE` | `/users/{userID}` | Remove user |
| `POST` | `/users/{userID}/suspend` | Suspend user |
| `POST` | `/users/{userID}/unsuspend` | Unsuspend user |
| `POST` | `/users/{userID}/undelete` | Undelete user |
| `POST` | `/users/{userID}/password` | Set's or changes user's password |
| `GET` | `/users/{userID}/membership` | Add member to a role |
| `POST` | `/users/{userID}/membership/{roleID}` | Add role to a user |
| `DELETE` | `/users/{userID}/membership/{roleID}` | Remove role from a user |

## Search users (Directory)

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/users/` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| userID | []string | GET | Filter by user ID | N/A | NO |
| roleID | []string | GET | Filter by role membership | N/A | NO |
| query | string | GET | Search query to match against users | N/A | NO |
| username | string | GET | Search username to match against users | N/A | NO |
| email | string | GET | Search email to match against users | N/A | NO |
| handle | string | GET | Search handle to match against users | N/A | NO |
| kind | types.UserKind | GET | Kind (normal, bot) | N/A | NO |
| incDeleted | bool | GET | [Deprecated] Include deleted users (requires 'access' permission) | N/A | NO |
| incSuspended | bool | GET | [Deprecated] Include suspended users | N/A | NO |
| deleted | uint | GET | Exclude (0, default), include (1) or return only (2) deleted users | N/A | NO |
| suspended | uint | GET | Exclude (0, default), include (1) or return only (2) suspended users | N/A | NO |
| sort | string | GET | Sort by (createdAt, updatedAt, deletedAt, suspendedAt, email, username, userID) | N/A | NO |
| page | uint | GET | Page number | N/A | NO |
| perPage | uint | GET | Returned items per page | N/A | NO |

## Create user

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/users/` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| email | string | POST | Email | N/A | YES |
| name | string | POST | Name | N/A | NO |
| handle | string | POST | Handle | N/A | NO |
| kind | types.UserKind | POST | Kind (normal, bot) | N/A | NO |

## Update user details

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/users/{userID}` | HTTP/S | PUT | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| userID | uint64 | PATH | User ID | N/A | YES |
| email | string | POST | Email | N/A | YES |
| name | string | POST | Name | N/A | YES |
| handle | string | POST | Handle | N/A | NO |
| kind | types.UserKind | POST | Kind (normal, bot) | N/A | NO |

## Read user details

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/users/{userID}` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| userID | uint64 | PATH | User ID | N/A | YES |

## Remove user

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/users/{userID}` | HTTP/S | DELETE | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| userID | uint64 | PATH | User ID | N/A | YES |

## Suspend user

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/users/{userID}/suspend` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| userID | uint64 | PATH | User ID | N/A | YES |

## Unsuspend user

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/users/{userID}/unsuspend` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| userID | uint64 | PATH | User ID | N/A | YES |

## Undelete user

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/users/{userID}/undelete` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| userID | uint64 | PATH | User ID | N/A | YES |

## Set's or changes user's password

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/users/{userID}/password` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| userID | uint64 | PATH | User ID | N/A | YES |
| password | string | POST | New password | N/A | YES |

## Add member to a role

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/users/{userID}/membership` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| userID | uint64 | PATH | User ID | N/A | YES |

## Add role to a user

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/users/{userID}/membership/{roleID}` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| roleID | uint64 | PATH | Role ID | N/A | YES |
| userID | uint64 | PATH | User ID | N/A | YES |

## Remove role from a user

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/users/{userID}/membership/{roleID}` | HTTP/S | DELETE | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| roleID | uint64 | PATH | Role ID | N/A | YES |
| userID | uint64 | PATH | User ID | N/A | YES |

---