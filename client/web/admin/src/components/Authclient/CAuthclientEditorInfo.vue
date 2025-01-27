<template>
  <b-card
    v-if="resource"
    data-test-id="card-auth-client-info"
    header-class="border-bottom"
    footer-class="border-top d-flex flex-wrap flex-fill-child gap-1"
    class="shadow-sm auth-clients"
  >
    <b-form
      @submit.prevent="submit"
    >
      <b-row>
        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('name')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="resource.meta.name"
              data-test-id="input-name"
              required
              :state="nameState"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('handle.label')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="resource.handle"
              data-test-id="input-handle"
              :disabled="resource.isDefault"
              :placeholder="$t('handle.placeholder-handle')"
              :state="handleState"
            />
            <b-form-invalid-feedback
              data-test-id="feedback-invalid-handle"
              :state="handleState"
            >
              {{ $t('handle.invalid-handle-characters') }}
            </b-form-invalid-feedback>

            <template
              v-if="resource.isDefault"
              #description
            >
              {{ $t('handle.disabledFootnote') }}
            </template>
          </b-form-group>
        </b-col>

        <b-col cols="12">
          <b-form-group
            :label="$t('redirectURI')"
            label-class="text-primary"
          >
            <c-form-table-wrapper
              :labels="{ addButton: $t('general:label.add') }"
              test-id="button-add-redirect-uris"
              @add-item="redirectURI.push('')"
            >
              <div
                v-if="redirectURI.length"
              >
                <b-input-group
                  v-for="(value, index) in redirectURI"
                  :key="index"
                  class="mb-2"
                >
                  <b-form-input
                    v-model="redirectURI[index]"
                    data-test-id="input-uri"
                    :placeholder="$t('uri')"
                  />
                  <b-button
                    data-test-id="button-remove-uri"
                    class="ml-1 text-danger"
                    variant="link"
                    @click="redirectURI.splice(index, 1)"
                  >
                    <font-awesome-icon
                      :icon="['fas', 'times']"
                    />
                  </b-button>
                </b-input-group>
              </div>
            </c-form-table-wrapper>
          </b-form-group>
        </b-col>

        <b-col cols="12">
          <b-form-group
            v-if="!fresh"
            :label="$t('secret')"
            label-class="text-primary"
          >
            <template #label>
              {{ $t('secret') }}
              <b-button
                v-if="!secretVisible"
                v-b-tooltip.noninteractive.hover="{ title: $t('tooltip.show-client-secret'), container: '#body' }"
                data-test-id="button-show-client-secret"
                variant="outline-light"
                size="sm"
                class="text-secondary border-0"
                @click="showSecret()"
              >
                <font-awesome-icon
                  :icon="['fas', 'eye']"
                />
              </b-button>

              <b-button
                v-else
                v-b-tooltip.noninteractive.hover="{ title: $t('tooltip.hide-client-secret'), container: '#body' }"
                data-test-id="button-hide-client-secret"
                variant="outline-light"
                size="sm"
                class="text-secondary border-0"
                @click="hideSecret()"
              >
                <font-awesome-icon
                  :icon="['fas', 'eye-slash']"
                />
              </b-button>
            </template>

            <b-input-group>
              <b-form-input
                v-model="secret"
                data-test-id="input-client-secret"
                disabled
                placeholder="****************************************************************"
              />

              <b-button
                v-b-tooltip.noninteractive.hover="{ title: $t('tooltip.regenerate-secret'), container: '#body' }"
                data-test-id="button-regenerate-client-secret"
                variant="outline-light"
                class="text-secondary border-0"
                @click="regenerateSecret()"
              >
                <font-awesome-icon
                  :icon="['fas', 'sync']"
                />
              </b-button>
            </b-input-group>
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group>
            <b-form-radio-group
              v-model="resource.validGrant"
              value="authorization_code"
              :options="[
                { value: 'authorization_code', text: $t('grant.authorization_code') },
                { value: 'client_credentials', text: $t('grant.client_credentials') },
              ]"
              @change="onGrantChange"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group>
            <b-form-checkbox
              data-test-id="checkbox-allow-access-to-user-profile"
              :checked="((resource.scope || []).includes('profile'))"
              @change="setScope($event, 'profile')"
            >
              {{ $t('profile') }}
            </b-form-checkbox>

            <b-form-checkbox
              data-test-id="checkbox-allow-access-to-corteza-api"
              :checked="((resource.scope || []).includes('api'))"
              @change="setScope($event, 'api')"
            >
              {{ $t('api') }}
            </b-form-checkbox>

            <b-form-checkbox
              data-test-id="checkbox-allow-client-to-use-oidc"
              :checked="((resource.scope || []).includes('openid'))"
              @change="setScope($event, 'openid')"
            >
              {{ $t('openid') }}
            </b-form-checkbox>

            <b-form-checkbox
              v-if="discoveryEnabled"
              data-test-id="checkbox-allow-client-access-to-discovery"
              :checked="((resource.scope || []).includes('discovery'))"
              @change="setScope($event, 'discovery')"
            >
              {{ $t('discovery') }}
            </b-form-checkbox>
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            data-test-id="valid-from"
            :label="$t('validFrom.label')"
            :description="$t('validFrom.description')"
            label-class="text-primary"
          >
            <c-input-date-time
              v-model="resource.validFrom"
              data-test-id="input-valid-from"
              :labels="{
                clear: $t('general:label.clear'),
                none: $t('general:label.none'),
                now: $t('general:label.now'),
                today: $t('general:label.today'),
              }"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            data-test-id="expires-at"
            :label="$t('expiresAt.label')"
            :description="$t('expiresAt.description')"
            label-class="text-primary"
          >
            <c-input-date-time
              v-model="resource.expiresAt"
              data-test-id="input-expires-at"
              :labels="{
                clear: $t('general:label.clear'),
                none: $t('general:label.none'),
                now: $t('general:label.now'),
                today: $t('general:label.today'),
              }"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group>
            <b-form-checkbox
              v-model="resource.enabled"
              data-test-id="checkbox-is-client-enabled"
              :disabled="resource.isDefault"
            >
              {{ $t('enabled.label') }}
            </b-form-checkbox>

            <b-form-text v-if="resource.isDefault">
              {{ $t('enabled.disabledFootnote') }}
            </b-form-text>
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group>
            <b-form-checkbox
              v-model="resource.trusted"
              data-test-id="checkbox-is-client-trusted"
            >
              {{ $t('trusted.label') }}
            </b-form-checkbox>
            <b-form-text>{{ $t('trusted.description') }}</b-form-text>
          </b-form-group>
        </b-col>
      </b-row>

      <b-row>
        <b-col
          v-show="isClientCredentialsGrant"
          cols="12"
          lg="6"
        >
          <b-form-group
            data-test-id="impersonate-user"
            :label="$t('security.impersonateUser.label')"
            :description="$t('security.impersonateUser.description')"
            label-class="text-primary"
          >
            <c-select-user
              :user-i-d="resource.security.impersonateUser"
              :placeholder="$t('security.impersonateUser.placeholder')"
              :clearable="false"
              @input="onUpdateUser"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            data-test-id="permitted-roles"
            :label="$t('security.permittedRoles.label')"
            label-class="text-primary"
          >
            <c-role-picker
              v-model="resource.security.permittedRoles"
            >
              <template #description>
                {{ $t('security.permittedRoles.description') }}
              </template>
            </c-role-picker>
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('security.prohibitedRoles.label')"
            data-test-id="prohibited-roles"
            label-class="text-primary"
          >
            <c-role-picker
              v-model="resource.security.prohibitedRoles"
            >
              <template #description>
                {{ $t('security.prohibitedRoles.description') }}
              </template>
            </c-role-picker>
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            data-test-id="forced-roles"
            :label="$t('security.forcedRoles.label')"
            label-class="text-primary"
          >
            <c-role-picker
              v-model="resource.security.forcedRoles"
              class="mb-3"
            >
              <template #description>
                {{ $t('security.forcedRoles.description') }}
              </template>
            </c-role-picker>
          </b-form-group>
        </b-col>

        <b-col
          v-if="!fresh && isClientCredentialsGrant"
          cols="12"
        >
          <b-form-group label-class="text-primary">
            <template #label>
              {{ $t('cUrl') }}
              <b-button
                v-b-tooltip.noninteractive.hover="{ title: $t('tooltip.copy-cURL'), container: '#body' }"
                data-test-id="button-copy-cURL"
                variant="outline-light"
                size="sm"
                class="text-secondary border-0"
                @click="copyToClipboard(exampleCurl)"
              >
                <font-awesome-icon
                  :icon="['far', 'copy']"
                />
              </b-button>
            </template>

            <b-textarea
              :value="exampleCurl"
              disabled
              data-test-id="cURL-string"
              rows="3"
              class="mb-2"
            />
          </b-form-group>

          <b-form-group label-class="text-primary">
            <template #label>
              {{ $t('accessToken') }}

              <b-button
                v-if="tokenRequest.token"
                v-b-tooltip.noninteractive.hover="{ title: $t('tooltip.copy-access-token'), container: '#body' }"
                data-test-id="copy-token-from-request"
                variant="outline-light"
                size="sm"
                class="text-secondary border-0"
                @click="copyToClipboard(tokenRequest.token)"
              >
                <font-awesome-icon
                  :icon="['far', 'copy']"
                />
              </b-button>
            </template>

            <b-textarea
              v-if="tokenRequest.token"
              data-test-id="cURL-string"
              :value="tokenRequest.token"
              disabled
              rows="5"
            />

            <b-button
              v-else
              data-test-id="button-test-cURL"
              variant="light"
              @click="getAccessTokenAPI()"
            >
              {{ $t('generateAccessToken') }}
            </b-button>

            <p
              v-if="tokenRequest.error"
              class="text-danger mt-2"
            >
              {{ tokenRequest.error }}
            </p>
          </b-form-group>
        </b-col>
      </b-row>

      <c-system-fields
        :resource="resource"
      />

      <!--
        include hidden input to enable
        trigger submit event w/ ENTER
      -->
      <input
        data-test-id="button-submit"
        type="submit"
        class="d-none"
        :disabled="saveDisabled"
      >
    </b-form>

    <template #header>
      <h4 class="m-0">
        {{ $t('title') }}
      </h4>
    </template>

    <template #footer>
      <template
        v-if="canDelete"
      >
        <c-input-confirm
          :data-test-id="isDeleted ? 'button-undelete': 'button-delete'"
          :disabled="processing"
          variant="danger"
          size="md"
          @confirmed="$emit(isDeleted ? 'undelete' : 'delete', resource.authClientID)"
        >
          {{ isDeleted ? $t('undelete') : $t('delete') }}
        </c-input-confirm>
      </template>

      <c-button-submit
        :disabled="saveDisabled"
        :processing="processing"
        :success="success"
        :text="$t('admin:general.label.submit')"
        class="ml-auto"
        @submit="submit"
      />
    </template>
  </b-card>
