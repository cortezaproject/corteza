import fs from 'fs'
import prettier from 'prettier'
import { fileURLToPath } from 'url'
import { dirname, join } from 'path'
import { loadEndpoints } from './spec.js'

const __dirname = dirname(fileURLToPath(import.meta.url))

const srcDir = join(__dirname, '../../src/corredor/helpers')
const lockPath = join(__dirname, 'helper-names.lock.json')

const check = process.argv.includes('--check')
const positional = process.argv.slice(2).filter(a => !a.startsWith('--'))
const serverPath = positional.length > 0 ? positional[0] : '../../server'

// ---------------------------------------------------------------------------
// The naming rule
// ---------------------------------------------------------------------------

// The six regular actions on a singular entity entrypoint get a friendly name.
// Every other action — and every action of an entrypoint that is not a singular
// entity — keeps the API client's own method name.
const FRIENDLY = {
  list: ep => `find${plural(pascal(ep))}`,
  read: ep => `find${pascal(ep)}ByID`,
  create: ep => `create${pascal(ep)}`,
  update: ep => `update${pascal(ep)}`,
  delete: ep => `delete${pascal(ep)}`,
  undelete: ep => `undelete${pascal(ep)}`,
}

// Entrypoints that are not a singular entity. A findXs/createX name over one of
// these reads as nonsense, or names a resource the endpoint does not address.
const NOT_AN_ENTITY = new Set([
  // a verb, not a resource
  'auth',
  'expression',
  'smtpConfigurationChecker',
  // one global thing, or a report drawn over many
  'settings',
  'stats',
  'actionlog',
  'locale',
  'mcp',
  'eventTypes',
  'constructLibrary',
  // the rule surface every resource shares
  'permissions',
  // several resources deep under one entrypoint
  'dml',
  'dataPrivacy',
  // keyed by the resource it labels, not by an ID of its own
  'label',
  // Corredor scripts; findAutomations would read as the automation service
  'automation',
  // already plural
  'queues',
  // gateway and DAL internals, addressed by the admin UI rather than by scripts
  'apigwRoute',
  'apigwFilter',
  'apigwProfiler',
  'dalDriver',
  'dalSensitivityLevel',
  'dalSchemaAlteration',
  'dalConnection',
])

// An action whose response carries the resource, so a class can be put around
// it. delete and undelete answer with nothing.
const CAST_ONE = new Set(['read', 'create', 'update'])

// The list filters a handle-shaped value is looked up by, in the order tried.
const LOOKUP_FIELDS = ['handle', 'slug', 'name', 'email']

function pascal(s) {
  return s.charAt(0).toUpperCase() + s.slice(1)
}

function plural(s) {
  if (/(s|x|z|ch|sh)$/.test(s)) return `${s}es`
  if (/[^aeiou]y$/.test(s)) return `${s.slice(0, -1)}ies`
  return `${s}s`
}

// ---------------------------------------------------------------------------
// Services
// ---------------------------------------------------------------------------

const services = [
  {
    namespace: 'system',
    className: 'System',
    apiProp: 'SystemAPI',
    curatedFile: 'system.ts',
    curatedClass: 'SystemHelper',
    // Entrypoint to the resource class its responses are cast into. An
    // entrypoint absent from here answers with the raw response.
    classes: {
      authClient: ['AuthClient', '../../system'],
      role: ['Role', '../../system'],
      userGroup: ['UserGroup', '../../system'],
      user: ['User', '../../system'],
      application: ['Application', '../../system'],
      reminder: ['Reminder', '../../system'],
      notification: ['Notification', '../../system'],
      template: ['Template', '../../system'],
      connection: ['Connection', '../../system'],
      llmProvider: ['LlmProvider', '../../system'],
      agent: ['Agent', '../../system'],
      chatbot: ['Chatbot', '../../system'],
      chatbotSession: ['ChatbotSession', '../../system'],
      project: ['Project', '../../system'],
      projectAiSystem: ['ProjectAiSystem', '../../system'],
      attachment: ['Attachment', '../../shared'],
    },
  },
  {
    namespace: 'compose',
    className: 'Compose',
    apiProp: 'ComposeAPI',
    curatedFile: 'compose.ts',
    curatedClass: 'ComposeHelper',
    classes: {
      namespace: ['Namespace', '../../compose'],
      page: ['Page', '../../compose'],
      pageLayout: ['PageLayout', '../../compose'],
      module: ['Module', '../../compose'],
      chart: ['Chart', '../../compose'],
      attachment: ['Attachment', '../../shared'],
      // Record needs its module to be built, which no endpoint response
      // carries; the curated helper casts records.
    },
  },
  {
    namespace: 'automation',
    className: 'Automation',
    apiProp: 'AutomationAPI',
    curatedFile: undefined,
    curatedClass: undefined,
    classes: {
      workflow: ['Workflow', '../../automation'],
      ngAutomation: ['NgAutomation', '../../automation'],
      function: ['Function as AutomationFunction', '../../automation'],
    },
  },
]

