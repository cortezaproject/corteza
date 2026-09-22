import { expect } from 'chai'
import { RoleEvent, SystemEvent, TriggerSystemServerScriptOnManual, UserEvent } from './events'
import { User } from './types/user'
import { Role } from './types/role'

const script = '/server-scripts/Greet.js:default'

function api() {
  const calls: Array<{ endpoint: string; params: object }> = []

  return {
    calls,
    automationTriggerScript(params: { script: string }): Promise<object> {
      calls.push({ endpoint: 'automation', params })
      return Promise.resolve({})
    },
    userTriggerScript(params: { userID: string; script: string }): Promise<object> {
      calls.push({ endpoint: 'user', params })
      return Promise.resolve({ userID: params.userID })
    },
    roleTriggerScript(params: { roleID: string; script: string }): Promise<object> {
      calls.push({ endpoint: 'role', params })
      return Promise.resolve({ roleID: params.roleID })
    },
  }
}

describe('system events', () => {
  describe('TriggerSystemServerScriptOnManual', () => {
    it('routes a user event to the user endpoint', async () => {
      const a = api()
      const user = new User({ userID: '1', email: 'someone@example.tld' })

      const rval = await TriggerSystemServerScriptOnManual(a)(UserEvent(user), script)

      expect(a.calls[0].endpoint).to.equal('user')
      expect(a.calls[0].params).to.include({ userID: '1', script })
      expect(rval).to.be.instanceOf(User)
    })

    it('routes a role event to the role endpoint', async () => {
      const a = api()
      const role = new Role({ roleID: '2', handle: 'support' })

      const rval = await TriggerSystemServerScriptOnManual(a)(RoleEvent(role), script)

      expect(a.calls[0].endpoint).to.equal('role')
      expect(a.calls[0].params).to.include({ roleID: '2', script })
      expect(rval).to.be.instanceOf(Role)
    })

    it('routes a system event to the service endpoint', async () => {
      const a = api()

      await TriggerSystemServerScriptOnManual(a)(SystemEvent(), script)

      expect(a.calls[0].endpoint).to.equal('automation')
    })

    it('refuses a resource it cannot trigger a script on', async () => {
      const a = api()
      const ev = { ...SystemEvent(), resourceType: 'system:application' }

      try {
        await TriggerSystemServerScriptOnManual(a)(ev, script)
        expect.fail('expected the handler to refuse the event')
      } catch (e) {
        expect((e as Error).message).to.contain('unknown resource type')
      }
      expect(a.calls).to.have.length(0)
    })
  })
})
