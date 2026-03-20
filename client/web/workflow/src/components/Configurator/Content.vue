<template>
  <c-rich-text-input
    v-model="label"
    :labels="{
      urlPlaceholder: $t('steps.content.configurator.urlPlaceholder'),
      ok: $t('steps.content.configurator.ok'),
    }"
    class="m-3"
    @input="$emit('update-value', $event)"
  />
</template>

<script>
import base from './base.vue'
import { components } from '@cortezaproject/corteza-vue-next'
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
