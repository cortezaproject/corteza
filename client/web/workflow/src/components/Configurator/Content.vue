<template>
  <c-rich-text-input
    v-model="label"
    :labels="{
      urlPlaceholder: $t('content.urlPlaceholder'),
      ok: $t('content.ok'),
      openLinkInNewTab: $t('content.openLinkInNewTab'),
    }"
    class="m-3"
    @input="$emit('update-value', $event)"
  />
</template>

<script>
import base from './base'
import { components } from '@cortezaproject/corteza-vue'
const { CRichTextInput } = components

export default {
  components: {
    CRichTextInput,
  },

  extends: base,

  computed: {
    label: {
      get () {
        return this.item.node.value
      },

      set (label) {
        this.item.node.value = label
      },
    },
  },
}
</script>
