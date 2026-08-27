<template>
  <c311-app-shell mode="public" :brand="t('portal.brand', 'City 311')" :title="t('portal.title', 'City 311')" :status-message="statusMessage">
    <template #nav>
      <div class="d-flex flex-wrap align-items-center gap-2">
        <c311-main-nav :items="navItems" :label="t('navigation.public', 'Public navigation')" />
        <c311-help-drawer help-key="public.request.submit" :label="t('help.label', 'Help')" :title="t('help.title', 'Help')" :close-label="t('action.close', 'Close')" :content="t('help.public.request.submit', 'Submit a service request.')" />
        <c311-language-selector :actor-id="actorID" />
      </div>
    </template>

    <c311-data-state :state="state" :messages="messages" :error="dataError" @retry="load">
      <template #populated>
        <c311-responsive-data
          :items="items"
          :columns="translatedColumns"
          row-key="request_id"
          :label="t('portal.myRequests', 'My service requests')"
        />
      </template>
    </c311-data-state>

    <section v-if="isSubmitRoute" aria-labelledby="submit-heading" class="mt-4">
      <h2 id="submit-heading" class="h4">{{ t('portal.submit.title', 'Submit a service request') }}</h2>
      <c311-error-summary :errors="formErrors" :title="t('error.review', 'Review your request')" />
      <form @submit.prevent="submit">
        <div class="form-group">
          <label for="c311-summary">{{ t('field.summary', 'Summary') }}</label>
          <input
            id="c311-summary"
            v-model.trim="form.summary"
            class="form-control"
            :aria-invalid="hasError('summary') ? 'true' : 'false'"
            aria-describedby="c311-summary-help"
            @input="markDirty"
          >
          <small id="c311-summary-help" class="form-text text-muted">{{ t('field.summaryHelp', 'Briefly describe the issue.') }}</small>
        </div>
        <div class="form-group">
          <label for="c311-description">{{ t('field.description', 'Description') }}</label>
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
          {{ t('action.submit', 'Submit request') }}
        </c311-capability-action>
      </form>
    </section>
  </c311-app-shell>
</template>

<script>
import * as C311JS from '@cortezaproject/corteza-js'
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
function formatDate (value) {
  try {
    return typeof C311JS.formatC311DateTime === 'function' ? C311JS.formatC311DateTime(value) : value || ''
  } catch (_error) {
    return value || ''
  }
}
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
      { key: 'request_number', labelKey: 'field.requestNumber' },
      { key: 'summary', labelKey: 'field.summary' },
      { key: 'status', labelKey: 'field.status' },
      { key: 'updated_at', labelKey: 'field.updated', format: value => formatDate(value) },
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
    navItems () {
      return [
        { route: '/c311/submit', label: this.t('navigation.submit', 'Submit a request') },
        { route: '/c311/status', label: this.t('navigation.status', 'Check status') },
      ]
    },
    translatedColumns () {
      return this.columns.map(column => ({ ...column, label: this.t(column.labelKey, column.labelKey) }))
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
    t (key, fallback) {
      const translated = this.$t?.(`c311:${key}`)
      return translated && translated !== `c311:${key}` && translated !== key ? translated : fallback
    },
    markDirty () {
      if (this.c311MarkDirty) this.c311MarkDirty(true)
    },
    hasError (field) {
      return this.formErrors.some(error => error.field === field || error.field === `/${field}`)
    },
    async submit () {
      this.formErrors = []
      if (!this.form.summary) {
        this.formErrors = [{ field: 'summary', code: 'REQUIRED', message: this.t('error.summaryRequired', 'Summary is required.') }]
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
        if (!response) throw new Error(this.t('error.providerUnavailable', 'The request provider is unavailable.'))
        this.form = { summary: '', description: '' }
        this.c311ClearDirtyDraft?.()
        this.statusMessage = this.t('status.submitted', `Request ${response.request_number} submitted.`).replace('{{number}}', response.request_number)
      } catch (error) {
        this.formErrors = error?.fieldErrors || error?.errors || [{ field: 'summary', code: 'SUBMIT_FAILED', message: error?.message || this.t('error.submitFailed', 'The request could not be submitted.') }]
      } finally {
        this.submitting = false
      }
    },
    async load () {
      this.state = 'loading'
      this.statusMessage = this.t('status.loadingRequests', 'Loading requests.')
      try {
        const provider = this.$C311?.provider
        const page = await provider?.listPortalRequests?.()
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
