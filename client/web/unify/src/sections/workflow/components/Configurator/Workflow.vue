<template>
  <div v-if="localWorkflow" class="flex flex-col gap-4">
    <div v-if="workflow.workflowID && workflow.workflowID !== '0'" class="flex gap-2 px-3 pt-3">
      <import
        data-test-id="button-import-workflow"
        :disabled="importProcessing"
        @import="$emit('import', $event)"
      />

      <export
        data-test-id="button-export-workflow"
        :workflows="[workflow.workflowID]"
        :file-name="workflow.meta.name || workflow.handle"
        size="small"
        severity="secondary"
      />

      <CPermissionsButton
        v-if="workflow.canGrant"
        :resource="`corteza::automation:workflow/${workflow.workflowID}`"
        :title="workflow.meta.name || workflow.handle || workflow.workflowID"
        :target="workflow.meta.name || workflow.handle || workflow.workflowID"
        icon="pi pi-lock"
        size="small"
        outlined
      />
    </div>

    <div class="flex flex-col gap-4 px-3">
      <CFormGroup :label="$t('configurator.name.label')" required>
        <InputText
          v-model="localWorkflow.meta.name"
          data-test-id="input-label"
          :placeholder="$t('configurator.name.placeholder')"
          :invalid="nameState === false"
        />
      </CFormGroup>

      <CFormGroup :label="$t('configurator.handle.label')">
        <InputText
          v-model="localWorkflow.handle"
          data-test-id="input-handle"
          :invalid="handleState === false"
          :placeholder="$t('configurator.handle.placeholder')"
        />
        <small
          v-if="handleState === false"
          class="text-red-500"
          data-test-id="input-handle-invalid-state"
        >
          {{ $t('configurator.handle.invalid-handle-characters') }}
        </small>
      </CFormGroup>

      <CFormGroup :label="$t('configurator.description.label')">
        <Textarea
          v-model="localWorkflow.meta.description"
          data-test-id="input-description"
          :placeholder="$t('configurator.description.placeholder')"
          rows="3"
          autoResize
        />
      </CFormGroup>

      <CFormGroup :label="$t('configurator.labels.label')">
        <namespace-module-selector
          :namespace-labels="localWorkflow?.labels?.ref_namespace || []"
          :module-labels="localWorkflow?.labels?.ref_module || []"
          @change="handleLabelsChange"
        />
      </CFormGroup>

      <Divider class="!my-1" />

      <CFormGroup :label="$t('configurator.run-as.label')">
        <c-input-user
          v-model="localWorkflow.runAs"
          data-test-id="select-run-as"
          :placeholder="$t('configurator.run-as.placeholder')"
        />
        <small class="text-muted-color">
          {{ $t('configurator.run-as.description') }}
        </small>
      </CFormGroup>

      <div class="flex items-center gap-2">
        <Checkbox
          v-model="localWorkflow.enabled"
          :binary="true"
          inputId="workflow-enabled"
          data-test-id="checkbox-enable-workflow"
        />
        <label for="workflow-enabled">
          {{ $t('general.enabled') }}
        </label>
      </div>

      <div class="flex flex-col gap-1">
        <div class="flex items-center gap-2">
          <Checkbox
            v-model="localWorkflow.meta.subWorkflow"
            :binary="true"
            inputId="workflow-sub"
            data-test-id="checkbox-sub-workflow"
          />
          <label for="workflow-sub">
            {{ $t('configurator.sub-workflow.label') }}
          </label>
        </div>
        <small class="text-muted-color ml-7">
          {{ $t('configurator.sub-workflow.description') }}
        </small>
      </div>

      <Divider class="!my-1" />

      <!-- Named input/output contract, read by the "Run Workflow" step -->
      <CFormGroup
        v-for="section in ioSections"
        :key="section.key"
        :label="section.label"
      >
        <small class="text-muted-color block mb-2">{{ section.description }}</small>

        <div class="flex items-center mb-2">
          <Button
            :label="section.addLabel"
            icon="pi pi-plus"
            severity="secondary"
            size="small"
            @click="addIODef(section.key)"
          />
        </div>

        <CFormList
          v-model="localWorkflow.meta[section.key]"
          :columns="ioColumns"
          :empty-message="section.emptyMessage"
          :confirm-remove="section.removeConfirm"
        >
          <template #row="{ item }">
            <InputText
              v-model="item.name"
              class="w-full"
              size="small"
              :placeholder="$t('configurator.io.columns.name')"
            />
            <InputText
              v-model="item.label"
              class="w-full"
              size="small"
              :placeholder="$t('configurator.io.columns.label')"
            />
            <Select
              v-model="item.types[0]"
              :options="ioTypeOptions"
              option-label="label"
              option-value="value"
              class="w-full"
              size="small"
            />
            <div class="flex justify-center">
              <Checkbox v-model="item.required" :binary="true" />
            </div>
          </template>
        </CFormList>
      </CFormGroup>
    </div>

    <div class="flex items-center w-full p-3 mt-auto border-t surface-border">
      <Button
        v-if="workflow.canDeleteWorkflow && !isDeleted"
        :label="$t('editor.delete')"
        severity="danger"
        text
        size="small"
        :loading="processingDelete"
        @click="handleDeleteClick"
      />
      <Button
        v-else-if="isDeleted"
        :label="$t('editor.undelete')"
        severity="secondary"
        text
        size="small"
        :loading="processingDelete"
        @click="$emit('undelete')"
      />

      <div class="flex-1" />

      <div class="flex items-center gap-2">
        <Button
          v-if="workflow.workflowID === '0'"
          :label="$t('editor.back')"
          severity="secondary"
          text
          size="small"
          @click="$router.back()"
        />
        <Button
          v-if="workflow.workflowID !== '0'"
          :label="$t('general.cancel')"
          severity="secondary"
          text
          size="small"
          @click="handleCancel"
        />

        <Button
          data-test-id="button-save-workflow"
          :label="$t('editor.save')"
          :disabled="isSaveDisabled"
          :loading="processingSave"
          size="small"
          @click="handleSave"
        />
      </div>
    </div>
  </div>