// ---------------------------------------------------------------------------
// The lock
// ---------------------------------------------------------------------------

const lockNotice = [
  'Method name per endpoint, keyed by "<entrypoint> <action>".',
  'A name here is reused verbatim: the generator appends, never recomputes, so a',
  'sharper naming rule can never rename a method a script already calls. Fix an',
  'awkward name by editing it here.',
  'An entry under "curated" collides with a hand-written helper method, which',
  'stands; no generated twin is emitted and the endpoint stays reachable through',
  'the raw API client.',
]

function readLock() {
  if (!fs.existsSync(lockPath)) {
    return { '//': lockNotice, names: {}, curated: {} }
  }

  return JSON.parse(fs.readFileSync(lockPath, 'utf8'))
}

// Members of a hand-written helper class, so a generated twin of one can be
// left out. Read from the file rather than listed here, so a curated method
// added later is honoured without touching the generator.
function curatedMembers(file) {
  if (!file) {
    return new Set()
  }

  const src = fs.readFileSync(join(srcDir, file), 'utf8')
  const member =
    /^ {2}(?:public |protected |private )?(?:readonly )?(?:static )?(?:async )?(?:get |set )?([A-Za-z_$][\w$]*)\s*[(:=<]/gm

  const out = new Set()
  let m
  while ((m = member.exec(src)) !== null) {
    if (m[1] !== 'constructor') {
      out.add(m[1])
    }
  }

  return out
}

// ---------------------------------------------------------------------------
// Reference resolution
// ---------------------------------------------------------------------------

// The entity's own ID is the last path parameter of whichever of these actions
// it has; the ones before it address its parents.
const OWN_ID_FROM = ['read', 'update', 'delete', 'undelete']

function describeEntities(endpoints) {
  const byEntrypoint = new Map()
  endpoints.forEach(e => {
    if (!byEntrypoint.has(e.entrypoint)) {
      byEntrypoint.set(e.entrypoint, new Map())
    }
    byEntrypoint.get(e.entrypoint).set(e.action, e)
  })

  const entities = {}

  byEntrypoint.forEach((actions, entrypoint) => {
    if (NOT_AN_ENTITY.has(entrypoint)) {
      return
    }

    const own = OWN_ID_FROM.map(a => actions.get(a)).find(e => e && e.pathParams.length > 0)
    if (!own) {
      return
    }

    const idField = own.pathParams[own.pathParams.length - 1]
    if (!/ID$/.test(idField)) {
      return
    }

    const list = actions.get('list')

    entities[entrypoint] = {
      idField,
      list: list ? list.fname : undefined,
      by: list ? LOOKUP_FIELDS.filter(f => list.params.includes(f)) : [],
      parents: list ? list.pathParams : [],
    }
  })

  return entities
}

// Which entity a path parameter names: the endpoint's own entrypoint when it is
// that entity's ID, otherwise the one entity whose ID goes by that name.
function refTargets(endpoints, entities) {
  const byIdField = new Map()
  Object.entries(entities).forEach(([entrypoint, e]) => {
    if (!byIdField.has(e.idField)) {
      byIdField.set(e.idField, [])
    }
    byIdField.get(e.idField).push(entrypoint)
  })

  const parentOf = new Map()
  byIdField.forEach((entrypoints, idField) => {
    if (entrypoints.length === 1) {
      parentOf.set(idField, entrypoints[0])
    }
  })

  return endpoints.map(e => {
    const own = entities[e.entrypoint]

    return e.pathParams
      .map(field => {
        if (own && field === own.idField) {
          return { field, entity: e.entrypoint }
        }

        const parent = parentOf.get(field)
        return parent ? { field, entity: parent } : undefined
      })
      .filter(Boolean)
  })
}

// ---------------------------------------------------------------------------
// Generation
// ---------------------------------------------------------------------------

const header = `// This is a generated file.
// Run \`pnpm run codegen\` in lib/js to rebuild it.
//
// A method here mirrors one endpoint of the REST spec, so a change the server
// makes to that endpoint reaches scripts as soon as this file is regenerated.
// Over the raw API client a method adds two things and nothing else: every ID
// in the endpoint's path also accepts the object it names or its handle, and a
// response is cast into its resource class where one exists. The curated helper
// that extends this class is where a breaking endpoint change is absorbed.`

function methodSource(e, svc, name, resolves, entity) {
  const lines = []

  if (e.title) {
    lines.push(`  // ${e.title}`)
  }
  if (e.description) {
    e.description
      .trim()
      .split('\n')
      .forEach(l => lines.push(`  // ${l.trim()}`))
  }
  if (name !== e.fname) {
    lines.push(`  // Mirrors ${svc.apiProp}.${e.fname}.`)
  }

  const hasArgs = e.fargs.length > 0
  const soleID =
    e.pathParams.length === 1 && /ID$/.test(e.pathParams[0]) ? e.pathParams[0] : undefined

  const cls = entity && svc.classes[e.entrypoint] ? svc.classes[e.entrypoint][0] : undefined
  const alias = cls ? cls.split(' as ').pop() : undefined

  let castKind
  if (alias && e.action === 'list') {
    castKind = 'set'
  } else if (alias && CAST_ONE.has(e.action)) {
    castKind = 'one'
  }

  let returns = 'Promise<KV>'
  if (castKind === 'set') {
    returns = `Promise<ListResponse<${alias}[], KV>>`
  } else if (castKind === 'one') {
    returns = `Promise<${alias}>`
  }

  const args = hasArgs ? `a: ${soleID ? 'IdArgs' : 'Args'} = {}, ` : ''
  const asyncWord = resolves.length > 0 ? 'async ' : ''
  lines.push(`  ${asyncWord}${name}(${args}extra: AxiosRequestConfig = {}): ${returns} {`)

  const body = []
  if (hasArgs) {
    body.push(`const args = argsOf(a${soleID ? `, '${soleID}'` : ''})`)
    resolves.forEach(r => {
      body.push(`args.${r.field} = await this.resolveRef('${r.entity}', args.${r.field}, args)`)
    })
  }

  const call = `this.${svc.apiProp}.${e.fname}(${hasArgs ? 'args, ' : ''}extra)`
  if (castKind === 'set') {
    body.push(`return ${call}.then(r => castSet(r, v => new ${alias}(v)))`)
  } else if (castKind === 'one') {
    body.push(`return ${call}.then(r => new ${alias}(r))`)
  } else {
    body.push(`return ${call}`)
  }

  body.forEach(l => lines.push(`    ${l}`))
  lines.push('  }')

  return lines.join('\n')
}

function generate(svc, endpoints, lock) {
  const entities = describeEntities(endpoints)
  const targets = refTargets(endpoints, entities)
  const curated = curatedMembers(svc.curatedFile)

  const names = lock.names[svc.namespace] || {}
  const collisions = {}
  const emitted = new Map()
  const used = new Set(['constructor', 'resolveRef', svc.apiProp])
  const methods = []
  const usedClasses = new Set()

  endpoints.forEach((e, i) => {
    const key = `${e.entrypoint} ${e.action}`
    const entity = !NOT_AN_ENTITY.has(e.entrypoint)

    let name = names[key]
    if (!name) {
      const friendly = entity ? FRIENDLY[e.action] : undefined
      name = friendly ? friendly(e.entrypoint) : e.fname
    }
    emitted.set(key, name)

    if (curated.has(name)) {
      collisions[key] = `${svc.curatedClass}.${name}`
      return
    }

    if (used.has(name)) {
      throw new Error(`${svc.namespace}: method name '${name}' is taken (endpoint '${key}')`)
    }
    used.add(name)

    const cls = svc.classes[e.entrypoint]
    const casts = e.action === 'list' || CAST_ONE.has(e.action)
    if (entity && cls && casts) {
      usedClasses.add(cls.join('\u0000'))
    }

    methods.push(methodSource(e, svc, name, targets[i], entity))
  })

  const imports = new Map()
  usedClasses.forEach(v => {
    const [cls, from] = v.split('\u0000')
    if (!imports.has(from)) {
      imports.set(from, [])
    }
    imports.get(from).push(cls)
  })

  const refEntries = Object.entries(entities).map(([entrypoint, r]) => {
    const list = r.list ? `'${r.list}'` : 'undefined'
    const by = r.by.map(f => `'${f}'`).join(', ')
    const parents = r.parents.map(f => `'${f}'`).join(', ')
    return `  ${entrypoint}: { idField: '${r.idField}', list: ${list}, by: [${by}], parents: [${parents}] },`
  })

  const out = [
    header,
    '',
    `import { AxiosRequestConfig } from 'axios'`,
    `import { ${svc.className} as ${svc.apiProp} } from '../../api-clients'`,
    ...[...imports.entries()]
      .sort()
      .map(([from, cc]) => `import { ${cc.sort().join(', ')} } from '${from}'`),
    `import { Args, IdArgs, KV, ListResponse, RefSpec, argsOf, castSet, resolveRef } from './shared'`,
    '',
    '// Where a handle, a slug or a name arriving in place of an ID is looked up.',
    'const refs: RefSpec = {',
    ...refEntries,
    '}',
    '',
    `export default class Generated${svc.className}Helper {`,
    `  readonly ${svc.apiProp}: ${svc.apiProp}`,
    '',
    `  constructor(ctx: { ${svc.apiProp}: ${svc.apiProp} }) {`,
    `    this.${svc.apiProp} = ctx.${svc.apiProp}`,
    '  }',
    '',
    '  // Turns an object, an ID or a handle into the ID the endpoint wants.',
    '  protected resolveRef(entity: string, value: unknown, context: KV = {}): Promise<unknown> {',
    `    return resolveRef(this.${svc.apiProp}, refs, entity, value, context)`,
    '  }',
    '',
    methods.join('\n\n'),
    '}',
    '',
  ].join('\n')

  return { source: out, names: emitted, collisions }
}

// ---------------------------------------------------------------------------
// Run
// ---------------------------------------------------------------------------

const lock = readLock()
const outputs = []
const nextLock = { '//': lockNotice, names: {}, curated: {} }

for (const svc of services) {
  const path = `${serverPath}/${svc.namespace}/rest.yaml`
  const endpoints = loadEndpoints(path)

  if (!endpoints) {
    console.error(`Could not read endpoints from '${path}'`)
    process.exit(1)
  }

  const { source, names, collisions } = generate(svc, endpoints, lock)

  // Existing entries keep their order, new ones are appended.
  const locked = lock.names[svc.namespace] || {}
  const ordered = {}
  Object.keys(locked).forEach(k => {
    if (names.has(k)) {
      ordered[k] = locked[k]
    }
  })
  names.forEach((v, k) => {
    if (!(k in ordered)) {
      ordered[k] = v
    }
  })

  nextLock.names[svc.namespace] = ordered
  nextLock.curated[svc.namespace] = collisions

  const file = join(srcDir, `${svc.namespace}.gen.ts`)
  const cfg = await prettier.resolveConfig(file)
  outputs.push({ file, content: await prettier.format(source, { ...cfg, filepath: file }) })
}

outputs.push({ file: lockPath, content: `${JSON.stringify(nextLock, null, 2)}\n` })

if (check) {
  const stale = outputs.filter(
    o => !fs.existsSync(o.file) || fs.readFileSync(o.file, 'utf8') !== o.content,
  )

  if (stale.length > 0) {
    console.error('The generated Corredor helper layer no longer matches its sources:')
    stale.forEach(o => console.error(`  ${o.file}`))
    console.error('\nRegenerate with:\n  cd lib/js && pnpm run codegen')
    process.exit(1)
  }

  console.log('Corredor helper layer is up to date')
} else {
  outputs.forEach(o => {
    fs.writeFileSync(o.file, o.content)
    console.log(`Generated ${o.file}`)
  })
}
