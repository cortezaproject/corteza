<template>
  <c311-app-shell mode="public" :brand="brandName" :title="pageTitle" :status-message="statusMessage">
    <template #nav>
      <div class="d-flex flex-wrap align-items-center gap-2">
        <c311-main-nav :items="navItems" :label="t('navigation.public', 'Public navigation')" />
        <c311-help-drawer
          help-key="public.request.submit"
          :label="t('help.label', 'Help')"
          :title="t('help.title', 'Help')"
          :close-label="t('action.close', 'Close')"
          :content="helpBody"
        />
        <c311-language-selector :actor-id="actorID" :value="forms.account.preferred_language.toLowerCase()" @change="languageChanged" />
      </div>
    </template>

    <div v-if="branding" class="c311-branding" data-c311-branding :style="brandStyle">
      <img v-if="safeLogoUrl" :src="safeLogoUrl" :alt="brandName" class="c311-branding__logo">
      <p v-if="identityShell && branding.login_header" class="c311-branding__header" data-c311-branding-login-header>{{ branding.login_header }}</p>
      <p v-else-if="branding.public_header" class="c311-branding__header">{{ branding.public_header }}</p>
    </div>

    <c311-data-state v-if="pageNeedsContent && contentState !== 'populated'" :state="contentState" :error="dataError" @retry="load" />

    <section v-if="page === 'home' && state === 'populated'" data-c311-page="home" :data-c311-content-key="contentKey" class="c311-public-page">
      <div v-if="contentBody" v-html="contentBody" />
      <p v-else>{{ t('portal.home.empty', 'Public information is not available yet.') }}</p>
      <div class="d-flex flex-wrap gap-2 mt-4">
        <router-link class="btn btn-primary" to="/c311/submit" data-c311-action="submit-request">{{ t('action.submit', 'Submit request') }}</router-link>
        <router-link class="btn btn-outline-primary" to="/c311/services">{{ t('navigation.services', 'Services') }}</router-link>
      </div>
    </section>

    <section v-else-if="page === 'services' && state === 'populated'" data-c311-page="services" :data-c311-content-key="contentKey" class="c311-public-page">
      <div v-if="contentBody" v-html="contentBody" />
      <p v-else>{{ t('portal.services.empty', 'No services are currently listed.') }}</p>
    </section>

    <section v-else-if="page === 'help' && state === 'populated'" data-c311-page="help" :data-c311-content-key="contentKey" class="c311-public-page">
      <div v-if="contentBody" v-html="contentBody" />
      <p v-if="!contentBody">{{ t('portal.help.empty', 'Help content is not available yet.') }}</p>
      <c311-data-state v-if="helpState !== 'populated'" :state="helpState" :error="helpError" @retry="load" />
      <p v-else-if="helpBody" v-html="helpBody" />
    </section>

    <section v-else-if="page === 'sign-in'" data-c311-page="sign-in" class="c311-public-page">
      <c311-error-summary :errors="formErrors" :field-targets="fieldTargets" :title="t('error.review', 'Review your information')" />
      <form @submit.prevent="signIn">
        <div class="form-group"><label for="c311-login-identifier">{{ t('field.loginIdentifier', 'Email or username') }}</label><input id="c311-login-identifier" v-model.trim="forms.signIn.login_identifier" class="form-control" autocomplete="username" aria-describedby="c311-error-summary" aria-errormessage="c311-error-summary" :aria-invalid="hasError('login_identifier') ? 'true' : 'false'" @input="markDirty"></div>
        <div class="form-group"><label for="c311-login-password">{{ t('field.password', 'Password') }}</label><input id="c311-login-password" v-model="forms.signIn.password" class="form-control" type="password" autocomplete="current-password" aria-describedby="c311-error-summary" aria-errormessage="c311-error-summary" :aria-invalid="hasError('password') ? 'true' : 'false'"></div>
        <button class="btn btn-primary" type="submit" data-c311-action="sign-in" :disabled="busy.signIn">{{ busy.signIn ? t('action.working', 'Working…') : t('action.signIn', 'Sign in') }}</button>
      </form>
      <div class="mt-3 d-flex flex-column align-items-start gap-2">
        <button class="btn btn-link p-0" type="button" data-c311-action="oidc-sign-in" :disabled="busy.federated" @click="federated('oidc')">{{ busy.federated ? t('action.working', 'Working…') : t('action.signInOidc', 'Continue with OIDC') }}</button>
        <router-link to="/c311/register">{{ t('navigation.register', 'Create an account') }}</router-link>
        <router-link to="/c311/forgot-password">{{ t('navigation.forgotPassword', 'Forgot password?') }}</router-link>
      </div>
      <p v-if="federatedMessage" class="alert alert-info mt-3" role="status">{{ federatedMessage }}</p>
    </section>

    <section v-else-if="page === 'register'" data-c311-page="register" class="c311-public-page">
      <c311-error-summary :errors="formErrors" :field-targets="fieldTargets" :title="t('error.review', 'Review your information')" />
      <form @submit.prevent="register">
        <div class="form-group"><label for="c311-register-name">{{ t('field.displayName', 'Name') }}</label><input id="c311-register-name" v-model.trim="forms.register.display_name" class="form-control" autocomplete="name" aria-describedby="c311-error-summary" aria-errormessage="c311-error-summary" :aria-invalid="hasError('display_name') ? 'true' : 'false'" @input="markDirty"></div>
        <div class="form-group"><label for="c311-register-email">{{ t('field.email', 'Email') }}</label><input id="c311-register-email" v-model.trim="forms.register.email" class="form-control" type="email" autocomplete="email" aria-describedby="c311-error-summary" aria-errormessage="c311-error-summary" :aria-invalid="hasError('email') ? 'true' : 'false'" @input="markDirty"></div>
        <div class="form-group"><label for="c311-register-login">{{ t('field.loginIdentifier', 'Username') }}</label><input id="c311-register-login" v-model.trim="forms.register.login_identifier" class="form-control" autocomplete="username" aria-describedby="c311-error-summary" aria-errormessage="c311-error-summary" :aria-invalid="hasError('login_identifier') ? 'true' : 'false'" @input="markDirty"></div>
        <div class="form-group"><label for="c311-register-password">{{ t('field.password', 'Password') }}</label><input id="c311-register-password" v-model="forms.register.password" class="form-control" type="password" autocomplete="new-password" aria-describedby="c311-error-summary" aria-errormessage="c311-error-summary" :aria-invalid="hasError('password') ? 'true' : 'false'"><small class="form-text text-muted">{{ passwordRuleText }}</small></div>
        <button class="btn btn-primary" type="submit" data-c311-action="register" :disabled="busy.register">{{ busy.register ? t('action.working', 'Working…') : t('action.register', 'Register') }}</button>
      </form>
      <p v-if="successMessage" class="alert alert-success mt-3" role="status">{{ successMessage }}</p>
    </section>

    <section v-else-if="page === 'forgot-password'" data-c311-page="forgot-password" class="c311-public-page">
      <c311-error-summary :errors="formErrors" :field-targets="fieldTargets" :title="t('error.review', 'Review your information')" />
      <form @submit.prevent="forgotPassword">
        <div class="form-group"><label for="c311-forgot-email">{{ t('field.email', 'Email') }}</label><input id="c311-forgot-email" v-model.trim="forms.forgot.email" class="form-control" type="email" autocomplete="email" aria-describedby="c311-error-summary" aria-errormessage="c311-error-summary" :aria-invalid="hasError('email') ? 'true' : 'false'" @input="markDirty"></div>
        <button class="btn btn-primary" type="submit" data-c311-action="forgot-password" :disabled="busy.forgot">{{ busy.forgot ? t('action.working', 'Working…') : t('action.sendReset', 'Send reset instructions') }}</button>
      </form>
      <p v-if="successMessage" class="alert alert-success mt-3" role="status">{{ successMessage }}</p>
    </section>

    <section v-else-if="page === 'reset-password'" data-c311-page="reset-password" class="c311-public-page">
      <c311-error-summary :errors="formErrors" :field-targets="fieldTargets" :title="t('error.review', 'Review your information')" />
      <form @submit.prevent="resetPassword">
        <div class="form-group"><label for="c311-reset-password">{{ t('field.password', 'New password') }}</label><input id="c311-reset-password" v-model="forms.reset.password" class="form-control" type="password" autocomplete="new-password" aria-describedby="c311-error-summary" aria-errormessage="c311-error-summary" :aria-invalid="hasError('password') ? 'true' : 'false'"><small class="form-text text-muted">{{ passwordRuleText }}</small></div>
        <button class="btn btn-primary" type="submit" data-c311-action="reset-password" :disabled="busy.reset">{{ busy.reset ? t('action.working', 'Working…') : t('action.resetPassword', 'Reset password') }}</button>
      </form>
      <p v-if="successMessage" class="alert alert-success mt-3" role="status">{{ successMessage }}</p>
    </section>

    <section v-else-if="page === 'account'" data-c311-page="account" class="c311-public-page">
      <c311-data-state v-if="state !== 'populated'" :state="state" :error="dataError" @retry="load" />
      <template v-else>
      <c311-error-summary :errors="formErrors" :field-targets="fieldTargets" :title="t('error.review', 'Review your information')" />
        <form @submit.prevent="updateAccount">
          <div class="form-group"><label for="c311-account-name">{{ t('field.displayName', 'Name') }}</label><input id="c311-account-name" v-model.trim="forms.account.display_name" class="form-control" autocomplete="name" aria-describedby="c311-error-summary" aria-errormessage="c311-error-summary" :aria-invalid="hasError('display_name') ? 'true' : 'false'" @input="markDirty"></div>
          <div class="form-group"><label for="c311-account-language">{{ t('field.language', 'Language') }}</label><select id="c311-account-language" v-model="forms.account.preferred_language" class="form-control" aria-describedby="c311-error-summary" aria-errormessage="c311-error-summary" :aria-invalid="hasError('preferred_language') ? 'true' : 'false'" @change="markDirty"><option value="EN">{{ t('language.english', 'English') }}</option><option value="ES">{{ t('language.spanish', 'Español') }}</option><option value="VI">{{ t('language.vietnamese', 'Tiếng Việt') }}</option></select></div>
          <fieldset class="form-group"><legend>{{ t('field.phones', 'Phone numbers') }}</legend><div v-for="(phone, index) in forms.account.phone_numbers" :key="`phone-${index}`" class="d-flex flex-wrap gap-2 mb-2"><label class="sr-only" :for="`c311-account-phone-${index}-label`">{{ t('field.phoneLabel', 'Phone label') }}</label><select :id="`c311-account-phone-${index}-label`" v-model="phone.label" class="form-control" :aria-invalid="hasError(`phone_numbers/${index}/label`) ? 'true' : 'false'" @change="markDirty"><option v-for="label in phoneLabels" :key="label" :value="label">{{ label }}</option></select><label class="sr-only" :for="`c311-account-phone-${index}-value`">{{ t('field.phoneNumber', 'Phone number') }}</label><input :id="`c311-account-phone-${index}-value`" v-model.trim="phone.value" class="form-control" autocomplete="tel" aria-describedby="c311-error-summary" aria-errormessage="c311-error-summary" :aria-invalid="hasError(`phone_numbers/${index}/value`) || hasError('phone_numbers') ? 'true' : 'false'" @input="markDirty"><button class="btn btn-outline-secondary" type="button" :aria-label="t('action.removePhone', 'Remove phone number')" @click="removePhone(index)">×</button></div><button class="btn btn-link p-0" type="button" :disabled="forms.account.phone_numbers.length >= 3" @click="addPhone">+ {{ t('action.addPhone', 'Add phone number') }}</button></fieldset>
          <fieldset class="form-group"><legend>{{ t('field.addresses', 'Addresses') }}</legend><div v-for="(address, index) in forms.account.addresses" :key="`address-${index}`" class="border rounded p-2 mb-2"><div class="d-flex justify-content-between align-items-center mb-2"><strong>{{ t('field.addressNumber', 'Address') }} {{ index + 1 }}</strong><button class="btn btn-outline-secondary" type="button" :aria-label="t('action.removeAddress', 'Remove address')" @click="removeAddress(index)">×</button></div><div class="form-group"><label :for="`c311-account-address-${index}-line1`">{{ t('field.addressLine1', 'Address line 1') }}</label><input :id="`c311-account-address-${index}-line1`" v-model.trim="address.line1" class="form-control" autocomplete="address-line1" aria-describedby="c311-error-summary" aria-errormessage="c311-error-summary" :aria-invalid="hasError(`addresses/${index}/line1`) || hasError('addresses') ? 'true' : 'false'" @input="markDirty"></div><div class="form-group"><label :for="`c311-account-address-${index}-line2`">{{ t('field.addressLine2', 'Address line 2') }}</label><input :id="`c311-account-address-${index}-line2`" v-model.trim="address.line2" class="form-control" autocomplete="address-line2" @input="markDirty"></div><div class="form-group"><label :for="`c311-account-address-${index}-city`">{{ t('field.city', 'City') }}</label><input :id="`c311-account-address-${index}-city`" v-model.trim="address.city" class="form-control" autocomplete="address-level2" aria-describedby="c311-error-summary" aria-errormessage="c311-error-summary" :aria-invalid="hasError(`addresses/${index}/city`) || hasError('addresses') ? 'true' : 'false'" @input="markDirty"></div><div class="form-group"><label :for="`c311-account-address-${index}-region`">{{ t('field.region', 'Region') }}</label><input :id="`c311-account-address-${index}-region`" v-model.trim="address.region" class="form-control" autocomplete="address-level1" @input="markDirty"></div><div class="form-group"><label :for="`c311-account-address-${index}-postal-code`">{{ t('field.postalCode', 'Postal code') }}</label><input :id="`c311-account-address-${index}-postal-code`" v-model.trim="address.postal_code" class="form-control" autocomplete="postal-code" aria-describedby="c311-error-summary" aria-errormessage="c311-error-summary" :aria-invalid="hasError(`addresses/${index}/postal_code`) || hasError('addresses') ? 'true' : 'false'" @input="markDirty"></div><div class="form-group"><label :for="`c311-account-address-${index}-country`">{{ t('field.country', 'Country') }}</label><input :id="`c311-account-address-${index}-country`" v-model.trim="address.country" class="form-control" autocomplete="country" maxlength="2" @input="markDirty"></div><label><input type="radio" name="c311-primary-address" :checked="address.primary" @change="setPrimaryAddress(index)"> {{ t('field.primaryAddress', 'Primary address') }}</label></div><button class="btn btn-link p-0" type="button" :disabled="forms.account.addresses.length >= 5" @click="addAddress">+ {{ t('action.addAddress', 'Add address') }}</button></fieldset>
          <div class="form-group"><label for="c311-account-category">{{ t('field.category', 'Primary category') }}</label><select id="c311-account-category" v-model="forms.account.primary_category" class="form-control" aria-describedby="c311-error-summary" aria-errormessage="c311-error-summary" :aria-invalid="hasError('primary_category') ? 'true' : 'false'" @change="markDirty"><option v-for="category in contactCategories" :key="category" :value="category">{{ category }}</option></select></div>
          <button class="btn btn-primary" type="submit" data-c311-action="update-account" :disabled="busy.account">{{ busy.account ? t('action.working', 'Working…') : t('action.save', 'Save changes') }}</button>
        </form>
        <hr>
        <form v-if="canChangeLoginIdentifier" @submit.prevent="changeLoginIdentifier">
          <h2>{{ t('account.loginIdentifier.title', 'Change login identifier') }}</h2>
          <div class="form-group"><label for="c311-account-login-identifier">{{ t('field.newLoginIdentifier', 'New login identifier') }}</label><input id="c311-account-login-identifier" v-model.trim="forms.account.login_identifier" class="form-control" autocomplete="username" aria-describedby="c311-error-summary" aria-errormessage="c311-error-summary" :aria-invalid="hasError('login_identifier') ? 'true' : 'false'" @input="markDirty"></div>
          <div class="form-group"><label for="c311-account-current-password">{{ t('field.currentPassword', 'Current password') }}</label><input id="c311-account-current-password" v-model="forms.account.current_password" class="form-control" type="password" autocomplete="current-password" aria-describedby="c311-error-summary" aria-errormessage="c311-error-summary" :aria-invalid="hasError('current_password') ? 'true' : 'false'"></div>
          <button class="btn btn-primary" type="submit" data-c311-action="change-login-identifier" :disabled="busy.loginIdentifier">{{ busy.loginIdentifier ? t('action.working', 'Working…') : t('action.changeLoginIdentifier', 'Save login identifier') }}</button>
        </form>
        <form v-if="canChangePassword" class="mt-4" @submit.prevent="changePassword">
          <h2>{{ t('account.password.title', 'Change password') }}</h2>
          <div class="form-group"><label for="c311-account-password-current-password">{{ t('field.currentPassword', 'Current password') }}</label><input id="c311-account-password-current-password" v-model="forms.account.current_password" class="form-control" type="password" autocomplete="current-password" aria-describedby="c311-error-summary" aria-errormessage="c311-error-summary" :aria-invalid="hasError('current_password') ? 'true' : 'false'"></div>
          <div class="form-group"><label for="c311-account-new-password">{{ t('field.newPassword', 'New password') }}</label><input id="c311-account-new-password" v-model="forms.account.new_password" class="form-control" type="password" autocomplete="new-password" aria-describedby="c311-error-summary" aria-errormessage="c311-error-summary" :aria-invalid="hasError('new_password') ? 'true' : 'false'"><small class="form-text text-muted">{{ passwordRuleText }}</small></div>
          <button class="btn btn-primary" type="submit" data-c311-action="change-password" :disabled="busy.password">{{ busy.password ? t('action.working', 'Working…') : t('action.changePassword', 'Save password') }}</button>
        </form>
        <p v-if="successMessage" class="alert alert-success mt-3" role="status">{{ successMessage }}</p>
      </template>
    </section>

    <section v-else-if="page === 'requests'" data-c311-page="requests" class="c311-public-page">
      <c311-data-state :state="state" :error="dataError" @retry="load">
        <template #populated><c311-responsive-data :items="items" :columns="requestColumns" row-key="request_id" :label="t('portal.myRequests', 'My requests')" /></template>
      </c311-data-state>
    </section>

    <section v-else-if="page === 'auth-callback'" data-c311-page="auth-callback" class="c311-public-page">
      <c311-error-summary :errors="formErrors" :field-targets="fieldTargets" :title="t('error.review', 'Review your information')" />
      <c311-data-state v-if="state !== 'populated'" :state="state" :error="dataError" @retry="load" />
      <p v-else role="status">{{ federatedMessage || t('identity.completing', 'Completing sign-in…') }}</p>
    </section>
    <section v-else-if="page === 'link-confirm'" data-c311-page="link-confirm" class="c311-public-page">
      <c311-error-summary :errors="formErrors" :field-targets="fieldTargets" :title="t('error.review', 'Review your information')" />
      <p v-if="linkState === 'pending'">{{ t('identity.linkConfirmationPrompt', 'Confirm linking this identity to your City 311 account before continuing.') }}</p>
      <p v-else-if="linkState === 'success'" class="alert alert-success" role="status">{{ federatedMessage }}</p>
      <c311-data-state v-else-if="linkState === 'error'" :state="dataErrorState" :error="dataError" @retry="confirmAccountLink" />
      <div v-if="linkState === 'pending' || linkState === 'error'" class="d-flex gap-2">
        <button class="btn btn-primary" type="button" data-c311-action="confirm-link" :disabled="busy.link" @click="confirmAccountLink">{{ busy.link ? t('action.working', 'Working…') : t('action.confirm', 'Confirm') }}</button>
        <button class="btn btn-secondary" type="button" data-c311-action="cancel-link-confirm" @click="cancelAccountLink">{{ t('action.cancel', 'Cancel') }}</button>
      </div>
    </section>
    <section v-else-if="page === 'logout-callback'" data-c311-page="logout-callback" class="c311-public-page">
      <c311-error-summary :errors="formErrors" :field-targets="fieldTargets" :title="t('error.review', 'Review your information')" />
      <c311-data-state v-if="state !== 'populated'" :state="state" :error="dataError" @retry="load" />
      <p v-else role="status">{{ successMessage || t('identity.loggedOut', 'You are signed out.') }}</p>
    </section>
    <footer v-if="branding && branding.public_footer" class="c311-branding__footer" data-c311-branding-footer>{{ branding.public_footer }}</footer>
  </c311-app-shell>
