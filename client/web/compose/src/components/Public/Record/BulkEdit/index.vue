<template>
  <div>
    <b-button
      :title="$t('recordList.bulkRecord.title')"
      variant="outline-light"
      class="text-primary border-0"
      size="sm"
      @click="showModal = true"
    >
      <font-awesome-icon
        :icon="['far', 'edit']"
      />
    </b-button>

    <b-modal
      :visible="showModal"
      :title="$t('recordList.bulkRecord.title')"
      body-class="p-0"
      footer-class="d-flex justify-content-between align-items-center"
      centered
      @hide="onModalHide"
    >
      <b-card class="pt-0">
        <field-editor
          v-for="(field, index) in fields"
          :key="index"
          :namespace="namespace"
          :module="module"
          :field="field"
          :errors="fieldErrors(field.name)"
          :record="record"
        />

        <hr
          v-if="fields.length"
          class="my-4"
        >

        <vue-select
          v-model="selectedField"
          :placeholder="$t('recordList.bulkRecord.searchFields')"
          :get-option-label="getFieldLabel"
          :options="moduleFields"
          append-to-body
          :calculate-position="calculatePosition"
          :selectable="option => !selectedFields.includes(option.name)"
          class="bg-white position-relative"
          @input="addField"
        />
      </b-card>

      <template #modal-footer>
        <b-button
          variant="light"
          :disabled="processing"
          @click="onReset"
        >
          {{ $t('recordList.bulkRecord.reset') }}
        </b-button>

        <div>
          <b-button
            variant="link"
            rounded
            class="text-decoration-none text-primary"
            @click="onModalHide"
          >
            {{ $t('general.label.cancel') }}
          </b-button>
          <b-button
            variant="primary"
            :disabled="!fields.length || processing"
            @click="handleBulkUpdateSelectedRecords(selectedRecords)"
          >
            {{ $t('general.label.save') }}
          </b-button>
        </div>
      </template>
    </b-modal>
  </div>
</template>

<script>
import FieldEditor from 'corteza-webapp-compose/src/components/ModuleFields/Editor'
import { VueSelect } from 'vue-select'
import { compose } from '@cortezaproject/corteza-js'
import calculatePosition from 'corteza-webapp-compose/src/mixins/vue-select-position'
import record from 'corteza-webapp-compose/src/mixins/record.js'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  name: 'BulkEdit',

  components: {
    VueSelect,
    FieldEditor,
  },

  mixins: [
    calculatePosition,
    record,
  ],

  props: {
    namespace: {
      type: compose.Namespace,
      required: true,
    },
    module: {
      type: compose.Module,
      required: true,
    },

    selectedRecords: {
      type: Array,
      required: true,
    },
  },

  data () {
    return {
      showModal: false,
      selectedField: undefined,
      fields: [],
    }
  },

  computed: {
    selectedFields () {
      return this.fields.map(({ name }) => name)
    },

    moduleFields () {
      return [
        ...[...this.module.fields].sort((a, b) =>
          (a.label || a.name).localeCompare(b.label || b.name),
        ),
        ...this.module.systemFields().filter(({ name }) => name === 'ownedBy'),
      ].filter((field) => this.isFieldEditable(field))
    },
  },

  created () {
    this.record = new compose.Record(this.module, {})
  },

  methods: {
    onModalHide () {
      this.showModal = false
    },

    getFieldLabel ({ kind, label, name }) {
      return label || name || kind
    },

    addField (field) {
      if (!field) return

      this.fields.push(field)
      this.selectedField = null
    },

    onReset () {
      this.record = new compose.Record(this.module, {})
      this.fields = []
    },

    isFieldEditable (field) {
      if (!field) return false

      const { canCreateOwnedRecord } = this.module || {}
      const { createdAt, canManageOwnerOnRecord } = this.record || {}
      const { name, canUpdateRecordValue, isSystem, expressions = {} } = field || {}

      if (!canUpdateRecordValue) return false

      if (isSystem) {
        // Make ownedBy field editable if correct permissions
        if (name === 'ownedBy') {
          // If not created we check module permissions, otherwise the canManageOwnerOnRecord
          return createdAt ? canManageOwnerOnRecord : canCreateOwnedRecord
        }

        return false
      }

      return !expressions.value
    },
  },
}
</script>

<style lang="scss">
.position-initial {
  position: initial;
}
</style>
