# Attachments

| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
| `GET` | `/namespace/{namespaceID}/attachment/{kind}/` | List, filter all page attachments |
| `GET` | `/namespace/{namespaceID}/attachment/{kind}/{attachmentID}` | Attachment details |
| `DELETE` | `/namespace/{namespaceID}/attachment/{kind}/{attachmentID}` | Delete attachment |
| `GET` | `/namespace/{namespaceID}/attachment/{kind}/{attachmentID}/original/{name}` | Serves attached file |
| `GET` | `/namespace/{namespaceID}/attachment/{kind}/{attachmentID}/preview.{ext}` | Serves preview of an attached file |

## List, filter all page attachments

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/attachment/{kind}/` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| pageID | uint64 | GET | Filter attachments by page ID | N/A | NO |
| moduleID | uint64 | GET | Filter attachments by module ID | N/A | NO |
| recordID | uint64 | GET | Filter attachments by record ID | N/A | NO |
| fieldName | string | GET | Filter attachments by field name | N/A | NO |
| page | uint | GET | Page number | N/A | NO |
| perPage | uint | GET | Returned items per page (default 50) | N/A | NO |
| sign | string | GET | Signature | N/A | NO |
| userID | uint64 | GET | User ID | N/A | NO |
| kind | string | PATH | Attachment kind | N/A | YES |
| namespaceID | uint64 | PATH | Namespace ID | N/A | YES |

## Attachment details

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/attachment/{kind}/{attachmentID}` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| attachmentID | uint64 | PATH | Attachment ID | N/A | YES |
| kind | string | PATH | Attachment kind | N/A | YES |
| namespaceID | uint64 | PATH | Namespace ID | N/A | YES |
| sign | string | GET | Signature | N/A | NO |
| userID | uint64 | GET | User ID | N/A | NO |

## Delete attachment

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/attachment/{kind}/{attachmentID}` | HTTP/S | DELETE | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| attachmentID | uint64 | PATH | Attachment ID | N/A | YES |
| kind | string | PATH | Attachment kind | N/A | YES |
| namespaceID | uint64 | PATH | Namespace ID | N/A | YES |
| sign | string | GET | Signature | N/A | NO |
| userID | uint64 | GET | User ID | N/A | NO |

## Serves attached file

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/attachment/{kind}/{attachmentID}/original/{name}` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| download | bool | GET | Force file download | N/A | NO |
| sign | string | GET | Signature | N/A | NO |
| userID | uint64 | GET | User ID | N/A | NO |
| attachmentID | uint64 | PATH | Attachment ID | N/A | YES |
| name | string | PATH | File name | N/A | YES |
| kind | string | PATH | Attachment kind | N/A | YES |
| namespaceID | uint64 | PATH | Namespace ID | N/A | YES |

## Serves preview of an attached file

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/attachment/{kind}/{attachmentID}/preview.{ext}` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| attachmentID | uint64 | PATH | Attachment ID | N/A | YES |
| ext | string | PATH | Preview extension/format | N/A | YES |
| kind | string | PATH | Attachment kind | N/A | YES |
| namespaceID | uint64 | PATH | Namespace ID | N/A | YES |
| sign | string | GET | Signature | N/A | NO |
| userID | uint64 | GET | User ID | N/A | NO |

---




# Compose automation scripts

| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
| `GET` | `/automation/` | List all available automation scripts for compose resources |
| `POST` | `/automation/trigger` | Triggers execution of a specific script on a system service level |

## List all available automation scripts for compose resources

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/automation/` | HTTP/S | GET |
Warning: implode(): Invalid arguments passed in /private/tmp/Users/darh/Work.crust/corteza-server/codegen/templates/README.tpl on line 32
 |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| resourceTypes | []string | GET | Filter by resource type | N/A | NO |
| eventTypes | []string | GET | Filter by event type | N/A | NO |
| excludeClientScripts | bool | GET | Do not include client scripts | N/A | NO |
| excludeServerScripts | bool | GET | Do not include server scripts | N/A | NO |

