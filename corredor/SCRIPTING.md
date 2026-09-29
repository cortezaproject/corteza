<!-- Generated from the sources listed below by corredor/tools/scripting-reference.mjs.
     Do not edit: run `pnpm --filter @planetcrust/human-corredor docs:scripting`. -->

# Corredor scripting reference

A Corredor script is an ES module with a default export. Corredor loads it
out of an extension directory, parses the export to learn what it binds to, and
runs its `exec` when the binding matches. Server scripts run inside Corredor
against the API; client scripts are bundled and run in the browser.

`dev/fixtures/corredor` is a working extension of one script per shape — a
before-create trigger, a manual button, an interval with `runAs`, an iterator, a
sink, and client scripts for both bundles. Start from the nearest one.

Generated from:

- `client/web/unify/src/plugins/corredor.js`
- `corredor/src/scripts/iterator.ts`
- `corredor/src/scripts/parser.ts`
- `corredor/src/scripts/trigger.ts`
- `corredor/src/server.ts`
- `lib/js/src/api-clients/automation.ts`
- `lib/js/src/api-clients/compose.ts`
- `lib/js/src/api-clients/index.ts`
- `lib/js/src/api-clients/system.ts`
- `lib/js/src/corredor/args-human.ts`
- `lib/js/src/corredor/ctx.ts`
- `lib/js/src/corredor/helpers/automation.gen.ts`
- `lib/js/src/corredor/helpers/compose.gen.ts`
- `lib/js/src/corredor/helpers/compose.ts`
- `lib/js/src/corredor/helpers/system.gen.ts`
- `lib/js/src/corredor/helpers/system.ts`
- `lib/js/tools/codegen/helper-names.lock.json`
- `lib/vue/src/corredor/compose-ui.ts`
- `lib/vue/src/corredor/script-bus.ts`
- `server/compose/service/event/events.yaml`
- `server/system/service/event/events.yaml`

## The script module

`label` and `description` are what the admin Scripts screen and a page
button show. Exactly one of `triggers` or `iterator` is required — `iterator`
wins when both are given. `exec` receives the event arguments as its first
parameter and the execution context as its second, and may be `async`.

What `exec` returns decides what the event does with it:

- a plain object is spread over the event's arguments, so `return $record`
  writes the record back;
- `false` aborts — the dispatching operation fails with `Aborted`, which is how
  a `before*` trigger refuses a save;
- anything else lands under `result`, and an event whose first mutable argument
  can take it decodes it there;
- `undefined` changes nothing.

A value is only read back when the argument is mutable. See **Arguments**.

```js
export default {
  label: 'What the Scripts screen and a page button call it',
  description: 'One line on what it does',

  security: { runAs: …, allow: …, deny: … },

  triggers({ on, before, after }) {
    return on(…).for(…).where(…)
  },

  exec(args, ctx) {
    // …
  },
}
```

The trigger definition is parsed on its own, ahead of the module: the parser
lifts the default export out of the file and evaluates that object literal in a
bare VM context. So `triggers`, `iterator`, `label`, `description` and
`security` cannot reference an import, a module-scope constant or anything else
outside the export — a script whose trigger does is loaded with an error and
never runs. `exec` is not parsed this way; it may use imports freely.

## Triggers

`triggers` is a function taking the trigger builder and returning one
trigger or an array of them; a generator function that `yield`s them works too.
Every method returns a new trigger, so a chain reads left to right and a script
can bind to several things at once.

- `on(…)`, `before(…)`, `after(…)` name the event types: each name is
  capitalised and prefixed, so `before('create')` is `beforeCreate`.
- `for(…)` names the resources. Without it the resource is the fallback below,
  which is why `on('manual')` alone binds to the whole system.
- `where(…)` adds a constraint: `where(name, value)`, `where(name, op, value)`,
  or a value array. Only the names the resource declares are matched.
- `uiProp(name, value)` decorates a manual trigger for the webapp: `app`,
  `page`, `slot`, `label`, `variant`. A trigger naming an app the webapp does
  not serve is skipped.
- `at(…)` and `every(…)` are deferred triggers and ignore the rest of the
  chain: each returns a fresh trigger on the fallback resource with the
  timestamps or crontab expressions as its constraint.

