<template>
  <div class="flex flex-col gap-4">
    <p v-if="message" class="m-0 whitespace-pre-wrap" v-html="message" />

    <label v-if="label" class="text-sm font-medium text-color">{{ label }}</label>

    <AutoComplete
      v-model="selected"
      :suggestions="options"
      option-label="label"
      class="w-full"
      dropdown
      :loading="processing"
      :placeholder="placeholder"
      @complete="search"
    />

    <div class="flex justify-end">
      <Button
        :disabled="loading"
        :label="pVal('buttonLabel', 'Submit')"
        @click="$emit('submit', { value: encodeValue() })"
      />
    </div>
  </div>
</template>

<script>
import base from './base.vue'

export default {
  name: 'CPromptComposeRecordPicker',
  extends: base,
  emits: ['submit'],
  data() {
    return {
      namespaceID: '',
      moduleID: '',
      processing: false,
      options: [],
      selected: null,
    }
  },
  computed: {
    placeholder() {
      return this.pVal('placeholder', 'Select a record')
    },
  },
  async created() {
    this.namespaceID = this.resolveNamespaceID(this.pVal('namespace'), this.pType('namespace'))
    this.moduleID = this.resolveModuleID(this.pVal('module'), this.pType('module'))
    await this.loadLatest()
  },
  methods: {
    resolveNamespaceID(namespace, type) {
      if (type === 'ComposeNamespace') return namespace?.namespaceID || ''
      return namespace || ''
    },
    resolveModuleID(module, type) {
      if (type === 'ComposeModule') return module?.moduleID || ''
      return module || ''
    },
    async loadLatest() {
      await this.fetchRecords('')
    },
    async search({ query = '' }) {
      await this.fetchRecords(query)
    },
    async fetchRecords(query = '') {
      if (!this.$ComposeAPI || !this.namespaceID || !this.moduleID) {
        return
      }

      this.processing = true
      try {
        const { set = [] } = await this.$ComposeAPI.recordList({
          namespaceID: this.namespaceID,
          moduleID: this.moduleID,
          limit: 10,
          query,
        })

        const labelField = this.pVal('labelField')
        this.options = set.map(record => ({
          record,
          label: record?.values?.[labelField] || record?.recordID,
        }))
      } finally {
        this.processing = false
      }
    },
    encodeValue() {
      if (!this.selected?.record) {
        return { '@type': 'Any', '@value': null }
      }

      return { '@type': 'ComposeRecord', '@value': this.selected.record }
    },
  },
}
</script>
