# DML API – curl examples

```bash
TOKEN="your-jwt-token-here"
```

---

## Connections

### List connections

```bash
curl -v -X GET \
  -H "Authorization: Bearer $TOKEN" \
  "http://localhost:1024/api/system/dml/connections/" | jq .
```

### Read connection

```bash
curl -v -X GET \
  -H "Authorization: Bearer $TOKEN" \
  "http://localhost:1024/api/system/dml/connections/1" | jq .
```

---

## Models

### List models for a connection

```bash
curl -v -X GET \
  -H "Authorization: Bearer $TOKEN" \
  "http://localhost:1024/api/system/dml/connections/1/models" | jq .
```

### Read a model by ident

```bash
curl -v -X GET \
  -H "Authorization: Bearer $TOKEN" \
  "http://localhost:1024/api/system/dml/connections/1/models/compose_record" | jq .
```

---

## Migrations

### List migrations

```bash
curl -v -X GET \
  -H "Authorization: Bearer $TOKEN" \
  "http://localhost:1024/api/system/dml/migrations/" | jq .
```

### Create migration

Simple full-copy with attribute passthrough:

```bash
curl -v -X POST \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "handle": "crm-to-warehouse",
    "meta": {
      "short": "CRM → Warehouse",
      "description": "Copy CRM contacts into the analytics warehouse."
    },
    "strategy": "full-copy",
    "mappings": [
      {
        "handle": "contacts",
        "sourceConnectionID": "3",
        "sourceModelIdent": "crm_contact",
        "targetConnectionID": "2",
        "targetModelIdent": "wh_record_fact",
        "attrs": [
          { "source": "ID",    "target": "ID" },
          { "source": "Email", "target": "Email" },
          { "source": "FirstName", "target": "Payload" }
        ]
      }
    ]
  }' \
  "http://localhost:1024/api/system/dml/migrations/" | jq .
```

With an AST transform expression and source-side filter:

```bash
curl -v -X POST \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "handle": "users-to-warehouse2",
    "meta": {
      "short": "Users to Warehouse",
      "description": "Sync system users to the analytics warehouse."
    },
    "strategy": "upsert",
    "mappings": [
      {
        "handle": "users",
        "sourceConnectionID": "1",
        "sourceModelIdent": "users",
        "targetConnectionID": "2",
        "targetModelIdent": "wh_record_fact",
        "attrs": [
          { "source": "ID",    "target": "ID" },
          { "source": "Email", "target": "Email" },
          {
            "source": "Name",
            "target": "Payload",
            "transform": {
              "ref": "ne",
              "args": [
                { "symbol": "Name" },
                { "value": { "@type": "String", "@value": "" } }
              ]
            }
          }
        ],
        "filter": {
          "ref": "isNull",
          "args": [
            { "symbol": "DeletedAt" }
          ]
        }
      }
    ]
  }' \
  "http://localhost:1024/api/system/dml/migrations/" | jq .
```

Multi-mapping migration (two models in one migration):

```bash
curl -v -X POST \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "handle": "full-sync",
    "meta": {
      "short": "Full sync",
      "description": "Copy both users and contacts to the warehouse."
    },
    "strategy": "truncate-then-copy",
    "mappings": [
      {
        "handle": "users",
        "sourceConnectionID": "1",
        "sourceModelIdent": "users",
        "targetConnectionID": "2",
        "targetModelIdent": "wh_record_fact",
        "attrs": [
          { "source": "ID",    "target": "ID" },
          { "source": "Email", "target": "Email" }
        ]
      },
      {
        "handle": "contacts",
        "sourceConnectionID": "3",
        "sourceModelIdent": "crm_contact",
        "targetConnectionID": "2",
        "targetModelIdent": "wh_record_fact",
        "attrs": [
          { "source": "ID",    "target": "ID" },
          { "source": "Email", "target": "Email" }
        ]
      }
    ]
  }' \
  "http://localhost:1024/api/system/dml/migrations/" | jq .
```

### Read migration

```bash
curl -v -X GET \
  -H "Authorization: Bearer $TOKEN" \
  "http://localhost:1024/api/system/dml/migrations/246800123456789" | jq .
```

## Strategy reference

| Value                | Behaviour                                      |
| -------------------- | ---------------------------------------------- |
| `full-copy`          | Append all source rows to target               |
| `upsert`             | Insert or update by primary key                |
| `truncate-then-copy` | Truncate target, then full-copy                |
| `dry-run`            | Read-only pass; reports counts, writes nothing |