## Triggers execution of a specific script on a system service level

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/automation/trigger` | HTTP/S | POST |
Warning: implode(): Invalid arguments passed in /private/tmp/Users/darh/Work.crust/corteza-server/codegen/templates/README.tpl on line 32
 |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| script | string | POST | Script to execute | N/A | YES |

---




# Charts

| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
| `GET` | `/namespace/{namespaceID}/chart/` | List/read charts |
| `POST` | `/namespace/{namespaceID}/chart/` | List/read charts  |
| `GET` | `/namespace/{namespaceID}/chart/{chartID}` | Read charts by ID |
| `POST` | `/namespace/{namespaceID}/chart/{chartID}` | Add/update charts |
| `DELETE` | `/namespace/{namespaceID}/chart/{chartID}` | Delete chart |

## List/read charts

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/chart/` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| query | string | GET | Search query to match against charts | N/A | NO |
| handle | string | GET | Search charts by handle | N/A | NO |
| page | uint | GET | Page number | N/A | NO |
| perPage | uint | GET | Returned items per page (default 50) | N/A | NO |
| sort | string | GET | Sort charts | N/A | NO |
| namespaceID | uint64 | PATH | Namespace ID | N/A | YES |

## List/read charts

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/chart/` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| config | sqlxTypes.JSONText | POST | Chart JSON | N/A | YES |
| name | string | POST | Chart name | N/A | YES |
| handle | string | POST | Chart handle | N/A | NO |
| namespaceID | uint64 | PATH | Namespace ID | N/A | YES |

## Read charts by ID

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/chart/{chartID}` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| chartID | uint64 | PATH | Chart ID | N/A | YES |
| namespaceID | uint64 | PATH | Namespace ID | N/A | YES |

## Add/update charts

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/chart/{chartID}` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| chartID | uint64 | PATH | Chart ID | N/A | YES |
| namespaceID | uint64 | PATH | Namespace ID | N/A | YES |
| config | sqlxTypes.JSONText | POST | Chart JSON | N/A | YES |
| name | string | POST | Chart name | N/A | YES |
| handle | string | POST | Chart handle | N/A | NO |
| updatedAt | *time.Time | POST | Last update (or creation) date | N/A | NO |

## Delete chart

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/chart/{chartID}` | HTTP/S | DELETE |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| chartID | uint64 | PATH | Chart ID | N/A | YES |
| namespaceID | uint64 | PATH | Namespace ID | N/A | YES |

---




# Modules

Compose module definitions

| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
| `GET` | `/namespace/{namespaceID}/module/` | List modules |
| `POST` | `/namespace/{namespaceID}/module/` | Create module |
| `GET` | `/namespace/{namespaceID}/module/{moduleID}` | Read module |
| `POST` | `/namespace/{namespaceID}/module/{moduleID}` | Update module |
| `DELETE` | `/namespace/{namespaceID}/module/{moduleID}` | Delete module |
| `POST` | `/namespace/{namespaceID}/module/{moduleID}/trigger` | Fire compose:module trigger |

## List modules

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/module/` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| query | string | GET | Search query | N/A | NO |
| name | string | GET | Search by name | N/A | NO |
| handle | string | GET | Search by handle | N/A | NO |
| page | uint | GET | Page number | N/A | NO |
| perPage | uint | GET | Returned items per page (default 50) | N/A | NO |
| sort | string | GET | Sort | N/A | NO |
| namespaceID | uint64 | PATH | Namespace ID | N/A | YES |

## Create module

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/module/` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| name | string | POST | Module Name | N/A | YES |
| handle | string | POST | Module handle | N/A | NO |
| fields | types.ModuleFieldSet | POST | Fields JSON | N/A | YES |
| meta | sqlxTypes.JSONText | POST | Module meta data | N/A | YES |
| namespaceID | uint64 | PATH | Namespace ID | N/A | YES |

## Read module

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/module/{moduleID}` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| moduleID | uint64 | PATH | Module ID | N/A | YES |
| namespaceID | uint64 | PATH | Namespace ID | N/A | YES |

## Update module

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/module/{moduleID}` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| moduleID | uint64 | PATH | Module ID | N/A | YES |
| namespaceID | uint64 | PATH | Namespace ID | N/A | YES |
| name | string | POST | Module Name | N/A | YES |
| handle | string | POST | Module Handle | N/A | NO |
| fields | types.ModuleFieldSet | POST | Fields JSON | N/A | YES |
| meta | sqlxTypes.JSONText | POST | Module meta data | N/A | YES |
| updatedAt | *time.Time | POST | Last update (or creation) date | N/A | NO |

