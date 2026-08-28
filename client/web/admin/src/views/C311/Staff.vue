<template>
  <c311-app-shell mode="staff" :brand="t('staff.brand', 'City 311 staff')" :title="t('staff.title', 'Requests')" :status-message="statusMessage">
    <template #nav>
      <div class="d-flex flex-wrap align-items-center gap-2">
        <c311-main-nav :items="navItems" :label="t('navigation.staff', 'Staff navigation')" />
        <c311-help-drawer help-key="staff.request.triage" :label="t('help.label', 'Help')" :title="t('help.title', 'Help')" :close-label="t('action.close', 'Close')" :content="t('help.staff.request.triage', 'Review and classify a request.')" />
        <c311-language-selector :actor-id="actorID" />
      </div>
    </template>

    <c311-data-state :state="state" :error="dataError" @retry="load">
      <template #populated>
        <c311-responsive-data
          :items="items"
          :columns="translatedColumns"
          row-key="request_id"
          :label="t('staff.queue', 'Staff request queue')"
        />
      </template>
    </c311-data-state>
  </c311-app-shell>
</template>

<script>
import * as C311JS from '@cortezaproject/corteza-js'
import { components, c311 } from '@cortezaproject/corteza-vue'

const {
  C311AppShell,
  C311DataState,
  C311HelpDrawer,
  C311LanguageSelector,
  C311MainNav,
  C311ResponsiveData,
} = components
const c311StateForError = c311?.c311StateForError
function formatDate (value) {
  try {
    return typeof C311JS.formatC311DateTime === 'function' ? C311JS.formatC311DateTime(value) : value || ''
  } catch (_error) {
    return value || ''
  }
}

export default {
  name: 'C311Staff',
  components: {
    C311AppShell,
    C311DataState,
    C311HelpDrawer,
    C311LanguageSelector,
    C311MainNav,
    C311ResponsiveData,
  },
  data: () => ({
    state: 'loading',
    dataError: null,
    statusMessage: '',
    items: [],
    columns: [
      { key: 'request_number', labelKey: 'field.requestNumber' },
      { key: 'summary', labelKey: 'field.summary' },
      { key: 'status', labelKey: 'field.status' },
      { key: 'owning_department', labelKey: 'field.department' },
      { key: 'updated_at', labelKey: 'field.updated', format: value => formatDate(value) },
    ],
  }),
  computed: {
    actorID () {
      return this.$C311 && this.$C311.session && this.$C311.session.actor ? this.$C311.session.actor.actor_id : ''
    },
    navItems () {
      return [
        { route: '/c311/staff', label: this.t('navigation.requests', 'Requests'), capability: 'staff_request_queue' },
        { route: '/c311/staff/reports', label: this.t('navigation.reports', 'Reports'), capability: 'report_catalogue' },
        { route: '/c311/staff/workflows', label: this.t('navigation.workflows', 'Workflows'), capability: 'workflow_list', scope: 'workflow.execute' },
      ]
    },
    translatedColumns () {
      return this.columns.map(column => ({ ...column, label: this.t(column.labelKey, column.labelKey) }))
    },
  },
  created () {
    this.load()
  },
  methods: {
    t (key, fallback) {
      const translated = this.$t?.(`c311:${key}`)
      return translated && translated !== `c311:${key}` && translated !== key ? translated : fallback
    },
    async load () {
      this.state = 'loading'
      this.statusMessage = this.t('status.loadingRequests', 'Loading requests.')
      try {
        const provider = this.$C311?.provider
        const page = await provider?.listStaffRequests?.()
        this.items = page?.items || []
        this.dataError = null
        this.state = this.items.length ? 'populated' : 'empty'
        this.statusMessage = this.items.length ? this.t('status.requestsLoaded', 'Requests loaded.') : this.t('status.noRequests', 'No requests found.')
      } catch (error) {
        this.dataError = error
        this.state = c311StateForError?.(error) || (error?.retryable ? 'retryable-error' : 'terminal-error')
        this.statusMessage = this.t('status.requestListUnavailable', 'Request list unavailable.')
      }
    },
  },
}
</script>
