<template>
  <b-card
    class="shadow-sm"
    data-test-id="card-application-info"
    header-bg-variant="white"
    footer-bg-variant="white"
  >
    <b-form
      @submit.prevent="$emit('submit', application)"
    >
      <b-form-group
        v-if="application.applicationID"
        :label="$t('id')"
        label-cols="2"
      >
        <b-form-input
          v-model="application.applicationID"
          data-test-id="input-application-id"
          plaintext
          disabled
        />
      </b-form-group>

      <b-form-group
        :label="$t('name')"
        label-cols="2"
      >
        <b-form-input
          v-model="application.name"
          data-test-id="input-name"
          :state="nameState"
          required
        />
      </b-form-group>

      <b-form-group
        label-cols="2"
        :class="{ 'mb-0': !application.applicationID }"
      >
        <b-form-checkbox
          v-model="application.enabled"
          data-test-id="checkbox-enabled"
        >
          {{ $t('enabled') }}
        </b-form-checkbox>
      </b-form-group>

      <b-form-group
        v-if="application.updatedAt"
        :label="$t('updatedAt')"
        label-cols="2"
      >
        <b-form-input
          v-model="application.updatedAt"
          data-test-id="input-updated-at"
          plaintext
          disabled
        />
      </b-form-group>

      <b-form-group
        v-if="application.deletedAt"
        :label="$t('deletedAt')"
        label-cols="2"
      >
        <b-form-input
          v-model="application.deletedAt"
          data-test-id="input-deleted-at"
          plaintext
          disabled
        />
      </b-form-group>

      <b-form-group
        v-if="application.createdAt"
        :label="$t('createdAt')"
        label-cols="2"
        class="mb-0"
      >
        <b-form-input
          v-model="application.createdAt"
          data-test-id="input-created-at"
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
        @submit="$emit('submit', application)"
      />

      <confirmation-toggle
        v-if="application && application.applicationID"
        :data-test-id="deleteButtonStatusCypressId"
        @confirmed="$emit('delete')"
      >
        {{ getDeleteStatus }}
      </confirmation-toggle>
    </template>
  </b-card>
</template>

<script>
import { NoID } from '@cortezaproject/corteza-js'
import ConfirmationToggle from 'corteza-webapp-admin/src/components/ConfirmationToggle'
import CSubmitButton from 'corteza-webapp-admin/src/components/CSubmitButton'

export default {
  name: 'CApplicationEditorInfo',

  i18nOptions: {
    namespaces: 'system.applications',
    keyPrefix: 'editor.info',
  },

  components: {
    ConfirmationToggle,
    CSubmitButton,
  },

  props: {
    application: {
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
    fresh () {
      return !this.application.applicationID || this.application.applicationID === NoID
    },

    editable () {
      return this.fresh ? this.canCreate : this.application.canUpdateApplication
    },

    saveDisabled () {
      return !this.editable || [this.nameState].includes(false)
    },

    nameState () {
      const { name } = this.application
      return name ? null : false
    },

    getDeleteStatus () {
      return this.application.deletedAt ? this.$t('undelete') : this.$t('delete')
    },

    deleteButtonStatusCypressId () {
      return `button-${this.getDeleteStatus.toLowerCase()}`
    },
  },
}
</script>
