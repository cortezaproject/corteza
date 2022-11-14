<template>
  <b-card
    class="shadow-sm"
    data-test-id="card-user-info"
    header-bg-variant="white"
    footer-bg-variant="white"
  >
    <b-form
      @submit.prevent="$emit('submit', user)"
    >
      <b-form-group
        :label="$t('email')"
        label-cols="2"
      >
        <b-form-input
          v-model="user.email"
          data-test-id="input-email"
          required
          :state="emailState"
          type="email"
        />
      </b-form-group>

      <b-form-group
        :label="$t('name')"
        label-cols="2"
      >
        <b-form-input
          v-model="user.name"
          data-test-id="input-name"
          required
        />
      </b-form-group>

      <b-form-group
        :label="$t('handle')"
        label-cols="2"
        :class="{ 'mb-0': !user.userID }"
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

      <b-form-group
        v-if="user.updatedAt"
        :label="$t('updatedAt')"
        label-cols="2"
      >
        <b-form-input
          data-test-id="input-updated-at"
          :value="user.updatedAt | locFullDateTime"
          plaintext
          disabled
        />
      </b-form-group>

      <b-form-group
        v-if="user.suspendedAt"
        :label="$t('suspendedAt')"
        label-cols="2"
      >
        <b-form-input
          data-test-id="input-suspended-at"
          :value="user.suspendedAt | locFullDateTime"
          plaintext
          disabled
        />
      </b-form-group>

      <b-form-group
        v-if="user.deletedAt"
        :label="$t('deletedAt')"
        label-cols="2"
      >
        <b-form-input
          data-test-id="input-deleted-at"
          :value="user.deletedAt | locFullDateTime"
          plaintext
          disabled
        />
      </b-form-group>

      <b-form-group
        v-if="user.createdAt"
        :label="$t('createdAt')"
        label-cols="2"
        class="mb-0"
      >
        <b-form-input
          data-test-id="input-created-at"
          :value="user.createdAt | locFullDateTime"
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
        :processing="processing"
        :success="success"
        :disabled="saveDisabled"
        @submit="$emit('submit', user)"
      />

      <confirmation-toggle
        v-if="!fresh && user.canDeleteUser"
        :data-test-id="deletedButtonStatusCypressId"
        @confirmed="$emit('delete')"
      >
        {{ getDeleteStatus }}
      </confirmation-toggle>

      <confirmation-toggle
        v-if="!fresh"
        :data-test-id="suspendButtonStatusCypressId"
        class="ml-1"
        cta-class="light"
        @confirmed="$emit('status')"
      >
        {{ getSuspendStatus }}
      </confirmation-toggle>

      <confirmation-toggle
        v-if="!fresh"
        data-test-id="button-sessions-revoke"
        :disabled="user.userID === userID"
        class="ml-1"
        cta-class="secondary"
        @confirmed="$emit('sessionsRevoke')"
      >
        {{ $t('revokeAllSession') }}
      </confirmation-toggle>

      <b-button
        v-if="!fresh && !user.emailConfirmed"
        variant="light"
        class="ml-1"
        @click="$emit('patch', '/emailConfirmed', true)"
      >
        {{ $t('confirmEmail') }}
      </b-button>

      <c-corredor-manual-buttons
        ui-page="user/editor"
        ui-slot="infoFooter"
        resource-type="system:user"
        default-variant="secondary"
        class="ml-2"
        @click="dispatchCortezaSystemUserEvent($event, { user })"
      />
    </template>
  </b-card>
</template>

<script>
import { NoID } from '@cortezaproject/corteza-js'
import { handle } from '@cortezaproject/corteza-vue'
import ConfirmationToggle from 'corteza-webapp-admin/src/components/ConfirmationToggle'
import CSubmitButton from 'corteza-webapp-admin/src/components/CSubmitButton'

export default {
  name: 'CUserEditorInfo',

  i18nOptions: {
    namespaces: 'system.users',
    keyPrefix: 'editor.info',
  },

  components: {
    ConfirmationToggle,
    CSubmitButton,
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
  },
}
</script>