</template>

<script>
import * as C311JS from '@cortezaproject/corteza-js'
import { components, c311, c311Identity, mixins } from '@cortezaproject/corteza-vue'

const { C311AppShell, C311DataState, C311ErrorSummary, C311HelpDrawer, C311LanguageSelector, C311MainNav, C311ResponsiveData } = components
const stateForError = c311?.c311StateForError
const validatePassword = c311Identity?.validatePassword || (() => [])
const c311DirtyGuard = mixins?.c311DirtyGuard || {
  data: () => ({ c311Dirty: false, c311DirtyStorageKey: '' }),
  methods: {
    c311MarkDirty (value = true) { this.c311Dirty = value },
    c311ReadDirtyDraft () { return null },
    c311SaveDirtyDraft () {},
    c311ClearDirtyDraft () { this.c311Dirty = false },
  },
}
const FORM_BY_PAGE = { 'sign-in': 'signIn', 'forgot-password': 'forgot', 'reset-password': 'reset' }
const LOGIN_IDENTIFIER_PATTERN = /^[a-z0-9._-]{3,64}$/

function formatDate (value) {
  try { return typeof C311JS.formatC311DateTime === 'function' ? C311JS.formatC311DateTime(value) : value || '' } catch (_error) { return value || '' }
}

function displayError (error) {
  if (!error || typeof error !== 'object') return { message: String(error || '') }
  return {
    status: error.status,
    code: error.code || error.error,
    error: error.error || error.code,
    retryable: !!error.retryable,
    message: error.message,
    errors: error.fieldErrors || error.errors || [],
    fieldErrors: error.fieldErrors || error.errors || [],
    current_version: error.current_version ?? error.currentVersion,
    currentVersion: error.currentVersion ?? error.current_version,
  }
}

