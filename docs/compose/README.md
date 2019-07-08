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
| page | uint | GET | Page number (0 based) | N/A | NO |
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
| page | uint | GET | Page number (0 based) | N/A | NO |
| perPage | uint | GET | Returned items per page (default 50) | N/A | NO |
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

## List modules

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/module/` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| query | string | GET | Search query | N/A | NO |
| page | uint | GET | Page number (0 based) | N/A | NO |
| perPage | uint | GET | Returned items per page (default 50) | N/A | NO |
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

---




# Namespaces

| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
| `GET` | `/namespace/` | List namespaces |
| `POST` | `/namespace/` | Create namespace |
| `GET` | `/namespace/{namespaceID}` | Read namespace |
| `POST` | `/namespace/{namespaceID}` | Update namespace |
| `DELETE` | `/namespace/{namespaceID}` | Delete namespace |

## List namespaces

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| query | string | GET | Search query | N/A | NO |
| page | uint | GET | Page number (0 based) | N/A | NO |
| perPage | uint | GET | Returned items per page (default 50) | N/A | NO |

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
| page | uint | GET | Page number (0 based) | N/A | NO |
| perPage | uint | GET | Returned items per page (default 50) | N/A | NO |
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
| `POST` | `/namespace/{namespaceID}/module/{moduleID}/record/` | Create record in module section |
| `GET` | `/namespace/{namespaceID}/module/{moduleID}/record/{recordID}` | Read records by ID from module section |
| `POST` | `/namespace/{namespaceID}/module/{moduleID}/record/{recordID}` | Update records in module section |
| `DELETE` | `/namespace/{namespaceID}/module/{moduleID}/record/{recordID}` | Delete record row from module section |
| `POST` | `/namespace/{namespaceID}/module/{moduleID}/record/attachment` | Uploads attachment and validates it against record field requirements |

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
| page | uint | GET | Page number (0 based) | N/A | NO |
| perPage | uint | GET | Returned items per page (default 50) | N/A | NO |
| sort | string | GET | Sort field (default id desc) | N/A | NO |
| namespaceID | uint64 | PATH | Namespace ID | N/A | YES |
| moduleID | uint64 | PATH | Module ID | N/A | YES |

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

---




# Triggers

Compose Triggers

| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
| `GET` | `/namespace/{namespaceID}/trigger/` | List available triggers |
| `POST` | `/namespace/{namespaceID}/trigger/` | Create trigger |
| `GET` | `/namespace/{namespaceID}/trigger/{triggerID}` | Get trigger details |
| `POST` | `/namespace/{namespaceID}/trigger/{triggerID}` | Update trigger |
| `Delete` | `/namespace/{namespaceID}/trigger/{triggerID}` | Delete trigger |

## List available triggers

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/trigger/` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| moduleID | uint64 | GET | Filter triggers by module | N/A | NO |
| query | string | GET | Search query | N/A | NO |
| page | uint | GET | Page number (0 based) | N/A | NO |
| perPage | uint | GET | Returned items per page (default 50) | N/A | NO |
| namespaceID | uint64 | PATH | Namespace ID | N/A | YES |

## Create trigger

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/trigger/` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| moduleID | uint64 | POST | Module ID | N/A | NO |
| name | string | POST | Name | N/A | YES |
| actions | []string | POST | Actions that trigger this trigger | N/A | NO |
| enabled | bool | POST | Enabled | N/A | NO |
| source | string | POST | Trigger source code | N/A | NO |
| updatedAt | *time.Time | POST | Last update (or creation) date | N/A | NO |
| namespaceID | uint64 | PATH | Namespace ID | N/A | YES |

## Get trigger details

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/trigger/{triggerID}` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| triggerID | uint64 | PATH | Trigger ID | N/A | YES |
| namespaceID | uint64 | PATH | Namespace ID | N/A | YES |

## Update trigger

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/trigger/{triggerID}` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| triggerID | uint64 | PATH | Trigger ID | N/A | YES |
| namespaceID | uint64 | PATH | Namespace ID | N/A | YES |
| moduleID | uint64 | POST | Module ID | N/A | NO |
| name | string | POST | Name | N/A | YES |
| actions | []string | POST | Actions that trigger this trigger | N/A | NO |
| enabled | bool | POST | Enabled | N/A | NO |
| source | string | POST | Trigger source code | N/A | NO |

## Delete trigger

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/namespace/{namespaceID}/trigger/{triggerID}` | HTTP/S | Delete |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| triggerID | uint64 | PATH | Trigger ID | N/A | YES |
| namespaceID | uint64 | PATH | Namespace ID | N/A | YES |

---