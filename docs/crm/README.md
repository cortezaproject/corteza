# Attachments

| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
| `GET` | `/attachment/{kind}/` | List, filter all page attachments |
| `GET` | `/attachment/{kind}/{attachmentID}` | Attachment details |
| `GET` | `/attachment/{kind}/{attachmentID}/original/{name}` | Serves attached file |
| `GET` | `/attachment/{kind}/{attachmentID}/preview.{ext}` | Serves preview of an attached file |

## List, filter all page attachments

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/attachment/{kind}/` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| pageID | uint64 | GET | Filter attachments by page ID | N/A | NO |
| moduleID | uint64 | GET | Filter attachments by mnodule ID | N/A | NO |
| recordID | uint64 | GET | Filter attachments by record ID | N/A | NO |
| fieldName | string | GET | Filter attachments by field name | N/A | NO |
| page | uint | GET | Page number (0 based) | N/A | NO |
| perPage | uint | GET | Returned items per page (default 50) | N/A | NO |
| kind | string | PATH | Attachment kind | N/A | YES |

## Attachment details

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/attachment/{kind}/{attachmentID}` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| attachmentID | uint64 | PATH | Attachment ID | N/A | YES |
| kind | string | PATH | Attachment kind | N/A | YES |

## Serves attached file

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/attachment/{kind}/{attachmentID}/original/{name}` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| download | bool | GET | Force file download | N/A | NO |
| attachmentID | uint64 | PATH | Attachment ID | N/A | YES |
| name | string | PATH | File name | N/A | YES |
| kind | string | PATH | Attachment kind | N/A | YES |

## Serves preview of an attached file

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/attachment/{kind}/{attachmentID}/preview.{ext}` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| attachmentID | uint64 | PATH | Attachment ID | N/A | YES |
| ext | string | PATH | Preview extension/format | N/A | YES |
| kind | string | PATH | Attachment kind | N/A | YES |

---




# Charts

| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
| `GET` | `/chart/` | List/read charts from module section |
| `POST` | `/chart/` | List/read charts from module section |
| `GET` | `/chart/{chartID}` | Read charts by ID from module section |
| `POST` | `/chart/{chartID}` | Add/update charts in module section |
| `DELETE` | `/chart/{chartID}` | Delete chart |

## List/read charts from module section

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/chart/` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |

## List/read charts from module section

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/chart/` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| config | sqlxTypes.JSONText | POST | Chart JSON | N/A | YES |
| name | string | POST | Chart name | N/A | YES |

## Read charts by ID from module section

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/chart/{chartID}` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| chartID | uint64 | PATH | Chart ID | N/A | YES |

## Add/update charts in module section

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/chart/{chartID}` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| chartID | uint64 | PATH | Chart ID | N/A | YES |
| config | sqlxTypes.JSONText | POST | Chart JSON | N/A | YES |
| name | string | POST | Chart name | N/A | YES |

## Delete chart

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/chart/{chartID}` | HTTP/S | DELETE |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| chartID | uint64 | PATH | Chart ID | N/A | YES |

---




# Modules

CRM module definitions

| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
| `GET` | `/module/` | List modules |
| `POST` | `/module/` | Create module |
| `GET` | `/module/{moduleID}` | Read module |
| `POST` | `/module/{moduleID}` | Update module |
| `DELETE` | `/module/{moduleID}` | Delete module |

## List modules

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/module/` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| query | string | GET | Search query | N/A | NO |

## Create module

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/module/` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| name | string | POST | Module Name | N/A | YES |
| fields | types.ModuleFieldSet | POST | Fields JSON | N/A | YES |
| meta | sqlxTypes.JSONText | POST | Module meta data | N/A | YES |

## Read module

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/module/{moduleID}` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| moduleID | uint64 | PATH | Module ID | N/A | YES |

## Update module

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/module/{moduleID}` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| moduleID | uint64 | PATH | Module ID | N/A | YES |
| name | string | POST | Module Name | N/A | YES |
| fields | types.ModuleFieldSet | POST | Fields JSON | N/A | YES |
| meta | sqlxTypes.JSONText | POST | Module meta data | N/A | YES |

## Delete module

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/module/{moduleID}` | HTTP/S | DELETE |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| moduleID | uint64 | PATH | Module ID | N/A | YES |

---




# Notifications

CRM Notifications

| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
| `POST` | `/notification/email` | Send email from the CRM |

## Send email from the CRM

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/notification/email` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| to | []string | POST | Email addresses or Crust user IDs | N/A | YES |
| cc | []string | POST | Email addresses or Crust user IDs | N/A | NO |
| replyTo | string | POST | Crust user ID or email address in reply-to field | N/A | NO |
| subject  | string | POST | Email subject | N/A | NO |
| content | sqlxTypes.JSONText | POST | Message content | N/A | YES |

---




# Pages

CRM module pages

| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
| `GET` | `/page/` | List available pages |
| `POST` | `/page/` | Create page |
| `GET` | `/page/{pageID}` | Get page details |
| `GET` | `/page/tree` | Get page all (non-record) pages, hierarchically |
| `POST` | `/page/{pageID}` | Update page |
| `POST` | `/page/{selfID}/reorder` | Reorder pages |
| `Delete` | `/page/{pageID}` | Delete page |
| `POST` | `/page/{pageID}/attachment` | Uploads attachment to page |

## List available pages

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/page/` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| selfID | uint64 | GET | Parent page ID | N/A | NO |

