<template>
  <b-card
    data-test-id="card-application-info"
    header-class="border-bottom"
    footer-class="border-top d-flex flex-wrap flex-fill-child gap-1"
    class="shadow-sm"
  >
    <b-form
      @submit.prevent="$emit('submit', application)"
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
              v-model="application.name"
              data-test-id="input-name"
              :state="nameState"
              required
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('enabled')"
            :class="{ 'mb-0': !application.applicationID }"
            label-class="text-primary"
          >
            <c-input-checkbox
              v-model="application.enabled"
              data-test-id="checkbox-enabled"
              :labels="checkboxLabel"
              switch
            />
          </b-form-group>
        </b-col>
      </b-row>

      <c-system-fields
        :id="application.applicationID"
        :resource="application"
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
        v-if="application && application.applicationID && application.canDeleteApplication"
        :data-test-id="deleteButtonStatusCypressId"
        variant="danger"
        size="md"
        @confirmed="$emit('delete')"
      >
        {{ getDeleteStatus }}
      </c-input-confirm>

      <c-button-submit
        :disabled="saveDisabled"
        :processing="processing"
        :success="success"
        :text="$t('admin:general.label.submit')"
        class="ml-auto"
        @submit="$emit('submit', application)"
      />
    </template>
  </b-card>
</template>

<script>
import { NoID } from '@cortezaproject/corteza-js'

export default {
  name: 'CApplicationEditorInfo',

  i18nOptions: {
    namespaces: 'system.applications',
    keyPrefix: 'editor.info',
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

  data () {
    return {
      checkboxLabel: {
        on: this.$t('general:label.general.yes'),
        off: this.$t('general:label.general.no'),
      },
    }
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
