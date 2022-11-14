import { extractID, genericPermissionUpdater, isFresh, kv, ListResponse, PermissionRole, PermissionResource } from './shared'
import { System as SystemAPI } from '../../api-clients'
import { User, Role, Application } from '../../system/'
import { IsCortezaID } from '../../cast'

interface SystemContext {
  SystemAPI: SystemAPI;
  $user?: User;
  $role?: Role;
  $application?: Application;
}

interface UserListFilter {
  [key: string]: string|boolean|number|{[key: string]: string}|undefined;
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
  labels?: {[key: string]: string};
  limit?: number;
  pageCursor?: string;
  sort?: string;
}

interface RoleListFilter {
  [key: string]: string|boolean|number|{[key: string]: string}|undefined;
  query?: string;
  deleted?: boolean;
  archived?: boolean;
  labels?: {[key: string]: string};
  limit?: number;
  pageCursor?: string;
  sort?: string;
}

/**
 * Helpers to determine if specific object looks like the type we are interested in.
 * It does not rely on instanceof, because of bundling issues.
 */
function isUser (o: any) {
  return o && !!o.userID
}
function isRole (o: any) {
  return o && !!o.roleID
}

/**
 * SystemHelper provides layer over System API and utilities that simplify automation script writing
 */
export default class SystemHelper {
  readonly SystemAPI: SystemAPI;
  readonly $user?: User;
  readonly $role?: Role;
  readonly $application?: Application;

  constructor (ctx: SystemContext) {
    this.SystemAPI = ctx.SystemAPI

    this.$user = ctx.$user
    this.$role = ctx.$role
    this.$application = ctx.$application
  }

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
  async findUsers (filter?: string|UserListFilter): Promise<ListResponse<User[], UserListFilter>> {
    if (typeof filter === 'string') {
      filter = { query: filter }
    }

    return this.SystemAPI
      .userList(filter || {})
      .then(res => {
        res.set = (res.set as any[]).map(u => new User(u))
        return res as unknown as ListResponse<User[], UserListFilter>
      })
  }

  /**
   * Finds user by ID
   *
   * @example
   * System.findUserByID()
   *
   * @param user
   */
  async findUserByID (user: string|User): Promise<User> {
    const userID = extractID(user, 'userID')
    return this.SystemAPI.userRead({ userID }).then(u => new User(u))
  }

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
  async findUserByEmail (email: string): Promise<User> {
    return this.findUsers({ email }).then(res => {
      if (!Array.isArray(res.set) || res.set.length === 0) {
        throw new Error('user not found')
      }

      return new User(res.set[0])
    })
  }

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
  async findUserByHandle (handle: string): Promise<User> {
    return this.findUsers({ handle }).then(res => {
      if (!Array.isArray(res.set) || res.set.length === 0 || !res.set) {
        throw new Error('user not found')
      }

      return new User(res.set[0])
    })
  }

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
  async saveUser (user: User): Promise<User> {
    return Promise.resolve(user).then(user => {
      if (isFresh(user.userID)) {
        return this.SystemAPI.userCreate(kv(user)).then(user => new User(user))
      } else {
        return this.SystemAPI.userUpdate(kv(user)).then(user => new User(user))
      }
    })
  }

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
  async setPassword (password: string, user: User|undefined = this.$user): Promise<User> {
    return this.resolveUser(user).then(user => {
      const { userID } = user
      if (isFresh(userID)) {
        throw new Error('Cannot set password for non existing user')
      }

      return this.SystemAPI.userSetPassword({ password, userID }).then(u => new User(u))
    })
  }

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
  async deleteUser (user: string|User): Promise<unknown> {
    return Promise.resolve(user).then(user => {
      const userID = extractID(user, 'userID')

      if (!isFresh(userID)) {
        return this.SystemAPI.userDelete({ userID })
      }
    })
  }

  /**
   * Searches for roles
   *
   * @param filter
   */
  async findRoles (filter?: string|RoleListFilter): Promise<ListResponse<Role[], RoleListFilter>> {
    if (typeof filter === 'string') {
      filter = { query: filter }
    }

    return this.SystemAPI
      .roleList(filter || {})
      .then(res => {
        res.set = (res.set as any[]).map(r => new Role(r))
        return res as unknown as ListResponse<Role[], RoleListFilter>
      })
  }

  /**
   * Finds user by ID
   *
   * @param role
   */
  async findRoleByID (role: string|Role): Promise<Role> {
    const roleID = extractID(role, 'roleID')
    return this.SystemAPI.roleRead({ roleID }).then(r => new Role(r))
  }

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
  async findRoleByHandle (handle: string): Promise<Role> {
    return this.findRoles(handle).then(res => {
      if (!Array.isArray(res.set) || res.set.length === 0 || !res.set) {
        throw new Error('role not found')
      }

      return new Role(res.set[0])
    })
  }

