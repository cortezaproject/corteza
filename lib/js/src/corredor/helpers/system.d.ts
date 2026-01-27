import { ListResponse, PermissionRole, PermissionResource } from './shared'
import { System as SystemAPI } from '../../api-clients'
import { User, Role, Application } from '../../system/'
interface SystemContext {
    SystemAPI: SystemAPI;
    $user?: User;
    $role?: Role;
    $application?: Application;
}
interface UserListFilter {
    [key: string]: string | boolean | number | {
        [key: string]: string;
    } | undefined;
    userID?: string;
    roleID?: string;
    query?: string;
    username?: string;
    email?: string;
    handle?: string;
    kind?: string;
    incDeleted?: boolean;
    incSuspended?: boolean;
    deleted?: boolean;
    suspended?: boolean;
    labels?: {
        [key: string]: string;
    };
    limit?: number;
    pageCursor?: string;
    sort?: string;
}
interface RoleListFilter {
    [key: string]: string | boolean | number | {
        [key: string]: string;
    } | undefined;
    query?: string;
    deleted?: boolean;
    archived?: boolean;
    labels?: {
        [key: string]: string;
    };
    limit?: number;
    pageCursor?: string;
    sort?: string;
}
/**
 * SystemHelper provides layer over System API and utilities that simplify automation script writing
 */
export default class SystemHelper {
  readonly SystemAPI: SystemAPI
  readonly $user?: User
  readonly $role?: Role
  readonly $application?: Application
  constructor(ctx: SystemContext);
  /**
     * Searches for users
     *
     * @example
     * System.findUsers('some-joe').then(({ set }) => {
     *   // do something with users (User[]) in set
     * })
     *
     * @param filter - filter object (or filtering conditions when string)
     * @property filter.query - Find %query% in email, handle, username, name...
     * @property filter.username - Filter by username
     * @property filter.handle - Filter by handle
     * @property filter.email - Filter by email
     * @property filter.kind - Filter by kind ('normal' - default, 'bot')
     * @property filter.incDeleted - Include deleted users
     * @property filter.incSuspended - Include suspended users
     * @property filter.sort - Sort results
     * @property filter.perPage - max returned records per page
     * @property filter.page - page to return (1-based)
     */
  findUsers(filter?: string | UserListFilter): Promise<ListResponse<User[], UserListFilter>>;
  /**
     * Finds user by ID
     *
     * @example
     * System.findUserByID()
     *
     * @param user
     */
  findUserByID(user: string | User): Promise<User>;
  /**
     * Finds user by email
     *
     * @example
     * System.findUserByEmail('name@example.tld').then(user => {
     *   // do something with user
     * })
     *
     * @param email
     */
  findUserByEmail(email: string): Promise<User>;
  /**
     * Finds user by handle
     *
     * @example
     * System.findUserByHandle('some-handle').then(user => {
     *   // do something with user
     * })
     *
     * @param handle
     */
  findUserByHandle(handle: string): Promise<User>;
  /**
     * Updates or creates user
     *
     * @example
     * System.findUserByHandle('some-handle').then(user => {
     *   user.handle = 'better-handle'
     *   return System.saveUser(user)
     * })
     *
     * @param user
     */
  saveUser(user: User): Promise<User>;
  /**
     * Sets/updates password for the user
     *
     * @example
     * System.findUserByHandle('some-handle').then(user => {
     *   user.handle = 'better-handle'
     *   return System.saveUser(user)
     * })
     *
     * @param password
     * @param user
     */
  setPassword(password: string, user?: User | undefined): Promise<User>;
  /**
     * Deletes user
     *
     * @example
     * System.findUserByHandle('soon-to-be-deleted').then(user => {
     *   return System.deleteUser(user)
     * })
     *
     * @param user
     */
  deleteUser(user: string | User): Promise<unknown>;
  /**
     * Searches for roles
     *
     * @param filter
     */
  findRoles(filter?: string | RoleListFilter): Promise<ListResponse<Role[], RoleListFilter>>;
  /**
     * Finds user by ID
     *
     * @param role
     */
  findRoleByID(role: string | Role): Promise<Role>;
  /**
     * Finds role by handle
     *
     * @example
     * System.findRoleByHandle('some-handle').then(user => {
     *   // do something with role
     * })
     *
     * @param handle
     */
  findRoleByHandle(handle: string): Promise<Role>;
  /**
     *
     * @param role
     */
  saveRole(role: Role): Promise<Role>;
  /**
     * Deletes a role
     *
     * @example
     * System.findUserByHandle('soon-to-be-deleted').then(user => {
     *   return System.deleteUser(user)
     * })
     *
     * @param role
     */
  deleteRole(role: Role): Promise<unknown>;
  /**
     * Assign role to user
     *
     * @example
     * addUserToRole('user-we-can-trust', 'admins')
     *
     * @param user resolvable user input
     * @param role resolvable role input
     */
  addUserToRole(user: User | string, role: Role | string): Promise<unknown>;
  /**
     * Remove role from user
     * @example
     * addUserToRole('user-we-can-trust', 'admins')
     *
     * @param user - resolvable user input
     * @param role - resolvable role input
     */
  removeUserFromRole(user: User | string, role: Role | string): Promise<unknown>;
  /**
     * Resolves users from the arguments and returns first valid
     *
     * Knows how to resolve from:
     *  - string that looks like an ID - find by id (fallback to find-by-handle)
     *  - string that looks like an email - find by email (fallback to find-by-handle)
     *  - string - find by handle
     *  - User object
     *  - object with userID or ownerID properties
     */
  resolveUser(...args: unknown[]): Promise<User>;
  /**
     * Resolves users from the arguments and returns first valid
     *
     * Knows how to resolve from:
     *  - string that looks like an ID - find by id (fallback to find-by-handle)
     *  - string - find by handle
     *  - Role object
     *  - object with roleID property
     */
  resolveRole(...args: unknown[]): Promise<Role>;
  /**
     * Allows access for the given role for the given System resource
     *
     * @example
     * // Allows users with `someRole` to access the newly created user
     * await Compose.allow({
     *    role: someRole,
     *    resource: newUser,
     *    operation: 'read',
     * })
     */
  allow(...pr: {
        role: PermissionRole;
        resource: PermissionResource;
        operation: string;
    }[]): Promise<void>;
    /**
     * Denies access for the given role for the given System resource
     *
     * @example
     * // Denies users with `someRole` from accessing the newly created user
     * await Compose.deny({
     *    role: someRole,
     *    resource: newUser,
     *    operation: 'read',
     * })
     */
  deny(...pr: {
        role: PermissionRole;
        resource: PermissionResource;
        operation: string;
    }[]): Promise<void>;
    /**
     * Inherits access for the given role for the given System resource
     *
     * @example
     * // Uses inherited permissions for the `sameRole` for the newly created user
     * await Compose.inherit({
     *    role: someRole,
     *    resource: newUser,
     *    operation: 'read',
     * })
     */
  inherit(...pr: {
        role: PermissionRole;
        resource: PermissionResource;
        operation: string;
    }[]): Promise<void>;
}
export {}
