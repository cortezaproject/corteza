interface KV {
  [_: string]: unknown
}
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
export declare function kv(a: unknown): KV
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
export declare function extractID(value?: unknown, prop?: string): string
export declare function isFresh(ID: string): boolean
export declare function genericPermissionUpdater(
  API: PermissionUpdater,
  rules: PermissionRule[],
): void
export {}