</template>

<script>
import { components, useConfirmDelete } from '@planetcrust/human-vue'
import { automation } from '@planetcrust/human-js'
import Import from '../Import.vue'
import Export from '../Export.vue'
import NamespaceModuleSelector from '../NamespaceModuleSelector.vue'

const { CInputUser } = components

const handleRe = /^[A-Za-z][0-9A-Za-z_\-.]*[A-Za-z0-9]$/

// Types a workflow input/output field can declare (label resolved via $t)
const IO_TYPES = [
  { key: 'text', value: 'String' },
  { key: 'number', value: 'Number' },
  { key: 'boolean', value: 'Boolean' },
  { key: 'dateTime', value: 'DateTime' },
  { key: 'any', value: 'Any' },
]

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

  setup() {
    const { confirmDelete } = useConfirmDelete()
    return { confirmDelete }
  },

  data() {
    return {
      localWorkflow: null,
    }
  },

  computed: {
    ioTypeOptions() {
      return IO_TYPES.map(t => ({ label: this.$t(`configurator.io.types.${t.key}`), value: t.value }))
    },

    ioColumns() {
      return [
        { label: this.$t('configurator.io.columns.name'), width: 'minmax(140px, 1.2fr)' },
        { label: this.$t('configurator.io.columns.label'), width: 'minmax(140px, 1.2fr)' },
        { label: this.$t('configurator.io.columns.type'), width: 'minmax(140px, 1fr)' },
        { label: this.$t('configurator.io.columns.required'), width: '90px', headerClass: 'text-center' },
      ]
    },

    ioSections() {
      return [
        {
          key: 'input',
          label: this.$t('configurator.io.inputs.label'),
          description: this.$t('configurator.io.inputs.description'),
          addLabel: this.$t('configurator.io.inputs.addField'),
          emptyMessage: this.$t('configurator.io.inputs.empty'),
          removeConfirm: this.$t('configurator.io.inputs.removeConfirm'),
        },
        {
          key: 'output',
          label: this.$t('configurator.io.outputs.label'),
          description: this.$t('configurator.io.outputs.description'),
          addLabel: this.$t('configurator.io.outputs.addField'),
          emptyMessage: this.$t('configurator.io.outputs.empty'),
          removeConfirm: this.$t('configurator.io.outputs.removeConfirm'),
        },
      ]
    },

    nameState() {
      return this.localWorkflow?.meta?.name ? null : false
    },

    handleState() {
      if (!this.localWorkflow) return null
      const h = this.localWorkflow.handle
      if (!h) return null
      return handleRe.test(h) ? null : false
    },

    canUpdateWorkflow() {
      return this.workflow.workflowID === '0' ? this.canCreate : this.workflow.canUpdateWorkflow
    },

    isSaveDisabled() {
      return !this.canUpdateWorkflow || [this.nameState, this.handleState].includes(false)
    },

    isDeleted() {
      return this.workflow.deletedAt
    },
  },

  watch: {
    workflow: {
      handler(newWorkflow) {
        if (!newWorkflow) {
          this.localWorkflow = null
          return
        }

        // Create a new Workflow instance from the existing workflow
        // This properly clones the workflow and avoids circular references
        this.localWorkflow = new automation.Workflow(newWorkflow)

        // Normalize the input/output declaration lists for direct editor binding
        for (const key of ['input', 'output']) {
          if (!Array.isArray(this.localWorkflow.meta[key])) {
            this.localWorkflow.meta[key] = []
          }
          for (const row of this.localWorkflow.meta[key]) {
            if (!Array.isArray(row.types) || row.types.length === 0) {
              row.types = ['String']
            }
          }
        }
      },
      immediate: true,
    },
  },

  methods: {
    addIODef(key) {
      this.localWorkflow.meta[key].push({ name: '', label: '', types: ['String'], required: false })
    },

    handleLabelsChange({ namespaceLabels, moduleLabels }) {
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

    handleSave() {
      // Emit save event with the local workflow copy
      this.$emit('save', this.localWorkflow)
      // Close the modal after save
      this.$emit('close')
    },

    handleCancel() {
      // Reset local workflow to original workflow data
      this.localWorkflow = new automation.Workflow(this.workflow)
      // Close the modal
      this.$emit('close')
    },

    handleDeleteClick() {
      this.confirmDelete({
        message: this.$t('editor.delete-confirm'),
        header: this.localWorkflow?.meta?.name || this.localWorkflow?.handle || '',
        onConfirm: () => this.$emit('delete'),
      })
    },
  },
}
</script>
