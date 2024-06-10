<template>
  <b-tab :title="$t('iframe.label')">
    <b-form-group
      v-if="enableFromRecordURL"
      :label="$t('iframe.srcFieldLabel')"
      :description="$t('iframe.srcFieldDesc')"
      :disabled="!fields.length"
      label-class="text-primary"
    >
      <b-select
        v-model="options.srcField"
        type="url"
        :options="fieldOptions"
      />
    </b-form-group>
    <b-form-group
      :label="$t('iframe.srcLabel')"
      :description="enableFromRecordURL ? $t('iframe.srcDesc') : ''"
      label-class="text-primary"
    >
      <c-input-expression
        v-model="options.src"
        auto-complete
        lang="text"
        :suggestion-params="recordAutoCompleteParams"
        height="2.375rem"
        type="url"
      />
      <i18next
        path="interpolationFootnote"
        tag="small"
        class="text-muted"
      >
        <code>${record.values.fieldName}</code>
        <code>${recordID}</code>
        <code>${ownerID}</code>
        <span><code>${userID}</code>, <code>${user.name}</code></span>
      </i18next>
    </b-form-group>
  </b-tab>
</template>
<script>
import base from './base'
import { components } from '@cortezaproject/corteza-vue'
import autocomplete from 'corteza-webapp-compose/src/mixins/autocomplete.js'

const { CInputExpression } = components

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  name: 'IFrame',

  components: { CInputExpression },

  extends: base,

  mixins: [autocomplete],

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
