import { mount } from '@vue/test-utils'
import { createPinia, getActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import { createRouter, createMemoryHistory } from 'vue-router'
import type { Component, DefineComponent } from 'vue'
import type { MountingOptions } from '@vue/test-utils'

export interface TestContext {
  composeAPI?: Record<string, unknown> | null
  systemAPI?: Record<string, unknown> | null
  automationAPI?: Record<string, unknown> | null
  auth?: Record<string, unknown> | null
  settings?: Record<string, unknown> | null
  userStore?: Record<string, unknown> | null
  recordStore?: Record<string, unknown> | null
  pageStore?: Record<string, unknown> | null
  recordRoutes?: Record<string, unknown> | null
  namespace?: Record<string, unknown> | null
  additionalProvides?: Record<string, unknown>
}

const defaultAuth = {
  user: { userID: 'test-user-id', roles: [] },
  accessTokenFn: () => 'test-token',
}

const defaultSettings = {
  get: () => undefined,
  attachment: () => '',
}

export function mountWithContext<T extends Component | DefineComponent>(
  component: T,
  ctx: TestContext = {},
  mountOptions: MountingOptions<any> = {},
) {
  const pinia = getActivePinia() ?? createPinia()
  const router = createRouter({
    history: createMemoryHistory(),
    routes: [{ path: '/', component: { template: '<div/>' } }],
  })
  const i18n = createI18n({ legacy: false, locale: 'en', messages: { en: {} } })

  const provides: Record<string, unknown> = {
    '$ComposeAPI': ctx.composeAPI !== undefined ? ctx.composeAPI : null,
    '$SystemAPI': ctx.systemAPI !== undefined ? ctx.systemAPI : null,
    '$AutomationAPI': ctx.automationAPI !== undefined ? ctx.automationAPI : null,
    '$Auth': ctx.auth !== undefined ? ctx.auth : defaultAuth,
    '$auth': ctx.auth !== undefined ? ctx.auth : defaultAuth,
    '$Settings': ctx.settings !== undefined ? ctx.settings : defaultSettings,
    '$userStore': ctx.userStore !== undefined ? ctx.userStore : null,
    '$recordStore': ctx.recordStore !== undefined ? ctx.recordStore : null,
    '$pageStore': ctx.pageStore !== undefined ? ctx.pageStore : null,
    '$recordRoutes': ctx.recordRoutes !== undefined ? ctx.recordRoutes : null,
    '$namespace': ctx.namespace !== undefined ? ctx.namespace : null,
    '$toast': { add: () => {}, remove: () => {} },
    '$eventBus': { emit: () => {}, on: () => () => {}, off: () => {} },
    ...ctx.additionalProvides,
  }

  const globalConfig = mountOptions.global ?? {}

  return mount(component, {
    ...mountOptions,
    global: {
      ...globalConfig,
      plugins: [pinia, router, i18n, ...(globalConfig.plugins ?? [])],
      provide: { ...provides, ...(globalConfig.provide ?? {}) },
      stubs: {
        Teleport: true,
        ...(globalConfig.stubs ?? {}),
      },
    },
  })
}
