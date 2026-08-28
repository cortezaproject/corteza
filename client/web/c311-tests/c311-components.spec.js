import Vue from 'vue'
import fs from 'fs'
import path from 'path'
import { mount } from '@vue/test-utils'
import { components, mixins, c311I18n } from './c311-components'
import { formatC311DateTime as mockFormatC311DateTime } from './time-test-helper'
import Portal from '../compose/src/views/C311/Portal.vue'
import Staff from '../admin/src/views/C311/Staff.vue'
import composeRoutes from '../compose/src/views/routes'
import adminRoutes from '../admin/src/views/routes'

jest.mock('@cortezaproject/corteza-js', () => ({
  formatC311DateTime: value => mockFormatC311DateTime(value, 'en-US'),
}))
jest.mock('@cortezaproject/corteza-vue', () => ({
  components: {
    C311AppShell: {},
    C311DataState: {},
    C311ErrorSummary: {},
    C311CapabilityAction: {},
    C311HelpDrawer: {},
    C311LanguageSelector: {},
    C311MainNav: {},
    C311ResponsiveData: {},
  },
  mixins: {
    c311DirtyGuard: {
      data: () => ({ c311Dirty: false, c311DirtyStorageKey: '' }),
      methods: {
        c311MarkDirty (value = true) { this.c311Dirty = value },
        c311ReadDirtyDraft () { return null },
        c311SaveDirtyDraft () {},
        c311ClearDirtyDraft () { this.c311Dirty = false },
      },
    },
  },
  c311: {
    c311StateForError (error) {
      if (error?.status === 503 || error?.retryable) return 'retryable-error'
      return 'terminal-error'
    },
  },
}))

const {
  C311CapabilityAction,
  C311DataState,
  C311ErrorSummary,
  C311FocusModal,
  C311HelpDrawer,
  C311AccessPage,
  C311LanguageSelector,
  C311MainNav,
  C311ResponsiveData,
  C311StatusAnnouncer,
} = components

const ButtonStub = {
  name: 'b-button',
  inheritAttrs: false,
  props: ['disabled'],
  template: '<button v-bind="$attrs" :disabled="disabled" @click="$emit(\'click\', $event)"><slot /></button>',
}

const mocks = {
  $t: (key, fallback) => {
    const values = {
      'c311:action.retry': 'Try again',
      'c311:language.label': 'Language',
      'c311:language.english': 'English',
      'c311:language.spanish': 'Español',
      'c311:language.vietnamese': 'Tiếng Việt',
    }
    return values[key] || fallback || key
  },
  $i18n: { i18next: { changeLanguage: jest.fn() } },
}

const stubs = { 'b-button': ButtonStub }
const RouterLinkStub = {
  props: ['to'],
  template: '<a v-bind="$attrs" :href="to"><slot /></a>',
}
const AppShellStub = { template: '<main><slot name="nav" /><slot /></main>' }
const DataStateStub = { template: '<section><slot name="populated" /></section>' }
const ChildStub = { template: '<span><slot /></span>' }

