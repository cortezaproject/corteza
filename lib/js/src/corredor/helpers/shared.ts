import { HumanID, IsHumanID, NoID } from '../../cast'

export interface KV {
  [_: string]: unknown
}

/**
 * What a generated helper method takes: the endpoint's arguments.
 *
 * Typed as object rather than KV so that a resource instance — a User, a Module —
 * is accepted where the endpoint wants its ID.
 */
export type Args = object

/**
 * Arguments of an endpoint with a single ID in its path: the arguments, or the
 * ID, handle or object that ID names.
 */
export type IdArgs = Args | string

interface PermissionUpdater {
  permissionsUpdate({ roleID, rules }: { roleID: string; rules: Array<object> }): void
}

export interface PermissionResource {
  resourceID: string
  [_: string]: any
}

export interface PermissionRole {
  roleID: string
  [_: string]: any
}

export interface PermissionRule {
  role: PermissionRole
  resource: PermissionResource
  operation: string
  access: string
}

export interface Permissions {
  [key: string]: {
    resource: string
    operation: string
    access: string
  }[]
}

export function kv(a: unknown): KV {
  return a as KV
}

export interface ListResponse<S, F> {
  set: S
  filter: F
}

/**
 * Extracts ID-like (numeric) value from string or object
 *
 * @param value - that stores ID in some way
 * @param prop - possible key lookup
 */
export function extractID(value?: unknown, prop?: string): string {
  if (value && typeof value === 'object') {
    if (!prop || !Object.prototype.hasOwnProperty.call(value, prop)) {
      return NoID
    }

    value = (value as { [_: string]: unknown })[prop]
  }

  return HumanID(value)
}

export function isFresh(ID: string): boolean {
  return !ID || ID === NoID
}

export function genericPermissionUpdater(API: PermissionUpdater, rules: PermissionRule[]): void {
  const g: Permissions = rules.reduce((acc: Permissions, p: PermissionRule) => {
    if (!acc[p.role.roleID]) {
      acc[p.role.roleID] = []
    }

    acc[p.role.roleID].push({
      resource: p.resource.resourceID,
      operation: p.operation,
      access: p.access,
    })
    return acc
  }, {})

  // @todo should return promise and stack all these into Promise.all()
  Object.keys(g).forEach(async roleID => {
    // permissions grouped per role
    await API.permissionsUpdate({ roleID, rules: g[roleID] })
  })
}

/**
 * How one entity is found when a handle arrives where its ID belongs.
 *
 * @property idField - name the entity's ID goes by
 * @property list - API client method that lists the entity
 * @property by - list filters a handle-shaped value is looked up by, in order
 * @property parents - path parameters the list call needs from its context
 */
export interface RefTarget {
  idField: string
  list?: string
  by: string[]
  parents: string[]
}

export interface RefSpec {
  [entity: string]: RefTarget
}

type Lister = (a: KV) => Promise<KV>

/**
 * Normalises a generated method's argument into the endpoint's argument object.
 *
 * A bare string stands for the single ID in the endpoint's path.
 */
export function argsOf(a: unknown, idField?: string): KV {
  if (typeof a === 'string') {
    return idField ? { [idField]: a } : {}
  }

  return { ...(a as KV) }
}

/**
 * Puts a resource class around every member of a list response.
 */
export function castSet<T>(res: KV, cast: (v: KV) => T): ListResponse<T[], KV> {
  const set = Array.isArray(res.set) ? (res.set as KV[]).map(cast) : []
  return { ...(res as object), set } as unknown as ListResponse<T[], KV>
}

/**
 * Turns an object, an ID or a handle into the ID an endpoint wants.
 *
 * An ID passes straight through. An object gives up the ID it carries, or its
 * handle when it carries no ID. Anything else is looked up through the entity's
 * own list endpoint, under each filter the entity is findable by.
 *
 * @param api - API client the lookup goes through
 * @param refs - what each entity is keyed and looked up by
 * @param entity - entity the value names
 * @param value - object, ID or handle
 * @param context - the endpoint's other arguments, for a list that needs a parent ID
 */
export async function resolveRef(
  api: object,
  refs: RefSpec,
  entity: string,
  value: unknown,
  context: KV = {},
): Promise<unknown> {
  const ref = refs[entity]
  if (!ref) {
    return value
  }

  const v = await value
  if (v === undefined || v === null || v === '') {
    return v
  }

  if (Array.isArray(v)) {
    return resolveRef(api, refs, entity, v[0], context)
  }

  if (typeof v === 'object') {
    const set = (v as KV).set
    if (Array.isArray(set)) {
      return resolveRef(api, refs, entity, set[0], context)
    }

    const ID = extractID(v, ref.idField)
    if (!isFresh(ID)) {
      return ID
    }

    const alt = (v as KV).handle || (v as KV).slug || (v as KV).name
    if (typeof alt !== 'string' || alt === '') {
      return ID
    }

    return resolveRef(api, refs, entity, alt, context)
  }

  if (typeof v !== 'string' || IsHumanID(v)) {
    return v
  }

  if (!ref.list || ref.by.length === 0) {
    return v
  }

  const list = (api as { [_: string]: Lister })[ref.list]
  const parents: KV = {}
  ref.parents.forEach(p => {
    parents[p] = context[p]
  })

  // An address can only be an email, whichever filters come before it.
  const by = v.includes('@') && ref.by.includes('email') ? ['email', ...ref.by] : ref.by

  for (const field of by) {
    const res = await list.call(api, { ...parents, [field]: v, limit: 1 })
    const set = res.set as KV[] | undefined
    if (Array.isArray(set) && set.length > 0) {
      return extractID(set[0], ref.idField)
    }
  }

  throw new Error(`${entity} not found`)
}
