<template>
  <b-card
    v-if="authClient"
    class="shadow-sm auth-clients"
    header-bg-variant="white"
    footer-bg-variant="white"
  >
    <b-form
      @submit.prevent="submit"
    >
      <b-form-group
        :label="$t('name')"
        label-cols="3"
      >
        <b-form-input
          v-model="authClient.meta.name"
          required
        />
      </b-form-group>

      <b-form-group
        :label="$t('handle.label')"
        label-cols="3"
      >
        <b-form-input
          v-model="authClient.handle"
          :disabled="authClient.isDefault"
          :placeholder="$t('handle.placeholder-handle')"
          :state="handleState"
        />
        <b-form-invalid-feedback :state="handleState">
          {{ $t('handle.invalid-handle-characters') }}
        </b-form-invalid-feedback>
        <template
          v-if="authClient.isDefault"
          #description
        >
          {{ $t('handle.disabledFootnote') }}
        </template>
      </b-form-group>

      <b-form-group
        :label="$t('redirectURI')"
        label-cols="3"
      >
        <b-button
          variant="light"
          class="align-top"
          @click="redirectURI.push('')"
        >
          + {{ $t('add') }}
        </b-button>

        <div
          v-if="redirectURI.length"
        >
          <b-input-group
            v-for="(value, index) in redirectURI"
            :key="index"
            class="mt-2"
          >
            <b-form-input
              v-model="redirectURI[index]"
              :placeholder="$t('uri')"
            />

            <b-button
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
      </b-form-group>

      <b-form-group
        v-if="!fresh"
        :label="$t('secret')"
        label-cols="3"
        class="mb-3"
      >
        <b-input-group>
          <b-form-input
            v-model="secret"
            disabled
            placeholder="****************************************************************"
          />

          <b-button
            v-if="!secretVisible"
            class="ml-1 text-primary"
            variant="link"
            @click="$emit('request-secret')"
          >
            <font-awesome-icon
              :icon="['fas', 'eye']"
            />
          </b-button>

          <b-button
            v-else
            class="ml-1 text-primary"
            variant="link"
            @click="$emit('regenerate-secret')"
          >
            <font-awesome-icon
              :icon="['fas', 'sync']"
            />
          </b-button>
        </b-input-group>
      </b-form-group>

      <b-form-group
        label-cols="3"
      >
        <b-form-radio-group
          v-model="authClient.validGrant"
          value="authorization_code"
          :options="[
            { value: 'authorization_code', text: $t('grant.authorization_code') },
            { value: 'client_credentials', text: $t('grant.client_credentials') },
          ]"
        />
      </b-form-group>

      <b-form-group
        :label="$t('validFrom.label')"
        label-cols="3"
        :description="$t('validFrom.description')"
      >
        <b-input-group>
          <b-form-datepicker
            v-model="validFrom.date"
            :placeholder="$t('choose-date')"
            locale="en"
          />

          <b-form-timepicker
            v-model="validFrom.time"
            :placeholder="$t('no-time')"
            locale="en"
          />

          <b-button
            class="ml-1 text-secondary"
            variant="link"
            @click="resetDateTime('validFrom')"
          >
            <font-awesome-icon
              :icon="['fas', 'sync']"
            />
          </b-button>
        </b-input-group>
      </b-form-group>

      <b-form-group
        :label="$t('expiresAt.label')"
        label-cols="3"
        :description="$t('expiresAt.description')"
      >
        <b-input-group>
          <b-form-datepicker
            v-model="expiresAt.date"
            :placeholder="$t('choose-date')"
            locale="en"
          />

          <b-form-timepicker
            v-model="expiresAt.time"
            :placeholder="$t('no-time')"
            locale="en"
          />

          <b-button
            class="ml-1 text-secondary"
            variant="link"
            @click="resetDateTime('expiresAt')"
          >
            <font-awesome-icon
              :icon="['fas', 'sync']"
            />
          </b-button>
        </b-input-group>
      </b-form-group>

      <b-form-group
        label-cols="3"
      >
        <b-form-checkbox
          :checked="(authClient.scope || []).includes('profile')"
          @change="setScope($event, 'profile')"
        >
          {{ $t('profile') }}
        </b-form-checkbox>
        <b-form-checkbox
          :checked="(authClient.scope || []).includes('api')"
          @change="setScope($event, 'api')"
        >
          {{ $t('api') }}
        </b-form-checkbox>
        <b-form-checkbox
          :checked="(authClient.scope || []).includes('openid')"
          @change="setScope($event, 'openid')"
        >
          {{ $t('openid') }}
        </b-form-checkbox>
        <b-form-checkbox
          v-if="discoveryEnabled"
          :checked="(authClient.scope || []).includes('discovery')"
          @change="setScope($event, 'discovery')"
        >
          {{ $t('discovery') }}
        </b-form-checkbox>
      </b-form-group>

      <b-form-group
        label-cols="3"
      >
        <b-form-checkbox
          v-model="authClient.trusted"
        >
          {{ $t('trusted.label') }}
        </b-form-checkbox>
        <b-form-text>{{ $t('trusted.description') }}</b-form-text>
      </b-form-group>

      <b-form-group
        label-cols="3"
      >
        <b-form-checkbox
          v-model="authClient.enabled"
          :disabled="authClient.isDefault"
        >
          {{ $t('enabled.label') }}
        </b-form-checkbox>

        <template
          v-if="authClient.isDefault"
          #description
        >
          {{ $t('enabled.disabledFootnote') }}
        </template>
      </b-form-group>

      <div v-if="isClientCredentialsGrant">
        <b-form-group
          label-cols="3"
          :label="$t('security.impersonateUser.label')"
          :description="$t('security.impersonateUser.description')"
        >
          <c-select-user
            :user-i-d="authClient.security.impersonateUser"
            @updateUser="onUpdateUser"
          />
        </b-form-group>
        <div v-if="!fresh">
          <b-form-group label-cols="3">
            <b-button
              variant="light"
              class="align-top"
              @click="toggleCurlSnippet()"
            >
              <template v-if="curlVisible">
                {{ $t('hideCurl') }}
              </template>
              <template v-else>
                {{ $t('generateCurl') }}
              </template>
            </b-button>
          </b-form-group>
          <b-form-group
            v-if="curlVisible"
            :label="$t('cUrl')"
            label-cols="3"
            class="curl"
          >
            <div class="w-100">
              <div class="d-flex">
                <pre
                  ref="cUrl"
                >