## Delete module

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/module/{moduleID}` | HTTP/S | DELETE |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| moduleID | uint64 | PATH | Module ID | N/A | YES |
| namespaceID | uint64 | PATH | Namespace ID | N/A | YES |

## Fire compose:module trigger

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/module/{moduleID}/trigger` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| moduleID | uint64 | PATH | ID | N/A | YES |
| namespaceID | uint64 | PATH | Namespace ID | N/A | YES |
| script | string | POST | Script to execute | N/A | YES |

---




# Namespaces

| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
| `GET` | `/namespace/` | List namespaces |
| `POST` | `/namespace/` | Create namespace |
| `GET` | `/namespace/{namespaceID}` | Read namespace |
| `POST` | `/namespace/{namespaceID}` | Update namespace |
| `DELETE` | `/namespace/{namespaceID}` | Delete namespace |
| `POST` | `/namespace/{namespaceID}/trigger` | Fire compose:namespace trigger |

## List namespaces

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| query | string | GET | Search query | N/A | NO |
| slug | string | GET | Search by namespace slug | N/A | NO |
| page | uint | GET | Page number | N/A | NO |
| perPage | uint | GET | Returned items per page (default 50) | N/A | NO |
| sort | string | GET | Sort namespaces | N/A | NO |

## Create namespace

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| name | string | POST | Name | N/A | YES |
| slug | string | POST | Slug (url path part) | N/A | NO |
| enabled | bool | POST | Enabled | N/A | NO |
| meta | sqlxTypes.JSONText | POST | Meta data | N/A | YES |

## Read namespace

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| namespaceID | uint64 | PATH | ID | N/A | YES |

## Update namespace

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| namespaceID | uint64 | PATH | ID | N/A | YES |
| name | string | POST | Name | N/A | YES |
| slug | string | POST | Slug (url path part) | N/A | NO |
| enabled | bool | POST | Enabled | N/A | NO |
| meta | sqlxTypes.JSONText | POST | Meta data | N/A | YES |
| updatedAt | *time.Time | POST | Last update (or creation) date | N/A | NO |

## Delete namespace

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}` | HTTP/S | DELETE |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| namespaceID | uint64 | PATH | ID | N/A | YES |

## Fire compose:namespace trigger

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/trigger` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| namespaceID | uint64 | PATH | ID | N/A | YES |
| script | string | POST | Script to execute | N/A | YES |

---




# Notifications

Compose Notifications

| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
| `POST` | `/notification/email` | Send email from the Compose |

## Send email from the Compose

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/notification/email` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| to | []string | POST | Email addresses | N/A | YES |
| cc | []string | POST | Email addresses | N/A | NO |
| replyTo | string | POST | Email address in reply-to field | N/A | NO |
| subject  | string | POST | Email subject | N/A | NO |
| content | sqlxTypes.JSONText | POST | Message content | N/A | YES |
| remoteAttachments | []string | POST | Remote files to attach to the email | N/A | NO |

---




# Pages

Compose pages

| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
| `GET` | `/namespace/{namespaceID}/page/` | List available pages |
| `POST` | `/namespace/{namespaceID}/page/` | Create page |
| `GET` | `/namespace/{namespaceID}/page/{pageID}` | Get page details |
| `GET` | `/namespace/{namespaceID}/page/tree` | Get page all (non-record) pages, hierarchically |
| `POST` | `/namespace/{namespaceID}/page/{pageID}` | Update page |
| `POST` | `/namespace/{namespaceID}/page/{selfID}/reorder` | Reorder pages |
| `Delete` | `/namespace/{namespaceID}/page/{pageID}` | Delete page |
| `POST` | `/namespace/{namespaceID}/page/{pageID}/attachment` | Uploads attachment to page |
| `POST` | `/namespace/{namespaceID}/page/{pageID}/trigger` | Fire compose:page trigger |

## List available pages

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/page/` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| selfID | uint64 | GET | Parent page ID | N/A | NO |
| query | string | GET | Search query | N/A | NO |
| handle | string | GET | Search by handle | N/A | NO |
| page | uint | GET | Page number | N/A | NO |
| perPage | uint | GET | Returned items per page (default 50) | N/A | NO |
| sort | string | GET | Sort | N/A | NO |
| namespaceID | uint64 | PATH | Namespace ID | N/A | YES |

