<template>
  <c311-app-shell mode="public" brand="City 311" title="City 311" :status-message="statusMessage">
    <template #nav>
      <div class="d-flex flex-wrap align-items-center gap-2">
        <c311-main-nav :items="navItems" label="Public navigation" />
        <c311-help-drawer help-key="public.request.submit" />
        <c311-language-selector :actor-id="actorID" />
      </div>
    </template>

    <c311-data-state :state="state" :messages="messages" :error="dataError" @retry="load">
      <template #populated>
        <c311-responsive-data
          :items="items"
          :columns="columns"
          row-key="request_id"
          label="My service requests"
        />
      </template>
    </c311-data-state>

    <section v-if="isSubmitRoute" aria-labelledby="submit-heading" class="mt-4">
      <h2 id="submit-heading" class="h4">Submit a service request</h2>
      <c311-error-summary :errors="formErrors" title="Review your request" />
      <form @submit.prevent="submit">
        <div class="form-group">
          <label for="c311-summary">Summary</label>
          <input
            id="c311-summary"
            v-model.trim="form.summary"
            class="form-control"
            :aria-invalid="hasError('summary') ? 'true' : 'false'"
            aria-describedby="c311-summary-help"
            @input="markDirty"
          >
          <small id="c311-summary-help" class="form-text text-muted">Briefly describe the issue.</small>
        </div>
        <div class="form-group">
          <label for="c311-description">Description</label>
          <textarea
            id="c311-description"
            v-model.trim="form.description"
            class="form-control"
            :aria-invalid="hasError('description') ? 'true' : 'false'"
            @input="markDirty"
          />
        </div>
        <c311-capability-action
          capability="portal_service_request_submit"
          :busy="submitting"
          :allow-anonymous="true"
          :explain-when-denied="true"
          type="submit"
        >
          Submit request
        </c311-capability-action>
      </form>
    </section>
  </c311-app-shell>
</template>

<script>
import { formatC311DateTime } from '@cortezaproject/corteza-js'
import { components, mixins, c311 } from '@cortezaproject/corteza-vue'

const {
  C311AppShell,
  C311DataState,
  C311ErrorSummary,
  C311CapabilityAction,
  C311HelpDrawer,
  C311LanguageSelector,
  C311MainNav,
  C311ResponsiveData,
} = components

const c311StateForError = c311?.c311StateForError
const c311DirtyGuard = mixins && mixins.c311DirtyGuard
  ? mixins.c311DirtyGuard
  : {
      data: () => ({ c311Dirty: false, c311DirtyStorageKey: '' }),
      methods: {
        c311MarkDirty (value = true) { this.c311Dirty = value },
        c311ReadDirtyDraft () { return null },
        c311SaveDirtyDraft () {},
        c311ClearDirtyDraft () { this.c311Dirty = false },
      },
    }

export default {
  name: 'C311Portal',
  components: {
    C311AppShell,
    C311DataState,
    C311ErrorSummary,
    C311CapabilityAction,
    C311HelpDrawer,
    C311LanguageSelector,
    C311MainNav,
    C311ResponsiveData,
  },
  mixins: [c311DirtyGuard],
  data: () => ({
    state: 'loading',
    items: [],
    messages: {},
    dataError: null,
    statusMessage: '',
    columns: [
      { key: 'request_number', label: 'Request number' },
      { key: 'summary', label: 'Summary' },
      { key: 'status', label: 'Status' },
      { key: 'updated_at', label: 'Updated', format: value => formatC311DateTime?.(value) || value || '' },
    ],
    navItems: [
      { route: '/c311/submit', label: 'Submit a request' },
      { route: '/c311/status', label: 'Check status' },
    ],
    form: {
      summary: '',
      description: '',
    },
    formErrors: [],
    submitting: false,
  }),
  computed: {
    actorID () {
      const actor = this.$C311 && this.$C311.session && this.$C311.session.actor
      return actor ? actor.actor_id : ''
    },
    isSubmitRoute () {
      return this.$route && this.$route.name === 'c311.submit'
    },
  },
  watch: {
    form: {
      deep: true,
      handler () {
        if (this.isSubmitRoute && this.c311Dirty) this.c311SaveDirtyDraft(this.form)
      },
    },
  },
  created () {
    this.c311DirtyStorageKey = 'c311.portal.submit'
    const draft = this.c311ReadDirtyDraft()
    if (draft) this.form = { ...this.form, ...draft }
    if (!this.isSubmitRoute) this.load()
    else {
      this.state = 'populated'
      this.statusMessage = ''
    }
  },
  methods: {
    markDirty () {
      if (this.c311MarkDirty) this.c311MarkDirty(true)
    },
    hasError (field) {
      return this.formErrors.some(error => error.field === field || error.field === `/${field}`)
    },
    async submit () {
      this.formErrors = []
      if (!this.form.summary) {
        this.formErrors = [{ field: 'summary', code: 'REQUIRED', message: 'Summary is required.' }]
        return
      }
      this.submitting = true
      try {
        const response = await this.$C311?.provider?.submitPortalRequest?.({
          summary: this.form.summary,
          description: this.form.description,
          service_type: 'GENERAL_INQUIRY',
          requester: { display_name: 'Portal visitor', email: 'visitor@example.test' },
        })
        if (!response) throw new Error('The request provider is unavailable.')
        this.form = { summary: '', description: '' }
        this.c311ClearDirtyDraft?.()
        this.statusMessage = `Request ${response.request_number} submitted.`
      } catch (error) {
        this.formErrors = error?.fieldErrors || error?.errors || [{ field: 'summary', code: 'SUBMIT_FAILED', message: error?.message || 'The request could not be submitted.' }]
      } finally {
        this.submitting = false
      }
    },
    async load () {
      this.state = 'loading'
      this.statusMessage = 'Loading requests.'
      try {
        const provider = this.$C311?.provider
        const page = await provider?.listPortalRequests?.()
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