Builder methods: `on`, `before`, `after`, `for`, `at`, `every`, `where`, `uiProp`.

Fallback resource when `for(…)` is left out: `system`.

## The iterator form

`iterator` replaces `triggers` when the script should walk a set of
resources rather than answer an event. It is a function taking `each` and
returning one iterator. `exec` then runs once per resource, with the same
arguments an event of that resource would carry.

`action` decides what happens to each resource after `exec`: `update` saves it,
`clone` makes a new one, `delete` removes it, and the empty default only reads.
Returning `false` skips that one resource. Without `every()` or `at()` the
iterator is manual and runs when something asks for it.

The `filter` is the resource's own list filter, and `query` is a server-side
expression, not JSON: string literals in it take **single** quotes, which the
surrounding double-quoted JavaScript string leaves alone.

```js
iterator(each) {
  return each({
    resourceType: …, // string
    action: …, // 'update' | 'delete' | 'clone' | ''
    filter: …, // IteratorFilter
  }).every('* * * * *')
},
```

`action`: `update`, `delete`, `clone`, or the empty string to only read.

Deferring methods: `every`, `at`.

Standard `filter` keys: `query`, `sort`, `limit`, `offset`, plus the resource's own.

## Security

`security` is optional. As a bare string it means `runAs`; as an object it
takes the keys below, each role a handle or an ID.

- `runAs` runs the script as that user instead of whoever triggered it. It
  needs `CORREDOR_RUN_AS_ENABLED` on the server (the default is on); with it
  off, the script errors instead of running. A deferred script — `every()`,
  `at()`, or an iterator with either — has no triggering user, so it can only
  reach anything with `runAs`.
- `allow` lists the roles that may run the script on top of whoever already
  may. `deny` lists the roles that may not, and wins over `allow` and over the
  roles that otherwise bypass the check.
- A script with no `security` is offered to anyone the resource's own
  permissions already let through.

A script's API calls carry the token of whoever triggered it, or of the
`runAs` user when the script declares one — the server mints it per execution
and the context's API clients send it on every call. So a script reaches exactly
what that user reaches, and a permission error inside a script is that user's,
not Corredor's. The argument `$authUser` is the identity the token belongs to,
while `$invoker` stays the user who triggered the script even under `runAs`.
(`ctx.$authUser` answers the same question a different way: it decodes the token
and fetches that user, so it is a promise and it costs a request.)

## The execution context

`exec`'s second parameter exposes:

| name                  | type                    |
| --------------------- | ----------------------- |
| `ctx.SystemAPI`       | `apiClients.System`     |
| `ctx.ComposeAPI`      | `apiClients.Compose`    |
| `ctx.AutomationAPI`   | `apiClients.Automation` |
| `ctx.console`         | `Logger`                |
| `ctx.log`             | `Logger`                |
| `ctx.$authUser`       | `Promise<User>`         |
| `ctx.System`          | `SystemHelper`          |
| `ctx.Compose`         | `ComposeHelper`         |
| `ctx.Automation`      | `AutomationHelper`      |
| `ctx.frontendBaseURL` | `string \| undefined`   |

API clients live in `lib/js/src/api-clients/` (`System`, `Compose`, `Federation`, `Automation`, `Discovery`); the context
configures `SystemAPI`, `ComposeAPI`, `AutomationAPI` with the script's own token.

A helper is a client with two conveniences: it takes an object, an ID or a
handle wherever the endpoint wants an ID, and it casts what comes back into a
library class. Every endpoint of the three services has one, generated from the
same definitions the clients are, so the layer cannot fall behind the API.

The generated method mirrors its endpoint, so a server change reaches scripts.
The hand-written methods listed below are the curated layer: they carry judgement
a generator cannot — defaults from the event, several ways to name a thing, and
shapes with no endpoint at all — and they are where a breaking change is absorbed.
Where both would answer to one name the hand-written one stands.

