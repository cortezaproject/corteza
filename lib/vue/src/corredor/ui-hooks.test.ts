import { describe, it, expect } from 'vitest'
import { UIHooks, type Script } from './ui-hooks'

function script(name: string, uiProps: Array<{ name: string; value: string }>): Script {
  return {
    name,
    label: `label of ${name}`,
    description: `description of ${name}`,
    triggers: [
      {
        eventTypes: ['onManual'],
        resourceTypes: ['system:user'],
        uiProps,
      },
    ],
  }
}

const app = { name: 'app', value: 'admin' }
const page = { name: 'page', value: 'user/editor' }
const slot = { name: 'slot', value: 'infoFooter' }

describe('UIHooks.Register', () => {
  it('reads label, page, slot, variant and app off the trigger', () => {
    const hooks = new UIHooks({ apps: ['admin'] })
    hooks.Register(script('s1', [app, page, slot, { name: 'variant', value: 'danger' }]))

    const [button] = hooks.Find('system:user', 'user/editor', 'infoFooter', 'admin')

    expect(button.script).toBe('s1')
    expect(button.label).toBe('label of s1')
    expect(button.description).toBe('description of s1')
    expect(button.page).toBe('user/editor')
    expect(button.slot).toBe('infoFooter')
    expect(button.variant).toBe('danger')
    expect(button.app).toBe('admin')
  })

  it('prefers the trigger label over the script label', () => {
    const hooks = new UIHooks({ apps: ['admin'] })
    hooks.Register(script('s1', [app, { name: 'label', value: 'From the trigger' }]))

    expect(hooks.FindByScript('s1')?.label).toBe('From the trigger')
  })

  it('skips a trigger naming an app it does not serve and keeps one naming none', () => {
    const hooks = new UIHooks({ apps: ['admin'] })
    hooks.Register(
      script('compose-only', [{ name: 'app', value: 'compose' }]),
      script('no-app', []),
    )

    expect(hooks.FindByScript('compose-only')).toBeUndefined()
    expect(hooks.FindByScript('no-app')?.app).toBe('')
  })

  it('registers a script once however often it is registered', () => {
    const hooks = new UIHooks({ apps: ['admin'] })
    hooks.Register(script('s1', [app]))
    hooks.Register(script('s1', [app]))

    expect(hooks.Find('system:user')).toHaveLength(1)
  })

  it('takes no trigger from a script with errors', () => {
    const hooks = new UIHooks({ apps: ['admin'] })
    hooks.Register({ ...script('broken', [app]), errors: ['does not compile'] })

    expect(hooks.Find('system:user')).toHaveLength(0)
  })

  it('takes no trigger that is not manual', () => {
    const hooks = new UIHooks({ apps: ['admin'] })
    hooks.Register({
      name: 'implicit',
      triggers: [{ eventTypes: ['beforeCreate'], resourceTypes: ['system:user'], uiProps: [app] }],
    })

    expect(hooks.Find('system:user')).toHaveLength(0)
  })

  it('orders buttons by trigger weight', () => {
    const hooks = new UIHooks({ apps: ['admin'] })
    hooks.Register(
      {
        ...script('heavy', [app]),
        triggers: [{ ...script('heavy', [app]).triggers![0], weight: 9 }],
      },
      {
        ...script('light', [app]),
        triggers: [{ ...script('light', [app]).triggers![0], weight: 1 }],
      },
    )

    expect(hooks.Find('system:user').map(b => b.script)).toEqual(['light', 'heavy'])
  })
})

describe('UIHooks.Find', () => {
  const hooks = new UIHooks({ apps: ['admin', 'compose'] })
  hooks.Register(
    script('admin-editor', [app, page, slot]),
    script('compose-one', [{ name: 'app', value: 'compose' }]),
  )

  it('matches on resource type, page and slot', () => {
    expect(hooks.Find('system:user', 'user/editor', 'infoFooter', 'admin')).toHaveLength(1)
    expect(hooks.Find('system:user', 'role/editor', 'infoFooter', 'admin')).toHaveLength(0)
    expect(hooks.Find('system:user', 'user/editor', 'toolbar', 'admin')).toHaveLength(0)
    expect(hooks.Find('system:role', 'user/editor', 'infoFooter', 'admin')).toHaveLength(0)
  })

  it('takes a list of resource types', () => {
    expect(
      hooks.Find(['system:role', 'system:user'], 'user/editor', 'infoFooter', 'admin'),
    ).toHaveLength(1)
  })

  it('matches any app when none is given', () => {
    expect(hooks.Find('system:user').map(b => b.script)).toEqual(['compose-one'])
    expect(hooks.Find('system:user', undefined, undefined, 'compose').map(b => b.script)).toEqual([
      'compose-one',
    ])
    expect(hooks.Find('system:user', undefined, undefined, 'admin')).toHaveLength(0)
  })
})
