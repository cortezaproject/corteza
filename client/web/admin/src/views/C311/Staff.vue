<template>
  <c311-app-shell mode="staff" brand="City 311 staff" title="Requests" :status-message="statusMessage">
    <template #nav>
      <div class="d-flex flex-wrap align-items-center gap-2">
        <c311-main-nav :items="navItems" label="Staff navigation" />
        <c311-help-drawer help-key="staff.request.triage" />
        <c311-language-selector :actor-id="actorID" />
      </div>
    </template>

    <c311-data-state :state="state" :error="dataError" @retry="load">
      <template #populated>
        <c311-responsive-data
          :items="items"
          :columns="columns"
          row-key="request_id"
          label="Staff request queue"
        />
      </template>
    </c311-data-state>
  </c311-app-shell>
</template>

<script>
import { formatC311DateTime } from '@cortezaproject/corteza-js'
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
      { key: 'request_number', label: 'Request number' },
      { key: 'summary', label: 'Summary' },
      { key: 'status', label: 'Status' },
      { key: 'owning_department', label: 'Department' },
      { key: 'updated_at', label: 'Updated', format: value => formatC311DateTime?.(value) || value || '' },
    ],
    navItems: [
      { route: '/c311/staff', label: 'Requests', capability: 'staff_request_queue' },
      { route: '/c311/staff/reports', label: 'Reports', capability: 'report_catalogue' },
      { route: '/c311/staff/workflows', label: 'Workflows', capability: 'workflow_list' },
    ],
  }),
  computed: {
    actorID () {
      return this.$C311 && this.$C311.session && this.$C311.session.actor ? this.$C311.session.actor.actor_id : ''
    },
  },
  created () {
    this.load()
  },
  methods: {
    async load () {
      this.state = 'loading'
      this.statusMessage = 'Loading requests.'
      try {
        const provider = this.$C311?.provider
        const page = await provider?.listStaffRequests?.()
        this.items = page?.items || []
        this.dataError = null
        this.state = this.items.length ? 'populated' : 'empty'
        this.statusMessage = this.items.length ? 'Requests loaded.' : 'No requests found.'
      } catch (error) {
        this.dataError = error
        this.state = c311StateForError?.(error) || (error?.retryable ? 'retryable-error' : 'terminal-error')
        this.statusMessage = 'Request list unavailable.'
      }
    },
  },
}
</script>
