<template>
  <b-card
    class="shadow-sm"
    header-bg-variant="white"
    footer-bg-variant="white"
  >
    <b-form
      @submit.prevent="$emit('submit', workflow)"
    >
      <b-form-group
        v-if="workflow.workflowID"
        :label="$t('id')"
        label-cols="2"
      >
        <b-form-input
          v-model="workflow.workflowID"
          plaintext
          disabled
        />
      </b-form-group>

      <b-form-group
        v-if="workflow.meta"
        :label="$t('name')"
        label-cols="2"
      >
        <b-form-input
          v-model="workflow.meta.name"
          required
        />
      </b-form-group>

      <b-form-group
        :label="$t('handle')"
        label-cols="2"
      >
        <b-form-input
          v-model="workflow.handle"
          :state="handleState"
        />
        <b-form-invalid-feedback :state="handleState">
          {{ $t('invalid-handle-characters') }}
        </b-form-invalid-feedback>
      </b-form-group>

      <b-form-group
        label-cols="2"
        :class="{ 'mb-0': !workflow.workflowID }"
      >
        <b-form-checkbox
          v-model="workflow.enabled"
        >
          {{ $t('enabled') }}
        </b-form-checkbox>
      </b-form-group>

      <b-form-group
        v-if="workflow.workflowID"
        label-cols="2"
      >
        <b-button
          variant="light"
          class="align-top"
          @click="openWorkflowBuilder()"
        >
          {{ $t('openBuilder') }}
        </b-button>
      </b-form-group>

      <b-form-group
        v-if="workflow.updatedAt"
        :label="$t('updatedAt')"
        label-cols="2"
      >
        <b-form-input
          v-model="workflow.updatedAt"
          plaintext
          disabled
        />
      </b-form-group>

      <b-form-group
        v-if="workflow.deletedAt"
        :label="$t('deletedAt')"
        label-cols="2"
      >
        <b-form-input
          v-model="workflow.deletedAt"
          plaintext
          disabled
        />
      </b-form-group>

      <b-form-group
        v-if="workflow.createdAt"
        :label="$t('createdAt')"
        label-cols="2"
        class="mb-0"
      >
        <b-form-input
          v-model="workflow.createdAt"
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
        @submit="$emit('submit', workflow)"
      />

      <confirmation-toggle
        v-if="workflow && workflow.workflowID"
        :disabled="deleteDisabled"
        @confirmed="$emit('delete')"
      >
        {{ getDeleteStatus }}
      </confirmation-toggle>
    </template>
  </b-card>
</template>

<script>
import { handleState } from 'corteza-webapp-admin/src/lib/handle'
import ConfirmationToggle from 'corteza-webapp-admin/src/components/ConfirmationToggle'
import CSubmitButton from 'corteza-webapp-admin/src/components/CSubmitButton'

export default {
  name: 'CWorkflowEditorInfo',

  i18nOptions: {
    namespaces: 'automation.workflows',
    keyPrefix: 'editor.info',
  },

  components: {
    ConfirmationToggle,
    CSubmitButton,
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

  computed: {
    editable () {
      return this.workflow.canUpdateWorkflow
    },

    handleState () {
      const { handle } = this.workflow

      return handle ? handleState(handle) : false
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
      window.open(`${window.location.origin}/workflow/${this.workflow.workflowID}/edit`, '_blank')
    },
  },
}
</script>
