import { describe, it } from 'mocha'
import { expect } from 'chai'
import { Ctx } from './ctx'
import { BaseArgs } from './shared'

// A script reaches an automation through the context, with the token of
// whoever triggered it
describe('ctx automation api', () => {
  const args = { $invoker: {}, authToken: 'a-token' } as unknown as BaseArgs

  it('builds the client from the configured base URL', () => {
    const ctx = new Ctx(args, console as never, {
      config: { cServers: { automation: { apiBaseURL: 'http://api/automation' } } },
    })

    expect(ctx.AutomationAPI).to.be.an('object')
    expect(ctx.AutomationAPI.baseURL).to.equal('http://api/automation')
  })

  it('says which configuration is missing rather than failing later', () => {
    const ctx = new Ctx(args, console as never, { config: { cServers: {} } })

    expect(() => ctx.AutomationAPI).to.throw(/automation server missing/)
  })
})
