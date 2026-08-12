import { automation } from '@planetcrust/human-js'

// Every scope property the dry-run form knows how to render and resolve.
// `widget` picks the input component in DryRunField; `deps` names sibling
// scope fields whose picked IDs are needed to render (cascade) and resolve
// this one; `resolve` fetches the full object injected into the scope.
// Scope properties without a definition have no sensible form input — they
// are hidden and auto-initialized as empty variables at encode time.
const SCOPE_DEFS = {
  namespace: {
    widget: 'namespace',
    resolve: ({ value }, deps, { ComposeAPI }) =>
      ComposeAPI.namespaceRead({ namespaceID: value }).then(ns => ({
        ...ns,
        resourceType: 'compose:namespace',
      })),
  },
  module: {
    widget: 'module',
    deps: ['namespace'],
    resolve: ({ value }, deps, { ComposeAPI }) =>
      deps.namespace &&
      ComposeAPI.moduleRead({ namespaceID: deps.namespace, moduleID: value }).then(m => ({
        ...m,
        resourceType: 'compose:module',
      })),
  },
  page: {
    widget: 'text',
    idInput: true,
    deps: ['namespace'],
    descriptionKey: 'editor.required-namespace',
    resolve: ({ value }, deps, { ComposeAPI }) =>
      deps.namespace &&
      ComposeAPI.pageRead({ namespaceID: deps.namespace, pageID: value }).then(p => ({ ...p })),
  },
  record: {
    widget: 'record',
    deps: ['namespace', 'module'],
    resolve: ({ value }, deps, { ComposeAPI }) =>
      deps.namespace &&
      deps.module &&
      ComposeAPI.recordRead({
        namespaceID: deps.namespace,
        moduleID: deps.module,
        recordID: value,
      }).then(r => ({ ...r, resourceType: 'compose:record' })),
  },
  user: {
    widget: 'user',
    resolve: ({ value }, deps, { SystemAPI }) =>
      SystemAPI.userRead({ userID: value }).then(u => ({ ...u, resourceType: 'User' })),
  },
  role: {
    widget: 'role',
    resolve: ({ value }, deps, { SystemAPI }) =>
      SystemAPI.roleRead({ roleID: value }).then(r => ({ ...r, resourceType: 'Role' })),
  },
  application: {
    widget: 'text',
    idInput: true,
    resolve: ({ value }, deps, { SystemAPI }) =>
      SystemAPI.applicationRead({ applicationID: value }).then(a => ({ ...a })),
  },
}

// old* properties (oldRecord, oldNamespace, ...) share the definition of
// their base property; their deps still point at the current (non-old)
// namespace/module fields, matching how the scope is resolved server-side.
function baseName(name) {
  return name.replace(/^old(.)/, (_, c) => c.toLowerCase())
}

function scopeDef(name) {
  return SCOPE_DEFS[name] || SCOPE_DEFS[baseName(name)]
}

// Event types list properties leaf-first (record, module, namespace); the
// form wants dependencies first (namespace, module, record), old* variants
// right after their base property
const SCOPE_ORDER = Object.keys(SCOPE_DEFS)

function scopeOrder(name) {
  const base = SCOPE_ORDER.indexOf(baseName(name))
  if (base < 0) return SCOPE_ORDER.length * 2
  return base * 2 + (name === baseName(name) ? 0 : 1)
}

// Build form field descriptors for the trigger's scope properties.
// `previous` (name → field) keeps already-entered values when re-opening.
export function buildScopeFields(properties = [], t, previous = {}) {
  return [...properties]
    .sort((a, b) => scopeOrder(a.name) - scopeOrder(b.name))
    .map(({ name }) => {
      const def = scopeDef(name) || {}

      return {
        name,
        section: 'scope',
        label: `${name}${def.idInput ? t('editor.id-parenthesis') : ''}`,
        description: def.descriptionKey ? t(def.descriptionKey) : '',
        widget: def.widget,
        deps: def.deps,
        value: (previous[name] || {}).value,
      }
    })
}

// Named-input types (workflow.meta.input) → form widget
const INPUT_WIDGETS = {
  Boolean: 'boolean',
  DateTime: 'datetime',
  Number: 'number',
}

// Build form field descriptors for the workflow's declared named inputs.
// `previous` (name → field) keeps already-entered values when re-opening,
// unless the input's declared type changed in the meantime.
export function buildInputFields(defs = [], previous = {}) {
  return defs.map(def => {
    const type = (def.types || [])[0] || 'String'
    const prev = (previous[def.name] || {}).type === type ? previous[def.name] : {}
    return {
      name: def.name,
      section: 'input',
      label: def.label || def.name,
      type,
      widget: INPUT_WIDGETS[type] || 'text',
      required: !!def.required,
      value: prev.value ?? (type === 'Boolean' ? false : undefined),
    }
  })
}

// Cast a filled-in named input to a typed expr value based on its declared type
function typedFromInput({ type, value }) {
  // Boolean always resolves — false is a valid value to inject
  if (type === 'Boolean') {
    return { '@type': 'Boolean', '@value': value === true || value === 'true' }
  }

  if (value === '' || value === null || value === undefined) return undefined

  switch (type) {
    case 'Number': {
      const n = Number(value)
      if (Number.isNaN(n)) return { '@type': 'String', '@value': value }
      return { '@type': 'Float', '@value': n }
    }
    case 'DateTime':
      return { '@type': 'DateTime', '@value': value }
    case 'Any':
      return { '@type': 'Any', '@value': value }
    // String, plus inputs declared before their type was retired
    default:
      return { '@type': 'String', '@value': String(value) }
  }
}

// Values of a field's dependency siblings, keyed by dependency name
function depValues(field, fields) {
  return (field.deps || []).reduce((deps, name) => {
    deps[name] = (fields.find(f => f.name === name) || {}).value
    return deps
  }, {})
}

// Resolve all filled-in fields into the Vars map sent to workflowExec.
// Scope fields fetch their full objects (unresolvable ones become empty
// variables), named inputs are cast to typed expr values and merged on top.
export async function encodeFields(fields, { ComposeAPI, SystemAPI }) {
  const args = {}

  for (const f of fields.filter(f => f.section === 'scope')) {
    const def = scopeDef(f.name)
    if (f.value && def) {
      // Unmet dependencies resolve to undefined; a failed fetch of a provided
      // value propagates so the dialog can surface it
      args[f.name] = await def.resolve(f, depValues(f, fields), { ComposeAPI, SystemAPI })
    }
    if (!args[f.name]) args[f.name] = {}
  }

  const vars = automation.Encode(args)

  for (const f of fields.filter(f => f.section === 'input')) {
    const typed = typedFromInput(f)
    if (typed !== undefined) vars[f.name] = typed
  }

  return vars
}