describe('C311 shared components', () => {
  afterEach(() => {
    localStorage.clear()
    sessionStorage.clear()
    document.body.innerHTML = ''
  })

  it('resolves the current local C311 package artifacts', () => {
    const jsEntry = require.resolve('@cortezaproject/corteza-js')
    const vueEntry = require.resolve('@cortezaproject/corteza-vue')
    expect(jsEntry).toContain(`${path.sep}node_modules${path.sep}@cortezaproject${path.sep}corteza-js${path.sep}`)
    expect(vueEntry).toContain(`${path.sep}node_modules${path.sep}@cortezaproject${path.sep}corteza-vue${path.sep}`)
    expect(fs.readFileSync(jsEntry, 'utf8')).toContain('MockC311Provider')
    expect(fs.readFileSync(vueEntry, 'utf8')).toContain('C311AppShell')
  })

  it('exposes stable route attributes independent of translated labels', () => {
    const wrapper = mount(C311MainNav, {
      propsData: {
        items: [
          { route: '/c311/submit', label: 'Enviar una solicitud' },
          { route: '/c311/status', label: 'Consultar estado' },
        ],
      },
      mocks: {
        ...mocks,
        $route: { path: '/c311/submit' },
        $C311: { can: () => true, hasScope: () => true },
      },
      stubs: { 'router-link': RouterLinkStub },
    })
    expect(wrapper.find('[data-c311-route="/c311/status"]').exists()).toBe(true)
    expect(wrapper.find('[data-c311-route="/c311/status"]').text()).toBe('Consultar estado')
  })

  it('renders and announces all eight data states', async () => {
    const wrapper = mount(C311DataState, { propsData: { state: 'loading' }, mocks, stubs, slots: { populated: '<span data-populated>rows</span>' } })
    for (const state of ['loading', 'empty', 'populated', 'forbidden', 'not-found', 'validation-error', 'retryable-error', 'terminal-error']) {
      await wrapper.setProps({ state })
      expect(wrapper.attributes('data-state')).toBe(state)
      if (state === 'populated') expect(wrapper.find('[data-populated]').exists()).toBe(true)
      else expect(wrapper.find('[role="status"]').exists()).toBe(true)
    }
  })

  it('enforces capability and scope decisions in the action UI', () => {
    const runtime = { can: value => value === 'allowed', hasScope: value => value === 'scope.ok', session: { authenticated: true } }
    const allowed = mount(C311CapabilityAction, { propsData: { capability: 'allowed', scope: 'scope.ok' }, mocks: { ...mocks, $C311: runtime }, stubs })
    const denied = mount(C311CapabilityAction, { propsData: { capability: 'allowed', scope: 'scope.missing', explainWhenDenied: true }, mocks: { ...mocks, $C311: runtime }, stubs })
    expect(allowed.find('button').exists()).toBe(true)
    expect(denied.find('button').exists()).toBe(false)
    expect(denied.text()).toContain('This action is unavailable')
  })

  it('focuses the error summary and then the linked invalid field', async () => {
    const wrapper = mount({
      components: { C311ErrorSummary },
      template: '<div><c311-error-summary :errors="errors" /><input id="c311-summary" /></div>',
      data: () => ({ errors: [] }),
    }, { mocks, attachTo: document.body })
    await wrapper.setData({ errors: [{ field: '/summary', code: 'REQUIRED', message: 'Required' }] })
    await Vue.nextTick()
    expect(document.activeElement).toBe(wrapper.find('[data-c311-error-summary]').element)
    await wrapper.find('[data-c311-error-summary] a').trigger('click')
    expect(document.activeElement.id).toBe('c311-summary')
  })

  it('announces status changes through aria-live', () => {
    const wrapper = mount(C311StatusAnnouncer, { propsData: { message: 'Loaded', assertive: true } })
    expect(wrapper.attributes('aria-live')).toBe('assertive')
    expect(wrapper.text()).toContain('Loaded')
  })

  it('renders translated 401, 403 and 404 access pages with a focusable heading', async () => {
    const wrapper = mount(C311AccessPage, { propsData: { status: 401 }, mocks, stubs, attachTo: document.body })
    expect(wrapper.find('h1').text()).toContain('Sign-in required')
    expect(wrapper.find('h1').attributes('tabindex')).toBe('-1')
    await wrapper.setProps({ status: 403 })
    expect(wrapper.text()).toContain('Access denied')
    await wrapper.setProps({ status: 404 })
    expect(wrapper.text()).toContain('Not found')
  })

  it('traps modal focus and returns focus to its opener', async () => {
    const opener = document.createElement('button')
    opener.id = 'modal-opener'
    document.body.appendChild(opener)
    opener.focus()
    const wrapper = mount(C311FocusModal, { propsData: { value: true, title: 'Details' }, mocks, stubs, attachTo: document.body, slots: { default: '<button id="first">First</button>' } })
    await Vue.nextTick()
    expect(document.activeElement).toBe(wrapper.find('[data-c311-focus-modal] .btn-secondary').element)
    const close = wrapper.find('[data-c311-focus-modal] button:last-child')
    await close.trigger('keydown', { key: 'Tab' })
    expect(document.activeElement).toBe(wrapper.find('[data-c311-focus-modal] button:first-of-type').element)
    await wrapper.setProps({ value: false })
    await Vue.nextTick()
    expect(document.activeElement).toBe(opener)
  })

  it('opens and closes the help drawer with focus return', async () => {
    const wrapper = mount(C311HelpDrawer, { propsData: { helpKey: 'public.request.submit', content: 'Help text' }, mocks, stubs, attachTo: document.body })
    const trigger = wrapper.find('button')
    trigger.element.focus()
    await trigger.trigger('click')
    await Vue.nextTick()
    await Vue.nextTick()
    expect(wrapper.find('[data-c311-help-drawer]').exists()).toBe(true)
    expect(document.activeElement).toBe(wrapper.find('[data-c311-help-drawer] button').element)
    await wrapper.find('[data-c311-help-drawer]').trigger('keydown', { key: 'Escape' })
    await Vue.nextTick()
    expect(wrapper.find('[data-c311-help-drawer]').exists()).toBe(false)
    expect(document.activeElement).toBe(trigger.element)
    expect(wrapper.emitted('open')[0]).toEqual(['public.request.submit'])
    expect(wrapper.emitted('close')[0]).toEqual(['public.request.submit'])
  })

  it('switches locale, falls back to English, and persists by actor', async () => {
    const wrapper = mount(C311LanguageSelector, { propsData: { actorID: 'actor-fixture-001' }, mocks })
    await wrapper.setData({ selected: 'es' })
    await wrapper.find('select').trigger('change')
    expect(localStorage.getItem('c311.locale.actor-fixture-001')).toBe('es')
    expect(mocks.$i18n.i18next.changeLanguage).toHaveBeenCalledWith('es')
    expect(c311I18n.readC311Locale('actor-fixture-001')).toBe('es')
    expect(c311I18n.readC311Locale('actor-fixture-unknown')).toBe(null)
    expect(c311I18n.C311_MESSAGES?.en?.['status.loading']).toBe('Loading')
  })

  it('keeps drafts free of secrets and restores valid values after refresh', () => {
    const Guard = Vue.extend({
      mixins: [mixins.c311DirtyGuard],
      template: '<input />',
    })
    const wrapper = mount(Guard, { attachTo: document.body })
    wrapper.vm.c311DirtyStorageKey = 'c311.test.draft'
    wrapper.vm.c311MarkDirty(true)
    wrapper.vm.c311SaveDirtyDraft({ summary: 'keep', password: 'drop', nested: { accessToken: 'drop', description: 'keep' } })
    expect(JSON.parse(sessionStorage.getItem('c311.test.draft'))).toEqual({ summary: 'keep', nested: { description: 'keep' } })
    expect(wrapper.vm.c311ReadDirtyDraft()).toEqual({ summary: 'keep', nested: { description: 'keep' } })
  })

  it('protects dirty navigation and reloads server state only after confirmation', async () => {
    const Guard = Vue.extend({
      mixins: [mixins.c311DirtyGuard],
      template: '<input />',
    })
    const wrapper = mount(Guard)
    const reload = jest.fn()
    wrapper.vm.c311ReloadServerState = reload
    wrapper.vm.c311MarkDirty(true)
    const originalConfirm = window.confirm
    window.confirm = jest.fn(() => false)
    const next = jest.fn()
    const leaveGuard = wrapper.vm.$options.beforeRouteLeave
    const invokeLeaveGuard = Array.isArray(leaveGuard) ? leaveGuard[0] : leaveGuard
    invokeLeaveGuard.call(wrapper.vm, {}, {}, next)
    expect(next).toHaveBeenCalledWith(false)
    expect(reload).not.toHaveBeenCalled()
    window.confirm = jest.fn(() => true)
    invokeLeaveGuard.call(wrapper.vm, { path: '/new' }, { path: '/old' }, next)
    expect(next).toHaveBeenCalledWith()
    expect(reload).toHaveBeenCalledWith({ path: '/new' }, { path: '/old' })
    window.confirm = originalConfirm
  })

  it('renders responsive table and card alternatives without changing data', () => {
    const wrapper = mount(C311ResponsiveData, {
      propsData: { items: [{ id: '1', name: 'Request' }], columns: [{ key: 'id', label: 'ID' }, { key: 'name', label: 'Name' }] },
    })
    expect(wrapper.find('[data-c311-responsive-table]').exists()).toBe(true)
    expect(wrapper.find('[data-c311-responsive-table]').attributes('role')).toBe('region')
    expect(wrapper.find('[data-c311-responsive-table]').attributes('tabindex')).toBe('-1')
    expect(wrapper.find('[data-c311-responsive-cards]').exists()).toBe(true)
    expect(wrapper.text()).toContain('Request')
  })

  it('formats instants in the fixed New York timezone with EST/EDT', () => {
    expect(mockFormatC311DateTime('2026-01-15T15:00:00.000Z', 'en-US')).toBe('01/15/2026 10:00 AM EST')
    expect(mockFormatC311DateTime('2026-07-15T15:00:00.000Z', 'en-US')).toBe('07/15/2026 11:00 AM EDT')
  })

  it('covers public portal loading, validation, submit success, and provider errors', async () => {
    const provider = {
      listPortalRequests: jest.fn().mockResolvedValue({ items: [{ request_id: 'request-1', updated_at: '2026-01-15T15:00:00.000Z' }] }),
      submitPortalRequest: jest.fn().mockResolvedValue({ request_number: 'SR-2026-00002' }),
    }
    const wrapper = mount(Portal, {
      mocks: { ...mocks, $route: { name: 'c311.submit' }, $C311: { provider, session: { actor: { actor_id: 'actor-1' } } } },
      stubs: {
        'c311-app-shell': AppShellStub,
        'c311-data-state': DataStateStub,
        'c311-error-summary': ChildStub,
        'c311-capability-action': ChildStub,
        'c311-help-drawer': ChildStub,
        'c311-language-selector': ChildStub,
        'c311-main-nav': ChildStub,
        'c311-responsive-data': ChildStub,
      },
    })
    await wrapper.vm.load()
    expect(wrapper.vm.state).toBe('populated')
    expect(wrapper.vm.actorID).toBe('actor-1')
    expect(wrapper.vm.translatedColumns[3].format('2026-01-15T15:00:00.000Z')).toBe('01/15/2026 10:00 AM EST')

    await wrapper.vm.submit()
    expect(wrapper.vm.formErrors[0].code).toBe('REQUIRED')
    wrapper.vm.form.summary = 'Pothole'
    await wrapper.vm.submit()
    expect(provider.submitPortalRequest).toHaveBeenCalledWith(expect.objectContaining({ summary: 'Pothole' }))
    expect(wrapper.vm.statusMessage).toContain('SR-2026-00002')

    provider.listPortalRequests.mockRejectedValueOnce({ status: 503, retryable: true })
    await wrapper.vm.load()
    expect(wrapper.vm.state).toBe('retryable-error')
  })

  it('covers staff queue empty, populated, and terminal provider states', async () => {
    const provider = { listStaffRequests: jest.fn().mockResolvedValue({ items: [] }) }
    const wrapper = mount(Staff, {
      mocks: { ...mocks, $route: { name: 'c311.staff' }, $C311: { provider, session: { actor: { actor_id: 'staff-1' } } } },
      stubs: {
        'c311-app-shell': AppShellStub,
        'c311-data-state': DataStateStub,
        'c311-help-drawer': ChildStub,
        'c311-language-selector': ChildStub,
        'c311-main-nav': ChildStub,
        'c311-responsive-data': ChildStub,
      },
    })
    await wrapper.vm.load()
    expect(wrapper.vm.state).toBe('empty')
    expect(wrapper.vm.actorID).toBe('staff-1')
    expect(wrapper.vm.navItems).toEqual(expect.arrayContaining([
      expect.objectContaining({ route: '/c311/staff/reports', capability: 'report_catalogue' }),
      expect.objectContaining({ route: '/c311/staff/workflows', scope: 'workflow.execute' }),
    ]))
    expect(wrapper.vm.translatedColumns[4].format('2026-07-15T15:00:00.000Z')).toBe('07/15/2026 11:00 AM EDT')

    provider.listStaffRequests.mockRejectedValueOnce(new Error('unavailable'))
    await wrapper.vm.load()
    expect(wrapper.vm.state).toBe('terminal-error')
    expect(wrapper.vm.dataError.message).toBe('unavailable')
  })

  it('defines public portal and guarded staff routes with mock interaction gating', () => {
    const composeStatus = composeRoutes.find(route => route.name === 'c311.status')
    expect(composeStatus.path).toBe('/c311/status')
    expect(composeStatus.meta.c311.public).toBe(true)
    const composeInteraction = composeRoutes.find(route => route.name === 'c311.test.interaction')
    const composeNext = jest.fn()
    window.C311Mode = 'mock'
    composeInteraction.beforeEnter({}, {}, composeNext)
    expect(composeNext).toHaveBeenCalledWith()
    window.C311Mode = 'http'
    composeInteraction.beforeEnter({}, {}, composeNext)
    expect(composeNext).toHaveBeenLastCalledWith({ name: 'c311.not-found' })

    const staff = adminRoutes.find(route => route.name === 'c311.staff')
    expect(staff.meta.c311).toEqual(expect.objectContaining({ requiresAuth: true, route: 'staff_request_queue', capabilities: ['staff_request_queue'], scopes: ['service_requests.write'] }))
    expect(adminRoutes.find(route => route.name === 'c311.unauthorized').path).toBe('/c311/401')
    expect(adminRoutes.find(route => route.name === 'c311.forbidden').path).toBe('/c311/403')
    expect(adminRoutes.find(route => route.name === 'c311.not-found').path).toBe('/c311/404')
  })
})