## Create page

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/page/` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| selfID | uint64 | POST | Parent Page ID | N/A | NO |
| moduleID | uint64 | POST | Module ID | N/A | NO |
| title | string | POST | Title | N/A | YES |
| handle | string | POST | Handle | N/A | NO |
| description | string | POST | Description | N/A | NO |
| visible | bool | POST | Visible in navigation | N/A | NO |
| blocks | sqlxTypes.JSONText | POST | Blocks JSON | N/A | YES |
| namespaceID | uint64 | PATH | Namespace ID | N/A | YES |

## Get page details

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/page/{pageID}` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| pageID | uint64 | PATH | Page ID | N/A | YES |
| namespaceID | uint64 | PATH | Namespace ID | N/A | YES |

## Get page all (non-record) pages, hierarchically

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/page/tree` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| namespaceID | uint64 | PATH | Namespace ID | N/A | YES |

## Update page

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/page/{pageID}` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| pageID | uint64 | PATH | Page ID | N/A | YES |
| namespaceID | uint64 | PATH | Namespace ID | N/A | YES |
| selfID | uint64 | POST | Parent Page ID | N/A | NO |
| moduleID | uint64 | POST | Module ID (optional) | N/A | NO |
| title | string | POST | Title | N/A | YES |
| handle | string | POST | Handle | N/A | NO |
| description | string | POST | Description | N/A | NO |
| visible | bool | POST | Visible in navigation | N/A | NO |
| blocks | sqlxTypes.JSONText | POST | Blocks JSON | N/A | YES |

## Reorder pages

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/page/{selfID}/reorder` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| selfID | uint64 | PATH | Parent page ID | N/A | YES |
| namespaceID | uint64 | PATH | Namespace ID | N/A | YES |
| pageIDs | []string | POST | Page ID order | N/A | YES |

## Delete page

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/page/{pageID}` | HTTP/S | Delete |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| pageID | uint64 | PATH | Page ID | N/A | YES |
| namespaceID | uint64 | PATH | Namespace ID | N/A | YES |

## Uploads attachment to page

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/page/{pageID}/attachment` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| pageID | uint64 | PATH | Page ID | N/A | YES |
| namespaceID | uint64 | PATH | Namespace ID | N/A | YES |
| upload | *multipart.FileHeader | POST | File to upload | N/A | YES |

## Fire compose:page trigger

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/page/{pageID}/trigger` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| pageID | uint64 | PATH | Page ID | N/A | YES |
| namespaceID | uint64 | PATH | Namespace ID | N/A | YES |
| script | string | POST | Script to execute | N/A | YES |

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




# Records

Compose records

| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
| `GET` | `/namespace/{namespaceID}/module/{moduleID}/record/report` | Generates report from module records |
| `GET` | `/namespace/{namespaceID}/module/{moduleID}/record/` | List/read records from module section |
| `POST` | `/namespace/{namespaceID}/module/{moduleID}/record/import` | Initiate record import session |
| `PATCH` | `/namespace/{namespaceID}/module/{moduleID}/record/import/{sessionID}` | Run record import |
| `GET` | `/namespace/{namespaceID}/module/{moduleID}/record/import/{sessionID}` | Get import progress |
| `GET` | `/namespace/{namespaceID}/module/{moduleID}/record/export{filename}.{ext}` | Exports records that match  |
| `POST` | `/namespace/{namespaceID}/module/{moduleID}/record/exec/{procedure}` | Executes server-side procedure over one or more module records |
| `POST` | `/namespace/{namespaceID}/module/{moduleID}/record/` | Create record in module section |
| `GET` | `/namespace/{namespaceID}/module/{moduleID}/record/{recordID}` | Read records by ID from module section |
| `POST` | `/namespace/{namespaceID}/module/{moduleID}/record/{recordID}` | Update records in module section |
| `DELETE` | `/namespace/{namespaceID}/module/{moduleID}/record/{recordID}` | Delete record row from module section |
| `POST` | `/namespace/{namespaceID}/module/{moduleID}/record/attachment` | Uploads attachment and validates it against record field requirements |
| `POST` | `/namespace/{namespaceID}/module/{moduleID}/record/{recordID}/trigger` | Fire compose:record trigger |

## Generates report from module records

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/module/{moduleID}/record/report` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| metrics | string | GET | Metrics (eg: 'SUM(money), MAX(calls)') | N/A | NO |
| dimensions | string | GET | Dimensions (eg: 'DATE(foo), status') | N/A | YES |
| filter | string | GET | Filter (eg: 'DATE(foo) > 2010') | N/A | NO |
| namespaceID | uint64 | PATH | Namespace ID | N/A | YES |
| moduleID | uint64 | PATH | Module ID | N/A | YES |

## List/read records from module section

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/module/{moduleID}/record/` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| filter | string | GET | Filtering condition | N/A | NO |
| page | uint | GET | Page number | N/A | NO |
| perPage | uint | GET | Returned items per page (default 50) | N/A | NO |
| sort | string | GET | Sort field (default id desc) | N/A | NO |
| namespaceID | uint64 | PATH | Namespace ID | N/A | YES |
| moduleID | uint64 | PATH | Module ID | N/A | YES |

