<template>
  <c311-app-shell mode="public" :brand="t('portal.brand', 'City 311')" :title="shellTitle" :status-message="statusMessage">
    <template #nav>
      <div class="d-flex flex-wrap align-items-center gap-2">
        <c311-main-nav :items="navItems" :label="t('navigation.public', 'Public navigation')" />
        <c311-help-drawer help-key="public.request.submit" :label="t('help.label', 'Help')" :title="t('help.title', 'Help')" :close-label="t('action.close', 'Close')" :content="t('help.public.request.submit', 'Submit a service request.')" />
        <c311-language-selector :actor-id="actorID" />
      </div>
    </template>

    <c311-data-state v-if="showRequestList" :state="state" :messages="messages" :error="dataError" @retry="load">
      <template #populated>
        <c311-responsive-data :items="items" :columns="translatedColumns" row-key="request_id" :label="t('portal.myRequests', 'My service requests')" />
      </template>
    </c311-data-state>

    <section v-if="isRequestFormRoute" aria-labelledby="submit-heading" class="mt-4 c311-request-page">
      <h2 id="submit-heading" class="h4">{{ isStaffAssistRoute ? t('portal.submit.staffTitle', 'Submit a request for a resident') : t('portal.submit.title', 'Submit a service request') }}</h2>
      <p v-if="isStaffAssistRoute" class="text-muted">{{ t('portal.submit.staffDescription', 'This request will be created using your staff permissions.') }}</p>
      <c311-error-summary :errors="formErrors" :title="t('error.review', 'Review your request')" />
      <div v-if="versionConflict" class="alert alert-warning" role="alert" data-c311-version-conflict>
        <p class="mb-2">{{ t('status.versionConflict', 'The draft changed on the server. Current server version: ' + versionConflict.current_version + '.', { version: versionConflict.current_version }) }}</p>
        <div class="d-flex flex-wrap gap-2">
          <button class="btn btn-outline-secondary" type="button" data-c311-action="reload-draft" @click="reloadDraft">{{ t('action.reloadDraft', 'Reload server version') }}</button>
          <button class="btn btn-outline-primary" type="button" data-c311-action="reapply-draft" @click="reapplyDraft">{{ t('action.reapplyDraft', 'Reapply my changes') }}</button>
        </div>
      </div>
      <form @submit.prevent="submit">
        <div class="form-group">
          <label for="c311-service-type">{{ t('field.serviceType', 'Service type') }}</label>
          <select id="c311-service-type" v-model="form.service_type" class="form-control" :aria-invalid="hasError('service_type') ? 'true' : 'false'" @change="markDirty">
            <option v-for="serviceType in serviceTypes" :key="serviceType" :value="serviceType">{{ serviceType }}</option>
          </select>
        </div>
        <div class="form-group">
          <label for="c311-summary">{{ t('field.summary', 'Summary') }}</label>
          <input id="c311-summary" v-model.trim="form.summary" class="form-control" :aria-invalid="hasError('summary') ? 'true' : 'false'" aria-describedby="c311-summary-help" @input="markDirty">
          <small id="c311-summary-help" class="form-text text-muted">{{ t('field.summaryHelp', 'Briefly describe the issue.') }}</small>
        </div>
        <div class="form-group">
          <label for="c311-description">{{ t('field.description', 'Description') }}</label>
          <textarea id="c311-description" v-model.trim="form.description" class="form-control" rows="4" :aria-invalid="hasError('description') ? 'true' : 'false'" @input="markDirty" />
        </div>

        <fieldset class="c311-request-section">
          <legend>{{ t('portal.submit.requester', 'Requester') }}</legend>
          <div class="form-group">
            <label for="c311-requester-name">{{ t('field.displayName', 'Display name') }}</label>
            <input id="c311-requester-name" v-model.trim="form.requester.display_name" class="form-control" autocomplete="name" :aria-invalid="hasError('requester.display_name') ? 'true' : 'false'" @input="markDirty">
          </div>
          <div class="form-group">
            <label for="c311-requester-email">{{ t('field.email', 'Email') }}</label>
            <input id="c311-requester-email" v-model.trim="form.requester.email" class="form-control" type="email" autocomplete="email" :aria-invalid="hasError('requester.email') ? 'true' : 'false'" @input="markDirty">
          </div>
          <div class="form-group">
            <label for="c311-requester-phone">{{ t('field.phone', 'Phone (optional)') }}</label>
            <input id="c311-requester-phone" v-model.trim="form.requester.phone" class="form-control" type="tel" autocomplete="tel" :aria-invalid="hasError('requester.phone') ? 'true' : 'false'" @input="markDirty">
          </div>
        </fieldset>

        <fieldset class="c311-request-section">
          <legend>{{ t('portal.submit.location', 'Location') }}</legend>
          <div class="form-group">
            <label for="c311-location-address">{{ t('field.address', 'Address') }}</label>
            <input id="c311-location-address" v-model.trim="form.location.address" class="form-control" autocomplete="address-line1" :aria-invalid="hasError('location.address') ? 'true' : 'false'" @input="markDirty">
          </div>
          <div class="c311-request-grid">
            <div class="form-group">
              <label for="c311-location-latitude">{{ t('field.latitude', 'Latitude (optional)') }}</label>
              <input id="c311-location-latitude" v-model.number="form.location.latitude" class="form-control" type="number" min="-90" max="90" step="0.0001" :aria-invalid="hasError('location.latitude') ? 'true' : 'false'" @input="markDirty">
            </div>
            <div class="form-group">
              <label for="c311-location-longitude">{{ t('field.longitude', 'Longitude (optional)') }}</label>
              <input id="c311-location-longitude" v-model.number="form.location.longitude" class="form-control" type="number" min="-180" max="180" step="0.0001" :aria-invalid="hasError('location.longitude') ? 'true' : 'false'" @input="markDirty">
            </div>
          </div>
        </fieldset>

        <fieldset class="c311-request-section">
          <legend>{{ t('portal.submit.attachments', 'Attachments') }}</legend>
          <div class="form-group">
            <label for="c311-attachment-file">{{ t('field.attachmentFile', 'Upload an attachment') }}</label>
            <input id="c311-attachment-file" class="form-control-file" type="file" accept="image/jpeg,image/png,application/pdf,text/plain,application/vnd.openxmlformats-officedocument.wordprocessingml.document" :disabled="attachmentUploading || form.attachment_tokens.length >= 5" @change="uploadAttachment">
            <small class="form-text text-muted">{{ t('portal.submit.attachmentFileHelp', 'JPEG, PNG, PDF, text, or DOCX up to 10 MB.') }}</small>
          </div>
          <div class="form-group">
            <label for="c311-attachment-token">{{ t('field.attachmentToken', 'Attachment token') }}</label>
            <div class="input-group">
              <input id="c311-attachment-token" v-model.trim="attachmentTokenInput" class="form-control" type="text" @input="markDirty" @keyup.enter.prevent="addAttachmentToken">
              <div class="input-group-append"><button class="btn btn-outline-secondary" type="button" data-c311-action="add-attachment-token" @click="addAttachmentToken">{{ t('action.add', 'Add') }}</button></div>
            </div>
          </div>
          <ul v-if="form.attachment_tokens.length" class="list-unstyled" data-c311-attachment-list>
            <li v-for="token in form.attachment_tokens" :key="token" class="d-flex align-items-center justify-content-between border-bottom py-1">
              <code>{{ token }}</code>
              <button class="btn btn-sm btn-link" type="button" :data-c311-action="`remove-attachment-${token}`" @click="removeAttachmentToken(token)">{{ t('action.remove', 'Remove') }}</button>
            </li>
          </ul>
          <small v-if="attachmentUploading" class="form-text text-info" role="status">{{ t('status.uploadingAttachment', 'Uploading attachment...') }}</small>
          <p v-if="attachmentRecoveryMessage" class="alert alert-warning" role="status" data-c311-attachment-recovery>{{ attachmentRecoveryMessage }}</p>
          <small class="form-text text-muted">{{ t('portal.submit.attachmentHelp', 'Only the returned attachment token is sent with your request.') }}</small>
        </fieldset>

        <div class="form-group">
          <label for="c311-custom-fields">{{ t('field.customFields', 'Configured custom fields (JSON)') }}</label>
          <textarea id="c311-custom-fields" v-model="customFieldsText" class="form-control" rows="3" :aria-invalid="hasError('custom_fields') ? 'true' : 'false'" @input="onCustomFieldsInput" />
        </div>

        <div class="form-check mb-3">
          <input id="c311-consent" v-model="form.consent" class="form-check-input" type="checkbox" :aria-invalid="hasError('consent') ? 'true' : 'false'" @change="markDirty">
          <label class="form-check-label" for="c311-consent">{{ t('field.consent', 'I confirm that the information provided is accurate.') }}</label>
        </div>

        <button v-if="!submissionResult && !isStaffAssistRoute && canSubmitCurrentRequest" class="btn btn-primary" type="submit" data-c311-action="submit-request" :disabled="submitting">
          <span v-if="submitting" class="spinner-border spinner-border-sm mr-1" aria-hidden="true" />
          {{ t('action.submit', 'Submit request') }}
        </button>
        <p v-else-if="!submissionResult && !isStaffAssistRoute && draftID" class="text-muted small" role="note" data-c311-draft-submit-denied>{{ t('action.draftSubmitDenied', 'You do not have permission to submit this draft.') }}</p>
        <c311-capability-action v-else-if="!submissionResult" capability="staff_service_request_create" scope="service_requests.write" :busy="submitting" :explain-when-denied="true" :denied-label="t('action.denied', 'This action is unavailable for your role.')" data-c311-action="submit-staff-request" type="submit">
          {{ t('action.submit', 'Submit request') }}
        </c311-capability-action>
      </form>

      <div v-if="isAuthenticated && !isStaffAssistRoute && (canCreateDraft || canUpdateDraft || canDeleteDraft)" class="mt-3 d-flex flex-wrap gap-2">
        <button v-if="!draftID && canCreateDraft" class="btn btn-outline-primary" type="button" data-c311-action="save-draft" :disabled="draftSaving" @click="saveDraft">{{ draftSaving ? t('status.saving', 'Saving...') : t('action.saveDraft', 'Save draft') }}</button>
        <button v-if="draftID && canUpdateDraft" class="btn btn-outline-primary" type="button" data-c311-action="save-draft" :disabled="draftSaving" @click="saveDraft">{{ draftSaving ? t('status.saving', 'Saving...') : t('action.saveDraft', 'Save draft') }}</button>
        <button v-if="draftID && canDeleteDraft" class="btn btn-outline-danger" type="button" data-c311-action="delete-draft" :disabled="draftSaving" @click="deleteDraft">{{ t('action.deleteDraft', 'Delete draft') }}</button>
      </div>
      <p v-else-if="isAuthenticated && !isStaffAssistRoute && draftID" class="text-muted small mt-3" role="note" data-c311-draft-actions-denied>{{ t('action.draftActionsDenied', 'You do not have permission to change this draft.') }}</p>

      <p v-if="successMessage" class="alert alert-success mt-3" role="status" data-c311-submit-success>{{ successMessage }}</p>
      <div v-if="submissionResult" class="alert alert-success" data-c311-submission-result>
        <strong>{{ t('portal.submit.confirmed', 'Request submitted') }}</strong>
        <span class="d-block">{{ t('field.requestNumber', 'Request number') }}: {{ submissionResult.request_number }}</span>
        <span class="d-block">{{ t('field.status', 'Status') }}: {{ submissionResult.status }}</span>
      </div>
    </section>
  </c311-app-shell>
