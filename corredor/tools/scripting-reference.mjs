#!/usr/bin/env node
//
// Generator for the Corredor scripting reference (corredor/SCRIPTING.md).
//
// Everything the reference states about the scripting surface is read out of
// the files that define it: the two event catalogues, the script parser with
// its trigger and iterator DSL, the execution context with its helpers and API
// clients, and the webapp's script bus. Prose that no file can state lives in
// this generator, so the output is reproducible.
//
//   node tools/scripting-reference.mjs           write the reference
//   node tools/scripting-reference.mjs --check   fail when it is out of date
//
// Every extractor asserts on what it found: a moved source or a renamed symbol
// stops the generator instead of thinning the reference.

import fs from 'fs'
import path from 'path'
import { fileURLToPath } from 'url'
import prettier from 'prettier'

const toolDir = path.dirname(fileURLToPath(import.meta.url))
const corredorDir = path.resolve(toolDir, '..')
const repoDir = path.resolve(corredorDir, '..')

const OUTPUT = 'corredor/SCRIPTING.md'

const SOURCES = {
  systemEvents: 'server/system/service/event/events.yaml',
  composeEvents: 'server/compose/service/event/events.yaml',
  caster: 'lib/js/src/corredor/args-human.ts',
  ctx: 'lib/js/src/corredor/ctx.ts',
  systemHelper: 'lib/js/src/corredor/helpers/system.ts',
  composeHelper: 'lib/js/src/corredor/helpers/compose.ts',
  apiClients: 'lib/js/src/api-clients/index.ts',
  systemClient: 'lib/js/src/api-clients/system.ts',
  automationClient: 'lib/js/src/api-clients/automation.ts',
  composeClient: 'lib/js/src/api-clients/compose.ts',
  trigger: 'corredor/src/scripts/trigger.ts',
  iterator: 'corredor/src/scripts/iterator.ts',
  parser: 'corredor/src/scripts/parser.ts',
  loaderPaths: 'corredor/src/server.ts',
  scriptBus: 'lib/vue/src/corredor/script-bus.ts',
  composeUI: 'lib/vue/src/corredor/compose-ui.ts',
  webappPlugin: 'client/web/unify/src/plugins/corredor.js',
}

function fail(msg) {
  throw new Error(`scripting-reference: ${msg}`)
}

function must(cond, msg) {
  if (!cond) {
    fail(msg)
  }
}

function read(key) {
  const rel = SOURCES[key]
  must(rel, `no source registered under ${key}`)
  const abs = path.join(repoDir, rel)
  must(fs.existsSync(abs), `source ${rel} is gone`)
  return fs.readFileSync(abs, 'utf8')
}

const cap = s => s.substring(0, 1).toUpperCase() + s.substring(1)

// ── event catalogues ────────────────────────────────────────────────────────

const EVENT_FIELDS = ['on', 'ba', 'props', 'constraints']
const PROP_KEYS = ['name', 'type', 'immutable', 'exprtype']

