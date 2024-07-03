<template>
  <b-card
    data-test-id="card-user-info"
    header-class="border-bottom"
    footer-class="border-top d-flex flex-wrap flex-fill-child gap-1"
    class="shadow-sm"
  >
    <b-form
      @submit.prevent="$emit('submit', user)"
    >
      <b-row>
        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('email')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="user.email"
              data-test-id="input-email"
              required
              :state="emailState"
              type="email"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('name')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="user.name"
              data-test-id="input-name"
              required
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('handle')"
            :class="{ 'mb-0': !user.userID }"
            label-class="text-primary"
          >
            <b-form-input
              v-model="user.handle"
              data-test-id="input-handle"
              :placeholder="$t('placeholder-handle')"
              :state="handleState"
            />
            <b-form-invalid-feedback :state="handleState">
              {{ $t('invalid-handle-characters') }}
            </b-form-invalid-feedback>
          </b-form-group>
        </b-col>
      </b-row>

      <c-system-fields
        :id="user.userID"
        :resource="user"
      />

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
      <h4 class="m-0">
        {{ $t('title') }}
      </h4>
    </template>

    <template #footer>
      <c-input-confirm
        v-if="!fresh && user.canDeleteUser"
        :data-test-id="deletedButtonStatusCypressId"
        variant="danger"
        size="md"
        @confirmed="$emit('delete')"
      >
        {{ getDeleteStatus }}
      </c-input-confirm>

      <c-input-confirm
        v-if="!fresh"
        :data-test-id="suspendButtonStatusCypressId"
        variant="light"
        size="md"
        @confirmed="$emit('status')"
      >
        {{ getSuspendStatus }}
      </c-input-confirm>

      <c-input-confirm
        v-if="!fresh"
        data-test-id="button-sessions-revoke"
        :disabled="user.userID === userID"
        variant="light"
        size="md"
        @confirmed="$emit('sessionsRevoke')"
      >
        {{ $t('revokeAllSession') }}
      </c-input-confirm>

      <b-button
        v-if="!fresh && !user.emailConfirmed"
        variant="light"
        @click="$emit('patch', '/emailConfirmed', true)"
      >
        {{ $t('confirmEmail') }}
      </b-button>

      <c-corredor-manual-buttons
        ui-page="user/editor"
        ui-slot="infoFooter"
        resource-type="system:user"
        default-variant="light"
        @click="dispatchCortezaSystemUserEvent($event, { user })"
      />

      <c-button-submit
        :disabled="saveDisabled"
        :processing="processing"
        :success="success"
        :text="$t('admin:general.label.submit')"
        class="ml-auto"
        @submit="$emit('submit', user)"
      />
    </template>
  </b-card>
</template>

<script>
import { NoID } from '@cortezaproject/corteza-js'
import { handle } from '@cortezaproject/corteza-vue'
import { getSystemFields } from 'corteza-webapp-admin/src/lib/sysFields'

export default {
  name: 'CUserEditorInfo',

  i18nOptions: {
    namespaces: 'system.users',
    keyPrefix: 'editor.info',
  },

  props: {
    user: {
      type: Object,
      required: true,
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

  computed: {
    getDeleteStatus () {
      return this.user.deletedAt ? this.$t('undelete') : this.$t('delete')
    },

    getSuspendStatus () {
      return this.user.suspendedAt ? this.$t('unsuspend') : this.$t('suspend')
    },

    userID () {
      if (this.$auth.user) {
        return this.$auth.user.userID
      }
      return undefined
    },

    fresh () {
      return !this.user.userID || this.user.userID === NoID
    },

    editable () {
      return this.fresh ? this.canCreate : this.user.canUpdateUser
    },

    emailState () {
      const { email } = this.user
      return email ? null : false
    },

    handleState () {
      return handle.handleState(this.user.handle)
    },

    saveDisabled () {
      return !this.editable || [this.emailState, this.handleState].includes(false)
    },

    deletedButtonStatusCypressId () {
      return `button-${this.getDeleteStatus.toLowerCase()}`
    },

    suspendButtonStatusCypressId () {
      return `button-${this.getSuspendStatus.toLowerCase()}`
    },

    systemFields () {
      return getSystemFields(this.role)
    },
  },
}
</script>