## Create page

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/page/` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| selfID | uint64 | POST | Parent Page ID | N/A | NO |
| moduleID | uint64 | POST | Module ID | N/A | NO |
| title | string | POST | Title | N/A | YES |
| description | string | POST | Description | N/A | NO |
| visible | bool | POST | Visible in navigation | N/A | NO |
| blocks | sqlxTypes.JSONText | POST | Blocks JSON | N/A | YES |

## Get page details

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/page/{pageID}` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| pageID | uint64 | PATH | Page ID | N/A | YES |

## Get page all (non-record) pages, hierarchically

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/page/tree` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |

## Update page

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/page/{pageID}` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| pageID | uint64 | PATH | Page ID | N/A | YES |
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
| `/page/{selfID}/reorder` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| selfID | uint64 | PATH | Parent page ID | N/A | YES |
| pageIDs | []string | POST | Page ID order | N/A | YES |

## Delete page

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/page/{pageID}` | HTTP/S | Delete |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| pageID | uint64 | PATH | Page ID | N/A | YES |

## Uploads attachment to page

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/page/{pageID}/attachment` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| pageID | uint64 | PATH | Page ID | N/A | YES |
| upload | *multipart.FileHeader | POST | File to upload | N/A | YES |

---




# Permissions

| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
| `GET` | `/permissions/effective` | Effective rules for current user |

## Effective rules for current user

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/permissions/effective` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| resource | string | GET | Show only rules for a specific resource | N/A | NO |

---




# Records

CRM records

| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
| `GET` | `/module/{moduleID}/record/report` | Generates report from module records |
| `GET` | `/module/{moduleID}/record/` | List/read records from module section |
| `POST` | `/module/{moduleID}/record/` | Create record in module section |
| `GET` | `/module/{moduleID}/record/{recordID}` | Read records by ID from module section |
| `POST` | `/module/{moduleID}/record/{recordID}` | Update records in module section |
| `DELETE` | `/module/{moduleID}/record/{recordID}` | Delete record row from module section |
| `POST` | `/module/{moduleID}/record/{recordID}/{fieldName}/attachment` | Uploads attachment and validates it against record field requirements |

## Generates report from module records

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/module/{moduleID}/record/report` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| metrics | string | GET | Metrics (eg: 'SUM(money), MAX(calls)') | N/A | NO |
| dimensions | string | GET | Dimensions (eg: 'DATE(foo), status') | N/A | YES |
| filter | string | GET | Filter (eg: 'DATE(foo) > 2010') | N/A | NO |
| moduleID | uint64 | PATH | Module ID | N/A | YES |

## List/read records from module section

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/module/{moduleID}/record/` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| filter | string | GET | Filtering condition | N/A | NO |
| page | int | GET | Page number (0 based) | N/A | NO |
| perPage | int | GET | Returned items per page (default 50) | N/A | NO |
| sort | string | GET | Sort field (default id desc) | N/A | NO |
| moduleID | uint64 | PATH | Module ID | N/A | YES |

## Create record in module section

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/module/{moduleID}/record/` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| values | types.RecordValueSet | POST | Record values | N/A | YES |
| moduleID | uint64 | PATH | Module ID | N/A | YES |

## Read records by ID from module section

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/module/{moduleID}/record/{recordID}` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| recordID | uint64 | PATH | Record ID | N/A | YES |
| moduleID | uint64 | PATH | Module ID | N/A | YES |

## Update records in module section

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/module/{moduleID}/record/{recordID}` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| recordID | uint64 | PATH | Record ID | N/A | YES |
| moduleID | uint64 | PATH | Module ID | N/A | YES |
| values | types.RecordValueSet | POST | Record values | N/A | YES |

## Delete record row from module section

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/module/{moduleID}/record/{recordID}` | HTTP/S | DELETE |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| recordID | uint64 | PATH | Record ID | N/A | YES |
| moduleID | uint64 | PATH | Module ID | N/A | YES |

## Uploads attachment and validates it against record field requirements

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/module/{moduleID}/record/{recordID}/{fieldName}/attachment` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| recordID | uint64 | PATH | Record ID | N/A | YES |
| fieldName | string | PATH | Field name | N/A | YES |
| moduleID | uint64 | PATH | Module ID | N/A | YES |
| upload | *multipart.FileHeader | POST | File to upload | N/A | YES |

---




# Triggers

CRM Triggers

| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
| `GET` | `/trigger/` | List available triggers |
| `POST` | `/trigger/` | Create trigger |
| `GET` | `/trigger/{triggerID}` | Get trigger details |
| `POST` | `/trigger/{triggerID}` | Update trigger |
| `Delete` | `/trigger/{triggerID}` | Delete trigger |

## List available triggers

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/trigger/` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| moduleID | uint64 | GET | Filter triggers by module | N/A | NO |

## Create trigger

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/trigger/` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| moduleID | uint64 | POST | Module ID | N/A | NO |
| name | string | POST | Name | N/A | YES |
| actions | []string | POST | Actions that trigger this trigger | N/A | NO |
| enabled | bool | POST | Enabled | N/A | NO |
| source | string | POST | Trigger source code | N/A | NO |

## Get trigger details

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/trigger/{triggerID}` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| triggerID | uint64 | PATH | Trigger ID | N/A | YES |

## Update trigger

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/trigger/{triggerID}` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| triggerID | uint64 | PATH | Trigger ID | N/A | YES |
| moduleID | uint64 | POST | Module ID | N/A | NO |
| name | string | POST | Name | N/A | YES |
| actions | []string | POST | Actions that trigger this trigger | N/A | NO |
| enabled | bool | POST | Enabled | N/A | NO |
| source | string | POST | Trigger source code | N/A | NO |

## Delete trigger

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/trigger/{triggerID}` | HTTP/S | Delete |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| triggerID | uint64 | PATH | Trigger ID | N/A | YES |

---