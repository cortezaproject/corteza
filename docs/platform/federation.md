---
title: Federation
description: Share records between two Human servers — turning federation on, pairing servers, exposing a module, mapping it on the receiving side, the read access a peer needs, and how syncing behaves.
---

# Federation

**Federation** lets two Human servers share records. One server **exposes** a
module, choosing which of its fields to share. The other server **maps** that
module onto one of its own modules and from then on receives its records:
new ones are created, changed ones updated and deleted ones removed.
Federation runs between servers of different organisations as easily as between
two servers of one, because each side decides for itself what it shares and
where what it receives lands.

Federation is off by default and still marked experimental: the server says so
in its log when it starts with federation on.

## Turning it on

Federation is configured per server, in its environment:

- [`FEDERATION_ENABLED`](/reference/environment#FEDERATION_ENABLED) turns it
  on. The Admin Area then shows a **Federation** group, with **Nodes** and
  **Permissions**, and the module editor gains a **Federation Settings**
  button.
- [`FEDERATION_HOST`](/reference/environment#FEDERATION_HOST) is the host
  (and port, if not the default) under which other servers reach this one. It
  goes into the pairing link this server hands out.
- [`FEDERATION_LABEL`](/reference/environment#FEDERATION_LABEL) is the name
  this server introduces itself by when another server pairs with it.

Servers talk to each other over HTTPS, at `/api/federation` under the host:
for example `https://crm.example.com/api/federation`.

## Pairing two servers

A **node** is one server's record of another server it federates with. Two
servers are paired when each has a node for the other and both nodes are
**paired**. Pairing takes three steps, by an administrator on each side:

1. **On the first server**, open **Admin Area → Federation → Nodes**, choose
   **New Federated Node** and enter the other server's name, its **Base URL**
   (its federation address, as above) and a contact email. Save, then choose
   **Generate URI** and copy the link. Send it to the other server's
   administrator. The node's status is **pending**.
2. **On the second server**, open **Admin Area → Federation → Nodes**, choose
   **Pair Node**, paste the link and choose **Pair**. This creates a node for
   the first server and asks it to pair. The link works once.
3. **Back on the first server**, the node's status is now **pair_requested**.
   Open its actions and choose **Confirm pending pair**. Both nodes become
   **paired**.

The pairing link always uses `https`. If the first server is reached at a
different address, correct the node's **Base URL** on the second server before
choosing **Pair**.

Pairing creates a user on each server that the other server acts as: its email
is `<nodeID>@federation.corteza`, and it is a member of the **Federation**
role. Everything a peer reads from this server, it reads as that user.

## Sharing a module

### Exposing it

On the server that has the records, open the module in the namespace's admin
panel and choose **Federation Settings**. On the **Expose** tab, pick the
paired server to share with and select the fields to share, or **All Fields**.
**Copy settings from** reuses the selection already made for another server.
Choose **Save and close**.

Only the selected fields ever leave the server. A field that is not exposed is
removed from every record before it is sent.

### Granting read access

The peer reads the records as the **Federation** role, and nothing grants that
role access to them by default. On the exposing server, give the
**Federation** role, for the exposed module:

- **List, search or filter records** on the module. Without it, no records
  reach the other server.
- **Read field value on records** on the exposed fields (or on all of the
  module's fields). Without it, records arrive with those fields empty, and
  nothing reports an error.

These come on top of the default read access to namespaces, modules and
records. On a server where that default access has been tightened, the role
needs read access to the namespace, the module and its records as well.
Permissions are explained in [Administration](./administration#the-permission-model).

### Mapping it

On the receiving server, an exposed module shows up after the next structure
sync — within two minutes by default. Open the module the records should land
in and choose **Federation Settings**. On the **Map** tab, pick the server,
pick the federated module, then match each federated field to one of this
module's fields. Choose **Save and close**.

Only mapped fields are filled in. The receiving module can be shaped
differently from the original: its fields can have other names, and fields that
are not mapped stay empty on received records.

## How syncing works

Each server runs two sync cycles against every server it is paired with:

- **Structure sync** fetches the modules the other server exposes to it, every
  [`FEDERATION_SYNC_STRUCTURE_MONITOR_INTERVAL`](/reference/environment#FEDERATION_SYNC_STRUCTURE_MONITOR_INTERVAL)
  (two minutes by default).
- **Data sync** fetches the records of every mapped module that changed since
  the last successful data sync, every
  [`FEDERATION_SYNC_DATA_MONITOR_INTERVAL`](/reference/environment#FEDERATION_SYNC_DATA_MONITOR_INTERVAL)
  (one minute by default), a page of
  [`FEDERATION_SYNC_DATA_PAGE_SIZE`](/reference/environment#FEDERATION_SYNC_DATA_PAGE_SIZE)
  records at a time.

A record created on the exposing server is created in the mapped module, one
that is changed is updated, and one that is deleted is deleted. A received
record remembers where it came from in its meta data: `federation` holds the
other server's address and `federation_extrecord` the original record's ID.
Records written by federation itself are never sent on, so two servers that
share in both directions do not send each other's records back.

If the exposing server changes which fields of an already shared module it
exposes, the receiving server stops syncing that module's records: the new
structure no longer matches the mapping made for the old one. Syncing resumes
once the exposed fields match again.

## Permissions

**Admin Area → Federation → Permissions** holds the rules for federation itself:
who may create and pair nodes, manage a node, expose a module and map a shared
module. The operations are listed in the
[permission reference](/reference/permissions/federation).
