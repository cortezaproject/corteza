<template>
  <b-form-group>
    <c-rich-text-input
      v-if="options"
      v-model="options.value"
      :labels="{
        urlPlaceholder: $t('label.urlPlaceholder'),
        ok: $t('label.ok'),
        openLinkInNewTab: $t('label.openLinkInNewTab'),
      }"
    />
  </b-form-group>
</template>

<script>
import base from './base'
import { components } from '@cortezaproject/corteza-vue'
const { CRichTextInput } = components

export default {
  i18nOptions: {
    namespaces: 'general',
  },

  components: {
    CRichTextInput,
  },

  extends: base,
}
</script>