</template>

<script>
import * as C311JS from '@cortezaproject/corteza-js'
import { components, mixins, c311 } from '@cortezaproject/corteza-vue'

const { C311AppShell, C311DataState, C311ErrorSummary, C311CapabilityAction, C311HelpDrawer, C311LanguageSelector, C311MainNav, C311ResponsiveData } = components
const c311StateForError = c311?.c311StateForError
const serviceTypes = C311JS.SERVICE_TYPES || ['TREE_MAINTENANCE', 'POTHOLE', 'MISSED_TRASH', 'GENERAL_INQUIRY']
const serviceTypeRules = C311JS.SERVICE_TYPE_RULES || {
  GENERAL_INQUIRY: { location_required: false, confirmed_coordinates_required: false },
  MISSED_TRASH: { location_required: true, confirmed_coordinates_required: true },
  POTHOLE: { location_required: true, confirmed_coordinates_required: true },
  TREE_MAINTENANCE: { location_required: true, confirmed_coordinates_required: true },
}

const c311DirtyGuard = mixins?.c311DirtyGuard || {
  data: () => ({ c311Dirty: false, c311DirtyStorageKey: '' }),
  methods: {
    c311MarkDirty (value = true) { this.c311Dirty = value },
    c311ReadDirtyDraft () { return null },
    c311SaveDirtyDraft () {},
    c311ClearDirtyDraft () { this.c311Dirty = false },
  },
}