</template>

<script>
import { NoID } from '@cortezaproject/corteza-js'
import { handle, components } from '@cortezaproject/corteza-vue'
import CRolePicker from 'corteza-webapp-admin/src/components/CRolePicker'
import CSelectUser from 'corteza-webapp-admin/src/components/Authclient/CSelectUser'
import copy from 'copy-to-clipboard'
import axios from 'axios'

const { CInputDateTime } = components

export default {
  name: 'CAuthclientEditorInfo',

  i18nOptions: {
    namespaces: 'system.authclients',
    keyPrefix: 'editor.info',
  },

  components: {
    CRolePicker,
    CSelectUser,
    CInputDateTime,
  },

  props: {
    resource: {
      type: Object,
      required: true,
    },

    canDelete: {
      type: Boolean,
      default: () => false,
    },

    processing: {
      type: Boolean,
      value: false,
    },

    success: {
      type: Boolean,
      value: false,
    },

    canCreate: {
      type: Boolean,
      required: true,
    },
  },

  data () {
    return {
      requestedSecret: '',
      secret: '',

      redirectURI: this.resource.redirectURI ? this.resource.redirectURI.split(' ') : [],

      curlVisible: false,

      tokenRequest: {
        token: '',
        error: '',
      },

      checkboxLabel: {
        on: this.$t('general:label.general.yes'),
        off: this.$t('general:label.general.no'),
      },
    }
  },

  computed: {
    fresh () {
      return !this.resource.authClientID || this.resource.authClientID === NoID
    },

    editable () {
      return this.fresh ? this.canCreate : this.resource.canUpdateAuthClient
    },

    isDeleted () {
      return this.resource.deletedAt && this.resource.canDeleteAuthClient
    },

    secretVisible () {
      return this.secret.length > 0
    },

    nameState () {
      return this.resource.meta.name ? null : false
    },

    handleState () {
      return handle.handleState(this.resource.handle)
    },

    isClientCredentialsGrant () {
      return this.resource.validGrant === 'client_credentials'
    },

    discoveryEnabled () {
      return this.$Settings.get('discovery.enabled', false)
    },

    saveDisabled () {
      return !this.editable || [this.nameState, this.handleState].includes(false)
    },

    curlURL () {
      return this.$auth.cortezaAuthURL + '/oauth2/token'
    },

    exampleCurl () {
      return `curl -X POST ${this.curlURL} -d grant_type=${this.resource.validGrant} -d scope='${this.resource.scope}' -u ${this.resource.authClientID}:${this.secret || 'PLACE-YOUR-CLIENT-SECRET-HERE'}`
    },
  },

  watch: {
    redirectURI: {
      handler (redirectURI) {
        this.resource.redirectURI = redirectURI.filter(ru => ru).join(' ')
      },
    },
  },

  methods: {
    onUpdateUser (userID) {
      this.resource.security.impersonateUser = userID
    },

    onGrantChange (grant) {
      if (grant === 'client_credentials' && (!this.resource.security.impersonateUser || this.resource.security.impersonateUser === NoID)) {
        this.resource.security.impersonateUser = this.$auth.user.userID
      }
    },

    copyToClipboard (value) {
      copy(value)
    },

    toggleCurlSnippet () {
      this.curlVisible = !this.curlVisible
    },

    submit () {
      if (!this.isClientCredentialsGrant || !this.resource.security.impersonateUser) {
        this.resource.security.impersonateUser = NoID
      }

      this.$emit('submit', this.resource)
    },

    setScope (value, target) {
      let items = this.resource.scope ? this.resource.scope.split(' ') : []

      if (value) {
        items.push(target)
      } else {
        items = items.filter(i => i !== target)
      }

      this.resource.scope = items.join(' ')
    },

    requestSecret () {
      const clientID = this.resource.authClientID

      return this.$SystemAPI.authClientExposeSecret(({ clientID })).then(secret => {
        this.requestedSecret = secret
      })
    },

    async showSecret () {
      if (!this.requestedSecret) {
        await this.requestSecret()
      }

      this.secret = this.requestedSecret
    },

    hideSecret () {
      this.secret = ''
    },

    async regenerateSecret () {
      const clientID = this.resource.authClientID

      this.$SystemAPI.authClientRegenerateSecret(({ clientID })).then(secret => {
        this.requestedSecret = secret
      })
    },

    async getAccessTokenAPI () {
      const clientID = this.resource.authClientID

      if (!this.requestedSecret) {
        await this.requestSecret()
      }

      const params = new URLSearchParams()

      params.append('grant_type', this.resource.validGrant)
      params.append('scope', this.resource.scope)

      axios.post(this.curlURL, params, { auth: { username: clientID, password: this.requestedSecret } }).then(response => {
        this.tokenRequest.token = (response.data || {}).access_token
        this.tokenRequest.error = ''
      }).catch(e => {
        const { error } = e.response.data || {}
        this.tokenRequest.error = error
        this.tokenRequest.token = ''
      })
    },
  },
}
</script>