`ctx.System` (19 methods): `findUsers`, `findUserByID`, `findUserByEmail`, `findUserByHandle`, `saveUser`, `setPassword`, `deleteUser`, `findRoles`, `findRoleByID`, `findRoleByHandle`, `saveRole`, `deleteRole`, `addUserToRole`, `removeUserFromRole`, `resolveUser`, `resolveRole`, `allow`, `deny`, `inherit`.

It defaults from `$user`, `$role`, `$application` when the event carries them.

`ctx.Compose` (35 methods): `makePage`, `savePage`, `deletePage`, `findPages`, `findPageByID`, `makeRecord`, `saveRecord`, `deleteRecord`, `findRecords`, `findLastRecord`, `findFirstRecord`, `findRecordByID`, `findAttachmentByID`, `moduleFieldNameFromLabel`, `makeModule`, `saveModule`, `findModules`, `findModuleByID`, `findModuleByName`, `findModuleByHandle`, `makeNamespace`, `saveNamespace`, `findNamespaces`, `findNamespaceByID`, `findNamespaceBySlug`, `sendMail`, `sendRecordToMail`, `walkFields`, `recordToHTML`, `recordToPlainText`, `resolveModule`, `resolveNamespace`, `allow`, `deny`, `inherit`.

It defaults from `$namespace`, `$module`, `$record` when the event carries them.

Generated alongside them, `ctx.System`, `ctx.Compose` and `ctx.Automation`
have one method per endpoint.
A regular list, read, create, update, delete or undelete reads as
`findUsers`, `findUserByID`, `createUser`, `updateUser`, `deleteUser`,
`undeleteUser`; everything else keeps the client method name, so
`ctx.Automation.ngAutomationExec()` runs a TAQ and `ctx.System.userSuspend()`
suspends a user. `lib/js/tools/codegen/helper-names.lock.json` is the list,
and a name in it never changes.

These 17 names belong to the hand-written layer, which takes its own
arguments rather than one object: `deletePage`, `deleteRecord`, `findAttachmentByID`, `findModuleByID`, `findModules`, `findNamespaceByID`, `findNamespaces`, `findPageByID`, `findPages`, `findRecordByID`, `findRecords`, `deleteRole`, `deleteUser`, `findRoleByID`, `findRoles`, `findUserByID`, `findUsers`.

## Arguments

Alongside the resource’s own, every execution carries:

| argument       | type           | meaning                                  |
| -------------- | -------------- | ---------------------------------------- |
| `$invoker`     | `User`, frozen | who triggered the script, `runAs` or not |
| `$authUser`    | `User`, frozen | whose token the script carries           |
| `authToken`    | raw `string`   | the JWT the context’s clients send       |
| `eventType`    | raw `string`   | the event type that matched              |
| `resourceType` | raw `string`   | the resource that dispatched it          |

The identity three are present only when the execution has one: a script
triggered by a user has `$invoker`, one with `runAs` has `$authUser` and
`authToken` for the run-as user, and a deferred script without `runAs` has
none of them — its API calls go out unauthenticated and fail.

An argument the library can cast arrives `$`-prefixed as that class, with the
uncast value beside it under `raw<Name>`; anything else arrives under its plain
name as decoded JSON. An immutable argument is never read back, and the service
dispatching the event may make the whole event immutable — every `compose:record`
`after*` event is, so an after trigger cannot change the record.

## Server resources

What a server script can bind a trigger to.

