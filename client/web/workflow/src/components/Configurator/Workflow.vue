<template>
  <b-modal
    id="workflow"
    :title="$t('editor:workflow-configuration')"
    size="lg"
    :hide-header-close="workflow.workflowID === '0'"
    :no-close-on-backdrop="workflow.workflowID === '0'"
    :no-close-on-esc="workflow.workflowID === '0'"
    no-fade
    body-class="p-0"
  >
    <template #modal-title>
      {{ $t('editor:workflow-configuration') }}
    </template>

    <div
      v-if="workflow.workflowID && workflow.workflowID !== '0'"
      class="d-flex p-3"
    >
      <import
        data-test-id="button-import-workflow"
        :disabled="importProcessing"
        @import="$emit('import', $event)"
      />

      <export
        data-test-id="button-export-workflow"
        :workflows="[workflow.workflowID]"
        :file-name="workflow.meta.name || workflow.handle"
        size="lg"
        class="ml-1"
      />

      <c-permissions-button
        v-if="workflow.canGrant"
        :title="workflow.meta.name || workflow.handle || workflow.workflowID"
        :target="workflow.meta.name || workflow.handle || workflow.workflowID"
        :resource="`corteza::automation:workflow/${workflow.workflowID}`"
        :button-label="$t('general:permissions')"
        class="btn-lg ml-1"
      />
    </div>

    <div v-if="localWorkflow">
      <b-tabs card>
        <b-tab
          :title="$t('general.label')"
          active
        >
          <b-form-group
            :label="$t('name.label')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="localWorkflow.meta.name"
              data-test-id="input-label"
              :placeholder="$t('name.placeholder')"
              :state="nameState"
            />
          </b-form-group>

          <b-form-group
            :label="$t('handle.label')"
            label-class="text-primary"
          >
            <b-form-input
              v-model="localWorkflow.handle"
              data-test-id="input-handle"
              :state="handleState"
              :placeholder="$t('handle.placeholder')"
            />
            <b-form-invalid-feedback
              data-test-id="input-handle-invalid-state"
              :state="handleState"
            >
              {{ $t('handle.invalid-handle-characters') }}
            </b-form-invalid-feedback>
          </b-form-group>

          <b-form-group
            :label="$t('description.label')"
            label-class="text-primary"
          >
            <b-form-textarea
              v-model="localWorkflow.meta.description"
              data-test-id="input-description"
              :placeholder="$t('description.placeholder')"
            />
          </b-form-group>

          <b-form-group
            :label="$t('run-as.label')"
            :description="$t('run-as.description')"
            label-class="text-primary"
          >
            <c-input-user
              v-model="localWorkflow.runAs"
              data-test-id="select-run-as"
              :placeholder="$t('run-as.placeholder')"
              clearable
            />
          </b-form-group>

          <b-form-group>
            <b-form-checkbox
              v-model="localWorkflow.enabled"
              data-test-id="checkbox-enable-workflow"
            >
              {{ $t('general:enabled') }}
            </b-form-checkbox>
          </b-form-group>

          <b-form-group
            :description="$t('sub-workflow.description')"
          >
            <b-form-checkbox
              v-model="localWorkflow.meta.subWorkflow"
              data-test-id="checkbox-sub-workflow"
            >
              {{ $t('sub-workflow.label') }}
            </b-form-checkbox>
          </b-form-group>
        </b-tab>

        <b-tab :title="$t('labels.label')">
          <namespace-module-selector
            :namespace-labels="localWorkflow?.labels?.ref_namespace || []"
            :module-labels="localWorkflow?.labels?.ref_module || []"
            @change="handleLabelsChange"
          />
        </b-tab>
      </b-tabs>
    </div>

    <template #modal-footer>
      <div class="d-flex w-100">
        <c-input-confirm
          v-if="workflow.canDeleteWorkflow && !isDeleted"
          size="md"
          size-confirm="md"
          variant="danger"
          :processing="processingDelete"
          :text="$t('editor:delete')"
          :borderless="false"
          @confirmed="$emit('delete')"
        />

        <c-input-confirm
          v-else-if="isDeleted"
          size="md"
          size-confirm="md"
          variant="light"
          :processing="processingDelete"
          :text="$t('editor:undelete')"
          :borderless="false"
          @confirmed="$emit('undelete')"
        />

        <b-button
          v-if="workflow.workflowID === '0'"
          variant="light"
          @click="$router.back()"
        >
          {{ $t('editor:back') }}
        </b-button>

        <div class="d-flex ml-auto">
          <b-button
            v-if="workflow.workflowID !== '0'"
            variant="light"
            class="ml-auto"
            @click="handleCancel"
          >
            {{ $t('general:cancel') }}
          </b-button>

          <c-button-submit
            data-test-id="button-save-workflow"
            :disabled="isSaveDisabled"
            :processing="processingSave"
            :text="$t('editor:save')"
            class="ml-1"
            @submit="handleSave"
          />
        </div>
      </div>
    </template>
  </b-modal>