function safeColour (value) {
  return typeof value === 'string' && /^(#[0-9a-f]{3,8}|rgba?\([\d\s.,%]+\))$/i.test(value) ? value : ''
}

function safeFontFamily (value) {
  return typeof value === 'string' && /^[a-z0-9 ,_-]+$/i.test(value) ? value : ''
}

function safeResourceURL (value) {
  if (typeof value !== 'string') return ''
  if (/^https:\/\//i.test(value) || /^\/(?!\/)/.test(value)) return value
  return ''
}

export default {
  name: 'C311PublicPortal',
  components: { C311AppShell, C311DataState, C311ErrorSummary, C311HelpDrawer, C311LanguageSelector, C311MainNav, C311ResponsiveData },
  mixins: [c311DirtyGuard],
  data: () => ({
    state: 'loading',
    contentState: 'loading',
    helpState: 'loading',
    dataError: null,
    helpError: null,
    statusMessage: '',
    successMessage: '',
    federatedMessage: '',
    contentBody: '',
    contentKey: '',
    helpBody: '',
    branding: null,
    items: [],
    formErrors: [],
    linkState: 'idle',
    profile: null,
    resetToken: '',
    loadGeneration: 0,
    sessionRevision: 0,
    restoredDraftFields: [],
    activeAccountAction: '',
    busy: { signIn: false, register: false, forgot: false, reset: false, federated: false, account: false, loginIdentifier: false, password: false, link: false },
    forms: {
      signIn: { login_identifier: '', password: '' },
      register: { display_name: '', email: '', login_identifier: '', password: '', preferred_language: 'EN' },
      forgot: { email: '' },
      reset: { password: '' },
      account: { display_name: '', preferred_language: 'EN', login_identifier: '', current_password: '', new_password: '', phone_numbers: [], addresses: [], primary_category: 'RESIDENT' },
    },
  }),
  computed: {
    provider () { return this.$C311?.provider },
    actorID () { return this.$C311?.session?.actor?.actor_id || '' },
    isAuthenticated () { return this.sessionRevision >= 0 && !!this.$C311?.session?.authenticated },
    page () {
      const names = { 'c311.portal': 'home', 'c311.services': 'services', 'c311.help': 'help', 'c311.sign-in': 'sign-in', 'c311.register': 'register', 'c311.forgot-password': 'forgot-password', 'c311.reset-password': 'reset-password', 'c311.account': 'account', 'c311.requests': 'requests', 'c311.auth.callback': 'auth-callback', 'c311.auth.link.confirm': 'link-confirm', 'c311.logout.callback': 'logout-callback' }
      return names[this.$route?.name] || 'home'
    },
    pageNeedsContent () { return ['home', 'services', 'help'].includes(this.page) },
    identityShell () { return ['sign-in', 'register', 'forgot-password', 'reset-password', 'auth-callback', 'link-confirm'].includes(this.page) },
    phoneLabels () { return C311JS.PHONE_LABELS || ['MOBILE', 'HOME', 'WORK'] },
    contactCategories () { return C311JS.CONTACT_CATEGORIES || ['RESIDENT', 'BUSINESS', 'BUSINESS_OWNER', 'VETERAN', 'NEIGHBORHOOD_ASSOCIATION', 'GOVERNMENT', 'OTHER'] },
    fieldTargets () {
      const targets = {
        login_identifier: 'c311-login-identifier',
        password: 'c311-login-password',
        display_name: 'c311-register-name',
        email: 'c311-register-email',
        current_password: this.activeAccountAction === 'password' ? 'c311-account-password-current-password' : 'c311-account-current-password',
        new_password: 'c311-account-new-password',
        phone_numbers: 'c311-account-phone-0-value',
        addresses: 'c311-account-address-0-line1',
        primary_category: 'c311-account-category',
        form: 'c311-login-identifier',
        token: 'c311-reset-password',
      }
      if (this.page === 'register') targets.form = 'c311-register-name'
      if (this.page === 'forgot-password') targets.form = 'c311-forgot-email'
      if (this.page === 'reset-password') targets.form = 'c311-reset-password'
      if (this.page === 'account') targets.form = 'c311-account-name'
      if (this.page === 'account') targets.display_name = 'c311-account-name'
      if (this.page === 'reset-password') targets.password = 'c311-reset-password'
      if (this.page === 'forgot-password') targets.email = 'c311-forgot-email'
      for (let index = 0; index < 3; index++) {
        targets[`phone_numbers/${index}/label`] = `c311-account-phone-${index}-label`
        targets[`phone_numbers/${index}/value`] = `c311-account-phone-${index}-value`
      }
      for (let index = 0; index < 5; index++) {
        targets[`addresses/${index}/line1`] = `c311-account-address-${index}-line1`
        targets[`addresses/${index}/city`] = `c311-account-address-${index}-city`
        targets[`addresses/${index}/postal_code`] = `c311-account-address-${index}-postal-code`
      }
      return targets
    },
    canChangeLoginIdentifier () { return this.canCapability('login_identifier_change') },
    canChangePassword () { return this.canCapability('password_change') },
    pageTitle () {
      return this.t(`portal.title.${this.page}`, {
        home: 'City 311', services: 'Services', help: 'Help center', 'sign-in': 'Sign in', register: 'Register', 'forgot-password': 'Forgot password', 'reset-password': 'Reset password', account: 'Account', requests: 'My requests', 'auth-callback': 'Sign in', 'link-confirm': 'Account linking', 'logout-callback': 'Signed out',
      }[this.page] || 'City 311')
    },
    brandName () { return this.branding?.organisation_name || this.t('portal.brand', 'City 311') },
    safeLogoUrl () { return safeResourceURL(this.branding?.logo_url) },
    brandStyle () {
      const style = {}
      const primary = safeColour(this.branding?.primary_colour)
      const accent = safeColour(this.branding?.accent_colour)
      const font = safeFontFamily(this.branding?.font_family)
      if (primary) style['--c311-primary-color'] = primary
      if (accent) style['--c311-accent-color'] = accent
      if (font) style.fontFamily = font
      return style
    },
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
    requestColumns () { return [{ key: 'request_number', label: this.t('field.requestNumber', 'Request number') }, { key: 'summary', label: this.t('field.summary', 'Summary') }, { key: 'status', label: this.t('field.status', 'Status') }, { key: 'updated_at', label: this.t('field.updated', 'Updated'), format: formatDate }] },
    passwordRuleText () { return this.t('password.rules', 'Use 12–128 characters with at least three of uppercase, lowercase, number, and symbol.') },
    dataErrorState () { return stateForError?.(this.dataError || {}) || 'terminal-error' },
  },
  watch: {
    $route: {
      async handler () {
        this.initializePage()
        await this.load()
      },
    },
    forms: {
      deep: true,
      handler () {
        if (this.c311Dirty) this.c311SaveDirtyDraft(this.identityDraft())
      },
    },
  },
  created () {
    this.initializePage()
    this.load()
  },
  methods: {
    t (key, fallback) { const value = this.$t?.(`c311:${key}`); return value && value !== `c311:${key}` && value !== key ? value : fallback },
    formatDate (value) { return formatDate(value) },
    canCapability (capability) {
      if (typeof this.$C311?.can === 'function') return this.$C311.can(capability)
      return !!this.$C311?.session?.actor?.capabilities?.includes(capability)
    },
    initializePage () {
      this.c311DirtyStorageKey = `c311.identity.${this.page}`
      this.restoredDraftFields = []
      this.linkState = this.page === 'link-confirm' ? 'pending' : 'idle'
      const draft = this.c311ReadDirtyDraft()
      const formName = FORM_BY_PAGE[this.page] || this.page
      if (draft && this.forms[formName] && typeof draft === 'object') {
        this.forms[formName] = { ...this.forms[formName], ...draft }
        this.restoredDraftFields = Object.keys(draft)
      }
      this.resetToken = typeof window !== 'undefined' && c311Identity?.resetTokenFromLocation ? c311Identity.resetTokenFromLocation(window.location) : ''
    },
    resetLoadState () {
      this.state = 'loading'
      this.contentState = 'loading'
      this.helpState = 'loading'
      this.dataError = null
      this.helpError = null
      this.statusMessage = ''
      this.successMessage = ''
      this.federatedMessage = ''
      this.contentBody = ''
      this.contentKey = ''
      this.helpBody = ''
      this.branding = null
      this.items = []
      this.formErrors = []
      this.profile = null
      this.linkState = 'idle'
    },
    normalizeField (field) { return String(field || '').replace(/^#/, '').replace(/^\//, '').replace(/\//g, '/') },
    hasError (field) {
      const target = this.normalizeField(field)
      return this.formErrors.some(error => this.normalizeField(error.field) === target || this.normalizeField(error.field).startsWith(`${target}/`))
    },
    identityDraft () {
      const formName = FORM_BY_PAGE[this.page] || this.page
      const form = this.forms[formName]
      if (!form) return {}
      return Object.fromEntries(Object.entries(form).filter(([key]) => !/(?:password|token|secret|credential)/i.test(key)))
    },
    markDirty () { this.c311MarkDirty?.(true) },
    setError (error) {
      const normalized = displayError(error)
      this.dataError = normalized
      const fieldErrors = normalized.fieldErrors.length ? normalized.fieldErrors : normalized.errors
      this.formErrors = fieldErrors.length
        ? fieldErrors
        : [{ field: 'form', code: normalized.code || normalized.error || 'OPERATION_FAILED', message: normalized.message || this.t('error.operationFailed', 'The operation could not be completed.') }]
    },
    async runBusy (name, task) {
      if (this.busy[name]) return
      this.busy[name] = true
      this.statusMessage = this.t('status.inProgress', 'Working…')
      try { return await task() } finally { this.busy[name] = false; this.statusMessage = '' }
    },
    async load () {
      const generation = ++this.loadGeneration
      const isCurrent = () => generation === this.loadGeneration
      this.resetLoadState()
      try {
        const session = typeof this.provider?.getSession === 'function' ? await this.provider.getSession() : undefined
        if (!isCurrent()) return
        if (session !== undefined) {
          this.$C311.session = session
          this.sessionRevision++
        }
        const branding = await this.provider?.getBranding?.()
        if (isCurrent()) this.$data.branding = branding
        if (this.pageNeedsContent) {
          if (!isCurrent()) return
          const contentKey = this.page === 'services' ? 'SERVICE_CATALOGUE' : this.page === 'help' ? 'HELP' : 'HOME'
          const language = (this.$i18n?.locale || 'en').toUpperCase()
          const [contentResult, helpResult] = await Promise.allSettled([
            this.provider?.getPublicContent?.(contentKey),
            this.provider?.getPublicHelp?.('public.request.submit', language),
          ])
          if (!isCurrent()) return
          if (contentResult.status === 'rejected') {
            const normalized = displayError(contentResult.reason)
            this.dataError = normalized
            this.contentState = stateForError?.(normalized) || (normalized.retryable ? 'retryable-error' : 'terminal-error')
            this.state = this.contentState
            return
          }
          this.contentBody = contentResult.value?.body || ''
          this.contentKey = contentResult.value?.content_key || contentKey
          this.contentState = this.contentBody ? 'populated' : 'empty'
          if (helpResult.status === 'fulfilled') {
            this.helpBody = helpResult.value?.body || ''
            this.helpState = this.helpBody ? 'populated' : 'empty'
          } else {
            this.helpError = displayError(helpResult.reason)
            this.helpState = stateForError?.(this.helpError) || (this.helpError.retryable ? 'retryable-error' : 'terminal-error')
          }
          this.state = 'populated'; return
        }
        if (this.page === 'requests') {
          const response = await this.provider?.listPortalRequests?.()
          if (!isCurrent()) return
          this.items = response?.items || []
          this.state = this.items.length ? 'populated' : 'empty'
          return
        }
        if (this.page === 'account') {
          const profile = await this.provider?.getProfile?.()
          if (!isCurrent()) return
          this.profile = profile
          const restored = new Set(this.restoredDraftFields)
          if (!restored.has('display_name')) this.forms.account.display_name = profile?.display_name || ''
          if (!restored.has('preferred_language')) this.forms.account.preferred_language = profile?.preferred_language || 'EN'
          if (!restored.has('login_identifier')) this.forms.account.login_identifier = profile?.login_identifier || ''
          if (!restored.has('phone_numbers')) this.forms.account.phone_numbers = profile?.phone_numbers || []
          if (!restored.has('addresses')) this.forms.account.addresses = profile?.addresses || []
          if (!restored.has('primary_category')) this.forms.account.primary_category = profile?.primary_category || 'RESIDENT'
          this.state = 'populated'
          return
        }
        if (this.page === 'auth-callback') { await this.completeFederated(generation); return }
        if (this.page === 'link-confirm') {
          this.linkState = 'pending'
          this.state = 'populated'
          return
        }
        if (this.page === 'logout-callback') {
          let currentLoad = false
          try {
            await this.provider?.signOut?.()
          } finally {
            currentLoad = isCurrent()
            if (currentLoad) {
              this.$C311?.clearSession?.()
              this.sessionRevision++
            }
          }
          if (!currentLoad) return
          this.successMessage = this.t('identity.loggedOut', 'You are signed out.')
          this.state = 'populated'
          return
        }
        this.state = 'populated'
      } catch (error) {
        if (!isCurrent()) return
        const normalized = displayError(error)
        this.dataError = normalized
        this.contentState = stateForError?.(normalized) || (normalized.retryable ? 'retryable-error' : 'terminal-error')
        this.state = this.contentState
      }
    },
    async signIn () {
      return this.runBusy('signIn', async () => {
        this.formErrors = []
        if (!this.forms.signIn.login_identifier || !this.forms.signIn.password) { this.formErrors = [{ field: !this.forms.signIn.login_identifier ? 'login_identifier' : 'password', code: 'REQUIRED', message: this.t('error.required', 'This field is required.') }]; return }
        try { const session = await this.provider?.signIn?.(this.forms.signIn); this.$C311.session = session; this.sessionRevision++; this.c311ClearDirtyDraft?.(); this.$router.push('/c311') } catch (error) { this.setError(error) }
      })
    },
    async register () {
      return this.runBusy('register', async () => {
        this.formErrors = []
        const errors = []
        if (!this.forms.register.display_name) errors.push({ field: 'display_name', code: 'REQUIRED', message: this.t('error.required', 'This field is required.') })
        if (!this.forms.register.login_identifier) errors.push({ field: 'login_identifier', code: 'REQUIRED', message: this.t('error.required', 'This field is required.') })
        else if (!LOGIN_IDENTIFIER_PATTERN.test(this.forms.register.login_identifier)) errors.push({ field: 'login_identifier', code: 'INVALID_FORMAT', message: this.t('error.loginIdentifier', 'Use 3-64 lowercase letters, numbers, dots, underscores, or hyphens.') })
        if (!this.forms.register.email || !/^[^@\s]+@[^@\s]+\.[^@\s]+$/.test(this.forms.register.email)) errors.push({ field: 'email', code: 'INVALID_FORMAT', message: this.t('error.email', 'Enter a valid email address.') })
        if (validatePassword(this.forms.register.password).length) errors.push({ field: 'password', code: 'INVALID_FORMAT', message: this.passwordRuleText })
        if (errors.length) { this.formErrors = errors; return }
        try { await this.provider?.registerAccount?.(this.forms.register); this.c311ClearDirtyDraft?.(); this.successMessage = this.t('identity.registrationAccepted', 'Registration received. Check your email to continue.') } catch (error) { this.setError(error) }
      })
    },
    async forgotPassword () {
      return this.runBusy('forgot', async () => {
        this.formErrors = []
        if (!this.forms.forgot.email || !/^[^@\s]+@[^@\s]+\.[^@\s]+$/.test(this.forms.forgot.email)) { this.formErrors = [{ field: 'email', code: 'INVALID_FORMAT', message: this.t('error.email', 'Enter a valid email address.') }]; return }
        try { const response = await this.provider?.requestPasswordReset?.(this.forms.forgot); this.c311ClearDirtyDraft?.(); this.successMessage = response?.message || this.t('identity.resetRequested', 'If the account exists, instructions have been sent.') } catch (error) { this.setError(error) }
      })
    },
    async resetPassword () {
      return this.runBusy('reset', async () => {
        this.formErrors = []
        const errors = validatePassword(this.forms.reset.password)
        if (!this.resetToken) errors.push('token')
        if (errors.length) { this.formErrors = [{ field: errors.includes('token') ? 'token' : 'password', code: errors.includes('token') ? 'INVALID_VALUE' : 'INVALID_FORMAT', message: errors.includes('token') ? this.t('error.resetToken', 'The reset link is invalid or expired.') : this.passwordRuleText }]; return }
        try { const response = await this.provider?.confirmPasswordReset?.({ token: this.resetToken, password: this.forms.reset.password }); this.successMessage = response?.message || this.t('identity.passwordReset', 'Your password has been reset.') } catch (error) { this.setError(error) }
      })
    },
    async federated (provider) {
      return this.runBusy('federated', async () => {
        try {
          const redirect = await this.provider?.startFederatedSignIn?.(provider)
          if (!redirect?.authorization_url) {
            this.federatedMessage = this.t('identity.unavailable', 'This sign-in method is unavailable.')
            return
          }
          this.federatedMessage = this.t('identity.redirecting', 'Redirecting to the identity provider…')
          if (typeof window !== 'undefined' && window.C311Mode === 'mock') {
            await this.$router?.push?.({ name: 'c311.auth.callback', query: { provider } })
            return
          }
          if (typeof window !== 'undefined' && window.C311Mode !== 'mock' && typeof window.location?.assign === 'function') window.location.assign(redirect.authorization_url)
        } catch (error) {
          this.setError(error)
          this.federatedMessage = error?.message || this.t('identity.failed', 'Identity sign-in failed. Try again.')
        }
      })
    },
    async confirmAccountLink () {
      return this.runBusy('link', async () => {
        this.formErrors = []
        this.linkState = 'loading'
        try {
          const session = await this.provider?.confirmAccountLink?.()
          this.$C311.pendingFederated = null
          this.$C311.session = session
          this.sessionRevision++
          this.linkState = 'success'
          this.federatedMessage = this.t('identity.linkConfirmed', 'Account linking confirmed.')
        } catch (error) {
          this.linkState = 'error'
          this.setError(error)
        }
      })
    },
    async cancelAccountLink () {
      if (this.$C311) this.$C311.pendingFederated = null
      this.federatedMessage = this.t('identity.linkCancelled', 'Account linking was cancelled.')
      await this.$router?.push?.({ name: 'c311.sign-in' })
    },
    async completeFederated (generation = this.loadGeneration) {
      try {
        const provider = this.$route?.query?.provider === 'saml' ? 'saml' : 'oidc'
        const query = { ...this.$route?.query }
        const result = await this.provider?.completeFederatedSignIn?.(provider, query)
        if (generation !== this.loadGeneration) return
        if (result?.outcome === 'link_confirmation_required') {
          this.$C311.pendingFederated = result.pending_link
          this.federatedMessage = this.t('identity.linkConfirmationRequired', 'This sign-in needs account-link confirmation.')
          await this.$router?.push?.({ name: 'c311.auth.link.confirm' })
          return
        }
        this.$C311.pendingFederated = null
        this.$C311.session = result?.session || result
        this.sessionRevision++
        if (typeof window !== 'undefined' && window.history?.replaceState) window.history.replaceState({}, '', '/c311/auth/callback')
        this.federatedMessage = this.t('identity.success', 'Sign-in completed.')
        this.state = 'populated'
      } catch (error) {
        if (generation !== this.loadGeneration) return
        const normalized = displayError(error)
        this.setError(normalized)
        this.state = stateForError?.(normalized) || (normalized.retryable ? 'retryable-error' : 'terminal-error')
        this.federatedMessage = normalized.message || this.t('identity.failed', 'Identity sign-in failed. Try again.')
      }
    },
    async updateAccount () {
      return this.runBusy('account', async () => {
        this.formErrors = []
        try {
          const phoneNumbers = this.forms.account.phone_numbers.filter(phone => phone.value?.trim())
          const addresses = this.forms.account.addresses.filter(address => Object.entries(address).some(([key, value]) => key !== 'primary' && String(value || '').trim()))
          const requiredAddressFields = ['line1', 'city', 'region', 'postal_code', 'country']
          const incompleteIndex = addresses.findIndex(address => requiredAddressFields.some(field => !String(address[field] || '').trim()))
          if (incompleteIndex >= 0) {
            this.formErrors = [{ field: `addresses/${incompleteIndex}`, code: 'REQUIRED', message: this.t('error.addressRequired', 'Complete every address field or remove the address.') }]
            return
          }
          if (addresses.length && addresses.filter(address => address.primary).length !== 1) {
            this.formErrors = [{ field: 'addresses', code: 'INVALID_VALUE', message: this.t('error.primaryAddress', 'Choose one primary address.') }]
            return
          }
          const profileInput = { display_name: this.forms.account.display_name, preferred_language: this.forms.account.preferred_language, phone_numbers: phoneNumbers, addresses, primary_category: this.forms.account.primary_category }
          this.profile = await this.provider?.updateProfile?.(profileInput, { expectedVersion: this.profile?.version })
          await this.provider?.updateLanguage?.(this.forms.account.preferred_language)
          this.c311ClearDirtyDraft?.()
          this.successMessage = this.t('account.saved', 'Account changes saved.')
        } catch (error) { this.setError(error) }
      })
    },
    async changeLoginIdentifier () {
      return this.runBusy('loginIdentifier', async () => {
        this.activeAccountAction = 'loginIdentifier'
        this.formErrors = []
        const errors = []
        if (!this.forms.account.login_identifier) errors.push({ field: 'login_identifier', code: 'REQUIRED', message: this.t('error.required', 'This field is required.') })
        if (!this.forms.account.current_password) errors.push({ field: 'current_password', code: 'REQUIRED', message: this.t('error.required', 'This field is required.') })
        if (errors.length) { this.formErrors = errors; return }
        try {
          const session = await this.provider?.changeLoginIdentifier?.({ current_password: this.forms.account.current_password, login_identifier: this.forms.account.login_identifier })
          this.$C311.session = session
          this.sessionRevision++
          this.c311ClearDirtyDraft?.()
          this.successMessage = this.t('account.loginIdentifier.saved', 'Login identifier changed.')
        } catch (error) { this.setError(error) }
      })
    },
    async changePassword () {
      return this.runBusy('password', async () => {
        this.activeAccountAction = 'password'
        this.formErrors = []
        const errors = []
        if (!this.forms.account.current_password) errors.push({ field: 'current_password', code: 'REQUIRED', message: this.t('error.required', 'This field is required.') })
        if (validatePassword(this.forms.account.new_password).length) errors.push({ field: 'new_password', code: 'INVALID_FORMAT', message: this.passwordRuleText })
        if (errors.length) { this.formErrors = errors; return }
        try { await this.provider?.changePassword?.({ current_password: this.forms.account.current_password, new_password: this.forms.account.new_password }); this.forms.account.new_password = ''; this.c311ClearDirtyDraft?.(); this.successMessage = this.t('account.password.saved', 'Password changed.') } catch (error) { this.setError(error) }
      })
    },
    async languageChanged (language) {
      language = String(language || 'en').toUpperCase()
      this.forms.account.preferred_language = language
      this.markDirty()
      const locale = language.toLowerCase()
      try {
        if (this.$i18n?.locale !== locale) this.$i18n.locale = locale
        if (this.$i18n?.i18next?.changeLanguage) await this.$i18n.i18next.changeLanguage(locale)
        if (this.isAuthenticated) await this.provider?.updateLanguage?.(language)
        await this.load()
      } catch (error) { this.setError(error) }
    },
    addPhone () {
      if (this.forms.account.phone_numbers.length >= 3) return
      this.forms.account.phone_numbers.push({ label: 'MOBILE', value: '' })
      this.markDirty()
    },
    removePhone (index) {
      this.forms.account.phone_numbers.splice(index, 1)
      this.markDirty()
    },
    addAddress () {
      if (this.forms.account.addresses.length >= 5) return
      this.forms.account.addresses.push({ line1: '', line2: '', city: '', region: '', postal_code: '', country: 'US', primary: this.forms.account.addresses.length === 0 })
      this.markDirty()
    },
    removeAddress (index) {
      const wasPrimary = this.forms.account.addresses[index]?.primary
      this.forms.account.addresses.splice(index, 1)
      if (wasPrimary && this.forms.account.addresses.length) this.forms.account.addresses[0].primary = true
      this.markDirty()
    },
    setPrimaryAddress (index) {
      this.forms.account.addresses.forEach((address, current) => { address.primary = current === index })
      this.markDirty()
    },
  },
}
</script>

<style scoped>
.c311-public-page { max-width: 56rem; }
.c311-public-page :deep(.c311-responsive-data) { min-width: 0; }
.c311-branding { border-left: 4px solid var(--c311-accent-color, #12b76a); padding: 0.75rem 1rem; margin-bottom: 1rem; }
.c311-branding__logo { max-height: 3rem; max-width: 12rem; object-fit: contain; }
.c311-branding__header, .c311-branding__footer { color: var(--c311-primary-color, inherit); }
.c311-branding__footer { margin-top: 2rem; }
</style>