| resource                      | event types                                                                                                                                                                                  |
| ----------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `system`                      | `onManual`, `onInterval`, `onTimestamp`                                                                                                                                                      |
| `system:sink`                 | `onRequest`                                                                                                                                                                                  |
| `system:mail`                 | `onManual`, `onReceive`, `onSend`                                                                                                                                                            |
| `system:auth`                 | `beforeLogin`, `afterLogin`, `beforeSignup`, `afterSignup`                                                                                                                                   |
| `system:auth-client`          | `onManual`, `beforeCreate`, `afterCreate`, `beforeUpdate`, `afterUpdate`, `beforeDelete`, `afterDelete`                                                                                      |
| `system:user`                 | `onManual`, `beforeCreate`, `afterCreate`, `beforeUpdate`, `afterUpdate`, `beforeDelete`, `afterDelete`, `beforeSuspend`, `afterSuspend`                                                     |
| `system:role`                 | `onManual`, `beforeCreate`, `afterCreate`, `beforeUpdate`, `afterUpdate`, `beforeDelete`, `afterDelete`                                                                                      |
| `system:role:member`          | `beforeAdd`, `afterAdd`, `beforeRemove`, `afterRemove`                                                                                                                                       |
| `system:user-group`           | `onManual`, `beforeCreate`, `afterCreate`, `beforeUpdate`, `afterUpdate`, `beforeDelete`, `afterDelete`, `beforeMemberAdd`, `afterMemberAdd`, `beforeMemberRemove`, `afterMemberRemove`      |
| `system:application`          | `onManual`, `beforeCreate`, `afterCreate`, `beforeUpdate`, `afterUpdate`, `beforeDelete`, `afterDelete`                                                                                      |
| `system:queue`                | `onMessage`                                                                                                                                                                                  |
| `system:data-privacy-request` | `onManual`, `beforeCreate`, `afterCreate`, `beforeUpdate`, `afterUpdate`, `beforeDelete`, `afterDelete`                                                                                      |
| `system:reminder`             | `onManual`, `beforeCreate`, `afterCreate`, `beforeUpdate`, `afterUpdate`, `beforeDelete`, `afterDelete`, `beforeDismiss`, `afterDismiss`, `beforeSnooze`, `afterSnooze`                      |
| `compose`                     | `onManual`, `onInterval`, `onTimestamp`                                                                                                                                                      |
| `compose:namespace`           | `onManual`, `beforeCreate`, `afterCreate`, `beforeUpdate`, `afterUpdate`, `beforeDelete`, `afterDelete`                                                                                      |
| `compose:page`                | `onManual`, `beforeCreate`, `afterCreate`, `beforeUpdate`, `afterUpdate`, `beforeDelete`, `afterDelete`                                                                                      |
| `compose:page-layout`         | `onManual`, `beforeCreate`, `afterCreate`, `beforeUpdate`, `afterUpdate`, `beforeDelete`, `afterDelete`                                                                                      |
| `compose:module`              | `onManual`, `beforeCreate`, `afterCreate`, `beforeUpdate`, `afterUpdate`, `beforeDelete`, `afterDelete`                                                                                      |
| `compose:record`              | `onManual`, `onIteration`, `beforeCreate`, `afterCreate`, `beforeUpdate`, `afterUpdate`, `beforeDelete`, `afterDelete`, `beforeUndelete`, `afterUndelete`, `beforeOrganize`, `afterOrganize` |

### System

#### `system`

Event types: `onManual`, `onInterval`, `onTimestamp`

No arguments of its own.

No constraints — `where(…)` has nothing to match on.

#### `system:sink`

Event types: `onRequest`

| argument    | type                                         | notes                                                              |
| ----------- | -------------------------------------------- | ------------------------------------------------------------------ |
| `$response` | `SinkResponse` (`*types.SinkResponse`)       | read back on a mutable event; `rawResponse` holds the uncast value |
| `$request`  | `SinkRequest`, frozen (`*types.SinkRequest`) | immutable — never read back; `rawRequest` holds the uncast value   |

Constraints for `where(…)`: `request.host`, `request.remote-address`, `request.method`, `request.path`, `request.username`, `request.password`, `request.content-type`, `request.get.*`, `request.post.*`, `request.header.*`

#### `system:mail`

Event types: `onManual`, `onReceive`, `onSend`

| argument  | type                     | notes                        |
| --------- | ------------------------ | ---------------------------- |
| `message` | raw `*types.MailMessage` | read back on a mutable event |

Constraints for `where(…)`: `message.header.subject`, `message.header.from`, `message.header.to`, `message.header.reply-to`, `message.header.cc`, `message.header.bcc`

#### `system:auth`

Event types: `beforeLogin`, `afterLogin`, `beforeSignup`, `afterSignup`

| argument   | type                      | notes                                                          |
| ---------- | ------------------------- | -------------------------------------------------------------- |
| `$user`    | `User` (`*types.User`)    | read back on a mutable event; `rawUser` holds the uncast value |
| `provider` | raw `*types.AuthProvider` | read back on a mutable event                                   |

