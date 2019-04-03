# Applications

| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
| `GET` | `/application/` | List applications |
| `POST` | `/application/` | Create application |
| `PUT` | `/application/{applicationID}` | Update user details |
| `GET` | `/application/{applicationID}` | Read application details |
| `DELETE` | `/application/{applicationID}` | Remove application |

## List applications

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/application/` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |

## Create application

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/application/` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| name | string | POST | Email | N/A | YES |
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

---




# Authentication

| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
| `GET` | `/auth/check` | Check JWT token |
| `POST` | `/auth/login` | Login user |
| `GET` | `/auth/logout` | Delete JWT token (Sign Out) |

## Check JWT token

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/auth/check` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |

## Login user

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/auth/login` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| username | string | POST | Username | N/A | YES |
| password | string | POST | Password | N/A | YES |

## Delete JWT token (Sign Out)

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/auth/logout` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |

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
| permissions | []rules.Rule | POST | List of permissions to set | N/A | YES |

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

## Update role details

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/roles/` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| name | string | POST | Name of Role | N/A | YES |
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




# Users

| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
| `GET` | `/users/` | Search users (Directory) |
| `POST` | `/users/` | Create user |
| `PUT` | `/users/{userID}` | Update user details |
| `GET` | `/users/{userID}` | Read user details and memberships |
| `DELETE` | `/users/{userID}` | Remove user |
| `POST` | `/users/{userID}/suspend` | Suspend user |
| `POST` | `/users/{userID}/unsuspend` | Unsuspend user |

## Search users (Directory)

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/users/` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| query | string | GET | Search query to match against users | N/A | NO |
| username | string | GET | Search username to match against users | N/A | NO |
| email | string | GET | Search email to match against users | N/A | NO |

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
| kind | string | POST | Kind (normal, bot) | N/A | NO |

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
| kind | string | POST | Kind (normal, bot) | N/A | NO |

## Read user details and memberships

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

---