function unquote(v) {
  return v.trim().replace(/^['"]|['"]$/g, '')
}

function parseFlowSeq(raw, where) {
  const m = /^\[(.*)\]$/.exec(raw.trim())
  must(m, `${where}: expected an inline list, got ${raw}`)
  const values = m[1]
    .split(',')
    .map(unquote)
    .filter(v => v.length > 0)
  must(values.length > 0, `${where}: empty list`)
  return values
}

function assignKV(item, raw, where, allowed) {
  const m = /^([a-zA-Z]+):\s*(.+)$/.exec(raw)
  must(m, `${where}: expected "key: value", got ${raw}`)
  const [, key, value] = m
  must(allowed.includes(key), `${where}: unknown key ${key}`)
  item[key] = value === 'true' ? true : value === 'false' ? false : unquote(value)
}

// Reads one events.yaml. The shape is narrow on purpose: anything the catalogue
// grows that this does not know about stops the generator.
function parseEvents(src, file) {
  const resources = []
  let current = null
  let list = null
  let item = null

  src.split('\n').forEach((raw, i) => {
    if (!raw.trim() || raw.trim().startsWith('#')) {
      return
    }

    const where = `${file}:${i + 1}`
    const indent = raw.length - raw.trimStart().length
    const line = raw.trim()

    if (indent === 0) {
      must(/^[a-z][a-z0-9:-]*:$/.test(line), `${where}: expected a resource key, got ${line}`)
      current = { resource: line.slice(0, -1), on: [], ba: [], props: [], constraints: [] }
      resources.push(current)
      list = null
      item = null
      return
    }

    must(current, `${where}: field outside a resource`)

    if (indent === 2) {
      item = null
      const m = /^([a-z]+):\s*(.*)$/.exec(line)
      must(m, `${where}: expected "field:", got ${line}`)
      const [, field, rest] = m
      must(EVENT_FIELDS.includes(field), `${where}: unknown resource field ${field}`)

      if (field === 'on' || field === 'ba') {
        current[field] = parseFlowSeq(rest, where)
        list = null
      } else {
        must(rest === '', `${where}: ${field} takes a block list`)
        list = field
      }
      return
    }

    must(list, `${where}: list item outside props/constraints`)

    if (indent === 4) {
      must(line.startsWith('- '), `${where}: expected a list item, got ${line}`)
      item = {}
      current[list].push(item)
      assignKV(item, line.slice(2), where, PROP_KEYS)
      return
    }

    must(indent === 6 && item, `${where}: unexpected indent ${indent}`)
    assignKV(item, line, where, PROP_KEYS)
  })

  must(resources.length > 0, `${file}: no resources`)

  resources.forEach(r => {
    must(
      r.on.length > 0 || r.ba.length > 0,
      `${file}: ${r.resource} declares neither on: nor ba: events`,
    )
    r.props.forEach(p =>
      must(p.name && p.type, `${file}: ${r.resource} has a prop without name/type`),
    )
    r.constraints.forEach(c =>
      must(c.name, `${file}: ${r.resource} has a constraint without a name`),
    )
    r.eventTypes = [
      ...r.on.map(e => `on${cap(e)}`),
      ...r.ba.flatMap(e => [`before${cap(e)}`, `after${cap(e)}`]),
    ]
  })

  return resources
}

// ── the rest of the surface ─────────────────────────────────────────────────

// Which argument names arrive $-prefixed and cast, and to what.
function parseCaster() {
  const src = read('caster')
  const casts = new Map()
  const re = /HumanTypes\.set\('([^']+)',\s*(\w+)(?:\((\w+)\))?\)/g

  for (const [, name, fn, cls] of src.matchAll(re)) {
    const frozen = /Freezer$/.test(fn)
    const type = cls ?? (/^record/.test(fn) ? 'Record' : undefined)
    must(type, `${SOURCES.caster}: cannot tell the class ${name} is cast to`)
    casts.set(name, { type, frozen })
  }

  must(casts.has('record') && casts.has('user'), `${SOURCES.caster}: caster map looks empty`)
  return casts
}

// Getters of the Ctx class — the second argument of exec().
function parseCtx() {
  const src = read('ctx')
  const getters = [...src.matchAll(/^ {2}get (\$?\w+)\(\): ([^{]+?)\s*\{/gm)].map(
    ([, name, type]) => ({
      name,
      type: type.trim(),
    }),
  )

  must(getters.length >= 5, `${SOURCES.ctx}: found ${getters.length} context getters`)
  ;['SystemAPI', 'ComposeAPI', 'System', 'Compose', 'log'].forEach(n =>
    must(
      getters.some(g => g.name === n),
      `${SOURCES.ctx}: context no longer exposes ${n}`,
    ),
  )

  return getters
}

// Public methods and contextual props of a helper class.
function parseHelper(key) {
  const src = read(key)
  const methods = [...src.matchAll(/^ {2}(?:async )?([a-zA-Z$][\w$]*)\(/gm)]
    .map(([, name]) => name)
    .filter(name => name !== 'constructor')
  const props = [...src.matchAll(/^ {2}readonly (\$\w+)\??: (\w+)/gm)].map(([, name, type]) => ({
    name,
    type,
  }))

  must(methods.length > 5, `${SOURCES[key]}: found ${methods.length} helper methods`)
  return { methods, props }
}

function clientMethods(key) {
  const src = read(key)
  const methods = [...src.matchAll(/^ {2}async (\w+)\(/gm)].map(([, name]) => name)
  must(methods.length > 20, `${SOURCES[key]}: found ${methods.length} client methods`)
  return methods
}

// Chainable methods of the Trigger class, in declaration order.
function parseTriggerDSL() {
  const src = read('trigger')
  const body = /export class Trigger \{([\s\S]*?)\n\}/.exec(src)
  must(body, `${SOURCES.trigger}: no Trigger class`)

  const methods = [...body[1].matchAll(/^ {2}(\w+)\(/gm)]
    .map(([, name]) => name)
    .filter(name => name !== 'constructor')

  ;['on', 'before', 'after', 'for', 'at', 'every', 'where', 'uiProp'].forEach(n =>
    must(methods.includes(n), `${SOURCES.trigger}: Trigger no longer has ${n}()`),
  )

  const fallback = /const defaultResource = '([^']+)'/.exec(src)
  must(fallback, `${SOURCES.trigger}: no default resource`)

  return { methods, defaultResource: fallback[1] }
}

function parseIterator() {
  const src = read('iterator')
  const fields = [...src.matchAll(/^ {2}public (\w+)!?: ([^=\n]+?)(?: =.*)?$/gm)].map(
    ([, name, type]) => ({ name, type: type.trim() }),
  )
  const methods = [...src.matchAll(/^ {2}(\w+)\(\.\.\.value: string\[\]\): Iterator/gm)].map(
    ([, name]) => name,
  )
  const filterKeys = [...src.matchAll(/^ {2}(\w+):\s*[^\n]+$/gm)]
    .map(([, name]) => name)
    .filter(name => name !== 'public')

  const actions = fields.find(f => f.name === 'action')
  must(actions, `${SOURCES.iterator}: Iterator has no action field`)
  must(
    fields.some(f => f.name === 'filter'),
    `${SOURCES.iterator}: Iterator has no filter field`,
  )
  must(methods.includes('every') && methods.includes('at'), `${SOURCES.iterator}: no every()/at()`)

  return {
    fields,
    methods,
    filterKeys: [...new Set(filterKeys)],
    actions: actions.type.split('|').map(v => v.trim().replace(/'/g, '')),
  }
}

function parseSecurityKeys() {
  const src = read('parser')
  const body = /interface RawScriptSecurity \{([\s\S]*?)\n\}/.exec(src)
  must(body, `${SOURCES.parser}: no RawScriptSecurity interface`)

  const keys = [...body[1].matchAll(/^ {2}(\w+)\??:/gm)].map(([, name]) => name)
  ;['runAs', 'allow', 'deny'].forEach(k =>
    must(keys.includes(k), `${SOURCES.parser}: security no longer takes ${k}`),
  )

  return keys
}

// Where the loader looks inside an extension directory.
function parseLoaderPatterns() {
  const src = read('loaderPaths')
  const patterns = [...src.matchAll(/^\/\/ (<search-path>\/[^\n]+)$/gm)].map(([, p]) => p)
  must(patterns.length === 2, `${SOURCES.loaderPaths}: expected 2 documented layout patterns`)
  return patterns
}

// Resource/event pairs the webapp dispatches to client scripts.
function parseScriptBus() {
  const src = read('scriptBus')
  const consts = new Map()

  for (const [, name, body] of src.matchAll(/^const (\w+) = \[([\s\S]*?)\]$/gm)) {
    consts.set(
      name,
      body
        .split(',')
        .map(unquote)
        .filter(v => v.length > 0),
    )
  }

  const onManual = /^const onManual = '([^']+)'$/m.exec(src)
  must(onManual, `${SOURCES.scriptBus}: no onManual constant`)

  const pairsBody = /wellKnownScriptPairs: eventbus\.WellKnownPairs = \{([\s\S]*?)\n\}/.exec(src)
  must(pairsBody, `${SOURCES.scriptBus}: no wellKnownScriptPairs`)

  const pairs = []
  for (const [, key, body] of pairsBody[1].matchAll(/^ {2}'?([\w:-]+)'?: \[(.*)\],$/gm)) {
    const eventTypes = body
      .split(',')
      .map(v => v.trim())
      .flatMap(v => {
        if (v === 'onManual') {
          return [onManual[1]]
        }
        const spread = /^\.\.\.(\w+)$/.exec(v)
        if (spread) {
          must(consts.has(spread[1]), `${SOURCES.scriptBus}: unknown spread ${v}`)
          return consts.get(spread[1])
        }
        return [unquote(v)]
      })
      .filter(v => v.length > 0)

    pairs.push({ resource: key, eventTypes })
  }

  must(pairs.length > 5, `${SOURCES.scriptBus}: found ${pairs.length} pairs`)
  must(
    pairs.some(p => p.resource.startsWith('ui:')),
    `${SOURCES.scriptBus}: no ui: resources`,
  )

  return pairs
}

function parseComposeUI() {
  const src = read('composeUI')
  const methods = [...src.matchAll(/^ {2}([a-zA-Z$][\w$]*)\(/gm)]
    .map(([, name]) => name)
    .filter(name => name !== 'constructor')
  must(methods.length > 2, `${SOURCES.composeUI}: found ${methods.length} methods`)
  return methods
}

function parseWebappBundles() {
  const src = read('webappPlugin')
  const apps = /^const apps = \[(.*)\]$/m.exec(src)
  must(apps, `${SOURCES.webappPlugin}: no apps list`)

  const bundles = /const bundles = \{([^}]*)\}/.exec(src)
  must(bundles, `${SOURCES.webappPlugin}: no bundle map`)

  return {
    apps: apps[1].split(',').map(unquote).filter(Boolean),
    bundles: [...bundles[1].matchAll(/(\w+):/g)].map(([, name]) => name),
  }
}

// ── the two claims the reference makes about reach ──────────────────────────

// What the convenience helpers cover, each with a method that proves it.
const HELPER_COVERAGE = [
  ['users', 'findUsers'],
  ['roles', 'findRoles'],
  ['permissions', 'allow'],
  ['records', 'findRecords'],
  ['modules', 'findModules'],
  ['namespaces', 'findNamespaces'],
  ['pages', 'findPages'],
  ['mail', 'sendMail'],
]

// What they do not cover, each with the client method that reaches it.
const RAW_ONLY = [
  ['automations (TAQ)', 'AutomationAPI', 'ngAutomationList'],
  ['workflows', 'AutomationAPI', 'workflowList'],
  ['agents', 'SystemAPI', 'agentList'],
  ['chatbots', 'SystemAPI', 'chatbotList'],
  ['projects', 'SystemAPI', 'projectList'],
  ['user groups', 'SystemAPI', 'userGroupList'],
  ['reminders', 'SystemAPI', 'reminderList'],
  ['notifications', 'SystemAPI', 'notificationList'],
  ['labels', 'SystemAPI', 'labelList'],
]

// Resource nouns the helpers must stay out of, as CRUD-shaped method names.
const NOT_IN_HELPERS =
  /^(?:find|save|make|delete|resolve|add|remove|set)\w*(Automation|Workflow|Agent|Chatbot|Project|UserGroup|Reminder|Notification|Label)s?$/

function checkReach(helperMethods, clients) {
  HELPER_COVERAGE.forEach(([what, method]) =>
    must(helperMethods.includes(method), `helpers no longer cover ${what} (${method}() is gone)`),
  )

  helperMethods.forEach(name => {
    const hit = NOT_IN_HELPERS.exec(name)
    must(!hit, `a helper now covers ${hit?.[1]} (${name}()) — the reference says it does not`)
  })

  RAW_ONLY.forEach(([what, client, method]) =>
    must(
      clients[client].includes(method),
      `${client}.${method}() is gone — the reference points at it for ${what}`,
    ),
  )
}

// ── prose ───────────────────────────────────────────────────────────────────
//
// The parts of the reference no source file states. Kept here so the whole
// output is one reproducible artefact.

const PROSE = {
  intro: `A Corredor script is an ES module with a default export. Corredor loads it
out of an extension directory, parses the export to learn what it binds to, and
runs its \`exec\` when the binding matches. Server scripts run inside Corredor
against the API; client scripts are bundled and run in the browser.

\`dev/fixtures/corredor\` is a working extension of one script per shape — a
before-create trigger, a manual button, an interval with \`runAs\`, an iterator, a
sink, and client scripts for both bundles. Start from the nearest one.`,

  module: `\`label\` and \`description\` are what the admin Scripts screen and a page
button show. Exactly one of \`triggers\` or \`iterator\` is required — \`iterator\`
wins when both are given. \`exec\` receives the event arguments as its first
parameter and the execution context as its second, and may be \`async\`.

What \`exec\` returns decides what the event does with it:

- a plain object is spread over the event's arguments, so \`return $record\`
  writes the record back;
- \`false\` aborts — the dispatching operation fails with \`Aborted\`, which is how
  a \`before*\` trigger refuses a save;
- anything else lands under \`result\`, and an event whose first mutable argument
  can take it decodes it there;
- \`undefined\` changes nothing.

A value is only read back when the argument is mutable. See **Arguments**.`,

  parsing: `The trigger definition is parsed on its own, ahead of the module: the parser
lifts the default export out of the file and evaluates that object literal in a
bare VM context. So \`triggers\`, \`iterator\`, \`label\`, \`description\` and
\`security\` cannot reference an import, a module-scope constant or anything else
outside the export — a script whose trigger does is loaded with an error and
never runs. \`exec\` is not parsed this way; it may use imports freely.`,

  triggers: `\`triggers\` is a function taking the trigger builder and returning one
trigger or an array of them; a generator function that \`yield\`s them works too.
Every method returns a new trigger, so a chain reads left to right and a script
can bind to several things at once.

- \`on(…)\`, \`before(…)\`, \`after(…)\` name the event types: each name is
  capitalised and prefixed, so \`before('create')\` is \`beforeCreate\`.
- \`for(…)\` names the resources. Without it the resource is the fallback below,
  which is why \`on('manual')\` alone binds to the whole system.
- \`where(…)\` adds a constraint: \`where(name, value)\`, \`where(name, op, value)\`,
  or a value array. Only the names the resource declares are matched.
- \`uiProp(name, value)\` decorates a manual trigger for the webapp: \`app\`,
  \`page\`, \`slot\`, \`label\`, \`variant\`. A trigger naming an app the webapp does
  not serve is skipped.
- \`at(…)\` and \`every(…)\` are deferred triggers and ignore the rest of the
  chain: each returns a fresh trigger on the fallback resource with the
  timestamps or crontab expressions as its constraint.`,

  iterator: `\`iterator\` replaces \`triggers\` when the script should walk a set of
resources rather than answer an event. It is a function taking \`each\` and
returning one iterator. \`exec\` then runs once per resource, with the same
arguments an event of that resource would carry.

\`action\` decides what happens to each resource after \`exec\`: \`update\` saves it,
\`clone\` makes a new one, \`delete\` removes it, and the empty default only reads.
Returning \`false\` skips that one resource. Without \`every()\` or \`at()\` the
iterator is manual and runs when something asks for it.

The \`filter\` is the resource's own list filter, and \`query\` is a server-side
expression, not JSON: string literals in it take **single** quotes, which the
surrounding double-quoted JavaScript string leaves alone.`,

  security: `\`security\` is optional. As a bare string it means \`runAs\`; as an object it
takes the keys below, each role a handle or an ID.

- \`runAs\` runs the script as that user instead of whoever triggered it. It
  needs \`CORREDOR_RUN_AS_ENABLED\` on the server (the default is on); with it
  off, the script errors instead of running. A deferred script — \`every()\`,
  \`at()\`, or an iterator with either — has no triggering user, so it can only
  reach anything with \`runAs\`.
- \`allow\` lists the roles that may run the script on top of whoever already
  may. \`deny\` lists the roles that may not, and wins over \`allow\` and over the
  roles that otherwise bypass the check.
- A script with no \`security\` is offered to anyone the resource's own
  permissions already let through.`,

  token: `A script's API calls carry the token of whoever triggered it, or of the
\`runAs\` user when the script declares one — the server mints it per execution
and the context's API clients send it on every call. So a script reaches exactly
what that user reaches, and a permission error inside a script is that user's,
not Corredor's. The argument \`$authUser\` is the identity the token belongs to,
while \`$invoker\` stays the user who triggered the script even under \`runAs\`.
(\`ctx.$authUser\` answers the same question a different way: it decodes the token
and fetches that user, so it is a promise and it costs a request.)`,

  helpers: `\`ctx.System\` and \`ctx.Compose\` are conveniences over the two clients:
they resolve handles, cast results to library classes, and default the module
and namespace from the event. They cover users, roles, permissions, records,
modules, namespaces, pages and mail — and nothing else. Automations, workflows,
agents, chatbots, projects, user groups, reminders, notifications and labels
have no helper; they are reached through the raw clients, by the method names in
\`lib/js/src/api-clients/\`.`,

  clientScripts: `A client script is the same module shape, bundled by Corredor and run in
the browser instead of inside Corredor. Its context is the webapp's own: the
API clients are the ones the session already holds, and a script in the
\`compose\` bundle additionally gets \`ComposeUI\`. The \`AutomationAPI\` getter is
not configured there and throws.

Client scripts bind to the server resources above through \`onManual\`, and to
the \`ui:\` resources the webapp dispatches itself. Those carry the page's
in-memory arguments, so a \`beforeFormSubmit\` script corrects the record that is
about to be saved, and returning \`false\` stops the save.`,

  deploy: `Corredor searches the paths in \`CORREDOR_EXT_SEARCH_PATHS\` for extension
directories and loads what it finds under them. It watches the directories it
found at start-up, so editing or adding a script under one of them reloads it;
a directory that did not exist then is picked up only on the next restart. The
server polls Corredor for the list, and the admin Scripts screen and
\`GET /api/system/automation/\` show what it holds — including a script that
failed to parse, with its error.`,
}

// ── rendering ───────────────────────────────────────────────────────────────

const HEADER = `<!-- Generated from the sources listed below by corredor/tools/scripting-reference.mjs.
     Do not edit: run \`pnpm --filter @planetcrust/human-corredor docs:scripting\`. -->`

const cell = v => String(v).replace(/\|/g, '\\|')

function table(head, rows) {
  return [
    `| ${head.join(' | ')} |`,
    `| ${head.map(() => '---').join(' | ')} |`,
    ...rows.map(r => `| ${r.map(cell).join(' | ')} |`),
  ].join('\n')
}

const code = list => list.map(v => `\`${v}\``).join(', ')

function argumentRow(prop, casts) {
  const cast = casts.get(prop.name)
  const name = cast ? `$${prop.name}` : prop.name
  const type = cast
    ? `\`${cast.type}\`${cast.frozen ? ', frozen' : ''} (\`${prop.type}\`)`
    : `raw \`${prop.type}\``
  const notes = []

  if (prop.immutable) {
    notes.push('immutable — never read back')
  } else {
    notes.push('read back on a mutable event')
  }

  if (cast) {
    notes.push(`\`raw${cap(prop.name)}\` holds the uncast value`)
  }

  return [`\`${name}\``, type, notes.join('; ')]
}

function resourceSection(r, casts) {
  const out = [`#### \`${r.resource}\``, '', `Event types: ${code(r.eventTypes)}`]

  if (r.props.length > 0) {
    out.push(
      '',
      table(
        ['argument', 'type', 'notes'],
        r.props.map(p => argumentRow(p, casts)),
      ),
    )
  } else {
    out.push('', 'No arguments of its own.')
  }

  if (r.constraints.length > 0) {
    out.push('', `Constraints for \`where(…)\`: ${code(r.constraints.map(c => c.name))}`)
  } else {
    out.push('', 'No constraints — `where(…)` has nothing to match on.')
  }

  return out.join('\n')
}

function render() {
  const system = parseEvents(read('systemEvents'), SOURCES.systemEvents)
  const compose = parseEvents(read('composeEvents'), SOURCES.composeEvents)
  const casts = parseCaster()
  const ctx = parseCtx()
  const systemHelper = parseHelper('systemHelper')
  const composeHelper = parseHelper('composeHelper')
  const dsl = parseTriggerDSL()
  const iterator = parseIterator()
  const securityKeys = parseSecurityKeys()
  const layout = parseLoaderPatterns()
  const bus = parseScriptBus()
  const composeUI = parseComposeUI()
  const webapp = parseWebappBundles()

  const clients = {
    SystemAPI: clientMethods('systemClient'),
    ComposeAPI: clientMethods('composeClient'),
    AutomationAPI: clientMethods('automationClient'),
  }

  checkReach([...systemHelper.methods, ...composeHelper.methods], clients)

  const clientNames = [...read('apiClients').matchAll(/^export \{ default as (\w+) \}/gm)].map(
    ([, n]) => n,
  )
  must(clientNames.length > 0, `${SOURCES.apiClients}: no clients exported`)

  const ctxClients = ctx.filter(g => g.type.startsWith('apiClients.'))
  const ctxRest = ctx.filter(g => !g.type.startsWith('apiClients.'))

  const commonArgs = [
    ['`$invoker`', '`User`, frozen', 'who triggered the script, `runAs` or not'],
    ['`$authUser`', '`User`, frozen', 'whose token the script carries'],
    ['`authToken`', 'raw `string`', 'the JWT the context’s clients send'],
    ['`eventType`', 'raw `string`', 'the event type that matched'],
    ['`resourceType`', 'raw `string`', 'the resource that dispatched it'],
  ]

  const md = []
  const add = (...blocks) => md.push(...blocks)

  add(HEADER, '', '# Corredor scripting reference', '', PROSE.intro, '')

  add(
    'Generated from:',
    '',
    ...Object.values(SOURCES)
      .sort()
      .map(s => `- \`${s}\``),
    '',
  )

  add('## The script module', '', PROSE.module, '', '```js', 'export default {')
  add("  label: 'What the Scripts screen and a page button call it',")
  add("  description: 'One line on what it does',", '')
  add(`  security: { ${securityKeys.map(k => `${k}: …`).join(', ')} },`, '')
  add('  triggers({ on, before, after }) {', '    return on(…).for(…).where(…)', '  },', '')
  add('  exec(args, ctx) {', '    // …', '  },', '}', '```', '')
  add(PROSE.parsing, '')

  add('## Triggers', '', PROSE.triggers, '')
  add(`Builder methods: ${code(dsl.methods)}.`, '')
  add(`Fallback resource when \`for(…)\` is left out: \`${dsl.defaultResource}\`.`, '')

  add('## The iterator form', '', PROSE.iterator, '')
  add(
    '```js',
    'iterator(each) {',
    '  return each({',
    ...iterator.fields
      .filter(f => ['resourceType', 'action', 'filter'].includes(f.name))
      .map(f => `    ${f.name}: …, // ${f.type}`),
    "  }).every('* * * * *')",
    '},',
    '```',
    '',
  )
  add(
    `\`action\`: ${code(iterator.actions.filter(Boolean))}, or the empty string to only read.`,
    '',
  )
  add(`Deferring methods: ${code(iterator.methods)}.`, '')
  add(`Standard \`filter\` keys: ${code(iterator.filterKeys)}, plus the resource's own.`, '')

  add('## Security', '', PROSE.security, '', PROSE.token, '')

  add('## The execution context', '', "`exec`'s second parameter exposes:", '')
  add(
    table(
      ['name', 'type'],
      [...ctxClients, ...ctxRest].map(g => [`\`ctx.${g.name}\``, `\`${g.type}\``]),
    ),
    '',
  )
  add(
    `API clients live in \`lib/js/src/api-clients/\` (${code(clientNames)}); the context`,
    `configures ${code(ctxClients.map(c => c.name))} with the script's own token.`,
    '',
  )
  add(PROSE.helpers, '')
  add(
    `\`ctx.System\` (${systemHelper.methods.length} methods): ${code(systemHelper.methods)}.`,
    '',
    `It defaults from ${code(systemHelper.props.map(p => p.name))} when the event carries them.`,
    '',
  )
  add(
    `\`ctx.Compose\` (${composeHelper.methods.length} methods): ${code(composeHelper.methods)}.`,
    '',
    `It defaults from ${code(composeHelper.props.map(p => p.name))} when the event carries them.`,
    '',
  )
  add(
    'Uncovered resources and where to reach them:',
    '',
    table(
      ['resource', 'client', 'example'],
      RAW_ONLY.map(([what, client, method]) => [what, `\`ctx.${client}\``, `\`${method}()\``]),
    ),
    '',
  )

  add('## Arguments', '', 'Alongside the resource’s own, every execution carries:', '')
  add(table(['argument', 'type', 'meaning'], commonArgs), '')
  add(
    'The identity three are present only when the execution has one: a script',
    'triggered by a user has `$invoker`, one with `runAs` has `$authUser` and',
    '`authToken` for the run-as user, and a deferred script without `runAs` has',
    'none of them — its API calls go out unauthenticated and fail.',
    '',
  )
  add(
    'An argument the library can cast arrives `$`-prefixed as that class, with the',
    'uncast value beside it under `raw<Name>`; anything else arrives under its plain',
    'name as decoded JSON. An immutable argument is never read back, and the service',
    'dispatching the event may make the whole event immutable — every `compose:record`',
    '`after*` event is, so an after trigger cannot change the record.',
    '',
  )

  add('## Server resources', '', 'What a server script can bind a trigger to.', '')
  add(
    table(
      ['resource', 'event types'],
      [...system, ...compose].map(r => [`\`${r.resource}\``, code(r.eventTypes)]),
    ),
    '',
  )
  add('### System', '')
  add(...system.map(r => `${resourceSection(r, casts)}\n`))
  add('### Compose', '')
  add(...compose.map(r => `${resourceSection(r, casts)}\n`))

  add('## Client scripts', '', PROSE.clientScripts, '')
  add(
    table(
      ['resource', 'event types'],
      bus.map(p => [`\`${p.resource}\``, code(p.eventTypes)]),
    ),
    '',
  )
  add(`\`ctx.ComposeUI\`: ${code(composeUI)}.`, '')
  add(`Bundles: ${code(webapp.bundles)}. \`uiProp('app', …)\` takes ${code(webapp.apps)}.`, '')

  add('## Where scripts live', '', PROSE.deploy, '')
  add(...layout.map(p => `- \`${p}\``), '')

  return `${md
    .join('\n')
    .replace(/\n{3,}/g, '\n\n')
    .trimEnd()}\n`
}

// ── entry point ─────────────────────────────────────────────────────────────

const target = path.join(repoDir, OUTPUT)
const generated = await prettier.format(render(), {
  ...(await prettier.resolveConfig(target)),
  parser: 'markdown',
})

if (process.argv.includes('--check')) {
  const current = fs.existsSync(target) ? fs.readFileSync(target, 'utf8') : ''

  if (current !== generated) {
    process.stderr.write(
      `${OUTPUT} is out of date.\n` +
        'Run: pnpm --filter @planetcrust/human-corredor docs:scripting\n',
    )
    process.exit(1)
  }

  process.stdout.write(`${OUTPUT} is up to date\n`)
} else {
  fs.writeFileSync(target, generated)
  process.stdout.write(`wrote ${OUTPUT}\n`)
}