Constraints for `where(…)`: `user.handle`, `user.email`

#### `system:auth-client`

Event types: `onManual`, `beforeCreate`, `afterCreate`, `beforeUpdate`, `afterUpdate`, `beforeDelete`, `afterDelete`

| argument        | type                    | notes                        |
| --------------- | ----------------------- | ---------------------------- |
| `authClient`    | raw `*types.AuthClient` | read back on a mutable event |
| `oldAuthClient` | raw `*types.AuthClient` | immutable — never read back  |

Constraints for `where(…)`: `auth-client.handle`

#### `system:user`

Event types: `onManual`, `beforeCreate`, `afterCreate`, `beforeUpdate`, `afterUpdate`, `beforeDelete`, `afterDelete`, `beforeSuspend`, `afterSuspend`

| argument   | type                           | notes                                                            |
| ---------- | ------------------------------ | ---------------------------------------------------------------- |
| `$user`    | `User` (`*types.User`)         | read back on a mutable event; `rawUser` holds the uncast value   |
| `$oldUser` | `User`, frozen (`*types.User`) | immutable — never read back; `rawOldUser` holds the uncast value |

Constraints for `where(…)`: `user.handle`, `user.email`

#### `system:role`

Event types: `onManual`, `beforeCreate`, `afterCreate`, `beforeUpdate`, `afterUpdate`, `beforeDelete`, `afterDelete`

| argument   | type                           | notes                                                            |
| ---------- | ------------------------------ | ---------------------------------------------------------------- |
| `$role`    | `Role` (`*types.Role`)         | read back on a mutable event; `rawRole` holds the uncast value   |
| `$oldRole` | `Role`, frozen (`*types.Role`) | immutable — never read back; `rawOldRole` holds the uncast value |

Constraints for `where(…)`: `role.handle`, `role.name`

#### `system:role:member`

Event types: `beforeAdd`, `afterAdd`, `beforeRemove`, `afterRemove`

| argument | type                   | notes                                                          |
| -------- | ---------------------- | -------------------------------------------------------------- |
| `$user`  | `User` (`*types.User`) | read back on a mutable event; `rawUser` holds the uncast value |
| `$role`  | `Role` (`*types.Role`) | read back on a mutable event; `rawRole` holds the uncast value |

Constraints for `where(…)`: `role.handle`, `role.name`, `user.handle`, `user.email`

#### `system:user-group`

Event types: `onManual`, `beforeCreate`, `afterCreate`, `beforeUpdate`, `afterUpdate`, `beforeDelete`, `afterDelete`, `beforeMemberAdd`, `afterMemberAdd`, `beforeMemberRemove`, `afterMemberRemove`

| argument       | type                   | notes                        |
| -------------- | ---------------------- | ---------------------------- |
| `userGroup`    | raw `*types.UserGroup` | read back on a mutable event |
| `oldUserGroup` | raw `*types.UserGroup` | immutable — never read back  |

Constraints for `where(…)`: `userGroup.handle`, `userGroup.name`

#### `system:application`

Event types: `onManual`, `beforeCreate`, `afterCreate`, `beforeUpdate`, `afterUpdate`, `beforeDelete`, `afterDelete`

| argument          | type                                         | notes                                                                   |
| ----------------- | -------------------------------------------- | ----------------------------------------------------------------------- |
| `$application`    | `Application` (`*types.Application`)         | read back on a mutable event; `rawApplication` holds the uncast value   |
| `$oldApplication` | `Application`, frozen (`*types.Application`) | immutable — never read back; `rawOldApplication` holds the uncast value |

Constraints for `where(…)`: `application.name`

#### `system:queue`

Event types: `onMessage`

| argument  | type                      | notes                        |
| --------- | ------------------------- | ---------------------------- |
| `payload` | raw `*types.QueueMessage` | read back on a mutable event |

Constraints for `where(…)`: `payload.queue`

#### `system:data-privacy-request`

Event types: `onManual`, `beforeCreate`, `afterCreate`, `beforeUpdate`, `afterUpdate`, `beforeDelete`, `afterDelete`