function formatDate (value) {
  try { return typeof C311JS.formatC311DateTime === 'function' ? C311JS.formatC311DateTime(value) : value || '' } catch (_error) { return value || '' }
}

function clone (value) {
  return value && typeof value === 'object' ? JSON.parse(JSON.stringify(value)) : value
}

function stableSerialize (value) {
  if (Array.isArray(value)) return `[${value.map(stableSerialize).join(',')}]`
  if (value && typeof value === 'object') return `{${Object.keys(value).sort((left, right) => left.localeCompare(right)).map(key => `${JSON.stringify(key)}:${stableSerialize(value[key])}`).join(',')}}`
  return JSON.stringify(value)
}

export default {
  name: 'C311Portal',
  components: { C311AppShell, C311DataState, C311ErrorSummary, C311CapabilityAction, C311HelpDrawer, C311LanguageSelector, C311MainNav, C311ResponsiveData },
  mixins: [c311DirtyGuard],
  data: () => ({
    state: 'loading',
    items: [],
    messages: {},
    dataError: null,
    statusMessage: '',
    successMessage: '',
    formErrors: [],
    submitting: false,
    draftSaving: false,
    draftID: '',
    draftVersion: null,
    submissionResult: null,
    idempotencyKey: '',
    idempotencyFingerprint: '',
    idempotencySerial: 0,
    attachmentTokenInput: '',
    attachmentUploading: false,
    attachmentRecoveryMessage: '',
    loadGeneration: 0,
    routeSignature: '',
    suppressRouteReload: false,
    customFieldsText: '{}',
    profile: null,
    versionConflict: null,
    conflictDraft: null,
    serviceTypes,
    form: { service_type: 'GENERAL_INQUIRY', summary: '', description: '', requester: { display_name: '', email: '', phone: '' }, location: { address: '', latitude: null, longitude: null }, attachment_tokens: [], custom_fields: {}, consent: false },
    columns: [
      { key: 'request_number', labelKey: 'field.requestNumber' },
      { key: 'summary', labelKey: 'field.summary' },
      { key: 'status', labelKey: 'field.status' },
      { key: 'updated_at', labelKey: 'field.updated', format: value => formatDate(value) },
    ],
  }),
  computed: {
    provider () { return this.$C311?.provider },
    actorID () { return this.$C311?.session?.actor?.actor_id || '' },
    isAuthenticated () { return !!this.$C311?.session?.authenticated },
    isSubmitRoute () { return this.$route?.name === 'c311.submit' },
    isStaffAssistRoute () { return this.$route?.name === 'c311.staff.submit' },
    isRequestFormRoute () { return this.isSubmitRoute || this.isStaffAssistRoute },
    showRequestList () { return this.$route?.name === 'c311.status' },
    canCapability () {
      return capability => {
        if (typeof this.$C311?.can === 'function') return this.$C311.can(capability)
        const capabilities = this.$C311?.session?.actor?.capabilities
        return !Array.isArray(capabilities) || capabilities.includes(capability)
      }
    },
    canCreateDraft () { return this.isAuthenticated && this.canCapability('portal_draft_create') },
    canUpdateDraft () { return !!this.draftID && this.canCapability('portal_draft_update') },
    canDeleteDraft () { return !!this.draftID && this.canCapability('portal_draft_delete') },
    canSubmitCurrentRequest () { return !this.draftID || this.canCapability('portal_draft_submit') },
    shellTitle () { return this.isStaffAssistRoute ? this.t('portal.submit.staffTitle', 'Submit a request for a resident') : this.isRequestFormRoute ? this.t('portal.submit.title', 'Submit a service request') : this.t('portal.title', 'City 311') },
    navItems () {
      const items = [
        { route: '/c311', label: this.t('navigation.home', 'Home') },
        { route: '/c311/services', label: this.t('navigation.services', 'Services') },
        { route: '/c311/help', label: this.t('navigation.help', 'Help') },
        { route: '/c311/submit', label: this.t('navigation.submit', 'Submit a request') },
        { route: '/c311/status', label: this.t('navigation.status', 'Check status') },
      ]
      if (this.isAuthenticated) items.push({ route: '/c311/requests', label: this.t('navigation.requests', 'My requests') }, { route: '/c311/account', label: this.t('navigation.account', 'Account') }, { route: '/c311/logout/callback', label: this.t('navigation.signOut', 'Sign out') })
      else items.push({ route: '/c311/sign-in', label: this.t('navigation.signIn', 'Sign in') }, { route: '/c311/register', label: this.t('navigation.register', 'Register') })
      return items
    },
    translatedColumns () { return this.columns.map(column => ({ ...column, label: this.t(column.labelKey, column.labelKey) })) },
  },
  watch: {
    form: { deep: true, handler () { if (this.isRequestFormRoute && this.c311Dirty) this.c311SaveDirtyDraft(this.draftStorageValue()) } },
    $route: {
      deep: true,
      handler (to) {
        const signature = this.routeKey(to)
        if (!signature || signature === this.routeSignature) return
        this.routeSignature = signature
        if (this.suppressRouteReload) {
          this.suppressRouteReload = false
          return
        }
        this.resetForRoute()
        this.load()
      },
    },
  },
  created () {
    this.routeSignature = this.routeKey(this.$route)
    this.configureRouteDraft()
    this.load()
  },
  methods: {
    t (key, fallback, params = {}) {
      const translated = this.$t?.(`c311:${key}`)
      const value = translated && translated !== `c311:${key}` && translated !== key ? translated : fallback
      return Object.entries(params).reduce((message, [name, replacement]) => message.replace(`{{${name}}}`, String(replacement)), value)
    },
    hasError (field) { return this.formErrors.some(error => error.field === field || error.field === `/${field}` || error.field?.replace(/^\//, '').replace(/\//g, '.') === field) },
    routeKey (route) {
      if (!route) return ''
      return `${route.name || ''}|${route.path || ''}|${JSON.stringify(route.query || {})}`
    },
    configureRouteDraft () {
      this.c311DirtyStorageKey = this.isRequestFormRoute ? `c311.portal.${this.isStaffAssistRoute ? 'staff-submit' : 'submit'}` : ''
      if (!this.isRequestFormRoute) return
      const draft = this.c311ReadDirtyDraft()
      if (draft && typeof draft === 'object') this.restoreLocalDraft(draft)
    },
    resetForRoute () {
      this.loadGeneration += 1
      this.state = 'loading'
      this.items = []
      this.messages = {}
      this.dataError = null
      this.statusMessage = ''
      this.successMessage = ''
      this.formErrors = []
      this.submitting = false
      this.draftSaving = false
      this.draftID = ''
      this.draftVersion = null
      this.submissionResult = null
      this.idempotencyKey = ''
      this.idempotencyFingerprint = ''
      this.attachmentTokenInput = ''
      this.attachmentUploading = false
      this.attachmentRecoveryMessage = ''
      this.profile = null
      this.versionConflict = null
      this.conflictDraft = null
      this.form = { service_type: 'GENERAL_INQUIRY', summary: '', description: '', requester: { display_name: '', email: '', phone: '' }, location: { address: '', latitude: null, longitude: null }, attachment_tokens: [], custom_fields: {}, consent: false }
      this.customFieldsText = '{}'
      this.c311MarkDirty?.(false)
      this.configureRouteDraft()
    },
    markDirty () { this.c311MarkDirty?.(true) },
    draftStorageValue () {
      return clone({
        service_type: this.form.service_type,
        summary: this.form.summary,
        description: this.form.description,
        requester: this.form.requester,
        location: this.form.location,
        ...(this.form.attachment_tokens.length ? { attachment_count: this.form.attachment_tokens.length } : {}),
        custom_fields: this.form.custom_fields,
        consent: this.form.consent,
      })
    },
    restoreLocalDraft (draft) {
      const source = draft && typeof draft === 'object' ? draft : {}
      const attachmentCount = Number(source.attachment_count) || (Array.isArray(source.attachment_tokens) ? source.attachment_tokens.length : 0)
      const safeDraft = Object.keys(source).reduce((result, key) => {
        if (key !== 'attachment_count' && key !== 'attachment_tokens') result[key] = source[key]
        return result
      }, {})
      this.form = { ...this.form, ...safeDraft, requester: { ...this.form.requester, ...(source.requester || {}) }, location: { ...this.form.location, ...(source.location || {}) }, attachment_tokens: [], custom_fields: source.custom_fields && typeof source.custom_fields === 'object' ? source.custom_fields : {} }
      this.customFieldsText = JSON.stringify(this.form.custom_fields, null, 2)
      this.attachmentRecoveryMessage = attachmentCount > 0 ? this.t('status.attachmentNeedsReupload', 'Previously selected attachments need to be uploaded again after this refresh.') : ''
    },
    hydrateRemoteDraft (draft) {
      const requester = draft?.primary_requester
      const location = draft?.location
      this.form = { ...this.form, service_type: draft?.service_type || this.form.service_type, summary: draft?.summary || '', description: draft?.description || '', requester: { display_name: requester?.display_name || '', email: requester?.emails?.[0] || '', phone: requester?.phone_numbers?.[0]?.value || '' }, location: { address: location?.address?.line1 || '', latitude: location?.latitude ?? null, longitude: location?.longitude ?? null }, custom_fields: draft?.custom_fields || {} }
      this.customFieldsText = JSON.stringify(this.form.custom_fields, null, 2)
      this.draftID = draft?.request_id || this.draftID
      this.draftVersion = typeof draft?.version === 'number' ? draft.version : this.draftVersion
    },
    isCurrentLoad (generation) { return generation === this.loadGeneration },
    async load () {
      const generation = ++this.loadGeneration
      this.state = 'loading'; this.dataError = null; this.formErrors = []; this.successMessage = ''; this.submissionResult = null
      if (this.isRequestFormRoute) {
        await this.loadRequestForm(generation)
        return
      }
      if (this.showRequestList) {
        await this.loadRequestList(generation)
        return
      }
      if (this.isCurrentLoad(generation)) this.state = 'populated'
    },
    setDataError (error) {
      this.dataError = error
      this.state = c311StateForError?.(error) || (error?.retryable ? 'retryable-error' : 'terminal-error')
    },
    async loadRequestForm (generation) {
      this.state = 'populated'
      if (this.isStaffAssistRoute) return
      if (!await this.loadProfile(generation) || !this.isCurrentLoad(generation)) return
      await this.loadRemoteDraft(generation)
    },
    async loadProfile (generation) {
      if (!this.isAuthenticated || !this.provider?.getProfile) return true
      try {
        const profile = await this.provider.getProfile()
        if (!this.isCurrentLoad(generation)) return false
        this.profile = profile
        if (!this.form.requester.display_name && this.profile) {
          this.form.requester.display_name = this.profile.display_name || ''
          this.form.requester.email = this.profile.emails?.[0] || ''
          this.form.requester.phone = this.profile.phone_numbers?.[0]?.value || ''
        }
        return true
      } catch (error) {
        if (!this.isCurrentLoad(generation)) return false
        this.setDataError(error)
        return false
      }
    },
    async loadRemoteDraft (generation) {
      const remoteDraftID = this.$route?.query?.draft_id
      if (!remoteDraftID || !this.isAuthenticated || !this.provider?.getDraft) return
      try {
        const draft = await this.provider.getDraft(remoteDraftID)
        if (this.isCurrentLoad(generation)) this.hydrateRemoteDraft(draft)
      } catch (error) {
        if (this.isCurrentLoad(generation)) this.setDataError(error)
      }
    },
    async loadRequestList (generation) {
      try {
        const page = await this.provider?.listPortalRequests?.()
        if (!this.isCurrentLoad(generation)) return
        this.items = page?.items || []
        this.state = this.items.length ? 'populated' : 'empty'
        this.statusMessage = this.items.length ? this.t('status.requestsLoaded', 'Requests loaded.') : this.t('status.noRequests', 'No requests found.')
      } catch (error) {
        if (!this.isCurrentLoad(generation)) return
        this.setDataError(error)
        this.statusMessage = this.t('status.requestListUnavailable', 'Request list unavailable.')
      }
    },
    addValidationError (errors, field, code, message) { errors.push({ field, code, message }) },
    validateBasics (errors) {
      const add = (field, code, message) => this.addValidationError(errors, field, code, message)
      if (!serviceTypes.includes(this.form.service_type)) add('service_type', 'INVALID_VALUE', this.t('error.serviceType', 'Choose a valid service type.'))
      if (!this.form.summary) add('summary', 'REQUIRED', this.t('error.summaryRequired', 'Summary is required.'))
      else if (this.form.summary.length < 5) add('summary', 'TOO_SHORT', this.t('error.summaryShort', 'Summary must be at least 5 characters.'))
      else if (this.form.summary.length > 160) add('summary', 'TOO_LONG', this.t('error.summaryLong', 'Summary must be 160 characters or fewer.'))
      if (!this.form.description) add('description', 'REQUIRED', this.t('error.descriptionRequired', 'Description is required.'))
      else if (this.form.description.length < 10) add('description', 'TOO_SHORT', this.t('error.descriptionShort', 'Description must be at least 10 characters.'))
      else if (this.form.description.length > 5000) add('description', 'TOO_LONG', this.t('error.descriptionLong', 'Description must be 5000 characters or fewer.'))
      else if (this.form.description.includes('<') && this.form.description.includes('>')) add('description', 'INVALID_FORMAT', this.t('error.descriptionPlainText', 'Description must be plain text.'))
    },
    validateRequester (errors) {
      const add = (field, code, message) => this.addValidationError(errors, field, code, message)
      const requester = this.form.requester
      if (!requester.display_name) add('requester.display_name', 'REQUIRED', this.t('error.requesterName', 'Requester name is required.'))
      else if (requester.display_name.length > 120) add('requester.display_name', 'TOO_LONG', this.t('error.requesterNameLong', 'Requester name must be 120 characters or fewer.'))
      if (!requester.email || !/^[^@\s]+@[^@\s]+\.[^@\s]+$/.test(requester.email)) add('requester.email', 'INVALID_FORMAT', this.t('error.email', 'Enter a valid email address.'))
      if (requester.phone && !/^\+[1-9]\d{1,14}$/.test(requester.phone)) add('requester.phone', 'INVALID_FORMAT', this.t('error.phone', 'Enter a phone number in international format.'))
    },
    validateLocation (errors) {
      const add = (field, code, message) => this.addValidationError(errors, field, code, message)
      const rule = serviceTypeRules[this.form.service_type]
      const location = this.form.location
      if (rule?.location_required && !location.address) add('location.address', 'LOCATION_REQUIRED', this.t('error.locationRequired', 'A location is required for this service type.'))
      const missingCoordinates = location.latitude === null || location.latitude === '' || location.longitude === null || location.longitude === ''
      if (rule?.confirmed_coordinates_required && missingCoordinates) add('location.latitude', 'COORDINATES_REQUIRED', this.t('error.coordinatesRequired', 'Latitude and longitude are required for this service type.'))
      this.validateCoordinate(errors, 'latitude', location.latitude, -90, 90, add)
      this.validateCoordinate(errors, 'longitude', location.longitude, -180, 180, add)
    },
    validateCoordinate (errors, field, value, minimum, maximum, add) {
      if (value === null || value === '') return
      const coordinate = Number(value)
      const label = field.charAt(0).toUpperCase() + field.slice(1)
      if (!Number.isFinite(coordinate)) add(`location.${field}`, 'INVALID_FORMAT', this.t(`error.${field}Format`, `${label} must be a number.`))
      else if (coordinate < minimum || coordinate > maximum) add(`location.${field}`, 'OUT_OF_RANGE', this.t(`error.${field}Range`, `${label} must be between ${minimum} and ${maximum}.`))
    },
    validateAttachments (errors) {
      if (this.form.attachment_tokens.length > 5) this.addValidationError(errors, 'attachment_tokens', 'TOO_MANY_ITEMS', this.t('error.attachmentsMax', 'You can add up to five attachments.'))
      if (!this.form.consent) this.addValidationError(errors, 'consent', 'REQUIRED', this.t('error.consentRequired', 'Please confirm the information is accurate.'))
    },
    validateCustomFields (errors) {
      try {
        const customFields = this.customFieldsText.trim() === '{}' && Object.keys(this.form.custom_fields).length ? this.form.custom_fields : JSON.parse(this.customFieldsText || '{}')
        if (!customFields || typeof customFields !== 'object' || Array.isArray(customFields)) throw new Error('invalid')
        this.form.custom_fields = customFields
      } catch (_error) {
        this.addValidationError(errors, 'custom_fields', 'INVALID_FORMAT', this.t('error.customFields', 'Custom fields must be a JSON object.'))
      }
    },
    validate () {
      const errors = []
      this.validateBasics(errors)
      this.validateRequester(errors)
      this.validateLocation(errors)
      this.validateAttachments(errors)
      this.validateCustomFields(errors)
      return errors
    },
    payload () {
      const payload = { summary: this.form.summary, description: this.form.description, service_type: this.form.service_type, requester: { display_name: this.form.requester.display_name, email: this.form.requester.email, ...(this.form.requester.phone ? { phone: this.form.requester.phone } : {}) }, ...(this.form.location.address ? { location: { address: this.form.location.address, ...(this.form.location.latitude !== null && this.form.location.latitude !== '' ? { latitude: Number(this.form.location.latitude) } : {}), ...(this.form.location.longitude !== null && this.form.location.longitude !== '' ? { longitude: Number(this.form.location.longitude) } : {}) } } : {}), ...(this.form.attachment_tokens.length ? { attachment_tokens: this.form.attachment_tokens.slice(0, 5) } : {}), ...(Object.keys(this.form.custom_fields).length ? { custom_fields: clone(this.form.custom_fields) } : {}) }
      return payload
    },
    draftPayload () {
      const payload = this.payload()
      if (!payload.summary) delete payload.summary
      if (!payload.description) delete payload.description
      if (!payload.requester.display_name && !payload.requester.email && !payload.requester.phone) delete payload.requester
      if (!payload.location) delete payload.location
      return payload
    },
    fingerprint (value) { return stableSerialize(value) },
    nextIdempotencyKey () {
      this.idempotencySerial += 1
      if (typeof crypto !== 'undefined' && typeof crypto.randomUUID === 'function') return crypto.randomUUID()
      if (typeof crypto !== 'undefined' && typeof crypto.getRandomValues === 'function') {
        const bytes = new Uint8Array(16)
        crypto.getRandomValues(bytes)
        return `c311-submit-${Array.from(bytes, byte => byte.toString(16).padStart(2, '0')).join('')}`
      }
      return `c311-submit-${Date.now()}-${this.idempotencySerial}`
    },
    setError (error) {
      this.dataError = error
      this.state = c311StateForError?.(error) || (error?.retryable ? 'retryable-error' : 'terminal-error')
      const code = error?.code || error?.error
      if (code === 'VERSION_CONFLICT' && !this.conflictDraft) {
        this.conflictDraft = { form: clone(this.form), customFieldsText: this.customFieldsText }
        this.versionConflict = { current_version: error?.current_version ?? error?.currentVersion ?? null }
      }
      const fieldErrors = error?.fieldErrors || error?.errors || []
      this.formErrors = fieldErrors.length ? fieldErrors.map(item => ({ ...item, message: item.message || this.t(`error.${item.code}`, item.code || this.t('error.operationFailed', 'The operation could not be completed.')) })) : [{ field: 'form', code: error?.code || error?.error || 'OPERATION_FAILED', message: error?.message || this.t('error.submitFailed', 'The request could not be submitted.') }]
    },
    async reloadDraft () {
      if (!this.draftID || !this.provider?.getDraft || this.draftSaving) return
      const generation = this.loadGeneration
      this.draftSaving = true
      try {
        const draft = await this.provider.getDraft(this.draftID)
        if (generation !== this.loadGeneration) return
        this.hydrateRemoteDraft(draft)
        this.versionConflict = null
        this.conflictDraft = null
        this.formErrors = []
        this.dataError = null
        this.successMessage = this.t('status.draftReloaded', 'Server draft reloaded.')
        this.c311MarkDirty?.(false)
      } catch (error) {
        if (generation !== this.loadGeneration) return
        this.setError(error)
      } finally {
        if (generation === this.loadGeneration) this.draftSaving = false
      }
    },
    reapplyDraft () {
      if (!this.conflictDraft) return
      this.form = clone(this.conflictDraft.form)
      this.customFieldsText = this.conflictDraft.customFieldsText
      this.versionConflict = null
      this.conflictDraft = null
      this.formErrors = []
      this.dataError = null
      this.successMessage = ''
      this.c311MarkDirty?.(true)
    },
    clearSubmissionState () {
      this.formErrors = []
      this.dataError = null
      this.successMessage = ''
    },
    validateSubmissionPermission () {
      if (!this.draftID || this.canSubmitCurrentRequest) return true
      this.setError({ status: 403, error: 'FORBIDDEN', message: this.t('error.draftSubmitDenied', 'You do not have permission to submit this draft.') })
      return false
    },
    async submitStaffRequest (request) {
      const constituent = this.profile?.constituent_id ? { constituent_id: this.profile.constituent_id } : { display_name: request.requester.display_name, email: request.requester.email }
      const detail = await this.provider?.createStaffServiceRequest?.({ constituent, request })
      return detail?.request || detail
    },
    async submitPortalRequest (request) {
      const fingerprint = this.fingerprint(request)
      if (!this.idempotencyKey || this.idempotencyFingerprint !== fingerprint) {
        this.idempotencyKey = this.nextIdempotencyKey()
        this.idempotencyFingerprint = fingerprint
      }
      return this.provider?.submitPortalRequest?.(request, { idempotencyKey: this.idempotencyKey })
    },
    async submitRequest () {
      const request = this.payload()
      if (this.isStaffAssistRoute) return this.submitStaffRequest(request)
      if (this.draftID) return this.provider?.submitDraft?.(this.draftID, { expectedVersion: this.draftVersion })
      return this.submitPortalRequest(request)
    },
    async applySubmissionSuccess (response) {
      if (!response) throw new Error(this.t('error.providerUnavailable', 'The request provider is unavailable.'))
      this.submissionResult = { request_number: response.request_number, status: response.status || 'SUBMITTED', version: response.version }
      this.state = 'populated'
      this.successMessage = this.t('status.submitted', `Request ${response.request_number} submitted.`).replace('{{number}}', response.request_number)
      this.statusMessage = this.successMessage
      this.c311ClearDirtyDraft?.()
      this.versionConflict = null
      this.conflictDraft = null
      if (this.draftID) {
        this.draftID = ''; this.draftVersion = null
        if (this.$router?.replace && this.$route?.query?.draft_id) {
          const query = { ...this.$route.query }
          delete query.draft_id
          await this.$router.replace({ query })
        }
      }
    },
    async submit () {
      if (this.submitting || this.submissionResult) return
      this.clearSubmissionState()
      if (!this.validateSubmissionPermission()) return
      const errors = this.validate()
      if (errors.length) {
        this.formErrors = errors
        this.state = 'validation-error'
        return
      }
      const generation = this.loadGeneration
      this.submitting = true
      try {
        const response = await this.submitRequest()
        if (!this.isCurrentLoad(generation)) return
        await this.applySubmissionSuccess(response)
      } catch (error) {
        if (this.isCurrentLoad(generation)) this.setError(error)
      } finally {
        if (this.isCurrentLoad(generation)) this.submitting = false
      }
    },
    async saveDraft () {
      if (!this.isAuthenticated || this.draftSaving || (!this.draftID && !this.canCreateDraft) || (this.draftID && !this.canUpdateDraft)) return
      const generation = this.loadGeneration
      this.formErrors = []; this.dataError = null; this.draftSaving = true
      try {
        const payload = this.draftPayload(); const response = this.draftID ? await this.provider?.updateDraft?.(this.draftID, payload, { expectedVersion: this.draftVersion }) : await this.provider?.createDraft?.(payload)
        if (generation !== this.loadGeneration) return
        this.hydrateRemoteDraft(response); this.c311MarkDirty?.(false); this.successMessage = this.t('status.draftSaved', 'Draft saved.'); this.versionConflict = null; this.conflictDraft = null
        if (this.$router?.replace && this.$route?.query && this.$route.query.draft_id !== this.draftID) {
          this.suppressRouteReload = true
          try {
            await this.$router.replace({ query: { ...this.$route.query, draft_id: this.draftID } })
          } finally {
            this.suppressRouteReload = false
          }
        }
      } catch (error) {
        if (generation !== this.loadGeneration) return
        this.setError(error)
      } finally {
        if (generation === this.loadGeneration) this.draftSaving = false
      }
    },
    async deleteDraft () {
      if (!this.draftID || this.draftSaving || !this.canDeleteDraft) return
      const generation = this.loadGeneration
      this.draftSaving = true
      try {
        await this.provider?.deleteDraft?.(this.draftID, { expectedVersion: this.draftVersion })
        if (generation !== this.loadGeneration) return
        this.draftID = ''; this.draftVersion = null; this.c311ClearDirtyDraft?.(); this.successMessage = this.t('status.draftDeleted', 'Draft deleted.')
        if (this.$router?.replace && this.$route?.query?.draft_id) { const query = { ...this.$route.query }; delete query.draft_id; await this.$router.replace({ query }) }
      } catch (error) {
        if (generation !== this.loadGeneration) return
        this.setError(error)
      } finally {
        if (generation === this.loadGeneration) this.draftSaving = false
      }
    },
    async uploadAttachment (event) {
      const file = event?.target?.files?.[0]
      if (!file) return
      const mediaTypes = ['image/jpeg', 'image/png', 'application/pdf', 'text/plain', 'application/vnd.openxmlformats-officedocument.wordprocessingml.document']
      const invalid = file.size > 10 * 1024 * 1024 || String(file.name || '').length < 1 || String(file.name || '').length > 120 || !mediaTypes.includes(file.type)
      if (invalid) {
        this.formErrors = [{ field: 'attachment_tokens', code: 'INVALID_FORMAT', message: this.t('error.attachmentInvalid', 'Choose a supported attachment up to 10 MB.') }]
        event.target.value = ''
        return
      }
      if (this.form.attachment_tokens.length >= 5 || !this.provider?.uploadPortalAttachment) return
      const generation = this.loadGeneration
      this.attachmentUploading = true
      this.formErrors = []
      try {
        const attachment = await this.provider.uploadPortalAttachment({ file, filename: file.name, media_type: file.type })
        if (generation !== this.loadGeneration) return
        if (attachment?.attachment_token && !this.form.attachment_tokens.includes(attachment.attachment_token)) this.form.attachment_tokens.push(attachment.attachment_token)
        this.statusMessage = this.t('status.attachmentUploaded', 'Attachment uploaded.')
        this.attachmentRecoveryMessage = ''
        this.markDirty()
      } catch (error) {
        if (generation !== this.loadGeneration) return
        this.setError(error)
      } finally {
        if (generation === this.loadGeneration) this.attachmentUploading = false
        event.target.value = ''
      }
    },
    addAttachmentToken () {
      const token = this.attachmentTokenInput.trim(); if (!token) return
      if (this.form.attachment_tokens.length >= 5) { this.formErrors = [{ field: 'attachment_tokens', code: 'TOO_MANY_ITEMS', message: this.t('error.attachmentsMax', 'You can add up to five attachments.') }]; return }
      if (!this.form.attachment_tokens.includes(token)) this.form.attachment_tokens.push(token); this.attachmentTokenInput = ''; this.attachmentRecoveryMessage = ''; this.markDirty()
    },
    removeAttachmentToken (token) { this.form.attachment_tokens = this.form.attachment_tokens.filter(item => item !== token); this.markDirty() },
    onCustomFieldsInput () {
      try { const value = JSON.parse(this.customFieldsText || '{}'); if (value && typeof value === 'object' && !Array.isArray(value)) this.form.custom_fields = value } catch (_error) {}
      this.markDirty()
    },
  },
}
</script>

<style scoped>
.c311-request-page { max-width: 56rem; }
.c311-request-section { border: 0; padding: 0; margin: 1.5rem 0; }
.c311-request-section legend { font-size: 1.1rem; font-weight: 600; }
.c311-request-grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 1rem; }
@media (max-width: 767px) { .c311-request-grid { grid-template-columns: 1fr; gap: 0; } }
</style>
