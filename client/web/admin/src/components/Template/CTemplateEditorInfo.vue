<template>
  <b-card
    data-test-id="card-template-info"
    header-class="border-bottom"
    footer-class="border-top d-flex flex-wrap flex-fill-child gap-1"
    class="shadow-sm"
  >
    <template #header>
      <h4 class="m-0">
        {{ $t('title') }}
      </h4>
    </template>

    <b-form
      @submit="$emit('submit', template)"
    >
      <b-row>
        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('meta.short')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="template.meta.short"
              data-test-id="input-short-name"
              required
              :state="shortState"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('handle')"
            label-class="text-primary"
            :class="{ 'mb-0': !template.templateID }"
          >
            <b-form-input
              v-model="template.handle"
              data-test-id="input-handle"
              :state="handleState"
              :placeholder="$t('placeholder-handle')"
            />

            <b-form-invalid-feedback :state="handleState">
              {{ $t('invalid-handle-characters') }}
            </b-form-invalid-feedback>
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('type')"
            label-class="text-primary"
          >
            <b-form-select
              v-model="template.type"
              data-test-id="select-template-type"
              :options="contentTypes"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('meta.description')"
            label-class="text-primary"
          >
            <b-form-textarea
              v-model="template.meta.description"
              data-test-id="textarea-description"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('partial')"
            :description="$t('partialDescription')"
            label-class="text-primary"
          >
            <c-input-checkbox
              v-model="template.partial"
              data-test-id="checkbox-is-partial-template"
              switch
              :labels="checkboxLabel"
              name="checkbox-1"
            />
          </b-form-group>
        </b-col>
      </b-row>

      <c-system-fields
        :resource="template"
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

    <template #footer>
      <c-input-confirm
        v-if="!fresh && template.canDeleteTemplate"
        :data-test-id="getDeletedButtonStatusCypressId"
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
        @submit="$emit('submit', template)"
      />
    </template>
  </b-card>
</template>

<script>
import { NoID } from '@cortezaproject/corteza-js'
import { handle } from '@cortezaproject/corteza-vue'

export default {
  name: 'CTemplateEditorInfo',

  i18nOptions: {
    namespaces: 'system.templates',
    keyPrefix: 'editor.info',
  },

  props: {
    template: {
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
      contentTypes: [
        { value: 'text/html', text: this.$t('contentType.text_html') },
        { value: 'text/plain', text: this.$t('contentType.text_plain') },
      ],
      checkboxLabel: {
        on: this.$t('general:label.general.yes'),
        off: this.$t('general:label.general.no'),
      },
    }
  },

  computed: {
    fresh () {
      return !this.template.templateID || this.template.templateID === NoID
    },

    editable () {
      return this.fresh ? this.canCreate : true // this.template.canUpdateRole
    },

    shortState () {
      return this.template.meta.short ? null : false
    },

    handleState () {
      return handle.handleState(this.template.handle)
    },

    saveDisabled () {
      return !this.editable || [this.shortState, this.handleState].includes(false)
    },

    getDeleteStatus () {
      return this.template.deletedAt ? this.$t('undelete') : this.$t('delete')
    },

    getArchiveStatus () {
      return this.template.archivedAt ? this.$t('unarchive') : this.$t('archive')
    },

    getDeletedButtonStatusCypressId () {
      return `button-${this.getDeleteStatus.toLowerCase()}`
    },
  },
}
</script>