| argument                | type                            | notes                        |
| ----------------------- | ------------------------------- | ---------------------------- |
| `dataPrivacyRequest`    | raw `*types.DataPrivacyRequest` | read back on a mutable event |
| `oldDataPrivacyRequest` | raw `*types.DataPrivacyRequest` | immutable — never read back  |

Constraints for `where(…)`: `role.name`

#### `system:reminder`

Event types: `onManual`, `beforeCreate`, `afterCreate`, `beforeUpdate`, `afterUpdate`, `beforeDelete`, `afterDelete`, `beforeDismiss`, `afterDismiss`, `beforeSnooze`, `afterSnooze`

| argument      | type                  | notes                        |
| ------------- | --------------------- | ---------------------------- |
| `reminder`    | raw `*types.Reminder` | read back on a mutable event |
| `oldReminder` | raw `*types.Reminder` | immutable — never read back  |

Constraints for `where(…)`: `reminder.resource`, `reminder.assigned-to`

### Compose

#### `compose`

Event types: `onManual`, `onInterval`, `onTimestamp`

No arguments of its own.

No constraints — `where(…)` has nothing to match on.

#### `compose:namespace`

Event types: `onManual`, `beforeCreate`, `afterCreate`, `beforeUpdate`, `afterUpdate`, `beforeDelete`, `afterDelete`

| argument        | type                                     | notes                                                                 |
| --------------- | ---------------------------------------- | --------------------------------------------------------------------- |
| `$namespace`    | `Namespace` (`*types.Namespace`)         | read back on a mutable event; `rawNamespace` holds the uncast value   |
| `$oldNamespace` | `Namespace`, frozen (`*types.Namespace`) | immutable — never read back; `rawOldNamespace` holds the uncast value |

Constraints for `where(…)`: `namespace.handle`, `namespace.name`

#### `compose:page`

Event types: `onManual`, `beforeCreate`, `afterCreate`, `beforeUpdate`, `afterUpdate`, `beforeDelete`, `afterDelete`

| argument     | type                             | notes                                                              |
| ------------ | -------------------------------- | ------------------------------------------------------------------ |
| `$page`      | `Page` (`*types.Page`)           | read back on a mutable event; `rawPage` holds the uncast value     |
| `$oldPage`   | `Page`, frozen (`*types.Page`)   | immutable — never read back; `rawOldPage` holds the uncast value   |
| `$namespace` | `Namespace` (`*types.Namespace`) | immutable — never read back; `rawNamespace` holds the uncast value |
| `selected`   | raw `[]interface{}`              | immutable — never read back                                        |

Constraints for `where(…)`: `namespace.handle`, `namespace.name`, `page.handle`, `page.name`

#### `compose:page-layout`

Event types: `onManual`, `beforeCreate`, `afterCreate`, `beforeUpdate`, `afterUpdate`, `beforeDelete`, `afterDelete`

| argument        | type                             | notes                                                              |
| --------------- | -------------------------------- | ------------------------------------------------------------------ |
| `pageLayout`    | raw `*types.PageLayout`          | read back on a mutable event                                       |
| `oldPageLayout` | raw `*types.PageLayout`          | immutable — never read back                                        |
| `$namespace`    | `Namespace` (`*types.Namespace`) | immutable — never read back; `rawNamespace` holds the uncast value |
| `selected`      | raw `[]interface{}`              | immutable — never read back                                        |

Constraints for `where(…)`: `namespace.handle`, `namespace.name`, `pageLayout.handle`, `pageLayout.name`

#### `compose:module`

Event types: `onManual`, `beforeCreate`, `afterCreate`, `beforeUpdate`, `afterUpdate`, `beforeDelete`, `afterDelete`

| argument     | type                               | notes                                                              |
| ------------ | ---------------------------------- | ------------------------------------------------------------------ |
| `$module`    | `Module` (`*types.Module`)         | read back on a mutable event; `rawModule` holds the uncast value   |
| `$oldModule` | `Module`, frozen (`*types.Module`) | immutable — never read back; `rawOldModule` holds the uncast value |
| `$namespace` | `Namespace` (`*types.Namespace`)   | immutable — never read back; `rawNamespace` holds the uncast value |

