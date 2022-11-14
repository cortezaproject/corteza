import { Event, GenericEventMaker, onManual } from '../eventbus/shared'
import { Role } from './types/role'
import { User } from './types/user'

interface TriggerEndpoints {
  automationTriggerScript (params: { script: string }): Promise<object>;
  roleTriggerScript (params: { roleID: string; script: string }): Promise<object>;
  userTriggerScript (params: { userID: string; script: string }): Promise<object>;
}

// @todo refactor this into more compose-like event structure (see compose/events.ts
export function SystemEvent (eventType = onManual): Event {
  return GenericEventMaker({ resourceType: 'system' }, eventType, () => true, {})
}

// @todo refactor this into more compose-like event structure (see compose/events.ts
export function UserEvent (user: User, eventType: string = onManual): Event {
  return GenericEventMaker(user, eventType, function (c) {
    switch (c.Name()) {
      case 'user':
      case 'user.handle':
        return c.Match(user.handle)
      case 'user.email':
        return c.Match(user.email)
    }

    return false
  }, { user })
}

// @todo refactor this into more compose-like event structure (see compose/events.ts
export function RoleEvent (role: Role, eventType = onManual): Event {
  return GenericEventMaker(role, eventType, function (c) {
    switch (c.Name()) {
      case 'role':
      case 'role.handle':
        return c.Match(role.handle)
      case 'role.name':
        return c.Match(role.name)
    }

    return false
  }, { role })
}

/**
 * Returns handler that routes onManual events for server script to the system API
 *
 * See makeAutomationScriptsRegistrator
 *
 * @param api
 * @return function
 */
export function TriggerSystemServerScriptOnManual (api: TriggerEndpoints) {
  return (ev: Event, script: string): Promise<unknown> => {
    const params = { script, args: ev.args }
    const { userID } = ev.args?.user as User
    const { roleID } = ev.args?.role as Role

    switch (ev.resourceType) {
      case 'system':
        return api.automationTriggerScript({ ...params })

      case 'system:user':
        return api.userTriggerScript({ userID, ...params }).then(rval => new User(rval))

      case 'system:role':
        return api.roleTriggerScript({ roleID, ...params }).then(rval => new Role(rval))

      default:
        throw Error(`cannot trigger server script: unknown resource type '${ev.resourceType}'`)
    }
  }
}