curl -X POST {{ curlURL }} \
-d grant_type=client_credentials \
-d scope='profile api' \
-u {{ authClient.authClientID }}:{{ secret || 'PLACE-YOUR-CLIENT-SECRET-HERE' }}
                </pre>
                <b-button
                  variant="link"
                  class="align-top ml-auto fit-content text-secondary"
                  @click="copyToClipboard('cUrl')"
                >
                  <font-awesome-icon
                    :icon="['far', 'copy']"
                  />
                </b-button>
              </div>
              <div class="d-flex">
                <div
                  class="overflow-wrap mr-2 mb-2"
                  :class="[tokenRequest.token ? 'text-success' : 'text-danger']"
                >
                  {{ tokenRequest.token || tokenRequest.error }}
                </div>
                <b-button
                  v-if="tokenRequest.token"
                  variant="link"
                  class="align-top ml-auto fit-content text-secondary"
                  @click="copyToClipboard('token')"
                >
                  <font-awesome-icon
                    :icon="['far', 'copy']"
                  />
                </b-button>
              </div>
            </div>
            <div
              v-if="secretVisible"
              class="d-flex mb-3"
            >
              <b-button
                variant="light"
                class="align-top fit-content"
                @click="getAccessTokenAPI()"
              >
                {{ $t('testCurl') }}
              </b-button>
            </div>
          </b-form-group>
        </div>
      </div>

      <b-form-group
        :label="$t('security.permittedRoles.label')"
        label-cols="3"
        class="mb-0"
      >
        <c-role-picker
          v-model="authClient.security.permittedRoles"
          class="mb-3"
        >
          <template #description>
            {{ $t('security.permittedRoles.description') }}
          </template>
        </c-role-picker>
      </b-form-group>

      <b-form-group
        :label="$t('security.forbiddenRoles.label')"
        label-cols="3"
        class="mb-0"
      >
        <c-role-picker
          v-model="authClient.security.forbiddenRoles"
          class="mb-3"
        >
          <template #description>
            {{ $t('security.forbiddenRoles.description') }}
          </template>
        </c-role-picker>
      </b-form-group>

      <b-form-group
        :label="$t('security.forcedRoles.label')"
        label-cols="3"
        class="mb-0"
      >
        <c-role-picker
          v-model="authClient.security.forcedRoles"
          class="mb-3"
        >
          <template #description>
            {{ $t('security.forcedRoles.description') }}
          </template>
        </c-role-picker>
      </b-form-group>

      <b-form-group
        v-if="authClient.createdAt"
        :label="$t('createdAt')"
        label-cols="3"
        class="mb-0"
      >
        <b-form-input
          :value="authClient.createdAt | locFullDateTime"
          plaintext
          disabled
        />
      </b-form-group>

      <b-form-group
        v-if="authClient.updatedAt"
        :label="$t('updatedAt')"
        label-cols="3"
      >
        <b-form-input
          :value="authClient.updatedAt | locFullDateTime"
          plaintext
          disabled
        />
      </b-form-group>

      <b-form-group
        v-if="authClient.deletedAt"
        :label="$t('deletedAt')"
        label-cols="3"
      >
        <b-form-input
          :value="authClient.deletedAt | locFullDateTime"
          plaintext
          disabled
        />
      </b-form-group>

      <!--
        include hidden input to enable
        trigger submit event w/ ENTER
      -->
      <input
        type="submit"
        class="d-none"
        :disabled="saveDisabled"
      >
    </b-form>

    <template #header>
      <h3 class="m-0">
        {{ $t('title') }}
      </h3>
    </template>

    <template #footer>
      <c-submit-button
        class="float-right"
        :disabled="saveDisabled"
        :processing="processing"
        :success="success"
        @submit="submit"
      />

      <template
        v-if="canDelete"
      >
        <confirmation-toggle
          v-if="isDeleted"
          :disabled="processing"
          @confirmed="$emit('undelete', authClient.authClientID)"
        >
          {{ $t('undelete') }}
        </confirmation-toggle>
        <confirmation-toggle
          v-else
          :disabled="processing"
          @confirmed="$emit('delete', authClient.authClientID)"
        >
          {{ $t('delete') }}
        </confirmation-toggle>
      </template>
    </template>
  </b-card>