Constraints for `where(…)`: `namespace.handle`, `namespace.name`, `module.handle`, `module.name`

#### `compose:record`

Event types: `onManual`, `onIteration`, `beforeCreate`, `afterCreate`, `beforeUpdate`, `afterUpdate`, `beforeDelete`, `afterDelete`, `beforeUndelete`, `afterUndelete`, `beforeOrganize`, `afterOrganize`

| argument            | type                               | notes                                                              |
| ------------------- | ---------------------------------- | ------------------------------------------------------------------ |
| `$record`           | `Record` (`*types.Record`)         | read back on a mutable event; `rawRecord` holds the uncast value   |
| `$oldRecord`        | `Record`, frozen (`*types.Record`) | immutable — never read back; `rawOldRecord` holds the uncast value |
| `$module`           | `Module` (`*types.Module`)         | immutable — never read back; `rawModule` holds the uncast value    |
| `$namespace`        | `Namespace` (`*types.Namespace`)   | immutable — never read back; `rawNamespace` holds the uncast value |
| `recordValueErrors` | raw `*types.RecordValueErrorSet`   | read back on a mutable event                                       |
| `selected`          | raw `[]interface{}`                | immutable — never read back                                        |
| `$page`             | `Page` (`*types.Page`)             | immutable — never read back; `rawPage` holds the uncast value      |
| `filter`            | raw `string`                       | immutable — never read back                                        |

Constraints for `where(…)`: `namespace.handle`, `namespace.name`, `module.handle`, `module.name`, `record.created-at`, `record.updated-at`, `record.deleted-at`, `record.values.*`

## Client scripts

A client script is the same module shape, bundled by Corredor and run in
the browser instead of inside Corredor. Its context is the webapp's own: the
API clients are the ones the session already holds, and a script in the
`compose` bundle additionally gets `ComposeUI`. The `AutomationAPI` getter is
not configured there and throws.

Client scripts bind to the server resources above through `onManual`, and to
the `ui:` resources the webapp dispatches itself. Those carry the page's
in-memory arguments, so a `beforeFormSubmit` script corrects the record that is
about to be saved, and returning `false` stops the save.

| resource                       | event types                                                                                                                              |
| ------------------------------ | ---------------------------------------------------------------------------------------------------------------------------------------- |
| `compose`                      | `onManual`                                                                                                                               |
| `compose:namespace`            | `onManual`                                                                                                                               |
| `compose:module`               | `onManual`                                                                                                                               |
| `compose:record`               | `onManual`                                                                                                                               |
| `compose:page`                 | `onManual`                                                                                                                               |
| `ui:compose`                   | `onManual`                                                                                                                               |
| `ui:compose:record-page`       | `onManual`, `beforeFormSubmit`, `onFormSubmitError`, `afterFormSubmit`, `beforeDelete`, `afterDelete`, `beforeUndelete`, `afterUndelete` |
| `ui:compose:admin-record-page` | `onManual`, `beforeFormSubmit`, `onFormSubmitError`, `afterFormSubmit`, `beforeDelete`, `afterDelete`, `beforeUndelete`, `afterUndelete` |
| `system`                       | `onManual`                                                                                                                               |
| `system:user`                  | `onManual`                                                                                                                               |
| `system:role`                  | `onManual`                                                                                                                               |

`ctx.ComposeUI`: `gotoRecordViewer`, `gotoRecordEditor`, `success`, `warning`.

Bundles: `compose`, `admin`, `unify`. `uiProp('app', …)` takes `compose`, `admin`, `unify`.

## Where scripts live

Corredor searches the paths in `CORREDOR_EXT_SEARCH_PATHS` for extension
directories and loads what it finds under them. It watches the directories it
found at start-up, so editing or adding a script under one of them reloads it;
a directory that did not exist then is picked up only on the next restart. The
server polls Corredor for the list, and the admin Scripts screen and
`GET /api/system/automation/` show what it holds — including a script that
failed to parse, with its error.

- `<search-path>/client-scripts/<bundle>/<path-to-script>/*.js`
- `<search-path>/server-scripts/<path-to-script>/*.js`