## Initiate record import session

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/module/{moduleID}/record/import` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| upload | *multipart.FileHeader | POST | File import | N/A | YES |
| namespaceID | uint64 | PATH | Namespace ID | N/A | YES |
| moduleID | uint64 | PATH | Module ID | N/A | YES |

## Run record import

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/module/{moduleID}/record/import/{sessionID}` | HTTP/S | PATCH |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| sessionID | uint64 | PATH | Import session | N/A | YES |
| namespaceID | uint64 | PATH | Namespace ID | N/A | YES |
| moduleID | uint64 | PATH | Module ID | N/A | YES |
| fields | json.RawMessage | POST | Fields defined by import file | N/A | YES |
| onError | string | POST | What happens if record fails to import | N/A | YES |

## Get import progress

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/module/{moduleID}/record/import/{sessionID}` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| sessionID | uint64 | PATH | Import session | N/A | YES |
| namespaceID | uint64 | PATH | Namespace ID | N/A | YES |
| moduleID | uint64 | PATH | Module ID | N/A | YES |

## Exports records that match

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/module/{moduleID}/record/export{filename}.{ext}` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| filter | string | GET | Filtering condition | N/A | NO |
| fields | []string | GET | Fields to export | N/A | YES |
| filename | string | PATH | Filename to use | N/A | NO |
| ext | string | PATH | Export format | N/A | YES |
| namespaceID | uint64 | PATH | Namespace ID | N/A | YES |
| moduleID | uint64 | PATH | Module ID | N/A | YES |

## Executes server-side procedure over one or more module records

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/module/{moduleID}/record/exec/{procedure}` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| procedure | string | PATH | Name of procedure to execute | N/A | YES |
| namespaceID | uint64 | PATH | Namespace ID | N/A | YES |
| moduleID | uint64 | PATH | Module ID | N/A | YES |
| args | []ProcedureArg | POST | Procedure arguments | N/A | NO |

## Create record in module section

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/module/{moduleID}/record/` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| values | types.RecordValueSet | POST | Record values | N/A | YES |
| namespaceID | uint64 | PATH | Namespace ID | N/A | YES |
| moduleID | uint64 | PATH | Module ID | N/A | YES |

## Read records by ID from module section

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/module/{moduleID}/record/{recordID}` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| recordID | uint64 | PATH | Record ID | N/A | YES |
| namespaceID | uint64 | PATH | Namespace ID | N/A | YES |
| moduleID | uint64 | PATH | Module ID | N/A | YES |

## Update records in module section

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/module/{moduleID}/record/{recordID}` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| recordID | uint64 | PATH | Record ID | N/A | YES |
| namespaceID | uint64 | PATH | Namespace ID | N/A | YES |
| moduleID | uint64 | PATH | Module ID | N/A | YES |
| values | types.RecordValueSet | POST | Record values | N/A | YES |

## Delete record row from module section

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/module/{moduleID}/record/{recordID}` | HTTP/S | DELETE |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| recordID | uint64 | PATH | Record ID | N/A | YES |
| namespaceID | uint64 | PATH | Namespace ID | N/A | YES |
| moduleID | uint64 | PATH | Module ID | N/A | YES |

## Uploads attachment and validates it against record field requirements

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/module/{moduleID}/record/attachment` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| recordID | uint64 | POST | Record ID | N/A | NO |
| fieldName | string | POST | Field name | N/A | YES |
| upload | *multipart.FileHeader | POST | File to upload | N/A | YES |
| namespaceID | uint64 | PATH | Namespace ID | N/A | YES |
| moduleID | uint64 | PATH | Module ID | N/A | YES |

## Fire compose:record trigger

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/module/{moduleID}/record/{recordID}/trigger` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| recordID | uint64 | PATH | ID | N/A | YES |
| namespaceID | uint64 | PATH | Namespace ID | N/A | YES |
| moduleID | uint64 | PATH | Module ID | N/A | YES |
| script | string | POST | Script to execute | N/A | YES |

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