</template>

<script>
import { handle, components } from '@cortezaproject/corteza-vue'
import { automation } from '@cortezaproject/corteza-js'
import Import from '../Import'
import Export from '../Export'
import NamespaceModuleSelector from '../NamespaceModuleSelector'

const { CInputUser } = components

export default {
  i18nOptions: {
    namespaces: 'configurator',
  },

  components: {
    CInputUser,
    Import,
    Export,
    NamespaceModuleSelector,
  },

  props: {
    workflow: {
      type: Object,
      default: () => {},
    },

    canCreate: {
      type: Boolean,
      default: false,
    },

    processingSave: {
      type: Boolean,
      default: false,
    },

    processingDelete: {
      type: Boolean,
      default: false,
    },

    importProcessing: {
      type: Boolean,
      default: false,
    },
  },

  data () {
    return {
      localWorkflow: null,
    }
  },

  computed: {
    nameState () {
      return this.localWorkflow?.meta?.name ? null : false
    },

    handleState () {
      return this.localWorkflow ? handle.handleState(this.localWorkflow.handle) : null
    },

    canUpdateWorkflow () {
      return this.workflow.workflowID === '0' ? this.canCreate : this.workflow.canUpdateWorkflow
    },

    isSaveDisabled () {
      return !this.canUpdateWorkflow || [this.nameState, this.handleState].includes(false)
    },

    isDeleted () {
      return this.workflow.deletedAt
    },

  },

  watch: {
    workflow: {
      handler (newWorkflow) {
        if (!newWorkflow) {
          this.localWorkflow = null
          return
        }

        // Create a new Workflow instance from the existing workflow
        // This properly clones the workflow and avoids circular references
        this.localWorkflow = new automation.Workflow(newWorkflow)
      },
      immediate: true,
    },
  },

  methods: {
    handleLabelsChange ({ namespaceLabels, moduleLabels }) {
      if (!this.localWorkflow.labels) {
        this.localWorkflow.labels = {}
      }

      // Store labels directly - no transformation needed
      if (namespaceLabels.length > 0) {
        this.localWorkflow.labels.ref_namespace = namespaceLabels
      } else {
        delete this.localWorkflow.labels.ref_namespace
      }

      if (moduleLabels.length > 0) {
        this.localWorkflow.labels.ref_module = moduleLabels
      } else {
        delete this.localWorkflow.labels.ref_module
      }
    },

    handleSave () {
      // Emit save event with the local workflow copy
      this.$emit('save', this.localWorkflow)
      // Close the modal after save
      this.$bvModal.hide('workflow')
    },

    handleCancel () {
      // Reset local workflow to original workflow data
      this.localWorkflow = new automation.Workflow(this.workflow)
      // Close the modal
      this.$bvModal.hide('workflow')
    },
  },
}
</script>