</template>

<script>
import { NoID } from '@cortezaproject/corteza-js'
import { handleState } from 'corteza-webapp-admin/src/lib/handle'
import Vue from 'vue'
import ConfirmationToggle from 'corteza-webapp-admin/src/components/ConfirmationToggle'
import CSubmitButton from 'corteza-webapp-admin/src/components/CSubmitButton'
import CRolePicker from 'corteza-webapp-admin/src/components/CRolePicker'
import CSelectUser from 'corteza-webapp-admin/src/components/Authclient/CSelectUser'
import copy from 'copy-to-clipboard'
import axios from 'axios'

const defSecurity = Object.freeze({
  impersonateUser: '0',
  permittedRoles: [],
  forbiddenRoles: [],
  forcedRoles: [],
})

export default {
  name: 'CAuthclientEditorInfo',

  i18nOptions: {
    namespaces: 'system.authclients',
    keyPrefix: 'editor.info',
  },

  components: {
    ConfirmationToggle,
    CSubmitButton,
    CRolePicker,
    CSelectUser,
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

    secret: {
      type: String,
      default: () => '',
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
    const authClient = Vue.util.extend({
      trusted: false,
      handle: '',
      meta: {
        name: '',
        description: '',
      },

      redirectURI: '',
      validGrant: '',

      // make sure all references are destroyed
    }, this.resource)

    if (!(authClient.security || {}).impersonateUser) {
      // hande empty security struct
      authClient.security = { ...defSecurity }
    }

    return {
      // setup all object props we need (reactivity)
      // when we migrate it to corteza-js using a proper Class this can remove it
      authClient,

      redirectURI: this.resource.redirectURI ? this.resource.redirectURI.split(' ') : [],

      // @todo should be handled via computed props
      validFrom: this.resource.validFrom ? {
        date: new Date(this.resource.validFrom).toISOString(),
        time: new Date(this.resource.validFrom).toTimeString().split(' ')[0],
      } : { date: null, time: null },

      // @todo should be handled via computed props
      expiresAt: this.resource.expiresAt ? {
        date: new Date(this.resource.expiresAt).toISOString(),
        time: new Date(this.resource.expiresAt).toTimeString().split(' ')[0],
      } : { date: null, time: null },

      curlVisible: false,
      curlURL: '',
      tokenRequest: {
        token: '',
        error: '',
      },
    }
  },

  computed: {
    fresh () {
      return !this.authClient.authClientID || this.authClient.authClientID === NoID
    },

    editable () {
      return this.fresh ? this.canCreate : this.authClient.canUpdateAuthClient
    },

    isDeleted () {
      return this.authClient.deletedAt
    },

    secretVisible () {
      return this.secret.length > 0
    },

    handleState () {
      const { handle } = this.authClient

      return handle ? handleState(handle) : false
    },

    isClientCredentialsGrant () {
      return this.authClient.validGrant === 'client_credentials'
    },

    discoveryEnabled () {
      return this.$Settings.get('discovery.enabled', false)
    },

    saveDisabled () {
      return !this.editable || [this.handleState].includes(false)
    },
  },

  watch: {
    'redirectURI': {
      handler (redirectURI) {
        this.authClient.redirectURI = redirectURI.filter(ru => ru).join(' ')
      },
    },
  },

  methods: {
    onUpdateUser (user) {
      this.authClient.security.impersonateUser = (user || {}).userID
    },

    getAccessTokenAPI () {
      const params = new URLSearchParams()
      params.append('grant_type', 'client_credentials')
      params.append('scope', 'profile api')
      axios.post(
        this.curlURL,
        params,
        { auth: { username: this.authClient.authClientID, password: this.secret } }
      ).then(response => {
        this.tokenRequest.token = (response.data || {}).access_token
      }).catch(error => {
        this.tokenRequest.error = error
      })
    },

    copyToClipboard (name) {
      if (name === 'cUrl') {
        copy(this.$refs.cUrl.innerHTML)
      } else {
        copy(this.tokenRequest.token)
      }
    },

    toggleCurlSnippet () {
      if (!this.curlVisible) {
        this.curlURL = this.$auth.cortezaAuthURL + '/oauth2/token'
      }
      this.curlVisible = !this.curlVisible
    },

    submit () {
      if (this.validFrom.date && this.validFrom.time) {
        this.authClient.validFrom = new Date(`${this.validFrom.date} ${this.validFrom.time}`).toISOString()
      } else {
        this.authClient.validFrom = undefined
      }

      if (!this.isClientCredentialsGrant || !this.authClient.security.impersonateUser) {
        this.authClient.security.impersonateUser = '0'
      }

      if (this.expiresAt.date && this.expiresAt.time) {
        this.authClient.expiresAt = new Date(`${this.expiresAt.date} ${this.expiresAt.time}`).toISOString()
      } else {
        this.authClient.expiresAt = undefined
      }

      this.$emit('submit', this.authClient)
    },

    setScope (value, target) {
      let items = this.authClient.scope ? this.authClient.scope.split(' ') : []

      if (value) {
        items.push(target)
      } else {
        items = items.filter(i => i !== target)
      }

      this.authClient.scope = items.join(' ')
    },

    resetDateTime (target) {
      if (target) {
        this[target].date = undefined
        this[target].time = undefined
      }
    },
  },
}
</script>
<style lang="scss">
.auth-clients{
  .fit-content{
    height:fit-content;
  }
  .overflow-wrap{
      overflow-wrap: anywhere;
  }
  .curl .form-row{
    flex-wrap: nowrap !important;
    .col{
      max-width: 84.3%;
    }
  }
}
</style>
