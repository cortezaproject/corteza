import { describe, it } from 'mocha'
import { expect } from 'chai'
import { execFileSync } from 'child_process'
import fs from 'fs'
import path from 'path'

import GeneratedComposeHelper from './compose.gen'
import GeneratedSystemHelper from './system.gen'
import GeneratedAutomationHelper from './automation.gen'
import ComposeHelper from './compose'
import SystemHelper from './system'
import { Compose as ComposeAPI, System as SystemAPI } from '../../api-clients'
import { Module } from '../../compose'
import { Agent, UserGroup } from '../../system'

const tools = path.join(__dirname, '..', '..', '..', 'tools', 'codegen')
const codegen = path.join(tools, 'corredor-helpers.js')

interface Lock {
  names: { [service: string]: { [endpoint: string]: string } }
  curated: { [service: string]: { [endpoint: string]: string } }
}

const lock: Lock = JSON.parse(fs.readFileSync(path.join(tools, 'helper-names.lock.json'), 'utf8'))

// Method names a generated class declares, read from its source: an instance
// answers to inherited members too, and only the file says what was generated.
function generatedMethods(service: string): Set<string> {
  const src = fs.readFileSync(path.join(__dirname, `${service}.gen.ts`), 'utf8')
  const decl = /^ {2}(?:async )?([A-Za-z_$][\w$]*)\(/gm

  const out = new Set<string>()
  let m
  while ((m = decl.exec(src)) !== null) {
    if (m[1] !== 'constructor') {
      out.add(m[1])
    }
  }

  return out
}

// Members a curated helper states itself, taken off its own prototype so that
// what it inherits from the generated base is left out.
function curatedMethods(helper: object): Set<string> {
  const own = new Set(Object.getOwnPropertyNames(Object.getPrototypeOf(helper)))
  own.delete('constructor')
  return own
}

describe('generated helper layer', () => {
  describe('output', () => {
    it('matches its sources', () => {
      try {
        execFileSync(process.execPath, [codegen, '--check'], { encoding: 'utf8', stdio: 'pipe' })
      } catch (e) {
        const { stdout, stderr } = e as { stdout?: string; stderr?: string }
        expect.fail(`${stdout ?? ''}${stderr ?? ''}`.trim() || 'the helper layer check failed')
      }
    })
  })

  describe('curated methods', () => {
    const curated = [
      { service: 'compose', helper: new ComposeHelper({ ComposeAPI: {} as ComposeAPI }) },
      { service: 'system', helper: new SystemHelper({ SystemAPI: {} as SystemAPI }) },
    ]

    curated.forEach(({ service, helper }) => {
      const own = curatedMethods(helper)

      it(`are not shadowed by a generated ${service} method`, () => {
        const shadowed = [...generatedMethods(service)].filter(m => own.has(m))
        expect(shadowed, `generated ${service} methods that shadow a curated one`).to.deep.equal([])
      })

      it(`serve every ${service} endpoint the lock records as colliding`, () => {
        const unserved = Object.entries(lock.curated[service])
          .filter(([, owner]) => !own.has(owner.split('.')[1]))
          .map(([endpoint]) => endpoint)

        expect(unserved, `${service} endpoints whose curated method is gone`).to.deep.equal([])
      })
    })
  })

  describe('locked names', () => {
    const services = [
      { name: 'system', helper: GeneratedSystemHelper },
      { name: 'compose', helper: GeneratedComposeHelper },
      { name: 'automation', helper: GeneratedAutomationHelper },
    ]

    services.forEach(({ name, helper }) => {
      it(`are all still produced by ${name}`, () => {
        const has = Object.getOwnPropertyNames(helper.prototype)
        const missing = Object.entries(lock.names[name])
          .filter(([endpoint]) => !(endpoint in lock.curated[name]))
          .map(([, method]) => method)
          .filter(method => !has.includes(method))

        expect(missing, `${name} methods the lock names but the class no longer has`).to.deep.equal(
          [],
        )
      })
    })
  })

  describe('behaviour', () => {
    it('resolves a handle to an ID and casts the response', async () => {
      const calls: Array<{ method: string; args: Record<string, unknown> }> = []

      const api = {
        namespaceList(a: Record<string, unknown>) {
          calls.push({ method: 'namespaceList', args: a })
          return Promise.resolve({ set: [{ namespaceID: '100', slug: 'crm' }] })
        },
        moduleList(a: Record<string, unknown>) {
          calls.push({ method: 'moduleList', args: a })
          return Promise.resolve({ set: [{ moduleID: '200', handle: 'Contact' }] })
        },
        moduleUpdate(a: Record<string, unknown>) {
          calls.push({ method: 'moduleUpdate', args: a })
          return Promise.resolve({ moduleID: '200', namespaceID: '100', handle: 'Contact' })
        },
      } as unknown as ComposeAPI

      const helper = new GeneratedComposeHelper({ ComposeAPI: api })
      const module = await helper.updateModule({ namespaceID: 'crm', moduleID: 'Contact' })

      expect(calls.map(c => c.method)).to.deep.equal([
        'namespaceList',
        'moduleList',
        'moduleUpdate',
      ])
      expect(calls[0].args).to.include({ slug: 'crm' })
      // The module lookup is scoped by the namespace resolved just before it
      expect(calls[1].args).to.include({ namespaceID: '100', handle: 'Contact' })
      expect(calls[2].args).to.include({ namespaceID: '100', moduleID: '200' })

      expect(module).to.be.instanceOf(Module)
      expect(module.handle).to.equal('Contact')
    })

    it('leaves an ID alone', async () => {
      const calls: string[] = []

      const api = {
        agentRead(a: Record<string, unknown>) {
          calls.push(`agentRead ${a.agentID}`)
          return Promise.resolve({ agentID: '300', handle: 'triage' })
        },
      } as unknown as SystemAPI

      const helper = new GeneratedSystemHelper({ SystemAPI: api })
      const agent = await helper.findAgentByID('300')

      expect(calls).to.deep.equal(['agentRead 300'])
      expect(agent).to.be.instanceOf(Agent)
      expect(agent.handle).to.equal('triage')
    })

    it('takes the ID off an object standing in for it', async () => {
      const calls: string[] = []

      const api = {
        agentRead(a: Record<string, unknown>) {
          calls.push(`agentRead ${a.agentID}`)
          return Promise.resolve({ agentID: '300' })
        },
      } as unknown as SystemAPI

      const helper = new GeneratedSystemHelper({ SystemAPI: api })
      await helper.findAgentByID(new Agent({ agentID: '300' }))

      expect(calls).to.deep.equal(['agentRead 300'])
    })

    it('casts every member of a list response', async () => {
      const api = {
        userGroupList() {
          return Promise.resolve({
            set: [{ userGroupID: '400' }, { userGroupID: '401' }],
            filter: { limit: 2 },
          })
        },
      } as unknown as SystemAPI

      const helper = new GeneratedSystemHelper({ SystemAPI: api })
      const { set, filter } = await helper.findUserGroups()

      expect(set).to.have.length(2)
      set.forEach(g => expect(g).to.be.instanceOf(UserGroup))
      expect(set[1].userGroupID).to.equal('401')
      expect(filter).to.deep.equal({ limit: 2 })
    })

    it('says which resource a handle did not find', async () => {
      const api = {
        namespaceList() {
          return Promise.resolve({ set: [] })
        },
      } as unknown as ComposeAPI

      const helper = new GeneratedComposeHelper({ ComposeAPI: api })

      await helper
        .deleteNamespace('nope')
        .then(() => expect.fail('expected the lookup to fail'))
        .catch(e => expect((e as Error).message).to.equal('namespace not found'))
    })

    it('reaches the automation service through its own helper', async () => {
      const helper = new GeneratedAutomationHelper({
        AutomationAPI: {
          workflowList() {
            return Promise.resolve({ set: [{ workflowID: '500' }], filter: {} })
          },
        } as never,
      })

      const { set } = await helper.findWorkflows()
      expect(set[0].workflowID).to.equal('500')
    })
  })
})
