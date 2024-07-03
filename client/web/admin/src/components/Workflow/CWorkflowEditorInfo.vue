<template>
  <b-card
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
      @submit.prevent="$emit('submit', workflow)"
    >
      <b-row>
        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            v-if="workflow.meta"
            :label="$t('name')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="workflow.meta.name"
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
            :label="$t('handle')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="workflow.handle"
              :state="handleState"
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
            :label="$t('enabled')"
            label-class="text-primary"
            :class="{ 'mb-0': !workflow.workflowID }"
          >
            <c-input-checkbox
              v-model="workflow.enabled"
              switch
              :labels="checkboxLabel"
            />
          </b-form-group>
        </b-col>
      </b-row>

      <c-system-fields
        :id="workflow.workflowID"
        :resource="workflow"
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
        v-if="workflow && workflow.workflowID && workflow.canDeleteWorkflow"
        :disabled="deleteDisabled"
        variant="danger"
        size="md"
        @confirmed="$emit('delete')"
      >
        {{ getDeleteStatus }}
      </c-input-confirm>

      <b-button
        v-if="workflow.workflowID"
        variant="light"
        @click="openWorkflowBuilder()"
      >
        {{ $t('openBuilder') }}
      </b-button>

      <c-button-submit
        :disabled="saveDisabled"
        :processing="processing"
        :success="success"
        :text="$t('admin:general.label.submit')"
        class="ml-auto"
        @submit="$emit('submit', workflow)"
      />
    </template>
  </b-card>
</template>

<script>
import { handle } from '@cortezaproject/corteza-vue'
import { NoID } from '@cortezaproject/corteza-js'

export default {
  name: 'CWorkflowEditorInfo',

  i18nOptions: {
    namespaces: 'automation.workflows',
    keyPrefix: 'editor.info',
  },

  props: {
    workflow: {
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
    editable () {
      return (!this.workflow.workflowID || this.workflow.workflowID === NoID) || this.workflow.canUpdateWorkflow
    },

    nameState () {
      return this.workflow.meta.name ? null : false
    },

    handleState () {
      return handle.handleState(this.workflow.handle)
    },

    saveDisabled () {
      return !this.editable || [this.nameState, this.handleState].includes(false)
    },

    deleteDisabled () {
      return !(this.workflow.deletedAt ? this.workflow.canUndeleteWorkflow : this.workflow.canDeleteWorkflow)
    },

    getDeleteStatus () {
      return this.workflow.deletedAt ? this.$t('undelete') : this.$t('delete')
    },
  },

  methods: {
    openWorkflowBuilder () {
      window.open(`/workflow/${this.workflow.workflowID}/edit`, '_blank')
    },
  },
}
</script>
