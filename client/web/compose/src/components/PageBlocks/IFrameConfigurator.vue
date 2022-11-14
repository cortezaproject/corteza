<template>
  <b-tab :title="$t('iframe.label')">
    <b-form-group
      v-if="enableFromRecordURL"
      class="form-group"
      :label="$t('iframe.srcFieldLabel')"
      :description="$t('iframe.srcFieldDesc')"
      :disabled="!fields.length"
    >
      <b-select
        v-model="options.srcField"
        type="url"
        class="form-control"
        :options="fieldOptions"
      />
    </b-form-group>
    <b-form-group
      class="form-group"
      :label="$t('iframe.srcLabel')"
      :description="enableFromRecordURL ? $t('iframe.srcDesc') : ''"
    >
      <b-form-input
        v-model="options.src"
        class="form-control"
        type="url"
      />
    </b-form-group>
  </b-tab>
</template>
<script>
import base from './base'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  name: 'IFrame',

  extends: base,

  computed: {
    fields () {
      if (!this.module) {
        return []
      }

      return this.module.fields
        .filter(({ kind }) => kind === 'Url')
        .map(({ name, label }) => ({ value: name, text: label }))
        .sort((a, b) => a.text.localeCompare(b.text))
    },

    fieldOptions () {
      return [
        { value: '', text: this.$t('iframe.pickURLField') },
        ...this.fields,
      ]
    },

    enableFromRecordURL () {
      return this.page.moduleID !== '0'
    },
  },
}
</script>