  /**
   *
   * @param role
   */
  async saveRole (role: Role): Promise<Role> {
    return Promise.resolve(role).then(role => {
      if (isFresh(role.roleID)) {
        return this.SystemAPI.roleCreate(kv(role)).then(role => new Role(role))
      } else {
        return this.SystemAPI.roleUpdate(kv(role)).then(role => new Role(role))
      }
    })
  }

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
  async deleteRole (role: Role): Promise<unknown> {
    return Promise.resolve(role).then(role => {
      const roleID = extractID(role, 'roleID')

      if (!isFresh(roleID)) {
        return this.SystemAPI.roleDelete({ roleID })
      }
    })
  }

  /**
   * Assign role to user
   *
   * @example
   * addUserToRole('user-we-can-trust', 'admins')
   *
   * @param user resolvable user input
   * @param role resolvable role input
   */
  async addUserToRole (user: User|string, role: Role|string): Promise<unknown> {
    let userID: string
    let roleID: string

    return this.resolveUser(user, this.$user).then(user => {
      userID = extractID(user, 'userID')
      return this.resolveRole(role, this.$role)
    }).then(role => {
      roleID = extractID(role, 'roleID')
      return this.SystemAPI.roleMemberAdd({ roleID, userID })
    })
  }

  /**
   * Remove role from user
   * @example
   * addUserToRole('user-we-can-trust', 'admins')
   *
   * @param user - resolvable user input
   * @param role - resolvable role input
   */
  async removeUserFromRole (user: User|string, role: Role|string): Promise<unknown> {
    let userID: string
    let roleID: string

    return this.resolveUser(user, this.$user).then(user => {
      userID = extractID(user, 'userID')
      return this.resolveRole(role, this.$role)
    }).then(role => {
      roleID = extractID(role, 'roleID')
      return this.SystemAPI.roleMemberRemove({ roleID, userID })
    })
  }

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
  async resolveUser (...args: unknown[]): Promise<User> {
    for (let u of args) {
      // Resolve pending promises if any...
      u = await u

      if (!u) {
        continue
      }

      if (typeof u === 'string') {
        try {
          if (IsCortezaID(u)) {
            // Looks like an ID, try to find it and fall back to handle
            return await this.findUserByID(u)
          } else if (u.indexOf('@') > 0) {
            return await this.findUserByEmail(u)
          }
        } catch (e) {}

        // Always fall back to handle
        return this.findUserByHandle(u)
      }

      if (typeof u !== 'object') {
        continue
      }

      if (isUser(u)) {
        // Already got what we need
        return Promise.resolve(u as User)
      }

      // Other kind of object with properties that might hold user ID
      const {
        userID,
        ownerID,
      } = u as { userID?: string; ownerID?: string}
      return this.resolveUser(userID, ownerID)
    }

    return Promise.reject(new Error('unexpected input type for user resolver'))
  }

  /**
   * Resolves users from the arguments and returns first valid
   *
   * Knows how to resolve from:
   *  - string that looks like an ID - find by id (fallback to find-by-handle)
   *  - string - find by handle
   *  - Role object
   *  - object with roleID property
   */
  async resolveRole (...args: unknown[]): Promise<Role> {
    for (let r of args) {
      // Resolve pending promises if any...
      r = await r

      if (!r) {
        continue
      }

      if (typeof r === 'string') {
        if (IsCortezaID(r)) {
          // Looks like an ID, try to find it and fall back to handle
          return this.findRoleByID(r).catch(() => this.findRoleByHandle(r as string))
        }

        return this.findRoleByHandle(r)
      }

      if (typeof r !== 'object') {
        continue
      }

      if (isRole(r)) {
        // Already got what we need
        return r as Role
      }

      // Other kind of object with properties that might hold role ID
      const {
        roleID,
      } = r as { roleID?: string}
      return this.resolveRole(roleID)
    }

    return Promise.reject(Error('unexpected input type for role resolver'))
  }

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
  async allow (...pr: { role: PermissionRole; resource: PermissionResource; operation: string }[]) {
    const rr = pr.map(p => ({
      role: p.role,
      resource: p.resource,
      operation: p.operation,
      access: 'allow',
    }))
    return genericPermissionUpdater(this.SystemAPI, rr)
  }

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
  async deny (...pr: { role: PermissionRole; resource: PermissionResource; operation: string }[]) {
    const rr = pr.map(p => ({
      role: p.role,
      resource: p.resource,
      operation: p.operation,
      access: 'deny',
    }))
    return genericPermissionUpdater(this.SystemAPI, rr)
  }

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
  async inherit (...pr: { role: PermissionRole; resource: PermissionResource; operation: string }[]) {
    const rr = pr.map(p => ({
      role: p.role,
      resource: p.resource,
      operation: p.operation,
      access: 'inherit',
    }))
    return genericPermissionUpdater(this.SystemAPI, rr)
  }
}
