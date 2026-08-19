import { describe, it, expect } from 'vitest'
import { computed } from 'vue'
import { Settings } from './settings'

function makeApi(payload: Record<string, any>) {
  return {
    baseURL: 'http://x/api/system',
    payload,
    attachmentOriginalEndpoint: ({ kind, attachmentID, name }: Record<string, string>) =>
      `/attachment/${kind}/${attachmentID}/original/${name}`,
    async settingsCurrent() {
      return this.payload
    },
  }
}

describe('Settings.get', () => {
  it('walks a dot-path and falls back to the default', async () => {
    const s = new Settings({ api: makeApi({ ui: { mainLogo: '/assets/logo.svg' } }) })
    await s.fetch()

    expect(s.get('ui.mainLogo')).toBe('/assets/logo.svg')
    expect(s.get('ui.iconLogo', 'fallback')).toBe('fallback')
  })
})

describe('Settings.attachment', () => {
  it('resolves an attachment: value to the original-file endpoint', async () => {
    const s = new Settings({ api: makeApi({ ui: { mainLogo: 'attachment:42' } }) })
    await s.fetch()

    expect(s.attachment('ui.mainLogo')).toBe(
      'http://x/api/system/attachment/settings/42/original/ui.mainLogo',
    )
  })

  it('serves a plain path off the instance root', async () => {
    const s = new Settings({ api: makeApi({ ui: { mainLogo: '/assets/logo.svg' } }) })
    await s.fetch()

    expect(s.attachment('ui.mainLogo')).toBe('http://x/assets/logo.svg')
  })
})

// A settings change reaches the app shell only if the effects that read the
// previous value re-run: components hold no copy, they read through $Settings.
describe('Settings.fetch', () => {
  it('invalidates effects that read a setting it replaces', async () => {
    const api = makeApi({ ui: { mainLogo: '/assets/logo.svg' } })
    const s = new Settings({ api })
    await s.fetch()

    const logo = computed(() => s.attachment('ui.mainLogo'))
    expect(logo.value).toBe('http://x/assets/logo.svg')

    api.payload = { ui: { mainLogo: 'attachment:42' } }
    await s.fetch()

    expect(logo.value).toBe('http://x/api/system/attachment/settings/42/original/ui.mainLogo')
  })

  it('drops keys the new payload no longer carries', async () => {
    const api = makeApi({ ui: { mainLogo: 'attachment:42' } })
    const s = new Settings({ api })
    await s.fetch()

    const logo = computed(() => s.get('ui.mainLogo', 'default'))
    expect(logo.value).toBe('attachment:42')

    api.payload = { ui: {} }
    await s.fetch()

    expect(logo.value).toBe('default')
  })
})
