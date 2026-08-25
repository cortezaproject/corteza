import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia } from 'pinia'
import { compose } from '@planetcrust/human-js'

// A translator hands the resource back to the editor that owns it. Handing back
// a plain object instead of the typed resource crashed both editors on save:
// the chart editor binds `report.yAxis.formatting`, which a raw read has no
// `yAxis` for, and the module editor calls `module.systemFields()`.

vi.mock('vue-i18n', () => ({
  useI18n: () => ({ t: (k: string) => k, locale: { value: 'en' } }),
}))
vi.mock('@planetcrust/human-vue', () => ({
  useRBACStore: () => ({ can: () => true }),
}))

import ChartTranslator from './Chart/ChartTranslator.vue'
import ModuleTranslator from './Module/ModuleTranslator.vue'

// The classes validate ids, so these have to look like real snowflakes.
const NS = '510738942305894401'
const namespace = { namespaceID: NS }

// The wrapper builds the config and hands it to the button; the dialog is what
// calls it, so the test takes it off the button the same way.
const updaterOf = (wrapper: any) =>
  wrapper.findComponent({ name: 'CTranslatorButton' }).props('updater')

function api(overrides: Record<string, any> = {}) {
  return {
    chartUpdateTranslations: vi.fn().mockResolvedValue({}),
    moduleUpdateTranslations: vi.fn().mockResolvedValue({}),
    moduleListTranslations: vi.fn().mockResolvedValue([]),
    ...overrides,
  }
}

const mountWrapper = (component: any, props: any, $ComposeAPI: any) =>
  mount(component, {
    props,
    global: {
      // CTranslatorButton reaches the translator store.
      plugins: [createPinia()],
      provide: { $ComposeAPI, $Settings: { get: () => ['en', 'fr'] } },
      stubs: { Button: true },
      directives: { tooltip: {} },
      mocks: { $t: (k: string) => k },
    },
  })

describe('ChartTranslator', () => {
  // What chartRead actually returns: no yAxis on the report at all.
  const raw = {
    chartID: '510738942338269185',
    namespaceID: NS,
    name: 'Pipeline',
    config: { reports: [{ metrics: [{ metricID: '11', field: 'amount', type: 'bar' }] }] },
  }

  it('hands back a constructed chart, not the raw read', async () => {
    const $ComposeAPI = api({ chartRead: vi.fn().mockResolvedValue(raw) })
    const wrapper = mountWrapper(
      ChartTranslator,
      { chart: new compose.Chart(raw), namespace },
      $ComposeAPI,
    )

    await updaterOf(wrapper)([])

    const [[emitted]] = wrapper.emitted('update:chart') as any[]
    expect(emitted).toBeInstanceOf(compose.Chart)
    expect(emitted.config.reports[0].yAxis).toBeTypeOf('object')
  })
})

describe('ModuleTranslator', () => {
  const raw = {
    moduleID: '510738942306287617',
    namespaceID: NS,
    name: 'Probe',
    handle: 'probe',
    fields: [],
  }

  it('hands back a constructed module, not the plain clone', async () => {
    const $ComposeAPI = api()
    const wrapper = mountWrapper(
      ModuleTranslator,
      { module: new compose.Module(raw), namespace },
      $ComposeAPI,
    )

    await updaterOf(wrapper)([])

    const [[emitted]] = wrapper.emitted('update:module') as any[]
    expect(emitted).toBeInstanceOf(compose.Module)
    expect(emitted.systemFields).toBeTypeOf('function')
  })
})
