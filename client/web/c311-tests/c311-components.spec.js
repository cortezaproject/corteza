import Vue from 'vue'
import fs from 'fs'
import path from 'path'
import { mount } from '@vue/test-utils'
import { components, mixins, c311I18n } from './c311-components'
import { formatC311DateTime as mockFormatC311DateTime } from './time-test-helper'
import Portal from '../compose/src/views/C311/Portal.vue'
import PublicPortal from '../compose/src/views/C311/PublicPortal.vue'
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
  c311Identity: {
    validatePassword: value => value === 'ValidPassword1!' ? [] : ['too-short'],
    resetTokenFromLocation: () => 'ephemeral-token',
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

  it('persists non-sensitive portal drafts in local storage for browser restart recovery', () => {
    const Guard = Vue.extend({
      mixins: [mixins.c311DirtyGuard],
      template: '<input />',
    })
    const wrapper = mount(Guard, { attachTo: document.body })
    wrapper.vm.c311DirtyStorageKey = 'c311.portal.submit'
    wrapper.vm.c311MarkDirty(true)
    wrapper.vm.c311SaveDirtyDraft({ summary: 'keep after restart', password: 'drop', nested: { access_token: 'drop', description: 'keep' } })
    sessionStorage.clear()
    expect(JSON.parse(localStorage.getItem('c311.portal.submit'))).toEqual({ summary: 'keep after restart', nested: { description: 'keep' } })
    expect(wrapper.vm.c311ReadDirtyDraft()).toEqual({ summary: 'keep after restart', nested: { description: 'keep' } })
  })

  it('restores request fields without attachment tokens and warns that uploads must be repeated', () => {
    const provider = {}
    const wrapper = mount(Portal, {
      mocks: { ...mocks, $route: { name: 'c311.submit', query: {} }, $C311: { provider, session: { authenticated: false } } },
      stubs: { 'c311-app-shell': AppShellStub, 'c311-data-state': DataStateStub, 'c311-error-summary': ChildStub, 'c311-capability-action': ChildStub, 'c311-help-drawer': ChildStub, 'c311-language-selector': ChildStub, 'c311-main-nav': ChildStub, 'c311-responsive-data': ChildStub, 'router-link': RouterLinkStub },
    })
    wrapper.vm.form = {
      service_type: 'GENERAL_INQUIRY', summary: 'Saved summary', description: 'Saved description for restart.',
      requester: { display_name: 'Saved Resident', email: 'saved@example.test', phone: '' }, location: { address: '1 Main St', latitude: 40, longitude: -73 }, attachment_tokens: ['opaque-token'], custom_fields: { ward: 'NORTH' }, consent: true,
    }
    const stored = wrapper.vm.draftStorageValue()
    expect(stored).not.toHaveProperty('attachment_tokens')
    expect(stored.attachment_count).toBe(1)
    expect(JSON.stringify(stored)).not.toContain('opaque-token')

    wrapper.vm.form = { service_type: 'GENERAL_INQUIRY', summary: '', description: '', requester: { display_name: '', email: '', phone: '' }, location: { address: '', latitude: null, longitude: null }, attachment_tokens: [], custom_fields: {}, consent: false }
    wrapper.vm.restoreLocalDraft(stored)
    expect(wrapper.vm.form.summary).toBe('Saved summary')
    expect(wrapper.vm.form.location.address).toBe('1 Main St')
    expect(wrapper.vm.form.custom_fields).toEqual({ ward: 'NORTH' })
    expect(wrapper.vm.form.consent).toBe(true)
    expect(wrapper.vm.form.attachment_tokens).toEqual([])
    expect(wrapper.vm.attachmentRecoveryMessage).toContain('uploaded again')
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
    wrapper.vm.form = {
      service_type: 'GENERAL_INQUIRY', summary: 'Pothole', description: 'A sufficiently long description.',
      requester: { display_name: 'Anonymous', email: 'anon@example.test', phone: '' }, location: { address: '', latitude: null, longitude: null }, attachment_tokens: [], custom_fields: {}, consent: true,
    }
    await wrapper.vm.submit()
    expect(provider.submitPortalRequest.mock.calls[0][0]).toEqual(expect.objectContaining({ summary: 'Pothole' }))
    expect(wrapper.vm.statusMessage).toContain('SR-2026-00002')

    wrapper.vm.$route.name = 'c311.status'
    provider.listPortalRequests.mockRejectedValue({ status: 503, retryable: true })
    await wrapper.vm.load()
    expect(wrapper.vm.state).toBe('retryable-error')
  })

  it('uses the provider for federated entry and keeps the mock browser on the portal', async () => {
    const provider = { startFederatedSignIn: jest.fn().mockResolvedValue({ authorization_url: 'https://identity.example.test/oidc/authorize' }) }
    window.C311Mode = 'mock'
    const wrapper = mount(PublicPortal, {
      mocks: { ...mocks, $route: { name: 'c311.sign-in', query: {} }, $C311: { provider, session: { authenticated: false } } },
      stubs: { 'c311-app-shell': AppShellStub, 'c311-error-summary': ChildStub, 'c311-help-drawer': ChildStub, 'c311-language-selector': ChildStub, 'c311-main-nav': ChildStub, 'c311-data-state': DataStateStub, 'c311-responsive-data': ChildStub, 'router-link': RouterLinkStub },
    })
    await wrapper.vm.federated('oidc')
    expect(provider.startFederatedSignIn).toHaveBeenCalledWith('oidc')
    expect(wrapper.vm.federatedMessage).toContain('Redirecting')
    window.C311Mode = undefined
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
    const staffSubmit = adminRoutes.find(route => route.name === 'c311.staff.submit')
    expect(staffSubmit.meta.c311).toEqual(expect.objectContaining({ requiresAuth: true, route: 'staff_service_request_create', capabilities: ['staff_service_request_create'], scopes: ['service_requests.write'] }))
    expect(adminRoutes.find(route => route.name === 'c311.unauthorized').path).toBe('/c311/401')
    expect(adminRoutes.find(route => route.name === 'c311.forbidden').path).toBe('/c311/403')
    expect(adminRoutes.find(route => route.name === 'c311.not-found').path).toBe('/c311/404')
    expect(composeRoutes.find(route => route.name === 'c311.security-notice')).toBeUndefined()
  })

  it('renders anonymous public identity navigation and protects private entries', async () => {
    const provider = {
      getBranding: jest.fn().mockResolvedValue({ organisation_name: 'Fixture City' }),
      getPublicContent: jest.fn().mockResolvedValue({ body: '<p>Welcome</p>' }),
      getPublicHelp: jest.fn().mockResolvedValue({ body: '<p>Help</p>' }),
    }
    const wrapper = mount(PublicPortal, {
      mocks: { ...mocks, $route: { name: 'c311.portal', query: {} }, $C311: { provider, session: { authenticated: false } } },
      stubs: {
        'c311-app-shell': AppShellStub,
        'c311-data-state': DataStateStub,
        'c311-error-summary': ChildStub,
        'c311-help-drawer': ChildStub,
        'c311-language-selector': ChildStub,
        'c311-main-nav': { props: ['items'], template: '<nav><a v-for="item in items" :key="item.route" :href="item.route" :data-c311-route="item.route">{{ item.label }}</a></nav>' },
        'c311-responsive-data': ChildStub,
        'router-link': RouterLinkStub,
      },
    })
    await wrapper.vm.load()
    expect(wrapper.vm.state).toBe('populated')
    expect(wrapper.vm.navItems.map(item => item.route)).toEqual(expect.arrayContaining(['/c311/sign-in', '/c311/register']))
    expect(wrapper.vm.navItems.map(item => item.route)).not.toContain('/c311/account')
    expect(wrapper.find('[data-c311-page="home"]').exists()).toBe(true)
    expect(provider.getPublicContent).toHaveBeenCalledWith('HOME')
  })

  it('loads public HELP content and contextual help independently', async () => {
    const provider = {
      getBranding: jest.fn().mockResolvedValue({ organisation_name: 'Fixture City' }),
      getPublicContent: jest.fn().mockResolvedValue({ content_key: 'HELP', body: '<p>Published help</p>' }),
      getPublicHelp: jest.fn().mockResolvedValue({ body: '<p>Context help</p>' }),
    }
    const wrapper = mount(PublicPortal, {
      mocks: { ...mocks, $route: { name: 'c311.help', query: {} }, $C311: { provider, session: { authenticated: false } } },
      stubs: {
        'c311-app-shell': AppShellStub,
        'c311-data-state': DataStateStub,
        'c311-error-summary': ChildStub,
        'c311-help-drawer': ChildStub,
        'c311-language-selector': ChildStub,
        'c311-main-nav': ChildStub,
        'c311-responsive-data': ChildStub,
        'router-link': RouterLinkStub,
      },
    })
    await wrapper.vm.load()
    expect(provider.getPublicContent).toHaveBeenCalledWith('HELP')
    expect(provider.getPublicHelp).toHaveBeenCalledWith('public.request.submit', 'EN')
    expect(wrapper.vm.contentBody).toContain('Published help')
    expect(wrapper.vm.helpBody).toContain('Context help')
    expect(wrapper.vm.contentState).toBe('populated')
    expect(wrapper.vm.helpState).toBe('populated')
  })

  it('renders safe branding fields and keeps HELP content when contextual help fails', async () => {
    const provider = {
      getBranding: jest.fn().mockResolvedValue({ organisation_name: 'Fixture City', public_header: 'City services', public_footer: 'Support', primary_colour: '#155eef', accent_colour: 'red; background:url(javascript:bad)', font_family: 'Inter, system-ui', logo_url: 'javascript:bad' }),
      getPublicContent: jest.fn().mockResolvedValue({ content_key: 'HELP', body: '<p>Published help</p>' }),
      getPublicHelp: jest.fn().mockRejectedValue({ status: 503, retryable: true, message: 'context unavailable' }),
    }
    const wrapper = mount(PublicPortal, {
      mocks: { ...mocks, $route: { name: 'c311.help', query: {} }, $C311: { provider, session: { authenticated: false } } },
      stubs: { 'c311-app-shell': AppShellStub, 'c311-data-state': DataStateStub, 'c311-error-summary': ChildStub, 'c311-help-drawer': ChildStub, 'c311-language-selector': ChildStub, 'c311-main-nav': ChildStub, 'c311-responsive-data': ChildStub, 'router-link': RouterLinkStub },
    })
    await wrapper.vm.load()
    expect(wrapper.vm.contentState).toBe('populated')
    expect(wrapper.vm.helpState).toBe('retryable-error')
    expect(wrapper.vm.contentBody).toContain('Published help')
    expect(wrapper.vm.safeLogoUrl).toBe('')
    expect(wrapper.vm.brandStyle).toEqual(expect.objectContaining({ '--c311-primary-color': '#155eef', fontFamily: 'Inter, system-ui' }))
    expect(wrapper.vm.brandStyle).not.toHaveProperty('--c311-accent-color')
    expect(wrapper.find('[data-c311-branding]').text()).toContain('City services')
    expect(wrapper.find('[data-c311-branding-footer]').text()).toContain('Support')
  })

  it('keeps account maintenance input and session unchanged on conflict', async () => {
    const session = { authenticated: true, actor: { actor_id: 'actor-1' } }
    const provider = {
      getProfile: jest.fn().mockResolvedValue({ display_name: 'Alex', preferred_language: 'EN', version: 3 }),
      changeLoginIdentifier: jest.fn().mockRejectedValue({ status: 409, error: 'VERSION_CONFLICT', current_version: 4 }),
      changePassword: jest.fn(),
    }
    const router = { push: jest.fn() }
    const wrapper = mount(PublicPortal, {
      mocks: { ...mocks, $route: { name: 'c311.account', query: {} }, $router: router, $C311: { provider, session } },
      stubs: { 'c311-app-shell': AppShellStub, 'c311-error-summary': ChildStub, 'c311-help-drawer': ChildStub, 'c311-language-selector': ChildStub, 'c311-main-nav': ChildStub, 'c311-data-state': DataStateStub, 'c311-responsive-data': ChildStub, 'router-link': RouterLinkStub },
    })
    await wrapper.vm.load()
    wrapper.vm.forms.account.login_identifier = 'alex.new'
    wrapper.vm.forms.account.current_password = 'Current-password-1!'
    await wrapper.vm.changeLoginIdentifier()
    expect(provider.changeLoginIdentifier).toHaveBeenCalledWith({ current_password: 'Current-password-1!', login_identifier: 'alex.new' })
    expect(wrapper.vm.forms.account.login_identifier).toBe('alex.new')
    expect(wrapper.vm.$C311.session).toBe(session)
    expect(wrapper.vm.successMessage).toBe('')
    expect(wrapper.vm.formErrors.length).toBeGreaterThan(0)
  })

  it('reloads content and profile when the reused portal route changes', async () => {
    const route = { name: 'c311.portal', query: {} }
    const provider = {
      getBranding: jest.fn().mockResolvedValue({ organisation_name: 'Fixture City' }),
      getPublicContent: jest.fn().mockImplementation(contentKey => Promise.resolve({ content_key: contentKey, body: `<p>${contentKey}</p>` })),
      getPublicHelp: jest.fn().mockResolvedValue({ body: '<p>Context help</p>' }),
      listPortalRequests: jest.fn().mockResolvedValue({ items: [{ request_id: 'request-1' }] }),
      getProfile: jest.fn().mockResolvedValue({ display_name: 'Alex Example', preferred_language: 'EN', login_identifier: 'alex' }),
    }
    const wrapper = mount(PublicPortal, {
      mocks: { ...mocks, $route: route, $C311: { provider, session: { authenticated: true, actor: { actor_id: 'actor-1' } } } },
      stubs: { 'c311-app-shell': AppShellStub, 'c311-error-summary': ChildStub, 'c311-help-drawer': ChildStub, 'c311-language-selector': ChildStub, 'c311-main-nav': ChildStub, 'c311-data-state': DataStateStub, 'c311-responsive-data': ChildStub, 'router-link': RouterLinkStub },
    })

    await wrapper.vm.load()
    expect(wrapper.vm.contentBody).toContain('HOME')
    expect(wrapper.vm.contentKey).toBe('HOME')

    route.name = 'c311.services'
    await wrapper.vm.$options.watch.$route.handler.call(wrapper.vm, route, { name: 'c311.portal' })
    expect(provider.getPublicContent).toHaveBeenLastCalledWith('SERVICE_CATALOGUE')
    expect(wrapper.vm.contentBody).toContain('SERVICE_CATALOGUE')
    expect(wrapper.vm.contentKey).toBe('SERVICE_CATALOGUE')

    route.name = 'c311.help'
    await wrapper.vm.$options.watch.$route.handler.call(wrapper.vm, route, { name: 'c311.services' })
    expect(provider.getPublicContent).toHaveBeenLastCalledWith('HELP')
    expect(provider.getPublicHelp).toHaveBeenCalledWith('public.request.submit', 'EN')
    expect(wrapper.vm.contentBody).toContain('HELP')
    expect(wrapper.vm.helpBody).toContain('Context help')
    expect(wrapper.vm.contentKey).toBe('HELP')

    route.name = 'c311.requests'
    await wrapper.vm.$options.watch.$route.handler.call(wrapper.vm, route, { name: 'c311.help' })
    expect(provider.listPortalRequests).toHaveBeenCalled()
    expect(wrapper.vm.items).toHaveLength(1)

    route.name = 'c311.account'
    await wrapper.vm.$options.watch.$route.handler.call(wrapper.vm, route, { name: 'c311.requests' })
    expect(provider.getProfile).toHaveBeenCalled()
    expect(wrapper.vm.forms.account.display_name).toBe('Alex Example')
    expect(wrapper.vm.contentBody).toBe('')
    expect(wrapper.vm.contentKey).toBe('')
    expect(wrapper.vm.items).toEqual([])
  })

  it('clears stale state before reset and callback routes', async () => {
    const route = { name: 'c311.services', query: {} }
    const provider = {
      getBranding: jest.fn().mockResolvedValue({ organisation_name: 'Fixture City' }),
      getPublicContent: jest.fn().mockResolvedValue({ body: '<p>SERVICE_CATALOGUE</p>' }),
      getPublicHelp: jest.fn().mockResolvedValue({ body: '<p>Context help</p>' }),
      completeFederatedSignIn: jest.fn().mockResolvedValue({ authenticated: true, actor: { actor_id: 'actor-1' } }),
    }
    const wrapper = mount(PublicPortal, {
      mocks: { ...mocks, $route: route, $C311: { provider, session: { authenticated: false } } },
      stubs: { 'c311-app-shell': AppShellStub, 'c311-error-summary': ChildStub, 'c311-help-drawer': ChildStub, 'c311-language-selector': ChildStub, 'c311-main-nav': ChildStub, 'c311-data-state': DataStateStub, 'c311-responsive-data': ChildStub, 'router-link': RouterLinkStub },
    })
    await wrapper.vm.load()
    expect(wrapper.vm.contentBody).toContain('SERVICE_CATALOGUE')

    route.name = 'c311.reset-password'
    await wrapper.vm.$options.watch.$route.handler.call(wrapper.vm, route, { name: 'c311.services' })
    expect(wrapper.vm.contentBody).toBe('')
    expect(wrapper.vm.helpBody).toBe('')
    expect(wrapper.vm.dataError).toBe(null)
    expect(wrapper.vm.successMessage).toBe('')

    route.name = 'c311.auth.callback'
    route.query = { provider: 'oidc', code: 'opaque-code' }
    await wrapper.vm.$options.watch.$route.handler.call(wrapper.vm, route, { name: 'c311.reset-password' })
    expect(provider.completeFederatedSignIn).toHaveBeenCalledWith('oidc', { provider: 'oidc', code: 'opaque-code' })
    expect(wrapper.vm.contentBody).toBe('')
  })

  it('signs out immediately on SPA navigation and removes authenticated navigation', async () => {
    const route = { name: 'c311.requests', query: {} }
    const session = { authenticated: true, actor: { actor_id: 'actor-1', capabilities: ['profile_get'] } }
    const provider = {
      listPortalRequests: jest.fn().mockResolvedValue({ items: [] }),
      signOut: jest.fn().mockResolvedValue(undefined),
    }
    const runtime = { provider, session, clearSession: jest.fn(() => { runtime.session = null }) }
    const wrapper = mount(PublicPortal, {
      mocks: { ...mocks, $route: route, $C311: runtime },
      stubs: { 'c311-app-shell': AppShellStub, 'c311-error-summary': ChildStub, 'c311-help-drawer': ChildStub, 'c311-language-selector': ChildStub, 'c311-main-nav': ChildStub, 'c311-data-state': DataStateStub, 'c311-responsive-data': ChildStub, 'router-link': RouterLinkStub },
    })
    await wrapper.vm.load()

    route.name = 'c311.logout.callback'
    await wrapper.vm.$options.watch.$route.handler.call(wrapper.vm, route, { name: 'c311.requests' })
    expect(provider.signOut).toHaveBeenCalledTimes(1)
    expect(runtime.clearSession).toHaveBeenCalledTimes(1)
    expect(wrapper.vm.navItems.map(item => item.route)).not.toEqual(expect.arrayContaining(['/c311/account', '/c311/requests']))
    expect(wrapper.vm.successMessage).toContain('signed out')
  })

  it.each([
    [['profile_get'], false, false],
    [['profile_get', 'login_identifier_change'], true, false],
    [['profile_get', 'password_change'], false, true],
    [[], false, false],
  ])('only renders account maintenance actions for capabilities %j', async (capabilities, canChangeLogin, canChangePassword) => {
    const provider = { getProfile: jest.fn().mockResolvedValue({ display_name: 'Alex', preferred_language: 'EN', login_identifier: 'alex' }) }
    const route = { name: 'c311.account', query: {} }
    const runtime = {
      provider,
      session: { authenticated: true, actor: { actor_id: 'actor-1', capabilities } },
      can: capability => capabilities.includes(capability),
    }
    const wrapper = mount(PublicPortal, {
      mocks: { ...mocks, $route: route, $C311: runtime },
      stubs: { 'c311-app-shell': AppShellStub, 'c311-error-summary': ChildStub, 'c311-help-drawer': ChildStub, 'c311-language-selector': ChildStub, 'c311-main-nav': ChildStub, 'c311-data-state': DataStateStub, 'c311-responsive-data': ChildStub, 'router-link': RouterLinkStub },
    })
    await wrapper.vm.load()
    expect(wrapper.find('[data-c311-action="change-login-identifier"]').exists()).toBe(canChangeLogin)
    expect(wrapper.find('[data-c311-action="change-password"]').exists()).toBe(canChangePassword)
  })

  it('updates the current session after successful login-identifier and password changes', async () => {
    const session = { authenticated: true, actor: { actor_id: 'actor-1' } }
    const nextSession = { authenticated: true, actor: { actor_id: 'actor-1', login_identifier: 'alex.new' } }
    const provider = {
      getProfile: jest.fn().mockResolvedValue({ display_name: 'Alex', login_identifier: 'alex', preferred_language: 'EN', version: 1 }),
      changeLoginIdentifier: jest.fn().mockResolvedValue(nextSession),
      changePassword: jest.fn().mockResolvedValue(undefined),
    }
    const wrapper = mount(PublicPortal, {
      mocks: { ...mocks, $route: { name: 'c311.account', query: {} }, $C311: { provider, session }, $router: { push: jest.fn() } },
      stubs: { 'c311-app-shell': AppShellStub, 'c311-error-summary': ChildStub, 'c311-help-drawer': ChildStub, 'c311-language-selector': ChildStub, 'c311-main-nav': ChildStub, 'c311-data-state': DataStateStub, 'c311-responsive-data': ChildStub, 'router-link': RouterLinkStub },
    })
    await wrapper.vm.load()
    wrapper.vm.forms.account.login_identifier = 'alex.new'
    wrapper.vm.forms.account.current_password = 'Current-password-1!'
    await wrapper.vm.changeLoginIdentifier()
    expect(wrapper.vm.$C311.session).toBe(nextSession)
    wrapper.vm.forms.account.new_password = 'ValidPassword1!'
    await wrapper.vm.changePassword()
    expect(provider.changePassword).toHaveBeenCalledWith({ current_password: 'Current-password-1!', new_password: 'ValidPassword1!' })
    expect(wrapper.vm.forms.account.new_password).toBe('')
  })

  it('routes callback pending-link confirmation through server state', async () => {
    const provider = {
      startFederatedSignIn: jest.fn().mockResolvedValue({ authorization_url: 'https://identity.example.test/oidc/authorize' }),
      completeFederatedSignIn: jest.fn().mockResolvedValue({ outcome: 'link_confirmation_required', pending_link: { expires_at: '2099-01-15T16:00:00.000Z' } }),
      confirmAccountLink: jest.fn().mockResolvedValue({ authenticated: true, actor: null, preferred_language: 'EN', expires_at: null }),
    }
    const router = { push: jest.fn() }
    const runtime = { provider, session: { authenticated: false }, pendingFederated: null }
    const wrapper = mount(PublicPortal, {
      mocks: { ...mocks, $route: { name: 'c311.sign-in', query: {} }, $router: router, $C311: runtime },
      stubs: { 'c311-app-shell': AppShellStub, 'c311-error-summary': ChildStub, 'c311-help-drawer': ChildStub, 'c311-language-selector': ChildStub, 'c311-main-nav': ChildStub, 'c311-data-state': DataStateStub, 'c311-responsive-data': ChildStub, 'router-link': RouterLinkStub },
    })
    await wrapper.vm.federated('oidc')
    wrapper.vm.$route.query = { provider: 'oidc', code: 'fixture' }
    await wrapper.vm.completeFederated()
    expect(router.push).toHaveBeenCalledWith({ name: 'c311.auth.link.confirm' })
    expect(runtime.pendingFederated).toEqual({ expires_at: '2099-01-15T16:00:00.000Z' })
    window.C311Mode = 'mock'
    await wrapper.vm.confirmAccountLink()
    window.C311Mode = undefined
    expect(provider.confirmAccountLink).toHaveBeenCalledWith()
    expect(wrapper.vm.linkState).toBe('success')
    expect(runtime.pendingFederated).toBe(null)
    expect(router.push).not.toHaveBeenCalledWith(expect.objectContaining({ query: expect.anything() }))
  })

  it('cancels pending account linking through the provider before returning to sign-in', async () => {
    const provider = {}
    const router = { push: jest.fn() }
    const runtime = { provider, session: { authenticated: false }, pendingFederated: { expires_at: '2099-01-15T16:00:00.000Z' } }
    const wrapper = mount(PublicPortal, {
      mocks: { ...mocks, $route: { name: 'c311.auth.link.confirm', query: {} }, $router: router, $C311: runtime },
      stubs: { 'c311-app-shell': AppShellStub, 'c311-error-summary': ChildStub, 'c311-help-drawer': ChildStub, 'c311-language-selector': ChildStub, 'c311-main-nav': ChildStub, 'c311-data-state': DataStateStub, 'c311-responsive-data': ChildStub, 'router-link': RouterLinkStub },
    })
    await wrapper.vm.cancelAccountLink()
    expect(runtime.pendingFederated).toBe(null)
    expect(router.push).toHaveBeenCalledWith({ name: 'c311.sign-in' })
  })

  it('validates registration and reset forms, binds provider errors, and keeps token out of the URL', async () => {
    const provider = {
      registerAccount: jest.fn().mockRejectedValue({ status: 422, errors: [{ field: '/email', code: 'INVALID_FORMAT', message: 'Invalid email' }] }),
      confirmPasswordReset: jest.fn().mockResolvedValue({ message: 'Reset complete' }),
    }
    const wrapper = mount(PublicPortal, {
      mocks: { ...mocks, $route: { name: 'c311.register', query: {} }, $C311: { provider, session: { authenticated: false } } },
      stubs: { 'c311-app-shell': AppShellStub, 'c311-error-summary': ChildStub, 'c311-help-drawer': ChildStub, 'c311-language-selector': ChildStub, 'c311-main-nav': ChildStub, 'c311-data-state': DataStateStub, 'c311-responsive-data': ChildStub, 'router-link': RouterLinkStub },
    })
    await wrapper.vm.register()
    expect(wrapper.vm.formErrors.map(error => error.field)).toEqual(expect.arrayContaining(['display_name', 'email', 'login_identifier', 'password']))
    for (const login_identifier of ['ab', 'a'.repeat(65), 'Bad Name']) {
      wrapper.vm.forms.register = { display_name: 'Example', email: 'example@example.test', login_identifier, password: 'ValidPassword1!', preferred_language: 'EN' }
      provider.registerAccount.mockClear()
      await wrapper.vm.register()
      expect(wrapper.vm.formErrors).toEqual(expect.arrayContaining([expect.objectContaining({ field: 'login_identifier', code: 'INVALID_FORMAT' })]))
      expect(provider.registerAccount).not.toHaveBeenCalled()
    }
    wrapper.vm.forms.register = { display_name: 'Example', email: 'example@example.test', login_identifier: 'example', password: 'ValidPassword1!', preferred_language: 'EN' }
    await wrapper.vm.register()
    expect(provider.registerAccount).toHaveBeenCalledWith(expect.objectContaining({ email: 'example@example.test' }))
    await wrapper.setData({ resetToken: 'ephemeral-token', forms: { ...wrapper.vm.forms, reset: { password: 'ValidPassword1!' } } })
    await wrapper.vm.resetPassword()
    expect(provider.confirmPasswordReset).toHaveBeenCalledWith({ token: 'ephemeral-token', password: 'ValidPassword1!' })
    expect(wrapper.vm.successMessage).toContain('Reset complete')
  })

  it('maps public content failures to retryable state and supports authenticated requests', async () => {
    const provider = {
      getBranding: jest.fn().mockRejectedValue({ status: 503, retryable: true }),
      listPortalRequests: jest.fn().mockResolvedValue({ items: [] }),
    }
    const wrapper = mount(PublicPortal, {
      mocks: { ...mocks, $route: { name: 'c311.requests', query: {} }, $C311: { provider, session: { authenticated: true, actor: { actor_id: 'actor-1' } } } },
      stubs: { 'c311-app-shell': AppShellStub, 'c311-data-state': DataStateStub, 'c311-error-summary': ChildStub, 'c311-help-drawer': ChildStub, 'c311-language-selector': ChildStub, 'c311-main-nav': ChildStub, 'c311-responsive-data': ChildStub, 'router-link': RouterLinkStub },
    })
    await wrapper.vm.load()
    expect(wrapper.vm.state).toBe('retryable-error')
    expect(wrapper.vm.navItems.map(item => item.route)).toEqual(expect.arrayContaining(['/c311/account', '/c311/requests', '/c311/logout/callback']))
  })

  it('renders the complete FE-03 form and sends only contract fields', async () => {
    const provider = {
      submitPortalRequest: jest.fn().mockResolvedValue({ request_number: 'SR-2026-00041', status: 'SUBMITTED', version: 1 }),
      getProfile: jest.fn().mockResolvedValue({ constituent_id: 'constituent-1', display_name: 'Alex Example', emails: ['alex@example.test'], phone_numbers: [] }),
    }
    const wrapper = mount(Portal, {
      mocks: { ...mocks, $route: { name: 'c311.submit', query: {} }, $C311: { provider, session: { authenticated: true, actor: { actor_id: 'actor-1' } } } },
      stubs: { 'c311-app-shell': AppShellStub, 'c311-data-state': DataStateStub, 'c311-error-summary': ChildStub, 'c311-capability-action': ChildStub, 'c311-help-drawer': ChildStub, 'c311-language-selector': ChildStub, 'c311-main-nav': ChildStub, 'c311-responsive-data': ChildStub, 'router-link': RouterLinkStub },
    })
    await wrapper.vm.load()
    expect(wrapper.find('#c311-service-type').exists()).toBe(true)
    expect(wrapper.find('#c311-requester-name').exists()).toBe(true)
    expect(wrapper.find('#c311-requester-email').exists()).toBe(true)
    expect(wrapper.find('#c311-requester-phone').exists()).toBe(true)
    expect(wrapper.find('#c311-location-address').exists()).toBe(true)
    expect(wrapper.find('#c311-attachment-token').exists()).toBe(true)
    expect(wrapper.find('#c311-custom-fields').exists()).toBe(true)
    expect(wrapper.find('#c311-consent').exists()).toBe(true)
    expect(wrapper.vm.form.requester.display_name).toBe('Alex Example')
    wrapper.vm.form = {
      service_type: 'GENERAL_INQUIRY', summary: 'Need help', description: 'Please help with this city service.',
      requester: { display_name: 'Changed Name', email: 'changed@example.test', phone: '+15550100' },
      location: { address: '', latitude: null, longitude: null }, attachment_tokens: ['upload-00031'], custom_fields: { ward: 'NORTH' }, consent: true,
    }
    await wrapper.vm.submit()
    const [input, options] = provider.submitPortalRequest.mock.calls[0]
    expect(input).toEqual(expect.objectContaining({ service_type: 'GENERAL_INQUIRY', requester: { display_name: 'Changed Name', email: 'changed@example.test', phone: '+15550100' }, attachment_tokens: ['upload-00031'], custom_fields: { ward: 'NORTH' } }))
    expect(input).not.toHaveProperty('consent')
    expect(input).not.toHaveProperty('source_channel')
    expect(options.idempotencyKey).toBeTruthy()
    expect(wrapper.vm.submissionResult.status).toBe('SUBMITTED')
    expect(wrapper.vm.submissionResult.request_number).toBe('SR-2026-00041')
  })

  it('validates conditional location and consent while preserving valid input', async () => {
    const provider = { submitPortalRequest: jest.fn() }
    const wrapper = mount(Portal, {
      mocks: { ...mocks, $route: { name: 'c311.submit', query: {} }, $C311: { provider, session: { authenticated: false } } },
      stubs: { 'c311-app-shell': AppShellStub, 'c311-data-state': DataStateStub, 'c311-error-summary': ChildStub, 'c311-capability-action': ChildStub, 'c311-help-drawer': ChildStub, 'c311-language-selector': ChildStub, 'c311-main-nav': ChildStub, 'c311-responsive-data': ChildStub, 'router-link': RouterLinkStub },
    })
    wrapper.vm.form.summary = 'Keep this summary'
    wrapper.vm.form.description = 'Keep this valid description.'
    wrapper.vm.form.requester = { display_name: 'Anonymous', email: 'anon@example.test', phone: '' }
    wrapper.vm.form.service_type = 'POTHOLE'
    await wrapper.vm.submit()
    expect(provider.submitPortalRequest).not.toHaveBeenCalled()
    expect(wrapper.vm.form.summary).toBe('Keep this summary')
    expect(wrapper.vm.formErrors.map(error => error.field)).toEqual(expect.arrayContaining(['location.address', 'consent']))
    expect(wrapper.vm.state).toBe('validation-error')

    wrapper.vm.form.service_type = 'GENERAL_INQUIRY'
    wrapper.vm.form.consent = true
    wrapper.vm.form.location.latitude = 'not-a-number'
    wrapper.vm.form.description = 'Description with <b>markup</b>.'
    await wrapper.vm.submit()
    expect(provider.submitPortalRequest).not.toHaveBeenCalled()
    expect(wrapper.vm.formErrors.map(error => error.field)).toEqual(expect.arrayContaining(['description', 'location.latitude']))
  })

  it('reuses one idempotency key for double submit and maps server field errors', async () => {
    let resolveSubmit
    const provider = {
      submitPortalRequest: jest.fn().mockImplementation(() => new Promise(resolve => { resolveSubmit = resolve })),
    }
    const wrapper = mount(Portal, {
      mocks: { ...mocks, $route: { name: 'c311.submit', query: {} }, $C311: { provider, session: { authenticated: false } } },
      stubs: { 'c311-app-shell': AppShellStub, 'c311-data-state': DataStateStub, 'c311-error-summary': ChildStub, 'c311-capability-action': ChildStub, 'c311-help-drawer': ChildStub, 'c311-language-selector': ChildStub, 'c311-main-nav': ChildStub, 'c311-responsive-data': ChildStub, 'router-link': RouterLinkStub },
    })
    wrapper.vm.form = {
      service_type: 'GENERAL_INQUIRY', summary: 'Need help', description: 'Please help with this city service.',
      requester: { display_name: 'Anonymous', email: 'anon@example.test', phone: '' }, location: { address: '', latitude: null, longitude: null }, attachment_tokens: [], custom_fields: {}, consent: true,
    }
    const first = wrapper.vm.submit()
    const second = wrapper.vm.submit()
    expect(provider.submitPortalRequest).toHaveBeenCalledTimes(1)
    const firstKey = provider.submitPortalRequest.mock.calls[0][1].idempotencyKey
    expect(firstKey).toBeTruthy()
    resolveSubmit({ request_number: 'SR-2026-00042', status: 'SUBMITTED', version: 1 })
    await Promise.all([first, second])

    provider.submitPortalRequest.mockRejectedValueOnce({ status: 422, error: 'VALIDATION_ERROR', errors: [{ field: '/requester/email', code: 'INVALID_FORMAT', message: 'Invalid email' }] })
    wrapper.vm.form.requester.email = 'valid@example.test'
    await wrapper.vm.submit()
    expect(wrapper.vm.formErrors[0].field).toBe('/requester/email')
    expect(wrapper.vm.form.requester.display_name).toBe('Anonymous')
    expect(provider.submitPortalRequest.mock.calls[1][1].idempotencyKey).not.toBe(firstKey)
  })

  it('uses the staff assist operation only on the staff route', async () => {
    const provider = { createStaffServiceRequest: jest.fn().mockResolvedValue({ request: { request_number: 'SR-2026-00043', status: 'SUBMITTED' } }) }
    const wrapper = mount(Portal, {
      mocks: { ...mocks, $route: { name: 'c311.staff.submit', query: {} }, $C311: { provider, session: { authenticated: true, actor: { actor_id: 'staff-1', capabilities: ['staff_service_request_create'], scopes: ['service_requests.write'] } }, can: () => true, hasScope: () => true } },
      stubs: { 'c311-app-shell': AppShellStub, 'c311-data-state': DataStateStub, 'c311-error-summary': ChildStub, 'c311-capability-action': ChildStub, 'c311-help-drawer': ChildStub, 'c311-language-selector': ChildStub, 'c311-main-nav': ChildStub, 'c311-responsive-data': ChildStub, 'router-link': RouterLinkStub },
    })
    wrapper.vm.form = {
      service_type: 'GENERAL_INQUIRY', summary: 'Staff request', description: 'Staff member is helping submit this request.',
      requester: { display_name: 'Resident', email: 'resident@example.test', phone: '' }, location: { address: '', latitude: null, longitude: null }, attachment_tokens: [], custom_fields: {}, consent: true,
    }
    await wrapper.vm.submit()
    expect(provider.createStaffServiceRequest).toHaveBeenCalledWith(expect.objectContaining({ constituent: expect.any(Object), request: expect.objectContaining({ summary: 'Staff request' }) }))
    expect(provider.createStaffServiceRequest.mock.calls[0][0].request).not.toHaveProperty('source_channel')
  })

  it('does not load constituent profile on the staff assist route', async () => {
    const provider = { getProfile: jest.fn().mockRejectedValue({ status: 403, error: 'FORBIDDEN' }) }
    const wrapper = mount(Portal, {
      mocks: { ...mocks, $route: { name: 'c311.staff.submit', query: {} }, $C311: { provider, session: { authenticated: true, actor: { actor_id: 'staff-1', capabilities: ['staff_service_request_create'], scopes: ['service_requests.write'] } } } },
      stubs: { 'c311-app-shell': AppShellStub, 'c311-data-state': DataStateStub, 'c311-error-summary': ChildStub, 'c311-capability-action': ChildStub, 'c311-help-drawer': ChildStub, 'c311-language-selector': ChildStub, 'c311-main-nav': ChildStub, 'c311-responsive-data': ChildStub, 'router-link': RouterLinkStub },
    })
    await wrapper.vm.load()
    expect(provider.getProfile).not.toHaveBeenCalled()
    expect(wrapper.vm.state).toBe('populated')
    expect(wrapper.vm.dataError).toBe(null)
    expect(wrapper.find('#c311-requester-name').exists()).toBe(true)
  })

  it('keeps user input and exposes reload/reapply actions for draft version conflicts', async () => {
    const provider = {
      updateDraft: jest.fn().mockRejectedValue({ status: 409, error: 'VERSION_CONFLICT', current_version: 2, message: 'The record changed before your update.' }),
      getDraft: jest.fn().mockResolvedValue({ request_id: 'draft-fixture-001', version: 2, summary: 'Server summary', description: 'Server description', service_type: 'GENERAL_INQUIRY', primary_requester: { display_name: 'Server Resident', emails: ['server@example.test'], phone_numbers: [] }, custom_fields: {} }),
    }
    const wrapper = mount(Portal, {
      mocks: { ...mocks, $route: { name: 'c311.submit', query: {} }, $C311: { provider, session: { authenticated: true, actor: { actor_id: 'actor-1' } } } },
      stubs: { 'c311-app-shell': AppShellStub, 'c311-data-state': DataStateStub, 'c311-error-summary': ChildStub, 'c311-capability-action': ChildStub, 'c311-help-drawer': ChildStub, 'c311-language-selector': ChildStub, 'c311-main-nav': ChildStub, 'c311-responsive-data': ChildStub, 'router-link': RouterLinkStub },
    })
    wrapper.vm.draftID = 'draft-fixture-001'
    wrapper.vm.draftVersion = 1
    wrapper.vm.form.summary = 'Keep user summary'
    wrapper.vm.form.description = 'Keep this user description.'
    wrapper.vm.form.requester = { display_name: 'User Resident', email: 'user@example.test', phone: '' }
    wrapper.vm.form.consent = true
    await wrapper.vm.saveDraft()
    expect(wrapper.vm.versionConflict.current_version).toBe(2)
    expect(wrapper.vm.form.summary).toBe('Keep user summary')
    await wrapper.vm.reapplyDraft()
    expect(wrapper.vm.form.summary).toBe('Keep user summary')
    expect(wrapper.vm.versionConflict).toBe(null)

    wrapper.vm.versionConflict = { current_version: 2 }
    wrapper.vm.conflictDraft = { form: { summary: 'Keep user summary' }, customFieldsText: '{}' }
    await wrapper.vm.reloadDraft()
    expect(provider.getDraft).toHaveBeenCalledWith('draft-fixture-001')
    expect(wrapper.vm.form.summary).toBe('Server summary')
    expect(wrapper.vm.draftVersion).toBe(2)
  })

  it('reloads the correct data when the shared portal route changes', async () => {
    const route = Vue.observable({ name: 'c311.submit', path: '/c311/submit', query: {} })
    const provider = {
      listPortalRequests: jest.fn().mockResolvedValue({ items: [{ request_id: 'status-1', request_number: 'SR-2026-00001' }] }),
      getProfile: jest.fn(),
    }
    const wrapper = mount(Portal, {
      mocks: { ...mocks, $route: route, $C311: { provider, session: { authenticated: false } } },
      stubs: { 'c311-app-shell': AppShellStub, 'c311-data-state': DataStateStub, 'c311-error-summary': ChildStub, 'c311-capability-action': ChildStub, 'c311-help-drawer': ChildStub, 'c311-language-selector': ChildStub, 'c311-main-nav': ChildStub, 'c311-responsive-data': ChildStub, 'router-link': RouterLinkStub },
    })
    wrapper.vm.submissionResult = { request_number: 'SR-2026-00009', status: 'SUBMITTED' }
    route.name = 'c311.status'
    route.path = '/c311/status'
    await Vue.nextTick()
    await Vue.nextTick()
    expect(provider.listPortalRequests).toHaveBeenCalled()
    expect(wrapper.vm.showRequestList).toBe(true)
    expect(wrapper.vm.submissionResult).toBe(null)
    expect(wrapper.vm.items).toHaveLength(1)
  })

  it('gates each portal draft action by its matching capability', async () => {
    const session = capabilities => ({ authenticated: true, actor: { actor_id: 'actor-1', capabilities } })
    const mountPortal = capabilities => mount(Portal, {
      mocks: { ...mocks, $route: { name: 'c311.submit', query: {} }, $C311: { provider: {}, session: session(capabilities) } },
      stubs: { 'c311-app-shell': AppShellStub, 'c311-data-state': DataStateStub, 'c311-error-summary': ChildStub, 'c311-capability-action': ChildStub, 'c311-help-drawer': ChildStub, 'c311-language-selector': ChildStub, 'c311-main-nav': ChildStub, 'c311-responsive-data': ChildStub, 'router-link': RouterLinkStub },
    })

    const createOnly = mountPortal(['portal_draft_create'])
    await Vue.nextTick()
    expect(createOnly.find('[data-c311-action="save-draft"]').exists()).toBe(true)
    createOnly.vm.draftID = 'draft-1'
    await Vue.nextTick()
    expect(createOnly.find('[data-c311-action="save-draft"]').exists()).toBe(false)
    expect(createOnly.find('[data-c311-action="delete-draft"]').exists()).toBe(false)
    expect(createOnly.find('[data-c311-draft-submit-denied]').exists()).toBe(true)

    const updateDeleteSubmit = mountPortal(['portal_draft_update', 'portal_draft_delete', 'portal_draft_submit'])
    updateDeleteSubmit.vm.draftID = 'draft-1'
    await Vue.nextTick()
    expect(updateDeleteSubmit.find('[data-c311-action="save-draft"]').exists()).toBe(true)
    expect(updateDeleteSubmit.find('[data-c311-action="delete-draft"]').exists()).toBe(true)
    expect(updateDeleteSubmit.find('[data-c311-draft-submit-denied]').exists()).toBe(false)

    const noDraftCapabilities = mountPortal([])
    await Vue.nextTick()
    expect(noDraftCapabilities.find('[data-c311-action="save-draft"]').exists()).toBe(false)
  })

  it('ignores a stale portal list failure after navigating to a new route', async () => {
    let rejectList
    const provider = {
      listPortalRequests: jest.fn().mockImplementation(() => new Promise((_resolve, reject) => { rejectList = reject })),
    }
    const route = Vue.observable({ name: 'c311.status', path: '/c311/status', query: {} })
    const wrapper = mount(Portal, {
      mocks: { ...mocks, $route: route, $C311: { provider, session: { authenticated: false } } },
      stubs: { 'c311-app-shell': AppShellStub, 'c311-data-state': DataStateStub, 'c311-error-summary': ChildStub, 'c311-capability-action': ChildStub, 'c311-help-drawer': ChildStub, 'c311-language-selector': ChildStub, 'c311-main-nav': ChildStub, 'c311-responsive-data': ChildStub, 'router-link': RouterLinkStub },
    })
    await Vue.nextTick()
    route.name = 'c311.submit'
    route.path = '/c311/submit'
    await Vue.nextTick()
    await Vue.nextTick()
    rejectList({ status: 503, retryable: true })
    await Vue.nextTick()
    await Vue.nextTick()
    expect(wrapper.vm.showRequestList).toBe(false)
    expect(wrapper.vm.dataError).toBe(null)
    expect(wrapper.vm.state).toBe('populated')
  })

  it('ignores a stale submit failure after navigating away from the form', async () => {
    let rejectSubmit
    const provider = {
      submitPortalRequest: jest.fn().mockImplementation(() => new Promise((_resolve, reject) => { rejectSubmit = reject })),
    }
    const route = Vue.observable({ name: 'c311.submit', path: '/c311/submit', query: {} })
    const wrapper = mount(Portal, {
      mocks: { ...mocks, $route: route, $C311: { provider, session: { authenticated: false } } },
      stubs: { 'c311-app-shell': AppShellStub, 'c311-data-state': DataStateStub, 'c311-error-summary': ChildStub, 'c311-capability-action': ChildStub, 'c311-help-drawer': ChildStub, 'c311-language-selector': ChildStub, 'c311-main-nav': ChildStub, 'c311-responsive-data': ChildStub, 'router-link': RouterLinkStub },
    })
    wrapper.vm.form = {
      service_type: 'GENERAL_INQUIRY', summary: 'Valid summary', description: 'A valid description for this request.',
      requester: { display_name: 'Resident', email: 'resident@example.test', phone: '' }, location: { address: '', latitude: null, longitude: null }, attachment_tokens: [], custom_fields: {}, consent: true,
    }
    const submit = wrapper.vm.submit()
    await Vue.nextTick()
    route.name = 'c311.status'
    route.path = '/c311/status'
    await Vue.nextTick()
    await Vue.nextTick()
    rejectSubmit({ status: 503, retryable: true })
    await submit
    await Vue.nextTick()
    expect(wrapper.vm.showRequestList).toBe(true)
    expect(wrapper.vm.dataError).toBe(null)
  })

  it('uploads an attachment through the provider and submits only its opaque token', async () => {
    const provider = {
      uploadPortalAttachment: jest.fn().mockResolvedValue({ attachment_token: 'opaque-upload-token' }),
      submitPortalRequest: jest.fn().mockResolvedValue({ request_number: 'SR-2026-00050', status: 'SUBMITTED', version: 1 }),
    }
    const wrapper = mount(Portal, {
      mocks: { ...mocks, $route: { name: 'c311.submit', query: {} }, $C311: { provider, session: { authenticated: false } } },
      stubs: { 'c311-app-shell': AppShellStub, 'c311-data-state': DataStateStub, 'c311-error-summary': ChildStub, 'c311-capability-action': ChildStub, 'c311-help-drawer': ChildStub, 'c311-language-selector': ChildStub, 'c311-main-nav': ChildStub, 'c311-responsive-data': ChildStub, 'router-link': RouterLinkStub },
    })
    const file = new File(['fixture'], 'fixture.txt', { type: 'text/plain' })
    const input = wrapper.find('#c311-attachment-file').element
    Object.defineProperty(input, 'files', { value: [file] })
    await wrapper.find('#c311-attachment-file').trigger('change')
    await Vue.nextTick()
    expect(provider.uploadPortalAttachment).toHaveBeenCalledWith(expect.objectContaining({ filename: 'fixture.txt', media_type: 'text/plain', file }))
    expect(wrapper.vm.form.attachment_tokens).toEqual(['opaque-upload-token'])
    wrapper.vm.form = { service_type: 'GENERAL_INQUIRY', summary: 'Need help', description: 'Please help with this city service.', requester: { display_name: 'Resident', email: 'resident@example.test', phone: '' }, location: { address: '', latitude: null, longitude: null }, attachment_tokens: ['opaque-upload-token'], custom_fields: {}, consent: true }
    await wrapper.vm.submit()
    expect(provider.submitPortalRequest.mock.calls[0][0].attachment_tokens).toEqual(['opaque-upload-token'])
    expect(provider.submitPortalRequest.mock.calls[0][0]).not.toHaveProperty('attachment_contents')
  })

  it('allows an attachment upload to be retried after a retryable failure', async () => {
    const provider = {
      uploadPortalAttachment: jest.fn()
        .mockRejectedValueOnce({ status: 503, error: 'TEMPORARILY_UNAVAILABLE', retryable: true })
        .mockResolvedValueOnce({ attachment_token: 'opaque-retry-token' }),
    }
    const wrapper = mount(Portal, {
      mocks: { ...mocks, $route: { name: 'c311.submit', query: {} }, $C311: { provider, session: { authenticated: false } } },
      stubs: { 'c311-app-shell': AppShellStub, 'c311-data-state': DataStateStub, 'c311-error-summary': ChildStub, 'c311-capability-action': ChildStub, 'c311-help-drawer': ChildStub, 'c311-language-selector': ChildStub, 'c311-main-nav': ChildStub, 'c311-responsive-data': ChildStub, 'router-link': RouterLinkStub },
    })
    const file = new File(['fixture'], 'fixture.txt', { type: 'text/plain' })
    const input = wrapper.find('#c311-attachment-file').element
    Object.defineProperty(input, 'files', { configurable: true, value: [file] })
    await wrapper.find('#c311-attachment-file').trigger('change')
    await Vue.nextTick()
    await Vue.nextTick()
    expect(provider.uploadPortalAttachment).toHaveBeenCalledTimes(1)
    expect(wrapper.vm.state).toBe('retryable-error')

    Object.defineProperty(input, 'files', { configurable: true, value: [file] })
    await wrapper.find('#c311-attachment-file').trigger('change')
    await Vue.nextTick()
    await Vue.nextTick()
    expect(provider.uploadPortalAttachment).toHaveBeenCalledTimes(2)
    expect(wrapper.vm.form.attachment_tokens).toEqual(['opaque-retry-token'])
  })

})
