import { CortezaID, NoID } from '../../cast'

interface KV {
  [_: string]: unknown;
}

interface PermissionUpdater {
  permissionsUpdate({ roleID, rules }: { roleID: string; rules: Array<object>}): void;
}

export interface PermissionResource {
  resourceID: string;
  [_: string]: any;
}

export interface PermissionRole {
  roleID: string;
  [_: string]: any;
}

export interface PermissionRule {
  role: PermissionRole;
  resource: PermissionResource;
  operation: string;
  access: string;
}

export interface Permissions {
  [key: string]: {
    resource: string;
    operation: string;
    access: string;
  }[];
}

export function kv (a: unknown): KV { return a as KV }

export interface ListResponse<S, F> {
  set: S;
  filter: F;
}

/**
 * Extracts ID-like (numeric) value from string or object
 *
 * @param value - that stores ID in some way
 * @param prop - possible key lookup
 */
export function extractID (value?: unknown, prop?: string): string {
  if (value && typeof value === 'object') {
    if (!prop || !Object.prototype.hasOwnProperty.call(value, prop)) {
      return NoID
    }

    value = (value as {[_: string]: unknown})[prop]
  }

  return CortezaID(value)
}

export function isFresh (ID: string): boolean {
  return !ID || ID === NoID
}

export function genericPermissionUpdater (API: PermissionUpdater, rules: PermissionRule[]): void {
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
